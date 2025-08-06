// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
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

package protocol

import (
	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-sqlparser/sql"
	"github.com/cybergarage/go-sqlparser/sql/system"
)

// MySQL: Column Definition
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset_column_definition.html

// NewColumnDefFromReader returns a new ColumnDef from the reader.
func NewColumnDefsFromResultSet(rs sql.ResultSet) ([]ColumnDef, error) {
	schema := rs.Schema()
	columns := schema.Columns()

	columnDefs := make([]ColumnDef, len(columns))
	for n, column := range columns {
		t, err := query.NewFieldTypeFrom(column.DataType())
		if err != nil {
			return nil, err
		}

		c, err := query.NewColumnDefFlagFrom(column.Constraint())
		if err != nil {
			return nil, err
		}

		columnDef := NewColumnDef(
			WithColumnDefSchema(schema.DatabaseName()),
			WithColumnDefTable(schema.TableName()),
			WithColumnDefName(column.Name()),
			WithColumnDefType(uint8(t)),
			WithColumnDefFlags(uint16(c)),
		)
		columnDefs[n] = columnDef
	}

	return columnDefs, nil
}

// NewColumnDefFromSystemSchemaColumn returns a new ColumnDef from the system schema column.
func NewColumnDefsFromSystemSchemaColumn(column system.SchemaColumn) (ColumnDef, error) {
	t, err := query.NewFieldTypeFrom(column.DataType())
	if err != nil {
		return nil, err
	}

	c, err := query.NewColumnDefFlagFrom(column.Constraint())
	if err != nil {
		return nil, err
	}

	columnDef := NewColumnDef(
		WithColumnDefSchema(column.Schema()),
		WithColumnDefTable(column.Table()),
		WithColumnDefName(column.Name()),
		WithColumnDefType(uint8(t)),
		WithColumnDefFlags(uint16(c)),
	)

	return columnDef, nil
}
