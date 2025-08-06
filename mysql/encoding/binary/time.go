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
	defaultTimeBytesLen = 12
)

// BytesToTime converts a byte slice to a time.Time.
func BytesToDuration(b []byte) (time.Duration, error) {
	if len(b) < 1 {
		return time.Duration(0), newErrInvalidTimeBytes(b)
	}

	l := int(b[0])
	if len(b) < (l + 1) {
		return time.Duration(0), newErrInvalidTimeBytes(b)
	}

	var isNegative, days, hour, minute, second, microsecond int

	switch l {
	case 0:
		return time.Duration(0), nil
	case 8:
		isNegative = int(b[1])
		days = int(b[2]) | int(b[3])<<8 | int(b[4])<<16 | int(b[5])<<24
		hour = int(b[6])
		minute = int(b[7])
		second = int(b[8])
	case 12:
		isNegative = int(b[1])
		days = int(b[2]) | int(b[3])<<8 | int(b[4])<<16 | int(b[5])<<24
		hour = int(b[6])
		minute = int(b[7])
		second = int(b[8])
		microsecond = int(b[9]) | int(b[10])<<8 | int(b[11])<<16 | int(b[12])<<24
	default:
		return time.Duration(0), newErrInvalidTimeBytes(b)
	}

	d := time.Duration(days)*24*time.Hour + time.Duration(hour)*time.Hour + time.Duration(minute)*time.Minute + time.Duration(second)*time.Second + time.Duration(microsecond)*time.Microsecond

	if isNegative != 0 {
		d = -d
	}

	return d, nil
}

// TimeToDatetimeBytes converts a time.Time to a datetime byte slice.
func DurationToTimeBytes(d time.Duration) []byte {
	var (
		isNegative                              byte
		days, hour, minute, second, microsecond int
	)

	if d < 0 {
		isNegative = 1
		d = -d
	}

	days = int(d / (24 * time.Hour))
	d -= time.Duration(days) * 24 * time.Hour
	hour = int(d / time.Hour)
	d -= time.Duration(hour) * time.Hour
	minute = int(d / time.Minute)
	d -= time.Duration(minute) * time.Minute
	second = int(d / time.Second)
	d -= time.Duration(second) * time.Second
	microsecond = int(d / 1000)
	b := make([]byte, defaultTimeBytesLen+1)
	b[0] = byte(defaultTimeBytesLen)
	b[1] = byte(isNegative)
	b[2] = byte(days & 0xFF)
	b[3] = byte((days >> 8) & 0xFF)
	b[4] = byte((days >> 16) & 0xFF)
	b[5] = byte((days >> 24) & 0xFF)
	b[6] = byte(hour)
	b[7] = byte(minute)
	b[8] = byte(second)
	b[9] = byte(microsecond & 0xFF)
	b[10] = byte((microsecond >> 8) & 0xFF)
	b[11] = byte((microsecond >> 16) & 0xFF)
	b[12] = byte((microsecond >> 24) & 0xFF)

	return b
}
