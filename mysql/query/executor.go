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

package query

import (
	"github.com/cybergarage/go-mysql/mysql/net"
)

// DDOExecutor defines a executor interface for DDO (Data Definition Operations).
type DDOExecutor interface {
	// CreateDatabase handles a CREATE DATABASE query.
	CreateDatabase(net.Conn, CreateDatabase) error
	// CreateTable handles a CREATE TABLE query.
	CreateTable(net.Conn, CreateTable) error
	// AlterDatabase handles a ALTER DATABASE query.
	AlterDatabase(net.Conn, AlterDatabase) error
	// AlterTable handles a ALTER TABLE query.
	AlterTable(net.Conn, AlterTable) error
	// DropDatabase handles a DROP DATABASE query.
	DropDatabase(net.Conn, DropDatabase) error
	// DropIndex handles a DROP INDEX query.
	DropTable(net.Conn, DropTable) error
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Insert handles a INSERT query.
	Insert(net.Conn, Insert) error
	// Select handles a SELECT query.
	Select(net.Conn, Select) (ResultSet, error)
	// Update handles a UPDATE query.
	Update(net.Conn, Update) (ResultSet, error)
	// Delete handles a DELETE query.
	Delete(net.Conn, Delete) (ResultSet, error)
}

// TCLExecutor defines a executor interface for TCL (Transaction Control Language).
type TCLExecutor interface {
	// Begin handles a BEGIN query.
	Begin(net.Conn, Begin) error
	// Commit handles a COMMIT query.
	Commit(net.Conn, Commit) error
	// Rollback handles a ROLLBACK query.
	Rollback(net.Conn, Rollback) error
}

// ExtraExecutor defines a executor interface for extra operations.
type ExtraExecutor interface {
	// Use handles a USE query.
	Use(net.Conn, Use) error
}

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
	ExtraExecutor
}

// ErrorHandler represents a user error handler.
type ErrorHandler interface {
	ParserError(net.Conn, string, error) error
}

// Executor represents a frontend message executor.
type Executor interface { // nolint: interfacebloat
	TCLExecutor
	QueryExecutor
	ErrorHandler
}
