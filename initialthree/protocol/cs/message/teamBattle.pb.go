// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/teamBattle.proto

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

type TeamBattleToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamBattleToS) Reset()         { *m = TeamBattleToS{} }
func (m *TeamBattleToS) String() string { return proto.CompactTextString(m) }
func (*TeamBattleToS) ProtoMessage()    {}
func (*TeamBattleToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_372a23e80ad9574a, []int{0}
}

func (m *TeamBattleToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamBattleToS.Unmarshal(m, b)
}
func (m *TeamBattleToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamBattleToS.Marshal(b, m, deterministic)
}
func (m *TeamBattleToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamBattleToS.Merge(m, src)
}
func (m *TeamBattleToS) XXX_Size() int {
	return xxx_messageInfo_TeamBattleToS.Size(m)
}
func (m *TeamBattleToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamBattleToS.DiscardUnknown(m)
}

var xxx_messageInfo_TeamBattleToS proto.InternalMessageInfo

type TeamBattleToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamBattleToC) Reset()         { *m = TeamBattleToC{} }
func (m *TeamBattleToC) String() string { return proto.CompactTextString(m) }
func (*TeamBattleToC) ProtoMessage()    {}
func (*TeamBattleToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_372a23e80ad9574a, []int{1}
}

func (m *TeamBattleToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamBattleToC.Unmarshal(m, b)
}
func (m *TeamBattleToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamBattleToC.Marshal(b, m, deterministic)
}
func (m *TeamBattleToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamBattleToC.Merge(m, src)
}
func (m *TeamBattleToC) XXX_Size() int {
	return xxx_messageInfo_TeamBattleToC.Size(m)
}
func (m *TeamBattleToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamBattleToC.DiscardUnknown(m)
}

var xxx_messageInfo_TeamBattleToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TeamBattleToS)(nil), "message.teamBattle_toS")
	proto.RegisterType((*TeamBattleToC)(nil), "message.teamBattle_toC")
}

func init() { proto.RegisterFile("cs/proto/message/teamBattle.proto", fileDescriptor_372a23e80ad9574a) }

var fileDescriptor_372a23e80ad9574a = []byte{
	// 99 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4c, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0x4d,
	0xcc, 0x75, 0x4a, 0x2c, 0x29, 0xc9, 0x49, 0xd5, 0x03, 0x4b, 0x08, 0xb1, 0x43, 0x65, 0x94, 0x04,
	0xb8, 0xf8, 0x10, 0x92, 0xf1, 0x25, 0xf9, 0xc1, 0x18, 0x22, 0xce, 0x4e, 0x4a, 0x51, 0x0a, 0x99,
	0x79, 0x99, 0x25, 0x99, 0x89, 0x39, 0x25, 0x19, 0x45, 0xa9, 0xa9, 0x10, 0xb3, 0x93, 0xf3, 0x73,
	0xf4, 0x93, 0x8b, 0x61, 0x36, 0x00, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb8, 0xa2, 0x6c, 0x34, 0x74,
	0x00, 0x00, 0x00,
}
