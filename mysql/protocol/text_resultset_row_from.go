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
	"fmt"

	"github.com/cybergarage/go-mysql/mysql/query"
)

// MySQL: Protocol::QueryResponse
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query.html
// MySQL: Text Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html
// MySQL: Text Resultset Row
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset_row.html
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/

// NewTextResultSetRowFromResultSetRow returns a new ResultSetRow from the specified ResultSetRow.
func NewTextResultSetRowFromResultSetRow(schema query.ResultSetSchema, rsRow query.ResultSetRow) (ResultSetRow, error) {
	schemaColumns := schema.Columns()
	schemaColumnCount := len(schemaColumns)
	rowColumns := make([]string, len(rsRow.Values()))
	for n, v := range rsRow.Values() {
		if schemaColumnCount <= n {
			return nil, fmt.Errorf("schema column count (%d) is less than row column count (%d)", schemaColumnCount, n)
		}
		columnType := schemaColumns[n].Type()
		switch columnType {
		default:
			rowColumns[n] = fmt.Sprintf("%s", v)
		}
	}
	row := NewTextResultSetRow(
		WithTextResultSetRowColmuns(rowColumns),
	)
	return row, nil
}

// NewTextResultSetRowsFromResultSet returns a new ResultSetRow list from the specified ResultSet.
func NewTextResultSetRowsFromResultSet(rs query.ResultSet) ([]ResultSetRow, error) {
	rows := []ResultSetRow{}
	for rs.Next() {
		row, err := NewTextResultSetRowFromResultSetRow(rs.Schema(), rs.Row())
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}
