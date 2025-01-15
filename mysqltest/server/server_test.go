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
	"testing"

	"github.com/cybergarage/go-authenticator/auth"
	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysql"
)

const (
	clientKey  = "../certs/key.pem"
	clientCert = "../certs/cert.pem"
)

var testQueries []string = []string{
	"CREATE DATABASE IF NOT EXISTS ycsb",
	"USE ycsb",
	"CREATE TABLE usertable (YCSB_KEY VARCHAR(255) PRIMARY KEY, FIELD0 TEXT, FIELD1 TEXT, FIELD2 TEXT, FIELD3 TEXT, FIELD4 TEXT, FIELD5 TEXT, FIELD6 TEXT, FIELD7 TEXT, FIELD8 TEXT, FIELD9 TEXT)",
	"DROP TABLE usertable",
	"DROP DATABASE ycsb",
}

func TestServer(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	const (
		username = "testuser"
		password = "testpassword"
	)

	settings := []struct {
		isTLSEnabled      bool
		isPasswordEnabled bool
	}{
		{
			isTLSEnabled:      true,
			isPasswordEnabled: false,
		},
		{
			isTLSEnabled:      false,
			isPasswordEnabled: false,
		},
		{
			isTLSEnabled:      false,
			isPasswordEnabled: true,
		},
		{
			isTLSEnabled:      true,
			isPasswordEnabled: true,
		},
	}

	for _, setting := range settings {
		t.Logf("TLS enabled: %v", setting.isTLSEnabled)
		t.Logf("Password enabled: %v", setting.isPasswordEnabled)

		server := NewServer()

		if setting.isPasswordEnabled {
			cred := auth.NewCredential(
				auth.WithCredentialUsername(username),
				auth.WithCredentialPassword(password),
			)
			server.SetCredentialStore(server)
			server.SetCredential(cred)
		}

		err := server.Start()
		if err != nil {
			t.Error(err)
			return
		}

		client := mysql.NewClient()
		if setting.isTLSEnabled {
			client.SetClientKeyFile(clientKey)
			client.SetClientCertFile(clientCert)
			client.SetRootCertFile(rootCert)
		}
		if setting.isPasswordEnabled {
			client.SetUser(username)
			client.SetPassword(password)
		}
		client.SetDatabase("ycsb")

		err = client.Open()
		defer client.Close()
		if err != nil {
			t.Error(err)
			return
		}

		err = client.Ping()
		if err != nil {
			t.Error(err)
			return
		}

		for n, query := range testQueries {
			t.Logf("[%d] %s", n, query)
			rows, err := client.Query(query)
			if err != nil {
				t.Error(err)
			}
			rows.Close()
		}

		err = server.Stop()
		if err != nil {
			t.Error(err)
			return
		}
	}
}
