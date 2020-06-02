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
}

// NewMemStore returns an in-memory storeinstance.
func NewMemStore() *MemStore {
	store := &MemStore{}

	return store
}

// CreateDatabase should handle a CREATE database statement.
func (store *MemStore) CreateDatabase(ctx context.Context, stmt *query.DBDDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// AlterDatabase should handle a ALTER database statement.
func (store *MemStore) AlterDatabase(ctx context.Context, stmt *query.DBDDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// DropDatabase should handle a DROP database statement.
func (store *MemStore) DropDatabase(ctx context.Context, stmt *query.DBDDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// CreateTable should handle a CREATE table statement.
func (store *MemStore) CreateTable(ctx context.Context, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// AlterTable should handle a ALTER table statement.
func (store *MemStore) AlterTable(ctx context.Context, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// DropTable should handle a DROP table statement.
func (store *MemStore) DropTable(ctx context.Context, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// RenameTable should handle a RENAME table statement.
func (store *MemStore) RenameTable(ctx context.Context, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// TruncateTable should handle a TRUNCATE table statement.
func (store *MemStore) TruncateTable(ctx context.Context, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// AnalyzeTable should handle a ANALYZE table statement.
func (store *MemStore) AnalyzeTable(ctx context.Context, stmt *query.DDL) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// Insert should handle a INSERT statement.
func (store *MemStore) Insert(ctx context.Context, stmt *query.Insert) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// Update should handle a UPDATE statement.
func (store *MemStore) Update(ctx context.Context, stmt *query.Update) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// Delete should handle a DELETE statement.
func (store *MemStore) Delete(ctx context.Context, stmt *query.Delete) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// Select should handle a SELECT statement.
func (store *MemStore) Select(ctx context.Context, stmt *query.Select) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}
