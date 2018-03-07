package yarf

import (
	"strings"
	"reflect"
	"runtime"
)

type Server struct {
	transporter Transporter
	namespace   string
}

func NewServer(t Transporter, namespace ... string) Server {
	s := Server{}
	s.transporter = t
	if len(namespace) > 0 {
		s.namespace = strings.Join(namespace, ".")
	} else {
		s.namespace = ""
	}
	return s
}

func toServerError(status int, response *Msg, errors ... string) (responseData []byte) {

	response.SetStatus(status)
	response.SetContent(NewRPCError(status, strings.Join(errors, ";")))
	responseData, _ = response.Marshal()
	return responseData
}

func (s *Server) HandleFunc(handler func(request *Msg, response *Msg) error) {
	parts := strings.Split(runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name(), ".")
	function := parts[len(parts) - 1]
	s.Handle(function, handler)
}

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



