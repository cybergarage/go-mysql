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

package binary

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
	// Create a buffer with some data
	buf := []byte{0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68}

	// Test PeekInt32 and ReadInt32
	reader := NewReaderWithReader(bytes.NewBuffer(buf))
	expectedInt4 := uint32(0x64636261)

	actualInt4, err := reader.ReadInt4()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualInt4 != expectedInt4 {
		t.Errorf("Expected %v, but got %v", expectedInt4, actualInt4)
	}

	expectedInt4 = uint32(0x68676665)

	actualInt4, err = reader.ReadInt4()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualInt4 != expectedInt4 {
		t.Errorf("Expected %v, but got %v", expectedInt4, actualInt4)
	}

	// Test ReadInt3
	reader = NewReaderWithReader(bytes.NewBuffer(buf))
	expectedInt3 := uint32(0x636261)

	actualInt3, err := reader.ReadInt3()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualInt3 != expectedInt3 {
		t.Errorf("Expected %v, but got %v", expectedInt3, actualInt3)
	}

	// Test ReadInt2
	reader = NewReaderWithReader(bytes.NewBuffer(buf))
	expectedInt2 := uint16(0x6261)

	actualInt2, err := reader.ReadInt2()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualInt2 != expectedInt2 {
		t.Errorf("Expected %v, but got %v", expectedInt2, actualInt2)
	}

	// Test ReadBytesUntil
	reader = NewReaderWithReader(bytes.NewBuffer(buf))
	expectedBytes := []byte{0x61, 0x62, 0x63}

	actualBytes, err := reader.ReadBytesUntil(0x64)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !bytes.Equal(actualBytes, expectedBytes) {
		t.Errorf("Expected %v, but got %v", expectedBytes, actualBytes)
	}

	// Test ReadNullTerminatedString
	reader = NewReaderWithReader(bytes.NewBuffer(append(buf, 0x00)))
	expectedString := "\x61\x62\x63\x64\x65\x66\x67\x68"

	actualString, err := reader.ReadNullTerminatedString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}

	// Test ReadEOFTerminatedString
	reader = NewReaderWithReader(bytes.NewBuffer(buf))
	expectedString = "\x61\x62\x63\x64\x65\x66\x67\x68"

	actualString, err = reader.ReadEOFTerminatedString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}

	// Test ReadFixedLengthString
	reader = NewReaderWithReader(bytes.NewBuffer(buf))
	expectedString = "\x61\x62\x63\x64"

	actualString, err = reader.ReadFixedLengthString(4)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}

	// Test ReadVariableLengthString
	reader = NewReaderWithReader(bytes.NewBuffer(buf))
	expectedString = "\x61\x62\x63\x64"

	actualString, err = reader.ReadVariableLengthString(4)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}
}
