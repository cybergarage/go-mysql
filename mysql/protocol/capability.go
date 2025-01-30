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

// Capability represents a MySQL Capability Flag.
type Capability uint32

const (
	// ClientLongPassword represents the CLIENT_LONG_PASSWORD capability flag.
	ClientLongPassword Capability = 1
	// ClientFoundRows represents the CLIENT_FOUND_ROWS capability flag.
	ClientFoundRows Capability = 2
	// ClientLongColumnFlag represents the CLIENT_LONG_FLAG capability flag.
	ClientLongColumnFlag Capability = 4
	// ClientConnectWithDB represents the CLIENT_CONNECT_WITH_DB capability flag.
	ClientConnectWithDB Capability = 8
	// ClientNoSchema represents the CLIENT_NO_SCHEMA capability flag.
	ClientNoSchema Capability = 16
	// ClientCompress represents the CLIENT_COMPRESS capability flag.
	ClientCompress Capability = 32
	// ClientODBC represents the CLIENT_ODBC capability flag.
	ClientODBC Capability = 64
	// ClientLocalFiles represents the CLIENT_LOCAL_FILES capability flag.
	ClientLocalFiles Capability = 128
	// ClientIgnoreSpace represents the CLIENT_IGNORE_SPACE capability flag.
	ClientIgnoreSpace Capability = 256
	// ClientProtocol41 represents the CLIENT_PROTOCOL_41 capability flag.
	ClientProtocol41 Capability = 512
	// ClientInteractive represents the CLIENT_INTERACTIVE capability flag.
	ClientInteractive Capability = 1024
	// ClientSSL represents the CLIENT_SSL capability flag.
	ClientSSL Capability = 2048
	// ClientIgnoreSIGPIPE represents the CLIENT_IGNORE_SIGPIPE capability flag.
	ClientIgnoreSIGPIPE Capability = 4096
	// ClientTransactions represents the CLIENT_TRANSACTIONS capability flag.
	ClientTransactions Capability = 8192
	// ClientReserved represents the CLIENT_RESERVED capability flag.
	ClientReserved Capability = 16384
	// ClientSecureConnection represents the CLIENT_SECURE_CONNECTION capability flag.
	ClientSecureConnection Capability = 32768
	// ClientMultiStatements represents the CLIENT_MULTI_STATEMENTS capability flag.
	ClientMultiStatements Capability = 65536
	// ClientMultiResults represents the CLIENT_MULTI_RESULTS capability flag.
	ClientMultiResults Capability = 131072
	// ClientPSMultiResults represents the CLIENT_PS_MULTI_RESULTS capability flag.
	ClientPSMultiResults Capability = 262144
	// ClientPluginAuth represents the CLIENT_PLUGIN_AUTH capability flag.
	ClientPluginAuth Capability = 524288
	// ClientConnectAttrs represents the CLIENT_CONNECT_ATTRS capability flag.
	ClientConnectAttrs Capability = 1048576
	// ClientPluginAuthLenencClientData represents the CLIENT_PLUGIN_AUTH_LENENC_CLIENT_DATA capability flag.
	ClientPluginAuthLenencClientData Capability = 2097152
	// ClientCanHandleExpiredPasswords represents the CLIENT_CAN_HANDLE_EXPIRED_PASSWORDS capability flag.
	ClientCanHandleExpiredPasswords Capability = 4194304
	// ClientSessionTrack represents the CLIENT_SESSION_TRACK capability flag.
	ClientSessionTrack Capability = 8388608
	// ClientDeprecateEOF represents the CLIENT_DEPRECATE_EOF capability flag.
	ClientDeprecateEOF Capability = 16777216
	// ClientCapabilitiyClientOptionalResultsetMetadata represents the CLIENT_OPTIONAL_RESULTSET_METADATA capability flag.
	ClientOptionalResultsetMetadata Capability = 33554432
	// ClientZstdCompressionAlgorithm represents the CLIENT_ZSTD_COMPRESSION_ALGORITHMS capability flag.
	ClientZstdCompressionAlgorithm Capability = 67108864
	// ClientQueryAttributes represents the CLIENT_QUERY_ATTRIBUTES capability flag.
	ClientQueryAttributes Capability = 134217728
	// CapabilityMultiFactoryAuth represents the CLIENT_MULTI_FACTORY_AUTH capability flag.
	CapabilityMultiFactoryAuth Capability = 67108864
	// ClientCapabilityExtension represents the CLIENT_CAPABILITY_EXTENSION capability flag.
	ClientCapabilityExtension Capability = 268435456
	// ClientSSLVerifyServerCert represents the CLIENT_SSL_VERIFY_SERVER_CERT capability flag.
	ClientSSLVerifyServerCert Capability = 536870912
	// ClientRemenberOptions represents the CLIENT_REMENBER_OPTIONS capability flag.
	ClientRemenberOptions Capability = 1073741824
)

// IsEnabled returns true if the specified flag is set.
func (c Capability) IsEnabled(flag Capability) bool {
	return (c & flag) != 0
}

// IsDisabled returns true if the specified flag is not set.
func (c Capability) IsDisabled(flag Capability) bool {
	return !c.IsEnabled(flag)
}

// ToBytes returns the capability flag as bytes.
func NewCapabilityFromBytes(data []byte) Capability {
	return Capability(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24)
}
