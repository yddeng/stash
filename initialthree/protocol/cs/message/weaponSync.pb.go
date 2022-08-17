// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/weaponSync.proto

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

type WeaponSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WeaponSyncToS) Reset()         { *m = WeaponSyncToS{} }
func (m *WeaponSyncToS) String() string { return proto.CompactTextString(m) }
func (*WeaponSyncToS) ProtoMessage()    {}
func (*WeaponSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_95ab88a3ee2deef7, []int{0}
}

func (m *WeaponSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeaponSyncToS.Unmarshal(m, b)
}
func (m *WeaponSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeaponSyncToS.Marshal(b, m, deterministic)
}
func (m *WeaponSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeaponSyncToS.Merge(m, src)
}
func (m *WeaponSyncToS) XXX_Size() int {
	return xxx_messageInfo_WeaponSyncToS.Size(m)
}
func (m *WeaponSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_WeaponSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_WeaponSyncToS proto.InternalMessageInfo

type WeaponSyncToC struct {
	IsAll                *bool     `protobuf:"varint,1,req,name=isAll" json:"isAll,omitempty"`
	Weapons              []*Weapon `protobuf:"bytes,2,rep,name=weapons" json:"weapons,omitempty"`
	UseCap               *int32    `protobuf:"varint,3,req,name=useCap" json:"useCap,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *WeaponSyncToC) Reset()         { *m = WeaponSyncToC{} }
func (m *WeaponSyncToC) String() string { return proto.CompactTextString(m) }
func (*WeaponSyncToC) ProtoMessage()    {}
func (*WeaponSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_95ab88a3ee2deef7, []int{1}
}

func (m *WeaponSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeaponSyncToC.Unmarshal(m, b)
}
func (m *WeaponSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeaponSyncToC.Marshal(b, m, deterministic)
}
func (m *WeaponSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeaponSyncToC.Merge(m, src)
}
func (m *WeaponSyncToC) XXX_Size() int {
	return xxx_messageInfo_WeaponSyncToC.Size(m)
}
func (m *WeaponSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_WeaponSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_WeaponSyncToC proto.InternalMessageInfo

func (m *WeaponSyncToC) GetIsAll() bool {
	if m != nil && m.IsAll != nil {
		return *m.IsAll
	}
	return false
}

func (m *WeaponSyncToC) GetWeapons() []*Weapon {
	if m != nil {
		return m.Weapons
	}
	return nil
}

func (m *WeaponSyncToC) GetUseCap() int32 {
	if m != nil && m.UseCap != nil {
		return *m.UseCap
	}
	return 0
}

type Weapon struct {
	ID                   *uint32  `protobuf:"varint,1,req,name=ID" json:"ID,omitempty"`
	ConfigID             *int32   `protobuf:"varint,2,opt,name=ConfigID" json:"ConfigID,omitempty"`
	Level                *int32   `protobuf:"varint,3,opt,name=level" json:"level,omitempty"`
	Exp                  *int32   `protobuf:"varint,4,opt,name=exp" json:"exp,omitempty"`
	RefineLevel          *int32   `protobuf:"varint,5,opt,name=refineLevel" json:"refineLevel,omitempty"`
	EquipCharacterID     *int32   `protobuf:"varint,6,opt,name=equipCharacterID" json:"equipCharacterID,omitempty"`
	BreakLevel           *int32   `protobuf:"varint,7,opt,name=breakLevel" json:"breakLevel,omitempty"`
	IsLock               *bool    `protobuf:"varint,8,opt,name=isLock" json:"isLock,omitempty"`
	IsRemove             *bool    `protobuf:"varint,9,opt,name=isRemove" json:"isRemove,omitempty"`
	GetTime              *int64   `protobuf:"varint,10,opt,name=getTime" json:"getTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Weapon) Reset()         { *m = Weapon{} }
func (m *Weapon) String() string { return proto.CompactTextString(m) }
func (*Weapon) ProtoMessage()    {}
func (*Weapon) Descriptor() ([]byte, []int) {
	return fileDescriptor_95ab88a3ee2deef7, []int{2}
}

func (m *Weapon) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Weapon.Unmarshal(m, b)
}
func (m *Weapon) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Weapon.Marshal(b, m, deterministic)
}
func (m *Weapon) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Weapon.Merge(m, src)
}
func (m *Weapon) XXX_Size() int {
	return xxx_messageInfo_Weapon.Size(m)
}
func (m *Weapon) XXX_DiscardUnknown() {
	xxx_messageInfo_Weapon.DiscardUnknown(m)
}

var xxx_messageInfo_Weapon proto.InternalMessageInfo

func (m *Weapon) GetID() uint32 {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return 0
}

func (m *Weapon) GetConfigID() int32 {
	if m != nil && m.ConfigID != nil {
		return *m.ConfigID
	}
	return 0
}

func (m *Weapon) GetLevel() int32 {
	if m != nil && m.Level != nil {
		return *m.Level
	}
	return 0
}

func (m *Weapon) GetExp() int32 {
	if m != nil && m.Exp != nil {
		return *m.Exp
	}
	return 0
}

func (m *Weapon) GetRefineLevel() int32 {
	if m != nil && m.RefineLevel != nil {
		return *m.RefineLevel
	}
	return 0
}

func (m *Weapon) GetEquipCharacterID() int32 {
	if m != nil && m.EquipCharacterID != nil {
		return *m.EquipCharacterID
	}
	return 0
}

