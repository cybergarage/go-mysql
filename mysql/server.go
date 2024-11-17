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

package mysql

import (
	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-tracing/tracer"
)

// SQLExecutor represents a SQL executor.
type SQLExecutor = query.SQLExecutor

// QueryExecutor represents a user query message executor.
type QueryExecutor interface {
	TCOExecutor
	DDOExecutor
	DMOExecutor
}

// ErrorHandler represents a user error handler.
type ErrorHandler interface {
	ParserError(Conn, string, error) (Response, error)
}

// ProtocolExecutor represents a frontend message executor.
type ProtocolExecutor interface {
	QueryExecutor
	ErrorHandler
}

// Server represents a MySQL-compatible server interface.
type Server interface {
	ServerConfig
	tracer.Tracer
	// SetTracer sets a tracing tracer.
	SetTracer(tracer.Tracer)
	// SetSQLExecutor sets an SQL executor to the server.
	SetSQLExecutor(executor SQLExecutor)
	// SetQueryExecutor sets a user query executor.
	SetQueryExecutor(QueryExecutor)
	// SetErrorHandler sets a user error handler.
	SetErrorHandler(ErrorHandler)
	// Start starts the server.
	Start() error
	// Stop stops the server.
	Stop() error
	// Restart restarts the server.
	Restart() error
}
