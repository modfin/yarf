package yarf

import (
	"context"
	"sync"
)

const (
	builderState = iota
	transitState
	requestState
	responseState
)

// NewClient create a new yarf client using a specific transporter
func NewClient(t CallTransporter) Client {
	s := Client{}
	s.transporter = t
	s.serializer = defaultSerializer()
	return s
}

// Client is a struct wrapping a transporting layer and methods for using yarf
type Client struct {
	transporter CallTransporter
	middleware  []Middleware
	serializer  Serializer
}

// WithMiddleware adds middleware to client request for pre and post processing
func (c *Client) WithMiddleware(middleware ...Middleware) {
	c.middleware = append(c.middleware, middleware...)
}

// WithSerializer setts the serializer used for transport.
func (c *Client) WithSerializer(serializer Serializer) {
	c.serializer = serializer
}

// Call is a short hand performs a request from function name, and req param. The response is unmarshaled into resp
func (c *Client) Call(function string, requestData interface{}, responseData interface{}, requestParams ...Param) error {
	return c.Request(function).
		SetParams(requestParams...).
		Content(requestData).
		Bind(responseData).
		Done()
}

// Request creates a request builder in yarf
func (c *Client) Request(function string) *RPC {
	return &RPC{
		client:      c,
		function:    function,
		requestMsg:  &Msg{serializer: c.serializer},
		responseMsg: &Msg{serializer: c.serializer},
		state:       builderState,
		done:        make(chan bool),
	}
}

// RPCTransit represents a request to in yarf when in transit and the purpose of it is to restricts the function that are
// allowed to be called when the request in transit in order to expose internal state and what can be done.
type RPCTransit struct {
	rpc *RPC
}

// RPC represents a request to in yarf and is used to build a request using the builder pattern
type RPC struct {
	client   *Client
	function string
	ctx      context.Context

	mutex sync.Mutex
	state int

	middleware []Middleware

	requestMsg         *Msg
	responseMsg        *Msg
	responseMsgContent interface{}

	callback      func(*Msg)
	errorCallback func(error)

	channel      chan *Msg
	errorChannel chan error

	err       error
	done      chan bool
	isDone    bool
	doneMutex sync.Mutex
}

// Content sets requests content, it does nothing if called after Exec()
func (r *RPC) Content(requestData interface{}) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != builderState {
		return r
	}

	r.requestMsg.SetContent(requestData)
	return r
}

// BinaryContent sets requests content as a binary format and marshaling will not be preformed, it does nothing if called after Exec()
func (r *RPC) BinaryContent(data []byte) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != builderState {
		return r
	}

	r.requestMsg.SetBinaryContent(data)
	return r
}

// WithContext sets context of request for outside control, it does nothing if called after Exec()
func (r *RPC) WithContext(ctx context.Context) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != builderState {
		return r
	}

	r.ctx = ctx
	return r
}

// WithCallback sets a callback function that will be called on success or failure, it does nothing if called after Exec()
func (r *RPC) WithCallback(callback func(*Msg), errorCallback func(error)) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != builderState {
		return r
	}

	r.callback = callback
	r.errorCallback = errorCallback
	return r
}

//WithMiddleware adds middleware to specific request.
func (r *RPC) WithMiddleware(middleware ...Middleware) *RPC {
	r.middleware = append(r.middleware, middleware...)
	return r
}

// UseChannels creates a chan *Msg and a chan error which can be used for a non blocking context.
// The channels creaded is closed once the request is completed, it does nothing if called after Exec()
func (r *RPC) UseChannels() *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != builderState {
		return r
	}

	r.channel = make(chan *Msg)
	r.errorChannel = make(chan error)
	return r
}

// Bind will unmarshal response into interface passed into method, it does nothing if called after Exec()
func (r *RPC) Bind(content interface{}) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != builderState {
		return r
	}

	r.responseMsgContent = content
	return r
}

// SetParam set a param that can be read by server side, like a query param in http requests, it does nothing if called after Exec()
func (r *RPC) SetParam(key string, value interface{}) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != builderState {
		return r
	}

	r.requestMsg.SetParam(key, value)
	return r
}

