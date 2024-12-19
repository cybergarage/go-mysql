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
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
	"github.com/google/uuid"
)

// ConnOption represents a connection option.
type ConnOption = func(*conn)

// conn represents a connection of MySQL binary.
type conn struct {
	net.Conn
	isClosed      bool
	msgReader     *PacketReader
	db            string
	ts            time.Time
	uuid          uuid.UUID
	id            uint64
	tracerContext tracer.Context
	tlsConn       *tls.Conn
	capabilities  Capability
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(netConn net.Conn, opts ...ConnOption) Conn {
	conn := &conn{
		Conn:          netConn,
		isClosed:      false,
		msgReader:     NewPacketReaderWith(netConn),
		db:            "",
		ts:            time.Now(),
		uuid:          uuid.New(),
		id:            0,
		tracerContext: nil,
		tlsConn:       nil,
		capabilities:  0,
	}
	conn.SetOptions(opts...)
	return conn
}

// WithConnDatabase sets a database name.
func WithConnDatabase(name string) func(*conn) {
	return func(conn *conn) {
		conn.db = name
	}
}

// WithConnTracer sets a tracer context.
func WithConnTracer(t tracer.Context) func(*conn) {
	return func(conn *conn) {
		conn.tracerContext = t
	}
}

// WithConnTLSConnectionState sets a TLS connection state.
func WithConnTLSConn(s *tls.Conn) func(*conn) {
	return func(conn *conn) {
		conn.tlsConn = s
	}
}

// Withuint64 sets a connection ID.
func WithConnID(id uint64) func(*conn) {
	return func(conn *conn) {
		conn.id = id
	}
}

// WithConnUUID sets a UUID.
func WithConnUUID(id uuid.UUID) func(*conn) {
	return func(conn *conn) {
		conn.uuid = id
	}
}

// WithConnCapabilities sets capabilities.
func WithConnCapabilities(c Capability) func(*conn) {
	return func(conn *conn) {
		conn.capabilities = c
	}
}

// Close closes the connection.
func (conn *conn) Close() error {
	if conn.isClosed {
		return nil
	}
	if err := conn.Conn.Close(); err != nil {
		return err
	}
	conn.isClosed = true
	return nil
}

// SetOptions sets the connection options.
func (conn *conn) SetOptions(opts ...ConnOption) {
	for _, opt := range opts {
		opt(conn)
	}
}

// SetDatabase sets the database name.
func (conn *conn) SetDatabase(db string) {
	conn.db = db
}

// Database returns the database name.
func (conn *conn) Database() string {
	return conn.db
}

// Timestamp returns the creation time of the connection.
func (conn *conn) Timestamp() time.Time {
	return conn.ts
}

// UUID returns the UUID of the connection.
func (conn *conn) UUID() uuid.UUID {
	return conn.uuid
}

// UUID returns the connection ID of the connection.
func (conn *conn) ID() uint64 {
	return conn.id
}

// Context returns the context of the connection.
func (conn *conn) Context() context.Context {
	return context.Background()
}

// SetSpanContext sets the tracer span context of the connection.
func (conn *conn) SetSpanContext(ctx tracer.Context) {
	conn.tracerContext = ctx
}

// SpanContext returns the tracer span context of the connection.
func (conn *conn) SpanContext() tracer.Context {
	return conn.tracerContext
}

// Span returns the current top tracer span on the tracer span stack.
func (conn *conn) Span() tracer.Span {
	return conn.tracerContext.Span()
}

// StartSpan starts a new child tracer span and pushes it onto the tracer span stack.
func (conn *conn) StartSpan(name string) bool {
	return conn.tracerContext.StartSpan(name)
}

// FinishSpan ends the current top tracer span and pops it from the tracer span stack.
func (conn *conn) FinishSpan() bool {
	return conn.tracerContext.FinishSpan()
}

// IsTLSConnection return true if the connection is enabled TLS.
func (conn *conn) IsTLSConnection() bool {
	return conn.tlsConn != nil
}

// TLSConn returns the TLS connection.
func (conn *conn) TLSConn() *tls.Conn {
	return conn.tlsConn
}

// SetCapabilities sets the capabilities.
func (conn *conn) SetCapabilities(c Capability) {
	conn.capabilities = c
}

// Capabilities returns the capabilities.
func (conn *conn) Capabilities() Capability {
	return conn.capabilities
}

// PacketReader returns a packet reader.
func (conn *conn) PacketReader() *PacketReader {
	return conn.msgReader
}

// ResponsePacket sends a response.
func (conn *conn) ResponsePacket(resMsg Response, opts ...ResponseOption) error {
	if resMsg == nil {
		return nil
	}
	for _, opt := range opts {
		opt(resMsg)
	}
	resBytes, err := resMsg.Bytes()
	if err != nil {
		return err
	}
	if _, err := conn.Conn.Write(resBytes); err != nil {
		return err
	}
	return nil
}

// ResponsePackets sends response packets.
func (conn *conn) ResponsePackets(resMsgs []Response, opts ...ResponseOption) error {
	if len(resMsgs) == 0 {
		return nil
	}
	for _, resMsg := range resMsgs {
		err := conn.ResponsePacket(resMsg, opts...)
		if err != nil {
			return err
		}
	}
	return nil
}

// ResponseOK sends an OK response.
func (conn *conn) ResponseOK(opts ...OKOption) error {
	pkt, err := NewOK(opts...)
	if err != nil {
		return err
	}
	return conn.ResponsePacket(pkt)
}

// ResponseError sends an error response.
func (conn *conn) ResponseError(err error, opts ...ERROption) error {
	pkt, err := NewERRFromError(err, opts...)
	if err != nil {
		return err
	}
	return conn.ResponsePacket(pkt)
}
