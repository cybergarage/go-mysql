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

// MySQL :: MySQL 8.0 Reference Manual :: 8.2.17 Pluggable Authentication
// https://dev.mysql.com/doc/refman/8.0/en/pluggable-authentication.html
// MySQL :: MySQL 8.0 Reference Manual :: 8.4.1.1 Native Pluggable Authentication
// https://dev.mysql.com/doc/refman/8.0/en/native-pluggable-authentication.html
// Authentication Plugins - MariaDB Knowledge Base
// https://mariadb.com/kb/en/authentication-plugins/

// nolint: gosec
const (
	MySQLNativePassword = "mysql_native_password"
	MySQLOldPassword    = "mysql_old_password"
)
