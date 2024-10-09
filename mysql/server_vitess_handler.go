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
	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-tracing/tracer"
	vitessmy "vitess.io/vitess/go/mysql"
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// NewConnection is called when a connection is created.
func (server *VitessServer) NewConnection(c *vitessmy.Conn) {
	conn := NewConnWith(tracer.NewNullTracer().StartSpan(""), c)
	server.AddConn(conn)
}

// ConnectionClosed is called when a connection is closed.
func (server *VitessServer) ConnectionClosed(c *vitessmy.Conn) {
	server.DeleteConnByUID(c.ConnectionID)
}

// ComInitDB is called once at the beginning to set db name, and subsequently for every ComInitDB event.
func (server *VitessServer) ComInitDB(c *vitessmy.Conn, dbName string) {
	conn, ok := server.GetConnByUID(c.ConnectionID)
	if ok {
		conn.SetDatabase(dbName)
	}
}

// ComQuery is called when a connection receives a query.
// nolint: exhaustive
func (server *VitessServer) ComQuery(c *vitessmy.Conn, q string, callback func(*Result) error) error {
	spanCtx := server.Tracer.StartSpan(PackageName)
	defer spanCtx.Span().Finish()

	conn, ok := server.GetConnByUID(c.ConnectionID)
	if ok {
		conn.SetSpanContext(spanCtx)
	} else {
		conn = NewConnWith(spanCtx, c)
	}

	conn.StartSpan("parse")
	parser := query.NewParser()
	stmt, err := parser.Parse(q)
	conn.FinishSpan()
	if err != nil {
		return err
	}

	var res *Result

	executor := server.queryExecutor
	if executor != nil {
		switch v := stmt.(type) {
		case (query.DBDDL):
			switch v := v.(type) {
			case *vitesssp.CreateDatabase:
				conn.StartSpan("CreateDatabase")
				defer conn.FinishSpan()
				res, err = executor.CreateDatabase(conn, query.NewDatabaseWithDBDDL(v))
			case *vitesssp.DropDatabase:
				conn.StartSpan("DropDatabase")
				defer conn.FinishSpan()
				res, err = executor.DropDatabase(conn, query.NewDatabaseWithDBDDL(v))
			case *vitesssp.AlterDatabase:
				conn.StartSpan("AlterDatabase")
				defer conn.FinishSpan()
				res, err = executor.AlterDatabase(conn, query.NewDatabaseWithDBDDL(v))
			}
		case (query.DDL):
			switch v := v.(type) {
			case *vitesssp.CreateTable:
				conn.StartSpan("CreateTable")
				defer conn.FinishSpan()
				res, err = executor.CreateTable(conn, query.NewSchemaWithDDL(v))
			case *vitesssp.DropTable:
				conn.StartSpan("DropTable")
				defer conn.FinishSpan()
				res, err = executor.DropTable(conn, query.NewSchemaWithDDL(v))
			case *vitesssp.AlterTable:
				conn.StartSpan("AlterTable")
				defer conn.FinishSpan()
				res, err = executor.AlterTable(conn, query.NewSchemaWithDDL(v))
			case *vitesssp.RenameTable:
				conn.StartSpan("RenameTable")
				defer conn.FinishSpan()
				res, err = executor.RenameTable(conn, query.NewSchemaWithDDL(v))
			case *vitesssp.TruncateTable:
				conn.StartSpan("TruncateTable")
				defer conn.FinishSpan()
				res, err = executor.TruncateTable(conn, query.NewSchemaWithDDL(v))
			}
		case (*query.Show):
			/* TODO: v.Type is deprecated
			switch v.Type {
			case "DATABASES":
				res, err = executor.ShowDatabases(conn)
			case "TABLES":
				res, err = executor.ShowTables(conn, conn.Database())
			}
			*/
		case (*vitesssp.Begin):
			conn.StartSpan("Begin")
			defer conn.FinishSpan()
			res, err = executor.Begin(conn, v)
		case (*vitesssp.Commit):
			conn.StartSpan("Commit")
			defer conn.FinishSpan()
			res, err = executor.Commit(conn, v)
		case (*vitesssp.Rollback):
			conn.StartSpan("Rollback")
			defer conn.FinishSpan()
			res, err = executor.Rollback(conn, v)
		case (*vitesssp.Insert):
			conn.StartSpan("Insert")
			defer conn.FinishSpan()
			res, err = executor.Insert(conn, query.NewInsertWithInsert(v))
		case (*vitesssp.Select):
			conn.StartSpan("Select")
			defer conn.FinishSpan()
			res, err = executor.Select(conn, query.NewSelectWithSelect(v))
		case (*vitesssp.Update):
			conn.StartSpan("Update")
			defer conn.FinishSpan()
			res, err = executor.Update(conn, query.NewUpdateWithUpdate(v))
		case (*vitesssp.Delete):
			conn.StartSpan("Delete")
			defer conn.FinishSpan()
			res, err = executor.Delete(conn, query.NewDeleteWithDelete(v))
		case (*query.Use):
			conn.StartSpan("Use")
			defer conn.FinishSpan()
			conn.SetDatabase(v.DBName.String())
		}
	}

	if err != nil {
		return err
	}

	if res == nil {
		res = &Result{ // nolint: exhaustruct
			Rows: [][]Value{},
		}
	}

	conn.StartSpan("response")
	err = callback(res)
	conn.FinishSpan()

	return err
}

// ComPrepare is called when a connection receives a prepared statement query.
func (server *VitessServer) ComPrepare(c *vitessmy.Conn, query string, bindVars map[string]*BindVariable) ([]*Field, error) {
	return nil, nil
}

// ComStmtExecute is called when a connection receives a statement execute query.
func (server *VitessServer) ComStmtExecute(c *vitessmy.Conn, prepare *PrepareData, callback func(*Result) error) error {
	return nil
}

// WarningCount is called at the end of each query to obtain the value to be returned to the client in the EOF packet.
func (server *VitessServer) WarningCount(c *vitessmy.Conn) uint16 {
	return 0
}

// ComResetConnection is called when the connection is reseted.
func (server *VitessServer) ComResetConnection(c *vitessmy.Conn) {
}
