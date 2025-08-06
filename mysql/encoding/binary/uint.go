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

// BytesToUint8 converts the specified byte array to an integer.
func BytesToUint8(b []byte) (uint64, error) {
	switch len(b) {
	case 1:
		v, err := BytesToUint1(b)
		return uint64(v), err
	case 2:
		v, err := BytesToUint2(b)
		return uint64(v), err
	case 3:
		v, err := BytesToUint3(b)
		return uint64(v), err
	case 4:
		v, err := BytesToUint4(b)
		return uint64(v), err
	case 8:
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
	return 0, newErrInvalidLength(8, len(b))
}

// Uint8ToBytes converts the specified integer to a byte array.
func Uint8ToBytes(v uint64) []byte {
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

// BytesToUint4 converts the specified byte array to an integer.
func BytesToUint4(b []byte) (uint32, error) {
	switch len(b) {
	case 1:
		v, err := BytesToUint1(b)
		return uint32(v), err
	case 2:
		v, err := BytesToUint2(b)
		return uint32(v), err
	case 3:
		v, err := BytesToUint3(b)
		return uint32(v), err
	case 4:
		v := uint32(b[3])<<24 |
			uint32(b[2])<<16 |
			uint32(b[1])<<8 |
			uint32(b[0])
		return v, nil
	}
	return 0, newErrInvalidLength(4, len(b))
}

// Uint4ToBytes converts the specified integer to a byte array.
func Uint4ToBytes(v uint32) []byte {
	b := make([]byte, 4)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[2] = byte((v >> 16) & 0xFF)
	b[3] = byte((v >> 24) & 0xFF)
	return b
}

// BytesToUint3 converts the specified byte array to an integer.
func BytesToUint3(b []byte) (uint32, error) {
	switch len(b) {
	case 1:
		v, err := BytesToUint1(b)
		return uint32(v), err
	case 2:
		v, err := BytesToUint2(b)
		return uint32(v), err
	case 3:
		v := uint32(b[2])<<16 |
			uint32(b[1])<<8 |
			uint32(b[0])
		return v, nil
	}
	return 0, newErrInvalidLength(3, len(b))
}

// Uint3ToBytes converts the specified integer to a byte array.
func Uint3ToBytes(v uint32) []byte {
	b := make([]byte, 3)
	b[2] = byte((v >> 16) & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	b[0] = byte(v & 0xFF)
	return b
}

// BytesToUint2 converts the specified byte array to an integer.
func BytesToUint2(b []byte) (uint16, error) {
	switch len(b) {
	case 1:
		v, err := BytesToUint1(b)
		return uint16(v), err
	case 2:
		v := uint16(b[1])<<8 |
			uint16(b[0])
		return v, nil
	}
	return 0, newErrInvalidLength(2, len(b))
}

// Uint2ToBytes converts the specified integer to a byte array.
func Uint2ToBytes(v uint16) []byte {
	b := make([]byte, 2)
	b[0] = byte(v & 0xFF)
	b[1] = byte((v >> 8) & 0xFF)
	return b
}

// BytesToUint8 converts the specified byte array to an integer.
func BytesToUint1(b []byte) (uint8, error) {
	if len(b) != 1 {
		return 0, newErrInvalidLength(1, len(b))
	}
	return b[0], nil
}

// Uint1ToBytes converts the specified integer to a byte array.
func Uint1ToBytes(v uint8) []byte {
	b := make([]byte, 1)
	b[0] = v
	return b
}
