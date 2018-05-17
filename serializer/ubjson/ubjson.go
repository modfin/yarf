package ubjson

import (
	"bitbucket.org/modfin/yarf"
	"github.com/jmank88/ubjson"
)

//Serializer is for encoding and decoding to ubjson
func Serializer() yarf.Serializer {
	return yarf.Serializer{
		Marshal:   ubjson.Marshal,
		Unmarshal: ubjson.Unmarshal,
	}
}
