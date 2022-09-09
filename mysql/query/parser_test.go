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

package query

import (
	"testing"
)

var testQueries []string = []string{
	"CREATE DATABASE IF NOT EXISTS ycsb",
	"USE ycsb",
	"CREATE TABLE IF NOT EXISTS usertable (YCSB_KEY VARCHAR(255) PRIMARY KEY, FIELD0 TEXT, FIELD1 TEXT, FIELD2 TEXT, FIELD3 TEXT, FIELD4 TEXT, FIELD5 TEXT, FIELD6 TEXT, FIELD7 TEXT, FIELD8 TEXT, FIELD9 TEXT)",
	"DROP TABLE ycsb",
	"DROP DATABASE ycsb",
}

func TestQuery(t *testing.T) {
	parser := NewParser()

	for n, query := range testQueries {
		_, err := parser.Parse(query)
		if err != nil {
			t.Logf("[%d] %s", n, query)
			t.Error(err)
		}
	}
}
