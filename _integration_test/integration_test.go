package integration

import (
	"github.com/modfin/yarf"
	_ "github.com/modfin/yarf/serializers"
	"github.com/modfin/yarf/serializers/msgpack"
	"testing"

	"github.com/modfin/yarf/serializers/jsoniterator"
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
	{"MSG_PACK", msgpack.Serializer()},
	//{"JSON", json.Serializer()},
	{"JSON_ITERATOR", jsoniterator.Serializer()},
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
