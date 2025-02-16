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

	"github.com/cybergarage/go-mysql/mysql/encoding/binary"
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

// ColumnDefOption represents a function to set a columnDef option.
type ColumnDefOption func(*columnDef)

// columnDef represents a MySQL Column Definition packet.
type columnDef struct {
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

func newColumnDefWith(p *packet, opts ...ColumnDefOption) *columnDef {
	pkt := &columnDef{
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
	return func(pkt *columnDef) {
		pkt.schema = schema
	}
}

// WithColumnDefTable returns a ColumnDefOption to set the table.
func WithColumnDefTable(table string) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.table = table
		pkt.orgTable = table
	}
}

// WithColumnDefOrgTable returns a ColumnDefOption to set the original table.
func WithColumnDefOrgTable(orgTable string) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.orgTable = orgTable
	}
}

// WithColumnDefName returns a ColumnDefOption to set the name.
func WithColumnDefName(name string) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.name = name
		pkt.orgName = name
	}
}

// WithColumnDefOrgName returns a ColumnDefOption to set the original name.
func WithColumnDefOrgName(orgName string) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.orgName = orgName
	}
}

// WithColumnDefColLength returns a ColumnDefOption to set the column length.
func WithColumnDefFixedFieldLength(fixedFieldLength uint64) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.fixedFieldLength = fixedFieldLength
	}
}

// WithColumnDefColLength returns a ColumnDefOption to set the column length.
func WithColumnDefCharSet(charSet uint16) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.charSet = charSet
	}
}

// WithColumnDefColLength returns a ColumnDefOption to set the column length.
func WithColumnDefType(colType uint8) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.colType = colType
	}
}

// WithColumnDefFlags returns a ColumnDefOption to set the flags.
func WithColumnDefFlags(flags uint16) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.flags = flags
	}
}

// WithColumnDefDecimals returns a ColumnDefOption to set the decimals.
func WithColumnDefDecimals(decimals uint8) ColumnDefOption {
	return func(pkt *columnDef) {
		pkt.decimals = decimals
	}
}

// NewColumnDef returns a new columnDef.
func NewColumnDef(opts ...ColumnDefOption) *columnDef {
	return newColumnDefWith(newPacket(), opts...)
}

// NewColumnDefFromReader returns a new columnDef from the reader.
func NewColumnDefFromReader(r io.Reader) (ColumnDef, error) {
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

// SetOptions sets the columnDef options.
func (pkt *columnDef) SetOptions(opts ...ColumnDefOption) {
	for _, opt := range opts {
		opt(pkt)
	}
}

// Catalog returns the column catalog.
func (pkt *columnDef) Catalog() string {
	return pkt.catalog
}

// Schema returns the column schema.
func (pkt *columnDef) Schema() string {
	return pkt.schema
}

// Table returns the column table.
func (pkt *columnDef) Table() string {
	return pkt.table
}

// OrgTable returns the column original table.
func (pkt *columnDef) OrgTable() string {
	return pkt.orgTable
}

// Name returns the column name.
func (pkt *columnDef) Name() string {
	return pkt.name
}

// OrgName returns the column original name.
func (pkt *columnDef) OrgName() string {
	return pkt.orgName
}

// CharSet returns the column character set.
func (pkt *columnDef) CharSet() uint16 {
	return pkt.charSet
}

// ColLength returns the column length.
func (pkt *columnDef) ColLength() uint32 {
	return pkt.colLength
}

// ColType returns the column type.
func (pkt *columnDef) ColType() uint8 {
	return pkt.colType
}

// Flags returns the column flags.
func (pkt *columnDef) Flags() uint16 {
	return pkt.flags
}

// Decimals returns the column decimals.
func (pkt *columnDef) Decimals() uint8 {
	return pkt.decimals
}

// Bytes returns the packet bytes.
func (pkt *columnDef) Bytes() ([]byte, error) {
	w := binary.NewWriter()

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
