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
	"github.com/cybergarage/go-mysql/mysql/errors"
	"github.com/cybergarage/go-safecast/safecast"
)

type resultsetRow struct {
	values []any
}

// ResultSetRowOptions represents a functional option for resultsetRow.
type ResultSetRowOptions func(*resultsetRow)

// WithResultSetRowValues returns a functional option for resultsetRow.
func WithResultSetRowValues(values []any) ResultSetRowOptions {
	return func(row *resultsetRow) {
		row.values = values
	}
}

// NewResultSetRow returns a new resultsetRow.
func NewResultSetRow(opts ...ResultSetRowOptions) ResultSetRow {
	row := &resultsetRow{
		values: []any{},
	}
	for _, opt := range opts {
		opt(row)
	}
	return row
}

// Values returns the values.
func (row *resultsetRow) Values() []any {
	return row.values
}

// ValueAt returns the value at the specified index.
func (row *resultsetRow) ValueAt(index int) (any, error) {
	if len(row.values) <= index {
		return nil, errors.ErrNotExist
	}
	return row.values[index], nil
}

// Scan scans the values.
func (row *resultsetRow) Scan(tos ...any) error {
	for i, to := range tos {
		if len(row.values) <= i {
			return errors.ErrNotExist
		}
		v := row.values[i]
		err := safecast.To(v, to)
		if err != nil {
			return err
		}
	}
	return nil
}

// ScanAt scans the value at the specified index.
func (row *resultsetRow) ScanAt(index int, to any) error {
	v, err := row.ValueAt(index)
	if err != nil {
		return err
	}
	return safecast.To(v, to)
}
