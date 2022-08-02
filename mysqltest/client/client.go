// Copyright (C) 2020 Satoshi Konno. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
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
)

// baseClient represents a base client.
type baseClient interface {
	// SetHost sets a host address.
	SetHost(host string)
	// SetPort sets a listen port.
	SetPort(port int)
	// SetDatabase sets a host database.
	SetDatabase(db string)
	// Open opens a database specified by the internal configuration.
	Open() error
	// Close closes opens a database specified by the internal configuration.
	Close() error
	// Query executes a query that returns rows.
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

// Client represents a client.
type Client interface {
	baseClient
	// CreateDatabase creates a specified database.
	CreateDatabase(name string) error
	// DropDatabase dtops a specified database.
	DropDatabase(name string) error
}

type deaultClient struct {
	baseClient
}

// NewDefaultClient returns a default client.
func NewDefaultClient() Client {
	client := &deaultClient{
		baseClient: NewClientDB(),
	}
	return client
}

// CreateDatabase creates a specified database.
func (client *deaultClient) CreateDatabase(name string) error {
	query := fmt.Sprintf("CREATE DATABASE %s", name)
	_, err := client.Query(query)
	return err
}

// DropDatabase dtops a specified database.
func (client *deaultClient) DropDatabase(name string) error {
	query := fmt.Sprintf("DROP DATABASE %s", name)
	_, err := client.Query(query)
	return err
}
