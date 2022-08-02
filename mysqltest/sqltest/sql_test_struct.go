// Copyright (C) 2020 Satoshi Konno. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqltest

import (
	"fmt"

	"github.com/cybergarage/go-mysql/mysqltest/client"
)

const (
	SQLTestFileExt = "test"
)

// SQLTest represents a SQL test.
type SQLTest struct {
	Scenario *SQLScenario
	client   *client.Client
}

// NewSQLTest returns a SQL test instance.
func NewSQLTest() *SQLTest {
	test := &SQLTest{}
	return test
}

// NewSQLTestWithFile return a SQL test instance for the specified test scenario file.
func NewSQLTestWithFile(filename string) (*SQLTest, error) {
	file := NewSQLTest()
	err := file.LoadFile(filename)
	return file, err
}

// SetClient sets a client for testing.
func (ct *SQLTest) SetClient(c *client.Client) {
	ct.client = c
}

// GetName returns the loaded senario name.
func (ct *SQLTest) GetName() string {
	return ct.Scenario.GetName()
}

// LoadFile loads a specified SQL test file.
func (ct *SQLTest) LoadFile(filename string) error {
	ct.Scenario = NewSQLScenario()
	return ct.Scenario.LoadFile(filename)
}

// LoadFileWithBasename loads a SQL test file which has specified basename.
func (ct *SQLTest) LoadFileWithBasename(basename string) error {
	return ct.LoadFile(basename + "." + SQLTestFileExt)
}

// Run runs a loaded scenario test.
func (ct *SQLTest) Run() error {
	scenario := ct.Scenario
	if scenario == nil {
		return nil
	}

	err := scenario.IsValid()
	if err != nil {
		return err
	}

	client := ct.client
	if client == nil {
		return fmt.Errorf(errorClientNotFound)
	}

	err = client.Open()
	if err != nil {
		return err
	}

	errTraceMsg := func(n int) string {
		errTraceMsg := ct.GetName() + "\n"
		for i := 0; i < n; i++ {
			errTraceMsg += fmt.Sprintf(goodQueryPrefix, i, scenario.Queries[i])
			errTraceMsg += "\n"
		}
		return errTraceMsg
	}

	for n, query := range scenario.Queries {
		rs, err := client.Query(query)
		if err != nil {
			return fmt.Errorf("%s%w", errTraceMsg(n), err)
		}
		defer rs.Close()

		names, err := rs.Columns()
		if err != nil {
			return err
		}
		columnCnt := len(names)

		rsRows := make([]interface{}, 0)
		for rs.Next() {
			values := make([]interface{}, columnCnt)
			err := rs.Scan(values...)
			if err != nil {
				return err
			}
			row := map[string]interface{}{}
			for i := 0; i < columnCnt; i++ {
				row[names[i]] = values[i]
			}
			rsRows = append(rsRows, row)
		}

		expectedRes := scenario.Results[n]
		expectedRows, err := expectedRes.GetRows()
		if err != nil {
			if len(rsRows) != 0 {
				return fmt.Errorf("%s"+errorJSONResponseHasUnexpectedRows, errTraceMsg(n), n, query, rsRows)
			}
		} else {
			if len(rsRows) != len(expectedRows) {
				return fmt.Errorf("%s"+errorJSONResponseUnmatchedRowCount, errTraceMsg(n), n, query, rsRows, expectedRows)
			}
		}

		for _, row := range rsRows {
			err = expectedRes.HasRow(row)
			if err != nil {
				return fmt.Errorf("%s"+errorQueryPrefix+"%s", errTraceMsg(n), n, query, err.Error())
			}
		}
	}

	err = client.Close()
	if err != nil {
		return err
	}

	return nil
}
