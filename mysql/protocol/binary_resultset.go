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
	*packet
	columnDefs []ColumnDef
	rows       []BinaryResultSetRow
}

func newBinaryResultSetWithPacket(pkt *packet, opts ...BinaryResultSetOption) *BinaryResultSet {
	q := &BinaryResultSet{
		packet:     pkt,
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
		pkt.SetCapability(c)
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
	pkt := newBinaryResultSetWithPacket(newPacket(), opts...)
	return pkt, nil
}

// NewBinaryResultSetFromReader returns a new binary resultset response packet from the specified reader.
func NewBinaryResultSetFromReader(reader io.Reader, opts ...BinaryResultSetOption) (*BinaryResultSet, error) {
	var err error

	pktHeader, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newBinaryResultSetWithPacket(pktHeader, opts...)

	columnCount, err := pkt.ReadLengthEncodedInt()
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

	numPeekBytes := 5
	nextBytes, err := pkt.PeekBytes(numPeekBytes)
	if err != nil {
		return nil, err
	}
	pkt.rows = []BinaryResultSetRow{}
	for nextBytes[4] != 0xFE {
		row, err := NewBinaryResultSetRowFromReader(pkt.Reader(),
			WithBinaryResultSetRowColumnDefs(pkt.columnDefs))
		if err != nil {
			return nil, err
		}
		pkt.rows = append(pkt.rows, *row)
		nextBytes, err = pkt.PeekBytes(numPeekBytes)
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

// Rows returns the rows.
func (pkt *BinaryResultSet) Rows() []BinaryResultSetRow {
	return pkt.rows
}

// Bytes returns the packet bytes.
func (pkt *BinaryResultSet) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	_, err := w.WriteBytes(pkt.HeaderBytes())
	if err != nil {
		return nil, err
	}

	err = w.WriteLengthEncodedInt(uint64(len(pkt.columnDefs)))
	if err != nil {
		return nil, err
	}

	seqID := pkt.SequenceID()
	for _, colDef := range pkt.columnDefs {
		seqID = seqID.Next()
		colDef.SetSequenceID(seqID)
		bytes, err := colDef.Bytes()
		if err != nil {
			return nil, err
		}
		_, err = w.WriteBytes(bytes)
		if err != nil {
			return nil, err
		}
	}

	for _, row := range pkt.rows {
		seqID = seqID.Next()
		row.SetSequenceID(seqID)
		bytes, err := row.Bytes()
		if err != nil {
			return nil, err
		}
		_, err = w.WriteBytes(bytes)
		if err != nil {
			return nil, err
		}
	}

	seqID = seqID.Next()
	if err := w.WriteEOF(pkt.Capability(), seqID); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
