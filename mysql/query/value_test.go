// Copyright (C) 2020 Satoshi Konno. All rights reserved.
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

package query

import (
	"testing"

	vitess "vitess.io/vitess/go/vt/sqlparser"
)

func TestNewValue(t *testing.T) {
	literals := []*vitess.Literal{
		vitess.NewStrLiteral("hello"),
		vitess.NewIntLiteral("1234"),
		vitess.NewFloatLiteral("1234"),
	}
	expValues := []interface{}{
		"hello",
		int64(1234),
		float64(1234),
	}

	for n, l := range literals {
		val, err := NewValueWithLiteral(l)
		if err != nil {
			t.Error(err)
			continue
		}
		if val != expValues[n] {
			t.Errorf("%s != %s", val, expValues[n])
		}
	}
}
