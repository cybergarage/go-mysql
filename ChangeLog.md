# Changelog

## v1.1.0 (2023-xx-xx)
- Updated executor interfaces to support major MySQL commands

## v1.0.3 (2023-10-03)
- Updated go-tracing package
- Updated connection to retrieve the UUID
- Added orderby and limit aliases

## v1.0.2 (2023-05-04)
- Updated Conn to embed tracer context
- Updated tracer spans

## v1.0.1 (2023-05-04)
- Updated Conn interface
- Remove debug log messages

## v1.0.0 (2023-05-04)
- Fixed executor interfaces for basic MySQL commands
- Updated logger functions to output more detail messages

## v0.9.4 (2023-04-23)
- Added tracing interface

## v0.9.3 (2023-04-02)
- Added connection logs
- Added Dockerfile

## v0.9.2 (2023-03-28)
- Updated mysqltest using go-sqltest
- Added sync.Map interface to mysql.Conn to store user data

## v0.9.1 (2023-02-23)
- Upgrade to go 1.20
- Upgrade to vitess v0.12.6
- Fixed compiler warnings

## v0.9.0 (2020-09-04)
- Initial public release  
