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

// HandshakeResponse represents a MySQL Handshake Response message.
type HandshakeResponse struct {
	*message
	capabilityFlags      CapabilityFlag
	maxPacketSize        uint32
	chanteSet            uint8
	username             string
	authResponseLength   uint8
	authResponse         string
	database             string
	clientPluginName     string
	attributes           map[string]string
	zstdCompressionLevel uint8
}

func newHandshakeResponseWithMessage(msg *message) *HandshakeResponse {
	return &HandshakeResponse{
		message:              msg,
		capabilityFlags:      0,
		maxPacketSize:        0,
		chanteSet:            0,
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
type HandshakeResponseOption func(*HandshakeResponse) error

// NewHandshakeResponse returns a new HandshakeResponse.
func NewHandshakeResponse(opts ...HandshakeResponseOption) (*HandshakeResponse, error) {
	h := newHandshakeResponseWithMessage(newMessage())
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

// NewHandshakeResponseFromReader returns a new HandshakeResponse from the reader.
func NewHandshakeResponseFromReader(reader io.Reader) (*HandshakeResponse, error) {
	var err error

	msg, err := NewMessageWithReader(reader)
	if err != nil {
		return nil, err
	}

	res := newHandshakeResponseWithMessage(msg)

	res.capabilityFlags, err = res.ReadCapabilityFlags()
	if err != nil {
		return nil, err
	}

	if !res.CapabilityFlags().IsEnabled(ClientProtocol41) {
		return nil, newErrNotSupported("HandshakeResponse320")
	}

	res.maxPacketSize, err = res.ReadInt4()
	if err != nil {
		return nil, err
	}

	res.chanteSet, err = res.ReadByte()
	if err != nil {
		return nil, err
	}

	err = res.SkipBytes(handshakeResponseFillerLen)
	if err != nil {
		return nil, err
	}

	res.username, err = res.ReadNullTerminatedString()
	if err != nil {
		return nil, err
	}

	if res.CapabilityFlags().IsEnabled(ClientPluginAuthLenencClientData) {
		res.authResponse, err = res.ReadLengthEncodedString()
		if err != nil {
			return nil, err
		}
	} else {
		res.authResponseLength, err = res.ReadByte()
		if err != nil {
			return nil, err
		}
		res.authResponse, err = res.ReadFixedLengthString(int(res.authResponseLength))
		if err != nil {
			return nil, err
		}
	}

	if res.CapabilityFlags().IsEnabled(ClientConnectWithDB) {
		res.database, err = res.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}
	}

	if res.CapabilityFlags().IsEnabled(ClientPluginAuth) {
		res.clientPluginName, err = res.ReadNullTerminatedString()
		if err != nil {
			return nil, err
		}
	}

	if res.CapabilityFlags().IsEnabled(ClientConnectAttrs) {
		nmap, err := res.ReadLengthEncodedInt()
		if err != nil {
			return nil, err
		}
		for i := uint64(0); i < nmap; i++ {
			key, err := res.ReadLengthEncodedString()
			if err != nil {
				return nil, err
			}
			value, err := res.ReadLengthEncodedString()
			if err != nil {
				return nil, err
			}
			res.attributes[key] = value
		}
	}

	if res.CapabilityFlags().IsEnabled(ClientZstdCompressionAlgorithm) {
		res.zstdCompressionLevel, err = res.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	return res, err
}

// CapabilityFlags returns the capability flags.
func (res *HandshakeResponse) CapabilityFlags() CapabilityFlag {
	return CapabilityFlag(res.capabilityFlags)
}

// MaxPacketSize returns the max packet size.
func (res *HandshakeResponse) MaxPacketSize() uint32 {
	return res.maxPacketSize
}

// Username returns the username.
func (res *HandshakeResponse) Username() string {
	return res.username
}

// AuthResponse returns the auth response.
func (res *HandshakeResponse) AuthResponse() string {
	return res.authResponse
}

// Database returns the database.
func (res *HandshakeResponse) Database() string {
	return res.database
}

// Bytes returns the message bytes.
func (res *HandshakeResponse) Bytes() ([]byte, error) {
	w := NewMessageWriter()

	if err := w.WriteCapabilityFlags(res.capabilityFlags); err != nil {
		return nil, err
	}

	if err := w.WriteInt4(res.maxPacketSize); err != nil {
		return nil, err
	}

	if err := w.WriteByte(res.chanteSet); err != nil {
		return nil, err
	}

	if err := w.WriteFillerBytes(0x00, handshakeResponseFillerLen); err != nil {
		return nil, err
	}

	if err := w.WriteNullTerminatedString(res.username); err != nil {
		return nil, err
	}

	if res.CapabilityFlags().IsEnabled(ClientPluginAuthLenencClientData) {
		if err := w.WriteLengthEncodedString(res.authResponse); err != nil {
			return nil, err
		}
	} else {
		if err := w.WriteByte(res.authResponseLength); err != nil {
			return nil, err
		}
		if err := w.WriteFixedLengthString(res.authResponse, int(res.authResponseLength)); err != nil {
			return nil, err
		}
	}

	if res.CapabilityFlags().IsEnabled(ClientConnectWithDB) {
		if err := w.WriteNullTerminatedString(res.database); err != nil {
			return nil, err
		}
	}

	if res.CapabilityFlags().IsEnabled(ClientPluginAuth) {
		if err := w.WriteNullTerminatedString(res.clientPluginName); err != nil {
			return nil, err
		}
	}

	if res.CapabilityFlags().IsEnabled(ClientConnectAttrs) {
		if err := w.WriteLengthEncodedInt(uint64(len(res.attributes))); err != nil {
			return nil, err
		}
		for key, value := range res.attributes {
			if err := w.WriteLengthEncodedString(key); err != nil {
				return nil, err
			}
			if err := w.WriteLengthEncodedString(value); err != nil {
				return nil, err
			}
		}
	}

	if res.CapabilityFlags().IsEnabled(ClientZstdCompressionAlgorithm) {
		if err := w.WriteByte(res.zstdCompressionLevel); err != nil {
			return nil, err
		}
	}

	res.message = NewMessage(
		MessageWithSequenceID(res.message.SequenceID()),
		MessageWithPayload(w.Bytes()),
	)

	return res.message.Bytes()
}
