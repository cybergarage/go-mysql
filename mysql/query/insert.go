// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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

package query

import (
	"fmt"

	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Insert represents an INSERT or REPLACE statement.
type Insert struct {
	*vitesssp.Insert
}

// NewInsertWithInsert creates a insert statement from the raw query.
func NewInsertWithInsert(stmt *vitesssp.Insert) *Insert {
	return &Insert{Insert: stmt}
}

// TableName returns the target table name.
func (stmt *Insert) TableName() string {
	return stmt.Insert.Table.Name.String()
}

// Columns returns all column values with nil when the all columns can be converted, otherwith with a last error.
func (stmt *Insert) Columns() (*Columns, error) {
	rows, ok := stmt.Rows.(vitesssp.Values)
	if !ok || len(rows) <= 0 {
		return nil, fmt.Errorf(errorColumnNotFound, rows)
	}

	columns := NewColumns()
	for n, expr := range rows[0] {
		switch v := expr.(type) {
		case *ComparisonExpr:
			col, ok := v.Left.(*ColName)
			if !ok {
				continue
			}
			val, ok := v.Right.(*Literal)
			if !ok {
				return nil, fmt.Errorf(errorUnexpectedExpression, expr)
			}
			column, err := NewColumnWithNameAndValue(col.Name.String(), val)
			if err != nil {
				return nil, fmt.Errorf(errorUnexpectedExpression, expr)
			}
			columns.AddColumn(column)
		case *Literal:
			name := stmt.Insert.Columns[n].String()
			column, err := NewColumnWithNameAndValue(name, expr)
			if err != nil {
				return nil, err
			}
			columns.AddColumn(column)
		default:
			return nil, fmt.Errorf(errorUnexpectedExpression, expr)
		}
	}

	return columns, nil
}
