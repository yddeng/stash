// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/teamDismissNotify.proto

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

type TeamDismissNotifyToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamDismissNotifyToS) Reset()         { *m = TeamDismissNotifyToS{} }
func (m *TeamDismissNotifyToS) String() string { return proto.CompactTextString(m) }
func (*TeamDismissNotifyToS) ProtoMessage()    {}
func (*TeamDismissNotifyToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_d118965bbe176755, []int{0}
}

func (m *TeamDismissNotifyToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamDismissNotifyToS.Unmarshal(m, b)
}
func (m *TeamDismissNotifyToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamDismissNotifyToS.Marshal(b, m, deterministic)
}
func (m *TeamDismissNotifyToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamDismissNotifyToS.Merge(m, src)
}
func (m *TeamDismissNotifyToS) XXX_Size() int {
	return xxx_messageInfo_TeamDismissNotifyToS.Size(m)
}
func (m *TeamDismissNotifyToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamDismissNotifyToS.DiscardUnknown(m)
}

var xxx_messageInfo_TeamDismissNotifyToS proto.InternalMessageInfo

type TeamDismissNotifyToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamDismissNotifyToC) Reset()         { *m = TeamDismissNotifyToC{} }
func (m *TeamDismissNotifyToC) String() string { return proto.CompactTextString(m) }
func (*TeamDismissNotifyToC) ProtoMessage()    {}
func (*TeamDismissNotifyToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_d118965bbe176755, []int{1}
}

func (m *TeamDismissNotifyToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamDismissNotifyToC.Unmarshal(m, b)
}
func (m *TeamDismissNotifyToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamDismissNotifyToC.Marshal(b, m, deterministic)
}
func (m *TeamDismissNotifyToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamDismissNotifyToC.Merge(m, src)
}
func (m *TeamDismissNotifyToC) XXX_Size() int {
	return xxx_messageInfo_TeamDismissNotifyToC.Size(m)
}
func (m *TeamDismissNotifyToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamDismissNotifyToC.DiscardUnknown(m)
}

var xxx_messageInfo_TeamDismissNotifyToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TeamDismissNotifyToS)(nil), "message.teamDismissNotify_toS")
	proto.RegisterType((*TeamDismissNotifyToC)(nil), "message.teamDismissNotify_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/teamDismissNotify.proto", fileDescriptor_d118965bbe176755)
}

var fileDescriptor_d118965bbe176755 = []byte{
	// 106 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x48, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0x4d,
	0xcc, 0x75, 0xc9, 0x2c, 0xce, 0xcd, 0x2c, 0x2e, 0xf6, 0xcb, 0x2f, 0xc9, 0x4c, 0xab, 0xd4, 0x03,
	0xcb, 0x0b, 0xb1, 0x43, 0x15, 0x28, 0x89, 0x73, 0x89, 0x62, 0xa8, 0x89, 0x2f, 0xc9, 0x0f, 0xc6,
	0x25, 0xe1, 0xec, 0xa4, 0x14, 0xa5, 0x90, 0x99, 0x97, 0x59, 0x92, 0x99, 0x98, 0x53, 0x92, 0x51,
	0x94, 0x9a, 0x0a, 0xb1, 0x30, 0x39, 0x3f, 0x47, 0x3f, 0xb9, 0x18, 0x66, 0x2d, 0x20, 0x00, 0x00,
	0xff, 0xff, 0x00, 0x94, 0x47, 0xb6, 0x89, 0x00, 0x00, 0x00,
}
