// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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

// MySQL: Protocol Basics
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basics.html
// MySQL: MySQL Packets
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_basic_packets.html
// MySQL: Messages
// https://dev.mysql.com/doc/dev/mysql-server/latest/mysqlx_protocol_messages.html
// MariaDB protocol difference with MySQL - MariaDB Knowledge Base
// https://mariadb.com/kb/en/mariadb-protocol-difference-with-mysql/

// Message represents a MySQL message.
type Message interface {
	// SetSequenceID sets the message sequence ID.
	SetSequenceID(n SequenceID)
	// PayloadLength returns the message payload length.
	PayloadLength() uint32
	// SequenceID returns the message sequence ID.
	SequenceID() SequenceID
	// Payload returns the message payload.
	Payload() []byte
	// Bytes returns the message bytes.
	Bytes() ([]byte, error)
}

// SequenceID represents a MySQL message sequence ID.
type SequenceID uint8

// message represents a MySQL message.
type message struct {
	*Reader
	payloadLength uint32
	sequenceID    SequenceID
	payload       []byte
}

func newMessage() *message {
	return &message{
		Reader:        nil,
		payloadLength: 0,
		sequenceID:    SequenceID(0),
		payload:       nil,
	}
}

// NewMessage returns a new MySQL message.
func NewMessageWithPayload(payload []byte) *message {
	msg := newMessage()
	msg.payloadLength = uint32(len(payload))
	msg.payload = payload
	return msg
}

// NewMessage returns a new MySQL message.
func NewMessageWithReader(reader io.Reader) (*message, error) {
	msg := newMessage()
	msg.Reader = NewReaderWith(reader)

	// Read the payload length

	payloadLengthBuf := make([]byte, 3)
	_, err := msg.ReadBytes(payloadLengthBuf)
	if err != nil {
		return nil, err
	}
	msg.payloadLength = uint32(payloadLengthBuf[0]) | uint32(payloadLengthBuf[1])<<8 | uint32(payloadLengthBuf[2])<<16

	// Read the sequence ID
	seqIDByte, err := msg.ReadByte()
	if err != nil {
		return nil, err
	}
	msg.sequenceID = SequenceID(seqIDByte)

	return msg, nil
}

// SetSequenceID sets the message sequence ID.
func (msg *message) SetSequenceID(n SequenceID) {
	msg.sequenceID = n
}

// PayloadLength returns the message payload length.
func (msg *message) PayloadLength() uint32 {
	return msg.payloadLength
}

// SequenceID returns the message sequence ID.
func (msg *message) SequenceID() SequenceID {
	return msg.sequenceID
}

// Payload returns the message payload.
func (msg *message) Payload() []byte {
	return msg.payload
}

// Bytes returns the message bytes.
func (msg *message) Bytes() ([]byte, error) {
	payloadLengthBuf := []byte{
		byte(msg.payloadLength),
		byte(msg.payloadLength >> 8),
		byte(msg.payloadLength >> 16),
	}
	seqIDByte := byte(msg.sequenceID)
	return append(append(payloadLengthBuf, seqIDByte), msg.payload...), nil
}
