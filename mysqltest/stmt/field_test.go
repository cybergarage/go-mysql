// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sttm

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/cybergarage/go-mysql/mysql/query"
	"github.com/cybergarage/go-mysql/mysql/stmt"
)

func TestMySQLBinaryProtocolExamples(t *testing.T) {
	// MySQL: Binary Protocol Resultset
	// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html

	tests := []struct {
		typ       stmt.FieldType
		hexString string
		expected  any
	}{
		{query.MySQLTypeString, "666f6f", "foo"},
		{query.MySQLTypeLongLong, "0100000000000000", int64(1)},
		{query.MySQLTypeLong, "01000000", int32(1)},
		{query.MySQLTypeShort, "0100", int16(1)},
		{query.MySQLTypeTiny, "01", int8(1)},
		{query.MySQLTypeDouble, "6666666666662440", float64(10.2)},
		{query.MySQLTypeFloat, "33332341", float32(10.2)},
		// datetime 2010-10-17 19:27:30.000 001
		{query.MySQLTypeDatetime, "0bda070a11131b1e01000000", time.Date(2010, 10, 17, 19, 27, 30, 1*1000, time.UTC)},
		// date 2010-10-17
		{query.MySQLTypeDate, "04da070a11", time.Date(2010, 10, 17, 0, 0, 0, 0, time.UTC)},
		// timestamp 2010-10-17 19:27:30.000 001
		{query.MySQLTypeTimestamp, "0bda070a11131b1e01000000", time.Date(2010, 10, 17, 19, 27, 30, 1*1000, time.UTC)},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s(%s)", test.typ.String(), test.hexString), func(t *testing.T) {
			bytes, err := hex.DecodeString(test.hexString)
			if err != nil {
				t.Errorf("Failed to decode the hex string (%s)", test.hexString)
				return
			}
			field := stmt.NewField(
				stmt.WithFieldType(test.typ),
				stmt.WithFieldBytes(bytes),
			)
			v, err := field.Value()
			if err != nil {
				t.Error(err)
				return
			}
			if v != test.expected {
				t.Errorf("Expected %v, but got %v", test.expected, v)
			}
		})
	}
}
