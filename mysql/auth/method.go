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

// MySQL: Authentication Methods
// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_protocol_connection_phase_authentication_methods.html#page_protocol_connection_phase_authentication_methods_old_password_authentication

// AuthMethod represents the authentication method.
type AuthMethod int

const (
	MySQLAuthenticationNone AuthMethod = iota
	MySQLOldPassword
	MySQLNativePassword
	MySQLCachingSHA2Password
)

const (
	MySQLOldPasswordID         = "mysql_old_password"
	MySQLNativePasswordID      = "mysql_native_password"
	MySQLCachingSHA2PasswordID = "caching_sha2_password"
)

// NewAuthMethodFromID creates a new authentication method from the ID.
func NewAuthMethodFromID(id string) (AuthMethod, error) {
	switch id {
	case MySQLOldPasswordID:
		return MySQLOldPassword, nil
	case MySQLNativePasswordID:
		return MySQLNativePassword, nil
	case MySQLCachingSHA2PasswordID:
		return MySQLCachingSHA2Password, nil
	default:
		return 0, newErrUnknownAuthenticationMethod(id)
	}
}

// EncryptFunc represents the function for encrypting a password.
func (method AuthMethod) EncryptFunc() (EncryptFunc, error) {
	switch method {
	case MySQLNativePassword:
		return NativeEncrypt, nil
	default:
		return nil, newErrNotSupported(method.String())
	}
}

// String returns the string representation of the authentication method.
func (method AuthMethod) String() string {
	switch method {
	case MySQLOldPassword:
		return MySQLOldPasswordID
	case MySQLNativePassword:
		return MySQLNativePasswordID
	case MySQLCachingSHA2Password:
		return MySQLCachingSHA2PasswordID
	default:
		return ""
	}
}
