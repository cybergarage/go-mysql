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

package sqltest

import (
	"os"
	"strings"
)

// Line represents a line of a SQL test file.
type Line = string

// SQLTestFile represents a SQL test file.
type SQLTestFile struct {
}

// NewSQLTestFile return a file instance.
func NewSQLTestFile() *SQLTestFile {
	file := &SQLTestFile{}
	return file
}

// LoadFile loads the specified test file.
func (file *SQLTestFile) LoadFile(filename string) ([]Line, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return file.ParseBytes(fileBytes)
}

// ParseBytes parses the specified line bytes.
func (file *SQLTestFile) ParseBytes(fileBytes []byte) ([]Line, error) {
	lines := make([]Line, 0)
	for _, line := range strings.Split(string(fileBytes), "\n") {
		// Skip blank or comment lines
		if len(line) == 0 || strings.HasPrefix(line, "-") {
			continue
		}
		lines = append(lines, line)
	}
	return lines, nil
}
