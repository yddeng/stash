// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss/proto/ssmessage/functionSwitchReload.proto

package ssmessage

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

type FunctionSwitchReload struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FunctionSwitchReload) Reset()         { *m = FunctionSwitchReload{} }
func (m *FunctionSwitchReload) String() string { return proto.CompactTextString(m) }
func (*FunctionSwitchReload) ProtoMessage()    {}
func (*FunctionSwitchReload) Descriptor() ([]byte, []int) {
	return fileDescriptor_521291f543de2ca5, []int{0}
}

func (m *FunctionSwitchReload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FunctionSwitchReload.Unmarshal(m, b)
}
func (m *FunctionSwitchReload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FunctionSwitchReload.Marshal(b, m, deterministic)
}
func (m *FunctionSwitchReload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FunctionSwitchReload.Merge(m, src)
}
func (m *FunctionSwitchReload) XXX_Size() int {
	return xxx_messageInfo_FunctionSwitchReload.Size(m)
}
func (m *FunctionSwitchReload) XXX_DiscardUnknown() {
	xxx_messageInfo_FunctionSwitchReload.DiscardUnknown(m)
}

var xxx_messageInfo_FunctionSwitchReload proto.InternalMessageInfo

func init() {
	proto.RegisterType((*FunctionSwitchReload)(nil), "ssmessage.functionSwitchReload")
}

func init() {
	proto.RegisterFile("ss/proto/ssmessage/functionSwitchReload.proto", fileDescriptor_521291f543de2ca5)
}

var fileDescriptor_521291f543de2ca5 = []byte{
	// 105 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2d, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x2e, 0xce, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f,
	0x2b, 0xcd, 0x4b, 0x2e, 0xc9, 0xcc, 0xcf, 0x0b, 0x2e, 0xcf, 0x2c, 0x49, 0xce, 0x08, 0x4a, 0xcd,
	0xc9, 0x4f, 0x4c, 0xd1, 0x03, 0x2b, 0x12, 0xe2, 0x84, 0xab, 0x52, 0x12, 0xe3, 0x12, 0xc1, 0xa6,
	0xd0, 0x49, 0x25, 0x4a, 0x29, 0x33, 0x2f, 0xb3, 0x24, 0x33, 0x31, 0xa7, 0x24, 0xa3, 0x28, 0x35,
	0x15, 0x62, 0x7a, 0x72, 0x7e, 0x8e, 0x7e, 0x71, 0x31, 0xc2, 0x0e, 0x40, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x2c, 0x70, 0x3c, 0x96, 0x78, 0x00, 0x00, 0x00,
}
