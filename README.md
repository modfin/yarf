[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/bitbucket.org/modfin/yarf)
[![Go Report Card](https://goreportcard.com/badge/bitbucket.org/modfin/yarf)](https://goreportcard.com/report/bitbucket.org/modfin/yarf)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE)


[comment]: ( https://codecov.io )
[comment]: ( Something ci )
[comment]: ( ![Yarf](yarf.png "Yarf") )

# yarf - Yet Another RPC Framework 
For simple comunication between services

## Motivation 
There are a few rpc frameworks out there, so why one more. The simple
answer is that may of them out there, such as gRPC, Twirp and so on did
not fit our need. They are often very opinionated, overly complex and
are in some cases much of a black box trying to solve every problem in
many languages.

What we found was that we were writing models i protobuf to be used in
frameworks such as gRPC and then on both client and server side had to
map them into local structs since protobuf was not expresive enough.
This instead of just having a versioned struct in a common kit repo
or such. In essens fighting with the rpc libs we tired in order for it
to work with our use cases.

## Overview
Yarf is a rpc framework focusing on ease of use and clear options of
how to use.


## Features
* Separation between protocol and transport
* Support for synchronise calls, channals and callbacks
* Support for middleware on both client and server side
* Support for context passing between client and server
* Expressiv builder pattern fron client calls


## Supported transport layers
* http
* nats

## Quickstart
See examples for more examples/simple

### Intallation
```
go get bitbucket.org/modfin/yarf
go get bitbucket.org/modfin/yarf/...
```

### Server
```go
package main
import (
    "bitbucket.org/modfin/yarf"
    "bitbucket.org/modfin/yarf/transport/thttp"
    "log"
)
func main(){

	transport, err := thttp.NewHttpTransporter(thttp.Options{})
	if err != nil {
        log.Fatal(err)
    }
    server := yarf.NewServer(transport, "a", "namespace")
    
    server.Handle("add", func(req *yarf.Msg, resp *yarf.Msg) (err error){
    
        resp.SetParam("res", req.Param("val1").IntOr(0)+req.Param("val2").IntOr(0))
        return nil
    })
    
    log.Fatal(transport.Start())
}
```

### Client
```go
package main
import (
    "bitbucket.org/modfin/yarf"
    "bitbucket.org/modfin/yarf/transport/thttp"
    "log"
    "fmt"
)
func main(){

    transport, err := thttp.NewHttpTransporter(thttp.Options{Discovery: &thttp.DiscoveryDefault{Host:"localhost"}})
    if err != nil {
        log.Fatal(err)
    }
    client := yarf.NewClient(transport)
    
    res, err := client.Request("a.namespace.add").
        SetParam("val1", 5).
        SetParam("val2", 7).
        Get()
    
    if err != nil{
        log.Fatal(err)
    }
    
    fmt.Println(" Result of 5 + 7 =", res.Param("res").IntOr(-1))
}

```

### Test
`go test -v ./...`

(NATS test might not work for bigger payload since a public server is
used for integration testing)


## TODO
* Testing
* More Examples
* Http Transport
  * Improving service discover on HTTP transport
  * Improving loadbalancing on HTTP transport
  * Support for http2 and tls transport
* Middlewares
  * Proper Logging
  * Statistics and latency collection
  * Circuit breakers
  * Caching
  * Authentictaion, JWT