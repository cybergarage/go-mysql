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

package vitess

import (
	"net"
	"strconv"
	"time"

	"github.com/cybergarage/go-logger/log"
	mysqlnet "github.com/cybergarage/go-mysql/mysql/net"
	"github.com/cybergarage/go-mysql/mysql/plugins"
	"github.com/cybergarage/go-tracing/tracer"
	vitessmy "vitess.io/vitess/go/mysql"
)

// Listener is the MySQL server protocol listener.
type Listener struct {
	*vitessmy.Listener
}

// A AuthHandler is an interface used for the user authentication.
type AuthHandler interface {
	vitessmy.AuthServer
}

// A QueryHandler is an interface used for the request queries.
type QueryHandler interface {
	vitessmy.Handler
}

// NewListener creates a new listener.
func NewListener(protocol, address string, authServer AuthHandler, handler QueryHandler, connReadTimeout time.Duration, connWriteTimeout time.Duration, proxyProtocol bool) (*Listener, error) {
	l, err := vitessmy.NewListener(protocol, address, authServer, handler, connReadTimeout, connWriteTimeout, proxyProtocol)
	if err != nil {
		return nil, err
	}
	return &Listener{Listener: l}, nil
}

// Server represents a MySQL-compatible server.
type Server struct {
	plugins.Executor
	tracer.Tracer
	plugins.Config
	*mysqlnet.ConnManager
	AuthHandler
	QueryHandler
	queryExecutor  QueryExecutor
	listener       *Listener
	productName    string
	productVersion string
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Tracer:         tracer.NullTracer,
		Config:         plugins.NewDefaultConfig(),
		ConnManager:    mysqlnet.NewConnManager(),
		AuthHandler:    NewDefaultAuthHandler(),
		QueryHandler:   nil,
		queryExecutor:  nil,
		listener:       nil,
		productName:    "vitess",
		productVersion: "x.x.x",
		Executor:       nil,
	}
	return server
}

// SetExecutor sets an executor to the server.
func (server *Server) SetExecutor(executor plugins.Executor) {
	server.Executor = executor
}

// SetProuctName sets a product name to the configuration.
func (server *Server) SetProductName(v string) {
	server.productName = v
}

// SetProductVersion sets a product version to the configuration.
func (server *Server) SetProductVersion(v string) {
	server.productVersion = v
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
	hostPort := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	l, err := NewListener("tcp", hostPort, server, server, 0, 0, false)
	if err != nil {
		return err
	}
	server.listener = l

	go server.listener.Accept()

	log.Infof("%s/%s (%s) started", server.productName, server.productVersion, hostPort)

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	if server.listener != nil {
		server.listener.Close()
		server.listener = nil
	}

	hostPort := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	log.Infof("%s/%s (%s) terminated", server.productName, server.productVersion, hostPort)

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
