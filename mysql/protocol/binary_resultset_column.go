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

// BinaryResultSetColumnOption represents a MySQL binary resultset row option.
type BinaryResultSetColumnOption func(*BinaryResultSetColumn)

// BinaryResultSetColumn represents a MySQL binary resultset row response packet.
type BinaryResultSetColumn struct {
	t     FieldType
	bytes []byte
}

// WithBinaryResultSetRowType returns a binary resultset row option to set the type.
func WithwBinaryResultSetColumnType(t FieldType) BinaryResultSetColumnOption {
	return func(row *BinaryResultSetColumn) {
		row.t = t
	}
}

// WithBinaryResultSetRowBytes returns a binary resultset row option to set the bytes.
func WithwBinaryResultSetColumnBytes(b []byte) BinaryResultSetColumnOption {
	return func(row *BinaryResultSetColumn) {
		row.bytes = b
	}
}

// NewBinaryResultSetColumn returns a new BinaryResultSetColumn.
func NewBinaryResultSetColumn(opts ...BinaryResultSetColumnOption) *BinaryResultSetColumn {
	column := &BinaryResultSetColumn{
		t:     0,
		bytes: nil,
	}
	for _, opt := range opts {
		opt(column)
	}
	return column
}

// NewBinaryResultSetColumnFromReader returns a new BinaryResultSetColumn from the reader.
func NewBinaryResultSetColumnFromReader(reader *PacketReader, opts ...BinaryResultSetColumnOption) (*BinaryResultSetColumn, error) {
	column := NewBinaryResultSetColumn(opts...)

	byteLen := 0

	switch column.t {
	case query.MySQLTypeString, query.MySQLTypeVarString, query.MySQLTypeVarchar:
		v, err := reader.ReadLengthEncodedString()
		if err != nil {
			return nil, err
		}
		column.bytes = []byte(v)
	case query.MySQLTypeTinyBlob, query.MySQLTypeMediumBlob, query.MySQLTypeLongBlob, query.MySQLTypeBlob:
		v, err := reader.ReadLengthEncodedBytes()
		if err != nil {
			return nil, err
		}
		column.bytes = v
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
		column.bytes, err = reader.ReadFixedLengthBytes(int(l))
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("%w field type: %s(%v)", ErrNotSupported, FieldType(column.t).String(), column.t)
	}

	if 0 < byteLen {
		column.bytes = make([]byte, byteLen)
		if _, err := reader.Read(column.bytes); err != nil {
			return nil, err
		}
	}

	return column, nil
}

// Type returns the type.
func (column *BinaryResultSetColumn) Type() FieldType {
	return column.t
}

// Bytes returns the bytes.
func (column *BinaryResultSetColumn) Bytes() ([]byte, error) {
	switch column.t {
	case query.MySQLTypeString, query.MySQLTypeVarString, query.MySQLTypeVarchar:
		w := NewPacketWriter()
		if err := w.WriteLengthEncodedString(string(column.bytes)); err != nil {
			return nil, err
		}
		return w.Bytes(), nil
	case query.MySQLTypeTinyBlob, query.MySQLTypeMediumBlob, query.MySQLTypeLongBlob, query.MySQLTypeBlob:
		w := NewPacketWriter()
		if err := w.WriteLengthEncodedBytes(column.bytes); err != nil {
			return nil, err
		}
		return w.Bytes(), nil
	case query.MySQLTypeDate, query.MySQLTypeTime, query.MySQLTypeDatetime, query.MySQLTypeTimestamp:
		w := NewPacketWriter()
		if err := w.WriteInt1(byte(len(column.bytes))); err != nil {
			return nil, err
		}
		if err := w.WriteFixedLengthBytes(column.bytes, len(column.bytes)); err != nil {
			return nil, err
		}
		return w.Bytes(), nil
	case query.MySQLTypeNull:
		return nil, nil
	}

	return column.bytes, nil
}
