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
	// stmtQueryMap is the map of prepared statements by query.
	stmtQueryMap map[string]PreparedStatement
	// stmtIDMap is the map of prepared statements by statement ID.
	stmtIDMap map[StatementID]PreparedStatement
}

// NewStatementManager creates a new statement instance.
func NewStatementManager() StatementManager {
	return &stmtManager{
		lastStatementID: 1,
		stmtQueryMap:    make(map[string]PreparedStatement),
		stmtIDMap:       make(map[StatementID]PreparedStatement),
	}
}

// NextPreparedStatementID returns the next prepared statement ID.
func (mgr *stmtManager) NextPreparedStatementID() (StatementID, error) {
	return mgr.lastStatementID.NextStatementID()
}

// RegisterPreparedStatement adds a prepared statement to the manager.
func (mgr *stmtManager) RegisterPreparedStatement(stmt PreparedStatement) error {
	mgr.stmtIDMap[stmt.StatementID()] = stmt
	mgr.stmtQueryMap[stmt.Query()] = stmt

	return nil
}

// LookupPreparedStatementByID returns a prepared statement by the statement ID.
func (mgr *stmtManager) LookupPreparedStatementByID(stmtID StatementID) (PreparedStatement, error) {
	stmt, ok := mgr.stmtIDMap[stmtID]
	if !ok {
		return nil, newErrInvalidStatementID(stmtID)
	}

	return stmt, nil
}

// LookupPreparedStatementByQuery returns a prepared statement by the query.
func (mgr *stmtManager) LookupPreparedStatementByQuery(query string) (PreparedStatement, error) {
	stmt, ok := mgr.stmtQueryMap[query]
	if !ok {
		return nil, newErrInvalidQuery(query)
	}

	return stmt, nil
}

// RemovePreparedStatement removes a prepared statement.
func (mgr *stmtManager) RemovePreparedStatement(stmt PreparedStatement) {
	delete(mgr.stmtQueryMap, stmt.Query())
	delete(mgr.stmtIDMap, stmt.StatementID())
}

// RemovePreparedStatementByID removes a prepared statement by the statement ID.
func (mgr *stmtManager) RemovePreparedStatementByID(stmtID StatementID) {
	stmt, err := mgr.LookupPreparedStatementByID(stmtID)
	if err != nil {
		return
	}

	mgr.RemovePreparedStatement(stmt)
}

// RemovePreparedStatementByQuery removes a prepared statement by the query.
func (mgr *stmtManager) RemovePreparedStatementByQuery(query string) {
	stmt, err := mgr.LookupPreparedStatementByQuery(query)
	if err != nil {
		return
	}

	mgr.RemovePreparedStatement(stmt)
}
