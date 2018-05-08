package thttp

import (
	"testing"
	"fmt"
	"bitbucket.org/modfin/yarf"
	"os"
	"bitbucket.org/modfin/yarf/example/simple"
	"time"
	"bitbucket.org/modfin/yarf/example/simple/integration"
)

var clientTransport yarf.Transporter
var client yarf.Client


func TestMain(m *testing.M){

	serverTransport, err := NewHTTPTransporter(Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServer(serverTransport, false)
	go serverTransport.Start()

	time.Sleep(200 * time.Millisecond)

	clientTransport, err = NewHTTPTransporter(Options{Discovery: &DiscoveryDNSA{Host: "127.0.0.1"}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client = yarf.NewClient(clientTransport)


	os.Exit(m.Run())
}



func TestIntegrationTransportHttp(t *testing.T){
	t.Run("Integration TransportHttp", integration.GetIntegrationTest(client))
}


func BenchmarkTransportHttp(b *testing.B){

	b.Run("Integration BenchmarkHttp", integration.GetBenchmarkAdd(client,  10, 11))
}

