package yarf

import (
	"context"
)

func init() {
	serializers = map[string]Serializer{}
	RegisterSerializer(SerializerJson())
	RegisterSerializer(MsgPackSerializer())
}

// Transporter is the interface that must be fulfilled for a transporter.
type Transporter interface {
	CallTransporter
	ListenTransporter
}

// CallTransporter is the interface that must be fulfilled for a transporter to be used as a client.
type CallTransporter interface {
	Call(ctx context.Context, function string, requestData []byte) (response []byte, err error)
}

// ListenTransporter is the interface that must be fulfilled for a transporter to be used as a server
type ListenTransporter interface {
	Listen(function string, toExec func(ctx context.Context, requestData []byte) (responseData []byte)) error
	Close() error
}

// Serializer is the interface that must be fulfilled for protocolSerializer of data before transport.
type Serializer struct {
	ContentType string
	Marshal     func(v interface{}) ([]byte, error)
	Unmarshal   func(data []byte, v interface{}) error
}

var serializers map[string]Serializer



// RegisterSerializer lets a user register a protocolSerializer for a specific content type
// this allow yarf to bind message content to that specific serial format.
// Yarf standard serializers can be registered by importing with side effect
func RegisterSerializer(serializer Serializer) {
	serializers[serializer.ContentType] = serializer
}

func serializer(contentType string) (serializer Serializer, ok bool) {
	serializer, ok = serializers[contentType]
	return
}

func defaultSerializer() Serializer {
	return MsgPackSerializer()
}
