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

// MySQL: include/field_types.h File Reference
// https://dev.mysql.com/doc/dev/mysql-server/latest/field__types_8h.html

// enum  	enum_field_types {
//   MYSQL_TYPE_DECIMAL , MYSQL_TYPE_TINY , MYSQL_TYPE_SHORT , MYSQL_TYPE_LONG ,
//   MYSQL_TYPE_FLOAT , MYSQL_TYPE_DOUBLE , MYSQL_TYPE_NULL , MYSQL_TYPE_TIMESTAMP ,
//   MYSQL_TYPE_LONGLONG , MYSQL_TYPE_INT24 , MYSQL_TYPE_DATE , MYSQL_TYPE_TIME ,
//   MYSQL_TYPE_DATETIME , MYSQL_TYPE_YEAR , MYSQL_TYPE_NEWDATE , MYSQL_TYPE_VARCHAR ,
//   MYSQL_TYPE_BIT , MYSQL_TYPE_TIMESTAMP2 , MYSQL_TYPE_DATETIME2 , MYSQL_TYPE_TIME2 ,
//   MYSQL_TYPE_TYPED_ARRAY , MYSQL_TYPE_VECTOR = 242 , MYSQL_TYPE_INVALID = 243 , MYSQL_TYPE_BOOL = 244 ,
//   MYSQL_TYPE_JSON = 245 , MYSQL_TYPE_NEWDECIMAL = 246 , MYSQL_TYPE_ENUM = 247 , MYSQL_TYPE_SET = 248 ,
//   MYSQL_TYPE_TINY_BLOB = 249 , MYSQL_TYPE_MEDIUM_BLOB = 250 , MYSQL_TYPE_LONG_BLOB = 251 , MYSQL_TYPE_BLOB = 252 ,
//   MYSQL_TYPE_VAR_STRING = 253 , MYSQL_TYPE_STRING = 254 , MYSQL_TYPE_GEOMETRY = 255
// }

// FieldType represents a MySQL field type.
type FieldType uint8

const (
	MYSQL_TYPE_DECIMAL FieldType = iota
	MYSQL_TYPE_TINY
	MYSQL_TYPE_SHORT
	MYSQL_TYPE_LONG
	MYSQL_TYPE_FLOAT
	MYSQL_TYPE_DOUBLE
	MYSQL_TYPE_NULL
	MYSQL_TYPE_TIMESTAMP
	MYSQL_TYPE_LONGLONG
	MYSQL_TYPE_INT24
	MYSQL_TYPE_DATE
	MYSQL_TYPE_TIME
	MYSQL_TYPE_DATETIME
	MYSQL_TYPE_YEAR
	MYSQL_TYPE_NEWDATE
	MYSQL_TYPE_VARCHAR
	MYSQL_TYPE_BIT
	MYSQL_TYPE_TIMESTAMP2
	MYSQL_TYPE_DATETIME2
	MYSQL_TYPE_TIME2
	MYSQL_TYPE_TYPED_ARRAY
	MYSQL_TYPE_VECTOR      = 242
	MYSQL_TYPE_INVALID     = 243
	MYSQL_TYPE_BOOL        = 244
	MYSQL_TYPE_JSON        = 245
	MYSQL_TYPE_NEWDECIMAL  = 246
	MYSQL_TYPE_ENUM        = 247
	MYSQL_TYPE_SET         = 248
	MYSQL_TYPE_TINY_BLOB   = 249
	MYSQL_TYPE_MEDIUM_BLOB = 250
	MYSQL_TYPE_LONG_BLOB   = 251
	MYSQL_TYPE_BLOB        = 252
	MYSQL_TYPE_VAR_STRING  = 253
	MYSQL_TYPE_STRING      = 254
	MYSQL_TYPE_GEOMETRY    = 255
)
