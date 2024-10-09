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

// config stores server configuration parammeters.
type config struct {
	address  string
	port     int
	database string
}

// NewDefaultConfig returns a default configuration instance.
func NewDefaultConfig() Config {
	config := &config{
		address:  defaultAddr,
		port:     defaultPort,
		database: "",
	}
	return config
}

// SetAddress sets a listen address.
func (config *config) SetAddress(host string) {
	config.address = host
}

// SetPort sets a listen port.
func (config *config) SetPort(port int) {
	config.port = port
}

// SetDatabase sets a host database.
func (config *config) SetDatabase(db string) {
	config.database = db
}

// Address returns a listen address.
func (config *config) Address() string {
	return config.address
}

// Port returns a listen port.
func (config *config) Port() int {
	return config.port
}

// Database returns a host database.
func (config *config) Database() string {
	return config.database
}
