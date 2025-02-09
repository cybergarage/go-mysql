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
	"fmt"

	"github.com/cybergarage/go-mysql/mysql/query"
)

// MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/

// BinaryResultSetRowOption represents a MySQL binary resultset row option.
type BinaryResultSetRowOption func(*BinaryResultSetRow)

// BinaryResultSetRow represents a MySQL binary resultset row response packet.
type BinaryResultSetRow struct {
	t     FieldType
	bytes []byte
}

func newBinaryResultSetRowWithPacket(opts ...BinaryResultSetRowOption) *BinaryResultSetRow {
	row := &BinaryResultSetRow{
		t:     0,
		bytes: nil,
	}
	row.SetOptions(opts...)
	return row
}

// WithBinaryResultSetRowType returns a binary resultset row option to set the type.
func WithwBinaryResultSetRowType(t FieldType) BinaryResultSetRowOption {
	return func(row *BinaryResultSetRow) {
		row.t = t
	}
}

// WithBinaryResultSetRowBytes returns a binary resultset row option to set the bytes.
func WithwBinaryResultSetRowBytes(b []byte) BinaryResultSetRowOption {
	return func(row *BinaryResultSetRow) {
		row.bytes = b
	}
}

// NewBinaryResultSetRow returns a new BinaryResultSetRow.
func NewBinaryResultSetRow(opts ...BinaryResultSetRowOption) *BinaryResultSetRow {
	return newBinaryResultSetRowWithPacket(opts...)
}

// NewBinaryResultSetRowFromReader returns a new BinaryResultSetRow from the reader.
func NewBinaryResultSetRowFromReader(reader *Reader, opts ...BinaryResultSetRowOption) (*BinaryResultSetRow, error) {
	row := newBinaryResultSetRowWithPacket(opts...)

	byteLen := 0

	switch row.t {
	case query.MySQLTypeString, query.MySQLTypeVarString, query.MySQLTypeVarchar:
		v, err := reader.ReadLengthEncodedString()
		if err != nil {
			return nil, err
		}
		row.bytes = []byte(v)
	case query.MySQLTypeTinyBlob, query.MySQLTypeMediumBlob, query.MySQLTypeLongBlob, query.MySQLTypeBlob:
		v, err := reader.ReadLengthEncodedBytes()
		if err != nil {
			return nil, err
		}
		row.bytes = v
	case query.MySQLTypeNull:
		byteLen = 0
	case query.MySQLTypeTiny:
		byteLen = 1
	case query.MySQLTypeShort, query.MySQLTypeYear:
		byteLen = 2
	case query.MySQLTypeLong, query.MySQLTypeFloat, query.MySQLTypeInt24:
		byteLen = 4
	case query.MySQLTypeLonglong, query.MySQLTypeDouble:
		byteLen = 8
	case query.MySQLTypeDate, query.MySQLTypeTime, query.MySQLTypeDatetime, query.MySQLTypeTimestamp:
		l, err := reader.ReadInt1()
		if err != nil {
			return nil, err
		}
		row.bytes, err = reader.ReadFixedLengthBytes(int(l))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w field type: %s(%v)", ErrNotSupported, FieldType(row.t).String(), row.t)
	}

	if 0 < byteLen {
		row.bytes = make([]byte, byteLen)
		if _, err := reader.Read(row.bytes); err != nil {
			return nil, err
		}
	}

	return row, nil
}

// SetOptions sets the options.
func (row *BinaryResultSetRow) SetOptions(opts ...BinaryResultSetRowOption) {
	for _, opt := range opts {
		opt(row)
	}
}

// Type returns the type.
func (row *BinaryResultSetRow) Type() FieldType {
	return row.t
}

// Bytes returns the bytes.
func (row *BinaryResultSetRow) Bytes() ([]byte, error) {
	return row.bytes, nil
}
