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

// CreateDatabase handles a CREATE DATABASE query.
func (server *server) CreateDatabase(conn Conn, stmt sql.CreateDatabase) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().CreateDatabase(conn, stmt))
}

// CreateTable handles a CREATE TABLE query.
func (server *server) CreateTable(conn Conn, stmt sql.CreateTable) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().CreateTable(conn, stmt))
}

// AlterDatabase handles a ALTER DATABASE query.
func (server *server) AlterDatabase(conn Conn, stmt sql.AlterDatabase) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().AlterDatabase(conn, stmt))
}

// AlterTable handles a ALTER TABLE query.
func (server *server) AlterTable(conn Conn, stmt sql.AlterTable) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().AlterTable(conn, stmt))
}

// DropDatabase handles a DROP DATABASE query.
func (server *server) DropDatabase(conn Conn, stmt sql.DropDatabase) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().DropDatabase(conn, stmt))
}

// DropIndex handles a DROP INDEX query.
func (server *server) DropTable(conn Conn, stmt sql.DropTable) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().DropTable(conn, stmt))
}

// Insert handles a INSERT query.
func (server *server) Insert(conn Conn, stmt sql.Insert) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().Insert(conn, stmt))
}

// Select handles a SELECT query.
func (server *server) Select(conn Conn, stmt sql.Select) (Response, error) {
	rs, err := server.SQLExecutor().Select(conn, stmt)
	if err != nil {
		return nil, err
	}
	return protocol.NewTextResultSetFromResultSet(rs)
}

// Update handles a UPDATE query.
func (server *server) Update(conn Conn, stmt sql.Update) (Response, error) {
	rs, err := server.SQLExecutor().Update(conn, stmt)
	if err != nil {
		return protocol.NewResponseWithError(err)
	}
	return protocol.NewOK(
		protocol.WithOKAffectedRows(uint64(rs.RowsAffected())),
	)
}

// Delete handles a DELETE query.
func (server *server) Delete(conn Conn, stmt sql.Delete) (Response, error) {
	rs, err := server.SQLExecutor().Delete(conn, stmt)
	if err != nil {
		return protocol.NewResponseWithError(err)
	}
	return protocol.NewOK(
		protocol.WithOKAffectedRows(uint64(rs.RowsAffected())),
	)
}

// Begin handles a BEGIN query.
func (server *server) Begin(conn Conn, stmt sql.Begin) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().Begin(conn, stmt))
}

// Commit handles a COMMIT query.
func (server *server) Commit(conn Conn, stmt sql.Commit) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().Commit(conn, stmt))
}

// Rollback handles a ROLLBACK query.
func (server *server) Rollback(conn Conn, stmt sql.Rollback) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().Rollback(conn, stmt))
}

// Use handles a USE query.
func (server *server) Use(conn Conn, stmt sql.Use) (Response, error) {
	return protocol.NewResponseWithError(server.SQLExecutor().Use(conn, stmt))
}

// ErrorHandler represents a user error handler.
func (server *server) ParserError(conn Conn, stmt string, err error) (Response, error) {
	return nil, err
}
