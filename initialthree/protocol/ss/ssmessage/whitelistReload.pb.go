// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss/proto/ssmessage/whitelistReload.proto

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

type WhitelistReload struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WhitelistReload) Reset()         { *m = WhitelistReload{} }
func (m *WhitelistReload) String() string { return proto.CompactTextString(m) }
func (*WhitelistReload) ProtoMessage()    {}
func (*WhitelistReload) Descriptor() ([]byte, []int) {
	return fileDescriptor_25456659dba0aec1, []int{0}
}

func (m *WhitelistReload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WhitelistReload.Unmarshal(m, b)
}
func (m *WhitelistReload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WhitelistReload.Marshal(b, m, deterministic)
}
func (m *WhitelistReload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WhitelistReload.Merge(m, src)
}
func (m *WhitelistReload) XXX_Size() int {
	return xxx_messageInfo_WhitelistReload.Size(m)
}
func (m *WhitelistReload) XXX_DiscardUnknown() {
	xxx_messageInfo_WhitelistReload.DiscardUnknown(m)
}

var xxx_messageInfo_WhitelistReload proto.InternalMessageInfo

func init() {
	proto.RegisterType((*WhitelistReload)(nil), "ssmessage.whitelistReload")
}

func init() {
	proto.RegisterFile("ss/proto/ssmessage/whitelistReload.proto", fileDescriptor_25456659dba0aec1)
}

var fileDescriptor_25456659dba0aec1 = []byte{
	// 99 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x28, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x2e, 0xce, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f,
	0xcf, 0xc8, 0x2c, 0x49, 0xcd, 0xc9, 0x2c, 0x2e, 0x09, 0x4a, 0xcd, 0xc9, 0x4f, 0x4c, 0xd1, 0x03,
	0xcb, 0x0b, 0x71, 0xc2, 0x15, 0x28, 0x09, 0x72, 0xf1, 0xa3, 0xa9, 0x71, 0x52, 0x89, 0x52, 0xca,
	0xcc, 0xcb, 0x2c, 0xc9, 0x4c, 0xcc, 0x29, 0xc9, 0x28, 0x4a, 0x4d, 0x85, 0x98, 0x99, 0x9c, 0x9f,
	0xa3, 0x5f, 0x5c, 0x8c, 0x30, 0x19, 0x10, 0x00, 0x00, 0xff, 0xff, 0x94, 0xf2, 0x29, 0xf9, 0x6e,
	0x00, 0x00, 0x00,
}
