// Copyright (C) 2020 Satoshi Konno. All rights reserved.
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

// Rows represents a row array.
type Rows struct {
	list []*Row
}

// NewRows returns a row array.
func NewRows() *Rows {
	return NewRowsWithRows(make([]*Row, 0))
}

// NewRowsWithRows returns a row array with the specified rows.
func NewRowsWithRows(rawRows []*Row) *Rows {
	rows := &Rows{
		list: rawRows,
	}
	return rows
}

// AddRow adds a row.
func (rows *Rows) AddRow(row *Row) bool {
	rows.list = append(rows.list, row)
	return true
}

// GetRows returns the all rows.
func (rows *Rows) GetRows() []*Row {
	return rows.list
}

// DeleteRow deletes the specified row.
func (rows *Rows) DeleteRow(targetRow *Row) int64 {
	for n, row := range rows.list {
		if row == targetRow {
			rows.list = append(rows.list[:n], rows.list[n+1:]...)
			return 1
		}
	}
	return 0
}

// GetRow returns a row of the specified index.
func (rows *Rows) GetRow(n int) (*Row, bool) {
	if len(rows.list) <= n {
		return nil, false
	}
	return rows.list[n], true
}
