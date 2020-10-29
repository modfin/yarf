package main

import (
	"github.com/modfin/yarf/example/simple"
	"github.com/modfin/yarf/transport/tnats"
	"fmt"
	"log"
	"time"
)

func main() {

	// Nats own demo servers is not well suited for large payloads, which is tested. So you might want to change this
	// if you want to accrual integration nats as the transport layer
	fmt.Println("Creating server transport")
	serverTransport, err := tnats.NewNatsTransporter("nats://demo.nats.io:4222", 10*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	simple.StartServer(serverTransport, true)

	time.Sleep(200 * time.Millisecond)

	clientTransport, err := tnats.NewNatsTransporter("nats://demo.nats.io:4222", 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	simple.RunClient(clientTransport)

}
