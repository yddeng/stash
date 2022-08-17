// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/weaponRefine.proto

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

type WeaponRefineToS struct {
	WeaponID             *uint32  `protobuf:"varint,1,req,name=weaponID" json:"weaponID,omitempty"`
	CostWeapons          *uint32  `protobuf:"varint,2,req,name=costWeapons" json:"costWeapons,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WeaponRefineToS) Reset()         { *m = WeaponRefineToS{} }
func (m *WeaponRefineToS) String() string { return proto.CompactTextString(m) }
func (*WeaponRefineToS) ProtoMessage()    {}
func (*WeaponRefineToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_83ff183410a3b868, []int{0}
}

func (m *WeaponRefineToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeaponRefineToS.Unmarshal(m, b)
}
func (m *WeaponRefineToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeaponRefineToS.Marshal(b, m, deterministic)
}
func (m *WeaponRefineToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeaponRefineToS.Merge(m, src)
}
func (m *WeaponRefineToS) XXX_Size() int {
	return xxx_messageInfo_WeaponRefineToS.Size(m)
}
func (m *WeaponRefineToS) XXX_DiscardUnknown() {
	xxx_messageInfo_WeaponRefineToS.DiscardUnknown(m)
}

var xxx_messageInfo_WeaponRefineToS proto.InternalMessageInfo

func (m *WeaponRefineToS) GetWeaponID() uint32 {
	if m != nil && m.WeaponID != nil {
		return *m.WeaponID
	}
	return 0
}

func (m *WeaponRefineToS) GetCostWeapons() uint32 {
	if m != nil && m.CostWeapons != nil {
		return *m.CostWeapons
	}
	return 0
}

type WeaponRefineToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *WeaponRefineToC) Reset()         { *m = WeaponRefineToC{} }
func (m *WeaponRefineToC) String() string { return proto.CompactTextString(m) }
func (*WeaponRefineToC) ProtoMessage()    {}
func (*WeaponRefineToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_83ff183410a3b868, []int{1}
}

func (m *WeaponRefineToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WeaponRefineToC.Unmarshal(m, b)
}
func (m *WeaponRefineToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WeaponRefineToC.Marshal(b, m, deterministic)
}
func (m *WeaponRefineToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WeaponRefineToC.Merge(m, src)
}
func (m *WeaponRefineToC) XXX_Size() int {
	return xxx_messageInfo_WeaponRefineToC.Size(m)
}
func (m *WeaponRefineToC) XXX_DiscardUnknown() {
	xxx_messageInfo_WeaponRefineToC.DiscardUnknown(m)
}

var xxx_messageInfo_WeaponRefineToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*WeaponRefineToS)(nil), "message.weaponRefine_toS")
	proto.RegisterType((*WeaponRefineToC)(nil), "message.weaponRefine_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/weaponRefine.proto", fileDescriptor_83ff183410a3b868)
}

var fileDescriptor_83ff183410a3b868 = []byte{
	// 138 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4e, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x4f, 0x4d,
	0x2c, 0xc8, 0xcf, 0x0b, 0x4a, 0x4d, 0xcb, 0xcc, 0x4b, 0xd5, 0x03, 0x4b, 0x09, 0xb1, 0x43, 0xe5,
	0x94, 0x02, 0xb8, 0x04, 0x90, 0xa5, 0xe3, 0x4b, 0xf2, 0x83, 0x85, 0xa4, 0xb8, 0x38, 0x20, 0x62,
	0x9e, 0x2e, 0x12, 0x8c, 0x0a, 0x4c, 0x1a, 0xbc, 0x41, 0x70, 0xbe, 0x90, 0x02, 0x17, 0x77, 0x72,
	0x7e, 0x71, 0x49, 0x38, 0x98, 0x5f, 0x2c, 0xc1, 0x04, 0x96, 0x46, 0x16, 0x52, 0x12, 0xc2, 0x30,
	0xd1, 0xd9, 0x49, 0x29, 0x4a, 0x21, 0x33, 0x2f, 0xb3, 0x24, 0x33, 0x31, 0xa7, 0x24, 0xa3, 0x28,
	0x35, 0x15, 0xe2, 0xbe, 0xe4, 0xfc, 0x1c, 0xfd, 0xe4, 0x62, 0x98, 0x2b, 0x01, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x3f, 0x93, 0xc4, 0x4b, 0xb8, 0x00, 0x00, 0x00,
}