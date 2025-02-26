// Copyright (C) 2025 The go-mysql Authors. All rights reserved.
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

package sysbench

import (
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysqltest/server"
	"github.com/cybergarage/go-sqltest/sqltest/sysbench"
)

func TestSysbench(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	cfg := NewDefaultConfig()

	server := server.NewServer()
	err := server.Start()
	if err != nil {
		t.Error(err)
		return
	}
	defer server.Stop()

	cmds := []string{
		sysbench.OltpReadWrite,
	}

	for _, cmd := range cmds {
		t.Run(cmd, func(t *testing.T) {
			sysbench.RunCommand(t, cmd, cfg.Config)
		})
	}
}
