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

package storage

import (
	"context"
	"fmt"
	"go-mysql/mysql"
	"go-mysql/mysql/query"
)

type MemStore struct {
	Databases
}

// NewMemStore returns an in-memory storeinstance.
func NewMemStore() *MemStore {
	store := &MemStore{
		Databases: NewDatabases(),
	}
	return store
}

// CreateDatabase should handle a CREATE database statement.
func (store *MemStore) CreateDatabase(ctx context.Context, conn *mysql.Conn, stmt *query.DBDDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	dbName := stmt.DBName
	_, ok := store.GetDatabase(dbName)
	if ok {
		return mysql.NewResult(), fmt.Errorf(errorDatabaseFound, dbName)
	}
	err := store.AddDatabase(NewDatabaseWithName(dbName))
	if err != nil {
		return mysql.NewResult(), err
	}
	return mysql.NewResult(), nil
}

// AlterDatabase should handle a ALTER database statement.
func (store *MemStore) AlterDatabase(ctx context.Context, conn *mysql.Conn, stmt *query.DBDDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// DropDatabase should handle a DROP database statement.
func (store *MemStore) DropDatabase(ctx context.Context, conn *mysql.Conn, stmt *query.DBDDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// CreateTable should handle a CREATE table statement.
func (store *MemStore) CreateTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	tableName := stmt.Table.Name.String()
	_ = NewTableWithName(tableName)
	return mysql.NewResult(), nil
}

// AlterTable should handle a ALTER table statement.
func (store *MemStore) AlterTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// DropTable should handle a DROP table statement.
func (store *MemStore) DropTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// RenameTable should handle a RENAME table statement.
func (store *MemStore) RenameTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// TruncateTable should handle a TRUNCATE table statement.
func (store *MemStore) TruncateTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// AnalyzeTable should handle a ANALYZE table statement.
func (store *MemStore) AnalyzeTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// Insert should handle a INSERT statement.
func (store *MemStore) Insert(ctx context.Context, conn *mysql.Conn, stmt *query.Insert) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// Update should handle a UPDATE statement.
func (store *MemStore) Update(ctx context.Context, conn *mysql.Conn, stmt *query.Update) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// Delete should handle a DELETE statement.
func (store *MemStore) Delete(ctx context.Context, conn *mysql.Conn, stmt *query.Delete) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}

// Select should handle a SELECT statement.
func (store *MemStore) Select(ctx context.Context, conn *mysql.Conn, stmt *query.Select) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return mysql.NewResult(), nil
}