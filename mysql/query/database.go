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

// Database represents a destination or source database of query.
type Database struct {
	value       string
	ifNotExists bool
	ifExists    bool
}

// NewDatabaseWithName returns a new database with the specified string.
func NewDatabaseWithName(name string) *Database {
	db := &Database{
		value:       name,
		ifNotExists: false,
		ifExists:    false,
	}
	return db
}

// NewDatabaseWithDBDDL returns a new database with the specified DBDDL.
func NewDatabaseWithDBDDL(dbddl *DBDDL) (*Database, error) {
	db := NewDatabaseWithName(dbddl.DBName)
	db.ifNotExists = dbddl.IfNotExists
	db.ifExists = dbddl.IfExists
	return db, nil
}

// NewDatabase returns a new database.
func NewDatabase() *Database {
	return NewDatabaseWithName("")
}

// Name returns the database name.
func (db *Database) Name() string {
	return db.value
}

// IfNotExists returns true when the IF NOT EXISTS option is enabled.
func (db *Database) IfNotExists() bool {
	return db.ifNotExists
}

// IfExists returns true when the IF EXISTS option is enabled.
func (db *Database) IfExists() bool {
	return db.ifExists
}
