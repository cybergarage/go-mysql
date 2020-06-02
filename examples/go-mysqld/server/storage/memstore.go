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

// Insert should handle INSERT queries.
func (store *MemStore) Insert(ctx context.Context, stmt *query.Insert) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// Update should handle UPDATE queries.
func (store *MemStore) Update(ctx context.Context, stmt *query.Update) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// Delete should handle DELETE queries.
func (store *MemStore) Delete(ctx context.Context, stmt *query.Delete) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}

// Select should handle SELECT queries.
func (store *MemStore) Select(ctx context.Context, stmt *query.Select) (*mysql.Result, error) {
	fmt.Printf("%v\n", stmt)
	return nil, nil
}
