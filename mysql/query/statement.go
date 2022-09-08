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

package query

import (
	vitess "vitess.io/vitess/go/vt/sqlparser"
)

// Statement represents a query statement.
type Statement = vitess.Statement

// StatementType encodes the type of a SQL statement.
type StatementType = vitess.StatementType

// These constants are used to identify the SQL statement type.
// Changing this list will require reviewing all calls to Preview.
const (
	StmtSelect   = vitess.StmtSelect
	StmtStream   = vitess.StmtStream
	StmtInsert   = vitess.StmtInsert
	StmtReplace  = vitess.StmtReplace
	StmtUpdate   = vitess.StmtUpdate
	StmtDelete   = vitess.StmtDelete
	StmtDDL      = vitess.StmtDDL
	StmtBegin    = vitess.StmtBegin
	StmtCommit   = vitess.StmtCommit
	StmtRollback = vitess.StmtRollback
	StmtSet      = vitess.StmtSet
	StmtShow     = vitess.StmtShow
	StmtUse      = vitess.StmtUse
	StmtOther    = vitess.StmtOther
	StmtUnknown  = vitess.StmtUnknown
	StmtComment  = vitess.StmtComment
	StmtPriv     = vitess.StmtPriv
)

const (
	CreateDDLAction           = vitess.CreateDDLAction
	AlterDDLAction            = vitess.AlterDDLAction
	DropDDLAction             = vitess.DropDDLAction
	RenameDDLAction           = vitess.RenameDDLAction
	TruncateDDLAction         = vitess.TruncateDDLAction
	CreateVindexDDLAction     = vitess.CreateVindexDDLAction
	DropVindexDDLAction       = vitess.DropVindexDDLAction
	AddVschemaTableDDLAction  = vitess.AddVschemaTableDDLAction
	DropVschemaTableDDLAction = vitess.DropVschemaTableDDLAction
	AddColVindexDDLAction     = vitess.AddColVindexDDLAction
	DropColVindexDDLAction    = vitess.DropColVindexDDLAction
	AddSequenceDDLAction      = vitess.AddSequenceDDLAction
	AddAutoIncDDLAction       = vitess.AddAutoIncDDLAction
)

// SelectStatement any SELECT statement.
type SelectStatement = vitess.SelectStatement

// Select represents a SELECT statement.
type Select = vitess.Select

// Union represents a UNION statement.
type Union = vitess.Union

// Stream represents a SELECT statement.
type Stream = vitess.Stream

// Insert represents an INSERT or REPLACE statement.
type Insert = vitess.Insert

// Update represents an UPDATE statement.
type Update = vitess.Update

// Delete represents a DELETE statement.
type Delete = vitess.Delete

// Set represents a SET statement.
type Set = vitess.Set

// SetTransaction represents a SET TRANSACTION statement.
type SetTransaction = vitess.SetTransaction

// Characteristic is a transaction related change.
type Characteristic = vitess.Characteristic

// IsolationLevel is self-explanatory in this context.
type IsolationLevel = vitess.IsolationLevel

// AccessMode is ReadOnly/ReadWrite.
type AccessMode = vitess.AccessMode

// DBDDL represents a CREATE, DROP, or ALTER database statement.
type DBDDL = vitess.DBDDL

// DDL represents a CREATE, ALTER, DROP, RENAME, TRUNCATE or ANALYZE statement.
type DDL = vitess.DDL

// Show represents a show statement.
type Show = vitess.Show

// Use represents a use statement.
type Use = vitess.Use

// Begin represents a Begin statement.
type Begin = vitess.Begin

// Commit represents a Commit statement.
type Commit = vitess.Commit

// Rollback represents a Rollback statement.
type Rollback = vitess.Rollback

// Explain represents an EXPLAIN statement
// type Explain = vitess.Explain
