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

//go:embed data/eof-001.hex
var eofPkt001 string

func TestEOFPacket(t *testing.T) {
	type expected struct {
		seqID    protocol.SequenceID
		warnings uint16
		status   protocol.StatusFlag
	}
	for _, test := range []struct {
		name string
		data string
		expected
	}{
		{
			"eof001",
			eofPkt001,
			expected{
				seqID:    protocol.SequenceID(5),
				warnings: 0,
				status:   0x0002,
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

			pkt, err := protocol.NewEOFFromReader(reader, protocol.WithEofCapability(protocol.ClientProtocol41))
			if err != nil {
				t.Error(err)
			}
			pkt.SetCapabilityEnabled(protocol.ClientProtocol41)

			if pkt.SequenceID() != test.expected.seqID {
				t.Errorf("expected %d, got %d", test.expected.seqID, pkt.SequenceID())
			}

			if pkt.Warnings() != test.expected.warnings {
				t.Errorf("expected %d, got %d", test.expected.warnings, pkt.Warnings())
			}

			if pkt.Status() != test.expected.status {
				t.Errorf("expected %d, got %d", test.expected.status, pkt.Status())
			}

			// Compare the packet bytes

			pktBytes, err := pkt.Bytes()
			if err != nil {
				t.Error(err)
				return
			}

			if !bytes.Equal(pktBytes, testBytes) {
				t.Errorf("expected %v, got %v", testBytes, pktBytes)
			}
		})
	}
}
