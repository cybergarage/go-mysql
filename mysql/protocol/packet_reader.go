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
)

// PacketReader represents a packet reader of MySQL protocol.
type PacketReader struct {
	*Reader
}

// NewPacketReader returns a new packet reader.
func NewPacketReaderWith(reader io.Reader) *PacketReader {
	return &PacketReader{
		Reader: NewReaderWithReader(reader),
	}
}

// ReadCapabilityFlags reads the capability flags.
func (reader *PacketReader) ReadCapability() (CapabilityFlag, error) {
	var capabilityFlags CapabilityFlag
	v, err := reader.ReadInt2()
	if err != nil {
		return 0, err
	}
	capabilityFlags = CapabilityFlag(v)

	if !capabilityFlags.IsEnabled(ClientProtocol41) {
		return capabilityFlags, nil
	}

	capabilityFlags3, err := reader.ReadInt1()
	if err != nil {
		return 0, err
	}
	capabilityFlags |= (CapabilityFlag)(capabilityFlags3) << 16

	capabilityFlags4, err := reader.ReadInt1()
	if err != nil {
		return 0, err
	}
	capabilityFlags |= (CapabilityFlag)(capabilityFlags4) << 24

	return capabilityFlags, nil
}

// PeekCapabilityFlags reads the capability flags.
func (reader *PacketReader) PeekCapability() (CapabilityFlag, error) {
	var capabilityFlags CapabilityFlag
	v, err := reader.PeekInt2()
	if err != nil {
		return 0, err
	}
	capabilityFlags = CapabilityFlag(v)

	if !capabilityFlags.IsEnabled(ClientProtocol41) {
		return capabilityFlags, nil
	}

	capabilityFlags3, err := reader.PeekInt1()
	if err != nil {
		return 0, err
	}
	capabilityFlags |= (CapabilityFlag)(capabilityFlags3) << 16

	capabilityFlags4, err := reader.PeekInt1()
	if err != nil {
		return 0, err
	}
	capabilityFlags |= (CapabilityFlag)(capabilityFlags4) << 24

	return capabilityFlags, nil
}
