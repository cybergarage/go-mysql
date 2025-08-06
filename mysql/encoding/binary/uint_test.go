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
	"testing"
)

func TestUint8Endode(t *testing.T) {
	ts := []uint64{
		0,
		1,
		math.MaxUint64 / 2,
		math.MaxUint64,
	}

	for _, tv := range ts {
		b := Uint8ToBytes(tv)

		v, err := BytesToUint8(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestUint4Endode(t *testing.T) {
	ts := []uint32{
		0,
		1,
		math.MaxUint32 / 2,
		math.MaxUint32,
	}

	for _, tv := range ts {
		b := Uint4ToBytes(tv)

		v, err := BytesToUint4(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if err != nil {
			t.Error(err)
			continue
		}

		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestUint3Endode(t *testing.T) {
	ts := []uint32{
		0,
		1,
		((1 << 24) - 1) / 2,
		((1 << 24) - 1),
	}

	for _, tv := range ts {
		b := Uint3ToBytes(tv)

		v, err := BytesToUint3(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestUint2Endode(t *testing.T) {
	ts := []uint16{
		0,
		1,
		math.MaxUint16 / 2,
		math.MaxUint16,
	}

	for _, tv := range ts {
		b := Uint2ToBytes(tv)

		v, err := BytesToUint2(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestUint1Endode(t *testing.T) {
	ts := []uint8{
		0,
		1,
		math.MaxUint8 / 2,
		math.MaxUint8,
	}

	for _, tv := range ts {
		b := Uint1ToBytes(tv)

		v, err := BytesToUint1(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}
