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

// MySQL: Protocol::QueryResponse
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query.html
// MySQL: Text Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html
// MySQL: Text Resultset Row
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset_row.html

// TextResultSetRowOption represents a MySQL text resultset row option.
type TextResultSetRowOption func(*TextResultSetRow)

// TextResultSetRow represents a MySQL text resultset row response packet.
type TextResultSetRow struct {
	*packet
	columns []string
}

func newTextResultSetRowWithPacket(pkt *packet, opts ...TextResultSetRowOption) *TextResultSetRow {
	row := &TextResultSetRow{
		packet:  pkt,
		columns: []string{},
	}
	row.SetOptions(opts...)
	return row
}

func WithTextResultSetRowColmunCount(c uint64) TextResultSetRowOption {
	return func(pkt *TextResultSetRow) {
		pkt.columns = make([]string, c)
	}
}

// NewTextResultSetRow returns a new TextResultSetRow.
func NewTextResultSetRow(opts ...TextResultSetRowOption) *TextResultSetRow {
	return newTextResultSetRowWithPacket(nil, opts...)
}

// NewTextResultSetRowFromReader returns a new TextResultSetRow from the reader.
func NewTextResultSetRowFromReader(reader *PacketReader, opts ...TextResultSetRowOption) (*TextResultSetRow, error) {
	pktReader, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	row := newTextResultSetRowWithPacket(pktReader, opts...)
	columnCount := len(row.columns)
	for n := 0; n < columnCount; n++ {
		column, err := row.ReadColumn()
		if err != nil {
			return nil, err
		}
		row.columns = append(row.columns, column)
	}

	return row, nil
}

// SetOptions sets the options.
func (row *TextResultSetRow) SetOptions(opts ...TextResultSetRowOption) {
	for _, opt := range opts {
		opt(row)
	}
}

// ReadColumn reads a column.
func (row *TextResultSetRow) ReadColumn() (string, error) {
	return row.ReadLengthEncodedString()
}

// Columns returns the columns.
func (row *TextResultSetRow) Columns() []string {
	return row.columns
}

// Bytes returns the packet bytes.
func (row *TextResultSetRow) Bytes() ([]byte, error) {
	w := NewPacketWriter()
	for _, column := range row.columns {
		err := w.WriteLengthEncodedString(column)
		if err != nil {
			return nil, err
		}
	}
	return w.Bytes(), nil
}
