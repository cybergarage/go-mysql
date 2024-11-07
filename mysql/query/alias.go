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
	"github.com/cybergarage/go-sqlparser/sql"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

const (
	EQ  = query.EQ
	NEQ = query.NEQ
	LT  = query.LT
	LE  = query.LE
	GT  = query.GT
	GE  = query.GE
	IN  = query.IN
	NIN = query.NIN
)

type (
	BindParam      = query.BindParam
	CreateDatabase = sql.CreateDatabase
	CreateTable    = sql.CreateTable
	AlterDatabase  = sql.AlterDatabase
	AlterTable     = sql.AlterTable
	DropDatabase   = sql.DropDatabase
	DropTable      = sql.DropTable
	Select         = sql.Select
	Insert         = sql.Insert
	Update         = sql.Update
	Delete         = sql.Delete
	Copy           = sql.Copy
	Begin          = sql.Begin
	Commit         = sql.Commit
	Rollback       = sql.Rollback
	Vacuum         = sql.Vacuum
	Truncate       = sql.Truncate
	Use            = sql.Use
	Schema         = sql.Schema
	Column         = query.Column
	ColumnList     = query.ColumnList
	Table          = query.Table
	TableList      = query.TableList
	Condition      = query.Condition
	Selector       = query.Selector
	SelectorList   = query.SelectorList
	Expr           = query.Expr
	CmpExpr        = query.CmpExpr
	AndExpr        = query.AndExpr
	OrExpr         = query.OrExpr
	SelectOption   = query.SelectOption
)

// Function represents a function.
type Function = query.Function

// FunctionExecutor represents a function executor.
type FunctionExecutor = query.FunctionExecutor

// AggregateFunction represents an aggregate function.
type AggregateFunction = query.AggregateFunction

// AggregateResultSet represents an aggregate result set.
type AggregateResultSet = query.AggregateResultSet