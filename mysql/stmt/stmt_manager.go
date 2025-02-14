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
type StatementManager interface {
	// NextPreparedStatementID returns the next prepared statement ID.
	NextPreparedStatementID() (StatementID, error)
	// RegisterPreparedStatement adds a prepared statement to the manager.
	RegisterPreparedStatement(stmt PreparedStatement) error
	// LookupPreparedStatementByID returns a prepared statement by the statement ID.
	LookupPreparedStatementByID(stmtID StatementID) (PreparedStatement, error)
	// RemovePreparedStatement removes a prepared statement by the statement ID.
	RemovePreparedStatement(stmtID StatementID)
}
