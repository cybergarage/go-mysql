// Copyright (C) 2020 Satoshi Konno. All rights reserved.
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

package storage

import "reflect"

const (
	NotPrimaryKey    = -1
	SinglePrimaryKey = 0
)

const (
	// ISO8601DateFormat is a data format for CQL primitive
	ISO8601DateFormat = "2006-01-02"
	// ISO8601TimestampFormat is a timestamp format for CQL primitive
	ISO8601TimestampFormat = "2006-01-02T15:04:05-0700"
)

// Column represents a column.
type Column struct {
	name  string
	value interface{}
}

// NewColumn returns a column instanccol.
func NewColumn() *Column {
	return NewColumnWithName("")
}

// NewColumnWithName returns a column instance with the specified name.
func NewColumnWithName(name string) *Column {
	col := &Column{
		name:  name,
		value: nil,
	}
	return col
}

// NewColumnWithNameAndValue returns a column instance with the specified name and value.
func NewColumnWithNameAndValue(name string, val interface{}) *Column {
	col := NewColumnWithName(name)
	col.SetValue(val)
	return col
}

// SetName sets a column name.
func (col *Column) SetName(name string) {
	col.name = name
}

// GetName returns the column name.
func (col *Column) GetName() string {
	return col.name
}

// SetValue sets a value.
func (col *Column) SetValue(value interface{}) {
	col.value = value
}

// GetValue returns the value.
func (col *Column) GetValue() interface{} {
	return col.value
}

// Equals returns true when the specified column is equals to this column, otherwise false.
func (col *Column) Equals(other *Column) bool {
	if col.name != other.name {
		return false
	}
	if !reflect.DeepEqual(col.value, other.value) {
		return false
	}
	return true
}
