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

// MySQL: Protocol::QueryResponse
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query.html

// QueryResponse represents a MySQL COM_QUERY response packet.
type QueryResponse struct {
	*packet
	capFlags CapabilityFlag
}

func newQueryResponseWithPacket(pkt *packet, opts ...QueryResponseOption) *QueryResponse {
	q := &QueryResponse{
		packet:   pkt,
		capFlags: 0,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// QueryResponseOption represents a COM_QUERY response option.
type QueryResponseOption func(*QueryResponse)

func WithQueryResponseCapabilities(c CapabilityFlag) QueryResponseOption {
	return func(pkt *QueryResponse) {
		pkt.capFlags = c
	}
}

// NewQueryResponse returns a new COM_QUERY response packet.
func NewQueryResponse(opts ...QueryResponseOption) (*QueryResponse, error) {
	pkt := newQueryResponseWithPacket(nil, opts...)
	return pkt, nil
}

// NewQueryResponseFromReader returns a new COM_QUERY response packet from the specified reader.
func NewQueryResponseFromReader(reader io.Reader, opts ...QueryResponseOption) (*QueryResponse, error) {
	var err error

	pktReader, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newQueryResponseWithPacket(pktReader, opts...)

	return pkt, nil
}

// Capabilities returns the capabilities.
func (pkt *QueryResponse) Capabilities() CapabilityFlag {
	return pkt.capFlags
}

// Bytes returns the packet bytes.
func (pkt *QueryResponse) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	pkt.packet =
		NewPacket(
			PacketWithSequenceID(pkt.SequenceID()),
			PacketWithPayload(w.Bytes()),
		)

	return pkt.Bytes()
}
