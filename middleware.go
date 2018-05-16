package yarf

// NextMiddleware simple alias for func() error
type NextMiddleware func() error

// Middleware the function that shall be implemented for a middleware
type Middleware func(request *Msg, response *Msg, next NextMiddleware) error

func processMiddleware(req *Msg, resp *Msg, handler func(request *Msg, response *Msg) error, middleware ...Middleware) error {

	if middleware == nil || len(middleware) == 0 {
		return handler(req, resp)
	}

	return middleware[0](req, resp, func() error {
		return processMiddleware(req, resp, handler, middleware[1:]...)
	})

}
