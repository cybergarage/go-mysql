// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apachce.org/licenses/LICENSE-2.0
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

	vitess "vitess.io/vitess/go/vt/sqlparser"
)

const (
	NotPrimaryKey    = -1
	SinglePrimaryKey = 0
)

const (
	// ISO8601DateFormat is a data format for SQL primitive
	ISO8601DateFormat = "2006-01-02"
	// ISO8601TimestampFormat is a timestamp format for SQL primitive
	ISO8601TimestampFormat = "2006-01-02T15:04:05-0700"
)

// Column represents a column.
type Column struct {
	name  string
	value interface{}
}

// NewColumn returns a column instanccol.
func NewColumn() *Column {
	return NewColumnWithName("")
}

// NewColumnWithName returns a column with the specified name.
func NewColumnWithName(name string) *Column {
	col := &Column{
		name:  name,
		value: nil,
	}
	return col
}

// NewColumnWithNameAndValue returns a column with the specified name and value.
func NewColumnWithNameAndValue(name string, val interface{}) *Column {
	col := NewColumnWithName(name)
	col.SetValue(val)
	return col
}

// NewColumnWithComparisonExpr returns a column with the specified ComparisonExpr.
func NewColumnWithComparisonExpr(expr interface{}) (*Column, error) {
	cmpExpr, ok := expr.(*vitess.ComparisonExpr)
	if !ok {
		return nil, fmt.Errorf(errorInvalidComparisonExpr, expr)
	}
	cmpCol, ok := cmpExpr.Left.(*vitess.ColName)
	if !ok {
		return nil, fmt.Errorf(errorInvalidComparisonExpr, cmpExpr)
	}
	cmpVal, ok := cmpExpr.Right.(*vitess.SQLVal)
	if !ok {
		return nil, fmt.Errorf(errorInvalidComparisonExpr, cmpExpr)
	}
	val, err := NewValueWithSQLVal(cmpVal)
	if err != nil {
		return nil, err
	}
	return NewColumnWithNameAndValue(cmpCol.Name.String(), val), nil
}

// SetName sets a column name.
func (col *Column) SetName(name string) {
	col.name = name
}

// Name returns the column name.
func (col *Column) Name() string {
	return col.name
}

// SetValue sets a value.
func (col *Column) SetValue(value interface{}) {
	col.value = value
}

// Value returns the value.
func (col *Column) Value() interface{} {
	return col.value
}

// Equals returns true when the specified column is equals to this column, otherwise false.
func (col *Column) Equals(other *Column) bool {
	if col.name != other.name {
		return false
	}
	if !reflect.DeepEqual(col.value, other.value) {
		return false
	}
	return true
}
