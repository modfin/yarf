package middleware

import (
	"bitbucket.org/modfin/yarf"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"time"
)

// OpenTracing is a middleware for doing open tracing
func OpenTracing(prefix string) func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

	return func(request *yarf.Msg, response *yarf.Msg, next yarf.NextMiddleware) error {

		funcname, _ := request.Function()

		start := time.Now()

		span := opentracing.StartSpan(prefix + funcname)
		defer span.Finish()

		uuid, _ := request.UUID()
		span.LogFields(
			log.String("type", "request"),
			log.String("uuid", uuid),
			log.Int("request_body_len", len(request.Content)),
		)

		err := next()

		uuid, _ = response.UUID()
		span.LogFields(
			log.String("type", "response"),
			log.String("uuid", uuid),
			log.Float64("response_duration", float64(time.Now().Sub(start).Nanoseconds())/1000000.0),
			log.Int("response_body_len", len(response.Content)),
		)

		return err
	}

}
