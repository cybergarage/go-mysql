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

package mysql

import (
	"fmt"

	"github.com/cybergarage/go-mysql/mysql/query"
	vitessst "vitess.io/vitess/go/sqltypes"
	vitesspq "vitess.io/vitess/go/vt/proto/query"
)

// Result represents a query result.
type Result = vitessst.Result

// Field represents a field.
type Field = vitesspq.Field

// BindVariable represents a single bind variable.
type BindVariable = vitesspq.BindVariable

// Value represents a value.
type Value = vitessst.Value

// NewValueWith builds a Value using typ and val. If the value and typ
// don't match, it returns an error.
func NewValueWith(typ vitesspq.Type, val []byte) (Value, error) {
	return vitessst.NewValue(typ, val)
}

// NULL represents the NULL value.
func NewNullValue() Value {
	return vitessst.NULL
}

// NewResult returns a successful result without any data.
func NewResult() *Result {
	res := &Result{
		Rows: [][]Value{},
	}
	return res
}

// NewResultWithRowsAffected returns a successful result with the specified rows affected count.
func NewResultWithRowsAffected(n uint64) *Result {
	res := &Result{
		Rows:         [][]Value{},
		RowsAffected: n,
	}
	return res
}

// NewResultWithRows returns a successful result with the specified rows.
func NewResultWithRows(db *query.Database, schema *query.Schema, rows *query.Rows) (*Result, error) {
	fields, err := schema.ToFields(db)
	if err != nil {
		return nil, err
	}

	res := NewResult()
	res.Fields = fields

	resRows := [][]Value{}
	for _, row := range rows.Rows() {
		resValues := []Value{}
		for _, field := range res.Fields {
			// FIXME: The current implementation have to use Row::ColumnByName()
			// because row columns might be unordered.
			column, ok := row.ColumnByName(field.Name)
			if !ok {
				return nil, fmt.Errorf("column (%s) is not found", field.Name)
			}
			resValue, err := column.ToValue()
			if err != nil {
				return nil, err
			}
			resValues = append(resValues, resValue)
		}
		resRows = append(resRows, resValues)
	}

	res.Rows = resRows

	return res, nil
}
