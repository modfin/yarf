package yarf

import (
	j "github.com/json-iterator/go"
	"github.com/vmihailenco/msgpack/v5"
)

//Serializer is for encoding and decoding to json
func SerializerJson() Serializer {
	return Serializer{
		ContentType: "application/json",
		Marshal: func(v interface{}) ([]byte, error) {
			return j.Marshal(v)
		},
		Unmarshal: func(data []byte, v interface{}) error {
			return j.Unmarshal(data, v)
		},
	}
}

//Serializer is for encoding and decoding to msgpack
func SerializerMsgPack() Serializer {
	return Serializer{
		ContentType: "application/msgpack",
		Marshal:     func(v interface{}) ([]byte, error) { return msgpack.Marshal(v) },
		Unmarshal:   func(data []byte, v interface{}) error { return msgpack.Unmarshal(data, v) },
	}
}
