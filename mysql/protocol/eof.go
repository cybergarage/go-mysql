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

// MySQL: EOF_Packet
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_ok_packet.html

const (
	eofPacketHeader = 0xFE
)

// EOF represents a MySQL EOF packet.
type EOF struct {
	*packet
	header   uint8
	warnings uint16
	status   StatusFlag
}

// EOFOption represents a MySQL EOF packet option.
type EOFOption func(*EOF) error

func newEOFPacket(p *packet, opts ...EOFOption) (*EOF, error) {
	pkt := &EOF{
		packet:   p,
		header:   eofPacketHeader,
		warnings: 0,
		status:   0,
	}
	for _, opt := range opts {
		if err := opt(pkt); err != nil {
			return nil, err
		}
	}
	return pkt, nil
}

// WithEOFCapability returns a EOFOption that sets the capability flag.
func WithEOFCapability(c CapabilityFlag) EOFOption {
	return func(pkt *EOF) error {
		pkt.SetCapabilityEnabled(c)
		return nil
	}
}

// WithEOFWarnings returns a EOFOption that sets the number of warnings.
func WithEOFWarnings(warnings uint16) EOFOption {
	return func(pkt *EOF) error {
		pkt.warnings = warnings
		return nil
	}
}

// WithEOFStatus returns a EOFOption that sets the status flag.
func WithEOFStatus(status StatusFlag) EOFOption {
	return func(pkt *EOF) error {
		pkt.status = status
		return nil
	}
}

// NewEOF returns a new EOF packet.
func NewEOF(opts ...EOFOption) (*EOF, error) {
	return newEOFPacket(nil, opts...)
}

// NewEOFFromReader returns a new EOF packet from the reader.
func NewEOFFromReader(reader io.Reader, opts ...EOFOption) (*EOF, error) {
	var err error

	pktReader, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt, err := newEOFPacket(pktReader, opts...)
	if err != nil {
		return nil, err
	}

	// header
	pkt.header, err = pkt.ReadByte()
	if err != nil {
		return nil, err
	}
	if pkt.header != eofPacketHeader {
		return nil, newErrInvalidHeader("EOF", pkt.header)
	}

	// warnings and status flags
	if pkt.CapabilityFlags().IsEnabled(ClientProtocol41) {
		// warnings
		pkt.warnings, err = pkt.ReadInt2()
		if err != nil {
			return nil, err
		}
		// status flags
		v, err := pkt.ReadInt2()
		if err != nil {
			return nil, err
		}
		pkt.status = StatusFlag(v)
	}

	return pkt, err
}

// Header returns the header.
func (pkt *EOF) Header() uint8 {
	return pkt.header
}

// Warnings returns the number of warnings.
func (pkt *EOF) Warnings() uint16 {
	return pkt.warnings
}

// Status returns the status flag.
func (pkt *EOF) Status() StatusFlag {
	return pkt.status
}

// Bytes returns a byte sequence of the EOF packet.
func (pkt *EOF) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	// header
	if err := w.WriteByte(pkt.header); err != nil {
		return nil, err
	}

	// warnings and status flags
	if pkt.CapabilityFlags().IsEnabled(ClientProtocol41) {
		if err := w.WriteInt2(pkt.warnings); err != nil {
			return nil, err
		}
		if err := w.WriteInt2(uint16(pkt.status)); err != nil {
			return nil, err
		}
	}

	pkt.packet = NewPacket(
		PacketWithSequenceID(pkt.packet.SequenceID()),
		PacketWithPayload(w.Bytes()),
	)

	return pkt.packet.Bytes()
}
