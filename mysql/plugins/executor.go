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

package plugins

import (
	"github.com/cybergarage/go-sqlparser/sql"
)

// DDOExecutor defines a executor interface for DDO (Data Definition Operations).
type DDOExecutor interface {
	// CreateDatabase handles a CREATE DATABASE query.
	CreateDatabase(Conn, sql.CreateDatabase) error
	// CreateTable handles a CREATE TABLE query.
	CreateTable(Conn, sql.CreateTable) error
	// AlterDatabase handles a ALTER DATABASE query.
	AlterDatabase(Conn, sql.AlterDatabase) error
	// AlterTable handles a ALTER TABLE query.
	AlterTable(Conn, sql.AlterTable) error
	// DropDatabase handles a DROP DATABASE query.
	DropDatabase(Conn, sql.DropDatabase) error
	// DropIndex handles a DROP INDEX query.
	DropTable(Conn, sql.DropTable) error
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Insert handles a INSERT query.
	Insert(Conn, sql.Insert) error
	// Select handles a SELECT query.
	Select(Conn, sql.Select) (ResultSet, error)
	// Update handles a UPDATE query.
	Update(Conn, sql.Update) (ResultSet, error)
	// Delete handles a DELETE query.
	Delete(Conn, sql.Delete) (ResultSet, error)
}

// TCLExecutor defines a executor interface for TCL (Transaction Control Language).
type TCLExecutor interface {
	// Begin handles a BEGIN query.
	Begin(Conn, sql.Begin) error
	// Commit handles a COMMIT query.
	Commit(Conn, sql.Commit) error
	// Rollback handles a ROLLBACK query.
	Rollback(Conn, sql.Rollback) error
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
}

// ErrorHandler represents a user error handler.
type ErrorHandler interface {
	ParserError(Conn, string, error) error
}

// Executor represents a frontend message executor.
type Executor interface { // nolint: interfacebloat
	TCLExecutor
	QueryExecutor
	ErrorHandler
}
