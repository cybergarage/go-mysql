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
	_ "embed"
	"strings"
	"testing"

	"github.com/cybergarage/go-mysql/mysql/protocol"
)

func TestServer(t *testing.T) {
	pkt := protocol.NewHandshake()

	if pkt.ProtocolVersion() != protocol.ProtocolVersion10 {
		t.Errorf("expected %d, got %d", protocol.ProtocolVersion10, pkt.ProtocolVersion())
	}

	if !strings.HasPrefix(pkt.ServerVersion(), protocol.SupportVersion) {
		t.Errorf("expected %s, got %s", protocol.SupportVersion, pkt.ServerVersion())
	}
}
