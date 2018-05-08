package yarf

import (
	"context"
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
		done:        make(chan bool),
	}
}

// RPC represents a request to in yarf and is a builder
type RPC struct {
	client   *Client
	function string
	ctx      context.Context

	requestMsg         *Msg
	responseMsg        *Msg
	responseMsgContent interface{}

	callback      func(*Msg)
	errorCallback func(error)

	channel      chan (*Msg)
	errorChannel chan (error)

	err    error
	done   chan (bool)
	isDone bool
}

// Content sets requests content
func (r *RPC) Content(requestData interface{}) *RPC {
	r.requestMsg.SetContent(requestData)
	return r
}

// BinaryContent sets requests content as a binary format and marshaling will not be preformed
func (r *RPC) BinaryContent(data []byte) *RPC {
	r.requestMsg.SetBinaryContent(data)
	return r
}

// Context sets context of request for outside control
func (r *RPC) Context(ctx context.Context) *RPC {
	r.ctx = ctx
	return r
}

// Callback sets a callback function that will be called on success or failure
func (r *RPC) Callback(callback func(*Msg), errorCallback func(error)) *RPC {
	r.callback = callback
	r.errorCallback = errorCallback
	return r
}

// SetChannel sets a channels that response and error will be passed to
func (r *RPC) SetChannel(channel chan (*Msg), errorChannel chan (error)) *RPC {
	r.channel = channel
	r.errorChannel = errorChannel
	return r
}

// MkChannel creates channels that can be used if external onec is not required
func (r *RPC) MkChannel() *RPC {
	r.channel = make(chan *Msg)
	r.errorChannel = make(chan error)
	return r
}

// Bind will unmarshal response into interface passed into method
func (r *RPC) Bind(content interface{}) *RPC {
	r.responseMsgContent = content
	return r
}

// SetParam set a param that can be read by server side, like a query param in http requests
func (r *RPC) SetParam(key string, value interface{}) *RPC {

	r.requestMsg.SetParam(key, value)

	return r
}

// Exec perform rpc request
func (r *RPC) Exec() *RPC {

	if r.ctx == nil {
		r.ctx = context.Background()
	}

	var reqBytes []byte
	reqBytes, r.err = r.requestMsg.Marshal()
	if r.err != nil {
		return r
	}

	go func() {
		defer func() {

			if r.err != nil {

				if r.errorChannel != nil {
					r.errorChannel <- r.err
				}

				if r.errorCallback != nil {
					r.errorCallback(r.err)
				}
			}

			r.done <- true

		}()

		var respBytes []byte
		respBytes, r.err = r.client.transporter.Call(r.ctx, r.function, reqBytes)

		if r.err != nil {
			return
		}

		r.err = r.responseMsg.Unmarshal(respBytes)
		if r.err != nil {
			return
		}

		if s, ok := r.responseMsg.Status(); s >= 500 && ok {
			err := RPCError{}
			r.responseMsg.Bind(&err)
			r.err = err
			return
		}

		if r.responseMsgContent != nil {
			r.err = r.responseMsg.Bind(r.responseMsgContent)
		}

		if r.channel != nil {
			r.channel <- r.responseMsg
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

// Channel returns channel associated with the request and is most likely to be used in conjunction with MkChannel()
func (r *RPC) Channel() (channel chan (*Msg), errorChannel chan (error)) {
	return r.channel, r.errorChannel
}

// Done waits until the rpc request is done and has returned a result
func (r *RPC) Done() error {
	if !r.isDone {
		r.isDone = <-r.done
	}
	return r.err
}

// Error return the error of the request, if any
func (r *RPC) Error() error {
	return r.err
}
