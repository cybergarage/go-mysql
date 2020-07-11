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

package sqltest

import (
	"path"
	"testing"

	"go-mysql/test/client"
	"go-mysql/test/server"
)

// TestSQLUntestedCases is a temporary debug test to check untested test cases.
func TestSQLUntestedCases(t *testing.T) {
	testFilenames := []string{
		// NOTE: Add your test files in 'untests' directory into the filename array
	}

	server := server.NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Stop()

	client := client.NewDefaultClient()

	for _, testFilename := range testFilenames {
		ct := NewSQLTest()
		err = ct.LoadFile(path.Join(SQLTestSuiteDefaultTestDirectory, testFilename))
		if err != nil {
			t.Error(err)
			continue
		}
		ct.SetClient(client)

		err = ct.Run()
		if err != nil {
			t.Error(err)
		}
	}

}