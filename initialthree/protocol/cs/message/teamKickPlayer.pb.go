// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/teamKickPlayer.proto

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

type TeamKickPlayerToS struct {
	KickPlayerID         *uint64  `protobuf:"varint,1,opt,name=kickPlayerID" json:"kickPlayerID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamKickPlayerToS) Reset()         { *m = TeamKickPlayerToS{} }
func (m *TeamKickPlayerToS) String() string { return proto.CompactTextString(m) }
func (*TeamKickPlayerToS) ProtoMessage()    {}
func (*TeamKickPlayerToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_d9321325724a2f7c, []int{0}
}

func (m *TeamKickPlayerToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamKickPlayerToS.Unmarshal(m, b)
}
func (m *TeamKickPlayerToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamKickPlayerToS.Marshal(b, m, deterministic)
}
func (m *TeamKickPlayerToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamKickPlayerToS.Merge(m, src)
}
func (m *TeamKickPlayerToS) XXX_Size() int {
	return xxx_messageInfo_TeamKickPlayerToS.Size(m)
}
func (m *TeamKickPlayerToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamKickPlayerToS.DiscardUnknown(m)
}

var xxx_messageInfo_TeamKickPlayerToS proto.InternalMessageInfo

func (m *TeamKickPlayerToS) GetKickPlayerID() uint64 {
	if m != nil && m.KickPlayerID != nil {
		return *m.KickPlayerID
	}
	return 0
}

type TeamKickPlayerToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamKickPlayerToC) Reset()         { *m = TeamKickPlayerToC{} }
func (m *TeamKickPlayerToC) String() string { return proto.CompactTextString(m) }
func (*TeamKickPlayerToC) ProtoMessage()    {}
func (*TeamKickPlayerToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_d9321325724a2f7c, []int{1}
}

func (m *TeamKickPlayerToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamKickPlayerToC.Unmarshal(m, b)
}
func (m *TeamKickPlayerToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamKickPlayerToC.Marshal(b, m, deterministic)
}
func (m *TeamKickPlayerToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamKickPlayerToC.Merge(m, src)
}
func (m *TeamKickPlayerToC) XXX_Size() int {
	return xxx_messageInfo_TeamKickPlayerToC.Size(m)
}
func (m *TeamKickPlayerToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamKickPlayerToC.DiscardUnknown(m)
}

var xxx_messageInfo_TeamKickPlayerToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TeamKickPlayerToS)(nil), "message.teamKickPlayer_toS")
	proto.RegisterType((*TeamKickPlayerToC)(nil), "message.teamKickPlayer_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/teamKickPlayer.proto", fileDescriptor_d9321325724a2f7c)
}

var fileDescriptor_d9321325724a2f7c = []byte{
	// 124 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4d, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0x4d,
	0xcc, 0xf5, 0xce, 0x4c, 0xce, 0x0e, 0xc8, 0x49, 0xac, 0x4c, 0x2d, 0xd2, 0x03, 0x4b, 0x0a, 0xb1,
	0x43, 0x65, 0x95, 0x2c, 0xb8, 0x84, 0x50, 0x15, 0xc4, 0x97, 0xe4, 0x07, 0x0b, 0x29, 0x71, 0xf1,
	0x64, 0xc3, 0x45, 0x3c, 0x5d, 0x24, 0x18, 0x15, 0x18, 0x35, 0x58, 0x82, 0x50, 0xc4, 0x94, 0x44,
	0xb0, 0xe8, 0x74, 0x76, 0x52, 0x8a, 0x52, 0xc8, 0xcc, 0xcb, 0x2c, 0xc9, 0x4c, 0xcc, 0x29, 0xc9,
	0x28, 0x4a, 0x4d, 0x85, 0xb8, 0x25, 0x39, 0x3f, 0x47, 0x3f, 0xb9, 0x18, 0xe6, 0x22, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0xfb, 0x47, 0x39, 0xa5, 0xa4, 0x00, 0x00, 0x00,
}
