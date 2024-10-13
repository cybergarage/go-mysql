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

package plugins

import "github.com/cybergarage/go-sqlparser/sql"

// Server represents a base executor server.
type Server struct {
}

// NewServer returns a base executor server instance.
func NewServer() *Server {
	s := &Server{}
	return s
}

// CreateDatabase handles a CREATE DATABASE query.
func (*Server) CreateDatabase(Conn, sql.CreateDatabase) (Response, error) {
	return nil, ErrNotImplemented
}

// CreateTable handles a CREATE TABLE query.
func (*Server) CreateTable(Conn, sql.CreateTable) (Response, error) {
	return nil, ErrNotImplemented
}

// AlterDatabase handles a ALTER DATABASE query.
func (*Server) AlterDatabase(Conn, sql.AlterDatabase) (Response, error) {
	return nil, ErrNotImplemented
}

// AlterTable handles a ALTER TABLE query.
func (*Server) AlterTable(Conn, sql.AlterTable) (Response, error) {
	return nil, ErrNotImplemented
}

// DropDatabase handles a DROP DATABASE query.
func (*Server) DropDatabase(Conn, sql.DropDatabase) (Response, error) {
	return nil, ErrNotImplemented
}

// DropIndex handles a DROP INDEX query.
func (*Server) DropTable(Conn, sql.DropTable) (Response, error) {
	return nil, ErrNotImplemented
}

// Insert handles a INSERT query.
func (*Server) Insert(Conn, sql.Insert) (Response, error) {
	return nil, ErrNotImplemented
}

// Select handles a SELECT query.
func (*Server) Select(Conn, sql.Select) (Response, error) {
	return nil, ErrNotImplemented
}

// Update handles a UPDATE query.
func (*Server) Update(Conn, sql.Update) (Response, error) {
	return nil, ErrNotImplemented
}

// Delete handles a DELETE query.
func (*Server) Delete(Conn, sql.Delete) (Response, error) {
	return nil, ErrNotImplemented
}

// Begin handles a BEGIN query.
func (*Server) Begin(Conn, sql.Begin) (Response, error) {
	return nil, ErrNotImplemented
}

// Commit handles a COMMIT query.
func (*Server) Commit(Conn, sql.Commit) (Response, error) {
	return nil, ErrNotImplemented
}

// Rollback handles a ROLLBACK query.
func (*Server) Rollback(Conn, sql.Rollback) (Response, error) {
	return nil, ErrNotImplemented
}

// ErrorHandler represents a user error handler.
func (*Server) ParserError(Conn, string, error) (Response, error) {
	return nil, ErrNotImplemented
}
