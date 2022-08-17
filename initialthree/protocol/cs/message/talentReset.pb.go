// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/talentReset.proto

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

type TalentResetToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TalentResetToS) Reset()         { *m = TalentResetToS{} }
func (m *TalentResetToS) String() string { return proto.CompactTextString(m) }
func (*TalentResetToS) ProtoMessage()    {}
func (*TalentResetToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_4056a2728ceabc0a, []int{0}
}

func (m *TalentResetToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TalentResetToS.Unmarshal(m, b)
}
func (m *TalentResetToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TalentResetToS.Marshal(b, m, deterministic)
}
func (m *TalentResetToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TalentResetToS.Merge(m, src)
}
func (m *TalentResetToS) XXX_Size() int {
	return xxx_messageInfo_TalentResetToS.Size(m)
}
func (m *TalentResetToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TalentResetToS.DiscardUnknown(m)
}

var xxx_messageInfo_TalentResetToS proto.InternalMessageInfo

type TalentResetToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TalentResetToC) Reset()         { *m = TalentResetToC{} }
func (m *TalentResetToC) String() string { return proto.CompactTextString(m) }
func (*TalentResetToC) ProtoMessage()    {}
func (*TalentResetToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_4056a2728ceabc0a, []int{1}
}

func (m *TalentResetToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TalentResetToC.Unmarshal(m, b)
}
func (m *TalentResetToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TalentResetToC.Marshal(b, m, deterministic)
}
func (m *TalentResetToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TalentResetToC.Merge(m, src)
}
func (m *TalentResetToC) XXX_Size() int {
	return xxx_messageInfo_TalentResetToC.Size(m)
}
func (m *TalentResetToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TalentResetToC.DiscardUnknown(m)
}

var xxx_messageInfo_TalentResetToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TalentResetToS)(nil), "message.talentReset_toS")
	proto.RegisterType((*TalentResetToC)(nil), "message.talentReset_toC")
}

func init() { proto.RegisterFile("cs/proto/message/talentReset.proto", fileDescriptor_4056a2728ceabc0a) }

var fileDescriptor_4056a2728ceabc0a = []byte{
	// 100 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4a, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0xcc,
	0x49, 0xcd, 0x2b, 0x09, 0x4a, 0x2d, 0x4e, 0x2d, 0xd1, 0x03, 0xcb, 0x08, 0xb1, 0x43, 0xa5, 0x94,
	0x04, 0xb9, 0xf8, 0x91, 0x64, 0xe3, 0x4b, 0xf2, 0x83, 0x31, 0x85, 0x9c, 0x9d, 0x94, 0xa2, 0x14,
	0x32, 0xf3, 0x32, 0x4b, 0x32, 0x13, 0x73, 0x4a, 0x32, 0x8a, 0x52, 0x53, 0x21, 0xc6, 0x27, 0xe7,
	0xe7, 0xe8, 0x27, 0x17, 0xc3, 0x2c, 0x01, 0x04, 0x00, 0x00, 0xff, 0xff, 0x6c, 0x8e, 0x47, 0xd9,
	0x77, 0x00, 0x00, 0x00,
}
