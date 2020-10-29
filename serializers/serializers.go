package serializers

import (
	"github.com/modfin/yarf"
	"github.com/modfin/yarf/serializers/jsoniterator"
	"github.com/modfin/yarf/serializers/msgpack"
)

func init() {
	yarf.RegisterSerializer(jsoniterator.Serializer())
	yarf.RegisterSerializer(msgpack.Serializer())
}
