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
// MySQL: Text Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html
// Result Set Packets - MariaDB Knowledge Base
// https://mariadb.com/kb/en/result-set-packets/
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/#tinyint-binary-encoding

// NewBinaryResultSetFromTextResultSet creates a new BinaryResultSet from a TextResultSet.
func NewBinaryResultSetFromTextResultSet(txtRs *TextResultSet, opts ...BinaryResultSetOption) (*BinaryResultSet, error) {
	opts = append(opts, WithBinaryResultSetColumnDefs(txtRs.ColumnDefs()))

	binRows := make([]BinaryResultSetRow, 0, len(txtRs.Rows()))
	for _, txtRow := range txtRs.Rows() {
		binRow, err := NewBinaryResultSetRowFromTextResultSetRow(txtRow)
		if err != nil {
			return nil, err
		}
		binRows = append(binRows, *binRow)
	}
	opts = append(opts, WithBinaryResultSetRows(binRows))

	return NewBinaryResultSet(opts...)
}
