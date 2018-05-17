package yarf

import (
	"context"
	"fmt"
)

// StatusOk rpc status ok
const StatusOk = 200

// StatusInternalError rpc status internal server error
const StatusInternalError = 500

// StatusInternalPanic rpc status internal server error when recovered from panic
const StatusInternalPanic = 501

// StatusHandlerError the handler function of request failed
const StatusHandlerError = 510

// StatusMarshalError could not marshal data
const StatusMarshalError = 550

// StatusUnmarshalError could not unmarshal data
const StatusUnmarshalError = 551

// HeaderStatus is the status header param name
const HeaderStatus = "status"

// HeaderFunction is the function name header param name
const HeaderFunction = "function"

// Msg represents a message that is being passed between client and server
type Msg struct {
	ctx        context.Context
	serializer Serializer

	Headers map[string]interface{}
	Binary  bool
	Content []byte
}

func (m *Msg) doMarshal() (data []byte, err error) {

	defer func() {
		if err != nil {
			fmt.Println("ERROR yarf.Msg Marshal", err)
		}
	}()

	data, err = m.serializer.Marshal(m)
	return

}

func (m *Msg) doUnmarshal(data []byte) (err error) {
	defer func() {
		if err != nil {
			fmt.Println("ERROR yarf.Msg Unmarshal", err)
		}
	}()
	err = m.serializer.Unmarshal(data, m)
	return
}

// Bind is userd to unmarshal/bind contetnt data to input interface
func (m *Msg) Bind(content interface{}) (err error) {

	err = m.serializer.Unmarshal(m.Content, content)
	if err != nil {
		return err
	}

	return nil
}

// Status returns the status header of the request, if one exist
func (m *Msg) Status() (status int, ok bool) {
	status64, ok := toInt(m.Headers[HeaderStatus])
	if ok {
		status = int(status64)
	}
	return
}

// Function returns the function name being called, if one exist
func (m *Msg) Function() (status string, ok bool) {
	status, ok = m.Headers[HeaderFunction].(string)
	return
}

// SetStatus sets the statues header of the message
func (m *Msg) SetStatus(code int) *Msg {
	m.SetHeader(HeaderStatus, code)
	return m
}

// Ok sets the status header to 200
func (m *Msg) Ok() *Msg {
	return m.SetStatus(StatusOk)
}

// InternalError sets the status header to 500
func (m *Msg) InternalError() *Msg {
	return m.SetStatus(StatusInternalError)
}

// SetHeader sets a generic header of the message
func (m *Msg) SetHeader(key string, value interface{}) *Msg {
	if m.Headers == nil {
		m.Headers = map[string]interface{}{}
	}

	m.Headers[key] = value
	return m
}

// SetContent sets the input interface as the content of the message
func (m *Msg) SetContent(content interface{}) *Msg {
	var err error
	m.Content, err = m.serializer.Marshal(content)
	m.Binary = false

	if err != nil {
		fmt.Println("Could not set Content", err)
	}

	return m
}

// SetBinaryContent sets the input data as content of the message
func (m *Msg) SetBinaryContent(content []byte) *Msg {
	m.Content = content
	m.Binary = true
	return m
}

// Context returns the context of the message. This is primarily for use on the server side, in order to monitor Done from client side
func (m *Msg) Context() context.Context {
	return m.ctx
}

// SetParam sets a param in the params header of the message. Which later provides helper methods of de/serializations and defaults.
func (m *Msg) SetParam(key string, value interface{}) *Msg {
	if m.Headers == nil {
		m.Headers = map[string]interface{}{}
	}
	if m.Headers["params"] == nil {
		m.Headers["params"] = map[string]interface{}{}
	}

	m.Headers["params"].(map[string]interface{})[key] = value
	return m
}

// SetParams params in the params header of the message. Which later provides helper methods of de/serializations and defaults.
func (m *Msg) SetParams(params ...Param) *Msg {
	if m.Headers == nil {
		m.Headers = map[string]interface{}{}
	}
	if m.Headers["params"] == nil {
		m.Headers["params"] = map[string]interface{}{}
	}

	for _, param := range params {
		m.Headers["params"].(map[string]interface{})[param.key] = param.value
	}
	return m
}

//Param receives a param from the params header, it is wrapped in a param struct which implements helper methods in how to access params.
func (m *Msg) Param(key string) *Param {
	p := Param{key: key}
	params := m.Headers["params"]

	if params == nil {
		return &p
	}

	pp, ok := params.(map[string]interface{})
	if !ok {
		return &p
	}

	p.value = pp[key]

	return &p
}
