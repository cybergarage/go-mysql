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

	vitess "vitess.io/vitess/go/vt/sqlparser"
)

const (
	StrVal   = vitess.StrVal
	IntVal   = vitess.IntVal
	FloatVal = vitess.FloatVal
	// TODO: Support other SQL types.
	// HexNum.
	// HexVal.
	// ValArg.
	// BitVal.
)

// NewValueWithLiteral converts the specified literal to Golang value.
func NewValueWithLiteral(literal *vitess.Literal) (interface{}, error) {
	switch literal.Type {
	case StrVal:
		return literal.Val, nil
	case IntVal:
		return strconv.ParseInt(literal.Val, 0, 64)
	case FloatVal:
		return strconv.ParseFloat(literal.Val, 64)
	default:
		return nil, fmt.Errorf(errorUnknownSQLValType, literal)
	}
}
