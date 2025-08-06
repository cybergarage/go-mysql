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

package protocol

import (
	"bytes"
	"io"

	"github.com/cybergarage/go-mysql/mysql/query"
)

// MySQL: COM_STMT_PREPARE
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_prepare.html
// COM_STMT_PREPARE - MariaDB Knowledge Base
// https://mariadb.com/kb/en/com_stmt_prepare/

// StmtPrepare represents a COM_STMT_PREPARE packet.
type StmtPrepare struct {
	Command

	dbName               string
	tblNames             []string
	stmt                 query.Statement
	query                string
	resultSetColumnNames []string
	parameterColumnNames []string
}

func newStmtPrepareWithCommand(cmd Command, opts ...StmtPrepareOption) *StmtPrepare {
	q := &StmtPrepare{
		Command:              cmd,
		query:                "",
		stmt:                 nil,
		dbName:               "",
		tblNames:             nil,
		resultSetColumnNames: nil,
		parameterColumnNames: nil,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// StmtPrepareOption represents a MySQL StmtPrepare option.
type StmtPrepareOption func(*StmtPrepare)

// WithStmtPrepareQuery sets the query string.
func WithStmtPrepareQuery(query string) StmtPrepareOption {
	return func(q *StmtPrepare) {
		q.query = query
	}
}

// WithStmtPrepareCapability sets the capability.
func WithStmtPrepareCapability(c Capability) StmtPrepareOption {
	return func(q *StmtPrepare) {
		q.Command.SetCapability(c)
	}
}

// WithStmtPrepareServerStatus sets the server status.
func WithStmtPrepareServerStatus(s ServerStatus) StmtPrepareOption {
	return func(q *StmtPrepare) {
		q.Command.SetServerStatus(s)
	}
}

// WithStmtPrepareDatabase sets the database name.
func WithStmtPrepareDatabase(dbName string) StmtPrepareOption {
	return func(q *StmtPrepare) {
		q.dbName = dbName
	}
}

// NewStmtPrepareFromReader reads a COM_STMT_PREPARE packet.
func NewStmtPrepareFromReader(reader io.Reader, opts ...StmtPrepareOption) (*StmtPrepare, error) {
	var err error

	cmd, err := NewCommandFromReader(reader)
	if err != nil {
		return nil, err
	}

	if err = cmd.IsType(ComStmtPrepare); err != nil {
		return nil, err
	}

	return NewStmtPrepareFromCommand(cmd, opts...)
}

// NewStmtPrepareFromCommand creates a new StmtPrepare from a Command.
func NewStmtPrepareFromCommand(cmd Command, opts ...StmtPrepareOption) (*StmtPrepare, error) {
	var err error

	pkt := newStmtPrepareWithCommand(cmd, opts...)

	payload := cmd.Payload()
	reader := NewPacketReaderWithReader(bytes.NewBuffer(payload[1:]))

	pkt.query, err = reader.ReadEOFTerminatedString()
	if err != nil {
		return nil, err
	}

	if err = pkt.parseQuery(); err != nil {
		return nil, err
	}

	return pkt, nil
}

// nolint: forcetypeassert
func (pkt *StmtPrepare) parseQuery() error {
	parser := query.NewParser()
	stmts, err := parser.ParseString(pkt.query)
	if err != nil {
		return err
	}
	if len(stmts) != 1 {
		return newInvalidStatement(pkt.query)
	}

	pkt.stmt = stmts[0]
	pkt.resultSetColumnNames = []string{}
	pkt.parameterColumnNames = []string{}

	var parameterColumnsFromExpr func(expr query.Expr) []string
	parameterColumnsFromExpr = func(expr query.Expr) []string {
		columns := []string{}
		switch expr := expr.(type) {
		case *query.CmpExpr:
			if expr.Right().IsPlaceHolder() {
				columns = append(columns, expr.Left().Name())
			}
		case *query.AndExpr:
			columns = append(columns, parameterColumnsFromExpr(expr.Left())...)
			columns = append(columns, parameterColumnsFromExpr(expr.Right())...)
		case *query.OrExpr:
			columns = append(columns, parameterColumnsFromExpr(expr.Left())...)
			columns = append(columns, parameterColumnsFromExpr(expr.Right())...)
		}
		return columns
	}

	parameterColumnsFromCondition := func(cond query.Condition) []string {
		if !cond.HasConditions() {
			return []string{}
		}
		return parameterColumnsFromExpr(cond.Expr())
	}

	switch pkt.stmt.StatementType() {
	case query.InsertStatement:
		stmt := pkt.stmt.(query.Insert)
		pkt.tblNames = []string{stmt.TableName()}
		for _, col := range stmt.Columns() {
			if !col.IsPlaceHolder() {
				continue
			}
			pkt.parameterColumnNames = append(pkt.parameterColumnNames, col.Name())
		}
	case query.SelectStatement:
		stmt := pkt.stmt.(query.Select)
		pkt.tblNames = stmt.From().TableNames()
		for _, sel := range stmt.Selectors() {
			pkt.resultSetColumnNames = append(pkt.resultSetColumnNames, sel.Name())
		}
		pkt.parameterColumnNames = parameterColumnsFromCondition(stmt.Where())
	case query.UpdateStatement:
		stmt := pkt.stmt.(query.Update)
		pkt.tblNames = []string{stmt.TableName()}
		pkt.parameterColumnNames = parameterColumnsFromCondition(stmt.Where())
	case query.DeleteStatement:
		stmt := pkt.stmt.(query.Delete)
		pkt.tblNames = []string{stmt.TableName()}
		pkt.parameterColumnNames = parameterColumnsFromCondition(stmt.Where())
	default:
		return newInvalidStatement(pkt.query)
	}

	return nil
}

// Query returns the query string.
func (pkt *StmtPrepare) Query() string {
	return pkt.query
}

// DatabaseName returns the database name.
func (pkt *StmtPrepare) DatabaseName() string {
	return pkt.dbName
}

// Statement returns the statement.
func (pkt *StmtPrepare) Statement() query.Statement {
	return pkt.stmt
}

// TableNames returns the table names.
func (pkt *StmtPrepare) TableNames() []string {
	return pkt.tblNames
}

// SetResultSetColumnNames sets the result set column names.
func (pkt *StmtPrepare) ResultSetColumnNames() []string {
	return pkt.resultSetColumnNames
}

// SetParameterColumnNames sets the parameter column names.
func (pkt *StmtPrepare) ParameterColumnNames() []string {
	return pkt.parameterColumnNames
}

// Bytes returns the packet bytes.
func (pkt *StmtPrepare) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteCommandType(pkt); err != nil {
		return nil, err
	}

	if _, err := w.WriteBytes([]byte(pkt.query)); err != nil {
		return nil, err
	}

	return pkt.Command.Bytes()
}
