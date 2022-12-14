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
	"fmt"

	"github.com/cybergarage/go-mysql/mysql"
)

// Client represents a client for MySQL server.
type Client struct {
	*mysql.Client
}

// NewDefaultClient returns a default client instance with the specified host and port.
func NewDefaultClient() *Client {
	client := &Client{
		Client: mysql.NewClient(),
	}
	return client
}

// CreateDatabase creates a specified database.
func (client *Client) CreateDatabase(name string) error {
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", name)
	_, err := client.Query(query)
	return err
}

// DropDatabase dtops a specified database.
func (client *Client) DropDatabase(name string) error {
	query := fmt.Sprintf("DROP DATABASE IF EXISTS %s", name)
	_, err := client.Query(query)
	return err
}

// Use sets a target database.
func (client *Client) Use(name string) error {
	client.Database = name
	return nil
}
