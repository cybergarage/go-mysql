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

func TestOKPacket(t *testing.T) {
	type expected struct {
		seqID        protocol.SequenceID
		affectedRows uint64
		lastInsertID uint64
	}
	for _, test := range []struct {
		name string
		opts []protocol.OKOption
		expected
	}{
		{
			"data/ok-001.hex",
			[]protocol.OKOption{protocol.WithOKCapability(protocol.ClientProtocol41)},
			expected{
				seqID:        protocol.SequenceID(2),
				affectedRows: 0,
				lastInsertID: 0,
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

			pkt, err := protocol.NewOKFromReader(reader, test.opts...)
			if err != nil {
				t.Error(err)
			}

			if pkt.SequenceID() != test.expected.seqID {
				t.Errorf("expected %d, got %d", test.expected.seqID, pkt.SequenceID())
			}

			if pkt.AffectedRows() != test.expected.affectedRows {
				t.Errorf("expected %d, got %d", test.expected.affectedRows, pkt.AffectedRows())
			}

			if pkt.LastInsertID() != test.expected.lastInsertID {
				t.Errorf("expected %d, got %d", test.expected.lastInsertID, pkt.LastInsertID())
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

func TestOKPacketWithClientProtocol41Capability(t *testing.T) {
	pkt, err := protocol.NewOK(
		protocol.WithOKCapability(protocol.ClientProtocol41),
		protocol.WithOKSecuenceID(protocol.SequenceID(2)),
	)
	if err != nil {
		t.Fatal(err)
	}

	pktBytes, err := pkt.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	const headerLength = 4
	const protocol41OKPayloadLength = 7
	if got := len(pktBytes) - headerLength; got != protocol41OKPayloadLength {
		t.Fatalf("expected payload length %d, got %d", protocol41OKPayloadLength, got)
	}

	reader := bytes.NewReader(pktBytes)
	parsedPkt, err := protocol.NewOKFromReader(reader, protocol.WithOKCapability(protocol.ClientProtocol41))
	if err != nil {
		t.Fatal(err)
	}

	if parsedPkt.SequenceID() != protocol.SequenceID(2) {
		t.Errorf("expected sequence ID %d, got %d", protocol.SequenceID(2), parsedPkt.SequenceID())
	}
}
