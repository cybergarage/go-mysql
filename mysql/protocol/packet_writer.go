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

// PacketWriter represents a packet writer of MySQL protocol.
type PacketWriter struct {
	*Writer
}

// NewPacketWriter returns a new packet writer.
func NewPacketWriter() *PacketWriter {
	return &PacketWriter{
		Writer: NewWriter(),
	}
}

// WriteCommandType writes a command type.
func (w *PacketWriter) WriteCommandType(cmd Command) error {
	return w.WriteByte(byte(cmd.Type()))
}

// WriteCapabilityFlags writes the capability flags.
func (w *PacketWriter) WriteCapabilityFlags(c CapabilityFlag) error {
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
