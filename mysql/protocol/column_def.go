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

// MySQL: Column Definition
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_com_query_response_text_resultset_column_definition.html

// ResultSet represents a MySQL text resultset packet.
type ColumnDef struct {
	*packet
	catalog   string
	schema    string
	table     string
	orgName   string
	charSet   uint16
	colLength uint32
	colType   uint8
	flags     uint16
	decimals  uint8
}

func NewColumnDefFromReader(r io.Reader) (*ColumnDef, error) {
	var err error

	pkt, err := NewPacketWithReader(r)
	if err != nil {
		return nil, err
	}

	colDef := &ColumnDef{
		packet: nil,
	}

	colDef.catalog, err = pkt.ReadLengthEncodedString()
	if err != nil {
		return nil, err
	}

	colDef.schema, err = pkt.ReadLengthEncodedString()
	if err != nil {
		return nil, err
	}

	colDef.table, err = pkt.ReadLengthEncodedString()
	if err != nil {
		return nil, err
	}

	colDef.orgName, err = pkt.ReadLengthEncodedString()
	if err != nil {
		return nil, err
	}

	colDef.charSet, err = pkt.ReadInt2()
	if err != nil {
		return nil, err
	}

	colDef.colLength, err = pkt.ReadInt4()
	if err != nil {
		return nil, err
	}

	colDef.colType, err = pkt.ReadInt1()
	if err != nil {
		return nil, err
	}

	colDef.flags, err = pkt.ReadInt2()
	if err != nil {
		return nil, err
	}

	colDef.decimals, err = pkt.ReadInt1()
	if err != nil {
		return nil, err
	}

	return colDef, nil
}

// Catalog returns the column catalog.
func (pkt *ColumnDef) Catalog() string {
	return pkt.catalog
}

// Schema returns the column schema.
func (pkt *ColumnDef) Schema() string {
	return pkt.schema
}

// Table returns the column table.
func (pkt *ColumnDef) Table() string {
	return pkt.table
}

// OrgName returns the column original name.
func (pkt *ColumnDef) OrgName() string {
	return pkt.orgName
}

// CharSet returns the column character set.
func (pkt *ColumnDef) CharSet() uint16 {
	return pkt.charSet
}

// ColLength returns the column length.
func (pkt *ColumnDef) ColLength() uint32 {
	return pkt.colLength
}

// ColType returns the column type.
func (pkt *ColumnDef) ColType() uint8 {
	return pkt.colType
}

// Flags returns the column flags.
func (pkt *ColumnDef) Flags() uint16 {
	return pkt.flags
}

// Decimals returns the column decimals.
func (pkt *ColumnDef) Decimals() uint8 {
	return pkt.decimals
}

// Bytes returns the packet bytes.
func (pkt *ColumnDef) Bytes() ([]byte, error) {
	w := NewWriter()

	if err := w.WriteLengthEncodedString(pkt.catalog); err != nil {
		return nil, err
	}

	if err := w.WriteLengthEncodedString(pkt.schema); err != nil {
		return nil, err
	}

	if err := w.WriteLengthEncodedString(pkt.table); err != nil {
		return nil, err
	}

	if err := w.WriteLengthEncodedString(pkt.orgName); err != nil {
		return nil, err
	}

	if err := w.WriteInt2(pkt.charSet); err != nil {
		return nil, err
	}

	if err := w.WriteInt4(pkt.colLength); err != nil {
		return nil, err
	}

	if err := w.WriteInt1(pkt.colType); err != nil {
		return nil, err
	}

	if err := w.WriteInt2(pkt.flags); err != nil {
		return nil, err
	}

	if err := w.WriteInt1(pkt.decimals); err != nil {
		return nil, err
	}

	pkt.packet = NewPacket(
		PacketWithSequenceID(pkt.packet.SequenceID()),
		PacketWithPayload(w.Bytes()),
	)

	return pkt.packet.Bytes()
}
