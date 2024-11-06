// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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
	vitessmy "vitess.io/vitess/go/mysql"
)

// AuthManager represents a manager for the user authentication.
type AuthManager struct {
	methods []AuthMethod
}

// NewAuthManager returns a new manager.
func NewAuthManager() *AuthManager {
	return &AuthManager{
		methods: make([]AuthMethod, 0),
	}
}

// AddAuthMethod adds a new authentication method.
func (mgr *AuthManager) AddAuthMethod(method AuthMethod) {
	mgr.methods = append(mgr.methods, method)
}

// AuthMethods returns the list of registered auth methods
// implemented by this auth server.
func (mgr *AuthManager) AuthMethods() []AuthMethod {
	return mgr.methods
}

// DefaultAuthMethodDescription returns MysqlNativePassword as the default
// authentication method for the auth server implementation.
func (mgr *AuthManager) DefaultAuthMethodDescription() AuthMethodDescription {
	return vitessmy.MysqlNativePassword
}
