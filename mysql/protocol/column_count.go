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

// MySQL: Text Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html
// Result Set Packets - MariaDB Knowledge Base
// https://mariadb.com/kb/en/result-set-packets/#column-count-packet

// ColumnDefOption represents a function to set a ColumnDef option.
type ColumnCountOption func(*ColumnCount)

// ColumnDef represents a MySQL Column Definition packet.
type ColumnCount struct {
	*packet
	capFlags        Capability
	metadataFollows ResultsetMetadata
	count           uint64
}

func newColumnCountWith(pkt *packet, opts ...ColumnCountOption) *ColumnCount {
	c := &ColumnCount{
		packet:          pkt,
		capFlags:        0,
		metadataFollows: 0,
		count:           0,
	}
	pkt.SetSequenceID(1)
	c.SetOptions(opts...)
	return c
}

// WithColumnCount returns a ColumnCountOption that sets the column count.
func WithColumnCount(c uint64) ColumnCountOption {
	return func(pkt *ColumnCount) {
		pkt.count = c
	}
}

func WithColumnCountCapabilities(c Capability) ColumnCountOption {
	return func(pkt *ColumnCount) {
		pkt.capFlags = c
	}
}

// NewColumnCount returns a new ColumnCount.
func NewColumnCount(opts ...ColumnCountOption) *ColumnCount {
	pkt := newColumnCountWith(newPacket(), opts...)
	return pkt
}

func NewColumnCountFromReader(r io.Reader, opts ...ColumnCountOption) (*ColumnCount, error) {
	var err error

	pktReader, err := NewPacketHeaderWithReader(r)
	if err != nil {
		return nil, err
	}

	pkt := newColumnCountWith(pktReader, opts...)

	if pkt.Capabilities().IsEnabled(ClientOptionalResultsetMetadata) {
		pkt.metadataFollows, err = pkt.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	pkt.count, err = pkt.ReadLengthEncodedInt()
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

// SetOptions sets the options.
func (pkt *ColumnCount) SetOptions(opts ...ColumnCountOption) {
	for _, opt := range opts {
		opt(pkt)
	}
}

// Capabilities returns the capabilities.
func (pkt *ColumnCount) Capabilities() Capability {
	return pkt.capFlags
}

// MetadataFollows returns the metadata follows.
func (pkt *ColumnCount) MetadataFollows() ResultsetMetadata {
	return pkt.metadataFollows
}

// ColumnCount returns the column count.
func (pkt *ColumnCount) ColumnCount() uint64 {
	return pkt.count
}

// Bytes returns the packet bytes.
func (pkt *ColumnCount) Bytes() ([]byte, error) {
	w := NewWriter()

	if pkt.Capabilities().IsEnabled(ClientOptionalResultsetMetadata) {
		err := w.WriteByte(pkt.metadataFollows)
		if err != nil {
			return nil, err
		}
	}

	err := w.WriteLengthEncodedInt(pkt.count)
	if err != nil {
		return nil, err
	}

	pkt.SetPayload(w.Bytes())

	return pkt.packet.Bytes()
}
