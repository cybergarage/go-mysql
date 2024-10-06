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
type OKOption func(*OK) error

// WithOKCapability returns a OKOption that sets the capability flag.
func WithOKCapability(c CapabilityFlag) OKOption {
	return func(pkt *OK) error {
		pkt.SetCapabilityEnabled(c)
		return nil
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
	if (pkt.header != okPacketHeader) && (pkt.header != errPacketHeader) {
		return nil, newErrInvalitHeader("OK", pkt.header)
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

	if pkt.CapabilityFlags().IsEnabled(ClientProtocol41) {
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
	} else if pkt.CapabilityFlags().IsEnabled(ClientTransactions) {
		// status
		pkt.status, err = pkt.ReadInt2()
		if err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientSessionTrack) {
		// info
		pkt.info, err = pkt.ReadLengthEncodedString()
		if err != nil {
			return nil, err
		}
		if pkt.Status().IsEnabled(StatusSessionStateChanged) {
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

// Status returns the status flag.
func (pkt *OK) Status() StatusFlag {
	return StatusFlag(pkt.status)
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
	if err := w.WriteByte(errPacketHeader); err != nil {
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

	if pkt.CapabilityFlags().IsEnabled(ClientProtocol41) {
		// status
		if err := w.WriteInt2(uint16(pkt.status)); err != nil {
			return nil, err
		}
		// warnings
		if err := w.WriteInt2(pkt.warnings); err != nil {
			return nil, err
		}
	} else if pkt.CapabilityFlags().IsEnabled(ClientTransactions) {
		// status
		if err := w.WriteInt2(pkt.status); err != nil {
			return nil, err
		}
	}

	if pkt.CapabilityFlags().IsEnabled(ClientSessionTrack) {
		// info
		if err := w.WriteLengthEncodedString(pkt.info); err != nil {
			return nil, err
		}
		if pkt.Status().IsEnabled(StatusSessionStateChanged) {
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

	pkt.packet = NewPacket(
		PacketWithSequenceID(pkt.packet.SequenceID()),
		PacketWithPayload(w.Bytes()),
	)

	return pkt.packet.Bytes()
}
