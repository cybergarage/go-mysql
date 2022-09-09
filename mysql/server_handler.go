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

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysql/query"

	vitessmy "vitess.io/vitess/go/mysql"
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// NewConnection is called when a connection is created.
func (server *Server) NewConnection(c *vitessmy.Conn) {
	log.Debugf("NewConnection %d", c.ConnectionID)
	server.AddConn(NewConnWithConn(c))
}

// ConnectionClosed is called when a connection is closed.
func (server *Server) ConnectionClosed(c *vitessmy.Conn) {
	log.Debugf("ConnectionClosed %d", c.ConnectionID)
	server.DeleteConnByUID(c.ConnectionID)
}

// ComInitDB is called once at the beginning to set db name, and subsequently for every ComInitDB event.
func (server *Server) ComInitDB(c *vitessmy.Conn, schemaName string) {
	log.Debugf("ComInitDB %v %d schema (%s)", c, c.ConnectionID, schemaName)
	conn, ok := server.GetConnByUID(c.ConnectionID)
	if ok {
		conn.Database = schemaName
	}
}

// ComQuery is called when a connection receives a query.
func (server *Server) ComQuery(c *vitessmy.Conn, q string, callback func(*Result) error) error {
	parser := query.NewParser()
	stmt, err := parser.Parse(q)
	if err != nil {
		return err
	}

	ctx := context.Background()
	conn, ok := server.GetConnByUID(c.ConnectionID)
	if !ok {
		conn = NewConnWithConn(c)
	}

	log.Debugf("ComQuery %v %s query (%s)", conn, conn.Database, q)

	var res *Result

	executor := server.queryExecutor
	if executor != nil {
		switch v := stmt.(type) {
		case (*query.DBDDL):
			switch v.Action {
			case vitesssp.CreateDBDDLAction:
				res, err = executor.CreateDatabase(ctx, conn, query.NewDatabaseWithDBDDL(v))
			case vitesssp.DropDBDDLAction:
				res, err = executor.DropDatabase(ctx, conn, query.NewDatabaseWithDBDDL(v))
			case vitesssp.AlterDBDDLAction:
				res, err = executor.AlterDatabase(ctx, conn, query.NewDatabaseWithDBDDL(v))
			}
		case (*query.DDL):
			switch v.Action {
			case vitesssp.CreateDDLAction:
				res, err = executor.CreateTable(ctx, conn, query.NewSchemaWithDDL(v))
			case vitesssp.DropDDLAction:
				res, err = executor.DropTable(ctx, conn, query.NewSchemaWithDDL(v))
			case vitesssp.AlterDDLAction:
				res, err = executor.AlterTable(ctx, conn, query.NewSchemaWithDDL(v))
			case vitesssp.RenameDDLAction:
				res, err = executor.RenameTable(ctx, conn, query.NewSchemaWithDDL(v))
			case vitesssp.TruncateDDLAction:
				res, err = executor.TruncateTable(ctx, conn, query.NewSchemaWithDDL(v))
			}
		case (*query.Show):
			switch v.Type {
			case "DATABASES":
				res, err = executor.ShowDatabases(ctx, conn)
			case "TABLES":
				res, err = executor.ShowTables(ctx, conn, conn.Database)
			}
		case (*vitesssp.Insert):
			res, err = executor.Insert(ctx, conn, query.NewInsertWithInsert(v))
		case (*vitesssp.Select):
			res, err = executor.Select(ctx, conn, query.NewSelectWithSelect(v))
		case (*vitesssp.Update):
			res, err = executor.Update(ctx, conn, query.NewUpdateWithUpdate(v))
		case (*vitesssp.Delete):
			res, err = executor.Delete(ctx, conn, query.NewDeleteWithDelete(v))
		case (*query.Use):
			conn.Database = v.DBName.String()
		}
	}

	if err != nil || res == nil {
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
