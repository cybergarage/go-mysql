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
// MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// MySQL: Binary Protocol Resultset Row
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/

// BinaryResultSetRowOption represents a MySQL binary resultset row option.
type BinaryResultSetRowOption func(*BinaryResultSetRow)

// BinaryResultSetRow represents a MySQL binary resultset row response packet.
type BinaryResultSetRow struct {
	*packet
	columns []string
}

func newBinaryResultSetRowWithPacket(pkt *packet, opts ...BinaryResultSetRowOption) *BinaryResultSetRow {
	row := &BinaryResultSetRow{
		packet:  pkt,
		columns: []string{},
	}
	row.SetOptions(opts...)
	return row
}

// WithBinaryResultSetRowColmunCount returns a binary resultset row option to set the column count.
func WithBinaryResultSetRowColmunCount(c uint64) BinaryResultSetRowOption {
	return func(pkt *BinaryResultSetRow) {
		pkt.columns = make([]string, c)
	}
}

// WithBinaryResultSetRowColmuns returns a binary resultset row option to set the columns.
func WithBinaryResultSetRowColmuns(columns []string) BinaryResultSetRowOption {
	return func(pkt *BinaryResultSetRow) {
		pkt.columns = columns
	}
}

// NewBinaryResultSetRow returns a new BinaryResultSetRow.
func NewBinaryResultSetRow(opts ...BinaryResultSetRowOption) *BinaryResultSetRow {
	return newBinaryResultSetRowWithPacket(newPacket(), opts...)
}

// NewBinaryResultSetRowFromReader returns a new BinaryResultSetRow from the reader.
func NewBinaryResultSetRowFromReader(reader *PacketReader, opts ...BinaryResultSetRowOption) (*BinaryResultSetRow, error) {
	pktReader, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	row := newBinaryResultSetRowWithPacket(pktReader, opts...)
	columnCount := len(row.columns)
	for n := 0; n < columnCount; n++ {
		column, err := row.ReadColumn()
		if err != nil {
			return nil, err
		}
		row.columns[n] = column
	}

	return row, nil
}

// SetOptions sets the options.
func (row *BinaryResultSetRow) SetOptions(opts ...BinaryResultSetRowOption) {
	for _, opt := range opts {
		opt(row)
	}
}

// ReadColumn reads a column.
func (row *BinaryResultSetRow) ReadColumn() (string, error) {
	return row.ReadLengthEncodedString()
}

// Columns returns the columns.
func (row *BinaryResultSetRow) Columns() []any {
	columns := make([]any, len(row.columns))
	for n, column := range row.columns {
		columns[n] = column
	}
	return columns
}

// Bytes returns the packet bytes.
func (row *BinaryResultSetRow) Bytes() ([]byte, error) {
	w := NewPacketWriter()
	for _, column := range row.columns {
		err := w.WriteLengthEncodedString(column)
		if err != nil {
			return nil, err
		}
	}
	row.SetPayload(w.Bytes())
	return row.packet.Bytes()
}
