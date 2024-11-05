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
	"github.com/cybergarage/go-mysql/mysql/net"
	"github.com/cybergarage/go-mysql/mysql/protocol"
	"github.com/cybergarage/go-sqlparser/sql"
)

// Conn represents a connection.
type Conn = net.Conn

// Response represents a response.
type Response = protocol.Response

// DDOExecutor defines a executor interface for DDO (Data Definition Operations).
type DDOExecutor interface {
	// CreateDatabase handles a CREATE DATABASE query.
	CreateDatabase(Conn, sql.CreateDatabase) (Response, error)
	// CreateTable handles a CREATE TABLE query.
	CreateTable(Conn, sql.CreateTable) (Response, error)
	// AlterDatabase handles a ALTER DATABASE query.
	AlterDatabase(Conn, sql.AlterDatabase) (Response, error)
	// AlterTable handles a ALTER TABLE query.
	AlterTable(Conn, sql.AlterTable) (Response, error)
	// DropDatabase handles a DROP DATABASE query.
	DropDatabase(Conn, sql.DropDatabase) (Response, error)
	// DropIndex handles a DROP INDEX query.
	DropTable(Conn, sql.DropTable) (Response, error)
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Insert handles a INSERT query.
	Insert(Conn, sql.Insert) (Response, error)
	// Select handles a SELECT query.
	Select(Conn, sql.Select) (Response, error)
	// Update handles a UPDATE query.
	Update(Conn, sql.Update) (Response, error)
	// Delete handles a DELETE query.
	Delete(Conn, sql.Delete) (Response, error)
}

// TCLExecutor defines a executor interface for TCL (Transaction Control Language).
type TCLExecutor interface {
	// Begin handles a BEGIN query.
	Begin(Conn, sql.Begin) (Response, error)
	// Commit handles a COMMIT query.
	Commit(Conn, sql.Commit) (Response, error)
	// Rollback handles a ROLLBACK query.
	Rollback(Conn, sql.Rollback) (Response, error)
}

// ExtraExecutor defines a executor interface for extra operations.
type ExtraExecutor interface {
	// Use handles a USE query.
	Use(net.Conn, sql.Use) (Response, error)
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
	ExtraExecutor
}

// ErrorHandler represents a user error handler.
type ErrorHandler interface {
	ParserError(Conn, string, error) (Response, error)
}

// Executor represents a frontend message executor.
type Executor interface { // nolint: interfacebloat
	TCLExecutor
	QueryExecutor
	ErrorHandler
}
