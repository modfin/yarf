[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/bitbucket.org/modfin/yarf)
[![Go Report Card](https://goreportcard.com/badge/bitbucket.org/modfin/yarf)](https://goreportcard.com/report/bitbucket.org/modfin/yarf)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/labstack/echo/master/LICENSE)


[comment]: ( https://codecov.io )
[comment]: ( Something ci )

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
* Support for Custom serializing, both of protocol and content level.
* Expressiv builder pattern fron client calls
* Client and server Middleware


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


func SHA256(req *yarf.Msg, resp *yarf.Msg) (err error) {

    hash := hashing.Sum256(req.Content)
    resp.SetParam("hash",
        base64.StdEncoding.EncodeToString(hash[:]))
    return
}

func main(){

	transport, err := thttp.NewHttpTransporter(thttp.Options{})
	if err != nil {
        log.Fatal(err)
    }
    server := yarf.NewServer(transport, "a", "namespace")

    server.HandleFunc(SHA256)
    server.Handle("add", func(req *yarf.Msg, resp *yarf.Msg) (err error){
        res := req.Param("val1").IntOr(0)+req.Param("val2").IntOr(0)
        resp.SetParam("res", res)
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
    
    msg, err := client.Request("a.namespace.add").
        SetParam("val1", 5).
        SetParam("val2", 7).
        Get()
    
    if err != nil{
        log.Fatal(err)
    }
    fmt.Println(" Result of 5 + 7 =", res.Param("res").IntOr(-1))


    msg, err = client.Request("a.namespace.SHA256").
        WithBinaryContent([]byte("Hello Yarf")).
        Get()

    if err != nil {
        return "", err
    }

    hash, ok := msg.Param("hash").String()
}

```

### Test
`go test -v ./...`
`./test.sh`, docker is requierd to run integration tests



## Design

Yarf consist in large of 5 things, API, Middleware, Protocol, Serlization
and Transport. We try to provide good defaults that both preform and are esay to use.
But, as most of us know, its not rare to end up in a corner case where things
have to be flexible. Therefore we try to seperate things in a clear way
that is easy to extend and change.


### API
Our main api that we provide for users of yarf is made up by a client, server
and the messages that pass between the two.


### Middleware
Yarf has support for middleware on both the client and server side of request
and both of them expect the same function type. Middleware can be use
when the need to decorate a message arise, do logging, caching, authentication or
anything else reqiering interception of messages


e.g. simple time logging
```go
 func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

		// Runs before client requst
        start = time.Now()

		err := next() // Running other middleware and request to server

		// Runs after client request
        fmt.Println("Request took", time.Now().Sub(start))

		return err
	}
```


e.g. a simple server side caching
```go
 func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

		// Runs before handler function

		key, ok := request.Param("cachekey").String()
		if !ok {
		    return errors.New("could not find cachekey in message")
		}

		b, found := getCachedItem(key)

		if found{
		    response.SetBinaryContent(b)
		    return nil
		}

		err := next() // Running other middleware and handler

		// Runs after handler function

        if err != nil{
            return err
        }

        setCachedItem("key", response.Content)

		return err
	}
```

On a client, middleware can be set on per request bases (local) or for all request going through
the client (global) and are run as follows

```
Client
   |
   V
   Global0 --> Global1 --> Local0 --> Local1 --> Call to transport layer
                                                        |
                                                        V
   Global0 <-- Global1 <-- Local0 <-- Local1 <-- Response from transport layer
   |
   V
Client

```

On the server, middleware can be set on per handler bases(local) or for all reguest going through
the server (global) and are run as follow

```
Incomming from transport layer
   |
   V
   Global0 --> Global1 --> Local0 --> Local1 --> Running handler function
                                                      |
                                                      V
   Global0 <-- Global1 <-- Local0 <-- Local1 <-- Response from handler function
   |
   V
Outgoing to transport layer
```


## Protocol
The yarf protocol is pretty straight forward but has a few layers to it.
It is not really that interesting unless you for some reason whant your
own serilization or need to change it

A message sent between client and server always start with a string prepended
by a newline, "\n", followed by bytes.

The first line describes the content type of the followig bytes, this
in order to deserlize the message into a yarf.Msg. This exist in order
to allow for different ways of (de)serlizing a message

The following bytes are then deserlized into the follwoing struct
```go
type Msg struct {
	Headers map[string]interface{}
	Content []byte
}
```




## TODO
* Unit testing
* More Examples
* Add support for reader and writers, for streaming requests/responses
* Http Transport
  * Improving service discover on HTTP transport
        * Consul
        * etcd
        * DNS
    * Improving loadbalancing on HTTP transport
    * Support for http2 and tls transport
* Middlewares
    * Proper Logging
    * Statistics and latency collection
    * Circuit breakers
    * Caching
    * Authentictaion, JWT