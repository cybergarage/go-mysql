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

// MySQL: COM_STMT_PREPARE Response
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_prepare.html
// COM_STMT_PREPARE response - MariaDB Knowledge Base
// https://mariadb.com/kb/en/com_stmt_prepare/
// 3 - Binary Protocol (Prepared Statements) - MariaDB Knowledge Base
// https://mariadb.com/kb/en/3-binary-protocol-prepared-statements/
// Server Response Packets (Binary Protocol) - MariaDB Knowledge Base
// https://mariadb.com/kb/en/server-response-packets-binary-protocol/

// StmtPrepareResponse represents a MySQL Prepare Response packet.
type StmtPrepareResponse struct {
	*packet
	status            Status
	stmtID            StatementID
	columns           []ColumnDef
	params            []ColumnDef
	warningCount      uint16
	resultSetMetadata ResultsetMetadata
}

func newStmtPrepareResponseWithPacket(pkt *packet, opts ...StmtPrepareResponseOption) *StmtPrepareResponse {
	prPkt := &StmtPrepareResponse{
		packet:            pkt,
		status:            Status(0),
		stmtID:            0,
		columns:           []ColumnDef{},
		params:            []ColumnDef{},
		warningCount:      0,
		resultSetMetadata: ResultsetMetadataNone,
	}

	prPkt.SetSequenceID(1)

	for _, opt := range opts {
		opt(prPkt)
	}
	return prPkt
}

// StmtPrepareResponseOption represents a StmtPrepareResponse option.
type StmtPrepareResponseOption func(*StmtPrepareResponse)

// WithStmtPrepareResponseCapability sets the capabilitys.
func WithStmtPrepareResponseCapability(c Capability) StmtPrepareResponseOption {
	return func(pkt *StmtPrepareResponse) {
		pkt.SetCapability(c)
	}
}

// WithStmtPrepareResponseServerStatus sets the server status.
func WithStmtPrepareResponseServerStatus(status ServerStatus) StmtPrepareResponseOption {
	return func(pkt *StmtPrepareResponse) {
		pkt.SetServerStatus(status)
	}
}

// WithStmtPrepareResponseStatementID sets the statement ID.
func WithStmtPrepareResponseStatementID(stmdID StatementID) StmtPrepareResponseOption {
	return func(pkt *StmtPrepareResponse) {
		pkt.stmtID = stmdID
	}
}

// WithStmtPrepareResponseColumns sets the columns.
func WithStmtPrepareResponseColumns(columns []ColumnDef) StmtPrepareResponseOption {
	return func(pkt *StmtPrepareResponse) {
		pkt.columns = columns
	}
}

// WithStmtPrepareResponseParams sets the params.
func WithStmtPrepareResponseParams(params []ColumnDef) StmtPrepareResponseOption {
	return func(pkt *StmtPrepareResponse) {
		pkt.params = params
	}
}

// WithStmtPrepareResponseWarningCount sets the warning count.
func WithStmtPrepareResponseWarningCount(warningCount uint16) StmtPrepareResponseOption {
	return func(pkt *StmtPrepareResponse) {
		pkt.warningCount = warningCount
	}
}

// WithStmtPrepareResponseResultSetMetadata sets the result set metadata.
func WithStmtPrepareResponseResultSetMetadata(metadata ResultsetMetadata) StmtPrepareResponseOption {
	return func(pkt *StmtPrepareResponse) {
		pkt.resultSetMetadata = metadata
	}
}

