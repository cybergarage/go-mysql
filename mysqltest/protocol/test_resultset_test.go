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

//go:embed data/text-resultset-001.hex
var textResultSetPkt001 string

//go:embed data/text-resultset-002.hex
var textResultSetPkt002 string

func TestTextResultSetPacket(t *testing.T) {
	type expected struct {
	}
	for _, test := range []struct {
		name     string
		data     string
		capFlags protocol.Capability
		expected
	}{
		{
			"text-resultset-001",
			textResultSetPkt001,
			(protocol.ClientProtocol41 | protocol.ClientQueryAttributes),
			expected{},
		},
		{
			"text-resultset-002",
			textResultSetPkt002,
			(protocol.ClientProtocol41 | protocol.ClientQueryAttributes),
			expected{},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			testBytes, err := hexdump.NewBytesWithHexdumpString(test.data)
			if err != nil {
				t.Error(err)
				return
			}
			reader := bytes.NewReader(testBytes)

			opts := []protocol.TextResultSetOption{
				protocol.WithTextResultSetCapabilities(test.capFlags),
			}
			pkt, err := protocol.NewTextResultSetFromReader(reader, opts...)
			if err != nil {
				t.Error(err)
				return
			}

			// Compare the packet bytes

			pktBytes, err := pkt.Bytes()
			if err != nil {
				t.Error(err)
				return
			}

			if !bytes.Equal(pktBytes, testBytes) {
				HexdumpErrors(t, testBytes, pktBytes)
			}
		})
	}
}
