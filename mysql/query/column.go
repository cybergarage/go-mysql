// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apachce.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package query

import (
	"fmt"
	"reflect"
	"time"

	"github.com/cybergarage/go-safecast/safecast"
	vitessst "vitess.io/vitess/go/sqltypes"
	vitesspq "vitess.io/vitess/go/vt/proto/query"
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Type defines the various supported data types in bind vars
// and query results.
type ColumnType = vitesspq.Type

// Field represents a field.
type Field = vitesspq.Field

const (
	NotPrimaryKey    = -1
	SinglePrimaryKey = 0
)

const (
	timestampFormat   = "2006-01-02 15:04:05"
	timestampFormatP3 = "2006-01-02 15:04:05.000"
	timestampFormatP6 = "2006-01-02 15:04:05.000000"
)

// Column represents a column.
type Column struct {
	name  string
	typ   ColumnType
	value *Value
}

// NewColumn returns a new column.
func NewColumn() *Column {
	return NewColumnWithName("")
}

// NewColumnWithNameAndType returns a column with the specified name and type.
func NewColumnWithNameAndType(name string, t ColumnType) *Column {
	col := &Column{
		name:  name,
		typ:   t,
		value: NewValue(),
	}
	return col
}

// NewColumnWithName returns a column with the specified name.
func NewColumnWithName(name string) *Column {
	return NewColumnWithNameAndType(name, Unknown)
}

// NewColumnWithNameAndValue returns a column with the specified name and value.
func NewColumnWithNameAndValue(name string, val interface{}) (*Column, error) {
	col := NewColumnWithName(name)
	err := col.SetValue(val)
	return col, err
}

// NewColumnWithValue returns a column with the specified value.
func NewColumnWithValue(val interface{}) (*Column, error) {
	col := NewColumn()
	err := col.SetValue(val)
	return col, err
}

// NewColumnWithTypeAndValue returns a column with the specified type and value.
func NewColumnWithTypeAndValue(typ ColumnType, val interface{}) (*Column, error) {
	col := NewColumn()
	col.typ = typ
	err := col.SetValue(val)
	return col, err
}

// NewColumnWithColumnDefinition returns a column with the specified vitess column definition.
func NewColumnWithColumnDefinition(column *vitesssp.ColumnDefinition) *Column {
	col := &Column{
		name: column.Name.String(),
		typ:  column.Type.SQLType(),
	}
	return col
}

// SetName sets a column name.
func (col *Column) SetName(name string) {
	col.name = name
}

// Name returns the column name.
func (col *Column) Name() string {
	return col.name
}

// SetValue sets a value.
func (col *Column) SetValue(val interface{}) error {
	return col.value.SetValue(val)
}

// Value returns the value.
func (col *Column) Value() interface{} {
	return col.value.Value()
}

// Type returns the SQL type.
func (col *Column) Type() ColumnType {
	return col.typ
}

// Equals returns true when the specified column is equals to this column, otherwise false.
func (col *Column) Equals(other *Column) bool {
	if col.name != other.name {
		return false
	}
	if !reflect.DeepEqual(col.Value(), other.Value()) {
		return false
	}
	return true
}

// ToValue converts a column to a vitess value.
func (col *Column) ToValue() (vitessst.Value, error) {
	value := col.Value()
	switch v := value.(type) {
	case int:
		return vitessst.InterfaceToValue(int64(v))
	case int8:
		return vitessst.InterfaceToValue(int64(v))
	case int16:
		return vitessst.InterfaceToValue(int64(v))
	case int32:
		return vitessst.InterfaceToValue(int64(v))
	case uint:
		return vitessst.InterfaceToValue(uint64(v))
	case uint8:
		return vitessst.InterfaceToValue(uint64(v))
	case uint16:
		return vitessst.InterfaceToValue(uint64(v))
	case uint32:
		return vitessst.InterfaceToValue(uint64(v))
	case float32:
		return vitessst.InterfaceToValue(float64(v))
	case time.Time:
		tv := v.Format(timestampFormat)
		return vitessst.InterfaceToValue([]byte(tv))
	case bool:
		if v {
			return vitessst.InterfaceToValue(int64(1))
		}
		return vitessst.InterfaceToValue(int64(0))
	}
	return vitessst.InterfaceToValue(value)
}

// ForValue converts a column to a vitess value for the specified SQL type.
func (col *Column) ForValue(t SQLType) (vitessst.Value, error) {
	switch t { // nolint: exhaustive
	case Timestamp, Datetime:
		var v time.Time
		err := safecast.ToTime(col.Value(), &v)
		if err == nil {
			ts := v.Format(timestampFormatP6)
			return vitessst.InterfaceToValue([]byte(ts))
		}
	}
	return col.ToValue()
}

// String returns the string representation.
func (col *Column) String() string {
	if col.value == nil {
		return col.name
	}
	return fmt.Sprintf("%s:%s", col.name, col.Value())
}
