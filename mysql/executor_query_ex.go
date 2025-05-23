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

package mysql

import (
	"github.com/cybergarage/go-mysql/mysql/protocol"
	"github.com/cybergarage/go-mysql/mysql/query"
	sql "github.com/cybergarage/go-sqlparser/sql/query"
)

type defaultExQueryExecutor struct {
	QueryExecutor
}

func NewDefaultExQueryExecutorWith(executor QueryExecutor) ExQueryExecutor {
	return &defaultExQueryExecutor{
		QueryExecutor: executor,
	}
}

// CreateIndex handles a CREATE INDEX query.
func (executor *defaultExQueryExecutor) CreateIndex(conn Conn, stmt query.CreateIndex) (Response, error) {
	alterStmt, err := sql.NewAlterTableFrom(stmt)
	if err != nil {
		return nil, err
	}
	return executor.QueryExecutor.AlterTable(conn, alterStmt)
}

// DropIndex handles a DROP INDEX query.
func (executor *defaultExQueryExecutor) DropIndex(conn Conn, stmt query.DropIndex) (Response, error) {
	alterStmt, err := sql.NewAlterTableFrom(stmt)
	if err != nil {
		return nil, err
	}
	return executor.QueryExecutor.AlterTable(conn, alterStmt)
}

// Truncate handles a TRUNCATE query.
func (executor *defaultExQueryExecutor) Truncate(conn Conn, stmt query.Truncate) (Response, error) {
	for _, table := range stmt.Tables() {
		stmt := sql.NewDeleteWith(table, sql.NewCondition())
		_, err := executor.QueryExecutor.Delete(conn, stmt)
		if err != nil {
			return nil, err
		}
	}
	return protocol.NewResponseWithError(nil)
}
