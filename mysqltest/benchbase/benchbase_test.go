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

package benchbase

import (
	"os"
	"testing"

	"github.com/cybergarage/go-logger/log"
	"github.com/cybergarage/go-mysql/mysqltest/server"
	"github.com/cybergarage/go-sqltest/sqltest/benchbase"
	// BenchBase integration placeholder: BenchBase is Java-based.
	// Actual invocation would require spawning BenchBase workload runner via CLI.
	// We keep this test minimal and focused on server startup lifecycle.
)

// TestBenchBase starts a test MySQL-compatible server and performs basic lifecycle checks.
func TestBenchBase(t *testing.T) {
	// Enable verbose debug logging to observe benchmark progress.
	log.EnableStdoutDebug(true)

	// Skip early if BenchBase tooling is not available on this system.
	if !benchbase.IsInstalled() {
		t.Skip("BenchBase is not installed; skipping test")
		return
	}

	wkdir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Working directory: %s", wkdir)

	// Start server
	srv := server.NewServer()
	if err := srv.Start(); err != nil {
		t.Fatalf("failed to start server: %v", err)
	}
	defer srv.Stop()

	// List of benches to execute; expand as needed.
	// Common BenchBase benches include: tpcc, tatp, smallbank, ycsb, epinions, etc.
	benches := []string{
		"tpcc",
	}

	// Each bench is run as a subtest for isolated reporting in go test output.
	for _, bench := range benches {
		// shadow for closure capture
		t.Run(bench, func(t *testing.T) {
			benchbase.RunWorkload(t, bench)
		})
	}
}
