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

package query

// MySQL: include/field_types.h File Reference
// https://dev.mysql.com/doc/dev/mysql-server/latest/field__types_8h.html

// FieldType represents a field type.
type FieldType uint8

const (
	// MysqlTypeDecimal represents a MYSQL_TYPE_DECIMAL type.
	MysqlTypeDecimal FieldType = iota
	// MysqlTypeTiny represents a MYSQL_TYPE_TINYINT type.
	MysqlTypeTiny
	// MysqlTypeShort represents a MYSQL_TYPE_SMALLINT type.
	MysqlTypeShort
	// MysqlTypeLong represents an INT type.
	MysqlTypeLong
	// MysqlTypeFloat represents a MYSQL_TYPE_FLOAT type.
	MysqlTypeFloat
	// MysqlTypeDouble represents a MYSQL_TYPE_DOUBLE type.
	MysqlTypeDouble
	// MysqlTypeNull represents a MYSQL_TYPE_NULL type.
	MysqlTypeNull
	// MysqlTypeTimestamp represents a MYSQL_TYPE_TIMESTAMP type.
	MysqlTypeTimestamp
	// MysqlTypeLonglong represents a MYSQL_TYPE_BIGINT type.
	MysqlTypeLonglong
	// MysqlTypeInt24 represents a MYSQL_TYPE_MEDIUMINT type.
	MysqlTypeInt24
	// MysqlTypeDate represents a MYSQL_TYPE_DATE type.
	MysqlTypeDate
	// MysqlTypeTime represents a MYSQL_TYPE_TIME type.
	MysqlTypeTime
	// MysqlTypeDatetime represents a MYSQL_TYPE_DATETIME type.
	MysqlTypeDatetime
	// MysqlTypeYear represents a MYSQL_TYPE_YEAR type.
	MysqlTypeYear
	// MysqlTypeNewdate represents a MYSQL_TYPE_NEWDATE type.
	MysqlTypeNewdate
	// MysqlTypeVarchar represents a MYSQL_TYPE_VARCHAR type.
	MysqlTypeVarchar
	// MysqlTypeBit represents a MYSQL_TYPE_BIT type.
	MysqlTypeBit
	// MysqlTypeTimestamp2 represents a MYSQL_TYPE_TIMESTAMP2 type.
	MysqlTypeTimestamp2
	// MysqlTypeDatetime2 represents a MYSQL_TYPE_DATETIME2 type.
	MysqlTypeDatetime2
	// MysqlTypeTime2 represents a MYSQL_TYPE_TIME2 type.
	MysqlTypeTime2
	// MysqlTypeTypedArray represents a MYSQL_TYPE_TYPED_ARRAY type.
	MysqlTypeTypedArray
	// MysqlTypeVector represents a MYSQL_TYPE_VECTOR type.
	MysqlTypeVector = 242
	// MysqlTypeInvalid represents an INVALID type.
	MysqlTypeInvalid = 243
	// MysqlTypeBool represents a MYSQL_TYPE_BOOL type.
	MysqlTypeBool = 244
	// MysqlTypeJSON represents a MYSQL_TYPE_JSON type.
	MysqlTypeJSON = 245
	// MysqlTypeNewdecimal represents a MYSQL_TYPE_NEWDECIMAL type.
	MysqlTypeNewdecimal = 246
	// MysqlTypeEnum represents an ENUM type.
	MysqlTypeEnum = 247
	// MysqlTypeSet represents a MYSQL_TYPE_SET type.
	MysqlTypeSet = 248
	// MysqlTypeTinyBlob represents a MYSQL_TYPE_TINYBLOB type.
	MysqlTypeTinyBlob = 249
	// MysqlTypeMediumBlob represents a MYSQL_TYPE_MEDIUMBLOB type.
	MysqlTypeMediumBlob = 250
	// MysqlTypeLongBlob represents a MYSQL_TYPE_LONGBLOB type.
	MysqlTypeLongBlob = 251
	// MysqlTypeBlob represents a MYSQL_TYPE_BLOB type.
	MysqlTypeBlob = 252
	// MysqlTypeVarString represents a MYSQL_TYPE_VAR_STRING type.
	MysqlTypeVarString = 253
	// MysqlTypeString represents a MYSQL_TYPE_STRING type.
	MysqlTypeString = 254
	// MysqlTypeGeometry represents a MYSQL_TYPE_GEOMETRY type.
	MysqlTypeGeometry = 255
)
