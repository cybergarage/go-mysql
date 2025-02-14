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
	"github.com/cybergarage/go-mysql/mysql/query"
)

// Statement is the type of statement.
type Statement = query.Statement

// PreparedStatement is the interface of prepared statement.
type PreparedStatement interface {
	// StatementID returns the statement ID.
	StatementID() StatementID
	// DatabaseName returns the database name of the statement.
	DatabaseName() string
	// Statement returns the statement.
	Statement() Statement
	// TableNames returns the table names of the statement.
	TableNames() []string
	// Query returns the query of the statement.
	Query() string
	// Parameters returns the parameters of the statement.
	Parameters() []Parameter
	// Bind binds the parameters to the statement.
	Bind([]Parameter) (Statement, error)
}
