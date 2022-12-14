// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/talentLevelUp.proto

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

type TalentLevelUpToS struct {
	InfiniteTalent       *bool    `protobuf:"varint,1,opt,name=infiniteTalent" json:"infiniteTalent,omitempty"`
	TalentID             *int32   `protobuf:"varint,2,opt,name=talentID" json:"talentID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TalentLevelUpToS) Reset()         { *m = TalentLevelUpToS{} }
func (m *TalentLevelUpToS) String() string { return proto.CompactTextString(m) }
func (*TalentLevelUpToS) ProtoMessage()    {}
func (*TalentLevelUpToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_c579bb37d48ff4a9, []int{0}
}

func (m *TalentLevelUpToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TalentLevelUpToS.Unmarshal(m, b)
}
func (m *TalentLevelUpToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TalentLevelUpToS.Marshal(b, m, deterministic)
}
func (m *TalentLevelUpToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TalentLevelUpToS.Merge(m, src)
}
func (m *TalentLevelUpToS) XXX_Size() int {
	return xxx_messageInfo_TalentLevelUpToS.Size(m)
}
func (m *TalentLevelUpToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TalentLevelUpToS.DiscardUnknown(m)
}

var xxx_messageInfo_TalentLevelUpToS proto.InternalMessageInfo

func (m *TalentLevelUpToS) GetInfiniteTalent() bool {
	if m != nil && m.InfiniteTalent != nil {
		return *m.InfiniteTalent
	}
	return false
}

func (m *TalentLevelUpToS) GetTalentID() int32 {
	if m != nil && m.TalentID != nil {
		return *m.TalentID
	}
	return 0
}

type TalentLevelUpToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TalentLevelUpToC) Reset()         { *m = TalentLevelUpToC{} }
func (m *TalentLevelUpToC) String() string { return proto.CompactTextString(m) }
func (*TalentLevelUpToC) ProtoMessage()    {}
func (*TalentLevelUpToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_c579bb37d48ff4a9, []int{1}
}

func (m *TalentLevelUpToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TalentLevelUpToC.Unmarshal(m, b)
}
func (m *TalentLevelUpToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TalentLevelUpToC.Marshal(b, m, deterministic)
}
func (m *TalentLevelUpToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TalentLevelUpToC.Merge(m, src)
}
func (m *TalentLevelUpToC) XXX_Size() int {
	return xxx_messageInfo_TalentLevelUpToC.Size(m)
}
func (m *TalentLevelUpToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TalentLevelUpToC.DiscardUnknown(m)
}

var xxx_messageInfo_TalentLevelUpToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TalentLevelUpToS)(nil), "message.talentLevelUp_toS")
	proto.RegisterType((*TalentLevelUpToC)(nil), "message.talentLevelUp_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/talentLevelUp.proto", fileDescriptor_c579bb37d48ff4a9)
}

var fileDescriptor_c579bb37d48ff4a9 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x49, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0xcc,
	0x49, 0xcd, 0x2b, 0xf1, 0x49, 0x2d, 0x4b, 0xcd, 0x09, 0x2d, 0xd0, 0x03, 0xcb, 0x09, 0xb1, 0x43,
	0x25, 0x95, 0xc2, 0xb9, 0x04, 0x51, 0xe4, 0xe3, 0x4b, 0xf2, 0x83, 0x85, 0xd4, 0xb8, 0xf8, 0x32,
	0xf3, 0xd2, 0x32, 0xf3, 0x32, 0x4b, 0x52, 0x43, 0xc0, 0x92, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x1c,
	0x41, 0x68, 0xa2, 0x42, 0x52, 0x5c, 0x1c, 0x10, 0xcd, 0x9e, 0x2e, 0x12, 0x4c, 0x0a, 0x8c, 0x1a,
	0xac, 0x41, 0x70, 0xbe, 0x92, 0x30, 0xa6, 0xc1, 0xce, 0x4e, 0x4a, 0x51, 0x0a, 0x20, 0xfd, 0x99,
	0x89, 0x39, 0x25, 0x19, 0x45, 0xa9, 0xa9, 0x10, 0x87, 0x26, 0xe7, 0xe7, 0xe8, 0x27, 0x17, 0xc3,
	0x9c, 0x0b, 0x08, 0x00, 0x00, 0xff, 0xff, 0xae, 0x6c, 0xf9, 0xe5, 0xc1, 0x00, 0x00, 0x00,
}
