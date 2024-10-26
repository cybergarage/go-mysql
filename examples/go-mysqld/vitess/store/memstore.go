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

package store

import (
	"fmt"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysql/errors"
	"github.com/cybergarage/go-mysql/mysql/plugins/vitess"
	"github.com/cybergarage/go-mysql/mysql/plugins/vitess/query"
)

type MemStore struct {
	*vitess.BaseExecutor
	Databases
}

// NewMemStore returns an in-memory storeinstance.
func NewMemStore() *MemStore {
	store := &MemStore{
		BaseExecutor: vitess.NewBaseExecutor(),
		Databases:    NewDatabases(),
	}
	return store
}

// Begin should handle a BEGIN statement.
func (store *MemStore) Begin(conn vitess.Conn, stmt *query.Begin) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	return vitess.NewResult(), nil
}

// Commit should handle a COMMIT statement.
func (store *MemStore) Commit(conn vitess.Conn, stmt *query.Commit) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	return vitess.NewResult(), nil
}

// Rollback should handle a ROLLBACK statement.
func (store *MemStore) Rollback(conn vitess.Conn, stmt *query.Rollback) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	return vitess.NewResult(), nil
}

// CreateDatabase should handle a CREATE database statement.
func (store *MemStore) CreateDatabase(conn vitess.Conn, stmt *query.Database) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	dbName := stmt.Name()
	_, ok := store.LookupDatabase(dbName)
	if ok {
		if stmt.IfNotExists() {
			return vitess.NewResult(), nil
		}
		return vitess.NewResult(), errors.NewErrDatabaseNotExist(dbName)
	}
	err := store.AddDatabase(NewDatabaseWithName(dbName))
	if err != nil {
		return vitess.NewResult(), err
	}
	return vitess.NewResult(), nil
}

// AlterDatabase should handle a ALTER database statement.
func (store *MemStore) AlterDatabase(conn vitess.Conn, stmt *query.Database) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	return vitess.NewResult(), nil
}

// DropDatabase should handle a DROP database statement.
func (store *MemStore) DropDatabase(conn vitess.Conn, stmt *query.Database) (*vitess.Result, error) {
	dbName := stmt.Name()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, errors.NewErrDatabaseNotExist(dbName)
	}

	if !store.Databases.DropDatabase(db) {
		return nil, fmt.Errorf("%s could not deleted", db.Name())
	}

	return vitess.NewResult(), nil
}

// CreateTable should handle a CREATE table statement.
func (store *MemStore) CreateTable(conn vitess.Conn, stmt *query.Schema) (*vitess.Result, error) {
	dbName := conn.Database()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, errors.NewErrDatabaseNotExist(dbName)
	}
	tableName := stmt.TableName()
	_, ok = db.LookupTable(tableName)
	if !ok {
		table := NewTableWith(tableName, stmt)
		db.AddTable(table)
	} else {
		if !stmt.GetIfNotExists() {
			return vitess.NewResult(), errors.NewCollectionExists(tableName)
		}
	}
	return vitess.NewResult(), nil
}

// AlterTable should handle a ALTER table statement.
func (store *MemStore) AlterTable(conn vitess.Conn, stmt *query.Schema) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	return vitess.NewResult(), nil
}

// DropTable should handle a DROP table statement.
func (store *MemStore) DropTable(conn vitess.Conn, stmt *query.Schema) (*vitess.Result, error) {
	dbName := conn.Database()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, errors.NewErrDatabaseNotExist(dbName)
	}
	tableName := stmt.TableName()
	table, ok := db.LookupTable(tableName)
	if !ok {
		return vitess.NewResult(), nil
	}

	if !db.DropTable(table) {
		return nil, fmt.Errorf("%s could not deleted", table.TableName())
	}

	return vitess.NewResult(), nil
}

// RenameTable should handle a RENAME table statement.
func (store *MemStore) RenameTable(conn vitess.Conn, stmt *query.Schema) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	return vitess.NewResult(), nil
}

