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

// StmtExecute represents a COM_STMT_EXECUTE packet.
type StmtExecute struct {
	Command
	stmdID StatementID
}

func newStmtExecuteWithCommand(cmd Command, opts ...StmtExecuteOption) *StmtExecute {
	q := &StmtExecute{
		Command: cmd,
		stmdID:  0,
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// StmtExecuteOption represents a MySQL StmtExecute option.
type StmtExecuteOption func(*StmtExecute)

// WithStmtExecuteStatementID sets the statement ID.
func WithStmtExecuteStatementID(stmdID StatementID) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.stmdID = stmdID
	}
}

// NewStmtExecuteFromReader reads a COM_STMT_EXECUTE packet.
func NewStmtExecuteFromReader(reader io.Reader, opts ...StmtExecuteOption) (*StmtExecute, error) {
	var err error

	cmd, err := NewCommandFromReader(reader)
	if err != nil {
		return nil, err
	}

	if err = cmd.IsType(ComStmtExecute); err != nil {
		return nil, err
	}

	return NewStmtExecuteFromCommand(cmd, opts...)
}

// NewStmtExecuteFromCommand creates a new StmtExecute from a Command.
func NewStmtExecuteFromCommand(cmd Command, opts ...StmtExecuteOption) (*StmtExecute, error) {
	var err error

	pkt := newStmtExecuteWithCommand(cmd, opts...)

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
func (pkt *StmtExecute) StatementID() StatementID {
	return pkt.stmdID
}

// Bytes returns the packet bytes.
func (pkt *StmtExecute) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if err := w.WriteCommandType(pkt); err != nil {
		return nil, err
	}

	if err := w.WriteInt4(uint32(pkt.stmdID)); err != nil {
		return nil, err
	}

	return pkt.Command.Bytes()
}
