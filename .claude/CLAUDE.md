# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

go-mysql is a Go framework for building MySQL-compatible servers. It handles the MySQL wire protocol, authentication, and SQL parsing automatically — developers implement the executor interfaces to handle DDL/DML queries.

## Commands

```bash
make format      # Regenerate mysql/version.go and run gofmt
make lint        # Run go vet + golangci-lint
make test        # Lint + full test suite with coverage (mysql-cover.out, mysql-cover.html)
make build       # Build the example go-mysqld binary
make run         # Install and start go-mysqld example server
make sysbench    # Run sysbench integration test only
```

Run a single test:
```bash
go test ./mysql/... -run TestFunctionName -v
go test ./mysql/protocol/... -run TestHandshake -v
```

## Architecture

The framework has four layers:

1. **Network layer** (`mysql/net/`) — TCP connections, packet I/O, capability negotiation
2. **Protocol layer** (`mysql/protocol/`) — MySQL wire protocol: handshake, auth, command dispatch, result set formatting (text and binary)
3. **Query execution layer** (`mysql/query/`, top-level `mysql/`) — Routes parsed SQL to executor interfaces; uses go-sqlparser for parsing
4. **Application layer** — User code implementing the executor interfaces

**Execution flow**: Client connects → handshake (`mysql/protocol/`) → auth (`mysql/auth/`) → command received → query parsed (`mysql/query/parser.go`) → routed to executor interface → results formatted and returned.

### Executor Interfaces (`mysql/executor.go`)

Developers implement these to handle queries:

- **`DDOExecutor`** / **`DDOExExecutor`**: DDL — CREATE/ALTER/DROP DATABASE/TABLE/INDEX
- **`DMOExecutor`** / **`DMOExExecutor`**: DML — INSERT, SELECT, UPDATE, DELETE, USE, TRUNCATE
- **`TCOExecutor`**: TCL — BEGIN, COMMIT, ROLLBACK

The framework automatically handles system queries (SHOW VARIABLES, SET, etc.) via `executor_query.go` / `executor_query_ex.go`.

### Key Packages

| Package | Role |
|---|---|
| `mysql/` | Server interface, executor interfaces, main entry point |
| `mysql/protocol/` | Wire protocol (67 files): handshake, auth, query/stmt commands |
| `mysql/auth/` | Authentication manager, credential store, plugins |
| `mysql/stmt/` | Prepared statement management and ID registry |
| `mysql/query/` | SQL parsing, execution routing, result set, field type mappings |
| `mysql/net/` | Connection abstraction and pool |
| `mysql/encoding/binary/` | Binary serialization for MySQL's binary protocol |
| `mysql/errors/` | Error type definitions |
| `examples/go-mysqld/` | Reference in-memory MySQL server implementation |
| `mysqltest/` | Integration tests, fixtures (pcapng/hex captures), benchmarks |

### Reference Implementation

`examples/go-mysqld/` is a complete MySQL-compatible in-memory server. `examples/go-mysqld/server/store/` shows how to implement data storage. Study this to understand how to use the framework.

## Dependencies

Key external packages (see `go.mod`):
- `cybergarage/go-sqlparser` — SQL parsing (ANTLR4-based, SQL92 target)
- `cybergarage/go-authenticator` — Authentication plugin system
- `cybergarage/go-sasl` — SASL authentication (caching_sha2_password)
- `cybergarage/go-sqltest` — SQL compliance testing
- `cybergarage/go-tracing` — Distributed tracing support

## Linting

golangci-lint config is in `.golangci.yaml`. The `make lint` target covers `mysql/...`, `mysqltest/...`, and `examples/...`. Most linters are enabled; notable disabled ones include `cyclop`, `wrapcheck`, and `exhaustruct`.

## Commit Style

Use scoped, imperative prefixes: `fix:`, `style:`, `test:`, `feat:`, `refactor:`. Keep subjects short and focused (see recent git log).
