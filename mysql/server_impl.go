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

package mysql

import (
	stderr "errors"

	"github.com/cybergarage/go-mysql/mysql/errors"
	"github.com/cybergarage/go-mysql/mysql/protocol"
	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-sqlparser/sql/system"
)

// server represents a base executor server.
type server struct {
	*protocol.Server
	sqlExecutor     SQLExecutor
	queryExecutor   QueryExecutor
	exQueryExecutor ExQueryExecutor
	errorHandler    ErrorHandler
}

// NewServer returns a base executor server instance.
func NewServer() Server {
	server := &server{
		Server:          protocol.NewServer(),
		sqlExecutor:     nil,
		queryExecutor:   NewDefaultQueryExecutor(),
		exQueryExecutor: nil,
		errorHandler:    nil,
	}

	server.exQueryExecutor = NewDefaultExQueryExecutorWith(
		server.queryExecutor,
	)

	server.Server.SetProductName(PackageName)
	server.Server.SetProductVersion(Version)
	server.Server.SetCommandHandler(server)

	return server
}

// SetExecutor sets an executor to the server.
func (server *server) SetSQLExecutor(sqlExeutor SQLExecutor) {
	server.sqlExecutor = sqlExeutor
	executors := []any{
		server.queryExecutor,
		server.exQueryExecutor,
		server.errorHandler,
	}
	for _, executor := range executors {
		if executor == nil {
			continue
		}
		if _, ok := executor.(Server); ok {
			continue
		}
		if setter, ok := executor.(SQLExecutorSetter); ok {
			setter.SetSQLExecutor(sqlExeutor)
		}
	}
}

// SetQueryExecutor sets a user query executor.
func (server *server) SetQueryExecutor(executor QueryExecutor) {
	server.queryExecutor = executor
}

// SetExQueryExecutor sets a user extended query executor.
func (server *server) SetExQueryExecutor(executor ExQueryExecutor) {
	server.exQueryExecutor = executor
}

// SetErrorHandler sets a user error handler.
func (server *server) SetErrorHandler(handler ErrorHandler) {
	server.errorHandler = handler
}

// SQLExecutor returns the executor of the server.
func (server *server) SQLExecutor() SQLExecutor {
	return server.sqlExecutor
}

// QueryExecutor returns the user query executor.
func (server *server) QueryExecutor() QueryExecutor {
	return server.queryExecutor
}

// ErrorHandler returns the user error handler.
func (server *server) ErrorHandler() ErrorHandler {
	return server.errorHandler
}