func (m *Weapon) GetBreakLevel() int32 {
	if m != nil && m.BreakLevel != nil {
		return *m.BreakLevel
	}
	return 0
}

func (m *Weapon) GetIsLock() bool {
	if m != nil && m.IsLock != nil {
		return *m.IsLock
	}
	return false
}

func (m *Weapon) GetIsRemove() bool {
	if m != nil && m.IsRemove != nil {
		return *m.IsRemove
	}
	return false
}

func (m *Weapon) GetGetTime() int64 {
	if m != nil && m.GetTime != nil {
		return *m.GetTime
	}
	return 0
}

func init() {
	proto.RegisterType((*WeaponSyncToS)(nil), "message.weaponSync_toS")
	proto.RegisterType((*WeaponSyncToC)(nil), "message.weaponSync_toC")
	proto.RegisterType((*Weapon)(nil), "message.weapon")
}

func init() { proto.RegisterFile("cs/proto/message/weaponSync.proto", fileDescriptor_95ab88a3ee2deef7) }

var fileDescriptor_95ab88a3ee2deef7 = []byte{
	// 315 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x91, 0x41, 0x4f, 0xb3, 0x40,
	0x10, 0x86, 0xc3, 0xf2, 0x51, 0xf8, 0xa6, 0xb1, 0x36, 0x1b, 0x63, 0x36, 0x1e, 0xcc, 0xca, 0x09,
	0x3d, 0x94, 0xc4, 0x7f, 0xa0, 0x70, 0x21, 0xe9, 0x69, 0xeb, 0xc9, 0x8b, 0x59, 0xc9, 0xb4, 0xdd,
	0x94, 0xb2, 0xc8, 0xd2, 0xaa, 0xff, 0xc4, 0x9f, 0x6b, 0xd8, 0xa5, 0xb5, 0xd1, 0x1b, 0xcf, 0xfb,
	0xbc, 0x19, 0x76, 0x32, 0x70, 0x53, 0x9a, 0xb4, 0x69, 0x75, 0xa7, 0xd3, 0x2d, 0x1a, 0x23, 0x57,
	0x98, 0xbe, 0xa3, 0x6c, 0x74, 0xbd, 0xf8, 0xac, 0xcb, 0x99, 0x15, 0x34, 0x1c, 0x4c, 0x3c, 0x85,
	0xc9, 0x8f, 0x7c, 0xe9, 0xf4, 0x22, 0x56, 0xbf, 0x92, 0x8c, 0x5e, 0x40, 0xa0, 0xcc, 0x43, 0x55,
	0x31, 0x8f, 0x93, 0x24, 0x12, 0x0e, 0xe8, 0x2d, 0x84, 0xae, 0x67, 0x18, 0xe1, 0x7e, 0x32, 0xbe,
	0x3f, 0x9f, 0x0d, 0x43, 0x67, 0x2e, 0x17, 0x07, 0x4f, 0x2f, 0x61, 0xb4, 0x33, 0x98, 0xc9, 0x86,
	0xf9, 0x9c, 0x24, 0x81, 0x18, 0x28, 0xfe, 0x22, 0x30, 0x72, 0x1d, 0x3a, 0x01, 0x52, 0xe4, 0xf6,
	0x07, 0x67, 0x82, 0x14, 0x39, 0xbd, 0x82, 0x28, 0xd3, 0xf5, 0x52, 0xad, 0x8a, 0x9c, 0x11, 0xee,
	0x25, 0x81, 0x38, 0x72, 0xff, 0x9e, 0x0a, 0xf7, 0x58, 0x31, 0xdf, 0x0a, 0x07, 0x74, 0x0a, 0x3e,
	0x7e, 0x34, 0xec, 0x9f, 0xcd, 0xfa, 0x4f, 0xca, 0x61, 0xdc, 0xe2, 0x52, 0xd5, 0x38, 0xb7, 0xed,
	0xc0, 0x9a, 0xd3, 0x88, 0xde, 0xc1, 0x14, 0xdf, 0x76, 0xaa, 0xc9, 0xd6, 0xb2, 0x95, 0x65, 0x87,
	0x6d, 0x91, 0xb3, 0x91, 0xad, 0xfd, 0xc9, 0xe9, 0x35, 0xc0, 0x6b, 0x8b, 0x72, 0xe3, 0x86, 0x85,
	0xb6, 0x75, 0x92, 0xf4, 0x4b, 0x2a, 0x33, 0xd7, 0xe5, 0x86, 0x45, 0xdc, 0x4b, 0x22, 0x31, 0x50,
	0xbf, 0x89, 0x32, 0x02, 0xb7, 0x7a, 0x8f, 0xec, 0xbf, 0x35, 0x47, 0xa6, 0x0c, 0xc2, 0x15, 0x76,
	0x4f, 0x6a, 0x8b, 0x0c, 0xb8, 0x97, 0xf8, 0xe2, 0x80, 0x8f, 0xf1, 0x33, 0x57, 0xb5, 0xea, 0x94,
	0xac, 0xba, 0x75, 0x8b, 0xe8, 0xee, 0x59, 0xea, 0x2a, 0x2d, 0xcd, 0xe1, 0xaa, 0xdf, 0x01, 0x00,
	0x00, 0xff, 0xff, 0xb4, 0x23, 0x57, 0x7b, 0xe8, 0x01, 0x00, 0x00,
}
