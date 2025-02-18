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

// FieldType represents a field type.
type FieldType uint8

// MySQL: include/field_types.h File Reference
// https://dev.mysql.com/doc/dev/mysql-server/latest/field__types_8h.html
// Type of the parameter value. See enum_field_type
// https://github.com/mysql/mysql-server/blob/trunk/include/field_types.h

const (
	// MySQLTypeDecimal represents a MYSQL_TYPE_DECIMAL type.
	MySQLTypeDecimal FieldType = iota
	// MySQLTypeTiny represents a MYSQL_TYPE_TINYINT type.
	MySQLTypeTiny
	// MySQLTypeShort represents a MYSQL_TYPE_SMALLINT type.
	MySQLTypeShort
	// MySQLTypeLong represents an INT type.
	MySQLTypeLong
	// MySQLTypeFloat represents a MYSQL_TYPE_FLOAT type.
	MySQLTypeFloat
	// MySQLTypeDouble represents a MYSQL_TYPE_DOUBLE type.
	MySQLTypeDouble
	// MySQLTypeNull represents a MYSQL_TYPE_NULL type.
	MySQLTypeNull
	// MySQLTypeTimestamp represents a MYSQL_TYPE_TIMESTAMP type.
	MySQLTypeTimestamp
	// MySQLTypeLonglong represents a MYSQL_TYPE_BIGINT type.
	MySQLTypeLonglong
	// MySQLTypeInt24 represents a MYSQL_TYPE_MEDIUMINT type.
	MySQLTypeInt24
	// MySQLTypeDate represents a MYSQL_TYPE_DATE type.
	MySQLTypeDate
	// MySQLTypeTime represents a MYSQL_TYPE_TIME type.
	MySQLTypeTime
	// MySQLTypeDatetime represents a MYSQL_TYPE_DATETIME type.
	MySQLTypeDatetime
	// MySQLTypeYear represents a MYSQL_TYPE_YEAR type.
	MySQLTypeYear
	// MySQLTypeNewdate represents a MYSQL_TYPE_NEWDATE type.
	MySQLTypeNewdate
	// MySQLTypeVarchar represents a MYSQL_TYPE_VARCHAR type.
	MySQLTypeVarchar
	// MySQLTypeBit represents a MYSQL_TYPE_BIT type.
	MySQLTypeBit
	// MySQLTypeTimestamp2 represents a MYSQL_TYPE_TIMESTAMP2 type.
	MySQLTypeTimestamp2
	// MySQLTypeDatetime2 represents a MYSQL_TYPE_DATETIME2 type.
	MySQLTypeDatetime2
	// MySQLTypeTime2 represents a MYSQL_TYPE_TIME2 type.
	MySQLTypeTime2
	// MySQLTypeTypedArray represents a MYSQL_TYPE_TYPED_ARRAY type.
	MySQLTypeTypedArray
	// MySQLTypeVector represents a MYSQL_TYPE_VECTOR type.
	MySQLTypeVector = 242
	// MySQLTypeInvalid represents an INVALID type.
	MySQLTypeInvalid = 243
	// MySQLTypeBool represents a MYSQL_TYPE_BOOL type.
	MySQLTypeBool = 244
	// MySQLTypeJSON represents a MYSQL_TYPE_JSON type.
	MySQLTypeJSON = 245
	// MySQLTypeNewdecimal represents a MYSQL_TYPE_NEWDECIMAL type.
	MySQLTypeNewdecimal = 246
	// MySQLTypeEnum represents an ENUM type.
	MySQLTypeEnum = 247
	// MySQLTypeSet represents a MYSQL_TYPE_SET type.
	MySQLTypeSet = 248
	// MySQLTypeTinyBlob represents a MYSQL_TYPE_TINYBLOB type.
	MySQLTypeTinyBlob = 249
	// MySQLTypeMediumBlob represents a MYSQL_TYPE_MEDIUMBLOB type.
	MySQLTypeMediumBlob = 250
	// MySQLTypeLongBlob represents a MYSQL_TYPE_LONGBLOB type.
	MySQLTypeLongBlob = 251
	// MySQLTypeBlob represents a MYSQL_TYPE_BLOB type.
	MySQLTypeBlob = 252
	// MySQLTypeVarString represents a MYSQL_TYPE_VAR_STRING type.
	MySQLTypeVarString = 253
	// MySQLTypeString represents a MYSQL_TYPE_STRING type.
	MySQLTypeString = 254
	// MySQLTypeGeometry represents a MYSQL_TYPE_GEOMETRY type.
	MySQLTypeGeometry = 255
)

