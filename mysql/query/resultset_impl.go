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
	rows         []ResultSetRow
	rowsAffected uint64
	rowCursor    int
}

// ResultSet represents a response resultset interface.
type ResultSetOption func(*resultset)

// WithResultSetRowsAffected returns a resultset option to set the rows affected.
func WithResultSetRowsAffected(rowsAffected uint64) ResultSetOption {
	return func(r *resultset) {
		r.rowsAffected = rowsAffected
	}
}

// WithResultSetSchema returns a resultset option to set the schema.
func WithResultSetSchema(schema ResultSetSchema) ResultSetOption {
	return func(r *resultset) {
		r.ResultSetSchema = schema
	}
}

// WithResultSetRows returns a resultset option to set the rows.
func WithResultSetRows(rows []ResultSetRow) ResultSetOption {
	return func(r *resultset) {
		r.rows = rows
	}
}

// NewResultSet returns a new ResultSet.
func NewResultSet(opts ...ResultSetOption) ResultSet {
	rs := &resultset{
		ResultSetSchema: nil,
		rows:            []ResultSetRow{},
		rowsAffected:    0,
		rowCursor:       0,
	}
	for _, opt := range opts {
		opt(rs)
	}
	return rs
}

// RowsAffected returns the number of rows affected.
func (r *resultset) RowsAffected() uint64 {
	return r.rowsAffected
}

// Next returns the next row.
func (r *resultset) Next() bool {
	r.rowCursor++
	return (r.rowCursor - 1) < len(r.rows)
}

// Row returns the current row.
func (r *resultset) Row() ResultSetRow {
	if len(r.rows) < (r.rowCursor - 1) {
		return nil
	}
	return r.rows[r.rowCursor-1]
}

// Schema returns the schema.
func (r *resultset) Schema() ResultSetSchema {
	return r.ResultSetSchema
}

// Close closes the resultset.
func (r *resultset) Close() error {
	return nil
}
