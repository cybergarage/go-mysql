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
	"fmt"
	"strconv"

	"vitess.io/vitess/go/vt/sqlparser"
	vitess "vitess.io/vitess/go/vt/sqlparser"
)

const (
	StrVal   = vitess.StrVal
	IntVal   = vitess.IntVal
	FloatVal = vitess.FloatVal
	// TODO: Support other SQL types
	// HexNum
	// HexVal
	// ValArg
	// BitVal
)

// NewValueWithSQLVal converts the specified SQL value to Golang value.
func NewValueWithSQLVal(sqlVal *sqlparser.SQLVal) (interface{}, error) {
	switch sqlVal.Type {
	case StrVal:
		return string(sqlVal.Val), nil
	case IntVal:
		return strconv.ParseInt(string(sqlVal.Val), 0, 64)
	case FloatVal:
		return strconv.ParseFloat(string(sqlVal.Val), 64)
	}
	return nil, fmt.Errorf(errorUnknownSQLValType, sqlVal)
}
