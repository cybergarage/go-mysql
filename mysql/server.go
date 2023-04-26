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

package mysql

import (
	"net"
	"strconv"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-tracing/tracer"
)

// Server represents a MySQL-compatible server.
type Server struct {
	tracer.Tracer
	AuthHandler
	QueryHandler
	*Config
	ConnMap
	queryExecutor QueryExecutor
	listener      *Listener
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Config:        NewDefaultConfig(),
		Tracer:        tracer.NullTracer,
		AuthHandler:   NewDefaultAuthHandler(),
		QueryHandler:  nil,
		queryExecutor: nil,
		ConnMap:       NewConnMap(),
	}
	return server
}

// SetTracer sets a tracing tracer.
func (server *Server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// SetAuthHandler sets a user authentication handler.
func (server *Server) SetAuthHandler(h AuthHandler) {
	server.AuthHandler = h
}

// SetQueryExecutor sets a query executor.
func (server *Server) SetQueryExecutor(e QueryExecutor) {
	server.queryExecutor = e
}

// Start starts the server.
func (server *Server) Start() error {
	hostPort := net.JoinHostPort(server.Address, strconv.Itoa(server.Port))
	l, err := NewListener("tcp", hostPort, server, server, 0, 0, false)
	if err != nil {
		return err
	}
	server.listener = l

	go server.listener.Accept()

	log.Infof("%s/%s (%s) started", PackageName, Version, hostPort)

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	if server.listener != nil {
		server.listener.Close()
		server.listener = nil
	}

	hostPort := net.JoinHostPort(server.Address, strconv.Itoa(server.Port))
	log.Infof("%s/%s (%s) terminated", PackageName, Version, hostPort)

	return nil
}

// Restart restarts the server.
func (server *Server) Restart() error {
	err := server.Stop()
	if err != nil {
		return err
	}

	return server.Start()
}
