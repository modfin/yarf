package yarf

import "context"

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
}

// Serializer is the interface that must be fulfilled for serializer of data before transport.
type Serializer struct {
	Marshal   func(v interface{}) ([]byte, error)
	Unmarshal func(data []byte, v interface{}) error
}
