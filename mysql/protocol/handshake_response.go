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
	"errors"
	"io"

	"github.com/cybergarage/go-mysql/mysql/auth"
	"github.com/cybergarage/go-mysql/mysql/encoding/binary"
)

// MySQL: Connection Phase
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase.html
// MySQL: Protocol::HandshakeResponse
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_handshake_response.html
// Connecting - MariaDB Knowledge Base
// https://mariadb.com/kb/en/connection/

const (
	handshakeResponseFillerLen = 23
)

// HandshakeResponse represents a MySQL Handshake Response packet.
type HandshakeResponse struct {
	*packet
	Capabilitys        Capability
	maxPacketSize      uint32
	charSet            uint8
	username           string
	authResponseLength uint8
	authResponse       []byte
	database           string
	clientPluginName   string
	*AttributeMap
	zstdCompressionLevel uint8
}

func newHandshakeResponseWithPacket(pkt *packet) *HandshakeResponse {
	return &HandshakeResponse{
		packet:               pkt,
		Capabilitys:          0,
		maxPacketSize:        0,
		charSet:              0,
		username:             "",
		authResponseLength:   0,
		authResponse:         []byte{},
		database:             "",
		clientPluginName:     "",
		AttributeMap:         NewAttributeMap(),
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

	resPkt, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	res := newHandshakeResponseWithPacket(resPkt)

	res.Capabilitys, err = res.ReadCapability()
	if err != nil {
		return nil, err
	}

	if !res.Capability().HasCapability(ClientProtocol41) {
		return nil, newErrNotSupported("HandshakeResponse320")
	}

	res.maxPacketSize, err = res.ReadInt4()
	if err != nil {
		return nil, err
	}

	res.charSet, err = res.ReadByte()
	if err != nil {
		return nil, err
	}

	err = res.SkipBytes(handshakeResponseFillerLen)
	if err != nil {
		return nil, err
	}

	res.username, err = res.ReadNullTerminatedString()
	if err != nil {
		if errors.Is(err, io.EOF) {
			res.username = ""
			return res, nil
		}
		return nil, err
	}

	if res.Capability().HasCapability(ClientPluginAuthLenencClientData) {
		res.authResponse, err = res.ReadLengthEncodedBytes()
		if err != nil {
			return nil, err
		}
	} else {
		res.authResponseLength, err = res.ReadByte()
		if err != nil {
			return nil, err
		}
		res.authResponse, err = res.ReadFixedLengthBytes(int(res.authResponseLength))
		if err != nil {
			return nil, err
		}
	}

	if res.Capability().HasCapability(ClientConnectWithDB) {
		res.database, err = res.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}
	}

	if res.Capability().HasCapability(ClientPluginAuth) {
		res.clientPluginName, err = res.ReadNullTerminatedString()
		if err != nil {
			if errors.Is(err, io.EOF) {
				res.clientPluginName = ""
				return res, nil
			}
			return nil, err
		}
	}

	if res.Capability().HasCapability(ClientConnectAttrs) {
		attrSize, err := res.ReadLengthEncodedInt()
		if err != nil {
			return nil, err
		}
		readAttrSize := 0
		for readAttrSize < int(attrSize) {
			key, err := res.ReadLengthEncodedString()
			if err != nil {
				return nil, err
			}
			keyLen := len(key)
			readAttrSize += binary.LengthEncodeIntSize(uint64(keyLen)) + keyLen

			value, err := res.ReadLengthEncodedString()
			if err != nil {
				return nil, err
			}
			valueLen := len(value)
			readAttrSize += binary.LengthEncodeIntSize(uint64(valueLen)) + valueLen

			res.AddAttribute(key, value)
		}
	}

	if res.Capability().HasCapability(ClientZstdCompressionAlgorithm) {
		res.zstdCompressionLevel, err = res.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

// Capabilitys returns the capability flags.
func (pkt *HandshakeResponse) Capability() Capability {
	return Capability(pkt.Capabilitys)
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
func (pkt *HandshakeResponse) AuthResponse() []byte {
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

// AutMethod returns the authentication method.
func (pkt *HandshakeResponse) AutMethod() (auth.AuthMethod, error) {
	if len(pkt.clientPluginName) == 0 {
		return auth.MySQLAuthenticationNone, nil
	}
	return auth.NewAuthMethodFromID(pkt.clientPluginName)
}

// ZstdCompressionLevel returns the Zstd compression level.
func (pkt *HandshakeResponse) ZstdCompressionLevel() uint8 {
	return pkt.zstdCompressionLevel
}

// Bytes returns the packet bytes.
func (pkt *HandshakeResponse) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteCapability(pkt.Capabilitys); err != nil {
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

	if pkt.Capability().HasCapability(ClientPluginAuthLenencClientData) {
		if err := w.WriteLengthEncodedBytes(pkt.authResponse); err != nil {
			return nil, err
		}
	} else {
		if err := w.WriteByte(pkt.authResponseLength); err != nil {
			return nil, err
		}
		if err := w.WriteFixedLengthBytes(pkt.authResponse, int(pkt.authResponseLength)); err != nil {
			return nil, err
		}
	}

	if pkt.Capability().HasCapability(ClientConnectWithDB) {
		if err := w.WriteNullTerminatedString(pkt.database); err != nil {
			return nil, err
		}
	}

	hasNextSection := pkt.Capability().HasCapability(ClientConnectAttrs) && (0 < len(pkt.AttributeKeys()))

	if pkt.Capability().HasCapability(ClientPluginAuth) {
		if 0 < len(pkt.clientPluginName) || hasNextSection {
			if err := w.WriteNullTerminatedString(pkt.clientPluginName); err != nil {
				return nil, err
			}
		}
	}

	hasNextSection = pkt.Capability().HasCapability(ClientZstdCompressionAlgorithm)

	if pkt.Capability().HasCapability(ClientConnectAttrs) {
		if 0 < len(pkt.AttributeKeys()) || hasNextSection {
			attrWriter := NewPacketWriter()
			for _, key := range pkt.AttributeKeys() {
				value, _ := pkt.LookupAttribute(key)
				if err := attrWriter.WriteLengthEncodedString(key); err != nil {
					return nil, err
				}
				if err := attrWriter.WriteLengthEncodedString(value); err != nil {
					return nil, err
				}
			}
			attrBytes := attrWriter.Bytes()
			attrSize := len(attrBytes)
			if err := w.WriteLengthEncodedInt(uint64(attrSize)); err != nil {
				return nil, err
			}
			if _, err := w.WriteBytes(attrBytes); err != nil {
				return nil, err
			}
		}
	}

	if pkt.Capability().HasCapability(ClientZstdCompressionAlgorithm) {
		if err := w.WriteByte(pkt.zstdCompressionLevel); err != nil {
			return nil, err
		}
	}

	pkt.SetPayload(w.Bytes())

	return pkt.packet.Bytes()
}
