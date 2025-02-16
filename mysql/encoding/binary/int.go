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

// BytesToInt64 converts the specified byte array to an integer.
func BytesToInt64(b []byte) (int64, error) {
	if len(b) != 8 {
		return 0, newErrInvalidLength()
	}
	v := int64(b[7])<<56 |
		int64(b[6])<<48 |
		int64(b[5])<<40 |
		int64(b[4])<<32 |
		int64(b[3])<<24 |
		int64(b[2])<<16 |
		int64(b[1])<<8 |
		int64(b[0])
	return v, nil
}

// Int64ToBytes converts the specified integer to a byte array.
func Int64ToBytes(v int64) []byte {
	b := make([]byte, 8)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[2] = byte((v >> 16) & 0xFF)
	b[3] = byte((v >> 24) & 0xFF)
	b[4] = byte((v >> 32) & 0xFF)
	b[5] = byte((v >> 40) & 0xFF)
	b[6] = byte((v >> 48) & 0xFF)
	b[7] = byte((v >> 56) & 0xFF)
	return b
}

// BytesToInt32 converts the specified byte array to an integer.
func BytesToInt32(b []byte) (int32, error) {
	if len(b) != 4 {
		return 0, newErrInvalidLength()
	}
	v := int32(b[3])<<24 |
		int32(b[2])<<16 |
		int32(b[1])<<8 |
		int32(b[0])
	return v, nil
}

// Int32ToBytes converts the specified integer to a byte array.
func Int32ToBytes(v int32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[2] = byte((v >> 16) & 0xFF)
	b[3] = byte((v >> 24) & 0xFF)
	return b
}

// BytesToInt24 converts the specified byte array to an integer.
func BytesToInt24(b []byte) (int32, error) {
	if len(b) != 3 {
		return 0, newErrInvalidLength()
	}
	v := int32(b[2])<<16 |
		int32(b[1])<<8 |
		int32(b[0])
	return v, nil
}

// Int24ToBytes converts the specified integer to a byte array.
func Int24ToBytes(v int32) []byte {
	b := make([]byte, 3)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[2] = byte((v >> 16) & 0xFF)
	return b
}

// BytesToInt16 converts the specified byte array to an integer.
func BytesToInt16(b []byte) (int16, error) {
	if len(b) != 2 {
		return 0, newErrInvalidLength()
	}
	v := int16(b[1])<<8 |
		int16(b[0])
	return v, nil
}

// Int16ToBytes converts the specified integer to a byte array.
func Int16ToBytes(v int16) []byte {
	b := make([]byte, 2)
	b[1] = byte((v >> 8) & 0xFF)
	b[0] = byte(v & 0xFF)
	return b
}

// BytesToInt8 converts the specified byte array to an integer.
func BytesToInt8(b []byte) (int8, error) {
	if len(b) != 1 {
		return 0, newErrInvalidLength()
	}
	return int8(b[0]), nil
}

// Int8ToBytes converts the specified integer to a byte array.
func Int8ToBytes(v int8) []byte {
	b := make([]byte, 1)
	b[0] = byte(v)
	return b
}
