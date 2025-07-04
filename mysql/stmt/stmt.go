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
	"github.com/cybergarage/go-sqlparser/sql/stmt"
)

// NewStatementsFromPreparedStatement creates a new statement from the prepared statement and parameters.
func NewStatementsFromPreparedStatement(prepStmt PreparedStatement, params []Parameter) ([]Statement, error) {
	if len(params) != len(prepStmt.Parameters()) {
		return nil, newErrInvalidParameters()
	}
	bindParams := make(stmt.BindParams, len(params))
	for n, param := range params {
		v, err := param.Value()
		if err != nil {
			return nil, err
		}
		bindParams[n] = stmt.NewBindParam(v)
	}
	bindStmt := stmt.NewBindStatement(
		stmt.WithBindStatementQuery(prepStmt.Query()),
		stmt.WithBindStatementParams(bindParams),
	)
	return bindStmt.Statements()
}
