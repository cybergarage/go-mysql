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
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysql/errors"
	"github.com/cybergarage/go-mysql/mysql/net"
	"github.com/cybergarage/go-mysql/mysql/query"
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

// Begin should handle a BEGIN statement.
func (store *MemStore) Begin(conn net.Conn, stmt query.Begin) error {
	log.Debugf("%v", stmt)
	return errors.ErrNotImplemented
}

// Commit should handle a COMMIT statement.
func (store *MemStore) Commit(conn net.Conn, stmt query.Commit) error {
	log.Debugf("%v", stmt)
	return errors.ErrNotImplemented
}

// Rollback should handle a ROLLBACK statement.
func (store *MemStore) Rollback(conn net.Conn, stmt query.Rollback) error {
	log.Debugf("%v", stmt)
	return errors.ErrNotImplemented
}

// CreateDatabase should handle a CREATE database statement.
func (store *MemStore) CreateDatabase(conn net.Conn, stmt query.CreateDatabase) error {
	log.Debugf("%v", stmt)
	dbName := stmt.DatabaseName()
	_, ok := store.LookupDatabase(dbName)
	if ok {
		if stmt.IfNotExists() {
			return nil
		}
		return errors.NewDatabaseExists(dbName)
	}
	return store.AddDatabase(NewDatabaseWithName(dbName))
}

// AlterDatabase should handle a ALTER database statement.
func (store *MemStore) AlterDatabase(conn net.Conn, stmt query.AlterDatabase) error {
	log.Debugf("%v", stmt)
	return errors.ErrNotImplemented
}

// DropDatabase should handle a DROP database statement.
func (store *MemStore) DropDatabase(conn net.Conn, stmt query.DropDatabase) error {
	dbName := stmt.DatabaseName()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		if stmt.IfExists() {
			return nil
		}
		return errors.NewDatabaseNotFound(dbName)
	}
	return store.Databases.DropDatabase(db)
}

// CreateTable should handle a CREATE table statement.
func (store *MemStore) CreateTable(conn net.Conn, stmt query.CreateTable) error {
	dbName := conn.Database()
	db, ok := store.LookupDatabase(dbName)
	if !ok {
		return errors.NewDatabaseNotFound(dbName)
	}
	tableName := stmt.TableName()
	_, ok = db.LookupTable(tableName)
	if !ok {
		/*
			table := NewTableWith(tableName, stmt)
			db.AddTable(table)
		*/
	} else {
		if !stmt.IfNotExists() {
			return errors.NewCollectionExists(tableName)
		}
	}
	return errors.ErrNotImplemented
}

// AlterTable should handle a ALTER table statement.
func (store *MemStore) AlterTable(conn net.Conn, stmt query.AlterTable) error {
	// log.Debugf("%v", stmt)
	return errors.ErrNotImplemented
}

// DropTable should handle a DROP table statement.
func (store *MemStore) DropTable(conn net.Conn, stmt query.DropTable) error {
	/*
		dbName := conn.Database()
		db, ok := store.LookupDatabase(dbName)
		if !ok {
			return nil, errors.NewDatabaseNotFound(dbName)
		}
		tableName := stmt.TableName()
		table, ok := db.LookupTable(tableName)
		if !ok {
			return errors.ErrNotImplemented
		}

		if !db.DropTable(table) {
			return nil, fmt.Errorf("%s could not deleted", table.TableName())
		}
	*/
	return errors.ErrNotImplemented
}

// Insert should handle a INSERT statement.
func (store *MemStore) Insert(conn net.Conn, stmt query.Insert) error {
	/*
		log.Debugf("%v", stmt)
		dbName := conn.Database()
		tableName := stmt.TableName()
		table, ok := store.GetTableWithDatabase(dbName, tableName)
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
	*/
	return errors.ErrNotImplemented
}

// Update should handle a UPDATE statement.
func (store *MemStore) Update(conn net.Conn, stmt query.Update) (query.ResultSet, error) {
	/*
		log.Debugf("%v", stmt)

		dbName := conn.Database()
		cond := stmt.Where

		database, ok := store.LookupDatabase(dbName)
		if !ok {
			return nil, errors.NewDatabaseNotFound(dbName)
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
	*/
	return nil, errors.ErrNotImplemented

}

// Delete should handle a DELETE statement.
func (store *MemStore) Delete(conn net.Conn, stmt query.Delete) (query.ResultSet, error) {
	/*
	   dbName := conn.Database()
	   cond := stmt.Where

	   database, ok := store.LookupDatabase(dbName)

	   	if !ok {
	   		return nil, errors.NewDatabaseNotFound(dbName)
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
	*/
	return nil, errors.ErrNotImplemented
}

// Select should handle a SELECT statement.
func (store *MemStore) Select(conn net.Conn, stmt query.Select) (query.ResultSet, error) {
	/*
		log.Debugf("%v", stmt)

		dbName := conn.Database()
		database, ok := store.LookupDatabase(dbName)
		if !ok {
			return nil, errors.NewDatabaseNotFound(dbName)
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
				return errors.ErrNotImplemented
			}
			return nil, errors.NewCollectionNotFound(tableName)
		}

		cond := stmt.Where
		matchedRows := table.FindMatchedRows(cond)

		return vitess.NewResultWithRows(database.Database, table.Schema, matchedRows)
	*/
	return nil, errors.ErrNotImplemented
}

// ParserError should handle a parser error.
func (store *MemStore) ParserError(conn net.Conn, q string, err error) error {
	log.Debugf("%v", err)
	return err
}
