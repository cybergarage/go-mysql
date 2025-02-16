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

package bytes

import (
	"math"
	"testing"

	"github.com/cybergarage/go-mysql/mysql/encoding/bytes"
)

func TestUint64Endode(t *testing.T) {
	ts := []uint64{
		0,
		1,
		math.MaxUint64 / 2,
		math.MaxUint64,
	}

	for _, tv := range ts {
		b := bytes.Uint64ToBytes(tv)
		v, err := bytes.BytesToUint64(b)
		if err != nil {
			t.Error(err)
			continue
		}
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestUint32Endode(t *testing.T) {
	ts := []uint32{
		0,
		1,
		math.MaxUint32 / 2,
		math.MaxUint32,
	}

	for _, tv := range ts {
		b := bytes.Uint32ToBytes(tv)
		v, err := bytes.BytesToUint32(b)
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

func TestUint24Endode(t *testing.T) {
	ts := []uint32{
		0,
		1,
		((1 << 24) - 1) / 2,
		((1 << 24) - 1),
	}

	for _, tv := range ts {
		b := bytes.Uint24ToBytes(tv)
		v, err := bytes.BytesToUint24(b)
		if err != nil {
			t.Error(err)
			continue
		}
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestUint16Endode(t *testing.T) {
	ts := []uint16{
		0,
		1,
		math.MaxUint16 / 2,
		math.MaxUint16,
	}

	for _, tv := range ts {
		b := bytes.Uint16ToBytes(tv)
		v, err := bytes.BytesToUint16(b)
		if err != nil {
			t.Error(err)
			continue
		}
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}

func TestUint8Endode(t *testing.T) {
	ts := []uint8{
		0,
		1,
		math.MaxUint8 / 2,
		math.MaxUint8,
	}

	for _, tv := range ts {
		b := bytes.Uint8ToBytes(tv)
		v, err := bytes.BytesToUint8(b)
		if err != nil {
			t.Error(err)
			continue
		}
		if tv != v {
			t.Errorf("Failed to convert (%d != %d)", tv, v)
		}
	}
}
