// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
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

package net

import (
	"github.com/cybergarage/go-sqlparser/sql/net"
	"github.com/google/uuid"
)

type connManager struct {
	net.ConnManager
}

// NewConnManager returns a new connection manager instance.
func NewConnManager() ConnManager {
	return &connManager{
		ConnManager: net.NewConnManager(),
	}
}

// AddConn adds the specified connection.
func (mgr *connManager) AddConn(c Conn) error {
	return mgr.ConnManager.AddConn(c)
}

// UpdateConn updates the specified connection.
func (mgr *connManager) UpdateConn(from Conn, to Conn) error {
	return mgr.ConnManager.UpdateConn(from, to)
}

// Conns returns the included connections.
func (mgr *connManager) Conns() []Conn {
	conns := mgr.ConnManager.Conns()
	ret := make([]Conn, len(conns))
	for i, c := range conns {
		ret[i] = c.(Conn) // nolint: forcetypeassert
	}
	return ret
}

// LookupConnByUID returns a connection and true when the specified connection exists by the connection ID, otherwise nil and false.
func (mgr *connManager) LookupConnByUID(cid uint64) (Conn, bool) {
	c, ok := mgr.ConnManager.LookupConnByUID(cid)
	if c == nil {
		return nil, ok
	}
	return c.(Conn), ok // nolint: forcetypeassert
}

// LookupConnByUUID returns the connection with the specified UUID.
func (mgr *connManager) LookupConnByUUID(uuid uuid.UUID) (Conn, bool) {
	c, ok := mgr.ConnManager.LookupConnByUUID(uuid)
	if c == nil {
		return nil, ok
	}
	return c.(Conn), ok // nolint: forcetypeassert
}

// RemoveConn deletes the specified connection from the map.
func (mgr *connManager) RemoveConn(conn Conn) error {
	return mgr.ConnManager.RemoveConn(conn)
}
