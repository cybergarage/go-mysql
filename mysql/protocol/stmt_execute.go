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

	"github.com/cybergarage/go-mysql/mysql/stmt"
)

// MySQL: COM_STMT_EXECUTE
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_execute.html
// COM_STMT_EXECUTE - MariaDB Knowledge Base
// https://mariadb.com/kb/en/com_stmt_execute/

// StatementBindSendType represents a statement bind send type.
type StatementBindSendType uint8

const (
	// StatementBindSendTypeToServer represents a statement bind send type to server.
	StatementBindSendTypeToServer StatementBindSendType = 0x01
)

// IsToServer returns true if the statement bind send type is to server.
func (t StatementBindSendType) IsToServer() bool {
	return t == StatementBindSendTypeToServer
}

// StmtExecute represents a COM_STMT_EXECUTE packet.
type StmtExecute struct {
	Command
	stmdID       StatementID
	cursorType   CursorType
	iterCnt      uint32
	numParams    uint16
	nullBitmap   *NullBitmap
	bindSendType StatementBindSendType
	paramNames   []string
	paramValues  [][]byte
	paramTypes   []FieldType
	stmtMgr      stmt.StatementManager
	params       []stmt.Parameter
}

func newStmtExecuteWithCommand(cmd Command, opts ...StmtExecuteOption) *StmtExecute {
	q := &StmtExecute{
		Command:      cmd,
		stmdID:       0,
		cursorType:   CursorTypeNoCursor,
		iterCnt:      1,
		numParams:    0,
		nullBitmap:   nil,
		bindSendType: 0,
		paramNames:   []string{},
		paramValues:  [][]byte{},
		paramTypes:   []FieldType{},
		stmtMgr:      nil,
		params:       []stmt.Parameter{},
	}
	for _, opt := range opts {
		opt(q)
	}
	return q
}

// StmtExecuteOption represents a MySQL StmtExecute option.
type StmtExecuteOption func(*StmtExecute)

// WithStmtExecuteStatementCapability sets the statement capability.
func WithStmtExecuteStatementCapability(c Capability) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.SetCapability(c)
	}
}

// WithStmtExecuteStatementID sets the statement ID.
func WithStmtExecuteStatementID(stmdID StatementID) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.stmdID = stmdID
	}
}

// WithStmtExecuteCursorType sets the cursor type.
func WithStmtExecuteCursorType(cursorType CursorType) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.cursorType = cursorType
	}
}

// WithStmtExecuteIterationCount sets the iteration count.
func WithStmtExecuteIterationCount(iterCnt uint32) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.iterCnt = iterCnt
	}
}

// WithStmtExecuteNumParams sets the number of parameters.
func WithStmtExecuteNumParams(numParams uint16) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.numParams = numParams
	}
}

// WithStmtExecuteBindSendType sets the bind send type.
func WithStmtExecuteStatementManager(stmtMgr stmt.StatementManager) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.stmtMgr = stmtMgr
	}
}

// WithStmtExecuteBindSendType sets the bind send type.
func WithStmtExecuteBindSendType(bindSendType StatementBindSendType) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.bindSendType = bindSendType
	}
}

