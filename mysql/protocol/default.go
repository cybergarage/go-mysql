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

package protocol

import (
	"github.com/cybergarage/go-mysql/mysql/auth"
)

const (
	DefaultPort                  = 3306
	DefaultMaxPacketSize         = 0
	DefaultCharset               = CharSetUTF8
	DefaultAuthPluginDataPartLen = 20

	SupportVersion = "5.7.9"

	DefaultServerCapability = ClientLongPassword |
		ClientFoundRows |
		ClientLongColumnFlag |
		ClientProtocol41 |
		ClientSecureConnection |
		ClientPluginAuth

	DefaultHandshakeServerCapabilities = DefaultServerCapability |
		ClientConnectWithDB

	DefaultSSLRequestCapabilities = DefaultServerCapability |
		ClientSSL

	DefaultServerStatus = ServerStatusAutocommit

	DefaultAuthPluginName = auth.MySQLNativePasswordID
)
