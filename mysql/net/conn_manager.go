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
	"github.com/google/uuid"
)

// ConnManager represents a connection map.
type ConnManager interface {
	// AddConn adds the specified connection.
	AddConn(c Conn) error
	// UpdateConn updates the specified connection.
	UpdateConn(from Conn, to Conn) error
	// Conns returns the included connections.
	Conns() []Conn
	// LookupConnByUID returns a connection and true when the specified connection exists by the connection ID, otherwise nil and false.
	LookupConnByUID(cid uint64) (Conn, bool)
	// LookupConnByUUID returns the connection with the specified UUID.
	LookupConnByUUID(uuid uuid.UUID) (Conn, bool)
	// RemoveConn deletes the specified connection from the map.
	RemoveConn(conn Conn) error
	// RemoveConnByUID deletes the specified connection by the connection ID.
	RemoveConnByUID(cid uint64)
	// RemoveConnByUID deletes the specified connection by the connection ID.
	RemoveConnByUUID(uuid uuid.UUID)
	// Start starts the connection manager.
	Start() error
	// Close closes the connection manager.
	Close() error
	// Stop closes all connections.
	Stop() error
}
