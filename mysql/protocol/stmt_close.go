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

// MySQL: COM_STMT_CLOSE
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_close.html
// COM_STMT_CLOSE - MariaDB Knowledge Base
// https://mariadb.com/kb/en/3-binary-protocol-prepared-statements-com_stmt_close/

// StmtClose represents a COM_STMT_CLOSE packet.
type StmtClose struct {
	Command

	stmdID StatementID
}

func newStmtCloseWithCommand(cmd Command, opts ...StmtCloseOption) *StmtClose {
	q := &StmtClose{
		Command: cmd,
		stmdID:  0,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// StmtCloseOption represents a MySQL StmtClose option.
type StmtCloseOption func(*StmtClose)

// WithStmtCloseStatementID sets the statement ID.
func WithStmtCloseStatementID(stmdID StatementID) StmtCloseOption {
	return func(q *StmtClose) {
		q.stmdID = stmdID
	}
}

// NewStmtCloseFromReader reads a COM_STMT_CLOSE packet.
func NewStmtCloseFromReader(reader io.Reader, opts ...StmtCloseOption) (*StmtClose, error) {
	var err error

	cmd, err := NewCommandFromReader(reader)
	if err != nil {
		return nil, err
	}

	if err = cmd.IsType(ComStmtClose); err != nil {
		return nil, err
	}

	return NewStmtCloseFromCommand(cmd, opts...)
}

// NewStmtCloseFromCommand creates a new StmtClose from a Command.
func NewStmtCloseFromCommand(cmd Command, opts ...StmtCloseOption) (*StmtClose, error) {
	var err error

	pkt := newStmtCloseWithCommand(cmd, opts...)

	payload := cmd.Payload()
	reader := NewPacketReaderWithReader(bytes.NewBuffer(payload[1:]))

	v, err := reader.ReadInt4()
	if err != nil {
		return nil, err
	}
	pkt.stmdID = StatementID(v)

	return pkt, nil
}

// StatementID returns the statement ID.
func (pkt *StmtClose) StatementID() StatementID {
	return pkt.stmdID
}

// Bytes returns the packet bytes.
func (pkt *StmtClose) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteCommandType(pkt); err != nil {
		return nil, err
	}

	if err := w.WriteInt4(uint32(pkt.stmdID)); err != nil {
		return nil, err
	}

	pkt.SetPayload(w.Bytes())

	return pkt.Command.Bytes()
}
