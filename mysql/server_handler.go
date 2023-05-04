// Copyright (C) 2020 Satoshi Konno. All rights reserved.
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
	vitessmy "vitess.io/vitess/go/mysql"
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// NewConnection is called when a connection is created.
func (server *Server) NewConnection(c *vitessmy.Conn) {
	server.AddConn(NewConnWithConn(c))
}

// ConnectionClosed is called when a connection is closed.
func (server *Server) ConnectionClosed(c *vitessmy.Conn) {
	server.DeleteConnByUID(c.ConnectionID)
}

// ComInitDB is called once at the beginning to set db name, and subsequently for every ComInitDB event.
func (server *Server) ComInitDB(c *vitessmy.Conn, dbName string) {
	conn, ok := server.GetConnByUID(c.ConnectionID)
	if ok {
		conn.SetDatabase(dbName)
	}
}

// ComQuery is called when a connection receives a query.
// nolint: exhaustive
func (server *Server) ComQuery(c *vitessmy.Conn, q string, callback func(*Result) error) error {
	parser := query.NewParser()
	stmt, err := parser.Parse(q)
	if err != nil {
		return err
	}

	conn, ok := server.GetConnByUID(c.ConnectionID)
	if !ok {
		conn = NewConnWithConn(c)
	}
	span := server.Tracer.StartSpan(PackageName)
	conn.SetSpanContext(span)
	defer span.Span().Finish()

	var res *Result

	executor := server.queryExecutor
	if executor != nil {
		switch v := stmt.(type) {
		case (query.DBDDL):
			switch v := v.(type) {
			case *vitesssp.CreateDatabase:
				s := conn.SpanContext().Span().StartSpan("CreateDatabase")
				defer s.Span().Finish()
				res, err = executor.CreateDatabase(conn, query.NewDatabaseWithDBDDL(v))
			case *vitesssp.DropDatabase:
				s := conn.SpanContext().Span().StartSpan("DropDatabase")
				defer s.Span().Finish()
				res, err = executor.DropDatabase(conn, query.NewDatabaseWithDBDDL(v))
			case *vitesssp.AlterDatabase:
				s := conn.SpanContext().Span().StartSpan("AlterDatabase")
				defer s.Span().Finish()
				res, err = executor.AlterDatabase(conn, query.NewDatabaseWithDBDDL(v))
			}
		case (query.DDL):
			switch v := v.(type) {
			case *vitesssp.CreateTable:
				s := conn.SpanContext().Span().StartSpan("CreateTable")
				defer s.Span().Finish()
				res, err = executor.CreateTable(conn, query.NewSchemaWithDDL(v))
			case *vitesssp.DropTable:
				s := conn.SpanContext().Span().StartSpan("DropTable")
				defer s.Span().Finish()
				res, err = executor.DropTable(conn, query.NewSchemaWithDDL(v))
			case *vitesssp.AlterTable:
				s := conn.SpanContext().Span().StartSpan("AlterTable")
				defer s.Span().Finish()
				res, err = executor.AlterTable(conn, query.NewSchemaWithDDL(v))
			case *vitesssp.RenameTable:
				s := conn.SpanContext().Span().StartSpan("RenameTable")
				defer s.Span().Finish()
				res, err = executor.RenameTable(conn, query.NewSchemaWithDDL(v))
			case *vitesssp.TruncateTable:
				s := conn.SpanContext().Span().StartSpan("TruncateTable")
				defer s.Span().Finish()
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
		case (*vitesssp.Insert):
			s := conn.SpanContext().Span().StartSpan("Insert")
			defer s.Span().Finish()
			res, err = executor.Insert(conn, query.NewInsertWithInsert(v))
		case (*vitesssp.Select):
			s := conn.SpanContext().Span().StartSpan("Select")
			defer s.Span().Finish()
			res, err = executor.Select(conn, query.NewSelectWithSelect(v))
		case (*vitesssp.Update):
			s := conn.SpanContext().Span().StartSpan("Update")
			defer s.Span().Finish()
			res, err = executor.Update(conn, query.NewUpdateWithUpdate(v))
		case (*vitesssp.Delete):
			s := conn.SpanContext().Span().StartSpan("Delete")
			defer s.Span().Finish()
			res, err = executor.Delete(conn, query.NewDeleteWithDelete(v))
		case (*query.Use):
			s := conn.SpanContext().Span().StartSpan("Use")
			defer s.Span().Finish()
			conn.SetDatabase(v.DBName.String())
		}
	}

	if err != nil {
		return err
	}

	if res == nil {
		res = &Result{
			Rows: [][]Value{},
		}
	}

	err = callback(res)

	return err
}

// ComPrepare is called when a connection receives a prepared statement query.
func (server *Server) ComPrepare(c *vitessmy.Conn, query string, bindVars map[string]*BindVariable) ([]*Field, error) {
	return nil, nil
}

// ComStmtExecute is called when a connection receives a statement execute query.
func (server *Server) ComStmtExecute(c *vitessmy.Conn, prepare *PrepareData, callback func(*Result) error) error {
	return nil
}

// WarningCount is called at the end of each query to obtain the value to be returned to the client in the EOF packet.
func (server *Server) WarningCount(c *vitessmy.Conn) uint16 {
	return 0
}

// ComResetConnection is called when the connection is reseted.
func (server *Server) ComResetConnection(c *vitessmy.Conn) {
}
