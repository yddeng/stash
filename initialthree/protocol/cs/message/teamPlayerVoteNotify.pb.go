// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/teamPlayerVoteNotify.proto

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

type TeamPlayerVoteNotifyToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamPlayerVoteNotifyToS) Reset()         { *m = TeamPlayerVoteNotifyToS{} }
func (m *TeamPlayerVoteNotifyToS) String() string { return proto.CompactTextString(m) }
func (*TeamPlayerVoteNotifyToS) ProtoMessage()    {}
func (*TeamPlayerVoteNotifyToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ed9c33e758f7ce8, []int{0}
}

func (m *TeamPlayerVoteNotifyToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamPlayerVoteNotifyToS.Unmarshal(m, b)
}
func (m *TeamPlayerVoteNotifyToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamPlayerVoteNotifyToS.Marshal(b, m, deterministic)
}
func (m *TeamPlayerVoteNotifyToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamPlayerVoteNotifyToS.Merge(m, src)
}
func (m *TeamPlayerVoteNotifyToS) XXX_Size() int {
	return xxx_messageInfo_TeamPlayerVoteNotifyToS.Size(m)
}
func (m *TeamPlayerVoteNotifyToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamPlayerVoteNotifyToS.DiscardUnknown(m)
}

var xxx_messageInfo_TeamPlayerVoteNotifyToS proto.InternalMessageInfo

type TeamPlayerVoteNotifyToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamPlayerVoteNotifyToC) Reset()         { *m = TeamPlayerVoteNotifyToC{} }
func (m *TeamPlayerVoteNotifyToC) String() string { return proto.CompactTextString(m) }
func (*TeamPlayerVoteNotifyToC) ProtoMessage()    {}
func (*TeamPlayerVoteNotifyToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_3ed9c33e758f7ce8, []int{1}
}

func (m *TeamPlayerVoteNotifyToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamPlayerVoteNotifyToC.Unmarshal(m, b)
}
func (m *TeamPlayerVoteNotifyToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamPlayerVoteNotifyToC.Marshal(b, m, deterministic)
}
func (m *TeamPlayerVoteNotifyToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamPlayerVoteNotifyToC.Merge(m, src)
}
func (m *TeamPlayerVoteNotifyToC) XXX_Size() int {
	return xxx_messageInfo_TeamPlayerVoteNotifyToC.Size(m)
}
func (m *TeamPlayerVoteNotifyToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamPlayerVoteNotifyToC.DiscardUnknown(m)
}

var xxx_messageInfo_TeamPlayerVoteNotifyToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TeamPlayerVoteNotifyToS)(nil), "message.teamPlayerVoteNotify_toS")
	proto.RegisterType((*TeamPlayerVoteNotifyToC)(nil), "message.teamPlayerVoteNotify_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/teamPlayerVoteNotify.proto", fileDescriptor_3ed9c33e758f7ce8)
}

var fileDescriptor_3ed9c33e758f7ce8 = []byte{
	// 110 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4e, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0x4d,
	0xcc, 0x0d, 0xc8, 0x49, 0xac, 0x4c, 0x2d, 0x0a, 0xcb, 0x2f, 0x49, 0xf5, 0xcb, 0x2f, 0xc9, 0x4c,
	0xab, 0xd4, 0x03, 0x2b, 0x11, 0x62, 0x87, 0xaa, 0x51, 0x92, 0xe2, 0x92, 0xc0, 0xa6, 0x2c, 0xbe,
	0x24, 0x3f, 0x18, 0x8f, 0x9c, 0xb3, 0x93, 0x52, 0x94, 0x42, 0x66, 0x5e, 0x66, 0x49, 0x66, 0x62,
	0x4e, 0x49, 0x46, 0x51, 0x6a, 0x2a, 0xc4, 0xe6, 0xe4, 0xfc, 0x1c, 0xfd, 0xe4, 0x62, 0x98, 0xfd,
	0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfb, 0x82, 0xc3, 0xa2, 0x92, 0x00, 0x00, 0x00,
}
