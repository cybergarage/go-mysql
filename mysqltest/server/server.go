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

package server

import (
	"crypto/tls"

	server "github.com/cybergarage/go-mysql/examples/go-mysqld/server"
	"github.com/cybergarage/go-mysql/mysql/auth"
)

const (
	serverKey  = "../certs/key.pem"
	serverCert = "../certs/cert.pem"
	rootCert   = "../certs/root_cert.pem"
)

// Server represents a test server.
type Server struct {
	*server.Server
	credStore map[string]auth.Credential
}

// NewServer returns a test server instance.
func NewServer() *Server {
	server := &Server{
		Server:    server.NewServer(),
		credStore: make(map[string]auth.Credential),
	}

	server.SetServerKeyFile(serverKey)
	server.SetServerCertFile(serverCert)
	server.SetRootCertFiles(rootCert)
	server.SetClientAuthType(tls.RequireAndVerifyClientCert)

	ca, err := auth.NewCertificateAuthenticator(
		auth.WithCommonNameRegexp("localhost"),
	)
	if err != nil {
		server.SetCertificateAuthenticator(ca)
	}

	return server
}

// SetCredential sets a credential.
func (server *Server) SetCredential(cred auth.Credential) {
	server.credStore[cred.Username()] = cred
}

// LookupCredential looks up a credential.
func (server *Server) LookupCredential(q auth.Query) (auth.Credential, bool, error) {
	user := q.Username()
	cred, ok := server.credStore[user]
	return cred, ok, nil
}
