package yarf

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
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

// HeaderUUID is the uuid header param name
const HeaderUUID = "status"

// HeaderFunction is the function name header param name
const HeaderFunction = "function"

// HeaderContentType is the function name header param name
const HeaderContentType = "content-type"

// Msg represents a message that is being passed between client and server
type Msg struct {
	ctx                context.Context
	protocolSerializer Serializer
	contentSerializer  Serializer

	builderError error

	Headers map[string]interface{}
	Content []byte
}

func (m *Msg) doMarshal() (data []byte, err error) {
	contentType := []byte(m.protocolSerializer.ContentType + "\n")
	content, err := m.protocolSerializer.Marshal(m)

	data = make([]byte, 0, len(contentType)+len(content))
	data = append(data, contentType...)
	data = append(data, content...)

	return

}

func (m *Msg) doUnmarshal(data []byte) (err error) {
	contentlen := -1

	for i := 0; i < 100 && i < len(data); i++ {

		if data[i] == 10 { // "\n"
			contentlen = i
			break
		}
	}

	if contentlen == -1 {
		return errors.New("Could not find content type")
	}

	contentType := string(data[:contentlen])

	ser, ok := serializer(contentType)

	if !ok {
		return errors.New("could not find a suitable protocolSerializer")
	}

	err = ser.Unmarshal(data[contentlen+1:], m)
	return
}

// BindContent is used to unmarshal/bind content data to input interface. It will look for a proper deserializer matching
// header content-type. Serializer can be registered by yarf.RegisterSerializer()
func (m *Msg) BindContent(content interface{}) (err error) {

	contentType, ok := m.ContentType()

	if !ok {
		return errors.New("could not find a content type to use for deserialization")
	}

	ser, ok := serializer(contentType)

	if !ok {
		return errors.New("could not find a protocolSerializer matching content type")
	}

	err = ser.Unmarshal(m.Content, content)

	return
}

// Status returns the status header of the request, if one exist
func (m *Msg) Status() (status int, ok bool) {
	status64, ok := toInt(m.Headers[HeaderStatus])
	if ok {
		status = int(status64)
	}
	return
}

// SetStatus sets the statues header of the message
func (m *Msg) SetStatus(code int) *Msg {
	m.SetHeader(HeaderStatus, code)
	return m
}

// ContentType returns the ContentType header of the request, if one exist
func (m *Msg) ContentType() (contentType string, ok bool) {
	contentType, ok = (m.Headers[HeaderContentType]).(string)
	return
}

// SetContentType sets the statues header of the message
func (m *Msg) SetContentType(contentType string) *Msg {
	m.SetHeader(HeaderContentType, contentType)
	return m
}

// Function returns the function name being called, if one exist
func (m *Msg) Function() (status string, ok bool) {
	status, ok = m.Headers[HeaderFunction].(string)
	return
}

// UUID returns the request uuid
func (m *Msg) UUID() (status string, ok bool) {
	status, ok = m.Headers[HeaderUUID].(string)
	return
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
	m.SetContentUsing(content, m.contentSerializer)
	return m
}

// SetContentUsing serializes the content with a specific contentSerializer and sets it as a binary payload
func (m *Msg) SetContentUsing(content interface{}, serializer Serializer) *Msg {
	var err error

	m.SetContentType(serializer.ContentType)
	m.Content, err = serializer.Marshal(content)

	if err != nil {
		m.builderError = err
		// TODO should panic? set an internal error?
		fmt.Println("Could not set SetContentUsing", err)
	}
	return m
}

// SetBinaryContent sets the input data as content of the message
func (m *Msg) SetBinaryContent(content []byte) *Msg {
	m.Content = content
	m.SetContentType("binary/octet-stream")
	return m
}

// Context returns the context of the message. This is primarily for use on the server side, in order to monitor Done from client side
func (m *Msg) Context() context.Context {
	return m.ctx
}

// WithContext sets the context of the message. The supplied context will replace the current one. If wrapping is intended get the current context first using Context
func (m *Msg) WithContext(ctx context.Context) *Msg {
	m.ctx = ctx
	return m
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
