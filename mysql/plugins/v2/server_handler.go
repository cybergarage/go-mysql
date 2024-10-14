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

package v2

import (
	"github.com/cybergarage/go-mysql/mysql/plugins"
	"github.com/cybergarage/go-mysql/mysql/protocol"
	"github.com/cybergarage/go-sqlparser/sql"
)

// CreateDatabase handles a CREATE DATABASE query.
func (server *Server) CreateDatabase(conn Conn, stmt sql.CreateDatabase) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().CreateDatabase(conn, stmt))
}

// CreateTable handles a CREATE TABLE query.
func (server *Server) CreateTable(conn Conn, stmt sql.CreateTable) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().CreateTable(conn, stmt))
}

// AlterDatabase handles a ALTER DATABASE query.
func (server *Server) AlterDatabase(conn Conn, stmt sql.AlterDatabase) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().AlterDatabase(conn, stmt))
}

// AlterTable handles a ALTER TABLE query.
func (server *Server) AlterTable(conn Conn, stmt sql.AlterTable) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().AlterTable(conn, stmt))
}

// DropDatabase handles a DROP DATABASE query.
func (server *Server) DropDatabase(conn Conn, stmt sql.DropDatabase) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().DropDatabase(conn, stmt))
}

// DropIndex handles a DROP INDEX query.
func (server *Server) DropTable(conn Conn, stmt sql.DropTable) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().DropTable(conn, stmt))
}

// Insert handles a INSERT query.
func (server *Server) Insert(conn Conn, stmt sql.Insert) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().Insert(conn, stmt))
}

// Select handles a SELECT query.
func (server *Server) Select(conn Conn, stmt sql.Select) (Response, error) {
	return nil, plugins.ErrNotImplemented
}

// Update handles a UPDATE query.
func (server *Server) Update(conn Conn, stmt sql.Update) (Response, error) {
	return nil, plugins.ErrNotImplemented
}

// Delete handles a DELETE query.
func (server *Server) Delete(conn Conn, stmt sql.Delete) (Response, error) {
	return nil, plugins.ErrNotImplemented
}

// Begin handles a BEGIN query.
func (server *Server) Begin(conn Conn, stmt sql.Begin) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().Begin(conn, stmt))
}

// Commit handles a COMMIT query.
func (server *Server) Commit(conn Conn, stmt sql.Commit) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().Commit(conn, stmt))
}

// Rollback handles a ROLLBACK query.
func (server *Server) Rollback(conn Conn, stmt sql.Rollback) (Response, error) {
	return protocol.NewResponseWithError(server.Executor().Rollback(conn, stmt))
}

// ErrorHandler represents a user error handler.
func (server *Server) ParserError(conn Conn, stmt string, err error) (Response, error) {
	return nil, plugins.ErrNotImplemented
}
