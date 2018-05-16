package yarf

import (
	"context"
	"sync"
)

const (
	initState = iota
	execState
	callState
	respState
)

// NewClient create a new yarf client using a specific transporter
func NewClient(t Transporter) Client {
	s := Client{}
	s.transporter = t

	return s
}

// Client is a struct wrapping a transporting layer and methods for using yarf
type Client struct {
	transporter Transporter
	middleware []Middleware
}


func (c *Client) WithMiddleware(middleware ...Middleware) {
	c.middleware = append(c.middleware, middleware...)
}

// Call performs a request from function name, and req param. The response is unmarshaled to resp
func (c *Client) Call(function string, req interface{}, resp interface{}) error {
	return c.Request(function).
		Content(req).
		Bind(resp).
		Exec().
		Done()
}

// Request creates a request builder in yarf
func (c *Client) Request(function string) *RPC {
	return &RPC{
		client:      c,
		function:    function,
		requestMsg:  &Msg{},
		responseMsg: &Msg{},
		state:       initState,
		done:        make(chan bool),
	}
}

// RPC represents a request to in yarf and is a builder
type RPC struct {
	client   *Client
	function string
	ctx      context.Context

	mutex sync.Mutex
	state int

	requestMsg         *Msg
	responseMsg        *Msg
	responseMsgContent interface{}

	callback      func(*Msg)
	errorCallback func(error)

	channel      chan *Msg
	errorChannel chan error

	err    error
	done   chan bool
	isDone bool
}

// Content sets requests content, it does nothing if called after Exec()
func (r *RPC) Content(requestData interface{}) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != initState {
		return r
	}

	r.requestMsg.SetContent(requestData)
	return r
}

// BinaryContent sets requests content as a binary format and marshaling will not be preformed, it does nothing if called after Exec()
func (r *RPC) BinaryContent(data []byte) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != initState {
		return r
	}

	r.requestMsg.SetBinaryContent(data)
	return r
}

// WithContext sets context of request for outside control, it does nothing if called after Exec()
func (r *RPC) WithContext(ctx context.Context) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != initState {
		return r
	}

	r.ctx = ctx
	return r
}

// WithCallback sets a callback function that will be called on success or failure, it does nothing if called after Exec()
func (r *RPC) WithCallback(callback func(*Msg), errorCallback func(error)) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != initState {
		return r
	}

	r.callback = callback
	r.errorCallback = errorCallback
	return r
}

// UseChannels creates a chan *Msg and a chan error which can be used for a non blocking context.
// The channels creaded is closed once the request is completed, it does nothing if called after Exec()
func (r *RPC) UseChannels() *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != initState {
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
	if r.state != initState {
		return r
	}

	r.responseMsgContent = content
	return r
}

// SetParam set a param that can be read by server side, like a query param in http requests, it does nothing if called after Exec()
func (r *RPC) SetParam(key string, value interface{}) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if r.state != initState {
		return r
	}

	r.requestMsg.SetParam(key, value)
	return r
}

// Exec perform rpc request
func (r *RPC) Exec(middleware...Middleware) *RPC {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.state = execState

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

		r.state = callState

		r.err = processMiddleware(r.requestMsg, r.responseMsg, func(request *Msg, response *Msg) error{

			var reqBytes []byte
			reqBytes, err := request.Marshal()

			if err != nil {
				return err
			}

			var respBytes []byte
			respBytes, err = r.client.transporter.Call(r.ctx, r.function, reqBytes)

			if err != nil {
				return err
			}

			err = response.Unmarshal(respBytes)
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
		}, append(r.client.middleware, middleware...)...)

		r.state = respState

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

	return r
}

// Get wait for request to be done before returning with the resulting message
func (r *RPC) Get() (*Msg, error) {
	r.Done()
	return r.responseMsg, r.err
}

// Channels returns channel associated with the request, these are created if UseChannels() is called before Exec()
func (r *RPC) Channels() (channel chan *Msg, errorChannel chan error) {
	return r.channel, r.errorChannel
}

// Done waits until the rpc request is done and has returned a result
func (r *RPC) Done() error {
	if !r.isDone {
		r.isDone = <-r.done
	}
	return r.err
}

// Error return the error of the request, if any at the point of calling it.
func (r *RPC) Error() error {
	return r.err
}
