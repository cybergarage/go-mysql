// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
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

// MySQL: Server Status Flag
// https://dev.mysql.com/doc/dev/mysql-server/latest/mysql__com_8h.html

type CursorType uint8

const (
	// CursorTypeNoCursor indicates that no cursor is used.
	CursorTypeNoCursor CursorType = 0
	// CursorTypeReadOnly indicates that the cursor is read-only.
	CursorTypeReadOnly CursorType = 1
	// CursorTypeForUpdate indicates that the cursor is for update.
	CursorTypeForUpdate CursorType = 2
	// CursorTypeScrollable indicates that the cursor is scrollable.
	CursorTypeScrollable CursorType = 4
)

// IsEnabled returns true if the status flag is enabled.
func (statFlag CursorType) IsEnabled(flag CursorType) bool {
	return (statFlag & flag) != 0
}
