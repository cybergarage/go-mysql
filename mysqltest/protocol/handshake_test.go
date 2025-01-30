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
	"strings"
	"testing"

	"github.com/cybergarage/go-logger/log/hexdump"
	"github.com/cybergarage/go-mysql/mysql/protocol"
)

func TestHandshakePacket(t *testing.T) {
	// Packet Length: 87
	// Packet Number: 0
	// Server Greeting
	// 	Protocol: 10
	// 	Version: 5.7.9-vitess-12.0.6
	// 	Thread ID: 1
	// 	Salt: Dn+\vB\x0En\x03
	// 	Server Capabilities: 0xa20f
	// 		.... .... .... ...1 = Long Password: Set
	// 		.... .... .... ..1. = Found Rows: Set
	// 		.... .... .... .1.. = Long Column Flags: Set
	// 		.... .... .... 1... = Connect With Database: Set
	// 		.... .... ...0 .... = Don't Allow database.table.column: Not set
	// 		.... .... ..0. .... = Can use compression protocol: Not set
	// 		.... .... .0.. .... = ODBC Client: Not set
	// 		.... .... 0... .... = Can Use LOAD DATA LOCAL: Not set
	// 		.... ...0 .... .... = Ignore Spaces before '(': Not set
	// 		.... ..1. .... .... = Speaks 4.1 protocol (new flag): Set
	// 		.... .0.. .... .... = Interactive Client: Not set
	// 		.... 0... .... .... = Switch to SSL after handshake: Not set
	// 		...0 .... .... .... = Ignore sigpipes: Not set
	// 		..1. .... .... .... = Knows about transactions: Set
	// 		.0.. .... .... .... = Speaks 4.1 protocol (old flag): Not set
	// 		1... .... .... .... = Can do 4.1 authentication: Set
	// 	Server Language: utf8mb3 COLLATE utf8mb3_general_ci (33)
	// 	Server Status: 0x0000
	// 		.... .... .... ...0 = In transaction: Not set
	// 		.... .... .... ..0. = AUTO_COMMIT: Not set
	// 		.... .... .... .0.. = Multi query / Unused: Not set
	// 		.... .... .... 0... = More results: Not set
	// 		.... .... ...0 .... = Bad index used: Not set
	// 		.... .... ..0. .... = No index used: Not set
	// 		.... .... .0.. .... = Cursor exists: Not set
	// 		.... .... 0... .... = Last row sent: Not set
	// 		.... ...0 .... .... = Database dropped: Not set
	// 		.... ..0. .... .... = No backslash escapes: Not set
	// 		.... .0.. .... .... = Metadata changed: Not set
	// 		.... 0... .... .... = Query was slow: Not set
	// 		...0 .... .... .... = PS Out Params: Not set
	// 		..0. .... .... .... = In Trans Readonly: Not set
	// 		.0.. .... .... .... = Session state changed: Not set
	// 	Extended Server Capabilities: 0x013b
	// 		.... .... .... ...1 = Multiple statements: Set
	// 		.... .... .... ..1. = Multiple results: Set
	// 		.... .... .... .0.. = PS Multiple results: Not set
	// 		.... .... .... 1... = Plugin Auth: Set
	// 		.... .... ...1 .... = Connect attrs: Set
	// 		.... .... ..1. .... = Plugin Auth LENENC Client Data: Set
	// 		.... .... .0.. .... = Client can handle expired passwords: Not set
	// 		.... .... 0... .... = Session variable tracking: Not set
	// 		.... ...1 .... .... = Deprecate EOF: Set
	// 		.... ..0. .... .... = Client can handle optional resultset metadata: Not set
	// 		.... .0.. .... .... = ZSTD Compression Algorithm: Not set
	// 		.... 0... .... .... = Query Attributes: Not set
	// 		...0 .... .... .... = Multifactor Authentication: Not set
	// 		..0. .... .... .... = Capability Extension: Not set
	// 		.0.. .... .... .... = Client verifies server's TLS/SSL certificate: Not set
	// 		0... .... .... .... = Unused: 0x0
	// 	Authentication Plugin Length: 21
	// 	Unused: 00000000000000000000
	// 	Salt: 2\x1Eg\ayx&\x18R\x1D\x01P
	// 	Authentication Plugin: mysql_native_password

	type expected struct {
		seqID          protocol.SequenceID
		protocolVer    protocol.ProtocolVersion
		serverVer      string
		conID          uint32
		capFlags       protocol.Capability
		charSet        protocol.CharSet
		serverStatus   protocol.ServerStatus
		authPluginName string
	}
	for _, test := range []struct {
		name string
		expected
	}{
		{
			"data/handshake-001.hex",
			expected{
				seqID:          protocol.SequenceID(0),
				protocolVer:    protocol.ProtocolVersion10,
				serverVer:      "5.7.9-vitess-12.0.6",
				conID:          1,
				capFlags:       protocol.Capability(0),
				charSet:        protocol.CharSet(protocol.CharSetUTF8),
				serverStatus:   protocol.ServerStatus(0),
				authPluginName: "mysql_native_password",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			testData, err := testPackettFiles.ReadFile(test.name)
			if err != nil {
				t.Error(err)
				return
			}
			testBytes, err := hexdump.NewBytesWithHexdumpBytes(testData)
			if err != nil {
				t.Error(err)
				return
			}
			reader := bytes.NewReader(testBytes)

			pkt, err := protocol.NewHandshakeFromReader(reader)
			if err != nil {
				t.Error(err)
			}

			if pkt.SequenceID() != test.expected.seqID {
				t.Errorf("expected %d, got %d", test.expected.seqID, pkt.SequenceID())
			}

			if pkt.ProtocolVersion() != test.expected.protocolVer {
				t.Errorf("expected %d, got %d", test.expected.protocolVer, pkt.ProtocolVersion())
			}

			if pkt.ServerVersion() != test.expected.serverVer {
				t.Errorf("expected %s, got %s", test.expected.serverVer, pkt.ServerVersion())
			}

			if pkt.ConnectionID() != test.expected.conID {
				t.Errorf("expected %d, got %d", test.expected.conID, pkt.ConnectionID())
			}

			if pkt.CharacterSet() != test.expected.charSet {
				t.Errorf("expected %d, got %d", test.expected.charSet, pkt.CharacterSet())
			}

			if pkt.ServerStatus() != test.expected.serverStatus {
				t.Errorf("expected %d, got %d", test.expected.serverStatus, pkt.ServerStatus())
			}

			if pkt.AuthPluginName() != test.expected.authPluginName {
				t.Errorf("expected %s, got %s", test.expected.authPluginName, pkt.AuthPluginName())
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

func TestServerHandshake(t *testing.T) {
	server := protocol.NewServer()

	err := server.Start()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
		return
	}

	defer func() {
		err := server.Stop()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
			return
		}
	}()

	conn := protocol.NewConnWith(nil)
	pkt, err := server.GenerateHandshakeForConn(conn)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
		return
	}

	if pkt.ProtocolVersion() != protocol.ProtocolVersion10 {
		t.Errorf("expected %d, got %d", protocol.ProtocolVersion10, pkt.ProtocolVersion())
	}

	if !strings.HasPrefix(pkt.ServerVersion(), protocol.SupportVersion) {
		t.Errorf("expected %s, got %s", protocol.SupportVersion, pkt.ServerVersion())
	}
}
