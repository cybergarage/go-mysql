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

package binary

import (
	"time"
)

// MySQL: Binary Protocol Resultset
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_binary_resultset.html
// Resultset row - MariaDB Knowledge Base
// https://mariadb.com/kb/en/resultset-row/#timestamp-binary-encoding

const (
	dateBytesLen = 4
)

// BytesToDate converts a date byte slice to a time.Time.
func BytesToDate(b []byte) (time.Time, error) {
	return BytesToTime(b)
}

// TimeToDatetimeBytes converts a time.Time to a datetime byte slice.
func TimeToDateBytes(t time.Time) []byte {
	year := t.Year()
	b := make([]byte, dateBytesLen+1)
	b[0] = byte(dateBytesLen)
	b[1] = byte(year & 0xFF)
	b[2] = byte((year >> 8) & 0xFF)
	b[3] = byte(t.Month())
	b[4] = byte(t.Day())
	return b
}
