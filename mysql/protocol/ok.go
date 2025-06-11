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
	header           uint8
	affectedRows     uint64
	lastInsertID     uint64
	status           uint16
	warnings         uint16
	info             string
	sessionStateInfo string
}

// OKOption represents a MySQL OK packet option.
type OKOption func(*OK)

// WithOKSecuenceID returns a OKOption that sets the sequence ID.
func WithOKSecuenceID(n SequenceID) OKOption {
	return func(pkt *OK) {
		pkt.SetSequenceID(n)
	}
}

// WithOKCapability returns a OKOption that sets the capability flag.
func WithOKCapability(c Capability) OKOption {
	return func(pkt *OK) {
		pkt.SetCapabilityEnabled(c)
	}
}

// WithOKAffectedRows returns a OKOption that sets the number of affected rows.
func WithOKAffectedRows(v uint64) OKOption {
	return func(pkt *OK) {
		pkt.affectedRows = v
	}
}

// WithOKLastInsertID returns a OKOption that sets the last insert ID.
func WithOKLastInsertID(v uint64) OKOption {
	return func(pkt *OK) {
		pkt.lastInsertID = v
	}
}

// WithOKServerStatus returns a OKOption that sets the server status flag.
func WithOKServerStatus(v ServerStatus) OKOption {
	return func(pkt *OK) {
		pkt.SetServerStatus(v)
	}
}

// WithOKWarnings returns a OKOption that sets the number of warnings.
func WithOKWarnings(v uint16) OKOption {
	return func(pkt *OK) {
		pkt.warnings = v
	}
}

func newOKPacket(p *packet, opts ...OKOption) (*OK, error) {
	pkt := &OK{
		packet:           p,
		header:           0,
		affectedRows:     0,
		lastInsertID:     0,
		status:           0,
		warnings:         0,
		info:             "",
		sessionStateInfo: "",
	}
	for _, opt := range opts {
		opt(pkt)
	}
	return pkt, nil
}

// NewOK returns a new OK packet.
func NewOK(opts ...OKOption) (*OK, error) {
	return newOKPacket(newPacket(), opts...)
}

// NewOKFromReader returns a new OK packet from the reader.
func NewOKFromReader(reader io.Reader, opts ...OKOption) (*OK, error) {
	var err error

	pktReader, err := NewPacketHeaderWithReader(reader)
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
	if (pkt.header != okPacketHeader) && (pkt.header != errPacketHeader) {
		return nil, newErrInvalidHeader("OK", pkt.header)
	}

	// affectedRows
	pkt.affectedRows, err = pkt.ReadLengthEncodedInt()
	if err != nil {
		return nil, err
	}

	// lastInsertID
	pkt.lastInsertID, err = pkt.ReadLengthEncodedInt()
	if err != nil {
		return nil, err
	}

	if pkt.Capability().HasCapability(ClientProtocol41) {
		// status
		v, err := pkt.ReadInt2()
		if err != nil {
			return nil, err
		}
		pkt.status = uint16(v)
		// warnings
		pkt.warnings, err = pkt.ReadInt2()
		if err != nil {
			return nil, err
		}
	} else if pkt.Capability().HasCapability(ClientTransactions) {
		// status
		pkt.status, err = pkt.ReadInt2()
		if err != nil {
			return nil, err
		}
	}

	if pkt.Capability().HasCapability(ClientSessionTrack) {
		// info
		pkt.info, err = pkt.ReadLengthEncodedString()
		if err != nil {
			return nil, err
		}
		if pkt.ServerStatus().IsEnabled(ServerSessionStateChanged) {
			// sessionStateInfo
			pkt.sessionStateInfo, err = pkt.ReadLengthEncodedString()
			if err != nil {
				return nil, err
			}
		}
	} else {
		// info
		pkt.info, err = pkt.ReadEOFTerminatedString()
		if err != nil {
			return nil, err
		}
	}

	return pkt, err
}

// Header returns the header.
func (pkt *OK) Header() uint8 {
	return pkt.header
}

// OK returns true if the packet is an OK packet.
func (pkt *OK) OK() bool {
	return pkt.header == okPacketHeader
}

// Err returns true if the packet is an ERR packet.
func (pkt *OK) Err() bool {
	return pkt.header == errPacketHeader
}

// AffectedRows returns the number of affected rows.
func (pkt *OK) AffectedRows() uint64 {
	return pkt.affectedRows
}

// LastInsertID returns the last insert ID.
func (pkt *OK) LastInsertID() uint64 {
	return pkt.lastInsertID
}

// Warnings returns the number of warnings.
func (pkt *OK) Warnings() uint16 {
	return pkt.warnings
}

// Info returns the info.
func (pkt *OK) Info() string {
	return pkt.info
}

// SessionStateInfo returns the session state info.
func (pkt *OK) SessionStateInfo() string {
	return pkt.sessionStateInfo
}

// Bytes returns a byte sequence of the OK packet.
func (pkt *OK) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	// header
	if err := w.WriteByte(pkt.header); err != nil {
		return nil, err
	}

	// affectedRows
	if err := w.WriteLengthEncodedInt(pkt.affectedRows); err != nil {
		return nil, err
	}

	// lastInsertID
	if err := w.WriteLengthEncodedInt(pkt.lastInsertID); err != nil {
		return nil, err
	}

	if pkt.Capability().HasCapability(ClientProtocol41) {
		// status
		if err := w.WriteInt2(uint16(pkt.status)); err != nil {
			return nil, err
		}
		// warnings
		if err := w.WriteInt2(pkt.warnings); err != nil {
			return nil, err
		}
	} else if pkt.Capability().HasCapability(ClientTransactions) {
		// status
		if err := w.WriteInt2(pkt.status); err != nil {
			return nil, err
		}
	}

	if pkt.Capability().HasCapability(ClientSessionTrack) {
		// info
		if err := w.WriteLengthEncodedString(pkt.info); err != nil {
			return nil, err
		}
		if pkt.ServerStatus().IsEnabled(ServerSessionStateChanged) {
			// sessionStateInfo
			if err := w.WriteLengthEncodedString(pkt.sessionStateInfo); err != nil {
				return nil, err
			}
		}
	} else {
		// info
		if err := w.WriteEOFTerminatedString(pkt.info); err != nil {
			return nil, err
		}
	}

	pkt.SetPayload(w.Bytes())

	return pkt.packet.Bytes()
}
