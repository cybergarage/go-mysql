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
	"go-mysql/mysql"
	"testing"
)

var testYCSBSetupQueries []string = []string{
	"CREATE DATABASE ycsb",
	"USE ycsb",
	"CREATE TABLE usertable (YCSB_KEY VARCHAR(255) PRIMARY KEY, FIELD0 TEXT, FIELD1 TEXT, FIELD2 TEXT, FIELD3 TEXT, FIELD4 TEXT, FIELD5 TEXT, FIELD6 TEXT, FIELD7 TEXT, FIELD8 TEXT, FIELD9 TEXT)",
	"SELECT * FROM usertable",
}

func testYCSBSetup(t *testing.T) {
	client := mysql.NewClient()
	err := client.Open()
	defer client.Close()
	if err != nil {
		t.Error(err)
	}

	for n, query := range testYCSBSetupQueries {
		t.Logf("[%d] %s", n, query)
		_, err := client.Query(query)
		if err != nil {
			t.Error(err)
		}
	}
}

func testYCSBLoadWorkload(t *testing.T) {
}

func testYCSBRunWorkload(t *testing.T) {
}

func TestYCSB(t *testing.T) {
	server := NewServer()
	server.SetQueryHandler(server)

	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	testYCSBSetup(t)
	testYCSBLoadWorkload(t)
	testYCSBRunWorkload(t)

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}
