// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/buyGold.proto

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

type BuyGoldToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BuyGoldToS) Reset()         { *m = BuyGoldToS{} }
func (m *BuyGoldToS) String() string { return proto.CompactTextString(m) }
func (*BuyGoldToS) ProtoMessage()    {}
func (*BuyGoldToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_840f3b082664b582, []int{0}
}

func (m *BuyGoldToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BuyGoldToS.Unmarshal(m, b)
}
func (m *BuyGoldToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BuyGoldToS.Marshal(b, m, deterministic)
}
func (m *BuyGoldToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BuyGoldToS.Merge(m, src)
}
func (m *BuyGoldToS) XXX_Size() int {
	return xxx_messageInfo_BuyGoldToS.Size(m)
}
func (m *BuyGoldToS) XXX_DiscardUnknown() {
	xxx_messageInfo_BuyGoldToS.DiscardUnknown(m)
}

var xxx_messageInfo_BuyGoldToS proto.InternalMessageInfo

type BuyGoldToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BuyGoldToC) Reset()         { *m = BuyGoldToC{} }
func (m *BuyGoldToC) String() string { return proto.CompactTextString(m) }
func (*BuyGoldToC) ProtoMessage()    {}
func (*BuyGoldToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_840f3b082664b582, []int{1}
}

func (m *BuyGoldToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BuyGoldToC.Unmarshal(m, b)
}
func (m *BuyGoldToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BuyGoldToC.Marshal(b, m, deterministic)
}
func (m *BuyGoldToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BuyGoldToC.Merge(m, src)
}
func (m *BuyGoldToC) XXX_Size() int {
	return xxx_messageInfo_BuyGoldToC.Size(m)
}
func (m *BuyGoldToC) XXX_DiscardUnknown() {
	xxx_messageInfo_BuyGoldToC.DiscardUnknown(m)
}

var xxx_messageInfo_BuyGoldToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*BuyGoldToS)(nil), "message.buyGold_toS")
	proto.RegisterType((*BuyGoldToC)(nil), "message.buyGold_toC")
}

func init() { proto.RegisterFile("cs/proto/message/buyGold.proto", fileDescriptor_840f3b082664b582) }

var fileDescriptor_840f3b082664b582 = []byte{
	// 96 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x2a, 0xad,
	0x74, 0xcf, 0xcf, 0x49, 0xd1, 0x03, 0x8b, 0x0a, 0xb1, 0x43, 0x85, 0x95, 0x78, 0xb9, 0xb8, 0xa1,
	0x32, 0xf1, 0x25, 0xf9, 0xc1, 0xa8, 0x5c, 0x67, 0x27, 0xa5, 0x28, 0x85, 0xcc, 0xbc, 0xcc, 0x92,
	0xcc, 0xc4, 0x9c, 0x92, 0x8c, 0xa2, 0xd4, 0x54, 0x88, 0x91, 0xc9, 0xf9, 0x39, 0xfa, 0xc9, 0xc5,
	0x30, 0x83, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa1, 0x7e, 0x16, 0xb6, 0x6b, 0x00, 0x00, 0x00,
}