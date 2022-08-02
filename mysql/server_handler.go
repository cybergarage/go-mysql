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
	"context"

	"github.com/cybergarage/go-mysql/mysql/query"

	"github.com/cybergarage/go-logger/log"

	vitess "vitess.io/vitess/go/mysql"
)

// NewConnection is called when a connection is created.
func (server *Server) NewConnection(c *vitess.Conn) {
	log.Debug("NewConnection %d\n", c.ConnectionID)
	server.AddConn(NewConnWithConn(c))
}

// ConnectionClosed is called when a connection is closed.
func (server *Server) ConnectionClosed(c *vitess.Conn) {
	log.Debug("ConnectionClosed %d\n", c.ConnectionID)
	server.DeleteConnByUID(c.ConnectionID)
}

// ComInitDB is called once at the beginning to set db name, and subsequently for every ComInitDB event.
func (server *Server) ComInitDB(c *vitess.Conn, schemaName string) {
	log.Debug("ComInitDB %v %d schema (%s)\n", c, c.ConnectionID, schemaName)
	conn, ok := server.GetConnByUID(c.ConnectionID)
	if ok {
		conn.Database = schemaName
	}
}

// ComQuery is called when a connection receives a query.
func (server *Server) ComQuery(c *vitess.Conn, q string, callback func(*Result) error) error {
	parser := query.NewParser()
	stmt, err := parser.Parse(q)
	if err != nil {
		return err
	}

	res := &Result{
		Rows: [][]Value{},
	}

	ctx := context.Background()
	conn, ok := server.GetConnByUID(c.ConnectionID)
	if !ok {
		conn = NewConnWithConn(c)
	}

	log.Debug("ComQuery %v %s query (%s)\n", conn, conn.Database, q)

	executor := server.QueryExecutor
	if executor != nil {
		switch v := stmt.(type) {
		case query.DDL:
			switch v.GetAction() {
			case query.CreateDDLAction:
				res, err = executor.CreateTable(ctx, conn, v)
			case query.DropDDLAction:
				res, err = executor.DropTable(ctx, conn, v)
			case query.AlterDDLAction:
				res, err = executor.AlterTable(ctx, conn, v)
			case query.RenameDDLAction:
				res, err = executor.RenameTable(ctx, conn, v)
			case query.TruncateDDLAction:
				res, err = executor.TruncateTable(ctx, conn, v)
			}
		// case (*query.Show):
		// 	switch v.Type {
		// 	case "DATABASES":
		// 		res, err = executor.ShowDatabases(ctx, conn)
		// 	case "TABLES":
		// 		res, err = executor.ShowTables(ctx, conn, conn.Database)
		// 	}
		case (*query.Insert):
			res, err = executor.Insert(ctx, conn, v)
		case (*query.Select):
			res, err = executor.Select(ctx, conn, v)
		case (*query.Update):
			res, err = executor.Update(ctx, conn, v)
		case (*query.Delete):
			res, err = executor.Delete(ctx, conn, v)
		case (*query.Use):
			conn.Database = v.DBName.String()
		default:
		}
	}

	if err != nil && res == nil {
		res = NewResult()
	}

	err = callback(res)
	return err
}

// ComPrepare is called when a connection receives a prepared statement query.
func (server *Server) ComPrepare(c *vitess.Conn, query string, bindVars map[string]*BindVariable) ([]*Field, error) {
	return nil, nil
}

// ComStmtExecute is called when a connection receives a statement execute query.
func (server *Server) ComStmtExecute(c *vitess.Conn, prepare *PrepareData, callback func(*Result) error) error {
	return nil
}

// WarningCount is called at the end of each query to obtain the value to be returned to the client in the EOF packet.
func (server *Server) WarningCount(c *vitess.Conn) uint16 {
	return 0
}

// ComResetConnection is called when the connection is reseted.
func (server *Server) ComResetConnection(c *vitess.Conn) {
}
