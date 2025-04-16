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

func TestQueryPacket(t *testing.T) {
	type expected struct {
		seqID protocol.SequenceID
		query string
	}
	for _, test := range []struct {
		name     string
		capFlags protocol.Capability
		expected
	}{
		{
			"data/query-001.hex",
			protocol.ClientQueryAttributes,
			expected{
				seqID: protocol.SequenceID(0),
				query: "select @@version_comment limit 1",
			},
		},
		{
			"data/query-002.hex",
			0,
			expected{
				seqID: protocol.SequenceID(0),
				query: "CREATE DATABASE IF NOT EXISTS sqltest1727254524366662000",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			testData, err := testEmbedPacketFiles.ReadFile(test.name)
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

			opts := []protocol.QueryOption{
				protocol.WithQueryCapability(test.capFlags),
			}
			pkt, err := protocol.NewQueryFromReader(reader, opts...)
			if err != nil {
				t.Error(err)
				return
			}

			if pkt.SequenceID() != test.expected.seqID {
				t.Errorf("expected %d, got %d", test.expected.seqID, pkt.SequenceID())
			}

			if pkt.Query() != test.expected.query {
				t.Errorf("expected %s, got %s", test.expected.query, pkt.Query())
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
