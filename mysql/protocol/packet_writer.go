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
	"time"

	"github.com/cybergarage/go-mysql/mysql/encoding/binary"
	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-safecast/safecast"
)

// PacketWriter represents a packet writer of MySQL protocol.
type PacketWriter struct {
	*binary.Writer
}

// NewPacketWriter returns a new packet writer.
func NewPacketWriter() *PacketWriter {
	return &PacketWriter{
		Writer: binary.NewWriter(),
	}
}

// WriteCommandType writes a command type.
func (w *PacketWriter) WriteCommandType(cmd Command) error {
	return w.WriteByte(byte(cmd.Type()))
}

// WriteCapabilitys writes the capability flags.
func (w *PacketWriter) WriteCapability(c Capability) error {
	if c.IsEnabled(ClientProtocol41) {
		return w.WriteInt4(uint32(c))
	}
	return w.WriteInt2(uint16(c >> 16))
}

// WriteFillerBytes writes the filler bytes.
func (w *PacketWriter) WriteFillerBytes(b byte, n int) error {
	for i := 0; i < n; i++ {
		if err := w.WriteByte(b); err != nil {
			return err
		}
	}
	return nil
}

// WritePacket writes a packet.
func (w *PacketWriter) WritePacket(pkt Response) error {
	pktBytes, err := pkt.Bytes()
	if err != nil {
		return err
	}
	_, err = w.WriteBytes(pktBytes)
	if err != nil {
		return err
	}
	return nil
}

// WriteOK writes a OK packet.
func (w *PacketWriter) WriteOK(opts ...any) error {
	okOpts := []OKOption{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case SequenceID:
			okOpts = append(okOpts, WithOKSecuenceID(v))
		case Capability:
			okOpts = append(okOpts, WithOKCapability(v))
		case ServerStatus:
			okOpts = append(okOpts, WithOKServerStatus(v))
		}
	}
	ok, err := NewOK(okOpts...)
	if err != nil {
		return err
	}
	okBytes, err := ok.Bytes()
	if err != nil {
		return err
	}
	_, err = w.WriteBytes(okBytes)
	if err != nil {
		return err
	}
	return nil
}

// WriteErr writes a ERR packet.
func (w *PacketWriter) WriteErr(opts ...any) error {
	errOpts := []ERROption{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case SequenceID:
			errOpts = append(errOpts, WithERRSecuenceID(v))
		}
	}
	pkt, err := NewERR(errOpts...)
	if err != nil {
		return err
	}
	errBytes, err := pkt.Bytes()
	if err != nil {
		return err
	}
	_, err = w.WriteBytes(errBytes)
	if err != nil {
		return err
	}
	return nil
}

// WriteEOF writes a EOF packet.
func (w *PacketWriter) WriteEOF(opts ...any) error {
	eofOpts := []EOFOption{}
	for _, opt := range opts {
		switch v := opt.(type) {
		case SequenceID:
			eofOpts = append(eofOpts, WithEOFCSecuenceID(v))
		case Capability:
			eofOpts = append(eofOpts, WithEOFCapability(v))
		case ServerStatus:
			eofOpts = append(eofOpts, WithEOFServerStatus(v))
		}
	}
	eof, err := NewEOF(eofOpts...)
	if err != nil {
		return err
	}
	eofBytes, err := eof.Bytes()
	if err != nil {
		return err
	}
	_, err = w.WriteBytes(eofBytes)
	if err != nil {
		return err
	}
	return nil
}

// WriteFieldBytes writes a field bytes.
func (w *PacketWriter) WriteFieldBytes(t FieldType, v any) error {
	// MySQL: Binary Protocol Resultset
	// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
	// Result Set Packets - MariaDB Knowledge Base
	// https://mariadb.com/kb/en/result-set-packets/

	switch t {
	case query.MySQLTypeString, query.MySQLTypeVarString, query.MySQLTypeVarchar:
		if s, ok := v.(string); ok {
			return w.WriteLengthEncodedString(s)
		}
	case query.MySQLTypeTinyBlob, query.MySQLTypeMediumBlob, query.MySQLTypeLongBlob, query.MySQLTypeBlob:
		if b, ok := v.([]byte); ok {
			return w.WriteLengthEncodedBytes(b)
		}
	case query.MySQLTypeNull:
		return nil
	case query.MySQLTypeTiny:
		var cv int8
		if err := safecast.ToInt8(v, &cv); err != nil {
			_, err := w.WriteBytes(binary.Int1ToBytes(cv))
			return err
		}
	case query.MySQLTypeShort, query.MySQLTypeYear:
		var cv int16
		if err := safecast.ToInt16(v, &cv); err != nil {
			_, err := w.WriteBytes(binary.Int2ToBytes(cv))
			return err
		}
	case query.MySQLTypeLong, query.MySQLTypeInt24:
		var cv int32
		if err := safecast.ToInt32(v, &cv); err != nil {
			_, err := w.WriteBytes(binary.Int4ToBytes(cv))
			return err
		}
	case query.MySQLTypeLonglong:
		var cv int64
		if err := safecast.ToInt64(v, &cv); err != nil {
			_, err := w.WriteBytes(binary.Int8ToBytes(cv))
			return err
		}
	case query.MySQLTypeFloat:
		var cv float32
		if err := safecast.ToFloat32(v, &cv); err != nil {
			_, err := w.WriteBytes(binary.Float4ToBytes(cv))
			return err
		}
	case query.MySQLTypeDouble:
		var cv float64
		if err := safecast.ToFloat64(v, &cv); err != nil {
			_, err := w.WriteBytes(binary.Float8ToBytes(cv))
			return err
		}
	case query.MySQLTypeDatetime, query.MySQLTypeTimestamp:
		var cv time.Time
		if err := safecast.ToTime(v, &cv); err != nil {
			_, err := w.WriteBytes(binary.TimeToBytes(cv))
			return err
		}

	}

	return newInvalidFieldValue(t, v)
}
