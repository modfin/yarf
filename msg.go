package yarf

import (
	"encoding/json"
	"github.com/gogo/protobuf/proto"
	"fmt"
)

var STATUS_OK int = 200
var STATUS_INTERNAL_ERROR int = 500

type Msg struct {
	Headers map[string]interface{}

	Binary  bool
	Content json.RawMessage
}

func (m *Msg) Marshal() (data []byte, err error) {

	defer func() {
		if err != nil {
			fmt.Println("ERROR Unmarshal", err)
		}
	}()

	var headerBytes []byte
	var contentBytes []byte

	headerBytes, err = json.Marshal(m.Headers)
	contentBytes = m.Content

	if err != nil {
		fmt.Println("Ma", 1)
		return nil, err
	}

	tmsg := &TMSG{
		Binary: &m.Binary,
		Headers: headerBytes,
		Content: contentBytes,
	}

	data, err = proto.Marshal(tmsg)

	if err != nil {
		fmt.Println("Ma", 2)
		return nil, err
	}

	return

}

func (m *Msg) Unmarshal(data []byte) (err error) {

	defer func() {
		if err != nil {
			fmt.Println("ERROR Unmarshal", err)
		}
	}()

	tmsg := &TMSG{}

	if err = proto.Unmarshal(data, tmsg); err != nil {
		fmt.Println("UnMa", 1)
		return
	}

	m.Binary = *tmsg.Binary

	headerBytes := tmsg.Headers
	contentBytes := tmsg.Content

	err = json.Unmarshal(headerBytes, &m.Headers)
	if err != nil {
		fmt.Println("UnMa", 3)
		return
	}

	m.Content = contentBytes

	return

}

func (m *Msg) Bind(content interface{}) (err error) {

	err = json.Unmarshal(m.Content, &content)
	if err != nil {
		return err
	}

	return nil
}

func (m *Msg) Status() (status int, ok bool) {
	status, ok = m.Headers["status"].(int)
	return
}

func (m *Msg) SetStatus(code int) *Msg {
	m.SetHeader("status", code)
	return m
}

func (m *Msg) Ok() *Msg {
	return m.SetStatus(STATUS_OK)
}

func (m *Msg) InternalError() *Msg {
	return m.SetStatus(STATUS_INTERNAL_ERROR)
}

func (m *Msg) SetHeader(key string, value interface{}) *Msg {
	if m.Headers == nil {
		m.Headers = map[string]interface{}{}
	}

	m.Headers[key] = value
	return m
}

func (m *Msg) SetContent(content interface{}) (*Msg) {
	m.Content, _ = json.Marshal(content)
	m.Binary = false
	return m
}

func (m *Msg) SetBinaryContent(content []byte) (*Msg) {
	m.Content = content
	m.Binary = true
	return m
}

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

func (m *Msg) Param(key string) *Param {
	p := Param{key:key}
	params := m.Headers["params"]

	if params == nil{
		return &p
	}

	pp, ok := params.(map[string]interface{})
	if !ok {
		return &p
	}

	p.value = pp[key]


	return &p
}






