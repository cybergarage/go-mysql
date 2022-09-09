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

package query

// NewNullColumn returns a null column.
func NewNullColumn() *Column {
	col, _ := NewColumnWithTypeAndValue(Null, NewValue())
	return col
}

// NewTextColumnWithValue returns a text column with the specified value.
func NewTextColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Text, val)
}

// NewTinyIntColumnWithValue returns a tiny int column with the specified value.
func NewTinyIntColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Int8, val)
}

// NewShortIntColumnWithValue returns a short int column with the specified value.
func NewShortIntColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Int16, val)
}

// NewIntegerColumnWithValue returns an integer column with the specified value.
func NewIntegerColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Int32, val)
}

// NewBigIntColumnWithValue returns a big int column with the specified value.
func NewBigIntColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Int64, val)
}

// NewFloatColumnWithValue returns a float column with the specified value.
func NewFloatColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Float32, val)
}

// NewDoubleColumnWithValue returns a double column with the specified value.
func NewDoubleColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Float64, val)
}

// NewTimestampColumnWithValue returns a timestamp column with the specified value.
func NewTimestampColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Timestamp, val)
}

// NewTimeColumnWithValue returns a time column with the specified value.
func NewTimeColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Time, val)
}

// NewDateColumnWithValue returns a date column with the specified value.
func NewDateColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Date, val)
}

// NewDatetimeColumnWithValue returns a date column with the specified value.
func NewDatetimeColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Datetime, val)
}

// NewBinaryColumnWithValue returns a binary column with the specified value.
func NewBinaryColumnWithValue(val interface{}) (*Column, error) {
	return NewColumnWithTypeAndValue(Binary, val)
}
