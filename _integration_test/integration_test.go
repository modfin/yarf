package integration

import (
	"context"
	"fmt"
	"github.com/modfin/yarf"
	"github.com/modfin/yarf/example/simple"
	"github.com/modfin/yarf/transport/tnats"
	"os"
	"testing"
	"time"
)

var transportTable = []struct {
	name  string
	start func(serializer yarf.Serializer, serverMiddleware ...yarf.Middleware) (client yarf.Client, stop func())
	extra bool
}{
	{"HTTP", CreateHTTP, true},
	{"NATS", CreateNats, false},
}

var serializerTable = []struct {
	name       string
	serializer yarf.Serializer
}{
	{"MSG_PACK", yarf.SerializerMsgPack()},
	{"JSON", yarf.SerializerJson()},
}

func TestServerClose(t *testing.T) {
	fmt.Println("Creating server transport")
	serverTransport, err := tnats.NewNatsTransporter("nats://localhost:4222", 10*time.Second)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	server := simple.StartServerWithSerializer(serverTransport, false, yarf.SerializerJson())
	defer server.Close()
	time.Sleep(200 * time.Millisecond)

	clientTransport, err := tnats.NewNatsTransporter("nats://localhost:4222", 10000*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := yarf.NewClient(clientTransport)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1500)
	defer cancel()
	transit := client.Request("a.integration.sleep").
		WithParam("sleep", 1000).
		WithContext(ctx).
		Async()
	time.Sleep(time.Millisecond*500)
	err = server.Close()
	if err != nil{
		t.Fatal("error closing server", err)
		return
	}
	_, err = transit.Get()

	if err == nil{
		t.Fatal("expected context timeout error")
		return
	}
	if err.Error() !="context deadline exceeded"{
		t.Fatal("expected error: msg, 'context deadline exceeded', got", err)
	}
}

func TestServerCloseGraceful(t *testing.T) {
	fmt.Println("Creating server transport")
	serverTransport, err := tnats.NewNatsTransporter("nats://localhost:4222", 10*time.Second)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	server := simple.StartServerWithSerializer(serverTransport, false, yarf.SerializerJson())
	defer server.Close()
	time.Sleep(200 * time.Millisecond)

	clientTransport, err := tnats.NewNatsTransporter("nats://localhost:4222", 10000*time.Second)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := yarf.NewClient(clientTransport)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1500)
	defer cancel()
	transit := client.Request("a.integration.sleep").
		WithParam("sleep", 1000).
		WithContext(ctx).
		Async()
	time.Sleep(time.Millisecond*500)
	go func() {
		err = server.CloseGraceful(time.Second*2)
		if err != nil{
			t.Fatal("error closing server", err)
			return
		}
	}()
	msg, err := transit.Get()

	if err != nil{
		t.Fatal("expected no error, but got", err)
		return
	}
	if msg.Param("res").IntOr(0) != 1000 {
		t.Fatal("expected res to be 1000, got", msg.Param("res").Value() )
	}

}

func TestMatrix(t *testing.T) {

	for _, ser := range serializerTable {
		for _, tran := range transportTable {

			client, stop := tran.start(ser.serializer)
			t.Run(tran.name+"/"+ser.name, GetIntegrationTest(client))

			t.Run(tran.name+"/"+ser.name+"/LARGE_PAYLOAD", GetExtraIntegrationTest(client))

			stop()
		}
	}
}

func TestMissMatchSerializers(t *testing.T) {
	tran := transportTable[1]

	client, stop := tran.start(yarf.SerializerJson())
	client.WithProtocolSerializer(yarf.SerializerMsgPack())
	client.WithSerializer(yarf.SerializerMsgPack())
	t.Run(tran.name+"/JSON/MSG_PACK", GetIntegrationTest(client))
	stop()

	tran = transportTable[1]

	client, stop = tran.start(yarf.SerializerMsgPack())
	client.WithProtocolSerializer(yarf.SerializerJson())
	client.WithSerializer(yarf.SerializerJson())
	t.Run(tran.name+"/MSG_PACK/JSON", GetIntegrationTest(client))
	stop()
}
