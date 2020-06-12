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

// ConnMap represents a connection map.
type ConnMap map[uint32]*Conn

// NewConnMap returns a connection map.
func NewConnMap() ConnMap {
	return ConnMap{}
}

// AddConn adds the specified connection.
func (cm ConnMap) AddConn(c *Conn) {
	cm[c.ConnectionID] = c
}

// GetByConnID returns a connection and true when the specified connection exists by the connection ID, otherwise nil and false.
func (cm ConnMap) GetByConnID(cid uint32) (*Conn, bool) {
	c, ok := cm[cid]
	return c, ok
}

// DeleteByConnID deletes the specified connection by the connection ID.
func (cm ConnMap) DeleteByConnID(cid uint32) {
	delete(cm, cid)
}

// Length returns the included connection count.
func (cm ConnMap) Length() int {
	return len(cm)
}
