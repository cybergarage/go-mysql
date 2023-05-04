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
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysql/query"
)

const (
	baseExecutorNotImplemented = "Not implemented %v"
)

// BaseExecutor is a stub listener to implement only target executor functions.
type BaseExecutor struct {
}

// NewBaseExecutor returns an in-memory executorinstance.
func NewBaseExecutor() *BaseExecutor {
	return &BaseExecutor{}
}

// CreateDatabase should handle a CREATE database statement.
func (executor *BaseExecutor) CreateDatabase(conn *Conn, stmt *query.Database) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// AlterDatabase should handle a ALTER database statement.
func (executor *BaseExecutor) AlterDatabase(conn *Conn, stmt *query.Database) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// DropDatabase should handle a DROP database statement.
func (executor *BaseExecutor) DropDatabase(conn *Conn, stmt *query.Database) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// CreateTable should handle a CREATE table statement.
func (executor *BaseExecutor) CreateTable(conn *Conn, stmt *query.Schema) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// AlterTable should handle a ALTER table statement.
func (executor *BaseExecutor) AlterTable(conn *Conn, stmt *query.Schema) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// DropTable should handle a DROP table statement.
func (executor *BaseExecutor) DropTable(conn *Conn, stmt *query.Schema) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// RenameTable should handle a RENAME table statement.
func (executor *BaseExecutor) RenameTable(conn *Conn, stmt *query.Schema) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// TruncateTable should handle a TRUNCATE table statement.
func (executor *BaseExecutor) TruncateTable(conn *Conn, stmt *query.Schema) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// Insert should handle a INSERT statement.
func (executor *BaseExecutor) Insert(conn *Conn, stmt *query.Insert) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// Update should handle a UPDATE statement.
func (executor *BaseExecutor) Update(conn *Conn, stmt *query.Update) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// Delete should handle a DELETE statement.
func (executor *BaseExecutor) Delete(conn *Conn, stmt *query.Delete) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// Select should handle a SELECT statement.
func (executor *BaseExecutor) Select(conn *Conn, stmt *query.Select) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, stmt)
	return NewResult(), nil
}

// ShowDatabases should handle a SHOW DATABASES statement.
func (executor *BaseExecutor) ShowDatabases(conn *Conn) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, "ShowDatabases")
	return NewResult(), nil
}

// ShowTables should handle a SHOW TABLES statement.
func (executor *BaseExecutor) ShowTables(conn *Conn, database string) (*Result, error) {
	log.Debugf(baseExecutorNotImplemented, "ShowTables")
	return NewResult(), nil
}
