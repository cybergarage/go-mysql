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
// MySQL: Protocol::HandshakeResponse
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_handshake_response.html

const (
	handshakeResponseFillerLen = 23
)

// HandshakeResponse represents a MySQL Handshake Response packet.
type HandshakeResponse struct {
	*packet
	capabilityFlags      CapabilityFlag
	maxPacketSize        uint32
	charSet              uint8
	username             string
	authResponseLength   uint8
	authResponse         string
	database             string
	clientPluginName     string
	attributes           map[string]string
	zstdCompressionLevel uint8
}

func newHandshakeResponseWithPacket(pkt *packet) *HandshakeResponse {
	return &HandshakeResponse{
		packet:               pkt,
		capabilityFlags:      0,
		maxPacketSize:        0,
		charSet:              0,
		username:             "",
		authResponseLength:   0,
		authResponse:         "",
		database:             "",
		clientPluginName:     "",
		attributes:           make(map[string]string),
		zstdCompressionLevel: 0,
	}
}

// HandshakeResponseOption represents a HandshakeResponse option.
type HandshakeResponseOption func(*HandshakeResponse)

// NewHandshakeResponse returns a new HandshakeResponse.
func NewHandshakeResponse(opts ...HandshakeResponseOption) *HandshakeResponse {
	h := newHandshakeResponseWithPacket(newPacket())
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// NewHandshakeResponseFromReader returns a new HandshakeResponse from the reader.
func NewHandshakeResponseFromReader(reader io.Reader) (*HandshakeResponse, error) {
	var err error

	pktReader, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newHandshakeResponseWithPacket(pktReader)

	pkt.capabilityFlags, err = pkt.ReadCapabilityFlags()
	if err != nil {
		return nil, err
	}

	if !pkt.CapabilityFlags().IsEnabled(ClientProtocol41) {
		return nil, newErrNotSupported("HandshakeResponse320")
	}

	pkt.maxPacketSize, err = pkt.ReadInt4()
	if err != nil {
		return nil, err
	}

	pkt.charSet, err = pkt.ReadByte()
	if err != nil {
		return nil, err
	}

	err = pkt.SkipBytes(handshakeResponseFillerLen)
	if err != nil {
		return nil, err
	}

	pkt.username, err = pkt.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}

	if pkt.CapabilityFlags().IsEnabled(ClientPluginAuthLenencClientData) {
		pkt.authResponse, err = pkt.ReadLengthEncodedString()
		if err != nil {
			return nil, err
		}
	} else {
		pkt.authResponseLength, err = pkt.ReadByte()
		if err != nil {
			return nil, err
		}
		pkt.authResponse, err = pkt.ReadFixedLengthString(int(pkt.authResponseLength))
		if err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientConnectWithDB) {
		pkt.database, err = pkt.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientPluginAuth) {
		pkt.clientPluginName, err = pkt.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientConnectAttrs) {
		nmap, err := pkt.ReadLengthEncodedInt()
		if err != nil {
			return nil, err
		}
		for i := uint64(0); i < nmap; i++ {
			key, err := pkt.ReadLengthEncodedString()
			if err != nil {
				return nil, err
			}
			value, err := pkt.ReadLengthEncodedString()
			if err != nil {
				return nil, err
			}
			pkt.attributes[key] = value
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientZstdCompressionAlgorithm) {
		pkt.zstdCompressionLevel, err = pkt.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	return pkt, nil
}

// CapabilityFlags returns the capability flags.
func (pkt *HandshakeResponse) CapabilityFlags() CapabilityFlag {
	return CapabilityFlag(pkt.capabilityFlags)
}

// MaxPacketSize returns the max packet size.
func (pkt *HandshakeResponse) MaxPacketSize() uint32 {
	return pkt.maxPacketSize
}

// CharSet returns the character set.
func (pkt *HandshakeResponse) CharSet() uint8 {
	return pkt.charSet
}

// Username returns the username.
func (pkt *HandshakeResponse) Username() string {
	return pkt.username
}

// AuthResponse returns the auth response.
func (pkt *HandshakeResponse) AuthResponse() string {
	return pkt.authResponse
}

// Database returns the database.
func (pkt *HandshakeResponse) Database() string {
	return pkt.database
}

// ClientPluginName returns the client plugin name.
func (pkt *HandshakeResponse) ClientPluginName() string {
	return pkt.clientPluginName
}

// Bytes returns the packet bytes.
func (pkt *HandshakeResponse) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteCapabilityFlags(pkt.capabilityFlags); err != nil {
		return nil, err
	}

	if err := w.WriteInt4(pkt.maxPacketSize); err != nil {
		return nil, err
	}

	if err := w.WriteByte(pkt.charSet); err != nil {
		return nil, err
	}

	if err := w.WriteFillerBytes(0x00, handshakeResponseFillerLen); err != nil {
		return nil, err
	}

	if err := w.WriteNullTerminatedString(pkt.username); err != nil {
		return nil, err
	}

	if pkt.CapabilityFlags().IsEnabled(ClientPluginAuthLenencClientData) {
		if err := w.WriteLengthEncodedString(pkt.authResponse); err != nil {
			return nil, err
		}
	} else {
		if err := w.WriteByte(pkt.authResponseLength); err != nil {
			return nil, err
		}
		if err := w.WriteFixedLengthString(pkt.authResponse, int(pkt.authResponseLength)); err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientConnectWithDB) {
		if err := w.WriteNullTerminatedString(pkt.database); err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientPluginAuth) {
		if err := w.WriteNullTerminatedString(pkt.clientPluginName); err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientConnectAttrs) {
		if err := w.WriteLengthEncodedInt(uint64(len(pkt.attributes))); err != nil {
			return nil, err
		}
		for key, value := range pkt.attributes {
			if err := w.WriteLengthEncodedString(key); err != nil {
				return nil, err
			}
			if err := w.WriteLengthEncodedString(value); err != nil {
				return nil, err
			}
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientZstdCompressionAlgorithm) {
		if err := w.WriteByte(pkt.zstdCompressionLevel); err != nil {
			return nil, err
		}
	}

	pkt.SetPayload(w.Bytes())

	return pkt.packet.Bytes()
}
