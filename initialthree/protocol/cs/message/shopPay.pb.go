// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/shopPay.proto

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

type ShopPayToS struct {
	PayID                *int32   `protobuf:"varint,1,opt,name=payID" json:"payID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShopPayToS) Reset()         { *m = ShopPayToS{} }
func (m *ShopPayToS) String() string { return proto.CompactTextString(m) }
func (*ShopPayToS) ProtoMessage()    {}
func (*ShopPayToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_feafa82453f7abbf, []int{0}
}

func (m *ShopPayToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShopPayToS.Unmarshal(m, b)
}
func (m *ShopPayToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShopPayToS.Marshal(b, m, deterministic)
}
func (m *ShopPayToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShopPayToS.Merge(m, src)
}
func (m *ShopPayToS) XXX_Size() int {
	return xxx_messageInfo_ShopPayToS.Size(m)
}
func (m *ShopPayToS) XXX_DiscardUnknown() {
	xxx_messageInfo_ShopPayToS.DiscardUnknown(m)
}

var xxx_messageInfo_ShopPayToS proto.InternalMessageInfo

func (m *ShopPayToS) GetPayID() int32 {
	if m != nil && m.PayID != nil {
		return *m.PayID
	}
	return 0
}

type ShopPayToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShopPayToC) Reset()         { *m = ShopPayToC{} }
func (m *ShopPayToC) String() string { return proto.CompactTextString(m) }
func (*ShopPayToC) ProtoMessage()    {}
func (*ShopPayToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_feafa82453f7abbf, []int{1}
}

func (m *ShopPayToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShopPayToC.Unmarshal(m, b)
}
func (m *ShopPayToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShopPayToC.Marshal(b, m, deterministic)
}
func (m *ShopPayToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShopPayToC.Merge(m, src)
}
func (m *ShopPayToC) XXX_Size() int {
	return xxx_messageInfo_ShopPayToC.Size(m)
}
func (m *ShopPayToC) XXX_DiscardUnknown() {
	xxx_messageInfo_ShopPayToC.DiscardUnknown(m)
}

var xxx_messageInfo_ShopPayToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ShopPayToS)(nil), "message.shopPay_toS")
	proto.RegisterType((*ShopPayToC)(nil), "message.shopPay_toC")
}

func init() { proto.RegisterFile("cs/proto/message/shopPay.proto", fileDescriptor_feafa82453f7abbf) }

var fileDescriptor_feafa82453f7abbf = []byte{
	// 116 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0xce, 0xc8,
	0x2f, 0x08, 0x48, 0xac, 0xd4, 0x03, 0x8b, 0x0a, 0xb1, 0x43, 0x85, 0x95, 0x94, 0xb9, 0xb8, 0xa1,
	0x32, 0xf1, 0x25, 0xf9, 0xc1, 0x42, 0x22, 0x5c, 0xac, 0x05, 0x89, 0x95, 0x9e, 0x2e, 0x12, 0x8c,
	0x0a, 0x8c, 0x1a, 0xac, 0x41, 0x10, 0x8e, 0x12, 0x2f, 0xb2, 0x22, 0x67, 0x27, 0xa5, 0x28, 0x85,
	0xcc, 0xbc, 0xcc, 0x92, 0xcc, 0xc4, 0x9c, 0x92, 0x8c, 0xa2, 0xd4, 0x54, 0x88, 0x45, 0xc9, 0xf9,
	0x39, 0xfa, 0xc9, 0xc5, 0x30, 0xeb, 0x00, 0x01, 0x00, 0x00, 0xff, 0xff, 0x50, 0x08, 0x6c, 0x05,
	0x81, 0x00, 0x00, 0x00,
}