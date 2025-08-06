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

// WithBinaryResultSetServerStatus returns a binary resultset option to set the server status.
func WithBinaryResultSetServerStatus(s ServerStatus) BinaryResultSetOption {
	return func(pkt *BinaryResultSet) {
		pkt.SetServerStatus(s)
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

	rsPkt := newBinaryResultSetWithPacket(pktHeader, opts...)

	columnCount, err := rsPkt.ReadLengthEncodedInt()
	if err != nil {
		return nil, err
	}

	// Column Definitions

	for range columnCount {
		colDef, err := NewColumnDefFromReader(reader)
		if err != nil {
			return nil, err
		}

		rsPkt.columnDefs = append(rsPkt.columnDefs, colDef)
	}

	// EOF

	_, err = NewEOFFromReader(reader, WithEOFCapability(rsPkt.Capability()))
	if err != nil {
		return nil, err
	}

	// Rows or EOF

	for {
		pkt, err := NewPacketWithReader(reader)
		if err != nil {
			return nil, err
		}

		pktHeader, err := pkt.PayloadHeaderByte()
		if err != nil {
			return nil, err
		}

		if pktHeader == 0xFE {
			break
		}

		pktBytes, err := pkt.Bytes()
		if err != nil {
			return nil, err
		}

		row, err := NewBinaryResultSetRowFromReader(
			NewPacketReaderWithBytes(pktBytes),
			WithBinaryResultSetRowColumnDefs(rsPkt.columnDefs))
		if err != nil {
			return nil, err
		}

		rsPkt.rows = append(rsPkt.rows, *row)
	}

	return rsPkt, nil
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

	// First Packet

	seqID := (SequenceID)(1)
	pkt.SetSequenceID(seqID)

	firstPktWriter := NewPacketWriter()

	err := firstPktWriter.WriteLengthEncodedInt(uint64(len(pkt.columnDefs)))
	if err != nil {
		return nil, err
	}

	pkt.SetPayload(firstPktWriter.Bytes())

	firstPktBytes, err := pkt.packet.Bytes()
	if err != nil {
		return nil, err
	}

	_, err = w.WriteBytes(firstPktBytes)
	if err != nil {
		return nil, err
	}

	// column_count * Column Definition

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

	// EOF

	seqID = seqID.Next()
	if err := w.WriteEOF(pkt.Capability(), pkt.ServerStatus(), seqID); err != nil {
		return nil, err
	}

	// None or many Binary Protocol Resultset Row

	for _, row := range pkt.rows {
		seqID = seqID.Next()
		row.SetSequenceID(seqID)

		rowBytes, err := row.Bytes()
		if err != nil {
			return nil, err
		}

		_, err = w.WriteBytes(rowBytes)
		if err != nil {
			return nil, err
		}
	}

	// EOF

	seqID = seqID.Next()
	if err := w.WriteEOF(pkt.Capability(), pkt.ServerStatus(), seqID); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}
