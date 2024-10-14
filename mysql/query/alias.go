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

package query

import (
	"github.com/cybergarage/go-sqlparser/sql"
)

// CreateDatabase represents a "CREATE DATABASE" statement interface.
type CreateDatabase = sql.CreateDatabase

// CreateTable represents a "CREATE TABLE" statement interface.
type CreateTable = sql.CreateTable

// AlterDatabase represents a "ALTER DATABASE" statement interface.
type AlterDatabase = sql.AlterDatabase

// AlterTable represents a "ALTER TABLE" statement interface.
type AlterTable = sql.AlterTable

// DropTable represents a "DROP TABLE" statement interface.
type DropDatabase = sql.DropDatabase

// DropTable represents a "DROP TABLE" statement interface.
type DropTable = sql.DropTable

// Insert represents a "INSERT" statement interface.
type Insert = sql.Insert

// Select represents a "SELECT" statement interface.
type Select = sql.Select

// Update represents a "UPDATE" statement interface.
type Update = sql.Update

// Delete represents a "DELETE" statement interface.
type Delete = sql.Delete

// Begin represents a "BEGIN" statement interface.
type Begin = sql.Begin

// Commit represents a "COMMIT" statement interface.
type Commit = sql.Commit

// Rollback represents a "ROLLBACK" statement interface.
type Rollback = sql.Rollback

// Copy represents a "COPY" statement interface.
type Copy = sql.Copy

// Vacuum represents a "VACUUM" statement interface.
type Vacuum = sql.Vacuum

// Truncate represents a "TRUNCATE" statement interface.
type Truncate = sql.Truncate
