// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/heartbeat.proto

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

type HeartbeatToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatToS) Reset()         { *m = HeartbeatToS{} }
func (m *HeartbeatToS) String() string { return proto.CompactTextString(m) }
func (*HeartbeatToS) ProtoMessage()    {}
func (*HeartbeatToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3b98b1d52029661, []int{0}
}

func (m *HeartbeatToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatToS.Unmarshal(m, b)
}
func (m *HeartbeatToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatToS.Marshal(b, m, deterministic)
}
func (m *HeartbeatToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatToS.Merge(m, src)
}
func (m *HeartbeatToS) XXX_Size() int {
	return xxx_messageInfo_HeartbeatToS.Size(m)
}
func (m *HeartbeatToS) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatToS.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatToS proto.InternalMessageInfo

type HeartbeatToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatToC) Reset()         { *m = HeartbeatToC{} }
func (m *HeartbeatToC) String() string { return proto.CompactTextString(m) }
func (*HeartbeatToC) ProtoMessage()    {}
func (*HeartbeatToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_d3b98b1d52029661, []int{1}
}

func (m *HeartbeatToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatToC.Unmarshal(m, b)
}
func (m *HeartbeatToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatToC.Marshal(b, m, deterministic)
}
func (m *HeartbeatToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatToC.Merge(m, src)
}
func (m *HeartbeatToC) XXX_Size() int {
	return xxx_messageInfo_HeartbeatToC.Size(m)
}
func (m *HeartbeatToC) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatToC.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*HeartbeatToS)(nil), "message.heartbeat_toS")
	proto.RegisterType((*HeartbeatToC)(nil), "message.heartbeat_toC")
}

func init() { proto.RegisterFile("cs/proto/message/heartbeat.proto", fileDescriptor_d3b98b1d52029661) }

var fileDescriptor_d3b98b1d52029661 = []byte{
	// 98 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x48, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0xcf, 0x48, 0x4d,
	0x2c, 0x2a, 0x49, 0x4a, 0x4d, 0x2c, 0xd1, 0x03, 0x8b, 0x0b, 0xb1, 0x43, 0x25, 0x94, 0xf8, 0xb9,
	0x78, 0xe1, 0x72, 0xf1, 0x25, 0xf9, 0xc1, 0xe8, 0x02, 0xce, 0x4e, 0x4a, 0x51, 0x0a, 0x99, 0x79,
	0x99, 0x25, 0x99, 0x89, 0x39, 0x25, 0x19, 0x45, 0xa9, 0xa9, 0x10, 0x83, 0x93, 0xf3, 0x73, 0xf4,
	0x93, 0x8b, 0x61, 0xc6, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0x14, 0xe1, 0x4b, 0x8f, 0x71, 0x00,
	0x00, 0x00,
}
