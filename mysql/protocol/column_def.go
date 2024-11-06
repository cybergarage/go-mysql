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
// Result Set Packets - MariaDB Knowledge Base
// https://mariadb.com/kb/en/result-set-packets/#column-definition-packet

const (
	defaultColumnDefCatalog       = "def"
	defaultColumnFixedFieldLength = 0x0c
	defaultColumnDefMaxLen        = 0xFFFF
	defaultColumnCharSet          = uint16(CharSetUTF8)
)

// ColumnDefOption represents a function to set a ColumnDef option.
type ColumnDefOption func(*ColumnDef)

// ColumnDef represents a MySQL Column Definition packet.
type ColumnDef struct {
	*packet
	catalog          string
	schema           string
	table            string
	orgTable         string
	name             string
	orgName          string
	fixedFieldLength uint64
	charSet          uint16
	colLength        uint32
	colType          uint8
	flags            uint16
	decimals         uint8
}

func newColumnDefWith(p *packet, opts ...ColumnDefOption) *ColumnDef {
	pkt := &ColumnDef{
		packet:           p,
		catalog:          defaultColumnDefCatalog,
		schema:           "",
		table:            "",
		orgTable:         "",
		name:             "",
		orgName:          "",
		fixedFieldLength: defaultColumnFixedFieldLength,
		charSet:          defaultColumnCharSet,
		colLength:        defaultColumnDefMaxLen,
		colType:          0,
		flags:            0,
		decimals:         0,
	}
	pkt.SetOptions(opts...)
	return pkt
}

// WithColumnDefSchema returns a ColumnDefOption to set the schema.
func WithColumnDefSchema(schema string) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.schema = schema
	}
}

// WithColumnDefTable returns a ColumnDefOption to set the table.
func WithColumnDefTable(table string) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.table = table
		pkt.orgTable = table
	}
}

// WithColumnDefOrgTable returns a ColumnDefOption to set the original table.
func WithColumnDefOrgTable(orgTable string) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.orgTable = orgTable
	}
}

// WithColumnDefName returns a ColumnDefOption to set the name.
func WithColumnDefName(name string) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.name = name
		pkt.orgName = name
	}
}

// WithColumnDefOrgName returns a ColumnDefOption to set the original name.
func WithColumnDefOrgName(orgName string) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.orgName = orgName
	}
}

// WithColumnDefColLength returns a ColumnDefOption to set the column length.
func WithColumnDefFixedFieldLength(fixedFieldLength uint64) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.fixedFieldLength = fixedFieldLength
	}
}

// WithColumnDefColLength returns a ColumnDefOption to set the column length.
func WithColumnDefCharSet(charSet uint16) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.charSet = charSet
	}
}

// WithColumnDefColLength returns a ColumnDefOption to set the column length.
func WithColumnDefType(colType uint8) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.colType = colType
	}
}

// WithColumnDefFlags returns a ColumnDefOption to set the flags.
func WithColumnDefFlags(flags uint16) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.flags = flags
	}
}

// WithColumnDefDecimals returns a ColumnDefOption to set the decimals.
func WithColumnDefDecimals(decimals uint8) ColumnDefOption {
	return func(pkt *ColumnDef) {
		pkt.decimals = decimals
	}
}

// NewColumnDef returns a new ColumnDef.
func NewColumnDef(opts ...ColumnDefOption) *ColumnDef {
	return newColumnDefWith(newPacket(), opts...)
}

// NewColumnDefFromReader returns a new ColumnDef from the reader.
func NewColumnDefFromReader(r io.Reader) (*ColumnDef, error) {
	var err error

	pkt, err := NewPacketHeaderWithReader(r)
	if err != nil {
		return nil, err
	}

	colDef := newColumnDefWith(pkt)

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

	colDef.orgTable, err = pkt.ReadLengthEncodedString()
	if err != nil {
		return nil, err
	}

	colDef.name, err = pkt.ReadLengthEncodedString()
	if err != nil {
		return nil, err
	}

	colDef.orgName, err = pkt.ReadLengthEncodedString()
	if err != nil {
		return nil, err
	}

	colDef.fixedFieldLength, err = pkt.ReadLengthEncodedInt()
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

	// unused

	_, err = pkt.ReadInt2()
	if err != nil {
		return nil, err
	}

	return colDef, nil
}

// SetOptions sets the ColumnDef options.
func (pkt *ColumnDef) SetOptions(opts ...ColumnDefOption) {
	for _, opt := range opts {
		opt(pkt)
	}
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

// OrgTable returns the column original table.
func (pkt *ColumnDef) OrgTable() string {
	return pkt.orgTable
}

// Name returns the column name.
func (pkt *ColumnDef) Name() string {
	return pkt.name
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

	if err := w.WriteLengthEncodedString(pkt.orgTable); err != nil {
		return nil, err
	}

	if err := w.WriteLengthEncodedString(pkt.name); err != nil {
		return nil, err
	}

	if err := w.WriteLengthEncodedString(pkt.orgName); err != nil {
		return nil, err
	}

	if err := w.WriteLengthEncodedInt(pkt.fixedFieldLength); err != nil {
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

	// unused

	if err := w.WriteInt2(0); err != nil {
		return nil, err
	}

	pkt.SetPayload(w.Bytes())

	return pkt.packet.Bytes()
}
