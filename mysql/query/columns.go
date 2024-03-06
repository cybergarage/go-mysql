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

import "strings"

// Columns represents a column array.
type Columns struct {
	allList   []*Column
	cachedMap map[string]*Column
}

// NewColumns returns a null columns.
func NewColumns() *Columns {
	cols := &Columns{
		allList:   make([]*Column, 0),
		cachedMap: map[string]*Column{},
	}
	return cols
}

// NewColumnsWithColumns returns a columns with the specified columns.
func NewColumnsWithColumns(columns []*Column) *Columns {
	cols := NewColumns()
	for _, column := range columns {
		cols.AddColumn(column)
	}
	return cols
}

// AddColumn adds a column.
func (cols *Columns) AddColumn(col *Column) {
	cols.allList = append(cols.allList, col)
}

// AddColumns adds columns.
func (cols *Columns) AddColumns(columns []*Column) {
	cols.allList = append(cols.allList, columns...)
}

// AllColumns returns the all columns.
func (cols *Columns) AllColumns() []*Column {
	return cols.allList
}

// Columns returns the all columns.
func (cols *Columns) Columns() []*Column {
	return cols.AllColumns()
}

// ColumnAt returns a column of the specified index.
func (cols *Columns) ColumnAt(n int) (*Column, bool) {
	if len(cols.allList) <= n {
		return nil, false
	}
	return cols.allList[n], true
}

// ColumnByName returns a column by the specified name.
func (cols *Columns) ColumnByName(name string) (*Column, bool) {
	lowerName := strings.ToLower(name)
	col, ok := cols.cachedMap[lowerName]
	if ok {
		return col, true
	}
	for _, col := range cols.allList {
		if lowerName == strings.ToLower(col.Name()) {
			cols.cachedMap[lowerName] = col
			return col, true
		}
	}
	return nil, false
}

// HasColumn returns true when the query has the specified column.
func (cols *Columns) HasColumn(name string) bool {
	_, ok := cols.ColumnByName(name)
	return ok
}

// IsAllColumn returns true when the query has only "*" column.
func (cols *Columns) IsAllColumn() bool {
	if len(cols.allList) == 0 {
		return true
	}
	return cols.HasColumn("*")
}

// Size returns the all column count.
func (cols *Columns) Size() int {
	return len(cols.allList)
}

// Clear clears all columns.
func (cols *Columns) Clear() {
	cols.allList = make([]*Column, 0)
	cols.cachedMap = map[string]*Column{}
}

// Equals returns true when the specified columns are equals to this columns, otherwise false.
func (cols *Columns) Equals(otherCols *Columns) bool {
	for _, col := range cols.allList {
		otherCol, ok := otherCols.ColumnByName(col.Name())
		if !ok {
			return false
		}
		if !col.Equals(otherCol) {
			return false
		}
	}
	return true
}

// String returns the string representation.
func (cols *Columns) String() string {
	str := "{"
	for n, col := range cols.allList {
		if 0 < n {
			str += ", "
		}
		str += col.String()
	}
	str += "}"
	return str
}
