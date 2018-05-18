package test

import (
	"bitbucket.org/modfin/yarf"
	"bitbucket.org/modfin/yarf/example/simple"
	"bitbucket.org/modfin/yarf/transport/thttp"
	"bitbucket.org/modfin/yarf/transport/tnats"
	"fmt"
	"golang.org/x/net/context"
	"os"
	"time"
)

// CreateHTTP returns a setup using HTTP as transport
func CreateHTTP(serializer yarf.Serializer, serverMiddleware ...yarf.Middleware) (client yarf.Client, stop func()) {

	serverTransport, err := thttp.NewHTTPTransporter(thttp.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServerWithSerializer(serverTransport, false, serializer, serverMiddleware...)

	go serverTransport.Start()
	time.Sleep(200 * time.Millisecond)

	clientTransport, err := thttp.NewHTTPTransporter(thttp.Options{Discovery: &thttp.DiscoveryDNSA{Host: "127.0.0.1"}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client = yarf.NewClient(clientTransport)
	client.WithSerializer(serializer)

	return client, func() {
		serverTransport.Stop(context.Background())
	}
}

// CreateNats returns a setup using Nats as transport
func CreateNats(serializer yarf.Serializer, serverMiddleware ...yarf.Middleware) (client yarf.Client, stop func()) {

	fmt.Println("Creating server transport")
	serverTransport, err := tnats.NewNatsTransporter("nats://demo.nats.io:4222", 10*time.Second)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServerWithSerializer(serverTransport, false, serializer, serverMiddleware...)
	time.Sleep(200 * time.Millisecond)

	clientTransport, err := tnats.NewNatsTransporter("nats://demo.nats.io:4222", 10000*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client = yarf.NewClient(clientTransport)

	return client, func() {}
}
