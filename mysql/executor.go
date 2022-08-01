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

package mysql

import (
	"context"
	"go-mysql/mysql/query"
)

// DDOExecutor defines a executor interface for DDO (Data Definition Operations).
type DDOExecutor interface {
	// CreateDatabase should handle a CREATE database statement.
	CreateDatabase(context.Context, *Conn, query.DBDDL) (*Result, error)
	// AlterDatabase should handle a ALTER database statement.
	AlterDatabase(context.Context, *Conn, query.DBDDL) (*Result, error)
	// DropDatabase should handle a DROP database statement.
	DropDatabase(context.Context, *Conn, query.DBDDL) (*Result, error)
	// CreateTable should handle a CREATE table statement.
	CreateTable(context.Context, *Conn, query.DDL) (*Result, error)
	// AlterTable should handle a ALTER table statement.
	AlterTable(context.Context, *Conn, query.DDL) (*Result, error)
	// DropTable should handle a DROP table statement.
	DropTable(context.Context, *Conn, query.DDL) (*Result, error)
	// RenameTable should handle a RENAME table statement.
	RenameTable(context.Context, *Conn, query.DDL) (*Result, error)
	// TruncateTable should handle a TRUNCATE table statement.
	TruncateTable(context.Context, *Conn, query.DDL) (*Result, error)
	// AnalyzeTable should handle a ANALYZE table statement.
	AnalyzeTable(context.Context, *Conn, query.DDL) (*Result, error)
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Insert should handle a INSERT statement.
	Insert(context.Context, *Conn, *query.Insert) (*Result, error)
	// Update should handle a UPDATE statement.
	Update(context.Context, *Conn, *query.Update) (*Result, error)
	// Delete should handle a DELETE statement.
	Delete(context.Context, *Conn, *query.Delete) (*Result, error)
	// Select should handle a SELECT statement.
	Select(context.Context, *Conn, *query.Select) (*Result, error)
}

// DCOExecutor defines a executor interface for DCO (Data Control Operations).
type DCOExecutor interface {
}

// DAOExecutor defines a executor interface for DAO (Database Administration Operations).
type DAOExecutor interface {
	// ShowDatabases should handle a SHOW DATABASES statement.
	ShowDatabases(context.Context, *Conn) (*Result, error)
	// ShowTables should handle a SHOW TABLES statement.
	ShowTables(context.Context, *Conn, string) (*Result, error)
}

// QueryExecutor represents an interface to execute all SQL queries.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
	DCOExecutor
	DAOExecutor
}
