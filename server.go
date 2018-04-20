package yarf

import (
	"reflect"
	"runtime"
	"strings"
)

// Server represents a yarf server with a particular transporter
type Server struct {
	transporter Transporter
	namespace   string
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
	return s
}

func toServerError(status int, response *Msg, errors ...string) (responseData []byte) {

	response.SetStatus(status)
	response.SetContent(NewRPCError(status, strings.Join(errors, ";")))
	responseData, _ = response.Marshal()
	return responseData
}

// HandleFunc creates a server endpoint for yarf using the handler function, the name of function will be on the format "namespace.FunctionName"
// e.g. my-namespace.Add, if a function named Add is passed into the function
func (s *Server) HandleFunc(handler func(request *Msg, response *Msg) error) {
	parts := strings.Split(runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(), ".")
	function := parts[len(parts)-1]
	s.Handle(function, handler)
}

// Handle creates a server endpoint for yarf using the handler function, the name of function will be on the format "namespace.function"
// e.g. my-namespace.Add, if a string "Add" is passed into the function coupled with a handler function
func (s *Server) Handle(function string, handler func(request *Msg, response *Msg) error) {
	if s.namespace != "" {
		function = s.namespace + "." + function
	}
	s.transporter.Listen(function, func(requestData []byte) (responseData []byte) {

		req := Msg{}
		resp := Msg{}

		err := req.Unmarshal(requestData)
		//err := json.Unmarshal(requestData, &req)

		if err != nil {
			return toServerError(500, &resp, err.Error())
		}

		err = handler(&req, &resp)

		if err != nil {
			return toServerError(500, &resp, err.Error())
		}

		responseData, err = resp.Marshal()

		if err != nil {
			return toServerError(500, &resp, err.Error())
		}

		return responseData
	})
}
