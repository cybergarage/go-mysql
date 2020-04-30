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

const (
	defaultPort = 3306
)

// Config stores server configuration parammeters.
type Config struct {
	Addr string
	Port int
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() *Config {
	config := &Config{
		Addr: "",
		Port: defaultPort,
	}
	return config
}

// SetAddress sets a bind interface address.
func (config *Config) SetAddress(addr string) {
	config.Addr = addr
}

// SetPort sets a listen port.
func (config *Config) SetPort(port int) {
	config.Port = port
}

// GetPort returns a listent port.
func (config *Config) GetPort() int {
	return config.Port
}
