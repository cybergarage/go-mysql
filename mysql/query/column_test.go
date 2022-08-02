// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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

func TestColumnEquals(t *testing.T) {
	columns := []*Column{
		NewColumn(),
		NewColumnWithName("a"),
		NewColumnWithNameAndValue("a", 1),
	}

	for _, colum := range columns {
		if !colum.Equals(colum) {
			t.Errorf("%s != %s", colum, colum)
		}
	}
}

func TestNewColumnWithComparisonExprs(t *testing.T) {
	cmpExprs := []*vitess.ComparisonExpr{
		{
			Operator: vitess.EqualOp,
			Left:     &vitess.ColName{Name: vitess.NewColIdent("a")},
			Right:    vitess.NewStrLiteral("hello"),
		},
		{
			Operator: vitess.EqualOp,
			Left:     &vitess.ColName{Name: vitess.NewColIdent("b")},
			Right:    vitess.NewIntLiteral("1234"),
		},
		{
			Operator: vitess.EqualOp,
			Left:     &vitess.ColName{Name: vitess.NewColIdent("c")},
			Right:    vitess.NewFloatLiteral("1234"),
		},
	}

	expNames := []interface{}{
		"a",
		"b",
		"c",
	}

	expValues := []interface{}{
		"hello",
		int64(1234),
		float64(1234),
	}

	for n, cmpExpr := range cmpExprs {
		col, err := NewColumnWithComparisonExpr(cmpExpr)
		if err != nil {
			t.Error(err)
		}
		name := col.Name()
		if name != expNames[n] {
			t.Errorf("%s != %s", name, expNames[n])
		}
		val := col.Value()
		if val != expValues[n] {
			t.Errorf("%s != %s", val, expValues[n])
		}
	}
}