// TruncateTable should handle a TRUNCATE table statement.
func (store *MemStore) TruncateTable(conn vitess.Conn, stmt *query.Schema) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	return vitess.NewResult(), nil
}

// Insert should handle a INSERT statement.
func (store *MemStore) Insert(conn vitess.Conn, stmt *query.Insert) (*vitess.Result, error) {
	log.Debugf("%v", stmt)
	dbName := conn.Database()
	tableName := stmt.TableName()
	table, ok := store.LookupTableWithDatabase(dbName, tableName)
	if !ok {
		return vitess.NewResult(), errors.NewCollectionNotFound(tableName)
	}

	row, err := query.NewRowWithInsert(stmt)
	if err != nil {
		return nil, err
	}

	if err := table.AddRow(row); err != nil {
		return nil, err
	}

	table.Dump()

	return vitess.NewResultWithRowsAffected(1), nil
}

// Update should handle a UPDATE statement.
func (store *MemStore) Update(conn vitess.Conn, stmt *query.Update) (*vitess.Result, error) {
	log.Debugf("%v", stmt)

	dbName := conn.Database()
	cond := stmt.Where

	database, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, errors.NewErrDatabaseNotExist(dbName)
	}

	nEffectedRows := uint64(0)
	for _, table := range stmt.Tables() {
		tableName, err := table.Name()
		if err != nil {
			return nil, err
		}
		table, ok := database.LookupTable(tableName)
		if !ok {
			return nil, errors.NewCollectionNotFound(tableName)
		}

		columns, err := stmt.Columns()
		if err != nil {
			return nil, err
		}

		nUpdatedRows, err := table.Update(columns, cond)
		if err != nil {
			return nil, err
		}
		nEffectedRows += uint64(nUpdatedRows)
	}

	return vitess.NewResultWithRowsAffected(nEffectedRows), nil
}

// Delete should handle a DELETE statement.
func (store *MemStore) Delete(conn vitess.Conn, stmt *query.Delete) (*vitess.Result, error) {
	dbName := conn.Database()
	cond := stmt.Where

	database, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, errors.NewErrDatabaseNotExist(dbName)
	}

	nEffectedRows := uint64(0)
	for _, table := range stmt.Tables() {
		tableName, err := table.Name()
		if err != nil {
			return nil, err
		}
		table, ok := database.LookupTable(tableName)
		if !ok {
			return nil, errors.NewCollectionNotFound(tableName)
		}

		nDeletedRows, err := table.Delete(cond)
		if err != nil {
			return nil, err
		}
		nEffectedRows += uint64(nDeletedRows)
	}

	return vitess.NewResultWithRowsAffected(nEffectedRows), nil
}

// Select should handle a SELECT statement.
func (store *MemStore) Select(conn vitess.Conn, stmt *query.Select) (*vitess.Result, error) {
	log.Debugf("%v", stmt)

	dbName := conn.Database()
	database, ok := store.LookupDatabase(dbName)
	if !ok {
		return nil, errors.NewErrDatabaseNotExist(dbName)
	}

	// NOTE: Select scans only a first table

	tables := stmt.From()
	tableName, err := tables[0].Name()
	if err != nil {
		return nil, err
	}

	table, ok := database.LookupTable(tableName)
	if !ok {
		// TODO: Support dummy dual table for MySQL connector 5.1.49
		if tableName == "dual" {
			return vitess.NewResult(), nil
		}
		return nil, errors.NewCollectionNotFound(tableName)
	}

	cond := stmt.Where
	matchedRows := table.FindMatchedRows(cond)

	return vitess.NewResultWithRows(database.Database, table.Schema, matchedRows)
}

// ShowDatabases should handle a SHOW DATABASES statement.
func (store *MemStore) ShowDatabases(conn vitess.Conn) (*vitess.Result, error) {
	return vitess.NewResult(), nil
}

// ShowTables should handle a SHOW TABLES statement.
func (store *MemStore) ShowTables(conn vitess.Conn, database string) (*vitess.Result, error) {
	return vitess.NewResult(), nil
}
