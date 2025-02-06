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

// StatementManager represents a prepared statement manager.
type StatementManager struct {
	lastStatementID StatementID
	// stmts is the map of prepared statements.
	stmts map[StatementID]PreparedStatement
}

// NewStatementManager creates a new PreparedManager instance.
func NewStatementManager() *StatementManager {
	return &StatementManager{
		lastStatementID: 1,
		stmts:           make(map[StatementID]PreparedStatement),
	}
}

// NextPreparedStatementID returns the next prepared statement ID.
func (mgr *StatementManager) NextPreparedStatementID() (StatementID, error) {
	return mgr.lastStatementID.NextStatementID()
}

// AddPreparedStatement adds a prepared statement to the manager.
func (mgr *StatementManager) AddPreparedStatement(stmt PreparedStatement) {
	mgr.stmts[stmt.StatementID()] = stmt
}

// PreparedStatement returns a prepared statement by the statement ID.
func (mgr *StatementManager) PreparedStatement(stmtID StatementID) (PreparedStatement, error) {
	stmt, ok := mgr.stmts[stmtID]
	if !ok {
		return nil, newInvalidStatementID(stmtID)
	}
	return stmt, nil
}

// RemovePreparedStatement removes a prepared statement by the statement ID.
func (mgr *StatementManager) RemovePreparedStatement(stmtID StatementID) {
	delete(mgr.stmts, stmtID)
}
