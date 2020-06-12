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

// Databases represents a collection of databases.
type Databases map[string]*Database

// NewDatabases returns a databases instance.
func NewDatabases() Databases {
	return Databases{}
}

// AddDatabase adds a specified database.
func (dbs Databases) AddDatabase(db *Database) error {
	dbName := db.GetName()
	dbs[dbName] = db
	return nil
}

// GetDatabase returns a database with the specified name.
func (dbs Databases) GetDatabase(name string) (*Database, bool) {
	ks, ok := dbs[name]
	return ks, ok
}