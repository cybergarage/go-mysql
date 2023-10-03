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
	"github.com/google/uuid"
	vitessmy "vitess.io/vitess/go/mysql"
)

// Conn represents a connection of MySQL binary protocol.
type Conn struct {
	*vitessmy.Conn
	db   string
	uid  uint32
	uuid uuid.UUID
	ts   time.Time
	sync.Map
	tracer.Context
}

// newConn returns a connection with a default empty connection.
func newConn() *Conn {
	return NewConnWith(
		tracer.NewNullTracer().StartSpan(""),
		&vitessmy.Conn{},
	)
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(ctx tracer.Context, c *vitessmy.Conn) *Conn {
	conn := &Conn{
		Conn:    c,
		uid:     0,
		ts:      time.Now(),
		uuid:    uuid.New(),
		Map:     sync.Map{},
		Context: ctx,
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

// UUID returns the UUID of the connection.
func (conn *Conn) UUID() uuid.UUID {
	return conn.uuid
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

// SetSpanContext sets the tracer span context to the connection.
func (conn *Conn) SetSpanContext(span tracer.Context) {
	conn.Context = span
}

// SpanContext returns the tracer span context of the connection.
func (conn *Conn) SpanContext() tracer.Context {
	return conn.Context
}
