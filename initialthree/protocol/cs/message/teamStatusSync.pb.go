// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/teamStatusSync.proto

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

type TeamStatusSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamStatusSyncToS) Reset()         { *m = TeamStatusSyncToS{} }
func (m *TeamStatusSyncToS) String() string { return proto.CompactTextString(m) }
func (*TeamStatusSyncToS) ProtoMessage()    {}
func (*TeamStatusSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_e19b1156c52adaba, []int{0}
}

func (m *TeamStatusSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamStatusSyncToS.Unmarshal(m, b)
}
func (m *TeamStatusSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamStatusSyncToS.Marshal(b, m, deterministic)
}
func (m *TeamStatusSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamStatusSyncToS.Merge(m, src)
}
func (m *TeamStatusSyncToS) XXX_Size() int {
	return xxx_messageInfo_TeamStatusSyncToS.Size(m)
}
func (m *TeamStatusSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamStatusSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_TeamStatusSyncToS proto.InternalMessageInfo

type TeamStatusSyncToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamStatusSyncToC) Reset()         { *m = TeamStatusSyncToC{} }
func (m *TeamStatusSyncToC) String() string { return proto.CompactTextString(m) }
func (*TeamStatusSyncToC) ProtoMessage()    {}
func (*TeamStatusSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_e19b1156c52adaba, []int{1}
}

func (m *TeamStatusSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamStatusSyncToC.Unmarshal(m, b)
}
func (m *TeamStatusSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamStatusSyncToC.Marshal(b, m, deterministic)
}
func (m *TeamStatusSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamStatusSyncToC.Merge(m, src)
}
func (m *TeamStatusSyncToC) XXX_Size() int {
	return xxx_messageInfo_TeamStatusSyncToC.Size(m)
}
func (m *TeamStatusSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamStatusSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_TeamStatusSyncToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TeamStatusSyncToS)(nil), "message.teamStatusSync_toS")
	proto.RegisterType((*TeamStatusSyncToC)(nil), "message.teamStatusSync_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/teamStatusSync.proto", fileDescriptor_e19b1156c52adaba)
}

var fileDescriptor_e19b1156c52adaba = []byte{
	// 103 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4d, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0x4d,
	0xcc, 0x0d, 0x2e, 0x49, 0x2c, 0x29, 0x2d, 0x0e, 0xae, 0xcc, 0x4b, 0xd6, 0x03, 0x4b, 0x0a, 0xb1,
	0x43, 0x65, 0x95, 0x44, 0xb8, 0x84, 0x50, 0x15, 0xc4, 0x97, 0xe4, 0x07, 0x63, 0x15, 0x75, 0x76,
	0x52, 0x8a, 0x52, 0xc8, 0xcc, 0xcb, 0x2c, 0xc9, 0x4c, 0xcc, 0x29, 0xc9, 0x28, 0x4a, 0x4d, 0x85,
	0xd8, 0x93, 0x9c, 0x9f, 0xa3, 0x9f, 0x5c, 0x0c, 0xb3, 0x0d, 0x10, 0x00, 0x00, 0xff, 0xff, 0x62,
	0x68, 0x7a, 0x58, 0x80, 0x00, 0x00, 0x00,
}