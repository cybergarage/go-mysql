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
	"slices"
)

// MySQL: Protocol Basics
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basics.html
// MySQL: MySQL Packets
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_packets.html
// MySQL: Packets
// https://dev.mysql.com/doc/dev/mysql-server/latest/mysqlx_protocol_packets.html
// MariaDB protocol difference with MySQL - MariaDB Knowledge Base
// https://mariadb.com/kb/en/mariadb-protocol-difference-with-mysql/

// Packet represents a MySQL packet.
type Packet interface {
	PacketIdentifier
	// SetSequenceID sets the packet sequence ID.
	SetSequenceID(n SequenceID)
	// SetPayload sets the packet payload.
	SetPayload(payload []byte)
	// PayloadLength returns the packet payload length.
	PayloadLength() uint32
	// SetCapability sets the packet capability flags.
	SetCapability(Capability)
	// SequenceID returns the packet sequence ID.
	SequenceID() SequenceID
	// Payload returns the packet payload.
	Payload() []byte
	// Capability returns the packet capability flags.
	Capability() Capability
	// Reader returns the packet reader.
	Reader() *PacketReader
	// Bytes returns the packet bytes.
	Bytes() ([]byte, error)
}

// packet represents a MySQL packet.
type packet struct {
	*PacketReader
	payloadLength uint32
	sequenceID    SequenceID
	payload       []byte
	Capabilitys   Capability
}

func newPacket() *packet {
	return &packet{
		PacketReader:  nil,
		payloadLength: 0,
		sequenceID:    SequenceID(0),
		payload:       nil,
		Capabilitys:   0,
	}
}

// PacketOption represents a packet option.
type PacketOption func(Packet)

// WithPacketPayload returns a packet option to set the payload.
func WithPacketPayload(payload []byte) PacketOption {
	return func(pkt Packet) {
		pkt.SetPayload(payload)
	}
}

// WithPacketSequenceID returns a packet option to set the sequence ID.
func WithPacketSequenceID(n SequenceID) PacketOption {
	return func(pkt Packet) {
		pkt.SetSequenceID(n)
	}
}

// WithPacketCapability returns a packet option to set the capability flags.
func WithPacketCapability(flags Capability) PacketOption {
	return func(pkt Packet) {
		pkt.SetCapability(flags)
	}
}

// NewPacket returns a new MySQL packet.
func NewPacket(opts ...PacketOption) *packet {
	pkt := newPacket()
	pkt.SetCapabilityEnabled(ClientProtocol41)
	pkt.SetOptions(opts...)
	return pkt
}

// NewPacketWithReader returns a new MySQL packet from the reader.
func NewPacketWithReader(reader io.Reader) (*packet, error) {
	return NewPacketWithPacketReader(NewPacketReaderWith(reader))
}

// NewPacketWithPacketReader returns a new MySQL packet from the packet reader.
func NewPacketWithPacketReader(reader *PacketReader) (*packet, error) {
	pkt := newPacket()
	pkt.PacketReader = reader
	err := pkt.ReadHeader()
	if err != nil {
		return nil, err
	}
	err = pkt.ReadPayload()
	if err != nil {
		return nil, err
	}
	return pkt, nil
}

// NewPacketHeaderWithReader returns a new MySQL packet from the reader.
func NewPacketHeaderWithReader(reader io.Reader) (*packet, error) {
	pkt := newPacket()
	pkt.PacketReader = NewPacketReaderWith(reader)
	err := pkt.ReadHeader()
	if err != nil {
		return nil, err
	}
	return pkt, nil
}

// SetOptions sets the options.
func (pkt *packet) SetOptions(opts ...PacketOption) {
	for _, opt := range opts {
		opt(pkt)
	}
}

// ReadHeader reads the packet header.
func (pkt *packet) ReadHeader() error {
	// Read the payload length
	payloadLengthBuf := make([]byte, 3)
	nread, err := pkt.ReadBytes(payloadLengthBuf)
	if err != nil {
		return err
	}
	if nread != 3 {
		return io.EOF
	}
	pkt.payloadLength = uint32(payloadLengthBuf[0]) | uint32(payloadLengthBuf[1])<<8 | uint32(payloadLengthBuf[2])<<16

	// Read the sequence ID
	seqIDByte, err := pkt.ReadByte()
	if err != nil {
		return err
	}
	pkt.sequenceID = SequenceID(seqIDByte)

	return nil
}

// ReadPayload reads the packet payload.
func (pkt *packet) ReadPayload() error {
	payload := make([]byte, pkt.payloadLength)
	nread, err := pkt.ReadBytes(payload)
	if err != nil {
		return err
	}
	if nread != int(pkt.payloadLength) {
		return io.EOF
	}
	pkt.payload = payload
	return nil
}

// SetPayload sets the packet payload.
func (pkt *packet) SetPayload(payload []byte) {
	pkt.payload = payload
	pkt.payloadLength = uint32(len(payload))
}

// SetSequenceID sets the packet sequence ID.
func (pkt *packet) SetSequenceID(n SequenceID) {
	pkt.sequenceID = n
}

// PayloadLength returns the packet payload length.
func (pkt *packet) PayloadLength() uint32 {
	return pkt.payloadLength
}

// SequenceID returns the packet sequence ID.
func (pkt *packet) SequenceID() SequenceID {
	return pkt.sequenceID
}

// Payload returns the packet payload.
func (pkt *packet) Payload() []byte {
	return pkt.payload
}

// SetCapabilitys sets the packet capability flags.
func (pkt *packet) SetCapability(flags Capability) {
	pkt.Capabilitys = flags
}

// Capabilitys returns the packet capability flags.
func (pkt *packet) Capability() Capability {
	return pkt.Capabilitys
}

// SetEnabled sets the specified flag.
func (pkt *packet) SetCapabilityEnabled(flag Capability) {
	pkt.Capabilitys |= flag
}

// SetDisabled unsets the specified flag.
func (pkt *packet) SetCapabilityDisabled(flag Capability) {
	pkt.Capabilitys &^= flag
}

// Reader returns the packet reader.
func (pkt *packet) Reader() *PacketReader {
	return pkt.PacketReader
}

// Bytes returns the packet bytes.
func (pkt *packet) Bytes() ([]byte, error) {
	payloadLengthBuf := []byte{
		byte(pkt.payloadLength & 0xFF),
		byte((pkt.payloadLength >> 8) & 0xFF),
		byte((pkt.payloadLength >> 16) & 0xFF),
	}
	seqIDByte := byte(pkt.sequenceID)
	return slices.Concat(payloadLengthBuf, []byte{seqIDByte}, pkt.payload), nil
}
