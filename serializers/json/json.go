package json

import (
	"github.com/modfin/yarf"
	j "encoding/json"
)

//Serializer is for encoding and decoding to json
func Serializer() yarf.Serializer {
	return yarf.Serializer{
		ContentType: "application/json",
		Marshal: func(v interface{}) ([]byte, error) {
			return j.Marshal(v)
		},
		Unmarshal: func(data []byte, v interface{}) error {
			return j.Unmarshal(data, v)
		},
	}
}
