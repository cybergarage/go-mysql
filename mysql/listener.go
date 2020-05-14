// Copyright (C) 2020 Satoshi Konno. All rights reserved.
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
	"time"

	vitess "vitess.io/vitess/go/mysql"
)

// Listener is the MySQL server protocol listener.
type Listener struct {
	*vitess.Listener
}

// NewListener creates a new listener.
func NewListener(protocol, address string, authServer AuthHandler, handler QueryHandler, connReadTimeout time.Duration, connWriteTimeout time.Duration, proxyProtocol bool) (*Listener, error) {
	l, err := vitess.NewListener(protocol, address, authServer, handler, connReadTimeout, connWriteTimeout, proxyProtocol)
	if err != nil {
		return nil, err
	}
	return &Listener{Listener: l}, nil
}
