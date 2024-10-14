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
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Tables represents a table array.
type Tables = []*Table

// NewTablesWitExprs returns a tables with the specified expressions.
func NewTablesWitExprs(exprs vitesssp.TableExprs) Tables {
	tables := make(Tables, len(exprs))
	for n, expr := range exprs {
		tables[n] = NewTableWitExpr(expr)
	}
	return tables
}
