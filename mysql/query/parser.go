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

package query

import (
	vitesssp "vitess.io/vitess/go/vt/sqlparser"
)

// Parser represents a SQL parser.
type Parser struct {
}

// NewParser returns a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// Parse parses a query string and returns the statement result.
func (parser *Parser) Parse(query string) (Statement, error) {
	stmt, err := vitesssp.Parse(query)
	return stmt, err
}
