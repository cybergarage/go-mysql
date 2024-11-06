// Copyright (C) 2024 The go-mysql Authors. All rights reserved.
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

type str struct {
	v      string
	isNull bool
}

// NewString returns a new string.
func NewString(v string) String {
	return &str{v: v, isNull: false}
}

// NewNull returns a new null string.
func NewNull(v string) String {
	return &str{v: "", isNull: true}
}

// Value returns the string value.
func (s *str) Value() string {
	return s.v
}

// IsNull returns true if the string is null.
func (s *str) IsNull() bool {
	return s.isNull
}
