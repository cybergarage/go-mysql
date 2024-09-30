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

// Handshake represents a MySQL Handshake message.
type Handshake struct {
	*message
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

func newHandshakeWithMessage(msg *message) *Handshake {
	return &Handshake{
		message:         msg,
		protocolVersion: uint8(ProtocolVersion10),
		serverVersion:   "",
		connectionID:    0,
		capabilityFlags: 0,
		characterSet:    uint8(CharacterSetUTF8),
		statusFlags:     0,
		authPluginData1: nil,
		authPluginData2: nil,
		authPluginName:  "",
	}
}

// HandshakeOption represents a MySQL Handshake option.
type HandshakeOption func(*Handshake) error

// WithHandshakeProtocolVersion sets the protocol version.
func WithHandshakeProtocolVersion(v ProtocolVersion) HandshakeOption {
	return func(h *Handshake) error {
		h.protocolVersion = uint8(v)
		return nil
	}
}

// WithHandshakeServerVersion sets the server version.
func WithHandshakeServerVersion(v string) HandshakeOption {
	return func(h *Handshake) error {
		h.serverVersion = v
		return nil
	}
}

// WithHandshakeConnectionID sets the connection ID.
func WithHandshakeConnectionID(v uint32) HandshakeOption {
	return func(h *Handshake) error {
		h.connectionID = v
		return nil
	}
}

// WithHandshakeCapabilityFlags sets the capability flags.
func WithHandshakeCapabilityFlags(v CapabilityFlag) HandshakeOption {
	return func(h *Handshake) error {
		h.capabilityFlags = uint32(v)
		return nil
	}
}

// WithHandshakeCharacterSet sets the character set.
func WithHandshakeCharacterSet(v CharacterSet) HandshakeOption {
	return func(h *Handshake) error {
		h.characterSet = uint8(v)
		return nil
	}
}

// WithHandshakeStatusFlags sets the status flags.
func WithHandshakeStatusFlags(v StatusFlag) HandshakeOption {
	return func(h *Handshake) error {
		h.statusFlags = uint16(v)
		return nil
	}
}

// WithHandshakeAuthPluginData1 sets the auth plugin data.
func WithHandshakeAuthPluginData(v []byte) HandshakeOption {
	return func(h *Handshake) error {
		if authPluginDataPartMaxLen < len(v) {
			return newInvalidLengthError("auth-plugin-data", len(v))
		}
		h.authPluginDataLen = uint8(len(v))
		if len(v) <= authPluginDataPart1Len {
			h.authPluginData1 = v
			h.authPluginData2 = nil
			return nil
		}
		h.authPluginData1 = v[:authPluginDataPart1Len]
		h.authPluginData2 = v[authPluginDataPart1Len:]
		return nil
	}
}

// WithHandshakeAuthPluginData2 sets the auth plugin name.
func WithHandshakeAuthPluginName(v string) HandshakeOption {
	return func(h *Handshake) error {
		h.authPluginName = v
		return nil
	}
}

// NewHandshake returns a new MySQL Handshake message.
func NewHandshake(opts ...HandshakeOption) (*Handshake, error) {
	h := newHandshakeWithMessage(newMessage())
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}
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

	h.authPluginData1, err = h.ReadFixedLengthBytes(authPluginDataPart1Len)
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

	h.authPluginDataLen = 0
	iv1, err := h.ReadByte()
	if err != nil {
		return nil, err
	}
	if h.CapabilityFlags().HasClientPluginAuth() {
		h.authPluginDataLen = iv1
	}

	_, err = h.ReadFixedLengthString(handshakeReservedLen) // Reserved
	if err != nil {
		return nil, err
	}

	if 0 < h.authPluginDataLen {
		authPluginDataLen := max(13, h.authPluginDataLen-8)
		h.authPluginData2, err = h.ReadFixedLengthBytes(int(authPluginDataLen))
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

func (h *Handshake) AuthPluginData() []byte {
	return append(h.authPluginData1, h.authPluginData2...)
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
	if err := w.WriteFixedLengthBytes(h.authPluginData1, authPluginDataPart1Len); err != nil {
		return nil, err
	}
	if err := w.WriteByte(0x00); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(uint16(h.capabilityFlags & 0xFFFF)); err != nil {
		return nil, err
	}
	if err := w.WriteByte(h.characterSet); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(h.statusFlags); err != nil {
		return nil, err
	}
	if err := w.WriteInt2(uint16(h.capabilityFlags >> 16)); err != nil {
		return nil, err
	}
	if h.CapabilityFlags().HasClientPluginAuth() {
		if err := w.WriteByte(uint8(len(h.authPluginData2) + authPluginDataPart1Len)); err != nil {
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
	if 0 < len(h.authPluginData2) {
		if err := w.WriteFixedLengthBytes(h.authPluginData2, len(h.authPluginData2)); err != nil {
			return nil, err
		}
	}
	if h.CapabilityFlags().HasClientPluginAuth() {
		if err := w.WriteNullTerminatedString(h.authPluginName); err != nil {
			return nil, err
		}
	}

	h.message = NewMessageWithPayload(w.Bytes())
	return h.message.Bytes()
}
