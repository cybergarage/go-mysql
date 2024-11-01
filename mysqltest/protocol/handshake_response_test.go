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
	"bytes"
	_ "embed"
	"testing"

	"github.com/cybergarage/go-logger/log/hexdump"
	"github.com/cybergarage/go-mysql/mysql/protocol"
)

//go:embed data/handshake-response-001.hex
var handshakeResponseMsg001 string

//go:embed data/handshake-response-002.hex
var handshakeResponseMsg002 string

func TestHandshakeResponsePacket(t *testing.T) {
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
		capFlags   protocol.CapabilityFlag
		maxPkt     uint32
		charSet    uint8
		username   string
		authRes    string
		database   string
		pluginName string
		attrs      map[string]string
		zstdLevel  uint8
	}
	for _, test := range []struct {
		name string
		data string
		expected
	}{
		{
			"handshake-response-001",
			handshakeResponseMsg001,
			expected{
				capFlags:   protocol.CapabilityFlag(0x000aa28d),
				maxPkt:     0,
				charSet:    45,
				username:   "skonno",
				authRes:    "",
				database:   "sqltest1727254524366662000",
				pluginName: "mysql_native_password",
				attrs:      map[string]string{},
				zstdLevel:  0,
			},
		},
		{
			"handshake-response-002",
			handshakeResponseMsg002,
			expected{
				capFlags:   protocol.CapabilityFlag(0x000fa68d),
				maxPkt:     0,
				charSet:    0,
				username:   "",
				authRes:    "",
				database:   "",
				pluginName: "",
				attrs:      map[string]string{},
				zstdLevel:  0,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			testBytes, err := hexdump.NewBytesWithHexdumpString(test.data)
			if err != nil {
				t.Error(err)
				return
			}
			reader := bytes.NewReader(testBytes)

			pkt, err := protocol.NewHandshakeResponseFromReader(reader)
			if err != nil {
				t.Error(err)
			}

			if test.expected.capFlags != 0 {
				if pkt.CapabilityFlags() != test.expected.capFlags {
					t.Errorf("expected %04X, got %04X", test.expected.capFlags, pkt.CapabilityFlags())
				}
			}

			if test.expected.maxPkt != 0 {
				if pkt.MaxPacketSize() != test.expected.maxPkt {
					t.Errorf("expected %d, got %d", test.expected.maxPkt, pkt.MaxPacketSize())
				}
			}

			if test.expected.charSet != 0 {
				if pkt.CharSet() != test.expected.charSet {
					t.Errorf("expected %d, got %d", test.expected.charSet, pkt.CharSet())
				}
			}

			if 0 < len(test.expected.username) {
				if pkt.Username() != test.expected.username {
					t.Errorf("expected %s, got %s", test.expected.username, pkt.Username())
				}
			}

			if 0 < len(test.expected.authRes) {
				if pkt.AuthResponse() != test.expected.authRes {
					t.Errorf("expected %s, got %s", test.expected.authRes, pkt.AuthResponse())
				}
			}

			if 0 < len(test.expected.database) {
				if pkt.Database() != test.expected.database {
					t.Errorf("expected %s, got %s", test.expected.database, pkt.Database())
				}
			}

			if 0 < len(test.expected.pluginName) {
				if pkt.ClientPluginName() != test.expected.pluginName {
					t.Errorf("expected %s, got %s", test.expected.pluginName, pkt.ClientPluginName())
				}
			}

			// Compare the packet bytes

			msgBytes, err := pkt.Bytes()
			if err != nil {
				t.Error(err)
				return
			}

			if !bytes.Equal(msgBytes, testBytes) {
				t.Errorf("expected %v, got %v", testBytes, msgBytes)
			}
		})
	}
}
