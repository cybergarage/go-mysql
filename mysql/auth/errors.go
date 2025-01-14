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

package auth

import (
	"errors"
	"fmt"
)

var (
	ErrNotSupported                = errors.New("not supported")
	ErrAccessDenied                = errors.New("access denied")
	ErrInvalidArgument             = errors.New("invalid argument")
	ErrUnknownAuthenticationMethod = errors.New("unknown authentication method")
)

func newErrNotSupported(s string) error {
	return fmt.Errorf("%w %s", ErrNotSupported, s)
}

func newErrUnknownAuthenticationMethod(id string) error {
	return fmt.Errorf("%w: %s", ErrUnknownAuthenticationMethod, id)
}
