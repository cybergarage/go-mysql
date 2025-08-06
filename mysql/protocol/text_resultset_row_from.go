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
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-safecast/safecast"
	"github.com/cybergarage/go-sqlparser/sql"
)

// MySQL: Protocol::QueryResponse
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query.html
// MySQL: Text Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html
// MySQL: Text Resultset Row
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset_row.html
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/

const (
	DateTimeFormat = "2006-01-02 15:04:05"
)

// NewTextResultSetRowsFromResultSet returns a new ResultSetRow list from the specified ResultSet.
func NewTextResultSetRowsFromResultSet(rs sql.ResultSet) ([]ResultSetRow, error) {
	rows := []ResultSetRow{}

	for rs.Next() {
		rsRow, err := rs.Row()
		if err != nil {
			return nil, err
		}

		row, err := NewTextResultSetRowFrom(rs.Schema(), rsRow)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return rows, nil
}

// NewTextResultSetRowFrom returns a new ResultSetRow from the specified ResultSetSchema and ResultSetRow.
func NewTextResultSetRowFrom(schema sql.ResultSetSchema, rsRow sql.ResultSetRow) (ResultSetRow, error) {
	schemaColumns := schema.Columns()
	schemaColumnCount := len(schemaColumns)

	rowColumns := make([]*string, len(rsRow.Values()))
	for n, v := range rsRow.Values() {
		if schemaColumnCount <= n {
			return nil, fmt.Errorf("schema column count (%d) is less than row column count (%d)", schemaColumnCount, n)
		}

		columnType := schemaColumns[n].DataType()

		rowValue, err := NewTextResultSetRowValueFrom(columnType, v)
		if err == nil {
			rowColumns[n] = &rowValue
		} else if errors.Is(err, ErrNull) {
			rowColumns[n] = nil
		} else {
			return nil, err
		}
	}

	row := NewTextResultSetRow(
		WithTextResultSetRowColmuns(rowColumns),
	)

	return row, nil
}

// NewTextResultSetRowValueFrom returns a new ResultSetRowValue from the specified DataType and value.
func NewTextResultSetRowValueFrom(t query.DataType, v any) (string, error) {
	if v == nil {
		return "", ErrNull
	}

	switch t {
	case query.CharData, query.CharacterData, query.VarCharData, query.VarCharacterData, query.TextData, query.TinyTextData, query.LongTextData:
		var rv string

		err := safecast.ToString(v, &rv)
		if err != nil {
			return "", err
		}

		return rv, nil
	case query.IntData, query.IntegerData, query.SmallIntData, query.MediumIntData, query.TinyIntData:
		var rv int

		err := safecast.ToInt(v, &rv)
		if err != nil {
			return "", err
		}

		return strconv.Itoa(rv), nil
	case query.FloatData, query.DoubleData, query.RealData:
		var rv float64

		err := safecast.ToFloat64(v, &rv)
		if err != nil {
			return "", err
		}

		return strconv.FormatFloat(rv, 'f', -1, 64), nil
	case query.TimeStampData, query.DateTimeData:
		var rv time.Time

		err := safecast.ToTime(v, &rv)
		if err != nil {
			return "", err
		}

		return rv.Format(DateTimeFormat), nil
	default:
	}

	return fmt.Sprintf("%s", v), nil
}
