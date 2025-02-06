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
	"bytes"
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
	capFlags   Capability
	columnCnt  *ColumnCount
	columnDefs []ColumnDef
	rows       []ResultSetRow
}

func newBinaryResultSetWithPacket(opts ...BinaryResultSetOption) *BinaryResultSet {
	q := &BinaryResultSet{
		capFlags:   0,
		columnCnt:  NewColumnCount(),
		columnDefs: []ColumnDef{},
		rows:       []ResultSetRow{},
	}
	q.SetOptions(opts...)
	return q
}

// BinaryResultSetOption represents a COM_QUERY binary resultset option.
type BinaryResultSetOption func(*BinaryResultSet)

// WithBinaryResultSetCapabilities returns a binary resultset option to set the capabilities.
func WithBinaryResultSetCapability(c Capability) BinaryResultSetOption {
	return func(pkt *BinaryResultSet) {
		pkt.capFlags = c
	}
}

// WithBinaryResultSetMetadataFollows returns a binary resultset option to set the metadata follows.
func WithBinaryResultSetMetadataFollows(m ResultsetMetadata) BinaryResultSetOption {
	return func(pkt *BinaryResultSet) {
		pkt.columnCnt.metadataFollows = m
	}
}

// WithBinaryResultSetColumnDefs returns a binary resultset option to set the column definitions.
func WithBinaryResultSetColumnDefs(colDefs []ColumnDef) BinaryResultSetOption {
	return func(pkt *BinaryResultSet) {
		pkt.columnCnt.count = uint64(len(colDefs))
		pkt.columnDefs = colDefs
	}
}

// WithBinaryResultSetRows returns a binary resultset option to set the rows.
func WithBinaryResultSetRows(rows []ResultSetRow) BinaryResultSetOption {
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
	var err error

	pkt := newBinaryResultSetWithPacket(opts...)

	columnCountOpts := []ColumnCountOption{
		WithColumnCountCapability(pkt.Capability()),
	}
	pkt.columnCnt, err = NewColumnCountFromReader(reader, columnCountOpts...)
	if err != nil {
		return nil, err
	}

	columnCount := pkt.columnCnt.ColumnCount()

	if pkt.Capability().IsDisabled(ClientOptionalResultsetMetadata) || pkt.columnCnt.MetadataFollows() == ResultsetMetadataFull {
		for i := 0; i < int(columnCount); i++ {
			colDef, err := NewColumnDefFromReader(reader)
			if err != nil {
				return nil, err
			}
			pkt.columnDefs = append(pkt.columnDefs, colDef)
		}
	}

	if pkt.Capability().IsDisabled(ClientDeprecateEOF) {
		_, err := NewEOFFromReader(reader, WithEOFCapability(pkt.Capability()))
		if err != nil {
			return nil, err
		}
	}

	// One or more Binary Resultset Row

	rowPkt, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	for !rowPkt.IsEOF() {
		rowPktBytes, err := rowPkt.Bytes()
		if err != nil {
			return nil, err
		}
		rowPktReader := NewPacketReaderWith(bytes.NewReader(rowPktBytes))
		row, err := NewBinaryResultSetRowFromReader(rowPktReader, WithBinaryResultSetRowColmunCount(columnCount))
		if err != nil {
			return nil, err
		}
		pkt.rows = append(pkt.rows, row)

		rowPkt, err = NewPacketWithReader(reader)
		if err != nil {
			return nil, err
		}
	}

	return pkt, nil
}

// SetOptions sets the options.
func (pkt *BinaryResultSet) SetOptions(opts ...BinaryResultSetOption) {
	for _, opt := range opts {
		opt(pkt)
	}
}

// SetCapability sets a capability flag.
func (pkt *BinaryResultSet) SetCapability(c Capability) {
	pkt.capFlags = c
}

// SetSequenceID sets the packet sequence ID.
func (pkt *BinaryResultSet) SetSequenceID(n SequenceID) {
	pkt.columnCnt.SetSequenceID(n)
	for _, colDef := range pkt.columnDefs {
		n = n.Next()
		colDef.SetSequenceID(n)
	}
	for _, row := range pkt.rows {
		n = n.Next()
		row.SetSequenceID(n)
	}
}

// Capability returns the capabilities.
func (pkt *BinaryResultSet) Capability() Capability {
	return pkt.capFlags
}

// Rows returns the rows.
func (pkt *BinaryResultSet) Rows() []ResultSetRow {
	return pkt.rows
}

// Bytes returns the packet bytes.
func (pkt *BinaryResultSet) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	columCntBytes, err := pkt.columnCnt.Bytes()
	if err != nil {
		return nil, err
	}
	_, err = w.WriteBytes(columCntBytes)
	if err != nil {
		return nil, err
	}

	secuenceID := pkt.columnCnt.SequenceID()
	secuenceID = secuenceID.Next()

	if pkt.Capability().IsDisabled(ClientOptionalResultsetMetadata) || pkt.columnCnt.MetadataFollows() == ResultsetMetadataFull {
		for _, colDef := range pkt.columnDefs {
			colDef.SetSequenceID(secuenceID)
			err := w.WritePacket(colDef)
			if err != nil {
				return nil, err
			}
			secuenceID = secuenceID.Next()
		}
	}

	if pkt.Capability().IsDisabled(ClientDeprecateEOF) {
		err := w.WriteEOF(secuenceID, pkt.Capability())
		if err != nil {
			return nil, err
		}
		secuenceID = secuenceID.Next()
	}

	// One or more Binary Resultset Row

	for _, row := range pkt.rows {
		row.SetSequenceID(secuenceID)
		rowBytes, err := row.Bytes()
		if err != nil {
			return nil, err
		}
		_, err = w.WriteBytes(rowBytes)
		if err != nil {
			return nil, err
		}
		secuenceID = secuenceID.Next()
	}

	if pkt.Capability().IsEnabled(ClientDeprecateEOF) {
		err := w.WriteOK(secuenceID, pkt.Capability())
		if err != nil {
			return nil, err
		}
	} else {
		err := w.WriteEOF(secuenceID, pkt.Capability())
		if err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}
