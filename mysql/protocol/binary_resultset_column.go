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

import "github.com/cybergarage/go-mysql/mysql/stmt"

// MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/

// BinaryResultSetColumnOption represents a MySQL binary resultset row option.
type BinaryResultSetColumnOption func(*BinaryResultSetColumn) error

// BinaryResultSetColumn represents a MySQL binary resultset row response packet.
type BinaryResultSetColumn struct {
	t     FieldType
	bytes []byte
}

// WithBinaryResultSetRowType returns a binary resultset row option to set the type.
func WithBinaryResultSetColumnType(t FieldType) BinaryResultSetColumnOption {
	return func(column *BinaryResultSetColumn) error {
		column.t = t
		return nil
	}
}

// WithBinaryResultSetRowBytes returns a binary resultset row option to set the bytes.
func WithBinaryResultSetColumnBytes(b []byte) BinaryResultSetColumnOption {
	return func(column *BinaryResultSetColumn) error {
		column.bytes = b
		return nil
	}
}

// WithBinaryResultSetRowValue returns a binary resultset row option to set the value.
func WithBinaryResultSetColumnValue(v any) BinaryResultSetColumnOption {
	return func(column *BinaryResultSetColumn) error {
		w := NewPacketWriter()
		err := w.WriteFieldValue(column.t, v)
		if err != nil {
			return err
		}

		reader := NewPacketReaderWithBytes(w.Bytes())
		v, err := reader.ReadFieldBytes(column.t)
		if err != nil {
			return err
		}
		column.bytes = v

		return nil
	}
}

// NewBinaryResultSetColumn returns a new BinaryResultSetColumn.
func NewBinaryResultSetColumn(opts ...BinaryResultSetColumnOption) (*BinaryResultSetColumn, error) {
	column := &BinaryResultSetColumn{
		t:     0,
		bytes: nil,
	}
	for _, opt := range opts {
		if err := opt(column); err != nil {
			return nil, err
		}
	}
	return column, nil
}

// NewBinaryResultSetColumnFromReader returns a new BinaryResultSetColumn from the reader.
func NewBinaryResultSetColumnFromReader(reader *PacketReader, opts ...BinaryResultSetColumnOption) (*BinaryResultSetColumn, error) {
	column, err := NewBinaryResultSetColumn(opts...)
	if err != nil {
		return nil, err
	}

	v, err := reader.ReadFieldBytes(column.t)
	if err != nil {
		return nil, err
	}
	column.bytes = v

	return column, nil
}

// Type returns the type.
func (column *BinaryResultSetColumn) Type() FieldType {
	return column.t
}

// Value returns the value.
func (column *BinaryResultSetColumn) Value() (any, error) {
	field := stmt.NewField(
		stmt.WithFieldType(column.t),
		stmt.WithFieldBytes(column.bytes),
	)
	return field.Value()
}

// Bytes returns the bytes.
func (column *BinaryResultSetColumn) Bytes() ([]byte, error) {
	w := NewPacketWriter()
	if err := w.WriteFieldBytes(column.t, column.bytes); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}
