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

package errors

import (
	"errors"
	"fmt"
)

// ErrNotImplemented is returned when the operation is not implemented.
var ErrNotImplemented = errors.New("not implemented")

// ErrNotSupported is returned when the operation is not supported.
var ErrNotSupported = errors.New("not supported")

// ErrNotExist is returned when the specified object is not exist.
var ErrNotExist = errors.New("not exist")

// ErrExist is returned when the specified object is exist.
var ErrExist = errors.New("exist")

// ErrNotEqual is returned when the specified object is not equal.
var ErrNotEqual = errors.New("not equal")

// ErrInvalid is returned when the specified object is invalid.
var ErrInvalid = errors.New("invalid")

// ErrNotFound is returned when the specified object is not found.
var ErrNotFound = errors.New("not found")

// NewErrNotImplemented returns an error for a not implemented.
func NewErrNotImplemented(s string) error {
	return fmt.Errorf("%s is %w", s, ErrNotImplemented)
}

// NewErrNotExis returns an error for a not exist.
func NewErrNotExis(name string) error {
	return fmt.Errorf("%s is not exists", name)
}
