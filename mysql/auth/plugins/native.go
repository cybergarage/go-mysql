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

import (
	"crypto/sha1"
)

// MySQL: Native Authentication
// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_protocol_connection_phase_authentication_methods_native_password_authentication.html
// Authentication Plugin - mysql_native_password - MariaDB Knowledge Base
// https://mariadb.com/kb/en/authentication-plugin-mysql_native_password/
// NativeEncrypt encrypts the password using the native MySQL password encryption algorithm.
func NativeEncrypt(passwd any, args ...any) (any, error) {
	nativeEncrypt := func(passwd any, rndData []byte) ([]byte, error) {
		// SHA1( password ) XOR SHA1( "20-bytes random data from server" <concat> SHA1( SHA1( password ) ) )
		xor := func(a, b []byte) []byte {
			minLength := min(len(a), len(b))
			result := make([]byte, minLength)
			for n := 0; n < minLength; n++ {
				result[n] = a[n] ^ b[n]
			}
			return result
		}

		var bytesPasswd []byte
		switch v := passwd.(type) {
		case string:
			bytesPasswd = []byte(v)
		case []byte:
			bytesPasswd = v
		default:
			return nil, ErrInvalidArgument
		}

		h := sha1.New()

		h.Write(bytesPasswd)
		passwdHash := h.Sum(nil)

		h.Reset()
		h.Write(passwdHash)
		passwdHashHash := h.Sum(nil)

		h.Reset()
		h.Write(rndData)
		h.Write(passwdHashHash)
		rndDataHash := h.Sum(nil)

		return xor(passwdHash, rndDataHash), nil
	}

	if len(args) == 0 {
		return "", ErrInvalidArgument
	}
	rndData, ok := args[0].([]byte)
	if !ok {
		return "", ErrInvalidArgument
	}
	if len(rndData) != 20 {
		return "", ErrInvalidArgument
	}

	return nativeEncrypt(passwd, rndData)
}
