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
// Connecting - MariaDB Knowledge Base
// https://mariadb.com/kb/en/connection/

// ProtocolVersion represents a MySQL Protocol Version.
type ProtocolVersion uint8

const (
	ProtocolVersion10 ProtocolVersion = 10
)

const (
	authPluginDataPart1Len               = 8
	authPluginDataPart2Len               = authPluginDataPart2MySQLLen
	authPluginDataPart2MySQLLen          = (authPluginDataPart2MariaDBMaxLen + authPluginDataPart2MariaDBReserveLen)
	authPluginDataPart2MariaDBMaxLen     = 12
	authPluginDataPart2MariaDBReserveLen = 1
	authPluginDataPartMaxLen             = authPluginDataPartMaxMySQLLen
	authPluginDataPartMaxMySQLLen        = (authPluginDataPart1Len + authPluginDataPart2MySQLLen)
	authPluginDataPartMaxMariaDBLen      = (authPluginDataPart1Len + authPluginDataPart2MariaDBMaxLen)
	handshakeReservedLen                 = 10
)

// Handshake represents a MySQL Handshake packet.
type Handshake struct {
	*packet
	protocolVersion   uint8
	serverVersion     string
	connectionID      uint32
	characterSet      uint8
	authPluginDataLen uint8
	authPluginData1   []byte
	authPluginData2   []byte
	authPluginName    string
}

func newHandshakeWithPacket(msg *packet) *Handshake {
	pkt := &Handshake{
		packet:            msg,
		protocolVersion:   uint8(ProtocolVersion10),
		serverVersion:     "",
		connectionID:      0,
		characterSet:      uint8(CharSetUTF8),
		authPluginDataLen: 0,
		authPluginData1:   nil,
		authPluginData2:   nil,
		authPluginName:    "",
	}
	pkt.SetCapability(DefaultServerCapability)
	pkt.SetServerStatus(DefaultServerStatus)
	return pkt
}

// HandshakeOption represents a MySQL Handshake option.
type HandshakeOption func(*Handshake)

// WithHandshakeProtocolVersion sets the protocol version.
func WithHandshakeProtocolVersion(v ProtocolVersion) HandshakeOption {
	return func(pkt *Handshake) {
		pkt.protocolVersion = uint8(v)
	}
}

// WithHandshakeServerVersion sets the server version.
func WithHandshakeServerVersion(v string) HandshakeOption {
	return func(pkt *Handshake) {
		pkt.serverVersion = v
	}
}

// WithHandshakeConnectionID sets the connection ID.
func WithHandshakeConnectionID(v uint32) HandshakeOption {
	return func(pkt *Handshake) {
		pkt.connectionID = v
	}
}

// WithHandshakeCapability sets the capability flags.
func WithHandshakeCapability(v Capability) HandshakeOption {
	return func(pkt *Handshake) {
		pkt.SetCapability(v)
	}
}

// WithHandshakeCharacterSet sets the character set.
func WithHandshakeCharacterSet(v CharSet) HandshakeOption {
	return func(pkt *Handshake) {
		pkt.characterSet = uint8(v)
	}
}

// WithHandshakeServerStatusFlags sets the status flags.
func WithHandshakeServerStatusFlags(v ServerStatus) HandshakeOption {
	return func(pkt *Handshake) {
		pkt.SetServerStatus(v)
	}
}

// WithHandshakeAuthPluginData1 sets the auth plugin data.
func WithHandshakeAuthPluginData(v []byte) HandshakeOption {
	// NOTE: mysql-server 5.7 send_server_handshake_packet()
	// https://github.com/mysql/mysql-server/blob/5.7/sql/auth/sql_authentication.cc#L512
	// " \0 byte, terminating the second part of a scramble"
	return func(pkt *Handshake) {
		if authPluginDataPartMaxLen < len(v) {
			v = v[:authPluginDataPartMaxLen]
		}
		pkt.authPluginDataLen = uint8(len(v) & 0xFF)
		if len(v) <= authPluginDataPart1Len {
			pkt.authPluginData1 = v
			pkt.authPluginData2 = nil
		}
		pkt.authPluginData1 = v[:authPluginDataPart1Len]
		pkt.authPluginData2 = v[authPluginDataPart1Len:]
	}
}

// WithHandshakeAuthPluginData2 sets the auth plugin name.
func WithHandshakeAuthPluginName(v string) HandshakeOption {
	return func(pkt *Handshake) {
		pkt.authPluginName = v
	}
}

// NewHandshake returns a new MySQL Handshake packet.
func NewHandshake(opts ...HandshakeOption) *Handshake {
	pkt := newHandshakeWithPacket(newPacket())
	pkt.serverVersion = SupportVersion
	for _, opt := range opts {
		opt(pkt)
	}
	return pkt
}

