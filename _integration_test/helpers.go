package integration

import (
	"github.com/modfin/yarf"
	"github.com/modfin/yarf/example/simple"
	"github.com/modfin/yarf/transport/thttp"
	"github.com/modfin/yarf/transport/tnats"
	"fmt"
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

	clientTransport, err := thttp.NewHTTPTransporter(thttp.Options{Discovery: &thttp.DiscoveryDNSA{Host: "localhost"}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client = yarf.NewClient(clientTransport)
	client.WithProtocolSerializer(serializer)
	client.WithSerializer(serializer)

	return client, func() { serverTransport.Close() }
}

// CreateNats returns a setup using Nats as transport
func CreateNats(serializer yarf.Serializer, serverMiddleware ...yarf.Middleware) (client yarf.Client, stop func()) {

	fmt.Println("Creating server transport")
	serverTransport, err := tnats.NewNatsTransporter("nats://localhost:4222", 10*time.Second)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServerWithSerializer(serverTransport, false, serializer, serverMiddleware...)
	time.Sleep(200 * time.Millisecond)

	clientTransport, err := tnats.NewNatsTransporter("nats://localhost:4222", 10000*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client = yarf.NewClient(clientTransport)
	client.WithProtocolSerializer(serializer)
	client.WithSerializer(serializer)

	return client, func() { serverTransport.Close() }
}
