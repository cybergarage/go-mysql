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
)

// nolint: gosec
const (
	// MySQL: Authentication Methods
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_protocol_connection_phase_authentication_methods.html#page_protocol_connection_phase_authentication_methods_old_password_authentication
	MySQLOldPassword = "mysql_old_password"

	// MySQL: Native Authentication
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_protocol_connection_phase_authentication_methods_native_password_authentication.html

	MySQLNativePassword = "mysql_native_password"

	// MySQL: Caching_sha2_password information
	// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_caching_sha2_authentication_exchanges.html
	MySQLCachingSHA2Password = "caching_sha2_password"
)

// CredentialAuthenticator is the interface for authenticating a client using credential.
type CredentialAuthenticator = auth.CredentialAuthenticator

// CredentialAuthenticator is the credential authenticator.
type CredentialStore = auth.CredentialStore
