// Copyright (C) 2020 Satoshi Konno. All rights reserved.
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
	"sync"

	"vitess.io/vitess/go/vt/sqlparser"
	vitess "vitess.io/vitess/go/vt/sqlparser"
)

// Row represents a row object which includes query execution results.
type Row struct {
	sync.Mutex
	*Columns
}

// NewRow return a row.
func NewRow() *Row {
	return NewRowWithColumns(NewColumns())
}

// NewRowWithColumns return a row with the specified columns.
func NewRowWithColumns(columns *Columns) *Row {
	row := &Row{
		Columns: columns,
	}
	return row
}

// NewRowWithValTuple return a row with the specified value tuple.
func NewRowWithValTuple(valTuple vitess.ValTuple) (*Row, error) {
	columns := NewColumns()
	for _, val := range valTuple {
		column, err := NewColumnWithComparisonExpr(val)
		if err != nil {
			return nil, err
		}
		columns.AddColumn(column)
	}
	return NewRowWithColumns(columns), nil
}

// Equals returns true when the specified row is equals to this row, otherwise false.
func (row *Row) Equals(other *Row) bool {
	return row.Columns.Equals(NewColumnsWithColumns(other.GetColumns()))
}

func (row *Row) hasMatchedColumn(name string, val interface{}) bool {
	return false
}

// IsMatched returns true when the row is satisfied with the specified condition, otherwise false.
func (row *Row) IsMatched(cond *Condition) bool {
	if cond == nil {
		return true
	}

	isMatched := false
	switch e := cond.Expr.(type) {
	case *sqlparser.ComparisonExpr:
		switch l := e.Left.(type) {
		case *sqlparser.ColName:
			switch r := e.Right.(type) {
			case sqlparser.ValTuple:
				for _, val := range r {
					switch v := val.(type) {
					case *sqlparser.SQLVal:
						value, err := NewValueWithSQLVal(v)
						if err != nil {
							continue
						}
						if row.hasMatchedColumn(l.Name.String(), value) {
							isMatched = true
						}
					}
				}
			}
		}
	}
	return isMatched
}

// String returns the string representation.
func (row *Row) String() string {
	str := ""
	for n, col := range row.GetColumns() {
		if 0 < n {
			str += ", "
		}
		str += fmt.Sprintf("%s", col.Value())
	}
	return str
}
