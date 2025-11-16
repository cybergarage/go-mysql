# Copyright (C) 2025 The go-mysql Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# BenchBase Setup Makefile
# Usage Guide BenchBase: https://github.com/cmu-db/benchbase?tab=readme-ov-file#usage-guide

# Configuration variables
BENCHBASE_GITURL := https://github.com/cybergarage
BENCHBASE_DIR := benchbase
BENCHBASE_TARGET := mysql
BENCHBASE_JAR := benchbase.jar
BENCHBASE_TARGET_DIR := $(BENCHBASE_DIR)/target/benchbase-$(BENCHBASE_TARGET)
BENCHBASE_TARGET_JAR := $(BENCHBASE_TARGET_DIR)/$(BENCHBASE_JAR)

# Default target
.DEFAULT_GOAL := setup

# Check if Java is available
.PHONY: check-java
check-java:
	@command -v java >/dev/null 2>&1 || { echo "Error: Java is not installed or not in PATH"; exit 1; }
	@command -v git >/dev/null 2>&1 || { echo "Error: Git is not installed or not in PATH"; exit 1; }

# Setup BenchBase (main target)
.PHONY: setup
setup: check-java $(BENCHBASE_JAR)
	@echo "BenchBase setup completed."
	java -jar $(BENCHBASE_JAR) -h

# Build BenchBase JAR if it doesn't exist
$(BENCHBASE_JAR): $(BENCHBASE_DIR)
	@echo "Building BenchBase..."
	cd $(BENCHBASE_DIR) && ./mvnw package -P $(BENCHBASE_TARGET)
	cd $(BENCHBASE_DIR)/target && tar xvzf benchbase-$(BENCHBASE_TARGET).tgz
	cp $(BENCHBASE_TARGET_JAR) .
	mv $(BENCHBASE_TARGET_DIR)/lib .
	mv $(BENCHBASE_TARGET_DIR)/data .
	mv $(BENCHBASE_TARGET_DIR)/config .
	mv $(BENCHBASE_TARGET_DIR)/scripts .

# Clone BenchBase repository if it doesn't exist
$(BENCHBASE_DIR):
	@echo "Cloning BenchBase repository..."
	git clone --depth 1 $(BENCHBASE_GITURL)/benchbase.git

# Clean all generated files and directories
.PHONY: clean
clean:
	@echo "Cleaning BenchBase files..."
	rm -rf $(BENCHBASE_DIR) lib data config scripts $(BENCHBASE_JAR)

# Force rebuild (clean + setup)
.PHONY: rebuild
rebuild: clean setup

# Run BenchBase help
.PHONY: help
help: $(BENCHBASE_JAR)
	java -jar $(BENCHBASE_JAR) -h

# Run TPCC benchmark (requires server to be running)
.PHONY: run-tpcc
run-tpcc: $(BENCHBASE_JAR)
	java -jar $(BENCHBASE_JAR) -b tpcc -c tpcc_config.xml --create=true --load=true --execute=true

# Show current status
.PHONY: status
status:
	@echo "BenchBase Setup Status:"
	@echo "  JAR file: $(shell [ -f $(BENCHBASE_JAR) ] && echo "✓ Present" || echo "✗ Missing")"
	@echo "  lib directory: $(shell [ -d lib ] && echo "✓ Present" || echo "✗ Missing")"
	@echo "  data directory: $(shell [ -d data ] && echo "✓ Present" || echo "✗ Missing")"
	@echo "  config directory: $(shell [ -d config ] && echo "✓ Present" || echo "✗ Missing")"
	@echo "  scripts directory: $(shell [ -d scripts ] && echo "✓ Present" || echo "✗ Missing")"
	@echo "  Java: $(shell command -v java >/dev/null 2>&1 && echo "✓ Available" || echo "✗ Not found")"
	@echo "  Git: $(shell command -v git >/dev/null 2>&1 && echo "✓ Available" || echo "✗ Not found")"

# Show make targets
.PHONY: list
list:
	@echo "Available targets:"
	@echo "  setup    - Setup BenchBase (default)"
	@echo "  clean    - Remove all BenchBase files"
	@echo "  rebuild  - Clean and setup BenchBase"
	@echo "  help     - Show BenchBase help"
	@echo "  run-tpcc - Run TPCC benchmark"
	@echo "  status   - Show current setup status"
	@echo "  list     - Show this help message"