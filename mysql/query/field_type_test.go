// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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

package query

import (
	"testing"
)

func TestFieldTypes(t *testing.T) {
	// MySQL: include/field_types.h File Reference
	// https://dev.mysql.com/doc/dev/mysql-server/latest/field__types_8h.html
	// Type of the parameter value. See enum_field_type
	// https://github.com/mysql/mysql-server/blob/trunk/include/field_types.h
	tests := []struct {
		fieldType FieldType
		fieldEnum uint8
		expected  string
	}{
		{MySQLTypeDecimal, 0, "MYSQL_TYPE_DECIMAL"},
		{MySQLTypeTiny, 1, "MYSQL_TYPE_TINYINT"},
		{MySQLTypeShort, 2, "MYSQL_TYPE_SMALLINT"},
		{MySQLTypeLong, 3, "MYSQL_TYPE_INT"},
		{MySQLTypeFloat, 4, "MYSQL_TYPE_FLOAT"},
		{MySQLTypeDouble, 5, "MYSQL_TYPE_DOUBLE"},
		{MySQLTypeNull, 6, "MYSQL_TYPE_NULL"},
		{MySQLTypeTimestamp, 7, "MYSQL_TYPE_TIMESTAMP"},
		{MySQLTypeLongLong, 8, "MYSQL_TYPE_BIGINT"},
		{MySQLTypeInt24, 9, "MYSQL_TYPE_MEDIUMINT"},
		{MySQLTypeDate, 10, "MYSQL_TYPE_DATE"},
		{MySQLTypeTime, 11, "MYSQL_TYPE_TIME"},
		{MySQLTypeDatetime, 12, "MYSQL_TYPE_DATETIME"},
		{MySQLTypeYear, 13, "MYSQL_TYPE_YEAR"},
		{MySQLTypeNewdate, 14, "MYSQL_TYPE_NEWDATE"},
		{MySQLTypeVarchar, 15, "MYSQL_TYPE_VARCHAR"},
		{MySQLTypeBit, 16, "MYSQL_TYPE_BIT"},
		{MySQLTypeTimestamp2, 17, "MYSQL_TYPE_TIMESTAMP2"},
		{MySQLTypeDatetime2, 18, "MYSQL_TYPE_DATETIME2"},
		{MySQLTypeTime2, 19, "MYSQL_TYPE_TIME2"},
		{MySQLTypeTypedArray, 20, "MYSQL_TYPE_TYPED_ARRAY"},
		{MySQLTypeVector, 242, "MYSQL_TYPE_VECTOR"},
		{MySQLTypeInvalid, 243, "MYSQL_TYPE_INVALID"},
		{MySQLTypeBool, 244, "MYSQL_TYPE_BOOL"},
		{MySQLTypeJSON, 245, "MYSQL_TYPE_JSON"},
		{MySQLTypeNewdecimal, 246, "MYSQL_TYPE_NEWDECIMAL"},
		{MySQLTypeEnum, 247, "MYSQL_TYPE_ENUM"},
		{MySQLTypeSet, 248, "MYSQL_TYPE_SET"},
		{MySQLTypeTinyBlob, 249, "MYSQL_TYPE_TINYBLOB"},
		{MySQLTypeMediumBlob, 250, "MYSQL_TYPE_MEDIUMBLOB"},
		{MySQLTypeLongBlob, 251, "MYSQL_TYPE_LONGBLOB"},
		{MySQLTypeBlob, 252, "MYSQL_TYPE_BLOB"},
		{MySQLTypeVarString, 253, "MYSQL_TYPE_VAR_STRING"},
		{MySQLTypeString, 254, "MYSQL_TYPE_STRING"},
		{MySQLTypeGeometry, 255, "MYSQL_TYPE_GEOMETRY"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			if got := test.fieldType; got != FieldType(test.fieldEnum) {
				t.Errorf("FieldType = %v, want %v", got, test.fieldEnum)
			}

			if got := test.fieldType.String(); got != test.expected {
				t.Errorf("FieldType.String() = %v, want %v", got, test.expected)
			}
		})
	}
}
