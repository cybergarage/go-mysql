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
	v2pg "github.com/cybergarage/go-mysql/examples/go-mysqld/v2"
	vitesspg "github.com/cybergarage/go-mysql/examples/go-mysqld/vitess"
	"github.com/cybergarage/go-tracing/tracer"
)

// Server represents a MySQL-compatible server interface.
type Server interface {
	ServerConfig
	tracer.Tracer
	Start() error
	Stop() error
	Restart() error
}

// NewServer creates a new server instance.
func NewServer() Server {
	v2Server := v2pg.NewServer()
	v2Server.SetProductName(PackageName)
	v2Server.SetProductVersion(Version)

	server := vitesspg.NewServer()
	server.SetProductName(PackageName)
	server.SetProductVersion(Version)

	return v2Server
	// return server
}