// WithStmtExecuteNullBitmap sets the null bitmap.
func WithStmtExecuteNullBitmap(nullBitmap *NullBitmap) StmtExecuteOption {
	return func(q *StmtExecute) {
		q.nullBitmap = nullBitmap
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
	pktReader := NewPacketReaderWith(bytes.NewBuffer(payload[1:]))

	iv4, err := pktReader.ReadInt4()
	if err != nil {
		return nil, err
	}
	pkt.stmdID = StatementID(iv4)

	iv1, err := pktReader.ReadInt1()
	if err != nil {
		return nil, err
	}
	pkt.cursorType = CursorType(iv1)

	pkt.iterCnt, err = pktReader.ReadInt4()
	if err != nil {
		return nil, err
	}

	if pkt.stmtMgr != nil {
		stmt, err := pkt.stmtMgr.PreparedStatement(pkt.stmdID)
		if err != nil {
			return nil, err
		}
		pkt.numParams = uint16(len(stmt.Parameters()))
	}

	if pkt.numParams == 0 {
		return pkt, nil
	}

	nullBitmapLen := int((pkt.numParams + 7) / 8)
	if 0 < nullBitmapLen {
		nullBitmapBytes := make([]byte, nullBitmapLen)
		if _, err := pktReader.ReadBytes(nullBitmapBytes); err != nil {
			return nil, err
		}
		pkt.nullBitmap = NewNullBitmap(
			WithNullBitmapNumFields(int(pkt.numParams)),
			WithNullBitmapBytes(nullBitmapBytes),
		)
	}

	iv1, err = pktReader.ReadInt1()
	if err != nil {
		return nil, err
	}
	pkt.bindSendType = StatementBindSendType(iv1)

	pkt.paramNames = make([]string, pkt.numParams)
	pkt.paramTypes = make([]FieldType, pkt.numParams)

	if pkt.bindSendType.IsToServer() {
		for n := 0; n < int(pkt.numParams); n++ {
			iv2, err := pktReader.ReadInt2()
			if err != nil {
				return nil, err
			}
			pkt.paramTypes[n] = FieldType(iv2)

			if pkt.Capability().IsEnabled(ClientQueryAttributes) {
				paramName, err := pktReader.ReadLengthEncodedString()
				if err != nil {
					return nil, err
				}
				pkt.paramNames[n] = paramName
			}
		}
	}

	pkt.paramValues = make([][]byte, pkt.numParams)
	for n := 0; n < int(pkt.numParams); n++ {
		if pkt.nullBitmap.IsNull(n) {
			continue
		}
		paramValue, err := pktReader.ReadLengthEncodedBytes()
		if err != nil {
			return nil, err
		}
		pkt.paramValues[n] = paramValue
	}

	// Create parameters

	pkt.params = make([]stmt.Parameter, pkt.numParams)
	for n := 0; n < int(pkt.numParams); n++ {
		paramOpts := []stmt.ParameterOption{}
		if n < len(pkt.paramNames) {
			paramOpts = append(paramOpts, stmt.WithParameterName(pkt.paramNames[n]))
		}
		if n < len(pkt.paramTypes) {
			paramOpts = append(paramOpts, stmt.WithParameterType(pkt.paramTypes[n]))
		}
		if n < len(pkt.paramValues) {
			paramOpts = append(paramOpts, stmt.WithParameterValue(pkt.paramValues[n]))
		}
		pkt.params[n] = stmt.NewParameter(paramOpts...)
	}

	return pkt, nil
}

// StatementID returns the statement ID.
func (pkt *StmtExecute) StatementID() StatementID {
	return pkt.stmdID
}

// CursorType returns the cursor type.
func (pkt *StmtExecute) CursorType() CursorType {
	return pkt.cursorType
}

// Parameters returns the parameters.
func (pkt *StmtExecute) Parameters() []stmt.Parameter {
	return pkt.params
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

	if err := w.WriteInt1(byte(pkt.cursorType)); err != nil {
		return nil, err
	}

	if err := w.WriteInt4(pkt.iterCnt); err != nil {
		return nil, err
	}

	if 0 < pkt.numParams {
		if _, err := w.WriteBytes(pkt.nullBitmap.Bytes()); err != nil {
			return nil, err
		}

		if err := w.WriteInt1(byte(pkt.bindSendType)); err != nil {
			return nil, err
		}

		if pkt.bindSendType.IsToServer() {
			for n := 0; n < int(pkt.numParams); n++ {
				if err := w.WriteInt2(uint16(pkt.paramTypes[n])); err != nil {
					return nil, err
				}
				if pkt.Capability().IsEnabled(ClientQueryAttributes) {
					if err := w.WriteLengthEncodedString(pkt.paramNames[n]); err != nil {
						return nil, err
					}
				}
			}
		}

		for n := 0; n < int(pkt.numParams); n++ {
			if pkt.nullBitmap.IsNull(n) {
				continue
			}
			if err := w.WriteLengthEncodedBytes(pkt.paramValues[n]); err != nil {
				return nil, err
			}
		}
	}

	pkt.SetPayload(w.Bytes())

	return pkt.Command.Bytes()
}
