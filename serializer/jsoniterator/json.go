package jsoniterator

import (
	"bitbucket.org/modfin/yarf"
	j "github.com/json-iterator/go"
)

//Serializer is for encoding and decoding to json
func Serializer() yarf.Serializer {
	return yarf.Serializer{
		Marshal: func(v interface{}) ([]byte, error) {
			return j.Marshal(v)
		},
		Unmarshal: func(data []byte, v interface{}) error {
			return j.Unmarshal(data, v)
		},
	}
}
