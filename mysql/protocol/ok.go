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

// MySQL: OK_Packet
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_ok_packet.html

const (
	okPacketHeader = 0x00
)

// OK represents a MySQL OK packet.
type OK struct {
	*packet
	header uint8
}

// OKOption represents a MySQL OK packet option.
type OKOption func(*OK) error

func newOKPacket(p *packet, opts ...OKOption) (*OK, error) {
	pkt := &OK{
		packet: p,
	}
	for _, opt := range opts {
		if err := opt(pkt); err != nil {
			return nil, err
		}
	}
	return pkt, nil
}

// NewOK returns a new OK packet.
func NewOK(opts ...OKOption) (*OK, error) {
	return newOKPacket(nil)
}

// NewOKFromReader returns a new OK packet from the reader.
func NewOKFromReader(reader io.Reader, opts ...OKOption) (*OK, error) {
	var err error

	pktReader, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt, err := newOKPacket(pktReader, opts...)
	if err != nil {
		return nil, err
	}

	// header
	pkt.header, err = pkt.ReadByte()
	if err != nil {
		return nil, err
	}
	if (pkt.header != okPacketHeader) && (pkt.header != okPacketHeader) {
		return nil, newErrInvalitHeader("OK", pkt.header)
	}

	return pkt, err
}

// Bytes returns a byte sequence of the OK packet.
func (pkt *OK) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	// header
	if err := w.WriteByte(errPacketHeader); err != nil {
		return nil, err
	}

	pkt.packet = NewPacket(
		PacketWithSequenceID(pkt.packet.SequenceID()),
		PacketWithPayload(w.Bytes()),
	)

	return pkt.packet.Bytes()
}
