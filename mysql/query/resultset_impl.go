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

type resultset struct {
	ResultSetSchema
	ResultSetRow
	rowsAffected int64
}

// ResultSet represents a response resultset interface.
type ResultSetOption func(*resultset)

// WithRowsAffected returns a resultset option to set the rows affected.
func WithRowsAffected(rowsAffected int64) ResultSetOption {
	return func(r *resultset) {
		r.rowsAffected = rowsAffected
	}
}

// NewResultSet returns a new ResultSet.
func NewResultSet() ResultSet {
	return &resultset{
		ResultSetSchema: nil,
		ResultSetRow:    nil,
		rowsAffected:    0,
	}
}

// RowsAffected returns the number of rows affected.
func (r *resultset) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

// Next returns the next row.
func (r *resultset) Next() bool {
	return false
}

// Row returns the current row.
func (r *resultset) Row() ResultSetRow {
	return nil
}

// Schema returns the schema.
func (r *resultset) Schema() ResultSetSchema {
	return nil
}

// Close closes the resultset.
func (r *resultset) Close() error {
	return nil
}
