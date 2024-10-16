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

package v2

import (
	"github.com/cybergarage/go-mysql/mysql/plugins"
	"github.com/cybergarage/go-mysql/mysql/protocol"
)

// Server represents a base executor server.
type Server struct {
	*protocol.Server
	executor plugins.QueryExecutor
}

// NewServer returns a base executor server instance.
func NewServer() *Server {
	s := &Server{
		Server:   protocol.NewServer(),
		executor: nil,
	}
	return s
}

// SetExecutor sets an executor to the server.
func (server *Server) SetExecutor(executor plugins.QueryExecutor) {
	server.executor = executor
}

// Executor returns the executor of the server.
func (server *Server) Executor() plugins.QueryExecutor {
	return server.executor
}

// HandleQuery handles a query.
func (server *Server) HandleQuery(q *protocol.Query) (protocol.Packet, error) {
	return nil, nil
}
