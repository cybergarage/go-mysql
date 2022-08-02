# Building and Testing

## Dependency Packages

### Vitess

`go-mysql` extends [Vitess](https://github.com/vitessio/vitess) package which is a sharding package for MySQL compatible servers. So `go-mysql` requires the [Vitess](https://github.com/vitessio/vitess) package for the compiling, but `go-mysql` has not succeeded to import the package directly yet due to the following compile errors. 

```
go build github.com/cybergarage/go-mysql/mysql
# pkg-config --cflags  -- gomysql
Package gomysql was not found in the pkg-config search path.
Perhaps you should add the directory containing `gomysql.pc'
to the PKG_CONFIG_PATH environment variable
No package 'gomysql' found
pkg-config: exit status 1
```

Therefore, `go-mysql` imports the package using the using replace directive of Go module, and so you have to check out the package into the same directory of `go-mysql` as the following.

```
git checkout https://github.com/vitessio/vitess
```


## Testing

### Unit Testing

For the unit test, `go-mysql` conforms to the following official testing of Go language.

- [The Go Programming Language - Package testing](https://golang.org/pkg/testing/)

To run all tests, use the following command.

```
make test
```
