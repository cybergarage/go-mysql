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

func TestErrPacket(t *testing.T) {
	type expected struct {
		seqID       protocol.SequenceID
		errCode     uint16
		stateMarker string
		state       protocol.ErrCode
		errMsg      string
	}
	for _, test := range []struct {
		name string
		expected
	}{
		{
			"data/err-001.hex",
			expected{
				seqID:       protocol.SequenceID(1),
				errCode:     0x0448,
				stateMarker: "#",
				state:       "HY000",
				errMsg:      "No tables used",
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

			pkt, err := protocol.NewERRFromReader(reader, protocol.WithERRCapability(protocol.ClientProtocol41))
			if err != nil {
				t.Error(err)
			}

			pkt.SetCapabilityEnabled(protocol.ClientProtocol41)

			if pkt.SequenceID() != test.seqID {
				t.Errorf("expected %d, got %d", test.seqID, pkt.SequenceID())
			}

			if pkt.Code() != test.errCode {
				t.Errorf("expected %d, got %d", test.errCode, pkt.Code())
			}

			if pkt.StateMarker() != test.stateMarker {
				t.Errorf("expected %s, got %s", test.stateMarker, pkt.StateMarker())
			}

			if pkt.State() != test.state {
				t.Errorf("expected %s, got %s", test.state, pkt.State())
			}

			if pkt.ErrMsg() != test.errMsg {
				t.Errorf("expected %s, got %s", test.errMsg, pkt.ErrMsg())
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
