// Copyright (C) 2020 Satoshi Konno. All rights reserved.
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
	vitesspq "vitess.io/vitess/go/vt/proto/query"
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

const (
	// NOTE: TestSchemaFindColumn checks whether ColKeyPrimary equals vt.sqlparser.colKeyPrimary
	// because it is an internal constant variable.
	ColKeyPrimary = vitesssp.ColumnKeyOption(1)
)

// Schema represents a table schema.
type Schema struct {
	*DDL
}

// NewSchemaWithDDL returns a schema with the specified DDL.
func NewSchemaWithDDL(ddl *DDL) *Schema {
	return &Schema{
		DDL: ddl,
	}
}

// NewSchemaWithName returns a schema with the specified table name.
func NewSchemaWithName(name string) *Schema {
	ddl := &vitesssp.DDL{
		Table: vitesssp.TableName{Name: vitesssp.NewTableIdent(name)},
	}
	return NewSchemaWithDDL(ddl)
}

// TableName returns the table name.
func (schema *Schema) TableName() string {
	return schema.DDL.Table.Name.String()
}

// FindPrimaryColumn returns the specified columns definition.
func (schema *Schema) FindPrimaryColumn() (*ColumnDefinition, bool) {
	for _, column := range schema.DDL.TableSpec.Columns {
		if column.Type.KeyOpt == ColKeyPrimary {
			return column, true
		}
	}
	return nil, false
}

// FindColumn returns the specified columns definition.
func (schema *Schema) FindColumn(name string) (*ColumnDefinition, bool) {
	for _, column := range schema.DDL.TableSpec.Columns {
		if column.Name.EqualString(name) {
			return column, true
		}
	}
	return nil, false
}

// ToFields converts a column definitions to a vitess fields.
// nolint: exhaustive
func (schema *Schema) ToFields(db *Database) ([]*Field, error) {
	fields := make([]*Field, 0)
	tblName := schema.TableName()
	for _, column := range schema.DDL.TableSpec.Columns {
		colName := column.Name.String()
		// FIXME: Set more appreciate column length to check official MySQL implementation
		colLen := 65535 // len(name) + 1
		field := &Field{
			Database:     db.Name(),
			Table:        tblName,
			OrgTable:     tblName,
			Name:         colName,
			OrgName:      colName,
			ColumnLength: uint32(colLen),
			Charset:      255, // utf8mb4,
			Type:         column.Type.SQLType(),
		}

		switch field.Type {
		case vitesspq.Type_BLOB, vitesspq.Type_TEXT:
			field.Flags = field.Flags | uint32(vitesspq.MySqlFlag_BLOB_FLAG)
		}
		if column.Type.NotNull {
			field.Flags = field.Flags | uint32(vitesspq.MySqlFlag_NOT_NULL_FLAG)
		}
		if column.Type.KeyOpt == ColKeyPrimary {
			field.Flags = field.Flags | uint32(vitesspq.MySqlFlag_PRI_KEY_FLAG)
			field.Flags = field.Flags | uint32(vitesspq.MySqlFlag_NOT_NULL_FLAG)
		}

		fields = append(fields, field)
	}
	return fields, nil
}
