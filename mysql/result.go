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
	vitesstypes "vitess.io/vitess/go/sqltypes"
	vitesspb "vitess.io/vitess/go/vt/proto/query"
)

// Result represents a query result.
type Result = vitesstypes.Result

// Field represents a field.
type Field = vitesspb.Field

// BindVariable represents a single bind variable.
type BindVariable = vitesspb.BindVariable

// Row represents a row.
type Row = vitesspb.Value

// Value represents a value.
type Value = vitesstypes.Value

// NewResult returns a blank result.
func NewResult() *Result {
	res := &Result{
		Rows: [][]Value{},
	}
	return res
}
