// Code generated by protoc-gen-go. DO NOT EDIT.
// source: tmsg.proto

/*
Package transport is a generated protocol buffer package.

It is generated from these files:
	tmsg.proto

It has these top-level messages:
	TMSG
*/
package yarf

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type TMSG struct {
	Binary           *bool  `protobuf:"varint,1,req,name=binary" json:"binary,omitempty"`
	Headers          []byte `protobuf:"bytes,2,opt,name=headers" json:"headers,omitempty"`
	Content          []byte `protobuf:"bytes,3,opt,name=content" json:"content,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *TMSG) Reset()                    { *m = TMSG{} }
func (m *TMSG) String() string            { return proto.CompactTextString(m) }
func (*TMSG) ProtoMessage()               {}
func (*TMSG) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *TMSG) GetBinary() bool {
	if m != nil && m.Binary != nil {
		return *m.Binary
	}
	return false
}

func (m *TMSG) GetHeaders() []byte {
	if m != nil {
		return m.Headers
	}
	return nil
}

func (m *TMSG) GetContent() []byte {
	if m != nil {
		return m.Content
	}
	return nil
}

func init() {
	proto.RegisterType((*TMSG)(nil), "transport.TMSG")
}

func init() { proto.RegisterFile("tmsg.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 108 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0xc9, 0x2d, 0x4e,
	0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2c, 0x29, 0x4a, 0xcc, 0x2b, 0x2e, 0xc8, 0x2f,
	0x2a, 0x51, 0x0a, 0xe2, 0x62, 0x09, 0xf1, 0x0d, 0x76, 0x17, 0x12, 0xe3, 0x62, 0x4b, 0xca, 0xcc,
	0x4b, 0x2c, 0xaa, 0x94, 0x60, 0x54, 0x60, 0xd2, 0xe0, 0x08, 0x82, 0xf2, 0x84, 0x24, 0xb8, 0xd8,
	0x33, 0x52, 0x13, 0x53, 0x52, 0x8b, 0x8a, 0x25, 0x98, 0x14, 0x18, 0x35, 0x78, 0x82, 0x60, 0x5c,
	0x90, 0x4c, 0x72, 0x7e, 0x5e, 0x49, 0x6a, 0x5e, 0x89, 0x04, 0x33, 0x44, 0x06, 0xca, 0x05, 0x04,
	0x00, 0x00, 0xff, 0xff, 0x2b, 0xf7, 0xe0, 0x0f, 0x6b, 0x00, 0x00, 0x00,
}
