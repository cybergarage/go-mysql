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
	"go-mysql/mysql/query"
)

// NewConnection is called when a connection is created.
func (server *Server) NewConnection(c *Conn) {
}

// ConnectionClosed is called when a connection is closed.
func (server *Server) ConnectionClosed(c *Conn) {
}

// ComInitDB is called once at the beginning to set db name, and subsequently for every ComInitDB event.
func (server *Server) ComInitDB(c *Conn, schemaName string) {
}

// ComQuery is called when a connection receives a query.
func (server *Server) ComQuery(c *Conn, q string, callback func(*Result) error) error {
	parser := query.NewParser()
	stmt, err := parser.Parse(q)
	if err != nil {
		return err
	}

	res := &Result{
		Rows: [][]Value{},
	}

	ctx := context.Background()

	executor := server.QueryExecutor
	if executor != nil {
		switch v := stmt.(type) {
		case (*query.DBDDL):
			switch v.Action {
			case "create":
				res, err = executor.CreateDatabase(ctx, v)
			case "drop":
				res, err = executor.DropDatabase(ctx, v)
			case "alter":
				res, err = executor.AlterDatabase(ctx, v)
			}
		case (*query.DDL):
			switch v.Action {
			case "create":
				res, err = executor.CreateTable(ctx, v)
			case "drop":
				res, err = executor.DropTable(ctx, v)
			case "alter":
				res, err = executor.AlterTable(ctx, v)
			case "rename":
				res, err = executor.RenameTable(ctx, v)
			case "truncate":
				res, err = executor.TruncateTable(ctx, v)
			case "analyze":
				res, err = executor.AnalyzeTable(ctx, v)
			}
		case (*query.Insert):
			res, err = executor.Insert(ctx, v)
		case (*query.Select):
			res, err = executor.Select(ctx, v)
		case (*query.Update):
			res, err = executor.Update(ctx, v)
		case (*query.Delete):
			res, err = executor.Delete(ctx, v)
		}
	}
	err = callback(res)
	return err
}

// ComPrepare is called when a connection receives a prepared statement query.
func (server *Server) ComPrepare(c *Conn, query string) ([]*Field, error) {
	return nil, nil
}

// ComStmtExecute is called when a connection receives a statement execute query.
func (server *Server) ComStmtExecute(c *Conn, prepare *PrepareData, callback func(*Result) error) error {
	return nil
}

// WarningCount is called at the end of each query to obtain the value to be returned to the client in the EOF packet.
func (server *Server) WarningCount(c *Conn) uint16 {
	return 0
}

// ComResetConnection is called when the connection is reseted.
func (server *Server) ComResetConnection(c *Conn) {
}
