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
	"go-mysql/mysql/log"
	"go-mysql/mysql/query"

	"vitess.io/vitess/go/vt/sqlparser"
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
	log.Debug("%v\n", stmt)
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
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// DropDatabase should handle a DROP database statement.
func (store *MemStore) DropDatabase(ctx context.Context, conn *mysql.Conn, stmt *query.DBDDL) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// CreateTable should handle a CREATE table statement.
func (store *MemStore) CreateTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	dbName := conn.Database
	db, ok := store.GetDatabase(dbName)
	if !ok {
		return mysql.NewResult(), fmt.Errorf(errorDatabaseNotFound, dbName)
	}
	tableName := stmt.Table.Name.String()
	_, ok = db.GetTable(tableName)
	if !ok {
		table := NewTableWithName(tableName)
		db.AddTable(table)
	} else {
		if !stmt.IfExists {
			return mysql.NewResult(), fmt.Errorf(errorTableFound, dbName, tableName)
		}
	}
	return mysql.NewResult(), nil
}

// AlterTable should handle a ALTER table statement.
func (store *MemStore) AlterTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// DropTable should handle a DROP table statement.
func (store *MemStore) DropTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// RenameTable should handle a RENAME table statement.
func (store *MemStore) RenameTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// TruncateTable should handle a TRUNCATE table statement.
func (store *MemStore) TruncateTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// AnalyzeTable should handle a ANALYZE table statement.
func (store *MemStore) AnalyzeTable(ctx context.Context, conn *mysql.Conn, stmt *query.DDL) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// Insert should handle a INSERT statement.
func (store *MemStore) Insert(ctx context.Context, conn *mysql.Conn, stmt *query.Insert) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	dbName := conn.Database
	tableName := stmt.Table.Name.String()
	table, ok := store.GetTableWithDatabase(dbName, tableName)
	if !ok {
		return mysql.NewResult(), fmt.Errorf(errorTableNotFound, dbName, tableName)
	}

	columns := query.NewColumns()

	//queryColumns := stmt.Columns

	rows := stmt.Rows
	log.Debug("%v\n", rows)
	node, _ := rows.(sqlparser.SQLNode)
	log.Debug("%v\n", node)
	queryRows, _ := node.(sqlparser.Values)
	log.Debug("%v\n", queryRows)

	// if len(queryColumns) != len(queryValues) {
	// 	// TODO: Return an aprociate errors
	// 	return mysql.NewResult(), nil
	// }

	// for n, queryColumn := range queryColumns {
	// 	name := queryColumn.String()
	// 	value := queryValues[n]
	// 	column := NewColumnWithNameAndValue(name, value)
	// 	columns.AddColumn(column)
	// 	log.Debug("[%d] %v\n", n, queryColumn)
	// }

	for _, queryRow := range queryRows {
		for _, expr := range queryRow {
			cmpExpr, ok := expr.(*sqlparser.ComparisonExpr)
			if !ok {
				continue
			}
			col, ok := cmpExpr.Left.(*sqlparser.ColName)
			if !ok {
				continue
			}
			val, ok := cmpExpr.Right.(*sqlparser.SQLVal)
			if !ok {
				continue
			}
			column := query.NewColumnWithNameAndValue(col.Name.String(), val.Val)
			columns.AddColumn(column)
		}

		table.AddRow(query.NewRowWithColumns(columns))

		// const (
		// 	StrVal = ValType(iota)
		// 	IntVal
		// 	FloatVal
		// 	HexNum
		// 	HexVal
		// 	ValArg
		// 	BitVal
		// )
	}

	return mysql.NewResult(), nil
}

// Update should handle a UPDATE statement.
func (store *MemStore) Update(ctx context.Context, conn *mysql.Conn, stmt *query.Update) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// Delete should handle a DELETE statement.
func (store *MemStore) Delete(ctx context.Context, conn *mysql.Conn, stmt *query.Delete) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// Select should handle a SELECT statement.
func (store *MemStore) Select(ctx context.Context, conn *mysql.Conn, stmt *query.Select) (*mysql.Result, error) {
	log.Debug("%v\n", stmt)
	return mysql.NewResult(), nil
}

// ShowDatabases should handle a SHOW DATABASES statement.
func (store *MemStore) ShowDatabases(ctx context.Context, conn *mysql.Conn) (*mysql.Result, error) {
	return mysql.NewResult(), nil
}

// ShowTables should handle a SHOW TABLES statement.
func (store *MemStore) ShowTables(ctx context.Context, conn *mysql.Conn, database string) (*mysql.Result, error) {
	return mysql.NewResult(), nil
}
