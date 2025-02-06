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

type preparedParameter struct {
	name string
	typ  stmt.FieldType
}

func (p *preparedParameter) Name() string {
	return p.name
}

func (p *preparedParameter) Type() stmt.FieldType {
	return p.typ
}

type preparedStmt struct {
	*StmtPrepare
	*StmtPrepareResponse
	params []stmt.Parameter
}

// NewPreparedStatmentWith creates a new prepared statement with the packet.
func NewPreparedStatmentWith(prePkt *StmtPrepare, resPkt *StmtPrepareResponse) stmt.PreparedStatement {
	preStmt := &preparedStmt{
		StmtPrepare:         prePkt,
		StmtPrepareResponse: resPkt,
	}

	resParams := resPkt.Params()
	preStmt.params = make([]stmt.Parameter, len(resParams))
	for n, resParam := range resParams {
		preStmt.params[n] = &preparedParameter{
			name: resParam.Name(),
			typ:  stmt.FieldType(resParam.ColType()),
		}
	}

	return preStmt
}

func (p *preparedStmt) Parameters() []stmt.Parameter {
	return p.params
}
