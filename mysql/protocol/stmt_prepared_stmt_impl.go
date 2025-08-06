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

package protocol

import (
	"github.com/cybergarage/go-mysql/mysql/stmt"
)

type preparedStmt struct {
	*StmtPrepare
	*StmtPrepareResponse

	params []stmt.Parameter
}

// NewPreparedStatmentWith creates a new prepared statement with the packet.
func NewPreparedStatmentWith(prePkt *StmtPrepare, resPkt *StmtPrepareResponse) stmt.PreparedStatement {
	resParams := resPkt.Params()
	params := make([]stmt.Parameter, len(resParams))
	for n, resParam := range resParams {
		params[n] = stmt.NewParameter(
			stmt.WithParameterName(resParam.Name()),
			stmt.WithParameterType(stmt.FieldType(resParam.ColType())))
	}
	return &preparedStmt{
		StmtPrepare:         prePkt,
		StmtPrepareResponse: resPkt,
		params:              params,
	}
}

// Parameters returns the parameters of the prepared statement.
func (p *preparedStmt) Parameters() []stmt.Parameter {
	return p.params
}

// Bind binds the parameters to the prepared statement.
func (p *preparedStmt) Bind(params []stmt.Parameter) ([]stmt.Statement, error) {
	return stmt.NewStatementsFromPreparedStatement(p, params)
}

// PrepareBytes returns the prepared packet bytes.
func (p *preparedStmt) PrepareBytes() []byte {
	bytes, _ := p.StmtPrepare.Bytes()
	return bytes
}

// PrepareResponseBytes returns the prepared response packet bytes.
func (p *preparedStmt) PrepareResponseBytes() []byte {
	bytes, _ := p.StmtPrepareResponse.Bytes()
	return bytes
}
