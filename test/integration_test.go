package test

import (
	"bitbucket.org/modfin/yarf"
	"bitbucket.org/modfin/yarf/serializer/json"
	"bitbucket.org/modfin/yarf/serializer/jsoniterator"
	"bitbucket.org/modfin/yarf/serializer/msgpack"
	"testing"
)

var transportTable = []struct {
	name  string
	start func(serializer yarf.Serializer, serverMiddleware ...yarf.Middleware) (client yarf.Client, stop func())
	extra bool
}{
	{"HTTP", CreateHTTP, true},
	{"NATS", CreateNats, false},
}

func TestTransport(t *testing.T) {

	for _, tran := range transportTable {
		client, stop := tran.start(msgpack.Serializer())
		t.Run(tran.name, GetIntegrationTest(client))

		if tran.extra {
			t.Run(tran.name+"/LARGE_PAYLOAD", GetExtraIntegrationTest(client))
		}

		stop()
	}

}

var serializerTable = []struct {
	name       string
	serializer yarf.Serializer
}{
	{"MSG_PACK", msgpack.Serializer()},
	{"JSON", json.Serializer()},
	{"JSON_ITERATOR", jsoniterator.Serializer()},
}

func TestSerializer(t *testing.T) {

	for _, ser := range serializerTable {
		client, stop := CreateHTTP(ser.serializer)
		t.Run(ser.name, GetIntegrationTest(client))
		t.Run(ser.name+"/LARGE_PAYLOAD", GetExtraIntegrationTest(client))
		stop()
	}

}
