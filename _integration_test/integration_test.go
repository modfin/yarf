package integration

import (
	"github.com/modfin/yarf"
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

var serializerTable = []struct {
	name       string
	serializer yarf.Serializer
}{
	{"MSG_PACK", yarf.MsgPackSerializer()},
	{"JSON_ITERATOR", yarf.SerializerJson()},
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
