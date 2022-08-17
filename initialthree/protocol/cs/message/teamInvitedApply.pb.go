// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/teamInvitedApply.proto

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

type TeamInvitedApplyToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamInvitedApplyToS) Reset()         { *m = TeamInvitedApplyToS{} }
func (m *TeamInvitedApplyToS) String() string { return proto.CompactTextString(m) }
func (*TeamInvitedApplyToS) ProtoMessage()    {}
func (*TeamInvitedApplyToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_950a75a60f6cdfe5, []int{0}
}

func (m *TeamInvitedApplyToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamInvitedApplyToS.Unmarshal(m, b)
}
func (m *TeamInvitedApplyToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamInvitedApplyToS.Marshal(b, m, deterministic)
}
func (m *TeamInvitedApplyToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamInvitedApplyToS.Merge(m, src)
}
func (m *TeamInvitedApplyToS) XXX_Size() int {
	return xxx_messageInfo_TeamInvitedApplyToS.Size(m)
}
func (m *TeamInvitedApplyToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamInvitedApplyToS.DiscardUnknown(m)
}

var xxx_messageInfo_TeamInvitedApplyToS proto.InternalMessageInfo

type TeamInvitedApplyToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamInvitedApplyToC) Reset()         { *m = TeamInvitedApplyToC{} }
func (m *TeamInvitedApplyToC) String() string { return proto.CompactTextString(m) }
func (*TeamInvitedApplyToC) ProtoMessage()    {}
func (*TeamInvitedApplyToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_950a75a60f6cdfe5, []int{1}
}

func (m *TeamInvitedApplyToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamInvitedApplyToC.Unmarshal(m, b)
}
func (m *TeamInvitedApplyToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamInvitedApplyToC.Marshal(b, m, deterministic)
}
func (m *TeamInvitedApplyToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamInvitedApplyToC.Merge(m, src)
}
func (m *TeamInvitedApplyToC) XXX_Size() int {
	return xxx_messageInfo_TeamInvitedApplyToC.Size(m)
}
func (m *TeamInvitedApplyToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamInvitedApplyToC.DiscardUnknown(m)
}

var xxx_messageInfo_TeamInvitedApplyToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TeamInvitedApplyToS)(nil), "message.teamInvitedApply_toS")
	proto.RegisterType((*TeamInvitedApplyToC)(nil), "message.teamInvitedApply_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/teamInvitedApply.proto", fileDescriptor_950a75a60f6cdfe5)
}

var fileDescriptor_950a75a60f6cdfe5 = []byte{
	// 105 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4f, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0x4d,
	0xcc, 0xf5, 0xcc, 0x2b, 0xcb, 0x2c, 0x49, 0x4d, 0x71, 0x2c, 0x28, 0xc8, 0xa9, 0xd4, 0x03, 0x4b,
	0x0b, 0xb1, 0x43, 0xe5, 0x95, 0xc4, 0xb8, 0x44, 0xd0, 0x95, 0xc4, 0x97, 0xe4, 0x07, 0xe3, 0x10,
	0x77, 0x76, 0x52, 0x8a, 0x52, 0xc8, 0xcc, 0xcb, 0x2c, 0xc9, 0x4c, 0xcc, 0x29, 0xc9, 0x28, 0x4a,
	0x4d, 0x85, 0xd8, 0x96, 0x9c, 0x9f, 0xa3, 0x9f, 0x5c, 0x0c, 0xb3, 0x13, 0x10, 0x00, 0x00, 0xff,
	0xff, 0xc1, 0xd8, 0x92, 0xa2, 0x86, 0x00, 0x00, 0x00,
}