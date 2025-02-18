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

import (
	"errors"
	"fmt"
)

// ErrInvalid is returned when the packet is invalid.
var ErrInvalid = errors.New("invalid")

// ErrNotSupported is returned when the packet is not supported.
var ErrNotSupported = errors.New("not supported")

// ErrNotExist is returned when the specified object is not exist.
var ErrNotExist = errors.New("not exist")

// ErrExist is returned when the specified object is exist.
var ErrExist = errors.New("exist")

// ErrTooManyConnections is returned when the connection is too many.
var ErrTooManyConnections = errors.New("too many connections")

func newErrNotSupported(v any) error {
	return fmt.Errorf("%v is %w", v, ErrNotSupported)
}

func newErrExist(v any) error {
	return fmt.Errorf("%v is %w", v, ErrExist)
}

func newErrNotExist(v any) error {
	return fmt.Errorf("%v is %w", v, ErrNotExist)
}

func newErrInvalidHeader(name string, v byte) error {
	return fmt.Errorf("%s is %w header (%X)", name, ErrInvalid, v)
}

func newErrInvalidCommandType(v CommandType, expected CommandType) error {
	return fmt.Errorf("%02X is %w code (%02X)", v, ErrInvalid, expected)
}

func newErrNotSupportedCommandType(v CommandType) error {
	return fmt.Errorf("command (%02X) is %w", v, ErrNotSupported)
}

func newErrInvalidPacketLength(v uint32) error {
	return fmt.Errorf("packet length is %w (%d)", ErrInvalid, v)
}

func newInvalidStatement(q string) error {
	return fmt.Errorf("%w statement : %s", ErrInvalid, q)
}

func newErrFieldTypeNotSupported(t FieldType) error {
	return fmt.Errorf("%v is %w", t.String(), ErrNotSupported)
}

func newErrInvalidColumnCount(actual int, expected int) error {
	return fmt.Errorf("%w column count (%d, %d)", ErrInvalid, actual, expected)
}
