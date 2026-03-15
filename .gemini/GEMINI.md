# Gemini Context: go-mysql

`go-mysql` is a database framework for easily implementing MySQL-compatible servers in Go. It handles the MySQL protocol and interprets major messages, allowing developers to focus on implementing Data Definition Language (DDL) and Data Manipulation Language (DML) query handlers.

## Project Overview

- **Core Functionality**: Provides a framework to build MySQL/MariaDB compatible servers.
- **Key Technologies**:
  - **Language**: Go (v1.25+)
  - **SQL Parser**: Based on [go-sqlparser](https://github.com/cybergarage/go-sqlparser), aiming for SQL92 compliance.
  - **Dependencies**: Uses various `cybergarage` projects like `go-logger`, `go-authenticator`, `go-tracing`, and `go-sqltest`.
- **Architecture**:
  - `mysql/`: Contains the core server logic, interfaces, and executor definitions.
  - `mysql/protocol/`: Implements the MySQL client/server protocol (packets, handshakes, result sets).
  - `mysql/auth/`: Handles authentication mechanisms (credentials, certificates).
  - `mysql/net/`: Connection management.
  - `examples/`: Contains sample implementations like `go-mysqld`.

## Building and Running

The project uses a `Makefile` for standard operations:

- **Build Examples**: `make build` (builds `go-mysqld`)
- **Run Example**: `make run` (installs and runs `go-mysqld`)
- **Test**: `make test` (runs all tests with coverage)
- **Lint**: `make lint` (runs `gofmt`, `go vet`, and `golangci-lint`)
- **Clean**: `make clean`
- **Docker**: `make image` (builds a Docker image for the example server)

## Development Conventions

- **Interface-Driven**: The framework relies heavily on interfaces. To implement a server, you typically implement `QueryExecutor` (which includes `TCOExecutor`, `DDOExecutor`, and `DMOExecutor`).
- **SQL Parsing**: Queries are parsed into ASTs using `github.com/cybergarage/go-sqlparser/sql`.
- **Code Quality**:
  - Always run `make format` (runs `gofmt`) before committing.
  - Use `golangci-lint` for static analysis.
  - New features or bug fixes should be accompanied by tests in `mysqltest/`.
- **Versioning**: The project uses a `version.gen` script in the `mysql/` directory to generate `version.go`. This is usually handled via `make version`.

## Key Interfaces

- **`mysql.Server`**: The main server interface.
- **`mysql.QueryExecutor`**: Composed of `TCOExecutor` (transactions), `DDOExecutor` (DDL), and `DMOExecutor` (DML).
- **`auth.Manager`**: Manages authentication and credential stores.
- **`mysql.Conn`**: Represents a client connection.
