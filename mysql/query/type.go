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
	vitessst "vitess.io/vitess/go/sqltypes"
)

// Vitess data types. These are idiomatically
// named synonyms for the vitess.Type values.
const (
	Unknown = -1
	// NULL_TYPE specifies a NULL type.
	Null = vitessst.Null
	// INT8 specifies a TINYINT type.
	Int8 = vitessst.Int8
	// UINT8 specifies a TINYINT UNSIGNED type.
	Uint8 = vitessst.Uint8
	// INT16 specifies a SMALLINT type.
	Int16 = vitessst.Int16
	// UINT16 specifies a SMALLINT UNSIGNED type.
	Uint16 = vitessst.Uint16
	// INT24 specifies a MEDIUMINT type.
	Int24 = vitessst.Int24
	// UINT24 specifies a MEDIUMINT UNSIGNED type.
	Uint24 = vitessst.Uint24
	// INT32 specifies a INTEGER type.
	Int32 = vitessst.Int32
	// UINT32 specifies a INTEGER UNSIGNED type.
	Uint32 = vitessst.Uint32
	// INT64 specifies a BIGINT type.
	Int64 = vitessst.Int64
	// UINT64 specifies a BIGINT UNSIGNED type.
	Uint64 = vitessst.Uint64
	// FLOAT32 specifies a FLOAT type.
	Float32 = vitessst.Float32
	// FLOAT64 specifies a DOUBLE or REAL type.
	Float64 = vitessst.Float64
	// TIMESTAMP specifies a TIMESTAMP type.
	Timestamp = vitessst.Timestamp
	// DATE specifies a DATE type.
	Date = vitessst.Date
	// TIME specifies a TIME type.
	Time = vitessst.Time
	// DATETIME specifies a DATETIME type.
	Datetime = vitessst.Datetime
	// YEAR specifies a YEAR type.
	Year = vitessst.Year
	// DECIMAL specifies a DECIMAL or NUMERIC type.
	Decimal = vitessst.Decimal
	// TEXT specifies a TEXT type.
	Text = vitessst.Text
	// BLOB specifies a BLOB type.
	Blob = vitessst.Blob
	// VARCHAR specifies a VARCHAR type.
	VarChar = vitessst.VarChar
	// VARBINARY specifies a VARBINARY type.
	VarBinary = vitessst.VarBinary
	// CHAR specifies a CHAR type.
	Char = vitessst.Char
	// BINARY specifies a BINARY type.
	Binary = vitessst.Binary
	// BIT specifies a BIT type.
	Bit = vitessst.Bit
	// ENUM specifies an ENUM type.
	Enum = vitessst.Enum
	// SET specifies a SET type.
	SetType = vitessst.Set
	// GEOMETRY specifies a GEOMETRY type.
	Geometry = vitessst.Geometry
	// JSON specifies a JSON type.
	TypeJSON = vitessst.TypeJSON
	// EXPRESSION specifies a SQL expression.
	Expression = vitessst.Expression
)
