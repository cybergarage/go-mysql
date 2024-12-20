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

// MySQL: Status Flag
// https://dev.mysql.com/doc/dev/mysql-server/latest/mysql__com_8h.html

// StatusFlag represents a MySQL Status Flag.
type StatusFlag uint16

const (
	StatusInTrans              StatusFlag = 1
	StatusAutoCommit           StatusFlag = 2
	StatusMoreResultsExists    StatusFlag = 8
	StatusQueryNoGoodIndexUsed StatusFlag = 16
	StatusQueryNoIndexUsed     StatusFlag = 32
	StatusCursorExists         StatusFlag = 64
	StatusLastRowSent          StatusFlag = 128
	StatusDBDropped            StatusFlag = 256
	StatusNoBackslashEscapes   StatusFlag = 512
	StatusMetadataChanged      StatusFlag = 1024
	StatusQueryWasSlow         StatusFlag = 2048
	StatusPsOutParams          StatusFlag = 4096
	StatusInTransReadOnly      StatusFlag = 8192
	StatusSessionStateChanged  StatusFlag = 1 << 14
)

// IsEnabled returns true if the status flag is enabled.
func (statFlag StatusFlag) IsEnabled(flag StatusFlag) bool {
	return (statFlag & flag) != 0
}
