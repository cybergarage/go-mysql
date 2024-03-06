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
	defaultAddr = ""
	defaultPort = 3306
)

// Config stores server configuration parammeters.
type Config struct {
	Address  string
	Port     int
	Database string
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() *Config {
	config := &Config{
		Address:  defaultAddr,
		Port:     defaultPort,
		Database: "",
	}
	return config
}

// SetAddress sets a listen address.
func (config *Config) SetAddress(host string) {
	config.Address = host
}

// SetPort sets a listen port.
func (config *Config) SetPort(port int) {
	config.Port = port
}

// SetDatabase sets a host database.
func (config *Config) SetDatabase(db string) {
	config.Database = db
}
