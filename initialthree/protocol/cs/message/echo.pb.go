// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/echo.proto

package message

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EchoToS struct {
	Msg                  *string  `protobuf:"bytes,1,req,name=msg" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoToS) Reset()         { *m = EchoToS{} }
func (m *EchoToS) String() string { return proto.CompactTextString(m) }
func (*EchoToS) ProtoMessage()    {}
func (*EchoToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4f19d1d89c76a95, []int{0}
}

func (m *EchoToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoToS.Unmarshal(m, b)
}
func (m *EchoToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoToS.Marshal(b, m, deterministic)
}
func (m *EchoToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoToS.Merge(m, src)
}
func (m *EchoToS) XXX_Size() int {
	return xxx_messageInfo_EchoToS.Size(m)
}
func (m *EchoToS) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoToS.DiscardUnknown(m)
}

var xxx_messageInfo_EchoToS proto.InternalMessageInfo

func (m *EchoToS) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

type EchoToC struct {
	Msg                  *string  `protobuf:"bytes,2,req,name=msg" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EchoToC) Reset()         { *m = EchoToC{} }
func (m *EchoToC) String() string { return proto.CompactTextString(m) }
func (*EchoToC) ProtoMessage()    {}
func (*EchoToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_e4f19d1d89c76a95, []int{1}
}

func (m *EchoToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EchoToC.Unmarshal(m, b)
}
func (m *EchoToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EchoToC.Marshal(b, m, deterministic)
}
func (m *EchoToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EchoToC.Merge(m, src)
}
func (m *EchoToC) XXX_Size() int {
	return xxx_messageInfo_EchoToC.Size(m)
}
func (m *EchoToC) XXX_DiscardUnknown() {
	xxx_messageInfo_EchoToC.DiscardUnknown(m)
}

var xxx_messageInfo_EchoToC proto.InternalMessageInfo

func (m *EchoToC) GetMsg() string {
	if m != nil && m.Msg != nil {
		return *m.Msg
	}
	return ""
}

func init() {
	proto.RegisterType((*EchoToS)(nil), "message.echo_toS")
	proto.RegisterType((*EchoToC)(nil), "message.echo_toC")
}

func init() { proto.RegisterFile("cs/proto/message/echo.proto", fileDescriptor_e4f19d1d89c76a95) }

var fileDescriptor_e4f19d1d89c76a95 = []byte{
	// 113 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4e, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x4d, 0xce,
	0xc8, 0xd7, 0x03, 0x0b, 0x09, 0xb1, 0x43, 0xc5, 0x94, 0x64, 0xb8, 0x38, 0x40, 0xc2, 0xf1, 0x25,
	0xf9, 0xc1, 0x42, 0x02, 0x5c, 0xcc, 0xb9, 0xc5, 0xe9, 0x12, 0x8c, 0x0a, 0x4c, 0x1a, 0x9c, 0x41,
	0x20, 0x26, 0x92, 0xac, 0x33, 0x4c, 0x96, 0x09, 0x2e, 0xeb, 0xa4, 0x14, 0xa5, 0x90, 0x99, 0x97,
	0x59, 0x92, 0x99, 0x98, 0x53, 0x92, 0x51, 0x94, 0x9a, 0x0a, 0xb1, 0x2d, 0x39, 0x3f, 0x47, 0x3f,
	0xb9, 0x18, 0x66, 0x27, 0x20, 0x00, 0x00, 0xff, 0xff, 0xc4, 0x2d, 0x8e, 0x19, 0x86, 0x00, 0x00,
	0x00,
}
