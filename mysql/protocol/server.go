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
	"bytes"
	"crypto/tls"
	"errors"
	"net"
	"strconv"

	"github.com/cybergarage/go-mysql/mysql/auth"
	mysqlnet "github.com/cybergarage/go-mysql/mysql/net"
	"github.com/cybergarage/go-tracing/tracer"
)

// Server represents a MySQL protocol server.
type Server struct {
	*Config
	*auth.Manager
	*mysqlnet.ConnManager
	tracer.Tracer
	lastConnID *Counter
	CommandHandler
	tcpListener net.Listener
}

// NewServer returns a new server instance.
func NewServer() *Server {
	server := &Server{
		Config:         NewDefaultConfig(),
		Manager:        auth.NewManager(),
		ConnManager:    mysqlnet.NewConnManager(),
		Tracer:         tracer.NullTracer,
		lastConnID:     NewCounter(),
		CommandHandler: nil,
		tcpListener:    nil,
	}
	server.SetCapability(DefaultServerCapabilities)
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

// GenerateHandshakeForConn retrun a handshake message for the specified connection.
func (server *Server) GenerateHandshakeForConn(conn mysqlnet.Conn) (*Handshake, error) {
	salt, err := auth.NewSalt(DefaultAuthPluginDataPartLen)
	if err != nil {
		return nil, err
	}
	return NewHandshake(
		WithHandshakeCharacterSet(CharSetUTF8),
		WithHandshakeCapability(server.Capability()),
		WithHandshakeServerVersion(server.ServerVersion()),
		WithHandshakeConnectionID(uint32(conn.ID())),
		WithHandshakeAuthPluginData(salt),
		WithHandshakeAuthPluginName(server.AuthPluginName()),
	), nil
}

// receive handles client packets.
func (server *Server) receive(netConn net.Conn) error { //nolint:gocyclo,maintidx
	// MySQL: Connection Lifecycle
	// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_lifecycle.html

	constructConnection := func(netConn net.Conn) (Conn, error) {
		server.lastConnID.Lock()
		defer server.lastConnID.Unlock()

		lastConnID := server.lastConnID.Count()
		nextConnID := server.lastConnID.Inc()

		conn := NewConnWith(netConn,
			WithConnID(uint64(nextConnID)),
		)

		for {
			if lastConnID == nextConnID {
				netConn.Close()
				return nil, ErrTooManyConnections
			}

			if err := server.AddConn(conn); err == nil {
				return conn, nil
			}

			nextConnID = server.lastConnID.Inc()

			conn = NewConnWith(netConn,
				WithConnID(uint64(nextConnID)),
			)
		}
	}

	conn, err := constructConnection(netConn)
	if err != nil {
		netConn.Close()
		return err
	}

	// MySQL: Connection Phase
	// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase.html

	reader := conn.PacketReader()

	// Initial Handshake Packet

	handshakeMsg, err := server.GenerateHandshakeForConn(conn)
	if err != nil {
		return err
	}

	err = conn.ResponsePacket(handshakeMsg)
	if err != nil {
		return err
	}

	// Initial Handshake Response Packet

	firstPkt, err := NewPacketWithReader(reader)
	if err != nil {
		return err
	}

	firstPktBytes, err := firstPkt.Bytes()
	if err != nil {
		return err
	}

	firstPktReader := bytes.NewBuffer(firstPktBytes)

	isSSLRequestPkt := func(pkt Packet) bool {
		if pkt.PayloadLength() < 4 {
			return false
		}
		capFlags := NewCapabilityFlagFromBytes(pkt.Payload()[0:4])
		return capFlags.IsEnabled(ClientSSL)
	}

	if isSSLRequestPkt(firstPkt) {
		// SSL Connection Request Packet
		_, err := NewSSLRequestFromReader(firstPktReader)
		if err != nil {
			return err
		}
		// SSL exchange
		tlsConfig, err := server.Config.TLSConfig()
		if err != nil {
			conn.ResponseError(err)
			return errors.Join(err, conn.Close())
		}
		tlsConn := tls.Server(conn, tlsConfig)
		if err := tlsConn.Handshake(); err != nil {
			conn.ResponseError(err)
			return errors.Join(err, conn.Close())
		}
		tlsConnState := tlsConn.ConnectionState()
		// Update TLS connection to the connection manager
		newConn := NewConnWith(
			tlsConn,
			WithConnID(conn.ID()),
			WithConnUUID(conn.UUID()),
			WithConnTLSConnectionState(&tlsConnState))
		if err := server.UpdateConn(conn, newConn); err != nil {
			conn.ResponseError(err)
			return errors.Join(err, conn.Close())
		}
		conn = newConn
	}

	defer func() {
		conn.Close()
		server.RemoveConn(conn)
	}()

	// Handshake Response Packet

	handshakeRes, err := NewHandshakeResponseFromReader(firstPktReader)
	if err != nil {
		return err
	}

	conn.SetCapabilities(handshakeRes.CapabilityFlags())

	authQuery := auth.NewQuery(
		auth.WithQueryUsername(handshakeRes.Username()),
		auth.WithQueryAuthResponse(handshakeRes.AuthResponse()),
	)

	ok := server.Authenticate(conn, authQuery)
	if !ok {
		conn.ResponseError(
			auth.ErrAccessDenied,
			WithERRSecuenceID(handshakeRes.SequenceID().Next()),
		)
		return errors.Join(err, conn.Close())
	}

	err = conn.ResponseOK(
		WithOKSecuenceID(handshakeRes.SequenceID().Next()),
	)
	if err != nil {
		return err
	}

	// MySQL: Command Phase
	// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_command_phase.html

	for {
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

		var res Response
		switch cmdType {
		case ComPing:
			res, err = NewOK()
		case ComQuery:
			if server.CommandHandler != nil {
				var q *Query
				q, err = NewQueryFromCommand(cmd)
				if err == nil {
					res, err = server.CommandHandler.HandleQuery(conn, q)
				}
			} else {
				err = newErrNotSupportedCommandType(cmdType)
			}
		case ComQuit:
			ok, err := NewOK()
			if err == nil {
				err = conn.ResponsePacket(ok)
			}
			finishSpans()
			return err
		default:
			err = cmd.SkipPayload()
			if err == nil {
				err = newErrNotSupportedCommandType(cmdType)
			}
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
