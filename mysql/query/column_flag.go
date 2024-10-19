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

// MySQL: Column Definition Flags
// https://dev.mysql.com/doc/dev/mysql-server/latest/group__group__cs__column__definition__flags.html

// ColumnDefFlag represents a MySQL column definition flag.
type ColumnDefFlag = uint16

const (
	NotNullFlag        ColumnDefFlag = 1
	PriKeyFlag         ColumnDefFlag = 2
	UniqueKeyFlag      ColumnDefFlag = 4
	MultipleKeyFlag    ColumnDefFlag = 8
	BlobFlag           ColumnDefFlag = 16
	UnsignedFlag       ColumnDefFlag = 32
	ZerofillFlag       ColumnDefFlag = 64
	BinaryFlag         ColumnDefFlag = 128
	EnumFlag           ColumnDefFlag = 256
	AutoIncrementFlag  ColumnDefFlag = 512
	TimestampFlag      ColumnDefFlag = 1024
	SetFlag            ColumnDefFlag = 2048
	NoDefaultValueFlag ColumnDefFlag = 4096
	OnUpdateNowFlag    ColumnDefFlag = 8192
	NumFlag            ColumnDefFlag = 32768
	PartKeyFlag        ColumnDefFlag = 16384
	GroupFlag          ColumnDefFlag = 32768
	// UniqueFlag                 ColumnDefFlag = 65536
	// BincmpFlag                 ColumnDefFlag = 131072
	// GetFixedFieldsFlag         ColumnDefFlag = 1 << 18
	// FieldInPartFuncFlag        ColumnDefFlag = 1 << 19
	// FieldInAddIndex            ColumnDefFlag = 1 << 20
	// FieldIsRenamed             ColumnDefFlag = 1 << 21
	// FieldFlagsStorageMedia     ColumnDefFlag = 22
	// FieldFlagsStorageMediaMask ColumnDefFlag = 3 << FieldFlagsStorageMedia
	// FieldFlagsColumnFormat     ColumnDefFlag = 24
	// FieldFlagsColumnFormatMask ColumnDefFlag = 3 << FieldFlagsColumnFormat
	// FieldIsDropped             ColumnDefFlag = 1 << 26
	// ExplicitNullFlag           ColumnDefFlag = 1 << 27
	// NotSecondaryFlag           ColumnDefFlag = 1 << 29
	// FieldIsInvisible           ColumnDefFlag = 1 << 30.
)
