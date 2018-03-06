package yarf

import "context"

type Transporter interface {
	Call(ctx context.Context, function string, requestData []byte) (response []byte, err error)
	Listen(function string, toExec func(requestData []byte) (responseData []byte)) error
}