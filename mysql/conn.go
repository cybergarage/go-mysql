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
	"sync"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
	vitessmy "vitess.io/vitess/go/mysql"
)

// Conn represents a connection of MySQL binary protocol.
type Conn struct {
	*vitessmy.Conn
	db  string
	uid uint32
	ts  time.Time
	sync.Map
	span tracer.SpanContext
}

// newConn returns a connection with a default empty connection.
func newConn() *Conn {
	return NewConnWithConn(&vitessmy.Conn{})
}

// NewConnWithConn returns a connection with a raw connection.
func NewConnWithConn(c *vitessmy.Conn) *Conn {
	conn := &Conn{
		Conn: c,
		uid:  0,
		ts:   time.Now(),
		Map:  sync.Map{},
	}

	if c != nil {
		conn.uid = c.ConnectionID
	}

	return conn
}

// SetDatabase sets th selected database to the connection.
func (conn *Conn) SetDatabase(name string) {
	conn.db = name
}

// Database returns the current selected database in the connection.
func (conn *Conn) Database() string {
	return conn.db
}

// ID returns the creation ID of the connection.
func (conn *Conn) ID() uint32 {
	return conn.uid
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

// SetSpanContext sets a span context to the connection.
func (conn *Conn) SetSpanContext(span tracer.SpanContext) {
	conn.span = span
}

// SpanContext returns the span context of the connection.
func (conn *Conn) SpanContext() tracer.SpanContext {
	return conn.span
}
