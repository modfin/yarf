package msgpack

import (
	"bitbucket.org/modfin/yarf"
	"bitbucket.org/modfin/yarf/example/simple"
	"bitbucket.org/modfin/yarf/example/simple/integration"
	"bitbucket.org/modfin/yarf/transport/thttp"
	"fmt"
	"os"
	"testing"
	"time"
)

var clientTransport yarf.Transporter
var client yarf.Client

func TestMain(m *testing.M) {

	serverTransport, err := thttp.NewHTTPTransporter(thttp.Options{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	simple.StartServerWithSerializer(serverTransport, false, Serializer())
	go serverTransport.Start()

	time.Sleep(200 * time.Millisecond)

	clientTransport, err = thttp.NewHTTPTransporter(thttp.Options{Discovery: &thttp.DiscoveryDNSA{Host: "127.0.0.1"}})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client = yarf.NewClient(clientTransport)
	client.WithSerializer(Serializer())

	os.Exit(m.Run())
}

func TestIntegrationSerializerMsgPack(t *testing.T) {
	t.Run("Integration SerializerMsgPack", integration.GetIntegrationTest(client))
}
