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

package server

import (
	"fmt"
	"go-mysql/mysql"
	"testing"
)

const (
	testYCSBDataSize = 100
)

func testYCSBSetup(t *testing.T, client *mysql.Client) {
	var setupQueries []string = []string{
		"CREATE DATABASE ycsb",
		"USE ycsb",
		"CREATE TABLE usertable (YCSB_KEY VARCHAR(255) PRIMARY KEY, FIELD0 TEXT, FIELD1 TEXT, FIELD2 TEXT, FIELD3 TEXT, FIELD4 TEXT, FIELD5 TEXT, FIELD6 TEXT, FIELD7 TEXT, FIELD8 TEXT, FIELD9 TEXT)",
	}

	for n, query := range setupQueries {
		t.Logf("[%d] %s", n, query)
		_, err := client.Query(query)
		if err != nil {
			t.Error(err)
		}
	}
}

func testYCSBLoadWorkload(t *testing.T, client *mysql.Client) {
	for n := 0; n < testYCSBDataSize; n++ {
		query := fmt.Sprintf("INSERT INTO usertable VALUES (YCSB_KEY =\"%d\", FIELD0 =\"%d\",FIELD1 =\"%d\",FIELD2 =\"%d\",FIELD3 =\"%d\",FIELD4 =\"%d\",FIELD5 =\"%d\",FIELD6 =\"%d\",FIELD7 =\"%d\",FIELD8 =\"%d\",FIELD9 =\"%d\")",
			n, n, n, n, n, n, n, n, n, n, n)
		_, err := client.Query(query)
		if err != nil {
			t.Error(err)
		}
	}
}

func testYCSBRunWorkload(t *testing.T, client *mysql.Client) {
	for n := 0; n < testYCSBDataSize; n++ {
		query := fmt.Sprintf("SELECT FIELD0, FIELD1, FIELD2, FIELD3, FIELD4, FIELD5, FIELD6, FIELD7, FIELD8, FIELD9 FROM usertable WHERE YCSB_KEY = \"%d\"",
			n)
		_, err := client.Query(query)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestYCSB(t *testing.T) {
	server := NewServer()
	server.SetQueryHandler(server)

	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}

	client := mysql.NewClient()
	err = client.Open()
	defer client.Close()
	if err != nil {
		t.Error(err)
	}

	testYCSBSetup(t, client)
	testYCSBLoadWorkload(t, client)
	testYCSBRunWorkload(t, client)

	err = server.Stop()
	if err != nil {
		t.Error(err)
		return
	}
}
