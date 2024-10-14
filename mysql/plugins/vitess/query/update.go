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
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Update represents an INSERT or REPLACE statement.
type Update struct {
	*vitesssp.Update
}

// NewUpdateWithUpdate creates a insert statement from the raw query.
func NewUpdateWithUpdate(stmt *vitesssp.Update) *Update {
	return &Update{Update: stmt}
}

// GetColumns returns all column values with nil when the all columns can be converted, otherwith with a last error.
func (stmt *Update) Columns() (*Columns, error) {
	columns := NewColumns()
	for _, expr := range stmt.Exprs {
		name := expr.Name.Name.String()
		column := NewColumnWithName(name)
		err := column.SetValue(expr.Expr)
		if err != nil {
			return nil, err
		}
		columns.AddColumn(column)
	}
	return columns, nil
}

// Tables returns all tables.
func (stmt *Update) Tables() Tables {
	return NewTablesWitExprs(stmt.TableExprs)
}
