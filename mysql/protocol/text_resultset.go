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

// MySQL: Protocol::QueryResponse
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query.html
// MySQL: Text Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset.html

// TextResultSet represents a MySQL text resultset response packet.
type TextResultSet struct {
	*packet
	capFlags        CapabilityFlag
	metadataFollows ResultsetMetadata
	columnCount     uint64
	columnDefs      []*ColumnDef
	rows            []Row
}

func newTextResultSetWithPacket(pkt *packet, opts ...TextResultSetOption) *TextResultSet {
	q := &TextResultSet{
		packet:          pkt,
		capFlags:        0,
		metadataFollows: 0,
		columnCount:     0,
		columnDefs:      []*ColumnDef{},
		rows:            []Row{},
	}
	q.SetOptions(opts...)
	return q
}

// TextResultSetOption represents a COM_QUERY text resultset option.
type TextResultSetOption func(*TextResultSet)

// WithTextResultSetCapabilities returns a text resultset option to set the capabilities.
func WithTextResultSetCapabilities(c CapabilityFlag) TextResultSetOption {
	return func(pkt *TextResultSet) {
		pkt.capFlags = c
	}
}

// WithTextResultSetMetadataFollows returns a text resultset option to set the metadata follows.
func WithTextResultSetMetadataFollows(m ResultsetMetadata) TextResultSetOption {
	return func(pkt *TextResultSet) {
		pkt.metadataFollows = m
	}
}

// WithTextResultSetColumnDefs returns a text resultset option to set the column definitions.
func WithTextResultSetColumnDefs(colDefs []*ColumnDef) TextResultSetOption {
	return func(pkt *TextResultSet) {
		pkt.columnCount = uint64(len(colDefs))
		pkt.columnDefs = colDefs
	}
}

// NewTextResultSet returns a new text resultset response packet.
func NewTextResultSet(opts ...TextResultSetOption) (*TextResultSet, error) {
	pkt := newTextResultSetWithPacket(nil, opts...)
	return pkt, nil
}

// NewTextResultSetFromReader returns a new text resultset response packet from the specified reader.
func NewTextResultSetFromReader(reader io.Reader, opts ...TextResultSetOption) (*TextResultSet, error) {
	var err error

	pktReader, err := NewPacketWithReader(reader)
	if err != nil {
		return nil, err
	}

	pkt := newTextResultSetWithPacket(pktReader, opts...)

	if pkt.Capabilities().IsEnabled(ClientOptionalResultsetMetadata) {
		pkt.metadataFollows, err = pkt.ReadByte()
		if err != nil {
			return nil, err
		}
	}

	pkt.columnCount, err = pkt.ReadLengthEncodedInt()
	if err != nil {
		return nil, err
	}

	if pkt.Capabilities().IsDisabled(ClientOptionalResultsetMetadata) || pkt.metadataFollows == ResultsetMetadataFull {
		for i := 0; i < int(pkt.columnCount); i++ {
			colDef, err := NewColumnDefFromReader(reader)
			if err != nil {
				return nil, err
			}
			pkt.columnDefs = append(pkt.columnDefs, colDef)
		}
	}

	if pkt.Capabilities().IsDisabled(ClientDeprecateEOF) {
		_, err := NewEOFFromReader(reader, WithEOFCapability(pkt.Capabilities()))
		if err != nil {
			return nil, err
		}
	}

	// One or more Text Resultset Row
	row, err := NewTextResultSetRowFromReader(pkt.Reader(), WithTextResultSetRowColmunCount(pkt.columnCount))
	if err != nil {
		return nil, err
	}
	pkt.rows = append(pkt.rows, row)

	if pkt.Capabilities().IsEnabled(ClientDeprecateEOF) {
		_, err := NewOKFromReader(reader, WithOKCapability(pkt.Capabilities()))
		if err != nil {
			return nil, err
		}
	} else {
		_, err := NewEOFFromReader(reader, WithEOFCapability(pkt.Capabilities()))
		if err != nil {
			return nil, err
		}
	}

	return pkt, nil
}

// SetOptions sets the options.
func (pkt *TextResultSet) SetOptions(opts ...TextResultSetOption) {
	for _, opt := range opts {
		opt(pkt)
	}
}

// Capabilities returns the capabilities.
func (pkt *TextResultSet) Capabilities() CapabilityFlag {
	return pkt.capFlags
}

// Rows returns the rows.
func (pkt *TextResultSet) Rows() []Row {
	return pkt.rows
}

// Bytes returns the packet bytes.
func (pkt *TextResultSet) Bytes() ([]byte, error) {
	w := NewPacketWriter()

	if pkt.Capabilities().IsEnabled(ClientOptionalResultsetMetadata) {
		err := w.WriteByte(pkt.metadataFollows)
		if err != nil {
			return nil, err
		}
	}

	err := w.WriteLengthEncodedInt(pkt.columnCount)
	if err != nil {
		return nil, err
	}

	if pkt.Capabilities().IsDisabled(ClientOptionalResultsetMetadata) || pkt.metadataFollows == ResultsetMetadataFull {
		for _, colDef := range pkt.columnDefs {
			colDefBytes, err := colDef.Bytes()
			if err != nil {
				return nil, err
			}
			_, err = w.WriteBytes(colDefBytes)
			if err != nil {
				return nil, err
			}
		}
	}

	if pkt.Capabilities().IsDisabled(ClientDeprecateEOF) {
		err := w.WriteEOF()
		if err != nil {
			return nil, err
		}
	}

	// One or more Text Resultset Row
	for _, row := range pkt.rows {
		rowBytes, err := row.Bytes()
		if err != nil {
			return nil, err
		}
		_, err = w.WriteBytes(rowBytes)
		if err != nil {
			return nil, err
		}
	}

	if pkt.Capabilities().IsEnabled(ClientDeprecateEOF) {
		err := w.WriteOK(pkt.Capabilities())
		if err != nil {
			return nil, err
		}
	} else {
		err := w.WriteEOF(pkt.Capabilities())
		if err != nil {
			return nil, err
		}
	}

	pkt.packet =
		NewPacket(
			PacketWithSequenceID(pkt.SequenceID()),
			PacketWithPayload(w.Bytes()),
		)

	return pkt.packet.Bytes()
}
