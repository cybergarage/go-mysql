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

package query

import (
	"github.com/cybergarage/go-sqlparser/sql/query"
)

type Statement = query.Statement
type StatementType = query.StatementType

const (
	CreateDatabaseStatement = query.CreateDatabaseStatement
	CreateTableStatement    = query.CreateTableStatement
	CreateIndexStatement    = query.CreateIndexStatement
	InsertStatement         = query.InsertStatement
	SelectStatement         = query.SelectStatement
	UpdateStatement         = query.UpdateStatement
	DeleteStatement         = query.DeleteStatement
	DropDatabaseStatement   = query.DropDatabaseStatement
	DropTableStatement      = query.DropTableStatement
	DropIndexStatement      = query.DropIndexStatement
	AlterDatabaseStatement  = query.AlterDatabaseStatement
	AlterTableStatement     = query.AlterTableStatement
	AlterIndexStatement     = query.AlterIndexStatement
	CopyStatement           = query.CopyStatement
	CommitStatement         = query.CommitStatement
	VacuumStatement         = query.VacuumStatement
	TruncateStatement       = query.TruncateStatement
	BeginStatement          = query.BeginStatement
	RollbackStatement       = query.RollbackStatement
)
