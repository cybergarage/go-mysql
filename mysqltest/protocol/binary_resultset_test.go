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

func TestBinaryResultSetPacket(t *testing.T) {
	type expected struct {
	}
	for _, test := range []struct {
		name string
		protocol.Capability
		protocol.ServerStatus
		expected
	}{
		{
			"data/binary-resultset-001.hex",
			protocol.DefaultServerCapability,
			protocol.DefaultServerStatus,
			expected{},
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

			pkt, err := protocol.NewBinaryResultSetFromReader(
				reader,
				protocol.WithBinaryResultSetCapability(test.Capability),
				protocol.WithBinaryResultSetServerStatus(test.ServerStatus),
			)
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
