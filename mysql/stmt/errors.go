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

package stmt

import (
	"errors"
	"fmt"
)

// ErrInvalid is the error for invalid statement ID.
var ErrInvalid = errors.New("invalid")

// ErrOverflow is returned when the value is overflow.
var ErrOverflow = errors.New("overflow")

// ErrNotSupported is returned when the feature is not supported.
var ErrNotSupported = errors.New("not supported")

func newErrInvalidStatementID(stmdID StatementID) error {
	return fmt.Errorf("%w statement ID: %d", ErrInvalid, stmdID)
}

func newErrInvalidQuery(query string) error {
	return fmt.Errorf("%w query: %s", ErrInvalid, query)
}

func newErrInvalidParameters() error {
	return fmt.Errorf("%w parameters", ErrInvalid)
}

func newErrNotSupportedFieldType(typ FieldType) error {
	return fmt.Errorf("%w field type: %s", ErrNotSupported, typ.String())
}
