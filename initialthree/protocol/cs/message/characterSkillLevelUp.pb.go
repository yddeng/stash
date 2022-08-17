// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/characterSkillLevelUp.proto

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

type CharacterSkillLevelUpToS struct {
	CharacterID          *int32   `protobuf:"varint,1,req,name=CharacterID" json:"CharacterID,omitempty"`
	SkillID              *int32   `protobuf:"varint,2,req,name=skillID" json:"skillID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CharacterSkillLevelUpToS) Reset()         { *m = CharacterSkillLevelUpToS{} }
func (m *CharacterSkillLevelUpToS) String() string { return proto.CompactTextString(m) }
func (*CharacterSkillLevelUpToS) ProtoMessage()    {}
func (*CharacterSkillLevelUpToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_25509715e7c4f175, []int{0}
}

func (m *CharacterSkillLevelUpToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CharacterSkillLevelUpToS.Unmarshal(m, b)
}
func (m *CharacterSkillLevelUpToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CharacterSkillLevelUpToS.Marshal(b, m, deterministic)
}
func (m *CharacterSkillLevelUpToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CharacterSkillLevelUpToS.Merge(m, src)
}
func (m *CharacterSkillLevelUpToS) XXX_Size() int {
	return xxx_messageInfo_CharacterSkillLevelUpToS.Size(m)
}
func (m *CharacterSkillLevelUpToS) XXX_DiscardUnknown() {
	xxx_messageInfo_CharacterSkillLevelUpToS.DiscardUnknown(m)
}

var xxx_messageInfo_CharacterSkillLevelUpToS proto.InternalMessageInfo

func (m *CharacterSkillLevelUpToS) GetCharacterID() int32 {
	if m != nil && m.CharacterID != nil {
		return *m.CharacterID
	}
	return 0
}

func (m *CharacterSkillLevelUpToS) GetSkillID() int32 {
	if m != nil && m.SkillID != nil {
		return *m.SkillID
	}
	return 0
}

type CharacterSkillLevelUpToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CharacterSkillLevelUpToC) Reset()         { *m = CharacterSkillLevelUpToC{} }
func (m *CharacterSkillLevelUpToC) String() string { return proto.CompactTextString(m) }
func (*CharacterSkillLevelUpToC) ProtoMessage()    {}
func (*CharacterSkillLevelUpToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_25509715e7c4f175, []int{1}
}

func (m *CharacterSkillLevelUpToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CharacterSkillLevelUpToC.Unmarshal(m, b)
}
func (m *CharacterSkillLevelUpToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CharacterSkillLevelUpToC.Marshal(b, m, deterministic)
}
func (m *CharacterSkillLevelUpToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CharacterSkillLevelUpToC.Merge(m, src)
}
func (m *CharacterSkillLevelUpToC) XXX_Size() int {
	return xxx_messageInfo_CharacterSkillLevelUpToC.Size(m)
}
func (m *CharacterSkillLevelUpToC) XXX_DiscardUnknown() {
	xxx_messageInfo_CharacterSkillLevelUpToC.DiscardUnknown(m)
}

var xxx_messageInfo_CharacterSkillLevelUpToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*CharacterSkillLevelUpToS)(nil), "message.characterSkillLevelUp_toS")
	proto.RegisterType((*CharacterSkillLevelUpToC)(nil), "message.characterSkillLevelUp_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/characterSkillLevelUp.proto", fileDescriptor_25509715e7c4f175)
}

var fileDescriptor_25509715e7c4f175 = []byte{
	// 146 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x49, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0xce, 0x48,
	0x2c, 0x4a, 0x4c, 0x2e, 0x49, 0x2d, 0x0a, 0xce, 0xce, 0xcc, 0xc9, 0xf1, 0x49, 0x2d, 0x4b, 0xcd,
	0x09, 0x2d, 0xd0, 0x03, 0xab, 0x11, 0x62, 0x87, 0x2a, 0x52, 0x0a, 0xe7, 0x92, 0xc4, 0xaa, 0x2e,
	0xbe, 0x24, 0x3f, 0x58, 0x48, 0x81, 0x8b, 0xdb, 0x19, 0x26, 0xe9, 0xe9, 0x22, 0xc1, 0xa8, 0xc0,
	0xa4, 0xc1, 0x1a, 0x84, 0x2c, 0x24, 0x24, 0xc1, 0xc5, 0x5e, 0x0c, 0xd2, 0xe5, 0xe9, 0x22, 0xc1,
	0x04, 0x96, 0x85, 0x71, 0x95, 0xa4, 0x71, 0x1b, 0xec, 0xec, 0xa4, 0x14, 0xa5, 0x90, 0x99, 0x97,
	0x59, 0x92, 0x99, 0x98, 0x53, 0x92, 0x51, 0x94, 0x9a, 0x0a, 0x71, 0x78, 0x72, 0x7e, 0x8e, 0x7e,
	0x72, 0x31, 0xcc, 0xf9, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x8a, 0x15, 0x67, 0x08, 0xd1, 0x00,
	0x00, 0x00,
}