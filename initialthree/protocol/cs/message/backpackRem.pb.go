// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/backpackRem.proto

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

type BackpackRemEntity struct {
	Type                 *int32   `protobuf:"varint,1,opt,name=type" json:"type,omitempty"`
	Id                   *uint32  `protobuf:"varint,2,opt,name=id" json:"id,omitempty"`
	Count                *int32   `protobuf:"varint,3,opt,name=count" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BackpackRemEntity) Reset()         { *m = BackpackRemEntity{} }
func (m *BackpackRemEntity) String() string { return proto.CompactTextString(m) }
func (*BackpackRemEntity) ProtoMessage()    {}
func (*BackpackRemEntity) Descriptor() ([]byte, []int) {
	return fileDescriptor_d23ae2a0877b239d, []int{0}
}

func (m *BackpackRemEntity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackpackRemEntity.Unmarshal(m, b)
}
func (m *BackpackRemEntity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackpackRemEntity.Marshal(b, m, deterministic)
}
func (m *BackpackRemEntity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackpackRemEntity.Merge(m, src)
}
func (m *BackpackRemEntity) XXX_Size() int {
	return xxx_messageInfo_BackpackRemEntity.Size(m)
}
func (m *BackpackRemEntity) XXX_DiscardUnknown() {
	xxx_messageInfo_BackpackRemEntity.DiscardUnknown(m)
}

var xxx_messageInfo_BackpackRemEntity proto.InternalMessageInfo

func (m *BackpackRemEntity) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return 0
}

func (m *BackpackRemEntity) GetId() uint32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *BackpackRemEntity) GetCount() int32 {
	if m != nil && m.Count != nil {
		return *m.Count
	}
	return 0
}

type BackpackRemToS struct {
	RemEntities          []*BackpackRemEntity `protobuf:"bytes,1,rep,name=remEntities" json:"remEntities,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *BackpackRemToS) Reset()         { *m = BackpackRemToS{} }
func (m *BackpackRemToS) String() string { return proto.CompactTextString(m) }
func (*BackpackRemToS) ProtoMessage()    {}
func (*BackpackRemToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_d23ae2a0877b239d, []int{1}
}

func (m *BackpackRemToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackpackRemToS.Unmarshal(m, b)
}
func (m *BackpackRemToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackpackRemToS.Marshal(b, m, deterministic)
}
func (m *BackpackRemToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackpackRemToS.Merge(m, src)
}
func (m *BackpackRemToS) XXX_Size() int {
	return xxx_messageInfo_BackpackRemToS.Size(m)
}
func (m *BackpackRemToS) XXX_DiscardUnknown() {
	xxx_messageInfo_BackpackRemToS.DiscardUnknown(m)
}

var xxx_messageInfo_BackpackRemToS proto.InternalMessageInfo

func (m *BackpackRemToS) GetRemEntities() []*BackpackRemEntity {
	if m != nil {
		return m.RemEntities
	}
	return nil
}

type BackpackRemToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BackpackRemToC) Reset()         { *m = BackpackRemToC{} }
func (m *BackpackRemToC) String() string { return proto.CompactTextString(m) }
func (*BackpackRemToC) ProtoMessage()    {}
func (*BackpackRemToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_d23ae2a0877b239d, []int{2}
}

func (m *BackpackRemToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BackpackRemToC.Unmarshal(m, b)
}
func (m *BackpackRemToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BackpackRemToC.Marshal(b, m, deterministic)
}
func (m *BackpackRemToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BackpackRemToC.Merge(m, src)
}
func (m *BackpackRemToC) XXX_Size() int {
	return xxx_messageInfo_BackpackRemToC.Size(m)
}
func (m *BackpackRemToC) XXX_DiscardUnknown() {
	xxx_messageInfo_BackpackRemToC.DiscardUnknown(m)
}

var xxx_messageInfo_BackpackRemToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*BackpackRemEntity)(nil), "message.backpackRemEntity")
	proto.RegisterType((*BackpackRemToS)(nil), "message.backpackRem_toS")
	proto.RegisterType((*BackpackRemToC)(nil), "message.backpackRem_toC")
}

func init() { proto.RegisterFile("cs/proto/message/backpackRem.proto", fileDescriptor_d23ae2a0877b239d) }

var fileDescriptor_d23ae2a0877b239d = []byte{
	// 190 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4a, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x4a, 0x4c,
	0xce, 0x2e, 0x48, 0x4c, 0xce, 0x0e, 0x4a, 0xcd, 0xd5, 0x03, 0xcb, 0x08, 0xb1, 0x43, 0xa5, 0x94,
	0x7c, 0xb9, 0x04, 0x91, 0x64, 0x5d, 0xf3, 0x4a, 0x32, 0x4b, 0x2a, 0x85, 0x84, 0xb8, 0x58, 0x4a,
	0x2a, 0x0b, 0x52, 0x25, 0x18, 0x15, 0x18, 0x35, 0x58, 0x83, 0xc0, 0x6c, 0x21, 0x3e, 0x2e, 0xa6,
	0xcc, 0x14, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xde, 0x20, 0xa6, 0xcc, 0x14, 0x21, 0x11, 0x2e, 0xd6,
	0xe4, 0xfc, 0xd2, 0xbc, 0x12, 0x09, 0x66, 0xb0, 0x22, 0x08, 0x47, 0xc9, 0x9f, 0x8b, 0x1f, 0xc9,
	0xb8, 0xf8, 0x92, 0xfc, 0x60, 0x21, 0x1b, 0x2e, 0xee, 0x22, 0xa8, 0xc9, 0x99, 0xa9, 0xc5, 0x12,
	0x8c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0x52, 0x7a, 0x50, 0x07, 0xe8, 0x61, 0xd8, 0x1e, 0x84, 0xac,
	0x5c, 0x49, 0x10, 0xdd, 0x40, 0x67, 0x27, 0xa5, 0x28, 0x85, 0xcc, 0xbc, 0xcc, 0x92, 0xcc, 0xc4,
	0x9c, 0x92, 0x8c, 0xa2, 0xd4, 0x54, 0x88, 0x5f, 0x93, 0xf3, 0x73, 0xf4, 0x93, 0x8b, 0x61, 0x3e,
	0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0x16, 0x96, 0x1e, 0xea, 0x04, 0x01, 0x00, 0x00,
}