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

package mysql

import (
	"testing"
)

func TestConnManager(t *testing.T) {
	TestLoopCnt := 10

	cm := NewConnManager()
	if cm.Length() != 0 {
		t.Errorf("%d != %d", cm.Length(), 0)
	}

	for n := 0; n < TestLoopCnt; n++ {
		if cm.Length() != n {
			t.Errorf("%d != %d", cm.Length(), n)
		}
		c := newConn()
		c.ConnectionID = uint32(n)
		cm.AddConn(c)
	}
	if cm.Length() != TestLoopCnt {
		t.Errorf("%d != %d", cm.Length(), TestLoopCnt)
	}

	for n := 0; n < TestLoopCnt; n++ {
		c, ok := cm.GetConnByUID(uint32(n))
		if !ok {
			t.Errorf("%d", n)
			continue
		}
		if c.ConnectionID != uint32(n) {
			t.Errorf("%d != %d", c.ConnectionID, n)
			continue
		}
	}

	for n := 0; n < TestLoopCnt; n++ {
		if cm.Length() != (TestLoopCnt - n) {
			t.Errorf("%d != %d", cm.Length(), (TestLoopCnt - n))
		}
		cm.DeleteConnByUID(uint32(n))
	}
	if cm.Length() != 0 {
		t.Errorf("%d != %d", cm.Length(), 0)
	}
}
