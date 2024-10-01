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

// MySQL: Capabilities Flags
// https://dev.mysql.com/doc/dev/mysql-server/latest/group__group__cs__capabilities__flags.html

// CapabilityFlag represents a MySQL Capability Flag.
type CapabilityFlag uint32

const (
	// CapabilityFlagClientLongPassword represents the CLIENT_LONG_PASSWORD capability flag.
	CapabilityFlagClientLongPassword CapabilityFlag = 1
	// CapabilityFlagClientFoundRows represents the CLIENT_FOUND_ROWS capability flag.
	CapabilityFlagClientFoundRows CapabilityFlag = 2
	// CapabilityFlagClientLongColumnFlag represents the CLIENT_LONG_FLAG capability flag.
	CapabilityFlagClientLongColumnFlag CapabilityFlag = 4
	// CapabilityFlagClientConnectWithDB represents the CLIENT_CONNECT_WITH_DB capability flag.
	CapabilityFlagClientConnectWithDB CapabilityFlag = 8
	// CapabilityFlagClientNoSchema represents the CLIENT_NO_SCHEMA capability flag.
	CapabilityFlagClientNoSchema CapabilityFlag = 16
	// CapabilityFlagClientCompress represents the CLIENT_COMPRESS capability flag.
	CapabilityFlagClientCompress CapabilityFlag = 32
	// CapabilityFlagClientODBC represents the CLIENT_ODBC capability flag.
	CapabilityFlagClientODBC CapabilityFlag = 64
	// CapabilityFlagClientLocalFiles represents the CLIENT_LOCAL_FILES capability flag.
	CapabilityFlagClientLocalFiles CapabilityFlag = 128
	// CapabilityFlagClientIgnoreSpace represents the CLIENT_IGNORE_SPACE capability flag.
	CapabilityFlagClientIgnoreSpace CapabilityFlag = 256
	// CapabilityFlagClientProtocol41 represents the CLIENT_PROTOCOL_41 capability flag.
	CapabilityFlagClientProtocol41 CapabilityFlag = 512
	// CapabilityFlagClientInteractive represents the CLIENT_INTERACTIVE capability flag.
	CapabilityFlagClientInteractive CapabilityFlag = 1024
	// CapabilityFlagClientSSL represents the CLIENT_SSL capability flag.
	CapabilityFlagClientSSL CapabilityFlag = 2048
	// CapabilityFlagClientIgnoreSIGPIPE represents the CLIENT_IGNORE_SIGPIPE capability flag.
	CapabilityFlagClientIgnoreSIGPIPE CapabilityFlag = 4096
	// CapabilityFlagClientTransactions represents the CLIENT_TRANSACTIONS capability flag.
	CapabilityFlagClientTransactions CapabilityFlag = 8192
	// CapabilityFlagClientReserved represents the CLIENT_RESERVED capability flag.
	CapabilityFlagClientReserved CapabilityFlag = 16384
	// CapabilityFlagClientSecureConnection represents the CLIENT_SECURE_CONNECTION capability flag.
	CapabilityFlagClientSecureConnection CapabilityFlag = 32768
	// CapabilityFlagClientMultiStatements represents the CLIENT_MULTI_STATEMENTS capability flag.
	CapabilityFlagClientMultiStatements CapabilityFlag = 65536
	// CapabilityFlagClientMultiResults represents the CLIENT_MULTI_RESULTS capability flag.
	CapabilityFlagClientMultiResults CapabilityFlag = 131072
	// CapabilityFlagClientPSMultiResults represents the CLIENT_PS_MULTI_RESULTS capability flag.
	CapabilityFlagClientPSMultiResults CapabilityFlag = 262144
	// CapabilityFlagClientPluginAuth represents the CLIENT_PLUGIN_AUTH capability flag.
	CapabilityFlagClientPluginAuth CapabilityFlag = 524288
	// CapabilityFlagClientConnectAttrs represents the CLIENT_CONNECT_ATTRS capability flag.
	CapabilityFlagClientConnectAttrs CapabilityFlag = 1048576
	// CapabilityFlagClientPluginAuthLenencClientData represents the CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA capability flag.
	CapabilityFlagClientPluginAuthLenencClientData CapabilityFlag = 2097152
	// CapabilityFlagClientCanHandleExpiredPasswords represents the CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS capability flag.
	CapabilityFlagClientCanHandleExpiredPasswords CapabilityFlag = 4194304
	// CapabilityFlagClientSessionTrack represents the CLIENT_SESSION_TRACK capability flag.
	CapabilityFlagClientSessionTrack CapabilityFlag = 8388608
	// CapabilityFlagClientDeprecateEOF represents the CLIENT_DEPRECATE_EOF capability flag.
	CapabilityFlagClientDeprecateEOF CapabilityFlag = 16777216
	// CapabilityFlagClientCapabilitiyClientOptionalResultsetMetadata represents the CLIENT_OPTIONAL_RESULTSET_METADATA capability flag.
	CapabilitiyClientOptionalResultsetMetadata CapabilityFlag = 33554432
	// CapabilityFlagClientZstdCompressionAlgorithms represents the CLIENT_ZSTD_COMPRESSION_ALGORITHMS capability flag.
	CapabilityClientQueryAttributes CapabilityFlag = 67108864
	// CapabilityMultiFactoryAuth represents the CLIENT_MULTI_FACTORY_AUTH capability flag.
	CapabilityMultiFactoryAuth CapabilityFlag = 134217728
	// CapabilityClientCapabilityExtension represents the CLIENT_CAPABILITY_EXTENSION capability flag.
	CapabilityClientCapabilityExtension CapabilityFlag = 268435456
	// CapabilityClientSSLVerifyServerCert represents the CLIENT_SSL_VERIFY_SERVER_CERT capability flag.
	CapabilityClientSSLVerifyServerCert CapabilityFlag = 536870912
	// CapabilityClientRemenberOptions represents the CLIENT_REMENBER_OPTIONS capability flag.
	CapabilityRemenberOptions CapabilityFlag = 1073741824
)

// IsEnabled returns true if the specified flag is set.
func (cap CapabilityFlag) IsEnabled(flag CapabilityFlag) bool {
	return (cap & flag) != 0
}
