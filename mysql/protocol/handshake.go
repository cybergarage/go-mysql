// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
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

// MySQL: Connection Phase
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase.html
// MySQL: Protocol::Handshake
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_handshake.html
// MySQL: Protocol::HandshakeV10
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_handshake_v10.html

// ProtocolVersion represents a MySQL Protocol Version.
type ProtocolVersion uint8

const (
	ProtocolVersion10 ProtocolVersion = 10
)

const (
	authPluginDataPart1Len = 8
)

// Handshake represents a MySQL Handshake message.
type Handshake struct {
	*message
	protocolVersion uint8
	serverVersion   string
	connectionID    uint32
	authPluginData1 string
	capabilityFlags uint32
	characterSet    uint8
	statusFlags     uint16
	authPluginData2 string
	authPluginName  string
}

func newHandshakeWithMessage(msg *message) *Handshake {
	return &Handshake{
		message: msg,
	}
}

// NewHandshake returns a new MySQL Handshake message.
func NewHandshake() (*Handshake, error) {
	h := newHandshakeWithMessage(newMessage())
	return h, nil
}

// NewHandshakeFromReader returns a new MySQL Handshake message from the specified reader.
func NewHandshakeFromReader(reader io.Reader) (*Handshake, error) {
	var err error

	msg, err := NewMessageWithReader(reader)
	if err != nil {
		return nil, err
	}

	h := newHandshakeWithMessage(msg)

	h.protocolVersion, err = h.ReadByte()
	if err != nil {
		return nil, err
	}

	h.serverVersion, err = h.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}

	h.connectionID, err = h.ReadInt4()
	if err != nil {
		return nil, err
	}

	h.authPluginData1, err = h.ReadFixedLengthString(authPluginDataPart1Len)
	if err != nil {
		return nil, err
	}

	_, err = h.ReadByte() // Filler
	if err != nil {
		return nil, err
	}

	iv2, err := h.ReadInt2()
	if err != nil {
		return nil, err
	}
	h.capabilityFlags = uint32(iv2)

	h.characterSet, err = h.ReadByte()
	if err != nil {
		return nil, err
	}

	h.statusFlags, err = h.ReadInt2()
	if err != nil {
		return nil, err
	}

	iv2, err = h.ReadInt2()
	if err != nil {
		return nil, err
	}
	h.capabilityFlags |= (uint32(iv2) << 16)

	authPluginDataLen := uint8(0)
	iv1, err := h.ReadByte()
	if err != nil {
		return nil, err
	}
	if h.CapabilityFlags().HasClientPluginAuth() {
		authPluginDataLen = iv1
	}

	_, err = h.ReadFixedLengthString(10) // Reserved
	if err != nil {
		return nil, err
	}

	if 0 < authPluginDataLen {
		authPluginDataLen = max(13, authPluginDataLen-8)
		h.authPluginData2, err = h.ReadFixedLengthString(int(authPluginDataLen))
		if err != nil {
			return nil, err
		}
	}

	if h.CapabilityFlags().HasClientPluginAuth() {
		h.authPluginName, err = h.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}
	}

	return h, err
}

func (h *Handshake) ProtocolVersion() ProtocolVersion {
	return ProtocolVersion(h.protocolVersion)
}

func (h *Handshake) ServerVersion() string {
	return h.serverVersion
}

func (h *Handshake) ConnectionID() uint32 {
	return h.connectionID
}

func (h *Handshake) AuthPluginData() string {
	return h.authPluginData1
}

func (h *Handshake) CapabilityFlags() CapabilityFlag {
	return CapabilityFlag(h.capabilityFlags)
}

func (h *Handshake) CharacterSet() CharacterSet {
	return CharacterSet(h.characterSet)
}

func (h *Handshake) StatusFlags() StatusFlag {
	return StatusFlag(h.statusFlags)
}

func (h *Handshake) AuthPluginName() string {
	return h.authPluginName
}

// Bytes returns the message bytes.
func (h *Handshake) Bytes() ([]byte, error) {
	w := NewWriter()
	if err := w.WriteByte(h.protocolVersion); err != nil {
		return nil, err
	}
	if err := w.WriteNullTerminatedString(h.serverVersion); err != nil {
		return nil, err
	}
	if err := w.WriteInt4(h.connectionID); err != nil {
		return nil, err
	}
	if err := w.WriteFixedLengthString(h.authPluginData1, authPluginDataPart1Len); err != nil {
		return nil, err
	}
	if err := w.WriteByte(0x00); err != nil {
		return nil, err
	}
	if err := w.WriteByte(0x00); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(h.statusFlags); err != nil {
		return nil, err
	}
	h.message = NewMessageWithPayload(w.Bytes())
	return h.message.Bytes()
}
