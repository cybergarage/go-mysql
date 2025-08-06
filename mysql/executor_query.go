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
	"github.com/cybergarage/go-sqlparser/sql"
)

// defaultQueryExecutor represents a base query message executor.
type defaultQueryExecutor struct {
	sqlExecutor SQLExecutor
}

// NewDefaultQueryExecutor returns a base query message executor instance.
func NewDefaultQueryExecutor() QueryExecutor {
	return &defaultQueryExecutor{
		sqlExecutor: nil,
	}
}

// SetSQLExecutor sets a SQL executor.
func (executor *defaultQueryExecutor) SetSQLExecutor(se SQLExecutor) {
	executor.sqlExecutor = se
}

// CreateDatabase handles a CREATE DATABASE query.
func (executor *defaultQueryExecutor) CreateDatabase(conn Conn, stmt sql.CreateDatabase) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.CreateDatabase(conn, stmt))
}

// CreateTable handles a CREATE TABLE query.
func (executor *defaultQueryExecutor) CreateTable(conn Conn, stmt sql.CreateTable) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.CreateTable(conn, stmt))
}

// AlterDatabase handles a ALTER DATABASE query.
func (executor *defaultQueryExecutor) AlterDatabase(conn Conn, stmt sql.AlterDatabase) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.AlterDatabase(conn, stmt))
}

// AlterTable handles a ALTER TABLE query.
func (executor *defaultQueryExecutor) AlterTable(conn Conn, stmt sql.AlterTable) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.AlterTable(conn, stmt))
}

// DropDatabase handles a DROP DATABASE query.
func (executor *defaultQueryExecutor) DropDatabase(conn Conn, stmt sql.DropDatabase) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.DropDatabase(conn, stmt))
}

// DropIndex handles a DROP INDEX query.
func (executor *defaultQueryExecutor) DropTable(conn Conn, stmt sql.DropTable) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.DropTable(conn, stmt))
}

// Insert handles a INSERT query.
func (executor *defaultQueryExecutor) Insert(conn Conn, stmt sql.Insert) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.Insert(conn, stmt))
}

// Select handles a SELECT query.
func (executor *defaultQueryExecutor) Select(conn Conn, stmt sql.Select) (Response, error) {
	rs, err := executor.sqlExecutor.Select(conn, stmt)
	if err != nil {
		return nil, err
	}

	return protocol.NewTextResultSetFromResultSet(rs)
}

// Update handles a UPDATE query.
func (executor *defaultQueryExecutor) Update(conn Conn, stmt sql.Update) (Response, error) {
	rs, err := executor.sqlExecutor.Update(conn, stmt)
	if err != nil {
		return protocol.NewResponseWithError(err)
	}

	return protocol.NewOK(
		protocol.WithOKAffectedRows(uint64(rs.RowsAffected())),
	)
}

// Delete handles a DELETE query.
func (executor *defaultQueryExecutor) Delete(conn Conn, stmt sql.Delete) (Response, error) {
	rs, err := executor.sqlExecutor.Delete(conn, stmt)
	if err != nil {
		return protocol.NewResponseWithError(err)
	}

	return protocol.NewOK(
		protocol.WithOKAffectedRows(uint64(rs.RowsAffected())),
	)
}

// Begin handles a BEGIN query.
func (executor *defaultQueryExecutor) Begin(conn Conn, stmt sql.Begin) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.Begin(conn, stmt))
}

// Commit handles a COMMIT query.
func (executor *defaultQueryExecutor) Commit(conn Conn, stmt sql.Commit) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.Commit(conn, stmt))
}

// Rollback handles a ROLLBACK query.
func (executor *defaultQueryExecutor) Rollback(conn Conn, stmt sql.Rollback) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.Rollback(conn, stmt))
}

// Use handles a USE query.
func (executor *defaultQueryExecutor) Use(conn Conn, stmt sql.Use) (Response, error) {
	return protocol.NewResponseWithError(executor.sqlExecutor.Use(conn, stmt))
}

// ErrorHandler represents a user error handler.
func (executor *defaultQueryExecutor) ParserError(conn Conn, stmt string, err error) (Response, error) {
	return nil, err
}
