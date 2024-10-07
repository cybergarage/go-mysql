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
	"net"
	"time"

	"github.com/cybergarage/go-tracing/tracer"
	"github.com/google/uuid"
)

// ConnOption represents a connection option.
type ConnOption = func(*Conn)

// Conn represents a connection of MySQL binary.
type Conn struct {
	net.Conn
	isClosed  bool
	msgReader *PacketReader
	db        string
	ts        time.Time
	uuid      uuid.UUID
	tracer.Context
	tlsState     *tls.ConnectionState
	capabilities CapabilityFlag
}

// NewConnWith returns a connection with a raw connection.
func NewConnWith(netConn net.Conn, opts ...ConnOption) *Conn {
	conn := &Conn{
		Conn:         netConn,
		isClosed:     false,
		msgReader:    NewPacketReaderWith(netConn),
		db:           "",
		ts:           time.Now(),
		uuid:         uuid.New(),
		Context:      nil,
		tlsState:     nil,
		capabilities: 0,
	}
	for _, opt := range opts {
		opt(conn)
	}
	return conn
}

// WithConnDatabase sets a database name.
func WithConnDatabase(name string) func(*Conn) {
	return func(conn *Conn) {
		conn.db = name
	}
}

// WithConnTracer sets a tracer context.
func WithConnTracer(t tracer.Context) func(*Conn) {
	return func(conn *Conn) {
		conn.Context = t
	}
}

// WithTLSConnectionState sets a TLS connection state.
func WithTLSConnectionState(s *tls.ConnectionState) func(*Conn) {
	return func(conn *Conn) {
		conn.tlsState = s
	}
}

// WithCapabilities sets capabilities.
func WithCapabilities(c CapabilityFlag) func(*Conn) {
	return func(conn *Conn) {
		conn.capabilities = c
	}
}

// Close closes the connection.
func (conn *Conn) Close() error {
	if conn.isClosed {
		return nil
	}
	if err := conn.Conn.Close(); err != nil {
		return err
	}
	conn.isClosed = true
	return nil
}

// SetDatabase sets the database name.
func (conn *Conn) SetDatabase(db string) {
	conn.db = db
}

// Database returns the database name.
func (conn *Conn) Database() string {
	return conn.db
}

// Timestamp returns the creation time of the connection.
func (conn *Conn) Timestamp() time.Time {
	return conn.ts
}

// UUID returns the UUID of the connection.
func (conn *Conn) UUID() uuid.UUID {
	return conn.uuid
}

// SetSpanContext sets the tracer span context of the connection.
func (conn *Conn) SetSpanContext(ctx tracer.Context) {
	conn.Context = ctx
}

// SpanContext returns the tracer span context of the connection.
func (conn *Conn) SpanContext() tracer.Context {
	return conn.Context
}

// IsTLSConnection return true if the connection is enabled TLS.
func (conn *Conn) IsTLSConnection() bool {
	return conn.tlsState != nil
}

// TLSConnectionState returns the TLS connection state.
func (conn *Conn) TLSConnectionState() (*tls.ConnectionState, bool) {
	return conn.tlsState, conn.tlsState != nil
}

// Capabilities returns the capabilities.
func (conn *Conn) Capabilities() CapabilityFlag {
	return conn.capabilities
}

// PacketReader returns a packet reader.
func (conn *Conn) PacketReader() *PacketReader {
	return conn.msgReader
}

// ResponsePacket sends a response.
func (conn *Conn) ResponsePacket(resMsg Packet) error {
	if resMsg == nil {
		return nil
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
func (conn *Conn) ResponsePackets(resMsgs []Packet) error {
	if len(resMsgs) == 0 {
		return nil
	}
	for _, resMsg := range resMsgs {
		err := conn.ResponsePacket(resMsg)
		if err != nil {
			return err
		}
	}
	return nil
}

// ResponseError sends an error response.
func (conn *Conn) ResponseError(err error) error {
	pkt, err := NewERRFromError(err)
	if err != nil {
		return err
	}
	return conn.ResponsePacket(pkt)
}
