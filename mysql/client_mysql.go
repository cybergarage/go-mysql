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

package mysql

import (
	"database/sql/driver"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

// Client represents a client for MySQL server.
type MySQLClient struct {
	*Config
	conn driver.Conn
}

// NewMySQLClient returns a client instance.
func NewMySQLClient() *MySQLClient {
	client := &MySQLClient{
		Config: NewDefaultConfig(),
		conn:   nil,
	}
	return client
}

// Open opens a database specified by the internal configuration.
func (client *MySQLClient) Open(dbName string) error {
	dbDrv := mysql.MySQLDriver{}
	dsName := fmt.Sprintf("root@tcp(127.0.0.1:3306)/%s", dbName)
	conn, err := dbDrv.Open(dsName)
	if err != nil {
		return err
	}
	client.conn = conn
	return nil
}

// Close closes opens a database specified by the internal configuration.
func (client *MySQLClient) Close() error {
	if client.conn == nil {
		return nil
	}
	return client.conn.Close()
}

// Query executes a query that returns rows.
func (client *MySQLClient) Query(query string, args ...interface{}) (driver.Rows, error) {
	if client.conn == nil {
		return nil, nil
	}
	stmt, err := client.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	queryArgs := []driver.Value{}
	rows, err := stmt.Query(queryArgs)
	return rows, err
}
