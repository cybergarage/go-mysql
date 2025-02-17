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
	"github.com/cybergarage/go-mysql/mysql/encoding/binary"
	"github.com/cybergarage/go-mysql/mysql/query"
)

// MySQL: COM_STMT_EXECUTE
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_execute.html
// MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// COM_STMT_EXECUTE - MariaDB Knowledge Base
// https://mariadb.com/kb/en/com_stmt_execute/
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/

// ParameterOption is the option of the parameter.
type ParameterOption func(*parameter)

type parameter struct {
	name string
	typ  FieldType
	v    []byte
}

// WithParameterName sets the name of the parameter.
func WithParameterName(name string) ParameterOption {
	return func(p *parameter) {
		p.name = name
	}
}

// WithParameterType sets the type of the parameter.
func WithParameterType(typ FieldType) ParameterOption {
	return func(p *parameter) {
		p.typ = typ
	}
}

// WithParameterBytes sets the value of the parameter.
func WithParameterBytes(v []byte) ParameterOption {
	return func(p *parameter) {
		p.v = v
	}
}

// NewParameter creates a new parameter with the options.
func NewParameter(opts ...ParameterOption) Parameter {
	param := &parameter{
		name: "",
		typ:  0,
		v:    nil,
	}
	for _, opt := range opts {
		opt(param)
	}
	return param
}

// Name returns the name of the parameter.
func (param *parameter) Name() string {
	return param.name
}

// Type returns the type of the parameter.
func (param *parameter) Type() FieldType {
	return param.typ
}

// Bytes returns the value of the parameter.
func (param *parameter) Bytes() []byte {
	return param.v
}

// Value returns the value of the parameter.
func (param *parameter) Value() (any, error) {
	switch param.typ {
	case query.MySQLTypeTiny:
		return binary.BytesToInt1(param.v)
	case query.MySQLTypeShort:
		return binary.BytesToInt2(param.v)
	case query.MySQLTypeLong:
		return binary.BytesToInt4(param.v)
	case query.MySQLTypeLonglong:
		return binary.BytesToInt8(param.v)
	case query.MySQLTypeFloat:
		return binary.BytesToFloat4(param.v)
	case query.MySQLTypeDouble:
		return binary.BytesToFloat8(param.v)
	case query.MySQLTypeNull:
		return nil, nil
	case query.MySQLTypeString, query.MySQLTypeVarString, query.MySQLTypeVarchar:
		return string(param.v), nil
	case query.MySQLTypeTinyBlob, query.MySQLTypeMediumBlob, query.MySQLTypeLongBlob, query.MySQLTypeBlob:
		return param.v, nil
	}

	return param.v, newErrNotSupportedFieldType(param.typ)
}
