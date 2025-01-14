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
	"crypto/sha1"
	"encoding/hex"

	"github.com/cybergarage/go-authenticator/auth"
)

// EncryptFunc is the interface for encrypting a password.
type EncryptFunc = auth.EncryptFunc

// MySQL: Native Authentication
// https://dev.mysql.com/doc/dev/mysql-server/8.4.2/page_protocol_connection_phase_authentication_methods_native_password_authentication.html
// NativeEncrypt encrypts the password using the native MySQL password encryption algorithm.
func NativeEncrypt(passwd string, args ...any) (string, error) {
	nativeEncrypt := func(passwd string, rndData []byte) string {
		// SHA1( password ) XOR SHA1( "20-bytes random data from server" <concat> SHA1( SHA1( password ) ) )
		xor := func(a, b []byte) []byte {
			minLength := min(len(a), len(b))
			result := make([]byte, minLength)
			for n := 0; n < minLength; n++ {
				result[n] = a[n] ^ b[n]
			}
			return result
		}

		h := sha1.New()

		h.Write([]byte(passwd))
		passwdHash := h.Sum(nil)

		h.Reset()
		h.Write(passwdHash)
		passwdHashHash := h.Sum(nil)

		h.Reset()
		h.Write(rndData)
		h.Write(passwdHashHash)
		rndDataHash := h.Sum(nil)

		return hex.EncodeToString(xor(passwdHash, rndDataHash))
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

	return nativeEncrypt(passwd, rndData), nil
}
