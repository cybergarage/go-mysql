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

var ErrInvalidRequest = errors.New("invalid request")
var ErrNotImplemented = errors.New("not implemented")
var ErrNotFound = errors.New("not found")
var ErrExists = errors.New("already exists")
var ErrInvalidParameter = errors.New("invalid parameter")

// NewSystemDatabaseNotFound returns an error for a system database not found.
func NewSystemDatabaseNotFound(name string) error {
	return fmt.Errorf("system database (%s) is %w", name, ErrNotFound)
}

// NewSystemTableNotFound returns an error for a system table not found.
func NewSystemTableNotFound(name string) error {
	return fmt.Errorf("system table (%s) is %w", name, ErrNotFound)
}

// NewDatabaseNotFound returns an error for a database not found.
func NewDatabaseNotFound(name string) error {
	return fmt.Errorf("database (%s) is %w", name, ErrNotFound)
}

// NewDatabaseExists returns an error for a database already exists.
func NewDatabaseExists(name string) error {
	return fmt.Errorf("database (%s) is %w", name, ErrExists)
}

// NewCollectionNotFound returns an error for a collection not found.
func NewCollectionNotFound(name string) error {
	return fmt.Errorf("collection (%s) is %w", name, ErrNotFound)
}

// NewCollectionExists returns an error for a collection already exists.
func NewCollectionExists(name string) error {
	return fmt.Errorf("collection (%s) is %w", name, ErrExists)
}
