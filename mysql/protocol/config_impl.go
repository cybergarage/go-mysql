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
	"fmt"

	"github.com/cybergarage/go-authenticator/auth/tls"
)

const (
	DefaultAddr        = ""
	DefaultProductName = "mysql"
)

// Config stores server configuration parammeters.
type config struct {
	addr string
	port int
	tls.CertConfig
	tlsEnabled     bool
	productName    string
	productVersion string
	capability     Capability
	serverStatus   ServerStatus
	autuPluginName string
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() Config {
	config := &config{
		addr:           DefaultAddr,
		port:           DefaultPort,
		CertConfig:     tls.NewCertConfig(),
		tlsEnabled:     true,
		productName:    DefaultProductName,
		productVersion: "",
		capability:     DefaultServerCapability,
		serverStatus:   DefaultServerStatus,
		autuPluginName: DefaultAuthPluginName,
	}
	return config
}

// SetProuctName sets a product name to the configuration.
func (config *config) SetProductName(v string) {
	config.productName = v
}

// SetProductVersion sets a product version to the configuration.
func (config *config) SetProductVersion(v string) {
	config.productVersion = v
}

// SetAddress sets a listen address to the configuration.
func (config *config) SetAddress(addr string) {
	config.addr = addr
}

// SetPort sets a listen port to the configuration.
func (config *config) SetPort(port int) {
	config.port = port
}

// Address returns the listen address from the configuration.
func (config *config) Address() string {
	return config.addr
}

// Port returns the listen port from the configuration.
func (config *config) Port() int {
	return config.port
}

// ProductName returns the product name from the configuration.
func (config *config) ProductName() string {
	return config.productName
}

// ProductVersion returns the product version from the configuration.
func (config *config) ProductVersion() string {
	return config.productVersion
}

// ServerVersion returns the server version for the handshake.
func (config *config) ServerVersion() string {
	ver := SupportVersion
	if 0 < len(config.productName) {
		ver = fmt.Sprintf("%s-%s", ver, config.productName)
	}
	if 0 < len(config.productVersion) {
		ver = fmt.Sprintf("%s-%s", ver, config.productVersion)
	}
	return ver
}

// SetCapability sets the capability flags to the configuration.
func (config *config) SetCapability(c Capability) {
	config.capability = c
}

// Capability returns the capability flags from the configuration.
func (config *config) Capability() Capability {
	return config.capability
}

// SetServerStatus sets the server status to the configuration.
func (config *config) SetServerStatus(status ServerStatus) {
	config.serverStatus = status
}

// ServerStatus returns the server status from the configuration.
func (config *config) ServerStatus() ServerStatus {
	return config.serverStatus
}

// SetAuthPluginName sets the auth plugin name to the configuration.
func (config *config) SetAuthPluginName(v string) {
	config.autuPluginName = v
}

// AuthPluginName returns the auth plugin name from the configuration.
func (config *config) AuthPluginName() string {
	return config.autuPluginName
}

// SetTLSEnabled sets a TLS enabled flag.
func (config *config) SetTLSEnabled(enabled bool) {
	config.tlsEnabled = enabled
}

// IsEnabled returns true if the TLS is enabled.
func (config *config) IsTLSEnabled() bool {
	return config.tlsEnabled
}
