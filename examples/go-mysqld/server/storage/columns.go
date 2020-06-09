// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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

// GetColumns returns the all columns.
func (cols *Columns) GetColumns() []*Column {
	return cols.allList
}

// GetColumn returns a column of the specified index.
func (cols *Columns) GetColumn(n int) (*Column, bool) {
	if len(cols.allList) <= n {
		return nil, false
	}
	return cols.allList[n], true
}

// GetColumnByName returns a column by the specified name.
func (cols *Columns) GetColumnByName(name string) (*Column, bool) {
	col, ok := cols.cachedMap[name]
	if ok {
		return col, true
	}
	for _, col := range cols.allList {
		if col.GetName() == name {
			cols.cachedMap[name] = col
			return col, true
		}
	}
	return nil, false
}

// HasColumn returns true when the query has the specified column.
func (cols *Columns) HasColumn(name string) bool {
	_, ok := cols.GetColumnByName(name)
	return ok
}

// IsAllColumn returns true when the query has only "*" column.
func (cols *Columns) IsAllColumn() bool {
	if len(cols.allList) == 0 {
		return true
	}
	return cols.HasColumn("*")
}

// Equals returns true when the specified columns are equals to this columns, otherwise false.
func (cols *Columns) Equals(otherCols *Columns) bool {
	for _, col := range cols.allList {
		otherCol, ok := otherCols.GetColumnByName(col.GetName())
		if !ok {
			return false
		}
		if !col.Equals(otherCol) {
			return false
		}
	}
	return true
}
