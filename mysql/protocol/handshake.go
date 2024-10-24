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
	authPluginDataPartMaxLen = 13
	authPluginDataPart1Len   = 8
	handshakeReservedLen     = 10
)

// Handshake represents a MySQL Handshake packet.
type Handshake struct {
	*packet
	protocolVersion   uint8
	serverVersion     string
	connectionID      uint32
	capabilityFlags   uint32
	characterSet      uint8
	statusFlags       uint16
	authPluginDataLen uint8
	authPluginData1   []byte
	authPluginData2   []byte
	authPluginName    string
}

func newHandshakeWithPacket(msg *packet) *Handshake {
	return &Handshake{
		packet:            msg,
		protocolVersion:   uint8(ProtocolVersion10),
		serverVersion:     "",
		connectionID:      0,
		capabilityFlags:   uint32(DefaultServerCapabilities),
		characterSet:      uint8(CharSetUTF8),
		statusFlags:       uint16(DefaultServerCapabilities),
		authPluginDataLen: 0,
		authPluginData1:   nil,
		authPluginData2:   nil,
		authPluginName:    "",
	}
}

// HandshakeOption represents a MySQL Handshake option.
type HandshakeOption func(*Handshake) error

// WithHandshakeProtocolVersion sets the protocol version.
func WithHandshakeProtocolVersion(v ProtocolVersion) HandshakeOption {
	return func(pkt *Handshake) error {
		pkt.protocolVersion = uint8(v)
		return nil
	}
}

// WithHandshakeServerVersion sets the server version.
func WithHandshakeServerVersion(v string) HandshakeOption {
	return func(pkt *Handshake) error {
		pkt.serverVersion = v
		return nil
	}
}

// WithHandshakeConnectionID sets the connection ID.
func WithHandshakeConnectionID(v uint32) HandshakeOption {
	return func(pkt *Handshake) error {
		pkt.connectionID = v
		return nil
	}
}

// WithHandshakeCapabilityFlags sets the capability flags.
func WithHandshakeCapabilityFlags(v CapabilityFlag) HandshakeOption {
	return func(pkt *Handshake) error {
		pkt.capabilityFlags = uint32(v)
		return nil
	}
}

// WithHandshakeCharacterSet sets the character set.
func WithHandshakeCharacterSet(v CharSet) HandshakeOption {
	return func(pkt *Handshake) error {
		pkt.characterSet = uint8(v)
		return nil
	}
}

// WithHandshakeStatusFlags sets the status flags.
func WithHandshakeStatusFlags(v StatusFlag) HandshakeOption {
	return func(pkt *Handshake) error {
		pkt.statusFlags = uint16(v)
		return nil
	}
}

// WithHandshakeAuthPluginData1 sets the auth plugin data.
func WithHandshakeAuthPluginData(v []byte) HandshakeOption {
	return func(pkt *Handshake) error {
		if authPluginDataPartMaxLen < len(v) {
			return newErrInvalidLength("auth-plugin-data", len(v))
		}
		pkt.authPluginDataLen = uint8(len(v))
		if len(v) <= authPluginDataPart1Len {
			pkt.authPluginData1 = v
			pkt.authPluginData2 = nil
			return nil
		}
		pkt.authPluginData1 = v[:authPluginDataPart1Len]
		pkt.authPluginData2 = v[authPluginDataPart1Len:]
		return nil
	}
}

// WithHandshakeAuthPluginData2 sets the auth plugin name.
func WithHandshakeAuthPluginName(v string) HandshakeOption {
	return func(pkt *Handshake) error {
		pkt.authPluginName = v
		return nil
	}
}

// NewHandshake returns a new MySQL Handshake packet.
func NewHandshake(opts ...HandshakeOption) (*Handshake, error) {
	pkt := newHandshakeWithPacket(newPacket())
	for _, opt := range opts {
		if err := opt(pkt); err != nil {
			return nil, err
		}
	}
	return pkt, nil
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
	pkt.capabilityFlags = uint32(iv2)

	pkt.characterSet, err = pkt.ReadByte()
	if err != nil {
		return nil, err
	}

	pkt.statusFlags, err = pkt.ReadInt2()
	if err != nil {
		return nil, err
	}

	iv2, err = pkt.ReadInt2()
	if err != nil {
		return nil, err
	}
	pkt.capabilityFlags |= (uint32(iv2) << 16)

	pkt.authPluginDataLen = 0
	iv1, err := pkt.ReadByte()
	if err != nil {
		return nil, err
	}
	if pkt.CapabilityFlags().IsEnabled(ClientPluginAuth) {
		pkt.authPluginDataLen = iv1
	}

	_, err = pkt.ReadFixedLengthString(handshakeReservedLen) // Reserved
	if err != nil {
		return nil, err
	}

	if 0 < pkt.authPluginDataLen {
		authPluginDataLen := max(13, pkt.authPluginDataLen-8)
		pkt.authPluginData2, err = pkt.ReadFixedLengthBytes(int(authPluginDataLen))
		if err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientPluginAuth) {
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

// CapabilityFlags returns the capability flags.
func (pkt *Handshake) CapabilityFlags() CapabilityFlag {
	return CapabilityFlag(pkt.capabilityFlags)
}

// CharacterSet returns the character set.
func (pkt *Handshake) CharacterSet() CharSet {
	return CharSet(pkt.characterSet)
}

// StatusFlags returns the status flags.
func (pkt *Handshake) StatusFlags() StatusFlag {
	return StatusFlag(pkt.statusFlags)
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
	if err := w.WriteInt2(uint16(pkt.capabilityFlags & 0xFFFF)); err != nil {
		return nil, err
	}
	if err := w.WriteByte(pkt.characterSet); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(pkt.statusFlags); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(uint16(pkt.capabilityFlags >> 16)); err != nil {
		return nil, err
	}
	if pkt.CapabilityFlags().IsEnabled(ClientPluginAuth) {
		if err := w.WriteByte(uint8(len(pkt.authPluginData2) + authPluginDataPart1Len)); err != nil {
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
		if err := w.WriteFixedLengthBytes(pkt.authPluginData2, len(pkt.authPluginData2)); err != nil {
			return nil, err
		}
	}
	if pkt.CapabilityFlags().IsEnabled(ClientPluginAuth) {
		if err := w.WriteNullTerminatedString(pkt.authPluginName); err != nil {
			return nil, err
		}
	}

	pkt.SetPayload(w.Bytes())

	return pkt.packet.Bytes()
}
