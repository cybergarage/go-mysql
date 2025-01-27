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

package protocol

// MySQL: COM_STMT_PREPARE
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_prepare.html

// StmtPrepare represents a COM_STMT_PREPARE packet.
type StmtPrepare struct {
	Command
	query string
}

func newStmtPrepareWithCommand(cmd Command, opts ...StmtPrepareOption) *StmtPrepare {
	q := &StmtPrepare{
		Command: cmd,
		query:   "",
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// StmtPrepareOption represents a MySQL StmtPrepare option.
type StmtPrepareOption func(*StmtPrepare)

// WithStmtPrepareQuery sets the query string.
func WithStmtPrepareQuery(query string) StmtPrepareOption {
	return func(q *StmtPrepare) {
		q.query = query
	}
}

// Query returns the query string.
func (stmt *StmtPrepare) Query() string {
	return stmt.query
}
