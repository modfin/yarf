package serializers

import (
	"bitbucket.org/modfin/yarf"
	"bitbucket.org/modfin/yarf/serializers/jsoniterator"
	"bitbucket.org/modfin/yarf/serializers/msgpack"
)

func init() {
	yarf.RegisterSerializer(jsoniterator.Serializer())
	yarf.RegisterSerializer(msgpack.Serializer())
}
