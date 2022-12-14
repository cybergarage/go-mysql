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
	"fmt"

	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Delete represents a DELETE statement.
type Delete struct {
	*vitesssp.Delete
}

// NewDeleteWithDelete creates a delete statement from the raw query.
func NewDeleteWithDelete(stmt *vitesssp.Delete) *Delete {
	return &Delete{Delete: stmt}
}

// NewDeleteWithName creates a delete statement from the table name.
func NewDeleteWithName(name string) (*Delete, error) {
	parser := NewParser()
	queryStr := fmt.Sprintf("DELETE FROM %s", name)
	parsedQuery, err := parser.Parse(queryStr)
	if err != nil {
		return nil, err
	}

	delQuery, ok := parsedQuery.(*vitesssp.Delete)
	if !ok {
		return nil, fmt.Errorf("%s", queryStr)
	}

	return NewDeleteWithDelete(delQuery), nil
}

// Tables returns all tables.
func (stmt *Delete) Tables() Tables {
	return NewTablesWitExprs(stmt.TableExprs)
}
