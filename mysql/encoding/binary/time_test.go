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

func TestTimeEncode(t *testing.T) {
	td := []time.Duration{
		0,
		1 * time.Second,
		1*time.Hour + 1*time.Minute + 1*time.Second,
		1*time.Hour + 1*time.Minute + 1*time.Second + 1*time.Millisecond,
		1*time.Hour + 1*time.Minute + 1*time.Second + 1*time.Microsecond,
		-1 * time.Second,
		-1*time.Hour - 1*time.Minute - 1*time.Second,
	}

	for _, tv := range td {
		t.Run(tv.String(), func(t *testing.T) {
			b := DurationToTimeBytes(tv)

			v, err := BytesToDuration(b)
			if err != nil {
				t.Error(err)
				return
			}

			if tv != v {
				t.Errorf("Failed to convert (%s != %s)", tv, v)
			}
		})
	}
}