// String returns a string representation of the field type.
func (ft FieldType) String() string {
	switch ft {
	case MySQLTypeDecimal:
		return "MYSQL_TYPE_DECIMAL"
	case MySQLTypeTiny:
		return "MYSQL_TYPE_TINYINT"
	case MySQLTypeShort:
		return "MYSQL_TYPE_SMALLINT"
	case MySQLTypeLong:
		return "MYSQL_TYPE_INT"
	case MySQLTypeFloat:
		return "MYSQL_TYPE_FLOAT"
	case MySQLTypeDouble:
		return "MYSQL_TYPE_DOUBLE"
	case MySQLTypeNull:
		return "MYSQL_TYPE_NULL"
	case MySQLTypeTimestamp:
		return "MYSQL_TYPE_TIMESTAMP"
	case MySQLTypeLonglong:
		return "MYSQL_TYPE_BIGINT"
	case MySQLTypeInt24:
		return "MYSQL_TYPE_MEDIUMINT"
	case MySQLTypeDate:
		return "MYSQL_TYPE_DATE"
	case MySQLTypeTime:
		return "MYSQL_TYPE_TIME"
	case MySQLTypeDatetime:
		return "MYSQL_TYPE_DATETIME"
	case MySQLTypeYear:
		return "MYSQL_TYPE_YEAR"
	case MySQLTypeNewdate:
		return "MYSQL_TYPE_NEWDATE"
	case MySQLTypeVarchar:
		return "MYSQL_TYPE_VARCHAR"
	case MySQLTypeBit:
		return "MYSQL_TYPE_BIT"
	case MySQLTypeTimestamp2:
		return "MYSQL_TYPE_TIMESTAMP2"
	case MySQLTypeDatetime2:
		return "MYSQL_TYPE_DATETIME2"
	case MySQLTypeTime2:
		return "MYSQL_TYPE_TIME2"
	case MySQLTypeTypedArray:
		return "MYSQL_TYPE_TYPED_ARRAY"
	case MySQLTypeVector:
		return "MYSQL_TYPE_VECTOR"
	case MySQLTypeInvalid:
		return "MYSQL_TYPE_INVALID"
	case MySQLTypeBool:
		return "MYSQL_TYPE_BOOL"
	case MySQLTypeJSON:
		return "MYSQL_TYPE_JSON"
	case MySQLTypeNewdecimal:
		return "MYSQL_TYPE_NEWDECIMAL"
	case MySQLTypeEnum:
		return "MYSQL_TYPE_ENUM"
	case MySQLTypeSet:
		return "MYSQL_TYPE_SET"
	case MySQLTypeTinyBlob:
		return "MYSQL_TYPE_TINYBLOB"
	case MySQLTypeMediumBlob:
		return "MYSQL_TYPE_MEDIUMBLOB"
	case MySQLTypeLongBlob:
		return "MYSQL_TYPE_LONGBLOB"
	case MySQLTypeBlob:
		return "MYSQL_TYPE_BLOB"
	case MySQLTypeVarString:
		return "MYSQL_TYPE_VAR_STRING"
	case MySQLTypeString:
		return "MYSQL_TYPE_STRING"
	case MySQLTypeGeometry:
		return "MYSQL_TYPE_GEOMETRY"
	default:
		return "UNKNOWN"
	}
}
