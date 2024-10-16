// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
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

package v2

import (
	"github.com/cybergarage/go-mysql/mysql/errors"
	"github.com/cybergarage/go-mysql/mysql/plugins"
	"github.com/cybergarage/go-mysql/mysql/protocol"
	"github.com/cybergarage/go-mysql/mysql/query"
)

// Server represents a base executor server.
type Server struct {
	*protocol.Server
	executor plugins.QueryExecutor
}

// NewServer returns a base executor server instance.
func NewServer() *Server {
	server := &Server{
		Server:   protocol.NewServer(),
		executor: nil,
	}
	server.Server.SetCommandHandler(server)
	return server
}

// SetExecutor sets an executor to the server.
func (server *Server) SetExecutor(executor plugins.QueryExecutor) {
	server.executor = executor
}

// Executor returns the executor of the server.
func (server *Server) Executor() plugins.QueryExecutor {
	return server.executor
}

// HandleQuery handles a query.
func (server *Server) HandleQuery(conn protocol.Conn, q *protocol.Query) (protocol.Response, error) {
	if server.executor == nil {
		return nil, errors.ErrNotImplemented
	}
	parser := query.NewParser()
	stmts, err := parser.ParseString(q.Query())
	if err != nil {
		return nil, server.executor.ParserError(conn, q.Query(), err)
	}
	for _, stmt := range stmts {
		var err error
		// err = stmt.Bind(msg.BindParams)
		// if err != nil {
		// 	return nil, err
		// }

		var res protocol.Response

		// nolint: forcetypeassert
		switch stmt.StatementType() {
		case query.BeginStatement:
			stmt := stmt.(query.Begin)
			err = server.executor.Begin(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.CommitStatement:
			stmt := stmt.(query.Commit)
			err = server.executor.Commit(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.RollbackStatement:
			stmt := stmt.(query.Rollback)
			err = server.executor.Rollback(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.CreateDatabaseStatement:
			stmt := stmt.(query.CreateDatabase)
			err = server.executor.CreateDatabase(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.CreateTableStatement:
			stmt := stmt.(query.CreateTable)
			err = server.executor.CreateTable(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.AlterDatabaseStatement:
			stmt := stmt.(query.AlterDatabase)
			err = server.executor.AlterDatabase(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.AlterTableStatement:
			stmt := stmt.(query.AlterTable)
			err = server.executor.AlterTable(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.DropDatabaseStatement:
			stmt := stmt.(query.DropDatabase)
			err = server.executor.DropDatabase(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.DropTableStatement:
			stmt := stmt.(query.DropTable)
			err = server.executor.DropTable(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.InsertStatement:
			stmt := stmt.(query.Insert)
			err = server.executor.Insert(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.SelectStatement:
			stmt := stmt.(query.Select)
			_, err = server.executor.Select(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.UpdateStatement:
			stmt := stmt.(query.Update)
			_, err = server.executor.Update(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		case query.DeleteStatement:
			stmt := stmt.(query.Delete)
			_, err = server.executor.Delete(conn, stmt)
			res, err = protocol.NewResponseWithError(err)
		}

		if err != nil {
			err = conn.ResponseError(err)
			if err != nil {
				return nil, err
			}
		} else if res != nil {
			err = conn.ResponsePacket(res)
			if err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}
