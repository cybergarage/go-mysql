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
	defaultDatetimeBytesLen = 11
)

// BytesToTime converts a byte slice to a time.Time.
func BytesToTime(b []byte) (time.Time, error) {
	if len(b) < 1 {
		return time.Time{}, newErrInvalidDatetimeBytes(b)
	}
	l := int(b[0])
	var year, month, day, hour, minute, second, microsecond int
	switch l {
	case 0:
		// 0 for special '0000-00-00 00:00:00' value.
		return time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC), nil
	case 4:
		year = int(b[1]) | int(b[2])<<8
		month = int(b[3])
		day = int(b[4])
		hour = 0
		minute = 0
		second = 0
		microsecond = 0
	case 7:
		year = int(b[1]) | int(b[2])<<8
		month = int(b[3])
		day = int(b[4])
		hour = int(b[5])
		minute = int(b[6])
		second = int(b[7])
		microsecond = 0
	case 11:
		year = int(b[1]) | int(b[2])<<8
		month = int(b[3])
		day = int(b[4])
		hour = int(b[5])
		minute = int(b[6])
		second = int(b[7])
		microsecond = int(b[8]) | int(b[9])<<8 | int(b[10])<<16 | int(b[11])<<24
	default:
		return time.Time{}, newErrInvalidDatetimeBytes(b)
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, microsecond*1000, time.UTC), nil
}

// BytesToDatetime converts a datetime byte slice to a time.Time.
func BytesToDatetime(b []byte) (time.Time, error) {
	return BytesToTime(b)
}

// TimeToDatetimeBytes converts a time.Time to a datetime byte slice.
func TimeToDatetimeBytes(t time.Time) []byte {
	year := t.Year()
	microsecond := t.Nanosecond() / 1000
	b := make([]byte, defaultDatetimeBytesLen+1)
	b[0] = byte(defaultDatetimeBytesLen)
	b[1] = byte(year & 0xFF)
	b[2] = byte((year >> 8) & 0xFF)
	b[3] = byte(t.Month())
	b[4] = byte(t.Day())
	b[5] = byte(t.Hour())
	b[6] = byte(t.Minute())
	b[7] = byte(t.Second())
	b[8] = byte(microsecond & 0xFF)
	b[9] = byte((microsecond >> 8) & 0xFF)
	b[10] = byte((microsecond >> 16) & 0xFF)
	b[11] = byte((microsecond >> 24) & 0xFF)
	return b
}
