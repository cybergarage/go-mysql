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

// NewHandshakeFromReader returns a new MySQL Handshake message from the specified reader.
func NewHandshakeFromReader(reader io.Reader) (*Handshake, error) {
	var err error

	msg, err := NewMessageWith(reader)
	if err != nil {
		return nil, err
	}

	h := &Handshake{
		message: msg,
	}

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

	h.authPluginData1, err = h.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}

	h.ReadByte() // Filler

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
	h.capabilityFlags &= (uint32(iv2) << 16)

	hasClientPluginAuthFlag := (CapabilityFlag(h.capabilityFlags) & CapabilityFlagClientPluginAuth) != 0
	authPluginDataLen := uint8(0)
	iv1, err := h.ReadByte()
	if err != nil {
		return nil, err
	}
	if hasClientPluginAuthFlag {
		authPluginDataLen = iv1
	}

	_, err = h.ReadFixedLengthString(10) // Reserved
	if err != nil {
		return nil, err
	}

	authPluginDataLen = max(13, authPluginDataLen-8)
	h.authPluginData2, err = h.ReadFixedLengthString(int(authPluginDataLen))
	if err != nil {
		return nil, err
	}

	if hasClientPluginAuthFlag {
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

func (h *Handshake) CharacterSet() uint8 {
	return h.characterSet
}

func (h *Handshake) StatusFlags() uint16 {
	return h.statusFlags
}

func (h *Handshake) AuthPluginName() string {
	return h.authPluginName
}
