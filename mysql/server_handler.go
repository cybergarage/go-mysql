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
	"fmt"
	"go-mysql/mysql/query"

	"vitess.io/vitess/go/mysql"
)

// NewConnection is called when a connection is created.
func (server *Server) NewConnection(c *mysql.Conn) {
}

// ConnectionClosed is called when a connection is closed.
func (server *Server) ConnectionClosed(c *mysql.Conn) {
}

// ComInitDB is called once at the beginning to set db name, and subsequently for every ComInitDB event.
func (server *Server) ComInitDB(c *mysql.Conn, schemaName string) {
	fmt.Printf("%v schema (%s)\n", c, schemaName)
}

// ComQuery is called when a connection receives a query.
func (server *Server) ComQuery(c *mysql.Conn, q string, callback func(*Result) error) error {
	fmt.Printf("%v query (%s)\n", c, q)

	parser := query.NewParser()
	stmt, err := parser.Parse(q)
	if err != nil {
		return err
	}

	res := &Result{
		Rows: [][]Value{},
	}

	ctx := context.Background()
	conn := NewConnWithConn(c)

	executor := server.QueryExecutor
	if executor != nil {
		switch v := stmt.(type) {
		case (*query.DBDDL):
			switch v.Action {
			case "create":
				res, err = executor.CreateDatabase(ctx, conn, v)
			case "drop":
				res, err = executor.DropDatabase(ctx, conn, v)
			case "alter":
				res, err = executor.AlterDatabase(ctx, conn, v)
			}
		case (*query.DDL):
			switch v.Action {
			case "create":
				res, err = executor.CreateTable(ctx, conn, v)
			case "drop":
				res, err = executor.DropTable(ctx, conn, v)
			case "alter":
				res, err = executor.AlterTable(ctx, conn, v)
			case "rename":
				res, err = executor.RenameTable(ctx, conn, v)
			case "truncate":
				res, err = executor.TruncateTable(ctx, conn, v)
			case "analyze":
				res, err = executor.AnalyzeTable(ctx, conn, v)
			}
		case (*query.Insert):
			res, err = executor.Insert(ctx, conn, v)
		case (*query.Select):
			res, err = executor.Select(ctx, conn, v)
		case (*query.Update):
			res, err = executor.Update(ctx, conn, v)
		case (*query.Delete):
			res, err = executor.Delete(ctx, conn, v)
		}
	}
	err = callback(res)
	return err
}

// ComPrepare is called when a connection receives a prepared statement query.
func (server *Server) ComPrepare(c *mysql.Conn, query string) ([]*Field, error) {
	return nil, nil
}

// ComStmtExecute is called when a connection receives a statement execute query.
func (server *Server) ComStmtExecute(c *mysql.Conn, prepare *PrepareData, callback func(*Result) error) error {
	return nil
}

// WarningCount is called at the end of each query to obtain the value to be returned to the client in the EOF packet.
func (server *Server) WarningCount(c *mysql.Conn) uint16 {
	return 0
}

// ComResetConnection is called when the connection is reseted.
func (server *Server) ComResetConnection(c *mysql.Conn) {
}
