package main

import (
	"bitbucket.org/modfin/yarf/transport/thttp"
	"fmt"
	"time"
	"bitbucket.org/modfin/yarf/example/simple"
	"os"
)


func main(){


	fmt.Println("Creating server transport")
	serverTransport, err := thttp.NewHttpTransporter(thttp.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServer(serverTransport)
	go serverTransport.Start()


	time.Sleep(200 * time.Millisecond)


	clientTransport, err := thttp.NewHttpTransporter(thttp.Options{Discovery: &thttp.DiscoveryDnsA{Host:"localhost"}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.RunClinet(clientTransport)


}