// HandleQuery handles a query.
func (server *server) HandleQuery(conn protocol.Conn, q *protocol.Query) (protocol.Response, error) {
	connCaps := conn.Capability()

	if server.sqlExecutor == nil {
		return nil, errors.ErrNotImplemented
	}

	parser := query.NewParser()
	stmts, err := parser.ParseString(q.Query())
	if err != nil {
		return server.errorHandler.ParserError(conn, q.Query(), err)
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
			err = conn.ResponsePacket(res,
				protocol.WithResponseSequenceID(seqID),
				protocol.WithResponseCapability(connCaps),
			)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, nil
}

// PrepareStatement prepares a statement.
func (server *server) PrepareStatement(conn protocol.Conn, stmtPrep *protocol.StmtPrepare) (*protocol.StmtPrepareResponse, error) {
	stmt, err := system.NewSchemaColumnsStatement(
		system.WithSchemaColumnsStatementDatabaseName(conn.Database()),
		system.WithSchemaColumnsStatementTableNames(stmtPrep.TableNames()),
	)
	if err != nil {
		return nil, err
	}
	rs, err := server.SQLExecutor().SystemSelect(conn, stmt.Statement())
	if err != nil {
		return nil, err
	}

	opts := []protocol.StmtPrepareResponseOption{}
	schemaColumnRs, err := system.NewSchemaColumnsResultSetFromResultSet(rs)
	if err != nil {
		return nil, err
	}

	// Create response columns for the prepare response.

	isAllColumnsLookup := func(columnNames []string) bool {
		switch {
		case len(columnNames) == 0:
			return true
		case len(columnNames) == 1 && columnNames[0] == "*":
			return true
		}
		return false
	}

	lookupColumDefs := func(schemaColumnRs system.SchemaColumnsResultSet, columnNames []string) ([]protocol.ColumnDef, error) {
		columnDefs := []protocol.ColumnDef{}
		if isAllColumnsLookup(columnNames) {
			for _, schemaColumn := range schemaColumnRs.Columns() {
				columnDef, err := protocol.NewColumnDefsFromSystemSchemaColumn(schemaColumn)
				if err != nil {
					return nil, err
				}
				columnDefs = append(columnDefs, columnDef)
			}
			return columnDefs, nil
		}
		for _, columnName := range columnNames {
			schemaColumn, err := schemaColumnRs.LookupColumn(columnName)
			if err != nil {
				return nil, err
			}
			columnDef, err := protocol.NewColumnDefsFromSystemSchemaColumn(schemaColumn)
			if err != nil {
				return nil, err
			}
			columnDefs = append(columnDefs, columnDef)
		}
		return columnDefs, nil
	}

	resultSetColumnDefs, err := lookupColumDefs(schemaColumnRs, stmtPrep.ResultSetColumnNames())
	if err != nil {
		return nil, err
	}
	opts = append(opts, protocol.WithStmtPrepareResponseColumns(resultSetColumnDefs))

	// Create response params for the prepare response.

	paramColumnDefs, err := lookupColumDefs(schemaColumnRs, stmtPrep.ParameterColumnNames())
	if err != nil {
		return nil, err
	}
	opts = append(opts, protocol.WithStmtPrepareResponseParams(paramColumnDefs))

	// Generate next statement ID and create prepare response.

	stmtID, err := server.NextPreparedStatementID()
	if err != nil {
		return nil, err
	}
	opts = append(opts, protocol.WithStmtPrepareResponseStatementID(stmtID))

	stmPrepRes := protocol.NewStmtPrepareResponse(opts...)

	// Register the prepared statement.

	premStmt := protocol.NewPreparedStatmentWith(stmtPrep, stmPrepRes)
	err = server.RegisterPreparedStatement(premStmt)
	if err != nil {
		return nil, err
	}

	// Return the prepared statement response.

	return stmPrepRes, nil
}

// ExecuteStatement executes a statement.
func (server *server) ExecuteStatement(conn protocol.Conn, stmtExec *protocol.StmtExecute) (protocol.Response, error) {
	preStmt, err := server.PreparedStatement(stmtExec.StatementID())
	if err != nil {
		return nil, err
	}
	stmt, err := preStmt.Bind(stmtExec.Parameters())
	if err != nil {
		return nil, err
	}
	res, err := server.HandleStatement(conn, stmt)
	if err != nil {
		return nil, err
	}
	switch res := res.(type) {
	case *protocol.TextResultSet:
		return protocol.NewBinaryResultSetFrom(res)
	}
	return res, nil
}

// CloseStatement closes a statement.
func (server *server) CloseStatement(conn protocol.Conn, stmt *protocol.StmtClose) (protocol.Response, error) {
	server.RemovePreparedStatement(stmt.StatementID())
	return nil, nil
}

func (server *server) HandleStatement(conn protocol.Conn, stmt query.Statement) (protocol.Response, error) {
	var err error
	var res protocol.Response

	// nolint: forcetypeassert
	switch stmt.StatementType() {
	case query.BeginStatement:
		stmt := stmt.(query.Begin)
		res, err = server.queryExecutor.Begin(conn, stmt)
	case query.CommitStatement:
		stmt := stmt.(query.Commit)
		res, err = server.queryExecutor.Commit(conn, stmt)
	case query.RollbackStatement:
		stmt := stmt.(query.Rollback)
		res, err = server.queryExecutor.Rollback(conn, stmt)
	case query.CreateDatabaseStatement:
		stmt := stmt.(query.CreateDatabase)
		res, err = server.queryExecutor.CreateDatabase(conn, stmt)
	case query.CreateTableStatement:
		stmt := stmt.(query.CreateTable)
		res, err = server.queryExecutor.CreateTable(conn, stmt)
	case query.CreateIndexStatement:
		stmt := stmt.(query.CreateIndex)
		res, err = server.exQueryExecutor.CreateIndex(conn, stmt)
	case query.AlterDatabaseStatement:
		stmt := stmt.(query.AlterDatabase)
		res, err = server.queryExecutor.AlterDatabase(conn, stmt)
	case query.AlterTableStatement:
		stmt := stmt.(query.AlterTable)
		res, err = server.queryExecutor.AlterTable(conn, stmt)
	case query.DropDatabaseStatement:
		stmt := stmt.(query.DropDatabase)
		res, err = server.queryExecutor.DropDatabase(conn, stmt)
	case query.DropTableStatement:
		stmt := stmt.(query.DropTable)
		res, err = server.queryExecutor.DropTable(conn, stmt)
	case query.DropIndexStatement:
		stmt := stmt.(query.DropIndex)
		res, err = server.exQueryExecutor.DropIndex(conn, stmt)
	case query.InsertStatement:
		stmt := stmt.(query.Insert)
		res, err = server.queryExecutor.Insert(conn, stmt)
	case query.SelectStatement:
		stmt := stmt.(query.Select)
		res, err = server.queryExecutor.Select(conn, stmt)
	case query.UpdateStatement:
		stmt := stmt.(query.Update)
		res, err = server.queryExecutor.Update(conn, stmt)
	case query.DeleteStatement:
		stmt := stmt.(query.Delete)
		res, err = server.queryExecutor.Delete(conn, stmt)
	case query.UseStatement:
		stmt := stmt.(query.Use)
		res, err = server.queryExecutor.Use(conn, stmt)
	case query.TruncateStatement:
		stmt := stmt.(query.Truncate)
		res, err = server.exQueryExecutor.Truncate(conn, stmt)
	}

	return res, err
}

// Start starts the server.
func (server *server) Start() error {
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
func (server *server) Stop() error {
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
func (server *server) Restart() error {
	err := server.Stop()
	if err != nil {
		return err
	}
	return server.Start()
}
