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
	"fmt"
	"strings"
)

// SQLScenario represents a scenario.
type SQLScenario struct {
	*SQLTestFile
	Filename string
	Queries  []string
	Results  []*SQLResponse
}

// NewSQLScenario return a scenario instance.
func NewSQLScenario() *SQLScenario {
	file := &SQLScenario{
		SQLTestFile: NewSQLTestFile(),
	}
	return file
}

// NewSQLScenarioWithFile return a scenario instance for the specified test scenario file.
func NewSQLScenarioWithFile(filename string) (*SQLScenario, error) {
	file := NewSQLScenario()
	err := file.LoadFile(filename)
	return file, err
}

// GetName returns the loaded scenario file name.
func (file *SQLScenario) GetName() string {
	return file.Filename
}

// IsValid checks whether the loaded scenario is available.
func (file *SQLScenario) IsValid() error {
	if len(file.Queries) != len(file.Results) {
		return fmt.Errorf(errorInvalidScenarioCases, len(file.Queries), len(file.Results))
	}
	return nil
}

// LoadFile loads the specified scenario.
func (file *SQLScenario) LoadFile(filename string) error {
	lines, err := file.SQLTestFile.LoadFile(filename)
	if err != nil {
		return err
	}

	err = file.ParseLineStrings(lines)
	if err != nil {
		return fmt.Errorf("%s : %w", filename, err)
	}

	file.Filename = filename

	return nil
}

// ParseLineStrings parses the specified scenario line strings.
func (file *SQLScenario) ParseLineStrings(lines []string) error {
	var queryStr, resultStr string
	file.Queries = make([]string, 0)
	file.Results = make([]*SQLResponse, 0)

	appendQuery := func() {
		if len(queryStr) <= 0 {
			return
		}
		file.Queries = append(file.Queries, strings.TrimSpace(queryStr))
		queryStr = ""
	}
	appendResult := func() error {
		if len(resultStr) <= 0 {
			return nil
		}
		result, err := NewSQLResponseWithString(strings.TrimSpace(resultStr))
		if err != nil {
			return err
		}
		file.Results = append(file.Results, result)
		resultStr = ""
		return nil
	}

	inJSON := false
	for _, line := range lines {
		if inJSON {
			if strings.HasPrefix(line, "}") {
				resultStr += line
				err := appendResult()
				if err != nil {
					return err
				}
				inJSON = false
				continue
			}
			resultStr += " " + strings.TrimSpace(line)
		} else {
			if strings.HasPrefix(line, "{") {
				appendQuery()
				resultStr = line
				inJSON = true
				continue
			}
			queryStr += " " + strings.TrimSpace(line)
		}
	}

	return file.IsValid()
}

// String returns the string representation.
func (file *SQLScenario) String() string {
	var str string
	nResults := len(file.Results)
	for n, query := range file.Queries {
		str += query + "\n"
		if n < nResults {
			str += file.Results[n].String() + "\n"
		}
	}
	return str
}
