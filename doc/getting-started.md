# Getting Started

This section describes how to implement your MySQL-compatible server using the go-mysql, and see  [Examples](doc/examples.md) about the sample implementation.

## STEP1: Inheriting Server

The go-mysql offers a core server, [mysql.Server](../mysql/server.go), and so inherit the core server in your instance as the following.

```
import (
	"github.com/cybergarage/go-mysql/mysql"
)

type MyServer struct {
	*mysql.Server
}

func NewMyServer() *MyServer {
	return &MyServer{
		Server: mysql.NewServer(),
	}
}
```

## STEP2: Preparing Query Handler

To handle queries to the your server, prepare a query handler according to [mysql.QueryExecutor](../mysql/executor.go) interface.

```
func NewMyServer() *MyServer {
	myserver := &MyServer{
		Server: mysql.NewServer(),
	}
    Myserver.SetQueryExecutor(myserver)
    return myserver
}

func (server *MyServer) Insert(context.Context, *Conn, *query.Insert) (*Result, error) {
    .....
}

....
```

The go-mysql offers the stub query executor, [mysql.BaseExecutor](../mysql/executor_base.go) which returns a success status for any query requests.
To inheriting the stub executor, you can start to implement only minimum query handle functions such as INSERT and SELECT.

## STEP3: Starting Server 

After implementing the query handler, start your server using  [mysql.Server::Start()](../mysql/server.go).

```
server := NewServer()

err := server.Start()
if err != nil {
	t.Error(err)
	return
}
defer server.Stop()

.... 
```
