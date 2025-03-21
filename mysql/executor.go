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
	"github.com/cybergarage/go-mysql/mysql/net"
	"github.com/cybergarage/go-mysql/mysql/protocol"
	"github.com/cybergarage/go-sqlparser/sql"
)

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

// DDOExExecutor defines a executor interface for extended DDO (Data Definition Operations).
type DDOExExecutor interface {
	// CreateIndex handles a CREATE INDEX query.
	CreateIndex(Conn, sql.CreateIndex) (Response, error)
	// DropIndex handles a DROP INDEX query.
	DropIndex(Conn, sql.DropIndex) (Response, error)
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Use handles a USE query.
	Use(net.Conn, sql.Use) (Response, error)
	// Insert handles a INSERT query.
	Insert(Conn, sql.Insert) (Response, error)
	// Select handles a SELECT query.
	Select(Conn, sql.Select) (Response, error)
	// Update handles a UPDATE query.
	Update(Conn, sql.Update) (Response, error)
	// Delete handles a DELETE query.
	Delete(Conn, sql.Delete) (Response, error)
}

// DMOExExecutor defines a executor interface for extended DMO (Data Manipulation Operations).
type DMOExExecutor interface {
	// Truncate handles a TRUNCATE query.
	Truncate(Conn, sql.Truncate) (Response, error)
}

// TCOExecutor defines a executor interface for TCL (Transaction Control Operations).
type TCOExecutor interface {
	// Begin handles a BEGIN query.
	Begin(Conn, sql.Begin) (Response, error)
	// Commit handles a COMMIT query.
	Commit(Conn, sql.Commit) (Response, error)
	// Rollback handles a ROLLBACK query.
	Rollback(Conn, sql.Rollback) (Response, error)
}
