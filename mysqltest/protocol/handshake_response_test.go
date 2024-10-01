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

//go:embed data/handshake-response-001.hex
var handshakeResponseMsg001 string

func TestHandshakeResponseMessage(t *testing.T) {
	// Packet Length: 89
	// Packet Number: 1
	// Login Request
	//     Client Capabilities: 0xa28d
	//         .... .... .... ...1 = Long Password: Set
	//         .... .... .... ..0. = Found Rows: Not set
	//         .... .... .... .1.. = Long Column Flags: Set
	//         .... .... .... 1... = Connect With Database: Set
	//         .... .... ...0 .... = Don't Allow database.table.column: Not set
	//         .... .... ..0. .... = Can use compression protocol: Not set
	//         .... .... .0.. .... = ODBC Client: Not set
	//         .... .... 1... .... = Can Use LOAD DATA LOCAL: Set
	//         .... ...0 .... .... = Ignore Spaces before '(': Not set
	//         .... ..1. .... .... = Speaks 4.1 protocol (new flag): Set
	//         .... .0.. .... .... = Interactive Client: Not set
	//         .... 0... .... .... = Switch to SSL after handshake: Not set
	//         ...0 .... .... .... = Ignore sigpipes: Not set
	//         ..1. .... .... .... = Knows about transactions: Set
	//         .0.. .... .... .... = Speaks 4.1 protocol (old flag): Not set
	//         1... .... .... .... = Can do 4.1 authentication: Set
	//     Extended Client Capabilities: 0x000a
	//     MAX Packet: 0
	//     Collation: utf8mb4 COLLATE utf8mb4_general_ci (45)
	//     Unused: 0000000000000000000000000000000000000000000000
	//     Username: skonno
	//     Schema: sqltest1727254524366662000
	//     Client Auth Plugin: mysql_native_password

	type expected struct {
		seqID          protocol.SequenceID
		protocolVer    protocol.ProtocolVersion
		serverVer      string
		conID          uint32
		capFlags       protocol.CapabilityFlag
		charSet        protocol.CharacterSet
		statusFlags    protocol.StatusFlag
		authPluginName string
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
				seqID:          protocol.SequenceID(0),
				protocolVer:    protocol.ProtocolVersion10,
				serverVer:      "5.7.9-vitess-12.0.6",
				conID:          1,
				capFlags:       protocol.CapabilityFlag(0),
				charSet:        protocol.CharacterSet(protocol.CharacterSetUTF8),
				statusFlags:    protocol.StatusFlag(0),
				authPluginName: "mysql_native_password",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			reader, err := hexdump.NewReaderFromHexdumpString(test.data)
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

			if msg.CharacterSet() != test.expected.charSet {
				t.Errorf("expected %d, got %d", test.expected.charSet, msg.CharacterSet())
			}

			if msg.StatusFlags() != test.expected.statusFlags {
				t.Errorf("expected %d, got %d", test.expected.statusFlags, msg.StatusFlags())
			}

			if msg.AuthPluginName() != test.expected.authPluginName {
				t.Errorf("expected %s, got %s", test.expected.authPluginName, msg.AuthPluginName())
			}
		})
	}
}
