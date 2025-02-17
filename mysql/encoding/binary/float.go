// Copyright (C) 2019 The go-mysql Authors. All rights reserved.
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
	"math"
)

// BytesToFloat8 converts the specified byte array to a float64.
func BytesToFloat8(b []byte) (float64, error) {
	if len(b) != 8 {
		return 0, newErrInvalidLength(8, len(b))
	}
	v := math.Float64frombits(
		uint64(b[7])<<56 |
			uint64(b[6])<<48 |
			uint64(b[5])<<40 |
			uint64(b[4])<<32 |
			uint64(b[3])<<24 |
			uint64(b[2])<<16 |
			uint64(b[1])<<8 |
			uint64(b[0]))
	return v, nil
}

// Float8ToBytes converts the specified float64 to a byte array.
func Float8ToBytes(v float64) []byte {
	b := make([]byte, 8)
	u := math.Float64bits(v)
	b[0] = byte(u & 0xFF)
	b[1] = byte((u >> 8) & 0xFF)
	b[2] = byte((u >> 16) & 0xFF)
	b[3] = byte((u >> 24) & 0xFF)
	b[4] = byte((u >> 32) & 0xFF)
	b[5] = byte((u >> 40) & 0xFF)
	b[6] = byte((u >> 48) & 0xFF)
	b[7] = byte((u >> 56) & 0xFF)
	return b
}

// BytesToFloat4 converts the specified byte array to a float32.
func BytesToFloat4(b []byte) (float32, error) {
	if len(b) != 4 {
		return 0, newErrInvalidLength(4, len(b))
	}
	v := math.Float32frombits(
		uint32(b[3])<<24 |
			uint32(b[2])<<16 |
			uint32(b[1])<<8 |
			uint32(b[0]))
	return v, nil
}

// Float4ToBytes converts the specified float32 to a byte array.
func Float4ToBytes(v float32) []byte {
	b := make([]byte, 4)
	u := math.Float32bits(v)
	b[0] = byte(u & 0xFF)
	b[1] = byte((u >> 8) & 0xFF)
	b[2] = byte((u >> 16) & 0xFF)
	b[3] = byte((u >> 24) & 0xFF)
	return b
}
