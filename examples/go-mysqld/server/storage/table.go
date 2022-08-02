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

package storage

import (
	"fmt"
	"sync"

	"github.com/cybergarage/go-mysql/mysql/query"
)

type Row = query.Row
type Rows = query.Rows
type Condition = query.Condition

// Table represents a destination or source database of query.
type Table struct {
	sync.Mutex
	value string
	*Rows
}

// NewTableWithName returns a new database with the specified string.
func NewTableWithName(name string) *Table {
	tbl := &Table{
		value: name,
		Rows:  query.NewRows(),
	}
	return tbl
}

// NewTable returns a new database.
func NewTable() *Table {
	return NewTableWithName("")
}

// GetName returns the database name.
func (tbl *Table) GetName() string {
	return tbl.value
}

// Insert adds a row.
func (tbl *Table) Insert(row *Row) bool {
	return tbl.AddRow(row)
}

// FindMatchedRows returns only matched and rows by the specified conditions.
func (tbl *Table) FindMatchedRows(cond *Condition) (*Rows, error) {
	matchedRows := query.NewRows()
	for _, row := range tbl.GetRows() {
		if row.IsMatched(cond) {
			matchedRows.AddRow(row)
		}
	}
	return matchedRows, nil
}

// Select returns only matched and projected rows by the specified conditions and the columns.
func (tbl *Table) Select(cond *Condition) (*Rows, error) {
	rows, err := tbl.FindMatchedRows(cond)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Update updates rows which are satisfied by the specified columns and conditions.
func (tbl *Table) Update(q string) int {
	return 0
}

// Delete deletes rows which are satisfied by the specified conditions.
func (tbl *Table) Delete() int {
	return 0
}

// DeleteAll deletes all rows in the table.
func (tbl *Table) DeleteAll() int {
	rows := tbl.GetRows()
	nRowsCnt := len(rows)
	tbl.Rows = query.NewRows()
	return nRowsCnt
}

// String returns the string representation.
func (tbl *Table) String() string {
	return tbl.value
}

// Dump outputs all row values for debug.
func (tbl *Table) Dump() {
	fmt.Printf("%s\n", tbl.GetName())
	for n, row := range tbl.GetRows() {
		fmt.Printf("[%d] %s\n", n, row.String())
	}
}
