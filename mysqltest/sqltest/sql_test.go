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
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysqltest/server"
	"github.com/cybergarage/go-sqltest/sqltest"
)

// TestSQLTest is a temporary debug test to check only the specified test cases.
func TestSQLTest(t *testing.T) {
	log.EnableStdoutDebug(true)

	server := server.NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Stop()

	// NOTE: Add your test files in 'untests' directory into the filename array
	testNames := []string{
		// "UpdateArithInt",
		// "SmplCrudDatetime",
		// "FuncAggrInt",
		// "FuncMathInt",
		// "SmplCrudInt",
		// "SmplCrudText",
		// "SmplCrudDatetime",
	}

	sqltest.RunEmbedSuites(t, sqltest.NewMySQLClient(), testNames...)
}
