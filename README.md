![Yarf](yarf.png "Yarf")

# yarf - Yet Another RPC Framework 


## Motivation 
There are a lot of rpc frameworks out there, so why one more. The simple answer is that may of them
 out there, such as gRPC, Twirp and so on did not fit our need. 
 They are often very opinionated, overly complex and 
 are in some cases much of a black box. 

## Overview
Yarf is a rpc framework focusing on ease of use and clear options of how to use.
 It provides an 


## Features
* Separation between protocol and transport
* Support for synchronise calls
* Support for callback functions
* Support for channel retrieval
 

## Supported transport layers
* http
* nats

## Quickstart
See examples for more examples 

### Intallation
```
go get bitbucket.org/modfin/yarf
go get bitbucket.org/modfin/yarf/transport/thttp
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
        Exec().
        Get()
    
    if err != nil{
        log.Fatal(err)
    }
    
    fmt.Println(" Result of 5 + 7 =", res.Param("res").IntOr(-1))
}

```


## TODO
* Testing
* More Examples
* Http Transport
  * Improving service discover on HTTP transport
  * Improving loadbalancing on HTTP transport
  * Support for http2 and tls transport
