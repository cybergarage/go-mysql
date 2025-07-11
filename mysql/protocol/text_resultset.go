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
// MySQL: Text Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html
// Result Set Packets - MariaDB Knowledge Base
// https://mariadb.com/kb/en/result-set-packets/

// TextResultSet represents a MySQL text resultset response packet.
type TextResultSet struct {
	capFlags   Capability
	columnCnt  *ColumnCount
	columnDefs []ColumnDef
	rows       []ResultSetRow
}

// TextResultSetOption represents a COM_QUERY text resultset option.
type TextResultSetOption func(*TextResultSet)

// WithTextResultSetCapabilities returns a text resultset option to set the capabilities.
func WithTextResultSetCapability(c Capability) TextResultSetOption {
	return func(pkt *TextResultSet) {
		pkt.capFlags = c
	}
}

// WithTextResultSetMetadataFollows returns a text resultset option to set the metadata follows.
func WithTextResultSetMetadataFollows(m ResultsetMetadata) TextResultSetOption {
	return func(pkt *TextResultSet) {
		pkt.columnCnt.metadataFollows = m
	}
}

// WithTextResultSetColumnDefs returns a text resultset option to set the column definitions.
func WithTextResultSetColumnDefs(colDefs []ColumnDef) TextResultSetOption {
	return func(pkt *TextResultSet) {
		pkt.columnCnt.count = uint64(len(colDefs))
		pkt.columnDefs = colDefs
	}
}

// WithTextResultSetRows returns a text resultset option to set the rows.
func WithTextResultSetRows(rows []ResultSetRow) TextResultSetOption {
	return func(pkt *TextResultSet) {
		pkt.rows = rows
	}
}

// NewTextResultSet returns a new text resultset response packet.
func NewTextResultSet(opts ...TextResultSetOption) (*TextResultSet, error) {
	pkt := &TextResultSet{
		capFlags:   0,
		columnCnt:  NewColumnCount(),
		columnDefs: []ColumnDef{},
		rows:       []ResultSetRow{},
	}
	pkt.SetOptions(opts...)
	return pkt, nil
}

// NewTextResultSetFromReader returns a new text resultset response packet from the specified reader.
func NewTextResultSetFromReader(reader io.Reader, opts ...TextResultSetOption) (*TextResultSet, error) {
	var err error

	pkt, err := NewTextResultSet(opts...)
	if err != nil {
		return nil, err
	}

	columnCountOpts := []ColumnCountOption{
		WithColumnCountCapability(pkt.Capability()),
	}
	pkt.columnCnt, err = NewColumnCountFromReader(reader, columnCountOpts...)
	if err != nil {
		return nil, err
	}

	columnCount := pkt.columnCnt.ColumnCount()

	if pkt.Capability().LacksCapability(ClientOptionalResultsetMetadata) || pkt.columnCnt.MetadataFollows() == ResultsetMetadataFull {
		for i := 0; i < int(columnCount); i++ {
			colDef, err := NewColumnDefFromReader(reader)
			if err != nil {
				return nil, err
			}
			pkt.columnDefs = append(pkt.columnDefs, colDef)
		}
	}

	if pkt.Capability().LacksCapability(ClientDeprecateEOF) {
		_, err := NewEOFFromReader(reader, WithEOFCapability(pkt.Capability()))
		if err != nil {
			return nil, err
		}
	}

	// One or more Text Resultset Row

	rowPkt, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	for !rowPkt.IsEOF() {
		rowPktBytes, err := rowPkt.Bytes()
		if err != nil {
			return nil, err
		}
		rowPktReader := NewPacketReaderWithReader(bytes.NewReader(rowPktBytes))
		row, err := NewTextResultSetRowFromReader(rowPktReader, WithTextResultSetRowColmunCount(columnCount))
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
func (pkt *TextResultSet) SetOptions(opts ...TextResultSetOption) {
	for _, opt := range opts {
		opt(pkt)
	}
}

// SetCapability sets a capability flag.
func (pkt *TextResultSet) SetCapability(c Capability) {
	pkt.capFlags = c
}

// SetSequenceID sets the packet sequence ID.
func (pkt *TextResultSet) SetSequenceID(n SequenceID) {
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
func (pkt *TextResultSet) Capability() Capability {
	return pkt.capFlags
}

// ColumnCount returns the column count.
func (pkt *TextResultSet) ColumnCount() uint64 {
	return pkt.columnCnt.ColumnCount()
}

// ColumnDefs returns the column definitions.
func (pkt *TextResultSet) ColumnDefs() []ColumnDef {
	return pkt.columnDefs
}

// Rows returns the rows.
func (pkt *TextResultSet) Rows() []ResultSetRow {
	return pkt.rows
}

// Bytes returns the packet bytes.
func (pkt *TextResultSet) Bytes() ([]byte, error) {
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

	if pkt.Capability().LacksCapability(ClientOptionalResultsetMetadata) || pkt.columnCnt.MetadataFollows() == ResultsetMetadataFull {
		for _, colDef := range pkt.columnDefs {
			colDef.SetSequenceID(secuenceID)
			err := w.WritePacket(colDef)
			if err != nil {
				return nil, err
			}
			secuenceID = secuenceID.Next()
		}
	}

	if pkt.Capability().LacksCapability(ClientDeprecateEOF) {
		err := w.WriteEOF(secuenceID, pkt.Capability())
		if err != nil {
			return nil, err
		}
		secuenceID = secuenceID.Next()
	}

	// One or more Text Resultset Row

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

	if pkt.Capability().HasCapability(ClientDeprecateEOF) {
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
