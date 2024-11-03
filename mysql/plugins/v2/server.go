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
	stderr "errors"

	"github.com/cybergarage/go-mysql/mysql/errors"
	"github.com/cybergarage/go-mysql/mysql/protocol"
	"github.com/cybergarage/go-mysql/mysql/query"
)

// Server represents a base executor server.
type Server struct {
	*protocol.Server
	executor      Executor
	queryExecutor query.Executor
}

// NewServer returns a base executor server instance.
func NewServer() *Server {
	server := &Server{
		Server:        protocol.NewServer(),
		executor:      nil,
		queryExecutor: nil,
	}
	server.executor = server
	server.Server.SetCommandHandler(server)
	return server
}

// SetExecutor sets an executor to the server.
func (server *Server) SetQueryExecutor(executor query.Executor) {
	server.queryExecutor = executor
}

// QueryExecutor returns the executor of the server.
func (server *Server) QueryExecutor() query.Executor {
	return server.queryExecutor
}

// HandleQuery handles a query.
func (server *Server) HandleQuery(conn protocol.Conn, q *protocol.Query) (protocol.Response, error) {
	connCaps := conn.Capabilities()

	if server.queryExecutor == nil {
		return nil, errors.ErrNotImplemented
	}

	parser := query.NewParser()
	stmts, err := parser.ParseString(q.Query())
	if err != nil {
		return nil, server.queryExecutor.ParserError(conn, q.Query(), err)
	}

	seqID := q.SequenceID().Next()
	for _, stmt := range stmts {
		res, err := server.HandleStatement(conn, stmt)
		if err != nil {
			err = conn.ResponseError(err,
				protocol.WithERRCapability(connCaps),
				protocol.WithERRSecuenceID(seqID),
			)
			if err != nil {
				return nil, err
			}
		} else if res != nil {
			res.SetSequenceID(seqID)
			err = conn.ResponsePacket(res)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

func (server *Server) HandleStatement(conn protocol.Conn, stmt query.Statement) (protocol.Response, error) {
	var err error
	var res protocol.Response

	// nolint: forcetypeassert
	switch stmt.StatementType() {
	case query.BeginStatement:
		stmt := stmt.(query.Begin)
		res, err = server.executor.Begin(conn, stmt)
	case query.CommitStatement:
		stmt := stmt.(query.Commit)
		res, err = server.executor.Commit(conn, stmt)
	case query.RollbackStatement:
		stmt := stmt.(query.Rollback)
		res, err = server.executor.Rollback(conn, stmt)
	case query.CreateDatabaseStatement:
		stmt := stmt.(query.CreateDatabase)
		res, err = server.executor.CreateDatabase(conn, stmt)
	case query.CreateTableStatement:
		stmt := stmt.(query.CreateTable)
		res, err = server.executor.CreateTable(conn, stmt)
	case query.AlterDatabaseStatement:
		stmt := stmt.(query.AlterDatabase)
		res, err = server.executor.AlterDatabase(conn, stmt)
	case query.AlterTableStatement:
		stmt := stmt.(query.AlterTable)
		res, err = server.executor.AlterTable(conn, stmt)
	case query.DropDatabaseStatement:
		stmt := stmt.(query.DropDatabase)
		res, err = server.executor.DropDatabase(conn, stmt)
	case query.DropTableStatement:
		stmt := stmt.(query.DropTable)
		res, err = server.executor.DropTable(conn, stmt)
	case query.InsertStatement:
		stmt := stmt.(query.Insert)
		res, err = server.executor.Insert(conn, stmt)
	case query.SelectStatement:
		stmt := stmt.(query.Select)
		res, err = server.executor.Select(conn, stmt)
	case query.UpdateStatement:
		stmt := stmt.(query.Update)
		res, err = server.executor.Update(conn, stmt)
	case query.DeleteStatement:
		stmt := stmt.(query.Delete)
		res, err = server.executor.Delete(conn, stmt)
	}

	return res, err
}

// Start starts the server.
func (server *Server) Start() error {
	type starter interface {
		Start() error
	}
	starters := []starter{
		server.Server,
	}
	for _, s := range starters {
		err := s.Start()
		if err != nil {
			return stderr.Join(err, server.Stop())
		}
	}
	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	type stopper interface {
		Stop() error
	}
	stoppers := []stopper{
		server.Server,
	}
	for _, s := range stoppers {
		err := s.Stop()
		if err != nil {
			return err
		}
	}
	return nil
}

// Restart restarts the server.
func (server *Server) Restart() error {
	err := server.Stop()
	if err != nil {
		return err
	}
	return server.Start()
}
