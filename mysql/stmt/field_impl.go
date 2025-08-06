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

package stmt

// MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// Result Set Packets - MariaDB Knowledge Base
// https://mariadb.com/kb/en/result-set-packets/

import (
	"time"

	"github.com/cybergarage/go-mysql/mysql/encoding/binary"
	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-safecast/safecast"
)

type field struct {
	t FieldType
	b []byte
	v any
}

// Field represents a field option.
type FieldOption func(*field)

// WithFieldType returns a field option to set the field type.
func WithFieldType(t FieldType) FieldOption {
	return func(f *field) {
		f.t = t
	}
}

// WithFieldBytes returns a field option to set the field bytes.
func WithFieldBytes(b []byte) FieldOption {
	return func(f *field) {
		f.b = b
	}
}

// WithFieldValue returns a field option to set the field value.
func WithFieldValue(v any) FieldOption {
	return func(f *field) {
		f.v = v
	}
}

// NewField creates a new field instance.
func NewField(opts ...FieldOption) Field {
	f := &field{
		t: FieldType(0),
		b: nil,
		v: nil,
	}
	for _, opt := range opts {
		opt(f)
	}
	return f
}

// Type returns the field type.
func (f *field) Type() FieldType {
	return f.t
}

// Bytes returns the field bytes.
func (f *field) Bytes() ([]byte, error) {
	if f.b != nil {
		return f.b, nil
	}

	switch f.t {
	case query.MySQLTypeString, query.MySQLTypeVarString, query.MySQLTypeVarchar:
		switch v := f.v.(type) {
		case string:
			f.b = []byte(v)
		case *string:
			if v == nil {
				f.b = []byte{binary.NullString}
			} else {
				f.b = []byte(*v)
			}
		case nil:
			f.b = []byte{binary.NullString}
		default:
			return nil, newErrInvalidField(f.t, f.v)
		}
	case query.MySQLTypeTinyBlob, query.MySQLTypeMediumBlob, query.MySQLTypeLongBlob, query.MySQLTypeBlob:
		b, ok := f.v.([]byte)
		if !ok {
			return nil, newErrInvalidField(f.t, f.v)
		}
		f.b = b
	case query.MySQLTypeNull:
		f.b = nil
	case query.MySQLTypeTiny:
		var cv int8
		err := safecast.ToInt8(f.v, &cv)
		if err != nil {
			return nil, err
		}
		f.b = binary.Int1ToBytes(cv)
	case query.MySQLTypeShort, query.MySQLTypeYear:
		var cv int16
		err := safecast.ToInt16(f.v, &cv)
		if err != nil {
			return nil, err
		}
		f.b = binary.Int2ToBytes(cv)
	case query.MySQLTypeLong, query.MySQLTypeInt24:
		var cv int32
		err := safecast.ToInt32(f.v, &cv)
		if err != nil {
			return nil, err
		}
		f.b = binary.Int4ToBytes(cv)
	case query.MySQLTypeLongLong:
		var cv int64
		err := safecast.ToInt64(f.v, &cv)
		if err != nil {
			return nil, err
		}
		f.b = binary.Int8ToBytes(cv)
	case query.MySQLTypeFloat:
		var cv float32
		err := safecast.ToFloat32(f.v, &cv)
		if err != nil {
			return nil, err
		}
		f.b = binary.Float4ToBytes(cv)
	case query.MySQLTypeDouble:
		var cv float64
		err := safecast.ToFloat64(f.v, &cv)
		if err != nil {
			return nil, err
		}
		f.b = binary.Float8ToBytes(cv)
	case query.MySQLTypeDatetime, query.MySQLTypeTimestamp:
		var cv time.Time
		err := safecast.ToTime(f.v, &cv)
		if err != nil {
			return nil, err
		}
		f.b = binary.TimeToDatetimeBytes(cv)
	case query.MySQLTypeDate:
		var cv time.Time
		err := safecast.ToTime(f.v, &cv)
		if err != nil {
			return nil, err
		}
		f.b = binary.TimeToDateBytes(cv)
	case query.MySQLTypeTime:
		cv, ok := f.v.(time.Duration)
		if !ok {
			return nil, newErrInvalidField(f.t, f.v)
		}
		f.b = binary.DurationToTimeBytes(cv)
	default:
		return nil, newErrNotSupportedFieldType(f.t)
	}

	return f.b, nil
}

// Value returns the field value.
func (f *field) Value() (any, error) {
	if f.v != nil || f.t == query.MySQLTypeNull {
		return f.v, nil
	}
	var err error
	if f.v == nil && (0 < len(f.b)) {
		switch f.t {
		case query.MySQLTypeTiny:
			f.v, err = binary.BytesToInt1(f.b)
		case query.MySQLTypeShort:
			f.v, err = binary.BytesToInt2(f.b)
		case query.MySQLTypeLong:
			f.v, err = binary.BytesToInt4(f.b)
		case query.MySQLTypeLongLong:
			f.v, err = binary.BytesToInt8(f.b)
		case query.MySQLTypeFloat:
			f.v, err = binary.BytesToFloat4(f.b)
		case query.MySQLTypeDouble:
			f.v, err = binary.BytesToFloat8(f.b)
		case query.MySQLTypeString, query.MySQLTypeVarString, query.MySQLTypeVarchar:
			f.v = string(f.b)
		case query.MySQLTypeTinyBlob, query.MySQLTypeMediumBlob, query.MySQLTypeLongBlob, query.MySQLTypeBlob:
			f.v = f.b
		case query.MySQLTypeDatetime, query.MySQLTypeTimestamp:
			f.v, err = binary.BytesToDatetime(f.b)
		case query.MySQLTypeDate:
			f.v, err = binary.BytesToDate(f.b)
		case query.MySQLTypeTime:
			f.v, err = binary.BytesToDuration(f.b)
		default:
			return nil, newErrNotSupportedFieldType(f.t)
		}
	}
	return f.v, err
}