// NewHandshakeFromReader returns a new MySQL Handshake packet from the specified reader.
func NewHandshakeFromReader(reader io.Reader) (*Handshake, error) {
	var err error

	pktReader, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newHandshakeWithPacket(pktReader)

	pkt.protocolVersion, err = pkt.ReadByte()
	if err != nil {
		return nil, err
	}

	pkt.serverVersion, err = pkt.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}

	pkt.connectionID, err = pkt.ReadInt4()
	if err != nil {
		return nil, err
	}

	pkt.authPluginData1, err = pkt.ReadFixedLengthBytes(authPluginDataPart1Len)
	if err != nil {
		return nil, err
	}

	_, err = pkt.ReadByte() // Filler
	if err != nil {
		return nil, err
	}

	iv2, err := pkt.ReadInt2()
	if err != nil {
		return nil, err
	}
	pkt.SetCapability(Capability(iv2))

	pkt.characterSet, err = pkt.ReadByte()
	if err != nil {
		return nil, err
	}

	iv2, err = pkt.ReadInt2()
	if err != nil {
		return nil, err
	}
	pkt.SetServerStatus(ServerStatus(iv2))

	iv2, err = pkt.ReadInt2()
	if err != nil {
		return nil, err
	}
	pkt.SetCapabilityEnabled(Capability((uint32(iv2) << 16)))

	pkt.authPluginDataLen = 0
	iv1, err := pkt.ReadByte()
	if err != nil {
		return nil, err
	}
	if pkt.Capability().IsEnabled(ClientPluginAuth) {
		pkt.authPluginDataLen = iv1
	}

	_, err = pkt.ReadFixedLengthString(handshakeReservedLen) // Reserved
	if err != nil {
		return nil, err
	}

	authPluginData2Len := max(13, pkt.authPluginDataLen-8)
	if 0 < authPluginData2Len {
		// mysql-server 5.7 send_server_handshake_packet()
		// https://github.com/mysql/mysql-server/blob/5.7/sql/auth/sql_authentication.cc#L512
		// NOTE: " \0 byte, terminating the second part of a scramble"
		pkt.authPluginData2, err = pkt.ReadFixedLengthBytes(int(authPluginData2Len))
		if err != nil {
			return nil, err
		}
		pkt.authPluginData2 = pkt.authPluginData2[:authPluginData2Len-1]
	}

	if pkt.Capability().IsEnabled(ClientPluginAuth) {
		pkt.authPluginName, err = pkt.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}
	}

	return pkt, err
}

// ProtocolVersion returns the protocol version.
func (pkt *Handshake) ProtocolVersion() ProtocolVersion {
	return ProtocolVersion(pkt.protocolVersion)
}

// ServerVersion returns the server version.
func (pkt *Handshake) ServerVersion() string {
	return pkt.serverVersion
}

// ConnectionID returns the connection ID.
func (pkt *Handshake) ConnectionID() uint32 {
	return pkt.connectionID
}

// AuthPluginData returns the auth plugin data.
func (pkt *Handshake) AuthPluginData() []byte {
	return append(pkt.authPluginData1, pkt.authPluginData2...)
}

// CharacterSet returns the character set.
func (pkt *Handshake) CharacterSet() CharSet {
	return CharSet(pkt.characterSet)
}

// AuthPluginName returns the auth plugin name.
func (pkt *Handshake) AuthPluginName() string {
	return pkt.authPluginName
}

// Bytes returns the packet bytes.
func (pkt *Handshake) Bytes() ([]byte, error) {
	w := NewWriter()
	if err := w.WriteByte(pkt.protocolVersion); err != nil {
		return nil, err
	}
	if err := w.WriteNullTerminatedString(pkt.serverVersion); err != nil {
		return nil, err
	}
	if err := w.WriteInt4(pkt.connectionID); err != nil {
		return nil, err
	}
	if err := w.WriteFixedLengthBytes(pkt.authPluginData1, authPluginDataPart1Len); err != nil {
		return nil, err
	}
	if err := w.WriteByte(0x00); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(uint16(pkt.Capability() & 0xFFFF)); err != nil {
		return nil, err
	}
	if err := w.WriteByte(pkt.characterSet); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(uint16(pkt.ServerStatus())); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(uint16(pkt.Capability() >> 16)); err != nil {
		return nil, err
	}
	if pkt.Capability().IsEnabled(ClientPluginAuth) {
		// mysql-server 5.7 send_server_handshake_packet()
		// https://github.com/mysql/mysql-server/blob/5.7/sql/auth/sql_authentication.cc#L512
		// NOTE: " \0 byte, terminating the second part of a scramble"
		authPluginDataLen := authPluginDataPart1Len + len(pkt.authPluginData2) + 1
		if err := w.WriteByte(uint8(authPluginDataLen)); err != nil {
			return nil, err
		}
	} else {
		if err := w.WriteByte(0x00); err != nil {
			return nil, err
		}
	}
	if err := w.WriteFixedLengthString("", handshakeReservedLen); err != nil {
		return nil, err
	}
	if 0 < len(pkt.authPluginData2) {
		// mysql-server 5.7 send_server_handshake_packet()
		// https://github.com/mysql/mysql-server/blob/5.7/sql/auth/sql_authentication.cc#L512
		// NOTE: " \0 byte, terminating the second part of a scramble"
		if err := w.WriteNullTerminatedBytes(pkt.authPluginData2); err != nil {
			return nil, err
		}
	}
	if pkt.Capability().IsEnabled(ClientPluginAuth) {
		if err := w.WriteNullTerminatedString(pkt.authPluginName); err != nil {
			return nil, err
		}
	}

	pkt.SetPayload(w.Bytes())

	return pkt.packet.Bytes()
}
