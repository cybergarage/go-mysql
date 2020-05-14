# Coding Guidelines

To develop `go-mysql`, you must follow the following coding guidelines basically.

- [Effective Go](https://golang.org/doc/effective_go.html#interface-names)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

## Static Analyzer

In addition to the above standard guideline, you must use the following static analyzers and fix all warnings.

- [go vet](https://golang.org/cmd/vet/)
- [GolangCI-Lint](https://github.com/golangci/golangci-lint)

For the setting of [GolangCI-Lint](https://github.com/golangci/golangci-lint), see the following setting file to know in more detail.

```
.golangci.yml
```