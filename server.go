package yarf

import (
	"context"
	"encoding/json"
	"reflect"
	"runtime"
	"strings"
)

// Server represents a yarf server with a particular transporter
type Server struct {
	transporter Transporter
	namespace   string
	middleware  []Middleware
	serializer  Serializer
}

// NewServer creates a new server with a particular server and name space of functions provided
func NewServer(t Transporter, namespace ...string) Server {
	s := Server{}
	s.transporter = t
	if len(namespace) > 0 {
		s.namespace = strings.Join(namespace, ".")
	} else {
		s.namespace = ""
	}
	s.serializer = Serializer{Marshal: json.Marshal, Unmarshal: json.Unmarshal}
	return s
}

func toServerErrorFrom(err RPCError, response *Msg) (responseData []byte) {
	response.SetStatus(err.Status)
	response.SetContent(err)
	responseData, _ = response.doMarshal()
	return responseData
}

func toServerError(status int, response *Msg, errors ...string) (responseData []byte) {
	return toServerErrorFrom(NewRPCError(status, strings.Join(errors, ";")), response)
}

// WithMiddleware add middleware to all requests
func (s *Server) WithMiddleware(middleware ...Middleware) {
	s.middleware = append(s.middleware, middleware...)
}

// WithSerializer setts the serializer used for transport.
func (s *Server) WithSerializer(serializer Serializer) {
	s.serializer = serializer
}

// HandleFunc creates a server endpoint for yarf using the handler function, the name of function will be on the format "namespace.FunctionName"
// e.g. my-namespace.Add, if a function named Add is passed into the function
func (s *Server) HandleFunc(handler func(request *Msg, response *Msg) error, middleware ...Middleware) {
	parts := strings.Split(runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(), ".")
	function := parts[len(parts)-1]
	s.Handle(function, handler, middleware...)
}

// Handle creates a server endpoint for yarf using the handler function, the name of function will be on the format "namespace.function"
// e.g. my-namespace.Add, if a string "Add" is passed into the function coupled with a handler function
func (s *Server) Handle(function string, handler func(request *Msg, response *Msg) error, middleware ...Middleware) {
	if s.namespace != "" {
		function = s.namespace + "." + function
	}

	s.transporter.Listen(function, func(ctx context.Context, requestData []byte) (responseData []byte) {

		req := Msg{serializer: s.serializer}
		resp := Msg{serializer: s.serializer}

		err := req.doUnmarshal(requestData)
		if err != nil {
			return toServerError(StatusUnmarshalError, &resp, err.Error())
		}
		req.ctx = ctx

		err = processMiddleware(&req, &resp, handler, append(s.middleware, middleware...)...)

		if err != nil {
			err2, ok := err.(RPCError)
			if ok {
				return toServerErrorFrom(err2, &resp)
			}

			return toServerError(StatusHandlerError, &resp, err.Error())
		}

		responseData, err = resp.doMarshal()

		if err != nil {
			return toServerError(StatusMarshalError, &resp, err.Error())
		}

		return responseData
	})
}
