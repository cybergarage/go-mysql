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
	"io"
)

// MySQL: COM_STMT_PREPARE Response
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_stmt_prepare.html
// COM_STMT_PREPARE response - MariaDB Knowledge Base
// https://mariadb.com/kb/en/com_stmt_prepare/
// Server Response Packets (Binary Protocol) - MariaDB Knowledge Base
// https://mariadb.com/kb/en/server-response-packets-binary-protocol/

// StatementID represents a statement ID.
type StatementID uint32

// StmtPrepareResponse represents a MySQL Prepare Response packet.
type StmtPrepareResponse struct {
	*packet
	caps              Capability
	status            Status
	stmtID            StatementID
	columns           []*ColumnDef
	params            []*ColumnDef
	warningCount      uint16
	resultSetMetadata ResultsetMetadata
}

func newStmtPrepareResponseWithPacket(pkt *packet) *StmtPrepareResponse {
	return &StmtPrepareResponse{
		packet:            pkt,
		caps:              DefaultServerCapabilities,
		status:            Status(0),
		stmtID:            0,
		columns:           []*ColumnDef{},
		params:            []*ColumnDef{},
		warningCount:      0,
		resultSetMetadata: ResultsetMetadataNone,
	}
}

// StmtPrepareResponseOption represents a StmtPrepareResponse option.
type StmtPrepareResponseOption func(*StmtPrepareResponse)

// WithPrepareResponseCapabilitys sets the capabilitys.
func WithPrepareResponseCapabilitys(cap Capability) StmtPrepareResponseOption {
	return func(h *StmtPrepareResponse) {
		h.caps = cap
	}
}

// WithPrepareResponseStatementID sets the statement ID.
func WithPrepareResponseStatementID(stmdID StatementID) StmtPrepareResponseOption {
	return func(h *StmtPrepareResponse) {
		h.stmtID = stmdID
	}
}

// WithPrepareResponseColumns sets the columns.
func WithPrepareResponseColumns(columns []*ColumnDef) StmtPrepareResponseOption {
	return func(h *StmtPrepareResponse) {
		h.columns = columns
	}
}

// WithPrepareResponseParams sets the params.
func WithPrepareResponseParams(params []*ColumnDef) StmtPrepareResponseOption {
	return func(h *StmtPrepareResponse) {
		h.params = params
	}
}

// WithPrepareResponseWarningCount sets the warning count.
func WithPrepareResponseWarningCount(warningCount uint16) StmtPrepareResponseOption {
	return func(h *StmtPrepareResponse) {
		h.warningCount = warningCount
	}
}

// WithPrepareResponseResultSetMetadata sets the result set metadata.
func WithPrepareResponseResultSetMetadata(metadata ResultsetMetadata) StmtPrepareResponseOption {
	return func(h *StmtPrepareResponse) {
		h.resultSetMetadata = metadata
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
func NewStmtPrepareResponseFromReader(reader io.Reader) (*StmtPrepareResponse, error) {
	var err error

	pktReader, err := NewPacketHeaderWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newStmtPrepareResponseWithPacket(pktReader)

	i1, err := pktReader.ReadInt1()
	if err != nil {
		return nil, err
	}
	pkt.status = Status(i1)

	i4, err := pktReader.ReadInt4()
	if err != nil {
		return nil, err
	}
	pkt.stmtID = StatementID(i4)

	numColumns, err := pktReader.ReadInt2()
	if err != nil {
		return nil, err
	}

	numParams, err := pktReader.ReadInt2()
	if err != nil {
		return nil, err
	}

	// 	reserved (1) -- 0x00
	_, err = pktReader.ReadByte()
	if err != nil {
		return nil, err
	}

	if pkt.PayloadLength() < 12 {
		return pkt, nil
	}

	pkt.warningCount, err = pktReader.ReadInt2()
	if err != nil {
		return nil, err
	}

	if pkt.Capabilitys.IsEnabled(ClientOptionalResultsetMetadata) {
		v, err := pktReader.ReadByte()
		if err != nil {
			return nil, err
		}
		pkt.resultSetMetadata = ResultsetMetadata(v)
	}

	if pkt.resultSetMetadata == ResultsetMetadataNone && numParams == 0 && numColumns == 0 {
		return pkt, nil
	}

	pkt.params = make([]*ColumnDef, 0, numParams)
	for n := 0; n < int(numParams); n++ {
		param, err := NewColumnDefFromReader(pktReader)
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

	pkt.columns = make([]*ColumnDef, 0, numColumns)
	for n := 0; n < int(numColumns); n++ {
		column, err := NewColumnDefFromReader(pktReader)
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

// StatementID returns the statement ID.
func (pkt *StmtPrepareResponse) StatementID() StatementID {
	return pkt.stmtID
}

// Columns returns the columns.
func (pkt *StmtPrepareResponse) Columns() []*ColumnDef {
	return pkt.columns
}

// Params returns the params.
func (pkt *StmtPrepareResponse) Params() []*ColumnDef {
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

// Capability returns the capability.
func (pkt *StmtPrepareResponse) Capability() Capability {
	return pkt.caps
}

// Bytes returns the packet bytes.
func (pkt *StmtPrepareResponse) Bytes() ([]byte, error) {
	w := NewPacketWriter()

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

	if pkt.Capabilitys.IsEnabled(ClientOptionalResultsetMetadata) {
		if err := w.WriteByte(byte(pkt.resultSetMetadata)); err != nil {
			return nil, err
		}
	}

	if pkt.resultSetMetadata == ResultsetMetadataNone && len(pkt.params) == 0 && len(pkt.columns) == 0 {
		return pkt.packet.Bytes()
	}

	for _, param := range pkt.params {
		if err := w.WritePacket(param); err != nil {
			return nil, err
		}
	}
	if pkt.Capability().IsDisabled(ClientDeprecateEOF) {
		if err := w.WriteEOF(pkt.Capability()); err != nil {
			return nil, err
		}
	}

	for _, column := range pkt.columns {
		if err := w.WritePacket(column); err != nil {
			return nil, err
		}
	}
	if pkt.Capability().IsDisabled(ClientDeprecateEOF) {
		if err := w.WriteEOF(pkt.Capability()); err != nil {
			return nil, err
		}
	}

	return pkt.packet.Bytes()
}
