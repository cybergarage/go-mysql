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
	"github.com/cybergarage/go-mysql/mysql/net"
	"github.com/cybergarage/go-sasl/sasl"
)

// Manager represents a MySQL auth manager.
type Manager struct {
	*sasl.Server
}

// NewManager returns a new SASL server.
func NewManager() *Manager {
	return &Manager{
		Server: sasl.NewServer(),
	}
}

// Authenticators returns the authenticators.
func (mgr *Manager) Authenticate(conn net.Conn, q *Query) bool {
	auths := mgr.Authenticators()
	if len(auths) == 0 {
		return true
	}
	for _, auth := range auths {
		if _, ok := auth.HasCredential(q); ok {
			return true
		}
	}
	return false
}
