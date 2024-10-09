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
	"time"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-tracing/tracer"
	vitessmy "vitess.io/vitess/go/mysql"
)

// VitessListener is the MySQL server protocol listener.
type VitessListener struct {
	*vitessmy.Listener
}

// A VitessAuthHandler is an interface used for the user authentication.
type VitessAuthHandler interface {
	vitessmy.AuthServer
}

// A VitessQueryHandler is an interface used for the request queries.
type VitessQueryHandler interface {
	vitessmy.Handler
}

// NewVitessListener creates a new listener.
func NewVitessListener(protocol, address string, authServer VitessAuthHandler, handler VitessQueryHandler, connReadTimeout time.Duration, connWriteTimeout time.Duration, proxyProtocol bool) (*VitessListener, error) {
	l, err := vitessmy.NewListener(protocol, address, authServer, handler, connReadTimeout, connWriteTimeout, proxyProtocol)
	if err != nil {
		return nil, err
	}
	return &VitessListener{Listener: l}, nil
}

// VitessServer represents a MySQL-compatible server.
type VitessServer struct {
	tracer.Tracer
	Config
	ConnManager
	VitessAuthHandler
	VitessQueryHandler
	queryExecutor QueryExecutor
	listener      *VitessListener
}

// NewServer returns a new server instance.
func NewServer() *VitessServer {
	server := &VitessServer{
		Tracer:             tracer.NullTracer,
		Config:             NewDefaultConfig(),
		ConnManager:        NewConnManager(),
		VitessAuthHandler:  NewDefaultAuthHandler(),
		VitessQueryHandler: nil,
		queryExecutor:      nil,
		listener:           nil,
	}
	return server
}

// SetTracer sets a tracing tracer.
func (server *VitessServer) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// SetAuthHandler sets a user authentication handler.
func (server *VitessServer) SetAuthHandler(h VitessAuthHandler) {
	server.VitessAuthHandler = h
}

// SetQueryExecutor sets a query executor.
func (server *VitessServer) SetQueryExecutor(e QueryExecutor) {
	server.queryExecutor = e
}

// Start starts the server.
func (server *VitessServer) Start() error {
	hostPort := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	l, err := NewVitessListener("tcp", hostPort, server, server, 0, 0, false)
	if err != nil {
		return err
	}
	server.listener = l

	go server.listener.Accept()

	log.Infof("%s/%s (%s) started", PackageName, Version, hostPort)

	return nil
}

// Stop stops the server.
func (server *VitessServer) Stop() error {
	if server.listener != nil {
		server.listener.Close()
		server.listener = nil
	}

	hostPort := net.JoinHostPort(server.Address(), strconv.Itoa(server.Port()))
	log.Infof("%s/%s (%s) terminated", PackageName, Version, hostPort)

	return nil
}

// Restart restarts the server.
func (server *VitessServer) Restart() error {
	err := server.Stop()
	if err != nil {
		return err
	}

	return server.Start()
}
