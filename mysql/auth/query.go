// Copyright (C) 2019 The go-mysql Authors. All rights reserved.
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
	"github.com/cybergarage/go-sasl/sasl/auth"
)

// Query represents a credential query.
type Query = auth.Query

// QueryOptionFn is a query option function.
type QueryOptionFn = auth.QueryOptionFn

// NewQuery returns a new query with options.
func NewQuery(opts ...QueryOptionFn) (Query, error) {
	return auth.NewQuery(opts...)
}

// WithQueryUsername returns an option to set the username.
func WithQueryUsername(username string) QueryOptionFn {
	return auth.WithQueryUsername(username)
}

// WithQueryAuthResponse returns an option to set the password.
func WithQueryAuthResponse(password string) QueryOptionFn {
	return auth.WithQueryPassword(password)
}

// WithQueryClientPluginName returns an option to set the client plugin name.
func WithQueryClientPluginName(clientPluginName string) QueryOptionFn {
	return func(q Query) error {
		method, err := NewAuthMethodFromID(clientPluginName)
		if err != nil {
			return err
		}
		if method == MySQLAuthenticationNone {
			return nil
		}
		encryptFunc, err := method.EncryptFunc()
		if err != nil {
			return err
		}
		q.SetEncryptFunc(encryptFunc)
		return nil
	}
}
