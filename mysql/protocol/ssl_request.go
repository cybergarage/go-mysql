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
// MySQL: Protocol::SSLRequest
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_SSLRequest.html
// MySQL: Protocol::SSLRequestV10
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_packets_protocol_SSLRequest_v10.html

const (
	sslRequestFillerLen = 23
)

// SSLRequest represents a MySQL SSLRequest packet.
type SSLRequest struct {
	*packet
	capabilityFlags CapabilityFlag
	characterSet    uint8
	maxPacketSize   uint32
}

func newSSLRequestWithPacket(msg *packet) *SSLRequest {
	return &SSLRequest{
		packet:          msg,
		capabilityFlags: DefaultSSLRequestCapabilities,
		characterSet:    uint8(DefaultCharset),
		maxPacketSize:   DefaultMaxPacketSize,
	}
}

// SSLRequestOption represents a MySQL SSLRequest option.
type SSLRequestOption func(*SSLRequest) error

// WithSSLRequestCapabilityFlags sets the capability flags.
func WithSSLRequestCapabilityFlags(v CapabilityFlag) SSLRequestOption {
	return func(h *SSLRequest) error {
		h.capabilityFlags = v
		return nil
	}
}

// WithSSLRequestCharacterSet sets the character set.
func WithSSLRequestCharacterSet(v CharacterSet) SSLRequestOption {
	return func(h *SSLRequest) error {
		h.characterSet = uint8(v)
		return nil
	}
}

// NewSSLRequest returns a new MySQL SSLRequest packet.
func NewSSLRequest(opts ...SSLRequestOption) (*SSLRequest, error) {
	h := newSSLRequestWithPacket(newPacket())
	for _, opt := range opts {
		if err := opt(h); err != nil {
			return nil, err
		}
	}
	return h, nil
}

// NewSSLRequestFromReader returns a new MySQL SSLRequest packet from the specified reader.
func NewSSLRequestFromReader(reader io.Reader) (*SSLRequest, error) {
	var err error

	msg, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newSSLRequestWithPacket(msg)

	pkt.capabilityFlags, err = pkt.ReadCapabilityFlags()
	if err != nil {
		return nil, err
	}

	if pkt.capabilityFlags.IsEnabled(ClientProtocol41) {
		pkt.maxPacketSize, err = pkt.ReadInt4()
		if err != nil {
			return nil, err
		}

		pkt.characterSet, err = pkt.ReadInt1()
		if err != nil {
			return nil, err
		}

		err = pkt.SkipBytes(sslRequestFillerLen)
		if err != nil {
			return nil, err
		}
	} else {
		pkt.maxPacketSize, err = pkt.ReadInt3()
		if err != nil {
			return nil, err
		}
	}

	return pkt, err
}

// CapabilityFlags returns the capability flags.
func (pkt *SSLRequest) CapabilityFlags() CapabilityFlag {
	return CapabilityFlag(pkt.capabilityFlags)
}

// CharacterSet returns the character set.
func (pkt *SSLRequest) CharacterSet() CharacterSet {
	return CharacterSet(pkt.characterSet)
}

// Bytes returns the packet bytes.
func (pkt *SSLRequest) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteCapabilityFlags(pkt.capabilityFlags); err != nil {
		return nil, err
	}

	if pkt.CapabilityFlags().IsEnabled(ClientProtocol41) {
		if err := w.WriteInt4(pkt.maxPacketSize); err != nil {
			return nil, err
		}
		if err := w.WriteInt1(pkt.characterSet); err != nil {
			return nil, err
		}
		if err := w.WriteFillerBytes(0x00, sslRequestFillerLen); err != nil {
			return nil, err
		}
	} else {
		if err := w.WriteInt3(pkt.maxPacketSize); err != nil {
			return nil, err
		}
	}

	pkt.packet = NewPacket(
		PacketWithSequenceID(pkt.packet.SequenceID()),
		PacketWithPayload(w.Bytes()),
	)

	return pkt.packet.Bytes()
}
