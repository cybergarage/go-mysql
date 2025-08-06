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

package plugins

// MySQL: Clear text client plugin
// https://dev.mysql.com/doc/dev/mysql-server/latest/page_protocol_connection_phase_authentication_methods_clear_text_password.html
func ClearEncrypt(passwd any, args ...any) (any, error) {
	var strPasswd string
	switch v := passwd.(type) {
	case string:
		strPasswd = v
	case []byte:
		strPasswd = string(v)
	default:
		return nil, ErrInvalidArgument
	}

	return strPasswd, nil
}
