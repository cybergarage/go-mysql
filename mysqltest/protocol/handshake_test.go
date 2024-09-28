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
	// Protocol: 10
	// Version: 5.7.9-vitess-12.0.6
	// Thread ID: 1
	// Salt: Dn+\vB\x0En\x03
	// Server Capabilities: 0xa20f
	// Server Language: utf8mb3 COLLATE utf8mb3_general_ci (33)
	// Server Status: 0x0000
	// Extended Server Capabilities: 0x013b
	// Authentication Plugin Length: 21
	// Unused: 00000000000000000000
	// Salt: 2\x1Eg\ayx&\x18R\x1D\x01P
	// Authentication Plugin: mysql_native_password

	type expected struct {
		seqID       protocol.SequenceID
		protocolVer protocol.ProtocolVersion
		serverVer   string
		conID       uint32
		capFlags    protocol.CapabilityFlag
		charSet     protocol.CharacterSet
	}
	for _, test := range []struct {
		name string
		data string
		expected
	}{
		{
			"handshake",
			handshakeMsg001,
			expected{
				seqID:       protocol.SequenceID(0),
				protocolVer: protocol.ProtocolVersion10,
				serverVer:   "5.7.9-vitess-12.0.6",
				conID:       1,
				capFlags:    protocol.CapabilityFlag(0),
				charSet:     protocol.CharacterSet(protocol.CharacterSetUTF8),
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			reader, err := hexdump.NewReaderWithString(test.data)
			if err != nil {
				t.Error(err)
				return
			}
			msg, err := protocol.NewHandshakeFromReader(reader)
			if err != nil {
				t.Error(err)
			}

			if msg.SequenceID() != test.expected.seqID {
				t.Errorf("expected %d, got %d", test.expected.seqID, msg.SequenceID())
			}

			if msg.ProtocolVersion() != test.expected.protocolVer {
				t.Errorf("expected %d, got %d", test.expected.protocolVer, msg.ProtocolVersion())
			}

			if msg.ServerVersion() != test.expected.serverVer {
				t.Errorf("expected %s, got %s", test.expected.serverVer, msg.ServerVersion())
			}

			if msg.ConnectionID() != test.expected.conID {
				t.Errorf("expected %d, got %d", test.expected.conID, msg.ConnectionID())
			}

			if msg.CapabilityFlags() != test.expected.capFlags {
				t.Errorf("expected %d, got %d", test.expected.capFlags, msg.CapabilityFlags())
			}
		})
	}
}
