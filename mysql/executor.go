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
}

// DMOExecutor defines a executor interface for DMO (Data Manipulation Operations).
type DMOExecutor interface {
	// Insert should handle INSERT queries.
	Insert(context.Context, *query.Insert) (*Result, error)
	// Update should handle UPDATE queries.
	Update(context.Context, *query.Update) (*Result, error)
	// Delete should handle DELETE queries.
	Delete(context.Context, *query.Delete) (*Result, error)
	// Select should handle SELECT queries.
	Select(context.Context, *query.Select) (*Result, error)
}

// DCOExecutor defines a executor interface for DCO (Data Control Operations).
type DCOExecutor interface {
}

// QueryExecutor represents an interface to execute all CQL queries.
type QueryExecutor interface {
	DDOExecutor
	DMOExecutor
	DCOExecutor
}
