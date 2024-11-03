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

//go:embed data/handshake-response-003.hex
var handshakeResponseMsg003 string

func TestHandshakeResponsePacket(t *testing.T) {
	type expected struct {
		capFlags   protocol.Capability
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
				capFlags:   protocol.Capability(0x000aa28d),
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
				capFlags:   protocol.Capability(0x000fa68d),
				maxPkt:     0,
				charSet:    0,
				username:   "",
				authRes:    "",
				database:   "test",
				pluginName: "mysql_native_password",
				attrs:      map[string]string{},
				zstdLevel:  0,
			},
		},
		{
			"handshake-response-003",
			handshakeResponseMsg003,
			expected{
				capFlags:   protocol.Capability(0x001ea285),
				maxPkt:     0,
				charSet:    0,
				username:   "root",
				authRes:    "",
				database:   "",
				pluginName: "mysql_native_password",
				attrs: map[string]string{
					"_os":             "debian6.0",
					"_client_name":    "libmysql",
					"_pid":            "22344",
					"_client_version": "5.6.6-m9",
					"_platform":       "x86_64",
					"foo":             "bar",
				},
				zstdLevel: 0,
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
				return
			}

			if test.expected.capFlags != 0 {
				if pkt.Capability() != test.expected.capFlags {
					t.Errorf("expected %04X, got %04X", test.expected.capFlags, pkt.Capability())
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

			for key, value := range test.expected.attrs {
				v, ok := pkt.LookupAttribute(key)
				if !ok || v != value {
					t.Errorf("expected %s, got %s", value, v)
				}
			}

			// Compare the packet bytes

			msgBytes, err := pkt.Bytes()
			if err != nil {
				t.Error(err)
				return
			}

			if !bytes.Equal(msgBytes, testBytes) {
				HexdumpErrors(t, testBytes, msgBytes)
			}
		})
	}
}
