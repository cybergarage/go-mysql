// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
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
// Result Set Packets - MariaDB Knowledge Base
// https://mariadb.com/kb/en/result-set-packets/

// NewBinaryResultSetRowFromTextResultSetRow creates a new BinaryResultSetRow from a TextResultSetRow.
func NewBinaryResultSetRowFromTextResultSetRow(columDefs []ColumnDef, txtRow ResultSetRow) (*BinaryResultSetRow, error) {
	txtColumns := txtRow.Columns()

	columnCnt := len(txtColumns)
	if columnCnt != len(columDefs) {
		return nil, newErrInvalidColumnCount(columnCnt, len(columDefs))
	}

	nullBitmap := NewNullBitmap(
		WithNullBitmapNumFields(columnCnt),
		WithNullBitmapOffset(0),
	)
	binColums := []*BinaryResultSetColumn{}

	for n, txtColum := range txtColumns {
		if txtColum == nil {
			nullBitmap.SetNull(n, true)
			continue
		}

		binColum, err := NewBinaryResultSetColumn(
			WithBinaryResultSetColumnType(FieldType(columDefs[n].ColType())),
			WithBinaryResultSetColumnValue(txtColum),
		)
		if err != nil {
			return nil, err
		}

		binColums = append(binColums, binColum)
	}

	return NewBinaryResultSetRow(
		WithBinaryResultSetRowColumnDefs(columDefs),
		WithBinaryResultSetRowNullBitmap(nullBitmap),
		WithBinaryResultSetRowColumns(binColums),
	), nil
}
