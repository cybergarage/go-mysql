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

package mysql

import (
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
