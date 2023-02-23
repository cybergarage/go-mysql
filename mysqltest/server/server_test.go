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

package server

import (
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysqltest/client"
)

var testQueries []string = []string{
	"CREATE DATABASE IF NOT EXISTS ycsb",
	"USE ycsb",
	"CREATE TABLE usertable (YCSB_KEY VARCHAR(255) PRIMARY KEY, FIELD0 TEXT, FIELD1 TEXT, FIELD2 TEXT, FIELD3 TEXT, FIELD4 TEXT, FIELD5 TEXT, FIELD6 TEXT, FIELD7 TEXT, FIELD8 TEXT, FIELD9 TEXT)",
	"DROP TABLE ycsb",
	"DROP DATABASE ycsb",
}

func TestServer(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	server := NewServer()

	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	client := client.NewDefaultClient()
	client.SetDatabase("ycsb")
	err = client.Open()
	defer client.Close()
	if err != nil {
		t.Error(err)
	}

	for n, query := range testQueries {
		t.Logf("[%d] %s", n, query)
		rows, err := client.Query(query)
		if err != nil {
			t.Error(err)
		}
		defer rows.Close()
	}

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}
