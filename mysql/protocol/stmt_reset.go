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

// StmtReset represents a COM_STMT_RESET packet.
type StmtReset struct {
	Command
	stmdID StatementID
}

func newStmtResetWithCommand(cmd Command, opts ...StmtResetOption) *StmtReset {
	q := &StmtReset{
		Command: cmd,
		stmdID:  0,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// StmtResetOption represents a MySQL StmtReset option.
type StmtResetOption func(*StmtReset)

// WithStmtResetStatementID sets the statement ID.
func WithStmtResetStatementID(stmdID StatementID) StmtResetOption {
	return func(q *StmtReset) {
		q.stmdID = stmdID
	}
}

// NewStmtResetFromReader reads a COM_STMT_RESET packet.
func NewStmtResetFromReader(reader io.Reader, opts ...StmtResetOption) (*StmtReset, error) {
	var err error

	cmd, err := NewCommandFromReader(reader)
	if err != nil {
		return nil, err
	}

	if err = cmd.IsType(ComStmtReset); err != nil {
		return nil, err
	}

	return NewStmtResetFromCommand(cmd, opts...)
}

// NewStmtResetFromCommand creates a new StmtReset from a Command.
func NewStmtResetFromCommand(cmd Command, opts ...StmtResetOption) (*StmtReset, error) {
	var err error

	pkt := newStmtResetWithCommand(cmd, opts...)

	payload := cmd.Payload()
	reader := NewPacketReaderWith(bytes.NewBuffer(payload[1:]))

	v, err := reader.ReadInt4()
	if err != nil {
		return nil, err
	}
	pkt.stmdID = StatementID(v)

	return pkt, nil
}

// StatementID returns the statement ID.
func (pkt *StmtReset) StatementID() StatementID {
	return pkt.stmdID
}

// Bytes returns the packet bytes.
func (pkt *StmtReset) Bytes() ([]byte, error) {
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
