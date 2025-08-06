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
	"testing"
	"time"
)

func TestDateEncode(t *testing.T) {
	ts := []time.Time{
		time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.February, 28, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.March, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.April, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.May, 5, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.June, 21, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.July, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.August, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.September, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.October, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.November, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
	}

	for _, tv := range ts {
		b := TimeToDateBytes(tv)

		v, err := BytesToDate(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if !tv.Equal(v) {
			t.Errorf("Failed to convert (%s != %s)", tv, v)
		}
	}
}
