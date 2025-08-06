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

func TestDatetimeEncode(t *testing.T) {
	ts := []time.Time{
		time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.February, 28, 23, 59, 59, 999999000, time.UTC),
		time.Date(2023, time.March, 15, 12, 30, 45, 123456000, time.UTC),
		time.Date(2023, time.April, 10, 6, 15, 30, 654321000, time.UTC),
		time.Date(2023, time.May, 5, 18, 45, 15, 987654000, time.UTC),
		time.Date(2023, time.June, 21, 9, 0, 0, 0, time.UTC),
		time.Date(2023, time.July, 4, 0, 0, 0, 0, time.UTC),
		time.Date(2023, time.August, 31, 23, 59, 59, 999999000, time.UTC),
		time.Date(2023, time.September, 30, 12, 0, 0, 0, time.UTC),
		time.Date(2023, time.October, 15, 6, 30, 30, 123456000, time.UTC),
		time.Date(2023, time.November, 11, 11, 11, 11, 111111000, time.UTC),
		time.Date(2023, time.December, 25, 0, 0, 0, 0, time.UTC),
	}

	for _, tv := range ts {
		b := TimeToDatetimeBytes(tv)

		v, err := BytesToDatetime(b)
		if err != nil {
			t.Error(err)
			continue
		}

		if !tv.Equal(v) {
			t.Errorf("Failed to convert (%s != %s)", tv, v)
		}
	}
}
