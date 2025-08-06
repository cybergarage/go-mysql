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

func TestStmtExecutePacket(t *testing.T) {
	type expected struct {
		stmtID protocol.StatementID
	}
	for _, test := range []struct {
		name      string
		numParams uint16
		expected
	}{
		{
			"data/stmt-execute-001.hex",
			1,
			expected{
				stmtID: 1,
			},
		},
		{
			"data/stmt-execute-002.hex",
			1,
			expected{
				stmtID: 2,
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

			pkt, err := protocol.NewStmtExecuteFromReader(reader,
				protocol.WithStmtExecuteNumParams(test.numParams),
			)
			if err != nil {
				t.Error(err)
				return
			}

			// Compare the packet fields

			if pkt.StatementID() != test.stmtID {
				t.Errorf("stmtID = %d, want %d", pkt.StatementID(), test.stmtID)
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
