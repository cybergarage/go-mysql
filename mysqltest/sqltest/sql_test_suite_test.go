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
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysqltest/client"
	"github.com/cybergarage/go-mysql/mysqltest/server"
)

const sqlTestDatabase = "test"

// TestSQLTestSuite runs already passed scenario test files.
func TestSQLTestSuite(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	server := server.NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Stop()

	cs, err := NewSQLTestSuiteWithDirectory(SQLTestSuiteDefaultTestDirectory)
	if err != nil {
		t.Error(err)
		return
	}

	client := client.NewDefaultClient()
	client.SetDatabase(sqlTestDatabase)
	err = client.CreateDatabase(sqlTestDatabase)
	if err != nil {
		t.Error(err)
	}

	cs.SetClient(client)

	for _, test := range cs.Tests {
		t.Run(test.Name(), func(t *testing.T) {
			test.SetClient(cs.client)
			err := test.Run()
			if err != nil {
				t.Errorf("%s : %s", test.Name(), err.Error())
			}
		})
	}

	err = client.DropDatabase(sqlTestDatabase)
	if err != nil {
		t.Error(err)
	}
}
