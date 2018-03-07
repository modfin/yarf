package main

import (
	"fmt"
	"time"
	"bitbucket.org/modfin/yarf/example/simple"
	"bitbucket.org/modfin/yarf/transport/tnats"
	"log"
)


func main(){

	// Nats own demo servers is not well suited for large payloads, which is tested. So you might want to change this
	// if you want to accrual test nats as the transport layer
	fmt.Println("Creating server transport")
	serverTransport, err := tnats.NewNatsTransporter("nats://demo.nats.io:4222", 10*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	simple.StartServer(serverTransport)


	time.Sleep(200 * time.Millisecond)

	clientTransport, err := tnats.NewNatsTransporter("nats://demo.nats.io:4222", 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	simple.RunClinet(clientTransport)


}


