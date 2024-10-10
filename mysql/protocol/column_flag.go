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

// MySQL: Column Definition Flags
// https://dev.mysql.com/doc/dev/mysql-server/latest/group__group__cs__column__definition__flags.html

// ColumnDefFlag represents a MySQL column definition flag.
type ColumnDefFlag uint16

// ColumnDefFlag constants
const (
	// ColumnDefNotNULL represents the NOT_NULL column definition flag.
	ColumnDefNotNULL ColumnDefFlag = 1
	// ColumnDefPriKey represents the PRI_KEY column definition flag.
	ColumnDefPriKey ColumnDefFlag = 2
	// ColumnDefUniqueKey represents the UNIQUE_KEY column definition flag.
	ColumnDefUniqueKey ColumnDefFlag = 4
	// ColumnDefMultipleKey represents the MULTIPLE_KEY column definition flag.
	ColumnDefMultipleKey ColumnDefFlag = 8
	// ColumnDefBlob represents the BLOB column definition flag.
	ColumnDefBlob ColumnDefFlag = 16
	// ColumnDefUnsigned represents the UNSIGNED column definition flag.
	ColumnDefUnsigned ColumnDefFlag = 32
	// ColumnDefZeroFill represents the ZEROFILL column definition flag.
	ColumnDefZeroFill ColumnDefFlag = 64
	// ColumnDefBinary represents the BINARY column definition flag.
	ColumnDefBinary ColumnDefFlag = 128
	// ColumnDefEnum represents the ENUM column definition flag.
	ColumnDefEnum ColumnDefFlag = 256
	// ColumnDefAutoIncrement represents the AUTO_INCREMENT column definition flag.
	ColumnDefAutoIncrement ColumnDefFlag = 512
	// ColumnDefTimestamp represents the TIMESTAMP column definition flag.
	ColumnDefTimestamp ColumnDefFlag = 1024
	// ColumnDefSet represents the SET column definition flag.
	ColumnDefSet ColumnDefFlag = 2048
	// ColumnDefNoDefaultValue represents the NO_DEFAULT_VALUE column definition flag.
	ColumnDefNoDefaultValue ColumnDefFlag = 4096
	// ColumnDefOnUpdateNow represents the ON_UPDATE_NOW column definition flag.
	ColumnDefOnUpdateNow ColumnDefFlag = 8192
	// ColumnDefPartKey represents the PART_KEY column definition flag.
	ColumnDefPartKey ColumnDefFlag = 16384
	// ColumnDefNum represents the NUM_FLAG column definition flag.
	ColumnDefNum ColumnDefFlag = 32768
)
