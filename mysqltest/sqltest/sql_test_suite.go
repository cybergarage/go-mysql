// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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
	"regexp"

	"github.com/cybergarage/go-mysql/mysqltest/client"
	"github.com/cybergarage/go-mysql/mysqltest/util"
)

const (
	SQLTestSuiteDefaultTestDirectory = "./tests"
)

// SQLTestSuite represents a SQL test suite.
type SQLTestSuite struct {
	Tests  []*SQLTest
	client *client.Client
}

// NewSQLTestSuite returns a SQL test suite instance.
func NewSQLTestSuite() *SQLTestSuite {
	suite := &SQLTestSuite{
		Tests: make([]*SQLTest, 0),
	}
	return suite
}

// NewSQLTestSuiteWithDirectory returns a SQL test suite instance which loads under the specified directory.
func NewSQLTestSuiteWithDirectory(dir string) (*SQLTestSuite, error) {
	suite := NewSQLTestSuite()
	err := suite.LoadDirectory(dir)
	return suite, err
}

// SetClient sets a client for testing.
func (suite *SQLTestSuite) SetClient(c *client.Client) {
	suite.client = c
}

// LoadDirectory loads all test files in the specified directory.
func (suite *SQLTestSuite) LoadDirectory(dir string) error {
	findPath := util.NewFileWithPath(dir)

	re, err := regexp.Compile(".*\\." + SQLTestFileExt)
	if err != nil {
		return err
	}

	files, err := findPath.ListFilesWithRegexp(re)
	if err != nil {
		return err
	}

	suite.Tests = make([]*SQLTest, 0)
	for _, file := range files {
		cs, err := NewSQLTestWithFile(file.Path)
		if err != nil {
			return err
		}
		suite.Tests = append(suite.Tests, cs)
	}

	return nil
}

// Run runs all loaded scenario tests. The method stops the testing when a scenario test is aborted, and the following tests are not run.
func (suite *SQLTestSuite) Run() error {
	for _, test := range suite.Tests {
		test.SetClient(suite.client)
		err := test.Run()
		if err != nil {
			return fmt.Errorf("%s : %w", test.Name(), err)
		}
	}
	return nil
}
