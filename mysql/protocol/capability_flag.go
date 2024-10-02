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
	// ClientLongPassword represents the CLIENT_LONG_PASSWORD capability flag.
	ClientLongPassword CapabilityFlag = 1
	// ClientFoundRows represents the CLIENT_FOUND_ROWS capability flag.
	ClientFoundRows CapabilityFlag = 2
	// ClientLongColumnFlag represents the CLIENT_LONG_FLAG capability flag.
	ClientLongColumnFlag CapabilityFlag = 4
	// ClientConnectWithDB represents the CLIENT_CONNECT_WITH_DB capability flag.
	ClientConnectWithDB CapabilityFlag = 8
	// ClientNoSchema represents the CLIENT_NO_SCHEMA capability flag.
	ClientNoSchema CapabilityFlag = 16
	// ClientCompress represents the CLIENT_COMPRESS capability flag.
	ClientCompress CapabilityFlag = 32
	// ClientODBC represents the CLIENT_ODBC capability flag.
	ClientODBC CapabilityFlag = 64
	// ClientLocalFiles represents the CLIENT_LOCAL_FILES capability flag.
	ClientLocalFiles CapabilityFlag = 128
	// ClientIgnoreSpace represents the CLIENT_IGNORE_SPACE capability flag.
	ClientIgnoreSpace CapabilityFlag = 256
	// ClientProtocol41 represents the CLIENT_PROTOCOL_41 capability flag.
	ClientProtocol41 CapabilityFlag = 512
	// ClientInteractive represents the CLIENT_INTERACTIVE capability flag.
	ClientInteractive CapabilityFlag = 1024
	// ClientSSL represents the CLIENT_SSL capability flag.
	ClientSSL CapabilityFlag = 2048
	// ClientIgnoreSIGPIPE represents the CLIENT_IGNORE_SIGPIPE capability flag.
	ClientIgnoreSIGPIPE CapabilityFlag = 4096
	// ClientTransactions represents the CLIENT_TRANSACTIONS capability flag.
	ClientTransactions CapabilityFlag = 8192
	// ClientReserved represents the CLIENT_RESERVED capability flag.
	ClientReserved CapabilityFlag = 16384
	// ClientSecureConnection represents the CLIENT_SECURE_CONNECTION capability flag.
	ClientSecureConnection CapabilityFlag = 32768
	// ClientMultiStatements represents the CLIENT_MULTI_STATEMENTS capability flag.
	ClientMultiStatements CapabilityFlag = 65536
	// ClientMultiResults represents the CLIENT_MULTI_RESULTS capability flag.
	ClientMultiResults CapabilityFlag = 131072
	// ClientPSMultiResults represents the CLIENT_PS_MULTI_RESULTS capability flag.
	ClientPSMultiResults CapabilityFlag = 262144
	// ClientPluginAuth represents the CLIENT_PLUGIN_AUTH capability flag.
	ClientPluginAuth CapabilityFlag = 524288
	// ClientConnectAttrs represents the CLIENT_CONNECT_ATTRS capability flag.
	ClientConnectAttrs CapabilityFlag = 1048576
	// ClientPluginAuthLenencClientData represents the CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA capability flag.
	ClientPluginAuthLenencClientData CapabilityFlag = 2097152
	// ClientCanHandleExpiredPasswords represents the CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS capability flag.
	ClientCanHandleExpiredPasswords CapabilityFlag = 4194304
	// ClientSessionTrack represents the CLIENT_SESSION_TRACK capability flag.
	ClientSessionTrack CapabilityFlag = 8388608
	// ClientDeprecateEOF represents the CLIENT_DEPRECATE_EOF capability flag.
	ClientDeprecateEOF CapabilityFlag = 16777216
	// ClientCapabilitiyClientOptionalResultsetMetadata represents the CLIENT_OPTIONAL_RESULTSET_METADATA capability flag.
	ClientOptionalResultsetMetadata CapabilityFlag = 33554432
	// ClientZstdCompressionAlgorithm represents the CLIENT_ZSTD_COMPRESSION_ALGORITHMS capability flag.
	ClientZstdCompressionAlgorithm CapabilityFlag = 67108864
	// CapabilityMultiFactoryAuth represents the CLIENT_MULTI_FACTORY_AUTH capability flag.
	CapabilityMultiFactoryAuth CapabilityFlag = 134217728
	// ClientCapabilityExtension represents the CLIENT_CAPABILITY_EXTENSION capability flag.
	ClientCapabilityExtension CapabilityFlag = 268435456
	// ClientSSLVerifyServerCert represents the CLIENT_SSL_VERIFY_SERVER_CERT capability flag.
	ClientSSLVerifyServerCert CapabilityFlag = 536870912
	// ClientRemenberOptions represents the CLIENT_REMENBER_OPTIONS capability flag.
	ClientRemenberOptions CapabilityFlag = 1073741824
)

// IsEnabled returns true if the specified flag is set.
func (capFlg CapabilityFlag) IsEnabled(flag CapabilityFlag) bool {
	return (capFlg & flag) != 0
}
