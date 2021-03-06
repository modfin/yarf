[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/modfin/yarf)
[![Go Report Card](https://goreportcard.com/badge/github.com/modfin/yarf)](https://goreportcard.com/report/github.com/modfin/yarf)
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
go get github.com/modfin/yarf
go get github.com/modfin/yarf/...
```

### Server
```go
package main
import (
    "github.com/modfin/yarf"
    "github.com/modfin/yarf/transport/thttp"
    "log"
)


func SHA256(req *yarf.Msg, resp *yarf.Msg) (err error) {

    hash := hashing.Sum256(req.Content)
    resp.SetParam("hash",
        base64.StdEncoding.EncodeToString(hash[:]))
    return
}

func join(req *yarf.Msg, resp *yarf.Msg) (err error) {

    slice, ok := req.Param("slice").StringSlice()
    if !ok{
        return errors.New("param arr was not a string slice")
    }

    joined := strings.Join(slice, "")

    resp.SetContent(joined)

    return nil
}

func main(){

	transport, err := thttp.NewHttpTransporter(thttp.Options{})
	if err != nil {
        log.Fatal(err)
    }
    server := yarf.NewServer(transport, "a", "namespace")

    server.HandleFunc(SHA256)
    server.HandleFunc(join)
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
    "github.com/modfin/yarf"
    "github.com/modfin/yarf/transport/thttp"
    "log"
    "fmt"
)
func main(){

    transport, err := thttp.NewHttpTransporter(thttp.Options{Discovery: &thttp.DiscoveryDefault{Host:"localhost"}})
    if err != nil {
        log.Fatal(err)
    }
    client := yarf.NewClient(transport)


    // Sending and reciving params
    msg, err := client.Request("a.namespace.add").
        WithParam("val1", 5).
        WithParam("val2", 7).
        Get()
    
    if err != nil{
        log.Fatal(err)
    }
    fmt.Println(" Result of 5 + 7 =", res.Param("res").IntOr(-1))


    //Sending binary content
    msg, err = client.Request("a.namespace.SHA256").
        WithBinaryContent([]byte("Hello Yarf")).
        Get()

    if err != nil {
        return "", err
    }

    hash, ok := msg.Param("hash").String()
    fmt.Println("ok", ok, "hash", hash)


    var joined string

    // Binding response content to string (works with structs, slices and so on)
    err = client.Request("a.namespace.join").
        WithParam("slice", []string{"jo", "in", "ed"}).
        BindResponseContent(&s).
        Done()

    fmt.Println("joined", joined, "err", err)

    // or
    err = client.Call("a.namespace.join", nil, &joined, yarf.NewParam("slice", []string{"joi", "ned"}))

    fmt.Println("joined", joined, "err", err)
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
		elapsed := time.Now().Sub(start)
        fmt.Println("Request to function", request.Function(), "took", elapsed)

		return err
	}
```


e.g. a simple server side caching
```go
 func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

		// Runs before handler function

		key, ok := request.Param("cachekey").String()
		if !ok {
		   // could not find cachekey in message, runs anyway
		   return next()
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

On a client, middleware can be set on per request bases (`local`) or for all request going through
the client (`global`) and are run as follows

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

On the server, middleware can be set on per handler bases (`local`) or for all reguest going through
the server (`global`) and are run as follow

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

The message struct header field contain mutiple keys that is usefull for the
yarf, but can also be used to pass paremeters. I shall howerver always contain
a content-type which helps yarf in dezerlization of Content, if needed.



## Serialization

This might be a some what confusing topic since there is a few layers to it.
But in general there are only two that has to be considerd. The serialization
of the protocol and the serialization of the content.

For both the the protocol and the content a content type shall aways be provided.
This helps the reciver of the message to deserialize the message and the content.

Serlizations can be done in different combinations and independet of each other.
The default is **msgpack** for both but can be changed per client, server or message basis.
Since it is might be hard to track what server has what and so on, serializers can
be regiserd in yarf by `yarf.RegisterSerializer(serializer)`. Yarf also provde some
extras ones, msgpack and json



## Transport
The transport layer works independetly of everything else and is responsibel
for service discovery, transport of data and provide to a context that can be
canceld from the client to the server.

A implmentation using Nats and one using HTTP is provided with yarf.

A transporter shall implment a rather simple api in order to work with yarf.
But since different transporters have different properties, some thinsgs may
vary. e.g. the function namespacing using Nats is a global and has no real need
for service discover, while HTTP has local namespace for each specific serivece.



## TODO
* Unit testing
* More documentation
* Add support for reader and writers, for streaming requests/responses
* Http Transport
    * Improving service discover on HTTP transport
        * Consul
        * etcd
        * DNS A
        * DNS SRV
    * Improving loadbalancing on HTTP transport
    * Support for http2 and tls transport
* Middlewares
    * Proper Logging
    * Statistics and latency collection
    * Circuit breakers
    * Caching
    * Authentictaion, JWT