# Changelog

## v1.2.0 (2025-01-xx)
- New Features:
  - Support for major authentication methods.
    - mysql_native_password
    - caching_sha2_password

- New Features:
## v1.1.4 (2024-12-31)
- Updates:
  - Authentication interfaces

## v1.1.3 (2024-12-28)
- New Features:
  - Supported for TLS connections.
  - Supported certificate authentication for TLS connection

## v1.1.2 (2024-12-11)
- New Features:
  - Supported converting `CREATE INDEX` and `DROP INDEX` commands to `ALTER TABLE`.
- Added:
  - Enabled secondary index tests of `go-sqltest` package.

## v1.1.1 (2024-11-24)
- Updates:
  - Adapted to the latest SQL executor interface provided by `go-sqlparser`.
  - Updated example to share a common SQL executor with `go-postgresql`

## v1.1.0 (2024-11-16)
- Improvements:
  - Re-implemented the protocol layer and parser independently.
  - Removed dependency on Vitess.

## v1.0.6 (2024-05-14)
- Bug Fixes:
  - Resolved security issues flagged by Dependabot alerts.

## v1.0.5 (2023-10-27)
- Enhancements:
  - Connection:
    - Added deadline methods to `Conn`.

## v1.0.4 (2023-10-03)
- Updates:
  - Enhanced connection to retrieve the UUID.
  - Updated the `go-tracing` package.

## v1.0.3 (2023-05-13)
- New Features:
  - Added support for `ORDER BY` and `LIMIT` aliases.

## v1.0.2 (2023-05-04)
- Improvements:
  - Embedded tracer context into `Conn`.
  - Enhanced tracer spans.

## v1.0.1 (2023-05-04)
- Updates:
  - Refined `Conn` interface.
  - Removed debug log messages.

## v1.0.0 (2023-05-04)
- Major Release:
  - Fixed executor interfaces for basic MySQL commands.
  - Updated logger functions to provide more detailed messages.

## v0.9.4 (2023-04-23)
- New Features:
  - Introduced a tracing interface.

## v0.9.3 (2023-04-02)
- Enhancements:
  - Added connection logs.
  - Introduced a `Dockerfile` for containerization.

## v0.9.2 (2023-03-28)
- Updates:
  - Enhanced `mysqltest` using `go-sqltest`.
  - Added `sync.Map` interface to `mysql.Conn` for user data storage.

## v0.9.1 (2023-02-23)
- Updates:
  - Upgraded to Go 1.20.
  - Upgraded to Vitess v0.12.6.
  - Fixed compiler warnings.

## v0.9.0 (2020-09-04)
- Initial Release:
  - First public version.
