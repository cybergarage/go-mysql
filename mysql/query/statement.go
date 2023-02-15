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

package query

import (
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Statement represents a query statement.
type Statement = vitesssp.Statement

// StatementType encodes the type of a SQL statement.
type StatementType = vitesssp.StatementType

// These constants are used to identify the SQL statement type.
// Changing this list will require reviewing all calls to Preview.
const (
	StmtSelect   = vitesssp.StmtSelect
	StmtStream   = vitesssp.StmtStream
	StmtInsert   = vitesssp.StmtInsert
	StmtReplace  = vitesssp.StmtReplace
	StmtUpdate   = vitesssp.StmtUpdate
	StmtDelete   = vitesssp.StmtDelete
	StmtDDL      = vitesssp.StmtDDL
	StmtBegin    = vitesssp.StmtBegin
	StmtCommit   = vitesssp.StmtCommit
	StmtRollback = vitesssp.StmtRollback
	StmtSet      = vitesssp.StmtSet
	StmtShow     = vitesssp.StmtShow
	StmtUse      = vitesssp.StmtUse
	StmtOther    = vitesssp.StmtOther
	StmtUnknown  = vitesssp.StmtUnknown
	StmtComment  = vitesssp.StmtComment
	StmtPriv     = vitesssp.StmtPriv
	// StmtExplain  = vitesssp.StmtExplain.
)

// SelectStatement any SELECT statement.
type SelectStatement = vitesssp.SelectStatement

// Union represents a UNION statement.
type Union = vitesssp.Union

// Stream represents a SELECT statement.
type Stream = vitesssp.Stream

// Set represents a SET statement.
type Set = vitesssp.Set

// SetTransaction represents a SET TRANSACTION statement.
type SetTransaction = vitesssp.SetTransaction

// Characteristic is a transaction related change.
type Characteristic = vitesssp.Characteristic

// IsolationLevel is self-explanatory in this context.
type IsolationLevel = vitesssp.IsolationLevel

// AccessMode is ReadOnly/ReadWrite.
type AccessMode = vitesssp.AccessMode

// DBDDL represents a CREATE, DROP, or ALTER database statement.
type DBDDL = vitesssp.DBDDLStatement

// DDL represents a CREATE, ALTER, DROP, RENAME, TRUNCATE or ANALYZE statement.
type DDL = vitesssp.DDLStatement

// Show represents a show statement.
type Show = vitesssp.Show

// Use represents a use statement.
type Use = vitesssp.Use

// Begin represents a Begin statement.
type Begin = vitesssp.Begin

// Commit represents a Commit statement.
type Commit = vitesssp.Commit

// Rollback represents a Rollback statement.
type Rollback = vitesssp.Rollback

// Explain represents an EXPLAIN statement.
type Explain = vitesssp.Explain

// ColumnDefinition describes a column in a CREATE TABLE statement.
type ColumnDefinition = vitesssp.ColumnDefinition

// UpdateExprs represents a list of update expressions.
type UpdateExprs = vitesssp.UpdateExprs

// Literal represents a fixed value.
type Literal = vitesssp.Literal
