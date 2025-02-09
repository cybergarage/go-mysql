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
// MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// Result Set Packets - MariaDB Knowledge Base
// https://mariadb.com/kb/en/result-set-packets/

// BinaryResultSet represents a MySQL binary resultset response packet.
type BinaryResultSet struct {
	capability Capability
	columnDefs []ColumnDef
	rows       []BinaryResultSetRow
}

func newBinaryResultSetWithPacket(opts ...BinaryResultSetOption) *BinaryResultSet {
	q := &BinaryResultSet{
		capability: 0,
		columnDefs: []ColumnDef{},
		rows:       []BinaryResultSetRow{},
	}
	q.SetOptions(opts...)
	return q
}

// BinaryResultSetOption represents a COM_QUERY binary resultset option.
type BinaryResultSetOption func(*BinaryResultSet)

// WithBinaryResultSetCapabilities returns a binary resultset option to set the capabilities.
func WithBinaryResultSetCapability(c Capability) BinaryResultSetOption {
	return func(pkt *BinaryResultSet) {
		pkt.capability = c
	}
}

// WithBinaryResultSetColumnDefs returns a binary resultset option to set the column definitions.
func WithBinaryResultSetColumnDefs(colDefs []ColumnDef) BinaryResultSetOption {
	return func(pkt *BinaryResultSet) {
		pkt.columnDefs = colDefs
	}
}

// WithBinaryResultSetRows returns a binary resultset option to set the rows.
func WithBinaryResultSetRows(rows []BinaryResultSetRow) BinaryResultSetOption {
	return func(pkt *BinaryResultSet) {
		pkt.rows = rows
	}
}

// NewBinaryResultSet returns a new binary resultset response packet.
func NewBinaryResultSet(opts ...BinaryResultSetOption) (*BinaryResultSet, error) {
	pkt := newBinaryResultSetWithPacket(opts...)
	return pkt, nil
}

// NewBinaryResultSetFromReader returns a new binary resultset response packet from the specified reader.
func NewBinaryResultSetFromReader(reader io.Reader, opts ...BinaryResultSetOption) (*BinaryResultSet, error) {
	pkt := newBinaryResultSetWithPacket(opts...)
	pktReader := NewReaderWithReader(reader)

	columnCount, err := pktReader.ReadLengthEncodedInt()
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(columnCount); i++ {
		colDef, err := NewColumnDefFromReader(reader)
		if err != nil {
			return nil, err
		}
		pkt.columnDefs = append(pkt.columnDefs, colDef)
	}

	nextByte, err := pktReader.PeekByte()
	if err != nil {
		return nil, err
	}

	for nextByte != 0xFE {
		nextByte, err = pktReader.PeekByte()
		if err != nil {
			return nil, err
		}
	}

	_, err = NewEOFFromReader(reader, WithEOFCapability(pkt.Capability()))
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

// SetOptions sets the options.
func (pkt *BinaryResultSet) SetOptions(opts ...BinaryResultSetOption) {
	for _, opt := range opts {
		opt(pkt)
	}
}

// SetSequenceID sets the packet sequence ID.
func (pkt *BinaryResultSet) SetSequenceID(n SequenceID) {
}

// SetCapability sets a capability flag.
func (pkt *BinaryResultSet) SetCapability(c Capability) {
	pkt.capability = c
}

// Capability returns the capabilities.
func (pkt *BinaryResultSet) Capability() Capability {
	return pkt.capability
}

// Rows returns the rows.
func (pkt *BinaryResultSet) Rows() []BinaryResultSetRow {
	return pkt.rows
}

// Bytes returns the packet bytes.
func (pkt *BinaryResultSet) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	return w.Bytes(), nil
}
