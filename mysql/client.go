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

package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Client represents a client for MySQL server.
type Client struct {
	*Config
	db *sql.DB
}

// NewClient returns a client instance.
func NewClient() *Client {
	client := &Client{
		Config: NewDefaultConfig(),
		db:     nil,
	}
	return client
}

// Open opens a database specified by the internal configuration.
func (client *Client) Open() error {
	dsName := fmt.Sprintf("tcp(%s:%d)/%s", client.Address, client.Port, client.Database)
	db, err := sql.Open("mysql", dsName)
	if err != nil {
		return err
	}

	// See: https://github.com/go-sql-driver/mysql
	// Important settings
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	client.db = db
	return nil
}

// Close closes opens a database specified by the internal configuration.
func (client *Client) Close() error {
	if client.db == nil {
		return nil
	}
	if err := client.db.Close(); err != nil {
		return err
	}
	client.db = nil
	return nil
}

// Query executes a query that returns rows.
func (client *Client) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if client.db == nil {
		err := client.Open()
		if err != nil {
			return nil, err
		}
	}
	return client.db.Query(query, args...)
}
