# Repository Guidelines

## Project Structure & Module Organization

Core library code lives in `mysql/`, split by concern such as `auth/`, `protocol/`, `query/`, `stmt/`, and `net/`. Tests are primarily under `mysqltest/`, with focused packages like `protocol/`, `sqltest/`, `server/`, `sysbench/`, and `ycsb/`. Example server code is in `examples/go-mysqld/`. Reference docs and diagrams live in `doc/`. Benchmark helpers and bundled assets are under `scripts/` and `mysqltest/benchbase/`; treat large `.pcapng`, `.jar`, and generated coverage files as support artifacts, not code.

## Build, Test, and Development Commands

Use the `Makefile` as the source of truth:

- `make format`: regenerates `mysql/version.go` and runs `gofmt -s` on library, tests, and examples.
- `make lint`: runs `go vet` and `golangci-lint` across `mysql/...`, `mysqltest/...`, and `examples/...`.
- `make test`: enforces linting, runs the full Go test suite with coverage, and writes `mysql-cover.out` and `mysql-cover.html`.
- `make build`: builds the example daemon binary after tests pass.
- `make run`: installs and starts `go-mysqld`.
- `make sysbench`: runs only the sysbench integration test.

CI runs `make test` on Go 1.25, so local changes should pass that path before review.

## Coding Style & Naming Conventions

Write Go compatible with `go 1.25`. Use tabs for indentation and keep formatting tool-driven; do not hand-format imports. Preferred tooling is `gofmt`, `goimports`, `gci`, and `golangci-lint` via `.golangci.yaml`. Follow Go naming: exported identifiers in `CamelCase`, unexported helpers in `camelCase`, and tests in `*_test.go`. Keep package names short and lowercase, matching the directory name.

## Testing Guidelines

Use the standard Go `testing` package. Name tests `TestXxx` and keep test fixtures near the package that owns them, such as `mysqltest/protocol/data/`. Run `make test` for full validation and `go test ./mysql/... ./mysqltest/...` for quicker iteration. Coverage output is expected from `make test`; avoid merging changes that reduce exercised paths without justification.

## Commit & Pull Request Guidelines

Recent history favors short, scoped subjects such as `style: fix staticcheck issues` and `Update .golangci.yaml`. Keep commit messages imperative and prefixed when useful (`style:`, `fix:`, `test:`). Pull requests should include a concise problem statement, the implementation summary, and the exact validation commands you ran. Link related issues when applicable; screenshots are only relevant for docs or example UI changes.
