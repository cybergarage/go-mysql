// Copyright (C) 2019 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	"io"

	"github.com/cybergarage/go-mysql/mysql/encoding/binary"
	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-mysql/mysql/stmt"
)

// PacketReader represents a packet reader of MySQL protocol.
type PacketReader struct {
	*binary.Reader
}

// NewPacketReaderWithReader returns a new packet reader with the specified reader.
func NewPacketReaderWithReader(reader io.Reader) *PacketReader {
	return &PacketReader{
		Reader: binary.NewReaderWithReader(reader),
	}
}

// NewPacketReaderWithBytes returns a new packet reader with the specified bytes.
func NewPacketReaderWithBytes(data []byte) *PacketReader {
	return &PacketReader{
		Reader: binary.NewReaderWithBytes(data),
	}
}

// ReadCapabilitys reads the capability flags.
func (reader *PacketReader) ReadCapability() (Capability, error) {
	var Capabilitys Capability
	v, err := reader.ReadInt2()
	if err != nil {
		return 0, err
	}
	Capabilitys = Capability(v)

	if !Capabilitys.HasCapability(ClientProtocol41) {
		return Capabilitys, nil
	}

	Capabilitys3, err := reader.ReadInt1()
	if err != nil {
		return 0, err
	}
	Capabilitys |= (Capability)(Capabilitys3) << 16

	Capabilitys4, err := reader.ReadInt1()
	if err != nil {
		return 0, err
	}
	Capabilitys |= (Capability)(Capabilitys4) << 24

	return Capabilitys, nil
}

// PeekCapabilitys reads the capability flags.
func (reader *PacketReader) PeekCapability() (Capability, error) {
	var Capabilitys Capability
	v, err := reader.PeekInt2()
	if err != nil {
		return 0, err
	}
	Capabilitys = Capability(v)

	if !Capabilitys.HasCapability(ClientProtocol41) {
		return Capabilitys, nil
	}

	Capabilitys3, err := reader.PeekInt1()
	if err != nil {
		return 0, err
	}
	Capabilitys |= (Capability)(Capabilitys3) << 16

	Capabilitys4, err := reader.PeekInt1()
	if err != nil {
		return 0, err
	}
	Capabilitys |= (Capability)(Capabilitys4) << 24

	return Capabilitys, nil
}

// ReadFieldBytes reads the field bytes.
func (reader *PacketReader) ReadDatetimeBytes() ([]byte, error) {
	v1, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	v2, err := reader.ReadFixedLengthBytes(int(v1))
	if err != nil {
		return nil, err
	}
	return append([]byte{v1}, v2...), nil
}

// ReadFieldBytes reads the field bytes.
func (reader *PacketReader) ReadFieldBytes(t FieldType) ([]byte, error) {
	switch t {
	case query.MySQLTypeString, query.MySQLTypeVarString, query.MySQLTypeVarchar:
		v, err := reader.ReadLengthEncodedString()
		if err != nil {
			return nil, err
		}
		return []byte(v), nil
	case query.MySQLTypeTinyBlob, query.MySQLTypeMediumBlob, query.MySQLTypeLongBlob, query.MySQLTypeBlob:
		return reader.ReadLengthEncodedBytes()
	case query.MySQLTypeNull:
		return nil, nil
	case query.MySQLTypeTiny:
		return reader.ReadNBytes(1)
	case query.MySQLTypeShort, query.MySQLTypeYear:
		return reader.ReadNBytes(2)
	case query.MySQLTypeLong, query.MySQLTypeFloat, query.MySQLTypeInt24:
		return reader.ReadNBytes(4)
	case query.MySQLTypeLongLong, query.MySQLTypeDouble:
		return reader.ReadNBytes(8)
	case query.MySQLTypeDate, query.MySQLTypeTime, query.MySQLTypeDatetime, query.MySQLTypeTimestamp:
		return reader.ReadDatetimeBytes()
	}

	return nil, newErrFieldTypeNotSupported(t)
}

// ReadFieldValue reads the field value.
func (reader *PacketReader) ReadFieldValue(t FieldType) (any, error) {
	b, err := reader.ReadFieldBytes(t)
	if err != nil {
		return nil, err
	}

	field := stmt.NewField(
		stmt.WithFieldType(t),
		stmt.WithFieldBytes(b),
	)

	return field.Value()
}
