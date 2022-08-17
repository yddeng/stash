// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss/proto/ssmessage/relayUnreachable.proto

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

type RelayUnreachable struct {
	PeerAddr             *string  `protobuf:"bytes,1,opt,name=peerAddr" json:"peerAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RelayUnreachable) Reset()         { *m = RelayUnreachable{} }
func (m *RelayUnreachable) String() string { return proto.CompactTextString(m) }
func (*RelayUnreachable) ProtoMessage()    {}
func (*RelayUnreachable) Descriptor() ([]byte, []int) {
	return fileDescriptor_f64de3a0fe86d79f, []int{0}
}

func (m *RelayUnreachable) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RelayUnreachable.Unmarshal(m, b)
}
func (m *RelayUnreachable) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RelayUnreachable.Marshal(b, m, deterministic)
}
func (m *RelayUnreachable) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RelayUnreachable.Merge(m, src)
}
func (m *RelayUnreachable) XXX_Size() int {
	return xxx_messageInfo_RelayUnreachable.Size(m)
}
func (m *RelayUnreachable) XXX_DiscardUnknown() {
	xxx_messageInfo_RelayUnreachable.DiscardUnknown(m)
}

var xxx_messageInfo_RelayUnreachable proto.InternalMessageInfo

func (m *RelayUnreachable) GetPeerAddr() string {
	if m != nil && m.PeerAddr != nil {
		return *m.PeerAddr
	}
	return ""
}

func init() {
	proto.RegisterType((*RelayUnreachable)(nil), "ssmessage.relayUnreachable")
}

func init() {
	proto.RegisterFile("ss/proto/ssmessage/relayUnreachable.proto", fileDescriptor_f64de3a0fe86d79f)
}

var fileDescriptor_f64de3a0fe86d79f = []byte{
	// 121 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2c, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x2e, 0xce, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f,
	0x4a, 0xcd, 0x49, 0xac, 0x0c, 0xcd, 0x2b, 0x4a, 0x4d, 0x4c, 0xce, 0x48, 0x4c, 0xca, 0x49, 0xd5,
	0x03, 0x2b, 0x10, 0xe2, 0x84, 0xab, 0x50, 0xd2, 0xe3, 0x12, 0x40, 0x57, 0x24, 0x24, 0xc5, 0xc5,
	0x51, 0x90, 0x9a, 0x5a, 0xe4, 0x98, 0x92, 0x52, 0x24, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0x04,
	0xe7, 0x3b, 0xa9, 0x44, 0x29, 0x65, 0xe6, 0x65, 0x96, 0x64, 0x26, 0xe6, 0x94, 0x64, 0x14, 0xa5,
	0xa6, 0x42, 0x6c, 0x4c, 0xce, 0xcf, 0xd1, 0x2f, 0x2e, 0x46, 0xd8, 0x0b, 0x08, 0x00, 0x00, 0xff,
	0xff, 0x90, 0xaf, 0x7f, 0x62, 0x8c, 0x00, 0x00, 0x00,
}
