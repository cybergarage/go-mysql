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
	"math"
	"testing"
)

func TestFloat8Endode(t *testing.T) {
	ts := []float64{
		math.MinInt64,
		math.MinInt64 / 2,
		-1,
		0,
		1,
		math.MaxInt64 / 2,
		math.MaxInt64,
	}

	for _, tv := range ts {
		b := Float8ToBytes(tv)

		v, err := BytesToFloat8(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if tv != v {
			t.Errorf("Failed to convert (%f != %f)", tv, v)
		}
	}
}

func TestFloat4Endode(t *testing.T) {
	ts := []float32{
		math.MinInt32,
		math.MinInt32 / 2,
		-1,
		0,
		1,
		math.MaxInt32 / 2,
		math.MaxInt32,
	}

	for _, tv := range ts {
		b := Float4ToBytes(tv)

		v, err := BytesToFloat4(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if tv != v {
			t.Errorf("Failed to convert (%f != %f)", tv, v)
		}
	}
}
