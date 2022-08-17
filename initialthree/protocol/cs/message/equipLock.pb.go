// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/equipLock.proto

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

type EquipLockToS struct {
	EquipID              *uint32  `protobuf:"varint,1,req,name=equipID" json:"equipID,omitempty"`
	IsLock               *bool    `protobuf:"varint,2,req,name=isLock" json:"isLock,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EquipLockToS) Reset()         { *m = EquipLockToS{} }
func (m *EquipLockToS) String() string { return proto.CompactTextString(m) }
func (*EquipLockToS) ProtoMessage()    {}
func (*EquipLockToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c1238e37443f3b1, []int{0}
}

func (m *EquipLockToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EquipLockToS.Unmarshal(m, b)
}
func (m *EquipLockToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EquipLockToS.Marshal(b, m, deterministic)
}
func (m *EquipLockToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EquipLockToS.Merge(m, src)
}
func (m *EquipLockToS) XXX_Size() int {
	return xxx_messageInfo_EquipLockToS.Size(m)
}
func (m *EquipLockToS) XXX_DiscardUnknown() {
	xxx_messageInfo_EquipLockToS.DiscardUnknown(m)
}

var xxx_messageInfo_EquipLockToS proto.InternalMessageInfo

func (m *EquipLockToS) GetEquipID() uint32 {
	if m != nil && m.EquipID != nil {
		return *m.EquipID
	}
	return 0
}

func (m *EquipLockToS) GetIsLock() bool {
	if m != nil && m.IsLock != nil {
		return *m.IsLock
	}
	return false
}

type EquipLockToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EquipLockToC) Reset()         { *m = EquipLockToC{} }
func (m *EquipLockToC) String() string { return proto.CompactTextString(m) }
func (*EquipLockToC) ProtoMessage()    {}
func (*EquipLockToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_9c1238e37443f3b1, []int{1}
}

func (m *EquipLockToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EquipLockToC.Unmarshal(m, b)
}
func (m *EquipLockToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EquipLockToC.Marshal(b, m, deterministic)
}
func (m *EquipLockToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EquipLockToC.Merge(m, src)
}
func (m *EquipLockToC) XXX_Size() int {
	return xxx_messageInfo_EquipLockToC.Size(m)
}
func (m *EquipLockToC) XXX_DiscardUnknown() {
	xxx_messageInfo_EquipLockToC.DiscardUnknown(m)
}

var xxx_messageInfo_EquipLockToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*EquipLockToS)(nil), "message.equipLock_toS")
	proto.RegisterType((*EquipLockToC)(nil), "message.equipLock_toC")
}

func init() { proto.RegisterFile("cs/proto/message/equipLock.proto", fileDescriptor_9c1238e37443f3b1) }

var fileDescriptor_9c1238e37443f3b1 = []byte{
	// 134 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x48, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x2d, 0x2c,
	0xcd, 0x2c, 0xf0, 0xc9, 0x4f, 0xce, 0xd6, 0x03, 0x8b, 0x0b, 0xb1, 0x43, 0x25, 0x94, 0x1c, 0xb9,
	0x78, 0xe1, 0x72, 0xf1, 0x25, 0xf9, 0xc1, 0x42, 0x12, 0x5c, 0xec, 0x60, 0x01, 0x4f, 0x17, 0x09,
	0x46, 0x05, 0x26, 0x0d, 0xde, 0x20, 0x18, 0x57, 0x48, 0x8c, 0x8b, 0x2d, 0xb3, 0x18, 0xa4, 0x4e,
	0x82, 0x49, 0x81, 0x49, 0x83, 0x23, 0x08, 0xca, 0x53, 0xe2, 0x47, 0x35, 0xc2, 0xd9, 0x49, 0x29,
	0x4a, 0x21, 0x33, 0x2f, 0xb3, 0x24, 0x33, 0x31, 0xa7, 0x24, 0xa3, 0x28, 0x35, 0x15, 0xe2, 0x94,
	0xe4, 0xfc, 0x1c, 0xfd, 0xe4, 0x62, 0x98, 0x83, 0x00, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8c, 0xb3,
	0xf6, 0x8f, 0xa3, 0x00, 0x00, 0x00,
}