// SetParams set params provided that can be read by server side, like a query param in http requests, it does nothing if called after Exec()
func (r *RPC) SetParams(params ...Param) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != builderState {
		return r
	}

	r.requestMsg.SetParams(params...)

	return r
}

// Exec perform rpc request and return a RPC transit struct. Done(), Get() and Channels() will call Exec() if it has not been called "manually".
func (r *RPC) Exec() *RPCTransit {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.state != builderState {
		return &RPCTransit{r}
	}

	r.state = transitState

	if r.ctx == nil {
		r.ctx = context.Background()
	}

	r.requestMsg.SetHeader(HeaderFunction, r.function)

	var cancel func()
	r.ctx, cancel = context.WithCancel(r.ctx)

	go func() {
		r.mutex.Lock()
		defer r.mutex.Unlock()
		defer func() {
			if r.err != nil {
				if r.errorChannel != nil {
					r.errorChannel <- r.err
				}

				if r.errorCallback != nil {
					r.errorCallback(r.err)
				}
			}
			cancel()
			r.done <- true
			close(r.done)
			if r.errorChannel != nil {
				close(r.errorChannel)
			}
		}()

		r.state = requestState

		r.err = processMiddleware(r.requestMsg, r.responseMsg, toClientRequestHandler(r), append(r.client.middleware, r.middleware...)...)

		r.state = responseState

		if r.err != nil {
			return
		}

		if r.channel != nil {
			r.channel <- r.responseMsg
			close(r.channel)
		}

		if r.callback != nil {
			r.callback(r.responseMsg)
		}

	}()

	return &RPCTransit{r}
}

func toClientRequestHandler(r *RPC) func(request *Msg, response *Msg) error {
	return func(request *Msg, response *Msg) error {

		var reqBytes []byte
		reqBytes, err := request.doMarshal()

		if err != nil {
			return err
		}

		var respBytes []byte
		respBytes, err = r.client.transporter.Call(r.ctx, r.function, reqBytes)

		if err != nil {
			return err
		}

		err = response.doUnmarshal(respBytes)
		if err != nil {
			return err
		}

		if s, ok := response.Status(); s >= 500 && ok {
			err := RPCError{}
			response.Bind(&err)
			return err
		}

		if r.responseMsgContent != nil {
			err = response.Bind(r.responseMsgContent)
			return err
		}

		return nil
	}
}

// Get wait for request to be done before returning with the resulting message. If Exec() has not been called, Get() will
// call it.
func (r *RPC) Get() (*Msg, error) {
	r.Done()
	return r.responseMsg, r.err
}

// Get wait for request to be done before returning with the resulting message. If Exec() has not been called, Get() will
// call it.
func (r *RPCTransit) Get() (*Msg, error) {
	return r.rpc.Get()
}

// Channels returns channel associated with the request, these are created if UseChannels() is called before Exec().
// Channels() will call UseChannel() and then Exec() if Exec() has not been called.
func (r *RPC) Channels() (<-chan *Msg, <-chan error) {
	r.UseChannels()
	r.Exec()

	return r.channel, r.errorChannel
}

// Channels returns channel associated with the request, these are created if UseChannels() is called before Exec().
// Channels() will call UseChannel() and then Exec() if Exec() has not been called.
func (r *RPCTransit) Channels() (<-chan *Msg, <-chan error) {
	return r.rpc.Channels()
}

// Done waits until the rpc request is done and has returned a result. Done will call Exec(), which performs the request,
// if Exec() has not been called prior to Done being called
func (r *RPC) Done() error {
	r.Exec()

	r.doneMutex.Lock()
	if !r.isDone {
		r.isDone = <-r.done
	}
	r.doneMutex.Unlock()

	return r.err
}

// Done waits until the rpc request is done and has returned a result. Done will call Exec(), which performs the request,
// if Exec() has not been called prior to Done being called
func (r *RPCTransit) Done() error {
	return r.rpc.Done()
}

// Error return the error of the request, if any at the point of calling it.
func (r *RPC) Error() error {
	return r.err
}

// Error return the error of the request, if any at the point of calling it.
func (r *RPCTransit) Error() error {
	return r.rpc.Error()
}
