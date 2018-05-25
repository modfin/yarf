package msgpack

import (
	"bitbucket.org/modfin/yarf"
	pack "github.com/vmihailenco/msgpack"
)

//Serializer is for encoding and decoding to msgpack
func Serializer() yarf.Serializer {
	return yarf.Serializer{
		ContentType: "application/msgpack",
		Marshal:     func(v interface{}) ([]byte, error) { return pack.Marshal(v) },
		Unmarshal:   func(data []byte, v interface{}) error { return pack.Unmarshal(data, v) },
	}
}
