// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/attrSync.proto

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

type AttrSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AttrSyncToS) Reset()         { *m = AttrSyncToS{} }
func (m *AttrSyncToS) String() string { return proto.CompactTextString(m) }
func (*AttrSyncToS) ProtoMessage()    {}
func (*AttrSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_80062de2a34aeccc, []int{0}
}

func (m *AttrSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AttrSyncToS.Unmarshal(m, b)
}
func (m *AttrSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AttrSyncToS.Marshal(b, m, deterministic)
}
func (m *AttrSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttrSyncToS.Merge(m, src)
}
func (m *AttrSyncToS) XXX_Size() int {
	return xxx_messageInfo_AttrSyncToS.Size(m)
}
func (m *AttrSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_AttrSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_AttrSyncToS proto.InternalMessageInfo

type AttrSyncToC struct {
	IsAll                *bool    `protobuf:"varint,1,req,name=isAll" json:"isAll,omitempty"`
	Attrs                []*Attr  `protobuf:"bytes,2,rep,name=attrs" json:"attrs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AttrSyncToC) Reset()         { *m = AttrSyncToC{} }
func (m *AttrSyncToC) String() string { return proto.CompactTextString(m) }
func (*AttrSyncToC) ProtoMessage()    {}
func (*AttrSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_80062de2a34aeccc, []int{1}
}

func (m *AttrSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AttrSyncToC.Unmarshal(m, b)
}
func (m *AttrSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AttrSyncToC.Marshal(b, m, deterministic)
}
func (m *AttrSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AttrSyncToC.Merge(m, src)
}
func (m *AttrSyncToC) XXX_Size() int {
	return xxx_messageInfo_AttrSyncToC.Size(m)
}
func (m *AttrSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_AttrSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_AttrSyncToC proto.InternalMessageInfo

func (m *AttrSyncToC) GetIsAll() bool {
	if m != nil && m.IsAll != nil {
		return *m.IsAll
	}
	return false
}

func (m *AttrSyncToC) GetAttrs() []*Attr {
	if m != nil {
		return m.Attrs
	}
	return nil
}

func init() {
	proto.RegisterType((*AttrSyncToS)(nil), "message.attrSync_toS")
	proto.RegisterType((*AttrSyncToC)(nil), "message.attrSync_toC")
}

func init() { proto.RegisterFile("cs/proto/message/attrSync.proto", fileDescriptor_80062de2a34aeccc) }

var fileDescriptor_80062de2a34aeccc = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x2c, 0x29,
	0x29, 0x0a, 0xae, 0xcc, 0x4b, 0xd6, 0x03, 0x0b, 0x0b, 0xb1, 0x43, 0xc5, 0xa5, 0x64, 0x31, 0x54,
	0x26, 0xe7, 0xe7, 0xe6, 0xe6, 0xe7, 0x41, 0xd4, 0x29, 0xf1, 0x71, 0xf1, 0xc0, 0x74, 0xc6, 0x97,
	0xe4, 0x07, 0x2b, 0x79, 0xa2, 0xf0, 0x9d, 0x85, 0x44, 0xb8, 0x58, 0x33, 0x8b, 0x1d, 0x73, 0x72,
	0x24, 0x18, 0x15, 0x98, 0x34, 0x38, 0x82, 0x20, 0x1c, 0x21, 0x65, 0x2e, 0x56, 0x90, 0xaa, 0x62,
	0x09, 0x26, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x5e, 0x3d, 0xa8, 0xd9, 0x7a, 0x8e, 0x25, 0x25, 0x45,
	0x41, 0x10, 0x39, 0x27, 0xa5, 0x28, 0x85, 0xcc, 0xbc, 0xcc, 0x92, 0xcc, 0xc4, 0x9c, 0x92, 0x8c,
	0xa2, 0xd4, 0x54, 0x88, 0x2b, 0x92, 0xf3, 0x73, 0xf4, 0x93, 0x8b, 0x61, 0x6e, 0x01, 0x04, 0x00,
	0x00, 0xff, 0xff, 0xbb, 0x40, 0xc7, 0xbf, 0xc8, 0x00, 0x00, 0x00,
}