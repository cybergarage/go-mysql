// Copyright (C) 2020 The go-mysql Authors. All rights reserved.
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

package mysql

import (
	"crypto/tls"
)

// CertConfig represents a TLS configuration interface.
type CertConfig interface {
	// SetClientAuthType sets a client authentication type.
	SetClientAuthType(authType tls.ClientAuthType)
	// SetServerKeyFile sets a SSL server key file.
	SetServerKeyFile(file string) error
	// SetServerCertFile sets a SSL server certificate file.
	SetServerCertFile(file string) error
	// SetRootCertFile sets a SSL root certificates.
	SetRootCertFiles(files ...string) error
	// SetServerKey sets a SSL server key.
	SetServerKey(key []byte)
	// SetServerCert sets a SSL server certificate.
	SetServerCert(cert []byte)
	// SetRootCerts sets a SSL root certificates.
	SetRootCerts(certs ...[]byte)
	// SetTLSConfig sets a TLS configuration.
	SetTLSConfig(tlsConfig *tls.Config)
	// TLSConfig returns a TLS configuration from the configuration.
	TLSConfig() (*tls.Config, error)
}

// TLSConfig represents a TLS configuration interface.
type TLSConfig interface {
	CertConfig

	// SetTLSEnabled sets a TLS enabled flag.
	SetTLSEnabled(enabled bool)
	// IsEnabled returns true if the TLS is enabled.
	IsTLSEnabled() bool
}

// Config represents a MySQL server configuration.
type Config interface {
	TLSConfig

	// SetAddress sets a listen address.
	SetAddress(host string)
	// SetPort sets a listen port.
	SetPort(port int)
	// Address returns a listen address.
	Address() string
	// Port returns a listen port.
	Port() int
}
