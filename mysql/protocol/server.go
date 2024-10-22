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

package protocol

import (
	"crypto/tls"
	"errors"
	"net"
	"strconv"

	mysqlnet "github.com/cybergarage/go-mysql/mysql/net"
	"github.com/cybergarage/go-tracing/tracer"
)

// Server represents a MySQL protocol server.
type Server struct {
	*Config
	*mysqlnet.ConnManager
	CommandHandler
	tracer.Tracer
	tcpListener net.Listener
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Config:         NewDefaultConfig(),
		ConnManager:    mysqlnet.NewConnManager(),
		CommandHandler: nil,
		Tracer:         tracer.NullTracer,
		tcpListener:    nil,
	}
	return server
}

// SetTracer sets a tracing tracer.
func (server *Server) SetTracer(t tracer.Tracer) {
	server.Tracer = t
}

// SetCommandHandler sets a command handler.
func (server *Server) SetCommandHandler(h CommandHandler) {
	server.CommandHandler = h
}

// Start starts the server.
func (server *Server) Start() error {
	err := server.ConnManager.Start()
	if err != nil {
		return err
	}

	err = server.open()
	if err != nil {
		return err
	}

	go server.serve()

	return nil
}

// Stop stops the server.
func (server *Server) Stop() error {
	if err := server.ConnManager.Stop(); err != nil {
		return err
	}

	err := server.close()
	if err != nil {
		return err
	}

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

// open opens a listen socket.
func (server *Server) open() error {
	var err error
	addr := net.JoinHostPort(server.addr, strconv.Itoa(server.port))
	server.tcpListener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return nil
}

// close closes a listening socket.
func (server *Server) close() error {
	if server.tcpListener != nil {
		err := server.tcpListener.Close()
		if err != nil {
			return err
		}
	}

	server.tcpListener = nil

	return nil
}

// serve handles client requests.
func (server *Server) serve() error {
	defer server.close()

	l := server.tcpListener
	for {
		if l == nil {
			break
		}
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		go server.receive(conn)
	}

	return nil
}

// receive handles client packets.
func (server *Server) receive(netConn net.Conn) error { //nolint:gocyclo,maintidx
	defer func() {
		netConn.Close()
	}()

	// MySQL: Connection Lifecycle
	// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_lifecycle.html

	conn := NewConnWith(netConn)
	defer func() {
		conn.Close()
	}()

	reader := conn.PacketReader()

	// MySQL: Connection Phase
	// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase.html

	// Initial Handshake Packet

	handshakeMsg, err := NewHandshake(
		WithHandshakeServerVersion(server.ServerVersion()))
	if err != nil {
		return err
	}

	err = conn.ResponsePacket(handshakeMsg)
	if err != nil {
		return err
	}

	capFlags, err := reader.PeekCapabilityFlags()
	if err != nil {
		return err
	}

	if capFlags.IsEnabled(ClientSSL) {
		// SSL Connection Request Packet
		_, err := NewSSLRequestFromReader(reader)
		if err != nil {
			return err
		}
		// SSL exchange
		tlsConfig, err := server.Config.TLSConfig()
		if err != nil {
			conn.ResponseError(err)
			return err
		}
		tlsConn := tls.Server(conn, tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			conn.ResponseError(err)
			return err
		}
		tlsConnState := tlsConn.ConnectionState()
		conn = NewConnWith(tlsConn, WithTLSConnectionState(&tlsConnState))
	}

	// Handshake Response Packet
	handshakeRes, err := NewHandshakeResponseFromReader(reader)
	if err != nil {
		return err
	}
	conn.SetCapabilities(handshakeRes.CapabilityFlags())

	// MySQL: Command Phase
	// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_command_phase.html

	server.AddConn(conn)
	defer func() {
		server.RemoveConn(conn)
	}()

	for {
		if server.CommandHandler == nil {
			return errors.New("no command handler")
		}

		var err error
		var cmd Command

		cmd, err = NewCommandFromReader(conn)
		if err != nil {
			return err
		}
		cmd.SetCapabilities(conn.Capabilities())

		cmdType := cmd.Type()

		loopSpan := server.Tracer.StartSpan(server.productName)
		conn.SetSpanContext(loopSpan)
		conn.StartSpan(cmdType.String())

		finishSpans := func() {
			conn.FinishSpan()
			loopSpan.Span().Finish()
		}

		var res Packet
		switch cmdType {
		case ComQuit:
			ok, err := NewOK()
			if err == nil {
				err = conn.ResponsePacket(ok)
			}
			finishSpans()
			return err
		case ComQuery:
			var q *Query
			q, err = NewQueryFromCommand(cmd)
			if err == nil {
				res, err = server.CommandHandler.HandleQuery(conn, q)
			}
		default:
			err = cmd.SkipPayload()
		}

		conn.FinishSpan()

		conn.StartSpan("response")

		if err == nil {
			if res != nil {
				err = conn.ResponsePacket(res)
			}
		} else {
			err = conn.ResponseError(err)
		}

		conn.FinishSpan()

		loopSpan.Span().Finish()

		if err != nil {
			return err
		}
	}
}
