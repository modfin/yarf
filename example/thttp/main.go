package main

import (
	"bitbucket.org/modfin/yarf/example/simple"
	"bitbucket.org/modfin/yarf/transport/thttp"
	"fmt"
	"os"
	"time"
	"bitbucket.org/modfin/yarf/transport/tdecoraters"
)

func main() {

	fmt.Println("Creating server transport")
	serverTransport, err := thttp.NewHTTPTransporter(thttp.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServer(serverTransport, true)
	go serverTransport.Start()

	time.Sleep(200 * time.Millisecond)

	clientTransport, err := thttp.NewHTTPTransporter(thttp.Options{Discovery: &thttp.DiscoveryDNSA{Host: "127.0.0.1"}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.RunClient(tdecoraters.ClientLogging(clientTransport))

}
