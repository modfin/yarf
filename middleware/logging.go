package middleware

import (
	"bitbucket.org/modfin/yarf"
	"encoding/json"
	"fmt"
	"time"
)

type loggingContent struct {
	At       time.Time       `json:"at"`
	Function string          `json:"function"`
	Headers  json.RawMessage `json:"headers"`
}

// Logging is a middleware for logging requests on the server side.
func Logging() func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

	return func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {
		l := loggingContent{}
		l.Function, _ = request.Function()
		l.At = time.Now()
		l.Headers, _ = json.Marshal(request.Headers)

		b, _ := json.Marshal(l)

		fmt.Println(string(b))

		err := next()

		// TODO log errors...

		return err
	}

}
