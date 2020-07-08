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

package sqltest

import (
	"testing"
)

func TestSQLScenario(t *testing.T) {
	testLines := []string{
		"USE system",
		"{",
		"}",
		"SELECT * FROM system.local",
		"{",
		"}",
	}

	s := NewSQLScenario()
	err := s.ParseLineStrings(testLines)
	if err != nil {
		t.Error(err)
		return
	}

	err = s.IsValid()
	if err != nil {
		t.Error(err)
	}
}
