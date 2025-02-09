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

// MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/

// BinaryResultSetRowOption represents a MySQL binary resultset row option.
type BinaryResultSetRowOption func(*BinaryResultSetRow)

// BinaryResultSetRow represents a MySQL binary resultset row response packet.
type BinaryResultSetRow struct {
	*packet
	columnDefs []ColumnDef
	nullBitmap *NullBitmap
	colums     []*BinaryResultSetColumn
}

func newBinaryResultSetRowWithPacket(pkt *packet, opts ...BinaryResultSetRowOption) *BinaryResultSetRow {
	row := &BinaryResultSetRow{
		packet:     pkt,
		columnDefs: nil,
		nullBitmap: nil,
		colums:     nil,
	}
	for _, opt := range opts {
		opt(row)
	}
	return row
}

// WithBinaryResultSetRowColumns returns a binary resultset row option to set the columns.
func WithBinaryResultSetRowColumns(columns []*BinaryResultSetColumn) BinaryResultSetRowOption {
	return func(row *BinaryResultSetRow) {
		row.colums = columns
	}
}

// WithBinaryResultSetRowColumnDefs returns a binary resultset row option to set the column definitions.
func WithBinaryResultSetRowColumnDefs(columnDefs []ColumnDef) BinaryResultSetRowOption {
	return func(row *BinaryResultSetRow) {
		row.columnDefs = columnDefs
	}
}

// WithBinaryResultSetRowNullBitmap returns a binary resultset row option to set the null bitmap.
func WithBinaryResultSetRowNullBitmap(nullBitmap *NullBitmap) BinaryResultSetRowOption {
	return func(row *BinaryResultSetRow) {
		row.nullBitmap = nullBitmap
	}
}

// NewBinaryResultSetRow returns a new BinaryResultSetRow.
func NewBinaryResultSetRow(opts ...BinaryResultSetRowOption) *BinaryResultSetRow {
	return newBinaryResultSetRowWithPacket(newPacket(), opts...)
}

// NewBinaryResultSetRowFromReader returns a new BinaryResultSetRow from the reader.
func NewBinaryResultSetRowFromReader(reader *PacketReader, opts ...BinaryResultSetRowOption) (*BinaryResultSetRow, error) {
	rowPkt, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	row := newBinaryResultSetRowWithPacket(rowPkt, opts...)

	numColumns := len(row.columnDefs)

	// 0x00 header

	_, err = reader.ReadByte()
	if err != nil {
		return nil, err
	}

	// NullBitmap

	nullBitmapBytes := make([]byte, CalculateNullBitmapLength(numColumns, 0))
	_, err = reader.ReadBytes(nullBitmapBytes)
	if err != nil {
		return nil, err
	}
	row.nullBitmap = NewNullBitmap(
		WithNullBitmapNumFields(numColumns),
		WithNullBitmapOffset(0),
		WithNullBitmapBytes(nullBitmapBytes),
	)

	// for each column

	row.colums = []*BinaryResultSetColumn{}

	for n := 0; n < numColumns; n++ {
		opts := []BinaryResultSetColumnOption{
			WithwBinaryResultSetColumnType(FieldType(row.columnDefs[n].ColType())),
		}
		var column *BinaryResultSetColumn
		if !row.nullBitmap.IsNull(n) {
			column, err = NewBinaryResultSetColumnFromReader(reader, opts...)
			if err != nil {
				return nil, err
			}
		} else {
			column = NewBinaryResultSetColumn(opts...)
		}
		row.colums = append(row.colums, column)
	}

	return row, nil
}

// Bytes returns the bytes.
func (row *BinaryResultSetRow) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	_, err := w.WriteBytes(row.HeaderBytes())
	if err != nil {
		return nil, err
	}

	// 0x00 header

	err = w.WriteByte(0x00)
	if err != nil {
		return nil, err
	}

	// NullBitmap

	_, err = w.WriteBytes(row.nullBitmap.Bytes())
	if err != nil {
		return nil, err
	}

	// for each column

	for n, colum := range row.colums {
		if row.nullBitmap.IsNull(n) {
			continue
		}
		bytes, err := colum.Bytes()
		if err != nil {
			return nil, err
		}
		_, err = w.WriteBytes(bytes)
		if err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}
