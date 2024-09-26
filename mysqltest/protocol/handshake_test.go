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
	_ "embed"
	"testing"

	"github.com/cybergarage/go-logger/log/hexdump"
	"github.com/cybergarage/go-mysql/mysql/protocol"
)

//go:embed data/handshake001.hex
var handshakeMsg001 string

func TestHandshakeMessage(t *testing.T) {
	for _, test := range []struct {
		name string
		data string
	}{
		{"handshake", handshakeMsg001},
	} {
		t.Run(test.name, func(t *testing.T) {
			reader, err := hexdump.NewReaderWithString(test.data)
			if err != nil {
				t.Error(err)
				return
			}
			_, err = protocol.NewHandshakeWith(reader)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
