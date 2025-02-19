// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
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
	"encoding/hex"
	"errors"
	"fmt"
)

var ErrInvalid = errors.New("invalid")

func newErrInvalidLength(expected int, actual int) error {
	return fmt.Errorf("%w byte length : %d < %d", ErrInvalid, actual, expected)
}

func newErrInvalidCode(name string, v uint) error {
	return fmt.Errorf("%s is %w code (%X)", name, ErrInvalid, v)
}

func newErrInvalidDatetimeBytes(b []byte) error {
	return fmt.Errorf("%w datetime bytes: %s", ErrInvalid, hex.EncodeToString(b))
}
