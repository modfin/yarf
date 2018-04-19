package yarf

import (
	"context"
)

func NewClient(t Transporter) Client {
	s := Client{}
	s.transporter = t

	return s
}

type Client struct {
	transporter Transporter
}

func (c *Client) Call(function string, req interface{}, resp interface{}) error {
	return c.Request(function).
		Content(req).
		Bind(resp).
		Exec().
		Done()
}

func (c *Client) Request(function string) *RPC {
	return &RPC{
		client:      c,
		function:    function,
		requestMsg:  &Msg{},
		responseMsg: &Msg{},
		done:        make(chan bool),
	}
}

type RPC struct {
	client   *Client
	function string
	context  context.Context

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

func (r *RPC) Content(requestData interface{}) *RPC {
	r.requestMsg.SetContent(requestData)
	return r
}
func (r *RPC) BinaryContent(data []byte) *RPC {
	r.requestMsg.SetBinaryContent(data)
	return r
}

func (r *RPC) Context(context context.Context) *RPC {
	r.context = context
	return r
}

func (r *RPC) Callback(callback func(*Msg), errorCallback func(error)) *RPC {
	r.callback = callback
	r.errorCallback = errorCallback
	return r
}

func (r *RPC) SetChannel(channel chan (*Msg), errorChannel chan (error)) *RPC {
	r.channel = channel
	r.errorChannel = errorChannel
	return r
}

func (r *RPC) MkChannel() *RPC {
	r.channel = make(chan *Msg)
	r.errorChannel = make(chan error)
	return r
}

func (r *RPC) Bind(content interface{}) *RPC {
	r.responseMsgContent = content
	return r
}

func (r *RPC) SetParam(key string, value interface{}) *RPC {

	r.requestMsg.SetParam(key, value)

	return r
}

func (r *RPC) Exec() *RPC {

	if r.context == nil {
		r.context = context.Background()
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
		respBytes, r.err = r.client.transporter.Call(r.context, r.function, reqBytes)

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

func (r *RPC) Get() (*Msg, error) {
	r.Done()
	return r.responseMsg, r.err
}

func (r *RPC) Channel() (channel chan (*Msg), errorChannel chan (error)) {
	return r.channel, r.errorChannel
}

func (r *RPC) Done() error {
	if !r.isDone {
		r.isDone = <-r.done
	}
	return r.err
}

func (r *RPC) Error() error {
	return r.err
}
