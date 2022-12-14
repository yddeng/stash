// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/teamChangeInfo.proto

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

type TeamChangeInfoToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamChangeInfoToS) Reset()         { *m = TeamChangeInfoToS{} }
func (m *TeamChangeInfoToS) String() string { return proto.CompactTextString(m) }
func (*TeamChangeInfoToS) ProtoMessage()    {}
func (*TeamChangeInfoToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_284830aaccf2fe5e, []int{0}
}

func (m *TeamChangeInfoToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamChangeInfoToS.Unmarshal(m, b)
}
func (m *TeamChangeInfoToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamChangeInfoToS.Marshal(b, m, deterministic)
}
func (m *TeamChangeInfoToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamChangeInfoToS.Merge(m, src)
}
func (m *TeamChangeInfoToS) XXX_Size() int {
	return xxx_messageInfo_TeamChangeInfoToS.Size(m)
}
func (m *TeamChangeInfoToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamChangeInfoToS.DiscardUnknown(m)
}

var xxx_messageInfo_TeamChangeInfoToS proto.InternalMessageInfo

type TeamChangeInfoToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamChangeInfoToC) Reset()         { *m = TeamChangeInfoToC{} }
func (m *TeamChangeInfoToC) String() string { return proto.CompactTextString(m) }
func (*TeamChangeInfoToC) ProtoMessage()    {}
func (*TeamChangeInfoToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_284830aaccf2fe5e, []int{1}
}

func (m *TeamChangeInfoToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamChangeInfoToC.Unmarshal(m, b)
}
func (m *TeamChangeInfoToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamChangeInfoToC.Marshal(b, m, deterministic)
}
func (m *TeamChangeInfoToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamChangeInfoToC.Merge(m, src)
}
func (m *TeamChangeInfoToC) XXX_Size() int {
	return xxx_messageInfo_TeamChangeInfoToC.Size(m)
}
func (m *TeamChangeInfoToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamChangeInfoToC.DiscardUnknown(m)
}

var xxx_messageInfo_TeamChangeInfoToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TeamChangeInfoToS)(nil), "message.teamChangeInfo_toS")
	proto.RegisterType((*TeamChangeInfoToC)(nil), "message.teamChangeInfo_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/teamChangeInfo.proto", fileDescriptor_284830aaccf2fe5e)
}

var fileDescriptor_284830aaccf2fe5e = []byte{
	// 103 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4d, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0x4d,
	0xcc, 0x75, 0xce, 0x48, 0xcc, 0x4b, 0x4f, 0xf5, 0xcc, 0x4b, 0xcb, 0xd7, 0x03, 0x4b, 0x0a, 0xb1,
	0x43, 0x65, 0x95, 0x44, 0xb8, 0x84, 0x50, 0x15, 0xc4, 0x97, 0xe4, 0x07, 0x63, 0x15, 0x75, 0x76,
	0x52, 0x8a, 0x52, 0xc8, 0xcc, 0xcb, 0x2c, 0xc9, 0x4c, 0xcc, 0x29, 0xc9, 0x28, 0x4a, 0x4d, 0x85,
	0xd8, 0x93, 0x9c, 0x9f, 0xa3, 0x9f, 0x5c, 0x0c, 0xb3, 0x0d, 0x10, 0x00, 0x00, 0xff, 0xff, 0x8c,
	0x23, 0xa6, 0xd4, 0x80, 0x00, 0x00, 0x00,
}
