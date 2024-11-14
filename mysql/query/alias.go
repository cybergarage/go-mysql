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
	Columns        = query.Columns
	Index          = query.Index
	Indexes        = query.Indexes
	Table          = query.Table
	TableList      = query.TableList
	Condition      = query.Condition
	Selector       = query.Selector
	Selectors      = query.Selectors
	Expr           = query.Expr
	CmpExpr        = query.CmpExpr
	AndExpr        = query.AndExpr
	OrExpr         = query.OrExpr
	SelectOption   = query.SelectOption
)

const (
	BigIntData       = query.BigIntData
	BinaryData       = query.BinaryData
	BitData          = query.BitData
	BlobData         = query.BlobData
	BooleanData      = query.BooleanData
	CharData         = query.CharData
	CharacterData    = query.CharacterData
	ClobData         = query.ClobData
	DateData         = query.DateData
	DateTimeData     = query.DateTimeData
	DecimalData      = query.DecimalData
	DoubleData       = query.DoubleData
	DoublePrecision  = query.DoublePrecision
	FloatData        = query.FloatData
	IntData          = query.IntData
	IntegerData      = query.IntegerData
	LongBlobData     = query.LongBlobData
	LongTextData     = query.LongTextData
	MediumBlobData   = query.MediumBlobData
	MediumIntData    = query.MediumIntData
	MediumTextData   = query.MediumTextData
	NumericData      = query.NumericData
	RealData         = query.RealData
	SetData          = query.SetData
	SmallIntData     = query.SmallIntData
	TextData         = query.TextData
	TimeData         = query.TimeData
	TimeStampData    = query.TimeStampData
	TinyBlobData     = query.TinyBlobData
	TinyIntData      = query.TinyIntData
	TinyTextData     = query.TinyTextData
	VarBinaryData    = query.VarBinaryData
	VarCharData      = query.VarCharData
	VarCharacterData = query.VarCharacterData
	YearData         = query.YearData
)

// Function represents a function.
type Function = query.Function

// FunctionExecutor represents a function executor.
type FunctionExecutor = query.FunctionExecutor

// AggregateFunction represents an aggregate function.
type AggregateFunction = query.AggregateFunction

// AggregateResultSet represents an aggregate result set.
type AggregateResultSet = query.AggregateResultSet
