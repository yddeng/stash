// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss/proto/ssmessage/worldBroadcastToMap.proto

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

type WorldBroadcastToMap struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WorldBroadcastToMap) Reset()         { *m = WorldBroadcastToMap{} }
func (m *WorldBroadcastToMap) String() string { return proto.CompactTextString(m) }
func (*WorldBroadcastToMap) ProtoMessage()    {}
func (*WorldBroadcastToMap) Descriptor() ([]byte, []int) {
	return fileDescriptor_7528a1b56454d194, []int{0}
}

func (m *WorldBroadcastToMap) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorldBroadcastToMap.Unmarshal(m, b)
}
func (m *WorldBroadcastToMap) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorldBroadcastToMap.Marshal(b, m, deterministic)
}
func (m *WorldBroadcastToMap) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorldBroadcastToMap.Merge(m, src)
}
func (m *WorldBroadcastToMap) XXX_Size() int {
	return xxx_messageInfo_WorldBroadcastToMap.Size(m)
}
func (m *WorldBroadcastToMap) XXX_DiscardUnknown() {
	xxx_messageInfo_WorldBroadcastToMap.DiscardUnknown(m)
}

var xxx_messageInfo_WorldBroadcastToMap proto.InternalMessageInfo

func init() {
	proto.RegisterType((*WorldBroadcastToMap)(nil), "ssmessage.worldBroadcastToMap")
}

func init() {
	proto.RegisterFile("ss/proto/ssmessage/worldBroadcastToMap.proto", fileDescriptor_7528a1b56454d194)
}

var fileDescriptor_7528a1b56454d194 = []byte{
	// 104 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x29, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x2e, 0xce, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f,
	0xcf, 0x2f, 0xca, 0x49, 0x71, 0x2a, 0xca, 0x4f, 0x4c, 0x49, 0x4e, 0x2c, 0x2e, 0x09, 0xc9, 0xf7,
	0x4d, 0x2c, 0xd0, 0x03, 0xab, 0x11, 0xe2, 0x84, 0x2b, 0x52, 0x12, 0xe5, 0x12, 0xc6, 0xa2, 0xce,
	0x49, 0x25, 0x4a, 0x29, 0x33, 0x2f, 0xb3, 0x24, 0x33, 0x31, 0xa7, 0x24, 0xa3, 0x28, 0x35, 0x15,
	0x62, 0x76, 0x72, 0x7e, 0x8e, 0x7e, 0x71, 0x31, 0xc2, 0x06, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x9c, 0x08, 0x81, 0xf2, 0x76, 0x00, 0x00, 0x00,
}