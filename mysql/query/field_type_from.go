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

import (
	"github.com/cybergarage/go-mysql/mysql/errors"
	"github.com/cybergarage/go-sqlparser/sql/query"
)

// MySQL: include/field_types.h File Reference
// https://dev.mysql.com/doc/dev/mysql-server/latest/field__types_8h.html

// NewFieldTypeFrom returns a new FieldType from the data type.
func NewFieldTypeFrom(t query.DataType) (FieldType, error) {
	switch t {
	case query.BigIntType:
		return MySQLTypeLongLong, nil
	case query.BinaryType:
		return MySQLTypeBlob, nil
	case query.BitType:
		return MySQLTypeBit, nil
	case query.BlobType:
		return MySQLTypeBlob, nil
	case query.BooleanType:
		return MySQLTypeBool, nil
	// case query.CharData, query.CharacterData:
	// 	return MySQLTypeChar, nil
	case query.ClobType:
		return MySQLTypeString, nil
	case query.DateType:
		return MySQLTypeDate, nil
	case query.DecimalType:
		return MySQLTypeDecimal, nil
	case query.DoubleType, query.DoublePrecisionType:
		return MySQLTypeDouble, nil
	case query.FloatType:
		return MySQLTypeFloat, nil
	case query.IntType, query.IntegerType:
		return MySQLTypeLong, nil
	case query.LongBlobType:
		return MySQLTypeLongBlob, nil
	case query.LongTextType:
		return MySQLTypeString, nil
	case query.MediumBlobType:
		return MySQLTypeMediumBlob, nil
	case query.MediumIntType:
		return MySQLTypeInt24, nil
	case query.MediumTextType:
		return MySQLTypeString, nil
	// case query.NumericData:
	// 	return MySQLTypeNumeric
	// case query.RealData:
	// 	return MySQLTypeReal
	case query.SetType:
		return MySQLTypeSet, nil
	case query.SmallIntType:
		return MySQLTypeTiny, nil
	case query.TextType:
		return MySQLTypeString, nil
	case query.TimeType:
		return MySQLTypeTime, nil
	case query.DateTimeType:
		return MySQLTypeDatetime, nil
	case query.TimeStampType:
		return MySQLTypeTimestamp, nil
	case query.TinyBlobType:
		return MySQLTypeTinyBlob, nil
	case query.TinyIntType:
		return MySQLTypeTiny, nil
	case query.TinyTextType:
		return MySQLTypeString, nil
	case query.VarBinaryType:
		return MySQLTypeBlob, nil
	case query.VarCharType, query.VarCharacterType:
		return MySQLTypeVarchar, nil
	case query.YearType:
		return MySQLTypeYear, nil
	}

	return 0, errors.NewErrUnsupported(t.String())
}
