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
	"time"

	vitessmy "vitess.io/vitess/go/mysql"
)

// Conn represents a connection of MySQL binary protocol.
type Conn struct {
	*vitessmy.Conn
	Database  string
	UID       uint32
	Timestamp time.Time
}

// newConn returns a connection with a default empty connection.
func newConn() *Conn {
	return NewConnWithConn(&vitessmy.Conn{})
}

// NewConnWithConn returns a connection with a raw connection.
func NewConnWithConn(c *vitessmy.Conn) *Conn {
	conn := &Conn{
		Conn:      c,
		UID:       0,
		Timestamp: time.Now(),
	}

	if c != nil {
		conn.UID = c.ConnectionID
	}

	return conn
}
