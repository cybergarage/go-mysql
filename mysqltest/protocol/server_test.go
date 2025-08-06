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
	"testing"

	"github.com/cybergarage/go-mysql/mysql"
	"github.com/cybergarage/go-mysql/mysql/protocol"
)

func TestServer(t *testing.T) {
	server := protocol.NewServer()

	err := server.Start()
	if err != nil {
		t.Errorf("expected nil, got %v", err)
		return
	}

	defer func() {
		err := server.Stop()
		if err != nil {
			t.Errorf("expected nil, got %v", err)
			return
		}
	}()

	client := mysql.NewClient()

	err = client.Open()
	if err != nil {
		t.Error(err)
		return
	}

	err = client.Ping()
	if err != nil {
		t.Error(err)
	}

	err = client.Close()
	if err != nil {
		t.Error(err)
	}
}
