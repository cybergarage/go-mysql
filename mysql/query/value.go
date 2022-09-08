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

	"vitess.io/vitess/go/sqltypes"
	vitesstypes "vitess.io/vitess/go/sqltypes"
	vitesspb "vitess.io/vitess/go/vt/proto/query"
	vitesssql "vitess.io/vitess/go/vt/sqlparser"
)

const (
	StrVal   = vitesssql.StrVal
	IntVal   = vitesssql.IntVal
	FloatVal = vitesssql.FloatVal
	HexNum   = vitesssql.HexNum
	HexVal   = vitesssql.HexVal
	BitVal   = vitesssql.BitVal
)

// SQLVal represents a single value.
type SQLVal = vitesstypes.Value

// ValType specifies the type for SQLVal.
type ValType = vitesspb.Type

// Value represents a query value.
type Value struct {
	typ   ValType
	value interface{}
}

// Expr represents an expression.
type Expr = vitesssql.Expr

// ValTuple represents a tuple of actual values.
type ValTuple = vitesssql.ValTuple

// ColName represents a column name.
type ColName = vitesssql.ColName

// AndExpr represents an AND expression.
type AndExpr = vitesssql.AndExpr

// OrExpr represents an OR expression.
type OrExpr = vitesssql.OrExpr

// XorExpr represents an XOR expression.
type XorExpr = vitesssql.XorExpr

// NotExpr represents a NOT expression.
type NotExpr = vitesssql.NotExpr

// ComparisonExpr represents a two-value comparison expression.
type ComparisonExpr = vitesssql.ComparisonExpr

// NewValue creates a query value .
func NewValue() *Value {
	return &Value{}
}

// NewValueWithType creates a query value with the specified type.
func NewValueWithType(val interface{}, valType ValType) *Value {
	return &Value{
		typ:   valType,
		value: val,
	}
}

// NewValueWithValue creates a query value from the raw SQL value.
func NewValueWithValue(val interface{}) (*Value, error) {
	value := NewValue()
	return value, value.SetValue(val)
}

// SetType sets a value type.
func (value *Value) SetType(t ValType) {
	value.typ = t
}

// Type returns the type.
func (value *Value) Type() ValType {
	return value.typ
}

// SetValue sets a value.
func (value *Value) SetValue(val interface{}) error {
	switch v := val.(type) {
	case SQLVal:
		return value.setSQLValue(v)
	case *Literal:
		return value.setLiteralValue(v)
	}

	value.value = val

	return nil
}

// setLiteralValue sets a literal value and type by the specified value, and
// returns SetLiteralValue error when the specified value is unknown.
// See: vitess.io/vitess/go/sqltypes::NewValue().
func (value *Value) setLiteralValue(lv *Literal) error {
	var v interface{}
	var vt vitesspb.Type
	var err error

	switch lv.Type {
	case StrVal:
		vt = sqltypes.VarBinary
		v = string(lv.Val)
	case IntVal:
		vt = sqltypes.Int64
		v, err = strconv.ParseInt(string(lv.Val), 0, 64)
	case FloatVal:
		vt = sqltypes.Float64
		v, err = strconv.ParseFloat(string(lv.Val), 64)
	default:
		err = fmt.Errorf(errorLiteralUnknownType, lv)
	}

	if err != nil {
		return err
	}

	value.typ = vt
	value.value = v

	return nil
}

// setSQLValue sets a query value and type by the specified value, and
// returns an error when the specified value is unknown.
// See: vitess.io/vitess/go/sqltypes::NewValue().
func (value *Value) setSQLValue(sv SQLVal) error {
	var v interface{}
	var err error

	vt := sv.Type()
	switch {
	case vitesstypes.IsSigned(vt):
		v, err = strconv.ParseInt(string(sv.Raw()), 0, 64)
	case vitesstypes.IsUnsigned(vt):
		v, err = strconv.ParseUint(string(sv.Raw()), 0, 64)
	case vitesstypes.IsFloat(vt) || vt == Decimal:
		v, err = strconv.ParseFloat(string(sv.Raw()), 64)
	case vitesstypes.IsQuoted(vt) || vt == Bit || vt == Null:
		v = string(sv.Raw())
	default:
		return fmt.Errorf(errorColumnUnexpectedType, sv.Raw(), sv.Type())
	}

	if err != nil {
		return err
	}

	value.typ = vt
	value.value = v

	return nil
}

// Value returns the value.
func (value *Value) Value() interface{} {
	return value.value
}
