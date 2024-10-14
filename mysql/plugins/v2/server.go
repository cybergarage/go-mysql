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
)

// Server represents a base executor server.
type Server struct {
	executor plugins.Executor
}

// NewServer returns a base executor server instance.
func NewServer() *Server {
	s := &Server{
		executor: nil,
	}
	return s
}

// SetExecutor sets an executor to the server.
func (server *Server) SetExecutor(executor plugins.Executor) {
	server.executor = executor
}

// Executor returns the executor of the server.
func (server *Server) Executor() plugins.Executor {
	return server.executor
}
