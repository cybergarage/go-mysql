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

// AuthenticationMethod represents the authentication method.
type AuthenticationMethod int

// Authentication methods.
const (
	MySQLOldPassword AuthenticationMethod = iota
	// MySQL: Native Authentication
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_protocol_connection_phase_authentication_methods_native_password_authentication.html
	MySQLNativePassword
	// MySQL: Caching_sha2_password information
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_caching_sha2_authentication_exchanges.html
	MySQLCachingSHA2Password
)

const (
	// MySQL: Authentication Methods
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_protocol_connection_phase_authentication_methods.html#page_protocol_connection_phase_authentication_methods_old_password_authentication
	MySQLOldPasswordID         = "mysql_old_password"
	MySQLNativePasswordID      = "mysql_native_password"
	MySQLCachingSHA2PasswordID = "caching_sha2_password"
)

// NewAuthenticationMethodFromID creates a new authentication method from the ID.
func NewAuthenticationMethodFromID(id string) (AuthenticationMethod, error) {
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
