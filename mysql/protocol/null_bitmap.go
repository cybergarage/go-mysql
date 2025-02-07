// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
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

// NULL-Bitmap - MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/

// NullBitmap represents a MySQL null bitmap.
type NullBitmap struct {
	numFields int
	offset    int
	bytes     []byte
}

// NullBitmapOption represents a MySQL NullBitmap option.
type NullBitmapOption func(*NullBitmap)

// WithNullBitmapOffset sets the offset of the NullBitmap.
func WithNullBitmapOffset(offset int) NullBitmapOption {
	return func(bmap *NullBitmap) {
		bmap.offset = offset
	}
}

// WithNullBitmapBytes sets the bytes of the NullBitmap.
func WithNullBitmapBytes(bytes []byte) NullBitmapOption {
	return func(bmap *NullBitmap) {
		bmap.bytes = bytes
	}
}

// WithNullBitmapNumFields sets the number of fields of the NullBitmap.
func WithNullBitmapNumFields(numFields int) NullBitmapOption {
	return func(bmap *NullBitmap) {
		bmap.numFields = numFields
	}
}

// NewNullBitmap creates a new NullBitmap with the given length.
func NewNullBitmap(opts ...NullBitmapOption) *NullBitmap {
	bmap := &NullBitmap{
		numFields: 0,
		offset:    0,
		bytes:     make([]byte, 0),
	}
	for _, opt := range opts {
		opt(bmap)
	}
	if bmap.bytes == nil {
		bmap.bytes = make([]byte, (bmap.numFields+7+bmap.offset)/8)
	}
	return bmap
}

// SetNumFields sets the number of fields of the NullBitmap.
func (bmap *NullBitmap) SetNumFields(numFields int) {
	bmap.numFields = numFields
}

// NumFields returns the number of fields of the NullBitmap.
func (bmap *NullBitmap) NumFields() int {
	return bmap.numFields
}

// SetOffset sets the offset of the NullBitmap.
func (bmap *NullBitmap) SetOffset(offset int) {
	bmap.offset = offset
}

// Offset returns the offset of the NullBitmap.
func (bmap *NullBitmap) Offset() int {
	return bmap.offset
}

// SetNull sets the null value of the NullBitmap.
func (bmap *NullBitmap) SetNull(i int, v bool) {
	// 	NULL-bitmap-byte = ((field-pos + offset) / 8)
	// NULL-bitmap-bit  = ((field-pos + offset) % 8)
	idx := bmap.offset + i
	if v {
		bmap.bytes[idx/8] |= 1 << uint(idx%8)
	} else {
		bmap.bytes[idx/8] &= ^(1 << uint(idx%8))
	}
}

// IsNull returns true if the i-th bit is null.
func (bmap *NullBitmap) IsNull(i int) bool {
	// 	NULL-bitmap-byte = ((field-pos + offset) / 8)
	// NULL-bitmap-bit  = ((field-pos + offset) % 8)
	idx := bmap.offset + i
	return bmap.bytes[idx/8]&(1<<uint(idx%8)) != 0
}

// Bytes returns the bytes of the NullBitmap.
func (bmap *NullBitmap) Bytes() []byte {
	return bmap.bytes
}
