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
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var ErrExist = errors.New("exist")
var ErrNotExist = errors.New("not exist")

func newErrConnExist(s string) error {
	return fmt.Errorf("connection (%s) is already %w", s, ErrExist)
}

func newErrConnNotExists(s string) error {
	return fmt.Errorf("connection (%s) is not %w", s, ErrNotExist)
}

// ConnManager represents a connection map.
type ConnManager struct {
	uidMap map[uint64]Conn
	m      map[uuid.UUID]Conn
	mutex  *sync.RWMutex
}

// NewConnManager returns a connection map.
func NewConnManager() *ConnManager {
	return &ConnManager{
		uidMap: map[uint64]Conn{},
		m:      map[uuid.UUID]Conn{},
		mutex:  &sync.RWMutex{},
	}
}

// AddConn adds the specified connection.
func (cm *ConnManager) AddConn(c Conn) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	uuid := c.UUID()
	if _, ok := cm.m[uuid]; ok {
		return newErrConnExist(uuid.String())
	}
	uid := c.ID()
	if _, ok := cm.uidMap[uid]; ok {
		return newErrConnExist(fmt.Sprintf("%d", uid))
	}

	cm.m[uuid] = c
	cm.uidMap[uid] = c

	return nil
}

// UpdateConn updates the specified connection.
func (cm *ConnManager) UpdateConn(from Conn, to Conn) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	uuid := from.UUID()
	if _, ok := cm.m[uuid]; !ok {
		return newErrConnNotExists(uuid.String())
	}
	uid := from.ID()
	if _, ok := cm.uidMap[uid]; !ok {
		return newErrConnNotExists(fmt.Sprintf("%d", uid))
	}

	cm.m[uuid] = to
	cm.uidMap[uid] = to

	return nil
}

// Conns returns the included connections.
func (cm *ConnManager) Conns() []Conn {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	conns := make([]Conn, 0, len(cm.m))
	for _, conn := range cm.m {
		conns = append(conns, conn)
	}
	return conns
}

// LookupConnByUID returns a connection and true when the specified connection exists by the connection ID, otherwise nil and false.
func (cm *ConnManager) LookupConnByUID(cid uint64) (Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	c, ok := cm.uidMap[cid]
	return c, ok
}

// LookupConnByUUID returns the connection with the specified UUID.
func (cm *ConnManager) LookupConnByUUID(uuid uuid.UUID) (Conn, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	c, ok := cm.m[uuid]
	return c, ok
}

// RemoveConn deletes the specified connection from the map.
func (cm *ConnManager) RemoveConn(conn Conn) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	delete(cm.uidMap, conn.ID())
	delete(cm.m, conn.UUID())
	return nil
}

// RemoveConnByUID deletes the specified connection by the connection ID.
func (cm *ConnManager) RemoveConnByUID(cid uint64) {
	conn, ok := cm.LookupConnByUID(cid)
	if !ok {
		return
	}
	cm.RemoveConn(conn)
}

// RemoveConnByUID deletes the specified connection by the connection ID.
func (cm *ConnManager) RemoveConnByUUID(uuid uuid.UUID) {
	conn, ok := cm.LookupConnByUUID(uuid)
	if !ok {
		return
	}
	cm.RemoveConn(conn)
}

// Start starts the connection manager.
func (cm *ConnManager) Start() error {
	return nil
}

// Close closes the connection manager.
func (cm *ConnManager) Close() error {
	var errs error
	conns := cm.Conns()
	for _, conn := range conns {
		err := conn.Close()
		if err == nil {
			if err := cm.RemoveConn(conn); err != nil {
				errs = errors.Join(errs, err)
			}
		} else {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

// Stop closes all connections.
func (cm *ConnManager) Stop() error {
	if err := cm.Close(); err != nil {
		return err
	}
	return nil
}
