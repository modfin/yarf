package main

import (
	"github.com/modfin/yarf"
	"github.com/modfin/yarf/example/simple"
	"github.com/modfin/yarf/middleware"
	"github.com/modfin/yarf/transport/thttp"
	"fmt"
	"github.com/opentracing/basictracer-go/examples/dapperish"
	"github.com/opentracing/opentracing-go"
	"os"
	"time"
)

func main() {

	opentracing.InitGlobalTracer(dapperish.NewTracer("dapperish_tester"))

	fmt.Println("Creating server transport")
	serverTransport, err := thttp.NewHTTPTransporter(thttp.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServerWithSerializer(serverTransport, true, yarf.MsgPackSerializer(), middleware.OpenTracing("Server> "))
	go serverTransport.Start()

	time.Sleep(200 * time.Millisecond)

	clientTransport, err := thttp.NewHTTPTransporter(thttp.Options{Discovery: &thttp.DiscoveryDNSA{Host: "127.0.0.1"}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.RunClient(clientTransport)

}
