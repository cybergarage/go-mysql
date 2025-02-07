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

import (
	"testing"

	"github.com/cybergarage/go-mysql/mysql/protocol"
)

func TestNullBitmap(t *testing.T) {
	// Test creating a new NullBitmap with default options
	bmap := protocol.NewNullBitmap()
	if bmap == nil {
		t.Fatal("expected non-nil NullBitmap")
	}
	if len(bmap.Bytes()) != 0 {
		t.Fatalf("expected empty bytes, got %v", bmap.Bytes())
	}

	for numFields := 0; numFields < 32; numFields++ {
		for offset := 0; offset <= 2; offset++ {
			bmap := protocol.NewNullBitmap(
				protocol.WithNullBitmapNumFields(numFields),
				protocol.WithNullBitmapOffset(offset),
			)
			if bmap == nil {
				t.Fatalf("expected non-nil NullBitmap for %d fields", numFields)
			}
			expectedLength := (bmap.NumFields() + 7 + bmap.Offset()) / 8
			if len(bmap.Bytes()) != expectedLength {
				t.Fatalf("expected bytes length %d, got %d", expectedLength, len(bmap.Bytes()))
			}

			// Test setting and getting null values
			for i := 0; i < numFields; i++ {
				if i < (numFields - 1) {
					if bmap.IsNull(i) {
						t.Fatalf("expected field %d to be not null", i)
					}
				}
				bmap.SetNull(i, true)
				if !bmap.IsNull(i) {
					t.Fatalf("expected field %d to be null", i)
				}
				bmap.SetNull(i, false)
				if bmap.IsNull(i) {
					t.Fatalf("expected field %d to be not null", i)
				}
			}
		}
	}
}
