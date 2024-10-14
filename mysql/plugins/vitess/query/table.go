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

// Table represents a destination or source database of query.
type Table struct {
	vitesssp.TableExpr
}

// NewTableWitExpr returns a table with the specified expression.
func NewTableWitExpr(expr vitesssp.TableExpr) *Table {
	return &Table{
		TableExpr: expr,
	}
}

// Name returns the table name.
func (tbl *Table) Name() (string, error) {
	return GetTableName(tbl.TableExpr)
}

// GetTableName returns the table name when the specified TableExpr is a known type.
func GetTableName(v vitesssp.TableExpr) (string, error) {
	switch table := interface{}(v).(type) {
	case (*vitesssp.AliasedTableExpr):
		switch tableData := interface{}(table.Expr).(type) {
		case (vitesssp.SimpleTableExpr):
			switch tableName := interface{}(tableData).(type) {
			case (vitesssp.TableName):
				return tableName.Name.String(), nil
			}
		}
	}

	return "", fmt.Errorf(errorTableUnknownType, v)
}
