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

func TestInt8Endode(t *testing.T) {
	ts := []int64{
		math.MinInt64,
		math.MinInt64 / 2,
		-1,
		0,
		1,
		math.MaxInt64 / 2,
		math.MaxInt64,
	}

	for _, tv := range ts {
		b := Int8ToBytes(tv)
		v, err := BytesToInt8(b)
		if err != nil {
			t.Error(err)
			continue
		}
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestInt4Endode(t *testing.T) {
	ts := []int32{
		math.MinInt32,
		math.MinInt32 / 2,
		-1,
		0,
		1,
		math.MaxInt32 / 2,
		math.MaxInt32,
	}

	for _, tv := range ts {
		b := Int4ToBytes(tv)
		v, err := BytesToInt4(b)
		if err != nil {
			t.Error(err)
			continue
		}
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestInt2Endode(t *testing.T) {
	ts := []int16{
		math.MinInt16,
		math.MinInt16 / 2,
		-1,
		0,
		1,
		math.MaxInt16 / 2,
		math.MaxInt16,
	}

	for _, tv := range ts {
		b := Int2ToBytes(tv)
		v, err := BytesToInt2(b)
		if err != nil {
			t.Error(err)
			continue
		}
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestInt1Endode(t *testing.T) {
	ts := []int8{
		math.MinInt8,
		math.MinInt8 / 2,
		-1,
		0,
		1,
		math.MaxInt8 / 2,
		math.MaxInt8,
	}

	for _, tv := range ts {
		b := Int1ToBytes(tv)
		v, err := BytesToInt1(b)
		if err != nil {
			t.Error(err)
			continue
		}
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}
