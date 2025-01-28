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
	"io"
)

// MySQL: COM_STMT_PREPARE
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_prepare.html

// StmtPrepare represents a COM_STMT_PREPARE packet.
type StmtPrepare struct {
	Command
	query string
}

func newStmtPrepareWithCommand(cmd Command, opts ...StmtPrepareOption) *StmtPrepare {
	q := &StmtPrepare{
		Command: cmd,
		query:   "",
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// StmtPrepareOption represents a MySQL StmtPrepare option.
type StmtPrepareOption func(*StmtPrepare)

// WithStmtPrepareQuery sets the query string.
func WithStmtPrepareQuery(query string) StmtPrepareOption {
	return func(q *StmtPrepare) {
		q.query = query
	}
}

// NewStmtPrepareFromReader reads a COM_STMT_PREPARE packet.
func NewStmtPrepareFromReader(reader io.Reader, opts ...StmtPrepareOption) (*StmtPrepare, error) {
	var err error

	cmd, err := NewCommandFromReader(reader)
	if err != nil {
		return nil, err
	}

	if err = cmd.IsType(ComQuery); err != nil {
		return nil, err
	}

	return NewStmtPrepareFromCommand(cmd, opts...)
}

// NewStmtPrepareFromCommand creates a new StmtPrepare from a Command.
func NewStmtPrepareFromCommand(cmd Command, opts ...StmtPrepareOption) (*StmtPrepare, error) {
	var err error

	pkt := newStmtPrepareWithCommand(cmd, opts...)

	payload := cmd.Payload()
	reader := NewPacketReaderWith(bytes.NewBuffer(payload[1:]))

	pkt.query, err = reader.ReadEOFTerminatedString()
	if err != nil {
		return nil, err
	}

	return pkt, nil
}

// Query returns the query string.
func (stmt *StmtPrepare) Query() string {
	return stmt.query
}

// Bytes returns the packet bytes.
func (pkt *StmtPrepare) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteCommandType(pkt); err != nil {
		return nil, err
	}

	if _, err := w.WriteBytes([]byte(pkt.query)); err != nil {
		return nil, err
	}

	return pkt.Command.Bytes()
}
