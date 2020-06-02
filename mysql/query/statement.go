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
	vitesssql "vitess.io/vitess/go/vt/sqlparser"
)

// Statement represents a query statement.
type Statement = vitesssql.Statement

// StatementType encodes the type of a SQL statement
type StatementType = vitesssql.StatementType

// These constants are used to identify the SQL statement type.
// Changing this list will require reviewing all calls to Preview.
const (
	StmtSelect   = vitesssql.StmtSelect
	StmtStream   = vitesssql.StmtStream
	StmtInsert   = vitesssql.StmtInsert
	StmtReplace  = vitesssql.StmtReplace
	StmtUpdate   = vitesssql.StmtUpdate
	StmtDelete   = vitesssql.StmtDelete
	StmtDDL      = vitesssql.StmtDDL
	StmtBegin    = vitesssql.StmtBegin
	StmtCommit   = vitesssql.StmtCommit
	StmtRollback = vitesssql.StmtRollback
	StmtSet      = vitesssql.StmtSet
	StmtShow     = vitesssql.StmtShow
	StmtUse      = vitesssql.StmtUse
	StmtOther    = vitesssql.StmtOther
	StmtUnknown  = vitesssql.StmtUnknown
	StmtComment  = vitesssql.StmtComment
	StmtPriv     = vitesssql.StmtPriv
	//StmtExplain  = vitesssql.StmtExplain
)

// SelectStatement any SELECT statement.
type SelectStatement = vitesssql.SelectStatement

// Select represents a SELECT statement.
type Select = vitesssql.Select

// Union represents a UNION statement.
type Union = vitesssql.Union

// Stream represents a SELECT statement.
type Stream = vitesssql.Stream

// Insert represents an INSERT or REPLACE statement.
type Insert = vitesssql.Insert

// Update represents an UPDATE statement.
type Update = vitesssql.Update

// Delete represents a DELETE statement.
type Delete = vitesssql.Delete

// Set represents a SET statement.
type Set = vitesssql.Set

// SetTransaction represents a SET TRANSACTION statement.
type SetTransaction = vitesssql.SetTransaction

// Characteristic is a transaction related change
type Characteristic = vitesssql.Characteristic

// IsolationLevel is self-explanatory in this context
type IsolationLevel = vitesssql.IsolationLevel

// AccessMode is ReadOnly/ReadWrite
type AccessMode = vitesssql.AccessMode

// DBDDL represents a CREATE, DROP, or ALTER database statement.
type DBDDL = vitesssql.DBDDL

// DDL represents a CREATE, ALTER, DROP, RENAME, TRUNCATE or ANALYZE statement.
type DDL = vitesssql.DDL

// ParenSelect is a parenthesized SELECT statement.
type ParenSelect = vitesssql.ParenSelect

// Show represents a show statement.
type Show = vitesssql.Show

// Use represents a use statement.
type Use = vitesssql.Use

// Begin represents a Begin statement.
type Begin = vitesssql.Begin

// Commit represents a Commit statement.
type Commit = vitesssql.Commit

// Rollback represents a Rollback statement.
type Rollback = vitesssql.Rollback

// Explain represents an EXPLAIN statement
//type Explain = vitesssql.Explain
