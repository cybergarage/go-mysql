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

// HandshakeResponse represents a MySQL Handshake Response message.
type HandshakeResponse struct {
	*message
}

func newHandshakeResponseWithMessage(msg *message) *HandshakeResponse {
	return &HandshakeResponse{
		message: msg,
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

	return res, err
}

// Bytes returns the message bytes.
func (res *HandshakeResponse) Bytes() ([]byte, error) {
	w := NewWriter()

	res.message = NewMessageWithPayload(w.Bytes())
	return res.message.Bytes()
}
