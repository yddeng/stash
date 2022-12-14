// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/characterTeamPrefabSet.proto

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

type CharacterTeamPrefabSetToS struct {
	TeamPrefabIdx        *int32               `protobuf:"varint,1,req,name=teamPrefabIdx" json:"teamPrefabIdx,omitempty"`
	Prefab               *CharacterTeamPrefab `protobuf:"bytes,2,req,name=prefab" json:"prefab,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CharacterTeamPrefabSetToS) Reset()         { *m = CharacterTeamPrefabSetToS{} }
func (m *CharacterTeamPrefabSetToS) String() string { return proto.CompactTextString(m) }
func (*CharacterTeamPrefabSetToS) ProtoMessage()    {}
func (*CharacterTeamPrefabSetToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_d73331a9391df712, []int{0}
}

func (m *CharacterTeamPrefabSetToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CharacterTeamPrefabSetToS.Unmarshal(m, b)
}
func (m *CharacterTeamPrefabSetToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CharacterTeamPrefabSetToS.Marshal(b, m, deterministic)
}
func (m *CharacterTeamPrefabSetToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CharacterTeamPrefabSetToS.Merge(m, src)
}
func (m *CharacterTeamPrefabSetToS) XXX_Size() int {
	return xxx_messageInfo_CharacterTeamPrefabSetToS.Size(m)
}
func (m *CharacterTeamPrefabSetToS) XXX_DiscardUnknown() {
	xxx_messageInfo_CharacterTeamPrefabSetToS.DiscardUnknown(m)
}

var xxx_messageInfo_CharacterTeamPrefabSetToS proto.InternalMessageInfo

func (m *CharacterTeamPrefabSetToS) GetTeamPrefabIdx() int32 {
	if m != nil && m.TeamPrefabIdx != nil {
		return *m.TeamPrefabIdx
	}
	return 0
}

func (m *CharacterTeamPrefabSetToS) GetPrefab() *CharacterTeamPrefab {
	if m != nil {
		return m.Prefab
	}
	return nil
}

type CharacterTeamPrefabSetToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CharacterTeamPrefabSetToC) Reset()         { *m = CharacterTeamPrefabSetToC{} }
func (m *CharacterTeamPrefabSetToC) String() string { return proto.CompactTextString(m) }
func (*CharacterTeamPrefabSetToC) ProtoMessage()    {}
func (*CharacterTeamPrefabSetToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_d73331a9391df712, []int{1}
}

func (m *CharacterTeamPrefabSetToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CharacterTeamPrefabSetToC.Unmarshal(m, b)
}
func (m *CharacterTeamPrefabSetToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CharacterTeamPrefabSetToC.Marshal(b, m, deterministic)
}
func (m *CharacterTeamPrefabSetToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CharacterTeamPrefabSetToC.Merge(m, src)
}
func (m *CharacterTeamPrefabSetToC) XXX_Size() int {
	return xxx_messageInfo_CharacterTeamPrefabSetToC.Size(m)
}
func (m *CharacterTeamPrefabSetToC) XXX_DiscardUnknown() {
	xxx_messageInfo_CharacterTeamPrefabSetToC.DiscardUnknown(m)
}

var xxx_messageInfo_CharacterTeamPrefabSetToC proto.InternalMessageInfo

type CharacterTeamPrefab struct {
	Name                 *string  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	CharacterList        []int32  `protobuf:"varint,2,rep,name=characterList" json:"characterList,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CharacterTeamPrefab) Reset()         { *m = CharacterTeamPrefab{} }
func (m *CharacterTeamPrefab) String() string { return proto.CompactTextString(m) }
func (*CharacterTeamPrefab) ProtoMessage()    {}
func (*CharacterTeamPrefab) Descriptor() ([]byte, []int) {
	return fileDescriptor_d73331a9391df712, []int{2}
}

func (m *CharacterTeamPrefab) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CharacterTeamPrefab.Unmarshal(m, b)
}
func (m *CharacterTeamPrefab) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CharacterTeamPrefab.Marshal(b, m, deterministic)
}
func (m *CharacterTeamPrefab) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CharacterTeamPrefab.Merge(m, src)
}
func (m *CharacterTeamPrefab) XXX_Size() int {
	return xxx_messageInfo_CharacterTeamPrefab.Size(m)
}
func (m *CharacterTeamPrefab) XXX_DiscardUnknown() {
	xxx_messageInfo_CharacterTeamPrefab.DiscardUnknown(m)
}

var xxx_messageInfo_CharacterTeamPrefab proto.InternalMessageInfo

func (m *CharacterTeamPrefab) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *CharacterTeamPrefab) GetCharacterList() []int32 {
	if m != nil {
		return m.CharacterList
	}
	return nil
}

func init() {
	proto.RegisterType((*CharacterTeamPrefabSetToS)(nil), "message.characterTeamPrefabSet_toS")
	proto.RegisterType((*CharacterTeamPrefabSetToC)(nil), "message.characterTeamPrefabSet_toC")
	proto.RegisterType((*CharacterTeamPrefab)(nil), "message.characterTeamPrefab")
}

func init() {
	proto.RegisterFile("cs/proto/message/characterTeamPrefabSet.proto", fileDescriptor_d73331a9391df712)
}

var fileDescriptor_d73331a9391df712 = []byte{
	// 197 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4d, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0xce, 0x48,
	0x2c, 0x4a, 0x4c, 0x2e, 0x49, 0x2d, 0x0a, 0x49, 0x4d, 0xcc, 0x0d, 0x28, 0x4a, 0x4d, 0x4b, 0x4c,
	0x0a, 0x4e, 0x2d, 0xd1, 0x03, 0x2b, 0x12, 0x62, 0x87, 0xaa, 0x52, 0xaa, 0xe0, 0x92, 0xc2, 0xae,
	0x30, 0xbe, 0x24, 0x3f, 0x58, 0x48, 0x85, 0x8b, 0xb7, 0x04, 0x2e, 0xe8, 0x99, 0x52, 0x21, 0xc1,
	0xa8, 0xc0, 0xa4, 0xc1, 0x1a, 0x84, 0x2a, 0x28, 0x64, 0xc2, 0xc5, 0x56, 0x00, 0xe6, 0x48, 0x30,
	0x29, 0x30, 0x69, 0x70, 0x1b, 0xc9, 0xe8, 0x41, 0x4d, 0xd7, 0xc3, 0x62, 0x74, 0x10, 0x54, 0xad,
	0x92, 0x0c, 0x1e, 0x9b, 0x9d, 0x95, 0xfc, 0xb9, 0x84, 0xb1, 0xc8, 0x0a, 0x09, 0x71, 0xb1, 0xe4,
	0x25, 0xe6, 0xa6, 0x4a, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06, 0x81, 0xd9, 0x20, 0x47, 0xc2, 0x95,
	0xfa, 0x64, 0x16, 0x97, 0x48, 0x30, 0x29, 0x30, 0x83, 0x1c, 0x89, 0x22, 0xe8, 0xa4, 0x14, 0xa5,
	0x90, 0x99, 0x97, 0x59, 0x92, 0x99, 0x98, 0x53, 0x92, 0x51, 0x94, 0x9a, 0x0a, 0x09, 0xac, 0xe4,
	0xfc, 0x1c, 0xfd, 0xe4, 0x62, 0x58, 0x90, 0x01, 0x02, 0x00, 0x00, 0xff, 0xff, 0x56, 0x83, 0x81,
	0x1b, 0x45, 0x01, 0x00, 0x00,
}
