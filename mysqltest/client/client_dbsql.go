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

package client

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// ClientDB represents a client.
type ClientDB struct {
	*Config
	db *sql.DB
}

// NewClientDB returns a client instance.
func NewClientDB() *ClientDB {
	client := &ClientDB{
		Config: NewDefaultConfig(),
		db:     nil,
	}
	return client
}

// Open opens a database specified by the internal configuration.
func (client *ClientDB) Open() error {
	dsName := fmt.Sprintf("root@tcp(%s:%d)/%s", client.Host, client.Port, client.Database)
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		return err
	}
	client.db = db
	return nil
}

// Close closes opens a database specified by the internal configuration.
func (client *ClientDB) Close() error {
	if client.db == nil {
		return nil
	}
	return client.db.Close()
}

// Query executes a query that returns rows.
func (client *ClientDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if client.db == nil {
		err := client.Open()
		if err != nil {
			return nil, err
		}
		defer client.Close()
	}
	return client.db.Query(query, args...)
}
