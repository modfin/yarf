package tnats

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

	fmt.Println("Creating server transport")
	serverTransport, err := NewNatsTransporter("nats://demo.nats.io:4222", 10*time.Second)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServer(serverTransport, false)

	time.Sleep(200 * time.Millisecond)

	clientTransport, err = NewNatsTransporter("nats://demo.nats.io:4222", 10*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client = yarf.NewClient(clientTransport)


	os.Exit(m.Run())
}



func TestIntegrationTransportNats(t *testing.T){
	t.Run("Integration TransportNats", integration.GetIntegrationTest(client))
}

func BenchmarkTransportHttp(b *testing.B){

	b.Run("Integration BenchmarkNats", integration.GetBenchmarkAdd(client,  10, 11))
}