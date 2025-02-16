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

// BytesToUint64 converts the specified byte array to an integer.
func BytesToUint64(b []byte) (uint64, error) {
	if len(b) != 8 {
		return 0, newErrInvalidLength()
	}
	v := uint64(b[7])<<56 |
		uint64(b[6])<<48 |
		uint64(b[5])<<40 |
		uint64(b[4])<<32 |
		uint64(b[3])<<24 |
		uint64(b[2])<<16 |
		uint64(b[1])<<8 |
		uint64(b[0])
	return v, nil
}

// Uint64ToBytes converts the specified integer to a byte array.
func Uint64ToBytes(v uint64) []byte {
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

// BytesToUint32 converts the specified byte array to an integer.
func BytesToUint32(b []byte) (uint32, error) {
	if len(b) != 4 {
		return 0, newErrInvalidLength()
	}
	v := uint32(b[3])<<24 |
		uint32(b[2])<<16 |
		uint32(b[1])<<8 |
		uint32(b[0])
	return v, nil
}

// Uint32ToBytes converts the specified integer to a byte array.
func Uint32ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[2] = byte((v >> 16) & 0xFF)
	b[3] = byte((v >> 24) & 0xFF)
	return b
}

// BytesToUint24 converts the specified byte array to an integer.
func BytesToUint24(b []byte) (uint32, error) {
	if len(b) != 3 {
		return 0, newErrInvalidLength()
	}
	v := uint32(b[2])<<16 |
		uint32(b[1])<<8 |
		uint32(b[0])
	return v, nil
}

// Uint24ToBytes converts the specified integer to a byte array.
func Uint24ToBytes(v uint32) []byte {
	b := make([]byte, 3)
	b[2] = byte((v >> 16) & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[0] = byte(v & 0xFF)
	return b
}

// BytesToUint16 converts the specified byte array to an integer.
func BytesToUint16(b []byte) (uint16, error) {
	if len(b) != 2 {
		return 0, newErrInvalidLength()
	}
	v := uint16(b[1])<<8 |
		uint16(b[0])
	return v, nil
}

// Uint16ToBytes converts the specified integer to a byte array.
func Uint16ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	return b
}

// BytesToUint8 converts the specified byte array to an integer.
func BytesToUint8(b []byte) (uint8, error) {
	if len(b) != 1 {
		return 0, newErrInvalidLength()
	}
	return b[0], nil
}

// Uint8ToBytes converts the specified integer to a byte array.
func Uint8ToBytes(v uint8) []byte {
	b := make([]byte, 1)
	b[0] = v
	return b
}
