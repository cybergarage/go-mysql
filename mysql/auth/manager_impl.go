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
	"github.com/cybergarage/go-authenticator/auth"
	"github.com/cybergarage/go-mysql/mysql/net"
)

// manager represents a MySQL auth manager.
type manager struct {
	auth.Manager
}

// Manager represents a MySQL auth manager.
func NewManager() Manager {
	return &manager{
		Manager: auth.NewManager(),
	}
}

// Authenticate	authenticates a connection with a query.
func (mgr *manager) Authenticate(conn net.Conn, q Query) bool {
	ok, err := mgr.Manager.VerifyCredential(conn, q)
	if err != nil {
		return false
	}
	return ok
}
