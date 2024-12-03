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
	case query.BigIntData:
		return MySQLTypeLonglong, nil
	case query.BinaryData:
		return MySQLTypeBlob, nil
	case query.BitData:
		return MySQLTypeBit, nil
	case query.BlobData:
		return MySQLTypeBlob, nil
	case query.BooleanData:
		return MySQLTypeBool, nil
	// case query.CharData, query.CharacterData:
	// 	return MySQLTypeChar, nil
	case query.ClobData:
		return MySQLTypeString, nil
	case query.DateData:
		return MySQLTypeDate, nil
	case query.DecimalData:
		return MySQLTypeDecimal, nil
	case query.DoubleData, query.DoublePrecision:
		return MySQLTypeDouble, nil
	case query.FloatData:
		return MySQLTypeFloat, nil
	case query.IntData, query.IntegerData:
		return MySQLTypeLong, nil
	case query.LongBlobData:
		return MySQLTypeLongBlob, nil
	case query.LongTextData:
		return MySQLTypeString, nil
	case query.MediumBlobData:
		return MySQLTypeMediumBlob, nil
	case query.MediumIntData:
		return MySQLTypeInt24, nil
	case query.MediumTextData:
		return MySQLTypeString, nil
	// case query.NumericData:
	// 	return MySQLTypeNumeric
	// case query.RealData:
	// 	return MySQLTypeReal
	case query.SetData:
		return MySQLTypeSet, nil
	case query.SmallIntData:
		return MySQLTypeTiny, nil
	case query.TextData:
		return MySQLTypeString, nil
	case query.TimeData:
		return MySQLTypeTime, nil
	case query.DateTimeData:
		return MySQLTypeDatetime, nil
	case query.TimeStampData:
		return MySQLTypeTimestamp, nil
	case query.TinyBlobData:
		return MySQLTypeTinyBlob, nil
	case query.TinyIntData:
		return MySQLTypeTiny, nil
	case query.TinyTextData:
		return MySQLTypeString, nil
	case query.VarBinaryData:
		return MySQLTypeBlob, nil
	case query.VarCharData, query.VarCharacterData:
		return MySQLTypeVarchar, nil
	case query.YearData:
		return MySQLTypeYear, nil
	}

	return 0, errors.NewErrUnsupported(t.String())
}
