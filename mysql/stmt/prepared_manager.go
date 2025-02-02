// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
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

package stmt

import (
	"sync"
)

type PreparedManager struct {
	mutex           sync.Mutex
	lastStatementID StatementID
	// stmts is the map of prepared statements.
	stmts map[StatementID]PreparedStatement
}

// NewPreparedManager creates a new PreparedManager instance.
func NewPreparedManager() *PreparedManager {
	return &PreparedManager{
		mutex:           sync.Mutex{},
		lastStatementID: 1,
		stmts:           make(map[StatementID]PreparedStatement),
	}
}

// AddPreparedStatement adds a prepared statement to the manager.
func (manager *PreparedManager) AddPreparedStatement(stmt PreparedStatement) {
	manager.stmts[stmt.StatementID()] = stmt
}

// PreparedStatement returns a prepared statement by the statement ID.
func (manager *PreparedManager) PreparedStatement(stmtID StatementID) (PreparedStatement, bool) {
	stmt, ok := manager.stmts[stmtID]
	return stmt, ok
}

// RemovePreparedStatement removes a prepared statement by the statement ID.
func (manager *PreparedManager) RemovePreparedStatement(stmtID StatementID) {
	delete(manager.stmts, stmtID)
}
