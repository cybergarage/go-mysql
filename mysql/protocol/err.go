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

	sql "github.com/cybergarage/go-sqlparser/sql/errors"
)

// MySQL: ERR_Packet
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_ok_packet.html

const (
	errPacketHeader = 0xFF
)

// ErrClass represents a standard error class.
type ErrClass = sql.Class

// ErrCode represents a standard error code.
type ErrCode = sql.Code

// ERR represents a MySQL ERR packet.
type ERR struct {
	*packet
	code        uint16
	stateMarker ErrClass
	state       ErrCode
	errMsg      string
}

// ERROption represents a MySQL ERR packet option.
type ERROption func(*ERR) error

func newERRPacket(p *packet, opts ...ERROption) (*ERR, error) {
	pkt := &ERR{
		packet:      p,
		code:        0,
		stateMarker: "",
		state:       "",
		errMsg:      "",
	}
	for _, opt := range opts {
		if err := opt(pkt); err != nil {
			return nil, err
		}
	}
	return pkt, nil
}

// WithErrCode sets the error code.
func WithErrCode(code uint16) ERROption {
	return func(pkt *ERR) error {
		pkt.code = code
		return nil
	}
}

// WithStateMarker sets the state marker.
func WithStateMarker(stateMarker string) ERROption {
	return func(pkt *ERR) error {
		pkt.stateMarker = ErrClass(stateMarker)
		return nil
	}
}

// WithState sets the state.
func WithState(state string) ERROption {
	return func(pkt *ERR) error {
		pkt.state = ErrCode(state)
		return nil
	}
}

// WithErrMsg sets the error message.
func WithErrMsg(errMsg string) ERROption {
	return func(pkt *ERR) error {
		pkt.errMsg = errMsg
		return nil
	}
}

// WithErrCapability sets the error capability.
func WithErrCapability(cap CapabilityFlag) ERROption {
	return func(pkt *ERR) error {
		pkt.SetCapabilityEnabled(cap)
		return nil
	}
}

// NewERR returns a new ERR packet.
func NewERR(opts ...ERROption) (*ERR, error) {
	pkt, err := newERRPacket(nil, opts...)
	if err != nil {
		return nil, err
	}
	return pkt, nil
}

// NewERRFromReader returns a new ERR packet from the reader.
func NewERRFromReader(reader io.Reader, opts ...ERROption) (*ERR, error) {
	var err error

	pktReader, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt, err := newERRPacket(pktReader, opts...)
	if err != nil {
		return nil, err
	}

	// header
	header, err := pkt.ReadByte()
	if err != nil {
		return nil, err
	}
	if header != errPacketHeader {
		return nil, newErrInvalitHeader("ERR", header)
	}

	// error_code
	pkt.code, err = pkt.ReadInt2()
	if err != nil {
		return nil, err
	}

	// sql_state_marker, sql_state
	if pkt.CapabilityFlags().IsEnabled(ClientProtocol41) {
		// sql_state_marker
		v, err := pkt.ReadFixedLengthString(1)
		if err != nil {
			return nil, err
		}
		pkt.stateMarker = ErrClass(v)
		// sql_state
		v, err = pkt.ReadFixedLengthString(5)
		if err != nil {
			return nil, err
		}
		pkt.state = ErrCode(v)
	}

	// error_message
	pkt.errMsg, err = pkt.ReadEOFTerminatedString()
	if err != nil {
		return nil, err
	}

	return pkt, err
}

// Code returns the error code.
func (pkt *ERR) Code() uint16 {
	return pkt.code
}

// StateMarker returns the state marker.
func (pkt *ERR) StateMarker() ErrClass {
	return pkt.stateMarker
}

// State returns the state.
func (pkt *ERR) State() ErrCode {
	return pkt.state
}

// ErrMsg returns the error message.
func (pkt *ERR) ErrMsg() string {
	return pkt.errMsg
}

// Bytes returns a byte sequence of the ERR packet.
func (pkt *ERR) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	// header
	if err := w.WriteByte(errPacketHeader); err != nil {
		return nil, err
	}

	// error_code
	if err := w.WriteInt2(pkt.code); err != nil {
		return nil, err
	}

	// sql_state_marker, sql_state
	if pkt.CapabilityFlags().IsEnabled(ClientProtocol41) {
		// sql_state_marker
		if err := w.WriteFixedLengthString(string(pkt.stateMarker), 1); err != nil {
			return nil, err
		}
		// sql_state
		if err := w.WriteFixedLengthString(string(pkt.state), 5); err != nil {
			return nil, err
		}
	}

	// error_message
	if err := w.WriteEOFTerminatedString(pkt.errMsg); err != nil {
		return nil, err
	}

	pkt.packet = NewPacket(
		PacketWithSequenceID(pkt.packet.SequenceID()),
		PacketWithPayload(w.Bytes()),
	)

	return pkt.packet.Bytes()
}
