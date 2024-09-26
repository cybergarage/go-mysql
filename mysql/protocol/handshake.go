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

// Handshake represents a MySQL Handshake message.
type Handshake struct {
	*Message
	protocolVersion  uint8
	serverVersion    string
	connectionID     uint32
	authPluginData   string
	capabilityFlags1 uint16
	characterSet     uint8
	statusFlags      uint16
	capabilityFlags2 uint16
	authPluginName   string
}

// NewHandshake returns a new MySQL Handshake message.
func NewHandshakeWith(reader io.Reader) (*Handshake, error) {
	var err error

	msg, err := NewMessageWith(reader)
	if err != nil {
		return nil, err
	}

	h := &Handshake{
		Message: msg,
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

	h.authPluginData, err = h.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}

	h.ReadByte() // Filler

	h.capabilityFlags1, err = h.ReadInt2()
	if err != nil {
		return nil, err
	}

	h.characterSet, err = h.ReadByte()
	if err != nil {
		return nil, err
	}

	h.statusFlags, err = h.ReadInt2()
	if err != nil {
		return nil, err
	}

	h.capabilityFlags2, err = h.ReadInt2()
	if err != nil {
		return nil, err
	}

	authPluginDataLen, err := h.ReadByte()
	if err != nil {
		return nil, err
	}

	h.ReadByte() // Filler

	h.authPluginName, err = h.ReadNullTerminatedString()

	return h, err
}
