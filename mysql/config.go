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

// ServerConfig stores server configuration parameters.
type ServerConfig interface {
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

// ClientConfig stores client configuration parameters.
type ClientConfig interface {
	ServerConfig
	// SetDatabase sets a host database.
	SetDatabase(db string)
	// Database returns a host database.
	Database() string
}

// Config stores client configuration parameters.
type Config = ClientConfig