// NewStmtPrepareResponse returns a new StmtPrepareResponse.
func NewStmtPrepareResponse(opts ...StmtPrepareResponseOption) *StmtPrepareResponse {
	h := newStmtPrepareResponseWithPacket(newPacket())
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// NewStmtPrepareResponseFromReader returns a new StmtPrepareResponse from the reader.
func NewStmtPrepareResponseFromReader(reader io.Reader, opts ...StmtPrepareResponseOption) (*StmtPrepareResponse, error) {
	var err error

	pktHeader, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newStmtPrepareResponseWithPacket(pktHeader, opts...)

	i1, err := pktHeader.ReadInt1()
	if err != nil {
		return nil, err
	}
	pkt.status = Status(i1)

	i4, err := pktHeader.ReadInt4()
	if err != nil {
		return nil, err
	}
	pkt.stmtID = StatementID(i4)

	numColumns, err := pktHeader.ReadInt2()
	if err != nil {
		return nil, err
	}

	numParams, err := pktHeader.ReadInt2()
	if err != nil {
		return nil, err
	}

	// 	reserved (1) -- 0x00
	_, err = pktHeader.ReadByte()
	if err != nil {
		return nil, err
	}

	pkt.warningCount, err = pktHeader.ReadInt2()
	if err != nil {
		return nil, err
	}

	if pkt.capability.IsEnabled(ClientOptionalResultsetMetadata) {
		v, err := pktHeader.ReadByte()
		if err != nil {
			return nil, err
		}
		pkt.resultSetMetadata = ResultsetMetadata(v)
	}

	if pkt.resultSetMetadata == ResultsetMetadataNone && numParams == 0 && numColumns == 0 {
		return pkt, nil
	}

	pkt.params = make([]ColumnDef, numParams)
	for n := 0; n < int(numParams); n++ {
		param, err := NewColumnDefFromReader(pktHeader)
		if err != nil {
			return nil, err
		}
		pkt.params[n] = param
	}
	if pkt.Capability().IsDisabled(ClientDeprecateEOF) {
		_, err := NewEOFFromReader(reader, WithEOFCapability(pkt.Capability()))
		if err != nil {
			return nil, err
		}
	}

	pkt.columns = make([]ColumnDef, numColumns)
	for n := 0; n < int(numColumns); n++ {
		column, err := NewColumnDefFromReader(pktHeader)
		if err != nil {
			return nil, err
		}
		pkt.columns[n] = column
	}
	if pkt.Capability().IsDisabled(ClientDeprecateEOF) {
		_, err := NewEOFFromReader(reader, WithEOFCapability(pkt.Capability()))
		if err != nil {
			return nil, err
		}
	}

	return pkt, nil
}

// NewStmtPrepareResponseFromBytes returns a new StmtPrepareResponse from the bytes.
func NewStmtPrepareResponseFromBytes(data []byte, opts ...StmtPrepareResponseOption) (*StmtPrepareResponse, error) {
	return NewStmtPrepareResponseFromReader(bytes.NewReader(data), opts...)
}

// SetStatementID sets the statement ID.
func (pkt *StmtPrepareResponse) SetStatementID(stmtID StatementID) {
	pkt.stmtID = stmtID
}

// StatementID returns the statement ID.
func (pkt *StmtPrepareResponse) StatementID() StatementID {
	return pkt.stmtID
}

// Columns returns the columns.
func (pkt *StmtPrepareResponse) Columns() []ColumnDef {
	return pkt.columns
}

// Params returns the params.
func (pkt *StmtPrepareResponse) Params() []ColumnDef {
	return pkt.params
}

// WarningCount returns the warning count.
func (pkt *StmtPrepareResponse) WarningCount() uint16 {
	return pkt.warningCount
}

// ResultSetMetadata returns the result set metadata.
func (pkt *StmtPrepareResponse) ResultSetMetadata() ResultsetMetadata {
	return pkt.resultSetMetadata
}

// Bytes returns the packet bytes.
func (pkt *StmtPrepareResponse) Bytes() ([]byte, error) {
	payloadLen := 1 + 4 + 2 + 2 + 1 + 2
	if pkt.Capability().IsEnabled(ClientOptionalResultsetMetadata) {
		payloadLen++
	}
	pkt.SetPayloadLength(payloadLen)

	w := NewPacketWriter()

	if _, err := w.WriteBytes(pkt.HeaderBytes()); err != nil {
		return nil, err
	}

	if err := w.WriteInt1(byte(pkt.status)); err != nil {
		return nil, err
	}

	if err := w.WriteInt4(uint32(pkt.stmtID)); err != nil {
		return nil, err
	}

	if err := w.WriteInt2(uint16(len(pkt.columns))); err != nil {
		return nil, err
	}

	if err := w.WriteInt2(uint16(len(pkt.params))); err != nil {
		return nil, err
	}

	// 	reserved (1) -- 0x00
	if err := w.WriteByte(0x00); err != nil {
		return nil, err
	}

	if err := w.WriteInt2(pkt.warningCount); err != nil {
		return nil, err
	}

	if pkt.Capability().IsEnabled(ClientOptionalResultsetMetadata) {
		if err := w.WriteByte(byte(pkt.resultSetMetadata)); err != nil {
			return nil, err
		}
	}

	if pkt.resultSetMetadata == ResultsetMetadataNone && len(pkt.params) == 0 && len(pkt.columns) == 0 {
		return w.Bytes(), nil
	}

	seqID := pkt.SequenceID()
	seqID = seqID.Next()

	for _, param := range pkt.params {
		param.SetSequenceID(seqID)
		if err := w.WritePacket(param); err != nil {
			return nil, err
		}
		seqID = seqID.Next()
	}
	if pkt.Capability().IsDisabled(ClientDeprecateEOF) {
		if err := w.WriteEOF(seqID, pkt.Capability(), pkt.ServerStatus()); err != nil {
			return nil, err
		}
		seqID = seqID.Next()
	}

	for _, column := range pkt.columns {
		column.SetSequenceID(seqID)
		if err := w.WritePacket(column); err != nil {
			return nil, err
		}
		seqID = seqID.Next()
	}
	if pkt.Capability().IsDisabled(ClientDeprecateEOF) {
		if err := w.WriteEOF(seqID, pkt.Capability(), pkt.ServerStatus()); err != nil {
			return nil, err
		}
	}

	return w.Bytes(), nil
}
