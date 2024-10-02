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

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriter(t *testing.T) {
	w := NewWriter()

	expectedInt1 := uint8(0x61)
	err := w.WriteInt1(expectedInt1)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedInt2 := uint16(0x6261)
	err = w.WriteInt2(expectedInt2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedInt3 := uint32(0x636261)
	err = w.WriteInt3(expectedInt3)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedInt4 := uint32(0x64636261)
	err = w.WriteInt4(expectedInt4)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedInt8 := uint64(0x6867666564636261)
	err = w.WriteInt8(expectedInt8)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedBytes := []byte{0x69, 0x6A, 0x6B, 0x6C}
	_, err = w.WriteBytes(expectedBytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = w.WriteFixedLengthBytes(expectedBytes, 20)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expectedString := "mnop"
	err = w.WriteNullTerminatedString(expectedString)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = w.WriteFixedLengthString(expectedString, 20)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	err = w.WriteLengthEncodedString(expectedString)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	_, err = w.WriteBytes([]byte(expectedString))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	//
	// Test written bytes
	//

	reader := NewReaderWithBytes(w.Bytes())
	actualInt1, err := reader.ReadInt1()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt1 != expectedInt1 {
		t.Errorf("Expected %v, but got %v", expectedInt1, actualInt1)
	}

	actualInt2, err := reader.ReadInt2()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt2 != expectedInt2 {
		t.Errorf("Expected %v, but got %v", expectedInt2, actualInt2)
	}

	actualInt3, err := reader.ReadInt3()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt3 != expectedInt3 {
		t.Errorf("Expected %v, but got %v", expectedInt3, actualInt3)
	}

	actualInt4, err := reader.ReadInt4()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt4 != expectedInt4 {
		t.Errorf("Expected %v, but got %v", expectedInt4, actualInt4)
	}

	actualInt8, err := reader.ReadInt8()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualInt8 != expectedInt8 {
		t.Errorf("Expected %v, but got %v", expectedInt8, actualInt8)
	}

	actualBytes := make([]byte, len(expectedBytes))
	_, err = reader.ReadBytes(actualBytes)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !bytes.Equal(actualBytes, expectedBytes) {
		t.Errorf("Expected %v, but got %v", expectedBytes, actualBytes)
	}

	actualBytes, err = reader.ReadFixedLengthBytes(20)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if bytes.HasPrefix(actualBytes, expectedBytes) == false {
		t.Errorf("Expected %v, but got %v", expectedString, actualBytes)
	}

	actualString, err := reader.ReadNullTerminatedString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}

	actualString, err = reader.ReadFixedLengthString(20)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if strings.HasPrefix(actualString, expectedString) == false {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}

	actualString, err = reader.ReadLengthEncodedString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}

	actualString, err = reader.ReadEOFTerminatedString()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if actualString != expectedString {
		t.Errorf("Expected %v, but got %v", expectedString, actualString)
	}
}
