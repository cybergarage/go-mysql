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

type resultsetColumn struct {
	name string
	t    DataType
	c    ColumnConstraint
}

// ResultSetColumnOptions represents a functional option for resultsetColumn.
type ResultSetColumnOption func(*resultsetColumn)

// WithResultSetColumnName returns a functional option for resultsetColumn.
func WithResultSetColumnName(name string) ResultSetColumnOption {
	return func(col *resultsetColumn) {
		col.name = name
	}
}

// WithResultSetColumnType returns a functional option for resultsetColumn.
func WithResultSetColumnType(t DataType) ResultSetColumnOption {
	return func(col *resultsetColumn) {
		col.t = t
	}
}

// WithResultSetColumnConstraint returns a functional option for resultsetColumn.
func WithResultSetColumnConstraint(c ColumnConstraint) ResultSetColumnOption {
	return func(col *resultsetColumn) {
		col.c = c
	}
}

// ResultSetColumn represents a resultset column interface.
func NewResultSetColumn(opts ...ResultSetColumnOption) ResultSetColumn {
	col := &resultsetColumn{
		name: "",
		t:    0,
		c:    0,
	}
	for _, opt := range opts {
		opt(col)
	}
	return col
}

// Name returns the column name.
func (col *resultsetColumn) Name() string {
	return col.name
}

// Type returns the column type.
func (col *resultsetColumn) Type() DataType {
	return col.t
}

// Constraint returns the column constraint.
func (col *resultsetColumn) Constraint() ColumnConstraint {
	return col.c
}
