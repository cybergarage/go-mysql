# Examples

## go-mysqld

`go-mysqld` is a simple MySQL-compatible implementation using `go-mysql`. The sample implementation is a in-memory MySQL-compatible server, and it supports only a table and do not support any JOIN queries.
```
 go-mysqld is an example of implementing a compatible MySQL server using go-mysql.
	NAME
	 go-mysqld

	SYNOPSIS
	 go-mysqld [OPTIONS]

	OPTIONS
	-v      : Enable verbose output.
	-p      : Enable profiling.

	RETURN VALUE
	  Return EXIT_SUCCESS or EXIT_FAILURE
```

To install the binary, use the following command. The install command installs the utility programs into `GO_PATH/bin`.

```
make install
```

The profile option enables pprof serves of Go which has the HTTP interface to observe `go-mysql` profile data.

- [The Go Programming Language - Package pprof](https://golang.org/pkg/net/http/pprof/)
- [The Go Blog - Profiling Go Programs](https://blog.golang.org/profiling-go-programs)
