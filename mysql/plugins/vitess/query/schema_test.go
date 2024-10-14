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
	"testing"
)

func TestNewSchema(t *testing.T) {
	NewSchemaWithDDL(nil)
}

func TestNewSchemaWithName(t *testing.T) {
	tblName := "foo"
	schema := NewSchemaWithName(tblName)

	if tblName != schema.TableName() {
		t.Errorf("%s != %s", tblName, schema.TableName())
	}
}

func TestSchemaFindColumn(t *testing.T) {
	testPrimaryColumnName := "name"
	testQuery := "CREATE TABLE " + testPrimaryColumnName + " (name VARCHAR(255) PRIMARY KEY, age INT)"

	parser := NewParser()
	stmt, err := parser.Parse(testQuery)
	if err != nil {
		t.Error(err)
		return
	}

	ddl, ok := stmt.(DDL)
	if !ok {
		t.Errorf("%v", stmt)
		return
	}

	schema := NewSchemaWithDDL(ddl)
	columnDef, ok := schema.FindPrimaryColumn()
	if !ok {
		t.Errorf("%v", stmt)
		return
	}

	if columnDef.Name.String() != testPrimaryColumnName {
		t.Errorf("%v", stmt)
		return
	}
}
