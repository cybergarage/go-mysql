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

//go:embed data/ssl-request-001.hex
var sslRequestMsg001 string

func TestSSLNewRequestPacket(t *testing.T) {
	req, err := protocol.NewSSLRequest()
	if err != nil {
		t.Error(err)
		return
	}
	reqBytes, err := req.Bytes()
	if err != nil {
		t.Error(err)
		return
	}

	req2, err := protocol.NewSSLRequestFromReader(bytes.NewReader(reqBytes))
	if err != nil {
		t.Error(err)
		return
	}
	req2Bytes, err := req2.Bytes()
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(reqBytes, req2Bytes) {
		t.Error("Invalid SSL request")
	}
}

func TestSSLRequestPacket(t *testing.T) {
	type expected struct {
	}
	for _, test := range []struct {
		name string
		data string
		expected
	}{
		{
			"ssl-request-001",
			sslRequestMsg001,
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

			pkt, err := protocol.NewSSLRequestFromReader(reader)
			if err != nil {
				t.Error(err)
				return
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
