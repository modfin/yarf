package main

import (
	"bitbucket.org/modfin/yarf/example/simple"
	"bitbucket.org/modfin/yarf/transport/thttp"
	"fmt"
	"os"
	"time"
)

func main() {

	fmt.Println("Creating server transport")
	serverTransport, err := thttp.NewHTTPTransporter(thttp.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServer(serverTransport)
	go serverTransport.Start()

	time.Sleep(200 * time.Millisecond)

	clientTransport, err := thttp.NewHTTPTransporter(thttp.Options{Discovery: &thttp.DiscoveryDNSA{Host: "localhost"}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.RunClient(clientTransport)

}
