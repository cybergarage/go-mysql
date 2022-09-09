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
	"reflect"
	"strconv"

	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Row represents a row object which includes query execution results.
type Row struct {
	*Columns
}

// NewRow return a row instance.
func NewRow() *Row {
	return NewRowWithColumns(NewColumns())
}

// NewRowWithColumns return a row instance with the specified columns.
func NewRowWithColumns(columns *Columns) *Row {
	row := &Row{
		Columns: columns,
	}
	return row
}

// NewRowWithInsert return a row instance with the specified INSERT statement.
func NewRowWithInsert(stmt *Insert) (*Row, error) {
	columns, err := stmt.Columns()
	if err != nil {
		return nil, err
	}
	return NewRowWithColumns(columns), nil
}

// HasMatchedColumn returns true when the row has the specified column, otherwise false.
func (row *Row) HasMatchedColumn(column *Column) bool {
	foundColumn, ok := row.ColumnByName(column.Name())
	if !ok {
		return false
	}

	deepEqual := func(iv1 interface{}, iv2 interface{}) bool {
		if reflect.DeepEqual(iv1, iv2) {
			return true
		}

		switch v1 := iv1.(type) {
		case string:
			switch v2 := iv2.(type) {
			case string:
				if v1 == v2 {
					return true
				}
			default:
				sv2 := fmt.Sprintf("%v", iv2)
				if v1 == sv2 {
					return true
				}
			}
		case int:
			switch v2 := iv2.(type) {
			case string:
				iv2, err := strconv.Atoi(v2)
				if err == nil && v1 == iv2 {
					return true
				}
			case int64:
				if v1 == int(v2) {
					return true
				}
			case float64:
				if v1 == int(v2) {
					return true
				}
			}
		case float64:
			switch v2 := iv2.(type) {
			case string:
				fv2, err := strconv.ParseFloat(v2, 64)
				if err == nil && v1 == fv2 {
					return true
				}
			case int:
				if int(v1) == v2 {
					return true
				}
			}
		}

		return false
	}

	return deepEqual(foundColumn.Value(), column.Value())
}

// IsMatched returns true when the row satisfies the specified condition, otherwise false.
func (row *Row) IsMatched(cond *Condition) bool {
	if cond == nil {
		return true
	}

	var isMatched func(row *Row, expr Expr) bool
	isMatched = func(row *Row, expr Expr) bool {
		switch v := expr.(type) {
		case *vitesssp.ComparisonExpr:
			col, ok := v.Left.(*ColName)
			if !ok {
				return false
			}
			val, ok := v.Right.(*Literal)
			if !ok {
				return false
			}
			c, err := NewColumnWithNameAndValue(col.Name.String(), val)
			if err != nil {
				return false
			}
			return row.HasMatchedColumn(c)
		case *vitesssp.AndExpr:
			l := isMatched(row, v.Left)
			r := isMatched(row, v.Right)
			if l && r {
				return true
			}
			return false
		case *vitesssp.OrExpr:
			if isMatched(row, v.Left) {
				return true
			}
			if isMatched(row, v.Right) {
				return true
			}
			return false
		}
		return false
	}

	return isMatched(row, cond.Expr)
}

// Update updates row the specified columns.
func (row *Row) Update(columns *Columns) error {
	for _, column := range columns.Columns() {
		name := column.Name()
		rowColumn, ok := row.ColumnByName(name)
		if !ok {
			continue
		}
		rowColumn.SetValue(column.Value())
	}
	return nil
}

// Equals returns true when the specified row is equals to this row, otherwise false.
func (row *Row) Equals(other *Row) bool {
	return row.Columns.Equals(NewColumnsWithColumns(other.AllColumns()))
}

// String returns the string representation.
func (row *Row) String() string {
	str := ""
	for n, col := range row.AllColumns() {
		if 0 < n {
			str += ", "
		}
		str += fmt.Sprintf("%v", col.Value())
	}
	return str
}
