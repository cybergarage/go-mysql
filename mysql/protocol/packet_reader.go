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

	if !Capabilitys.IsEnabled(ClientProtocol41) {
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

	if !Capabilitys.IsEnabled(ClientProtocol41) {
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
