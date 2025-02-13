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

// stmtManager represents a statement manager.
type stmtManager struct {
	lastStatementID StatementID
	// stmts is the map of prepared statements.
	stmts map[StatementID]PreparedStatement
}

// NewStatementManager creates a new statement instance.
func NewStatementManager() StatementManager {
	return &stmtManager{
		lastStatementID: 1,
		stmts:           make(map[StatementID]PreparedStatement),
	}
}

// NextPreparedStatementID returns the next prepared statement ID.
func (mgr *stmtManager) NextPreparedStatementID() (StatementID, error) {
	return mgr.lastStatementID.NextStatementID()
}

// RegisterPreparedStatement adds a prepared statement to the manager.
func (mgr *stmtManager) RegisterPreparedStatement(stmt PreparedStatement) error {
	mgr.stmts[stmt.StatementID()] = stmt
	return nil
}

// PreparedStatement returns a prepared statement by the statement ID.
func (mgr *stmtManager) PreparedStatement(stmtID StatementID) (PreparedStatement, error) {
	stmt, ok := mgr.stmts[stmtID]
	if !ok {
		return nil, newErrInvalidStatementID(stmtID)
	}
	return stmt, nil
}

// RemovePreparedStatement removes a prepared statement by the statement ID.
func (mgr *stmtManager) RemovePreparedStatement(stmtID StatementID) {
	delete(mgr.stmts, stmtID)
}
