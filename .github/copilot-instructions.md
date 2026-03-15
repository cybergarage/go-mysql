# GitHub Copilot Instructions for go-mysql

## Project Overview

`go-mysql` (`github.com/cybergarage/go-mysql`) is a Go framework for building MySQL-compatible servers. It handles the MySQL wire protocol, connection lifecycle, authentication, and SQL parsing — so implementors only need to provide DDL and DML handlers.

## Build, Test, and Lint

```bash
# Full pipeline (format → vet → lint → test)
make

# Run tests only (enforces linting first)
make test

# Format code (regenerates mysql/version.go + gofmt -s)
make format

# Lint only (go vet + golangci-lint via .golangci.yaml)
make lint

# Run a single test by name
go test -v -run ^TestName ./mysql/...
go test -v -run ^TestName ./mysqltest/server/...

# Quicker iteration without the full pipeline
go test ./mysql/... ./mysqltest/...

# Run sysbench benchmark tests
make sysbench
# equivalent to:
go test -v -p 1 -run ^TestSysbench github.com/cybergarage/go-mysql/mysqltest/sysbench

# Build example binary (go-mysqld)
make build
```

Tests run sequentially (`-p 1`) with a 10-minute timeout and generate coverage reports (`mysql-cover.out`, `mysql-cover.html`). CI runs `make test` on Go 1.25.

The `mysqltest/sqltest/` package contains `TestSQLTest` — a debug harness where you can uncomment specific test names from the `testNames` slice to run individual SQL test cases.

## Architecture

```
mysql/           # Core library (import as github.com/cybergarage/go-mysql/mysql)
  protocol/      # MySQL wire protocol: packets, handshake, OK/ERR/EOF, resultsets, capabilities
  net/           # Connection interface (wraps go-sqlparser net.Conn + StatementManager)
  query/         # SQL executor interfaces, ResultSet/ResultSetRow/ResultSetColumn types
  encoding/      # Binary encoding utilities
  errors/        # Sentinel errors and constructors
  stmt/          # Prepared statement management
  auth/          # Authentication (delegates to go-authenticator)

mysqltest/       # Test infrastructure (not part of the library API)
  server/        # In-memory test server used by all test suites
  sqltest/       # SQL correctness test runner (uses go-sqltest)
  protocol/      # Protocol-level tests
  benchbase/     # BenchBase workload tests
  sysbench/      # Sysbench workload tests
  ycsb/          # YCSB workload tests

examples/go-mysqld/  # Reference server implementation (in-memory store)
  server/            # Server struct composing mysql.Server + store.Store
```

## Implementing a MySQL-Compatible Server

There are two levels of integration:

### Option 1 — SQLExecutor (simplest)
Implement `sql.Executor` from `go-sqlparser` and register it. The framework dispatches all SQL operations through it.

```go
server := mysql.NewServer()
server.SetSQLExecutor(myStore)  // myStore implements sql.Executor
server.Start()
```

### Option 2 — QueryExecutor (fine-grained)
Implement the composed executor interfaces for full control over individual SQL operations:

- `TCOExecutor` — `Begin`, `Commit`, `Rollback`
- `DDOExecutor` — `CreateDatabase`, `CreateTable`, `AlterDatabase`, `AlterTable`, `DropDatabase`, `DropTable`
- `DMOExecutor` — `Use`, `Insert`, `Select`, `Update`, `Delete`
- `ExQueryExecutor` — `CreateIndex`, `DropIndex`, `Truncate`

```go
server.SetQueryExecutor(myExecutor)
```

All executor methods return `(protocol.Response, error)`.

## Coding Style

- **Indentation**: tabs (enforced by `gofmt`)
- **Import formatting**: tool-driven via `goimports` and `gci` — do not hand-format imports
- **Linting**: `.golangci.yaml` configures `golangci-lint`; run `make lint` before committing
- **Go version**: 1.25

## Commit Conventions

Short, scoped, imperative subjects. Use prefixes when useful: `style:`, `fix:`, `test:`. Example: `style: fix staticcheck issues`. PRs should include a problem statement, implementation summary, and the validation commands run.

## Key Conventions

### Returning responses from executors

Use `protocol.NewResponseWithError(err)` for DDL/DML operations that don't return rows — it returns `OK` if `err == nil`, or `ERR` otherwise:

```go
func (s *Store) CreateTable(conn mysql.Conn, stmt sql.CreateTable) (mysql.Response, error) {
    err := s.doCreate(stmt)
    return protocol.NewResponseWithError(err)
}
```

For `SELECT`, return a `protocol.TextResultSet` built from a `query.ResultSet`.

### Functional options pattern for protocol packets

All protocol packet constructors use functional options:

```go
protocol.NewOK(protocol.WithOKAffectedRows(1))
protocol.NewERR(protocol.WithERRMessage("table not found"))
protocol.NewOK(protocol.WithOKSecuenceID(seq), protocol.WithOKCapability(cap))
```

### Error types

Use sentinel errors from `mysql/errors/` and wrap them with constructors:

```go
errors.ErrNotImplemented
errors.ErrNotSupported
errors.ErrNotExist
errors.ErrNotFound
errors.ErrInvalid

errors.NewErrNotImplemented("ALTER TABLE")
errors.NewErrUnsupported("subqueries")
errors.NewErrNotExis("database")
```

### Conn interface

`mysql.Conn` (= `net.Conn`) extends `go-sqlparser`'s `sql/net.Conn` with `stmt.StatementManager` for prepared statement tracking. Pass it through to your store/executor as-is.

### Test server pattern

All `mysqltest/` test suites use `mysqltest/server.Server`, which composes `mysql.Server` with an in-memory store and registers itself as the `SQLExecutor`. New test suites should follow this pattern.

## Dependencies

Key external dependencies and their roles:
- `go-sqlparser` — SQL parsing; provides `sql.Executor`, `sql.CreateTable`, `sql.Select`, etc.
- `go-sqltest` — SQL test suite runner used in `mysqltest/sqltest/`
- `go-authenticator` / `go-sasl` — Authentication and SASL mechanism support
- `go-tracing` — Distributed tracing integration via `SetTracer()`
- `go-logger` — Structured logging (use `log.EnableStdoutDebug(true)` in tests)
