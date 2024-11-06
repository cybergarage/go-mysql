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

package protocol

import (
	"math"
	"sync"
)

// Counter is a counter.
type Counter struct {
	sync.Mutex
	count int32
}

// NewCounter returns a new counter.
func NewCounter() *Counter {
	return &Counter{
		Mutex: sync.Mutex{},
		count: 0,
	}
}

// NewCounterWith returns a new counter with the specified count.
func NewCounterWith(count int32) *Counter {
	return &Counter{
		Mutex: sync.Mutex{},
		count: count,
	}
}

// Count returns the count.
func (counter *Counter) Count() int32 {
	return counter.count
}

// Inc increments the counter and returns the new value.
func (counter *Counter) Inc() int32 {
	if counter.count == math.MaxInt32 {
		counter.count = 0
	}
	counter.count++
	return counter.count
}
