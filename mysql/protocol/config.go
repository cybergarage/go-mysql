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
)

const (
	DefaultAddr        = ""
	DefaultProductName = "mysql"
)

// Config stores server configuration parammeters.
type Config struct {
	addr string
	port int
	*TLSConf
	productName    string
	productVersion string
	capability     Capability
	autuPluginName string
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() *Config {
	config := &Config{
		addr:           DefaultAddr,
		port:           DefaultPort,
		TLSConf:        NewTLSConf(),
		productName:    DefaultProductName,
		productVersion: "",
		capability:     DefaultServerCapabilities,
		autuPluginName: DefaultAuthPluginName,
	}
	return config
}

// SetProuctName sets a product name to the configuration.
func (config *Config) SetProductName(v string) {
	config.productName = v
}

// SetProductVersion sets a product version to the configuration.
func (config *Config) SetProductVersion(v string) {
	config.productVersion = v
}

// SetAddress sets a listen address to the configuration.
func (config *Config) SetAddress(addr string) {
	config.addr = addr
}

// SetPort sets a listen port to the configuration.
func (config *Config) SetPort(port int) {
	config.port = port
}

// Address returns the listen address from the configuration.
func (config *Config) Address() string {
	return config.addr
}

// Port returns the listen port from the configuration.
func (config *Config) Port() int {
	return config.port
}

// ProductName returns the product name from the configuration.
func (config *Config) ProductName() string {
	return config.productName
}

// ProductVersion returns the product version from the configuration.
func (config *Config) ProductVersion() string {
	return config.productVersion
}

// ServerVersion returns the server version for the handshake.
func (config *Config) ServerVersion() string {
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
func (config *Config) SetCapability(c Capability) {
	config.capability = c
}

// Capability returns the capability flags from the configuration.
func (config *Config) Capability() Capability {
	return config.capability
}

// SetAuthPluginName sets the auth plugin name to the configuration.
func (config *Config) SetAuthPluginName(v string) {
	config.autuPluginName = v
}

// AuthPluginName returns the auth plugin name from the configuration.
func (config *Config) AuthPluginName() string {
	return config.autuPluginName
}
