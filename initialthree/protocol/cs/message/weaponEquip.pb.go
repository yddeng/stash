// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/weaponEquip.proto

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

type WeaponEquipToS struct {
	CharacterID          *int32   `protobuf:"varint,1,req,name=characterID" json:"characterID,omitempty"`
	WeaponID             *uint32  `protobuf:"varint,2,req,name=weaponID" json:"weaponID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WeaponEquipToS) Reset()         { *m = WeaponEquipToS{} }
func (m *WeaponEquipToS) String() string { return proto.CompactTextString(m) }
func (*WeaponEquipToS) ProtoMessage()    {}
func (*WeaponEquipToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_a97dca4dbd77ba12, []int{0}
}

func (m *WeaponEquipToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeaponEquipToS.Unmarshal(m, b)
}
func (m *WeaponEquipToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeaponEquipToS.Marshal(b, m, deterministic)
}
func (m *WeaponEquipToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeaponEquipToS.Merge(m, src)
}
func (m *WeaponEquipToS) XXX_Size() int {
	return xxx_messageInfo_WeaponEquipToS.Size(m)
}
func (m *WeaponEquipToS) XXX_DiscardUnknown() {
	xxx_messageInfo_WeaponEquipToS.DiscardUnknown(m)
}

var xxx_messageInfo_WeaponEquipToS proto.InternalMessageInfo

func (m *WeaponEquipToS) GetCharacterID() int32 {
	if m != nil && m.CharacterID != nil {
		return *m.CharacterID
	}
	return 0
}

func (m *WeaponEquipToS) GetWeaponID() uint32 {
	if m != nil && m.WeaponID != nil {
		return *m.WeaponID
	}
	return 0
}

type WeaponEquipToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WeaponEquipToC) Reset()         { *m = WeaponEquipToC{} }
func (m *WeaponEquipToC) String() string { return proto.CompactTextString(m) }
func (*WeaponEquipToC) ProtoMessage()    {}
func (*WeaponEquipToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_a97dca4dbd77ba12, []int{1}
}

func (m *WeaponEquipToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeaponEquipToC.Unmarshal(m, b)
}
func (m *WeaponEquipToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeaponEquipToC.Marshal(b, m, deterministic)
}
func (m *WeaponEquipToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeaponEquipToC.Merge(m, src)
}
func (m *WeaponEquipToC) XXX_Size() int {
	return xxx_messageInfo_WeaponEquipToC.Size(m)
}
func (m *WeaponEquipToC) XXX_DiscardUnknown() {
	xxx_messageInfo_WeaponEquipToC.DiscardUnknown(m)
}

var xxx_messageInfo_WeaponEquipToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*WeaponEquipToS)(nil), "message.weaponEquip_toS")
	proto.RegisterType((*WeaponEquipToC)(nil), "message.weaponEquip_toC")
}

func init() { proto.RegisterFile("cs/proto/message/weaponEquip.proto", fileDescriptor_a97dca4dbd77ba12) }

var fileDescriptor_a97dca4dbd77ba12 = []byte{
	// 143 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4a, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x4f, 0x4d,
	0x2c, 0xc8, 0xcf, 0x73, 0x2d, 0x2c, 0xcd, 0x2c, 0xd0, 0x03, 0xcb, 0x08, 0xb1, 0x43, 0xa5, 0x94,
	0xfc, 0xb9, 0xf8, 0x91, 0x64, 0xe3, 0x4b, 0xf2, 0x83, 0x85, 0x14, 0xb8, 0xb8, 0x93, 0x33, 0x12,
	0x8b, 0x12, 0x93, 0x4b, 0x52, 0x8b, 0x3c, 0x5d, 0x24, 0x18, 0x15, 0x98, 0x34, 0x58, 0x83, 0x90,
	0x85, 0x84, 0xa4, 0xb8, 0x38, 0x20, 0x9a, 0x3c, 0x5d, 0x24, 0x98, 0x14, 0x98, 0x34, 0x78, 0x83,
	0xe0, 0x7c, 0x25, 0x41, 0x74, 0x03, 0x9d, 0x9d, 0x94, 0xa2, 0x14, 0x32, 0xf3, 0x32, 0x4b, 0x32,
	0x13, 0x73, 0x4a, 0x32, 0x8a, 0x52, 0x53, 0x21, 0x8e, 0x4b, 0xce, 0xcf, 0xd1, 0x4f, 0x2e, 0x86,
	0x39, 0x11, 0x10, 0x00, 0x00, 0xff, 0xff, 0x6c, 0xaa, 0x2f, 0x05, 0xb5, 0x00, 0x00, 0x00,
}
