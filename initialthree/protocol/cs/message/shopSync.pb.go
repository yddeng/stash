// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/shopSync.proto

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

type ShopSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ShopSyncToS) Reset()         { *m = ShopSyncToS{} }
func (m *ShopSyncToS) String() string { return proto.CompactTextString(m) }
func (*ShopSyncToS) ProtoMessage()    {}
func (*ShopSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_f914a2d718b1e907, []int{0}
}

func (m *ShopSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShopSyncToS.Unmarshal(m, b)
}
func (m *ShopSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShopSyncToS.Marshal(b, m, deterministic)
}
func (m *ShopSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShopSyncToS.Merge(m, src)
}
func (m *ShopSyncToS) XXX_Size() int {
	return xxx_messageInfo_ShopSyncToS.Size(m)
}
func (m *ShopSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_ShopSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_ShopSyncToS proto.InternalMessageInfo

type ShopSyncToC struct {
	IsAll                *bool      `protobuf:"varint,1,req,name=isAll" json:"isAll,omitempty"`
	Products             []*Product `protobuf:"bytes,2,rep,name=products" json:"products,omitempty"`
	Shops                []*Shop    `protobuf:"bytes,3,rep,name=shops" json:"shops,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *ShopSyncToC) Reset()         { *m = ShopSyncToC{} }
func (m *ShopSyncToC) String() string { return proto.CompactTextString(m) }
func (*ShopSyncToC) ProtoMessage()    {}
func (*ShopSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_f914a2d718b1e907, []int{1}
}

func (m *ShopSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ShopSyncToC.Unmarshal(m, b)
}
func (m *ShopSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ShopSyncToC.Marshal(b, m, deterministic)
}
func (m *ShopSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ShopSyncToC.Merge(m, src)
}
func (m *ShopSyncToC) XXX_Size() int {
	return xxx_messageInfo_ShopSyncToC.Size(m)
}
func (m *ShopSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_ShopSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_ShopSyncToC proto.InternalMessageInfo

func (m *ShopSyncToC) GetIsAll() bool {
	if m != nil && m.IsAll != nil {
		return *m.IsAll
	}
	return false
}

func (m *ShopSyncToC) GetProducts() []*Product {
	if m != nil {
		return m.Products
	}
	return nil
}

func (m *ShopSyncToC) GetShops() []*Shop {
	if m != nil {
		return m.Shops
	}
	return nil
}

type Shop struct {
	Id                   *int32   `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	AlreadyRefreshTimes  *int32   `protobuf:"varint,2,req,name=alreadyRefreshTimes" json:"alreadyRefreshTimes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Shop) Reset()         { *m = Shop{} }
func (m *Shop) String() string { return proto.CompactTextString(m) }
func (*Shop) ProtoMessage()    {}
func (*Shop) Descriptor() ([]byte, []int) {
	return fileDescriptor_f914a2d718b1e907, []int{2}
}

func (m *Shop) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Shop.Unmarshal(m, b)
}
func (m *Shop) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Shop.Marshal(b, m, deterministic)
}
func (m *Shop) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Shop.Merge(m, src)
}
func (m *Shop) XXX_Size() int {
	return xxx_messageInfo_Shop.Size(m)
}
func (m *Shop) XXX_DiscardUnknown() {
	xxx_messageInfo_Shop.DiscardUnknown(m)
}

var xxx_messageInfo_Shop proto.InternalMessageInfo

func (m *Shop) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Shop) GetAlreadyRefreshTimes() int32 {
	if m != nil && m.AlreadyRefreshTimes != nil {
		return *m.AlreadyRefreshTimes
	}
	return 0
}

type Product struct {
	Id                   *int32   `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	AlreadyBuyTimes      *int32   `protobuf:"varint,2,req,name=alreadyBuyTimes" json:"alreadyBuyTimes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Product) Reset()         { *m = Product{} }
func (m *Product) String() string { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()    {}
func (*Product) Descriptor() ([]byte, []int) {
	return fileDescriptor_f914a2d718b1e907, []int{3}
}

func (m *Product) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Product.Unmarshal(m, b)
}
func (m *Product) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Product.Marshal(b, m, deterministic)
}
func (m *Product) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Product.Merge(m, src)
}
func (m *Product) XXX_Size() int {
	return xxx_messageInfo_Product.Size(m)
}
func (m *Product) XXX_DiscardUnknown() {
	xxx_messageInfo_Product.DiscardUnknown(m)
}

var xxx_messageInfo_Product proto.InternalMessageInfo

func (m *Product) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Product) GetAlreadyBuyTimes() int32 {
	if m != nil && m.AlreadyBuyTimes != nil {
		return *m.AlreadyBuyTimes
	}
	return 0
}

func init() {
	proto.RegisterType((*ShopSyncToS)(nil), "message.shopSync_toS")
	proto.RegisterType((*ShopSyncToC)(nil), "message.shopSync_toC")
	proto.RegisterType((*Shop)(nil), "message.shop")
	proto.RegisterType((*Product)(nil), "message.product")
}

func init() { proto.RegisterFile("cs/proto/message/shopSync.proto", fileDescriptor_f914a2d718b1e907) }

var fileDescriptor_f914a2d718b1e907 = []byte{
	// 230 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x90, 0x3f, 0x4f, 0xc3, 0x30,
	0x14, 0xc4, 0x45, 0x4a, 0xd4, 0xea, 0x01, 0x05, 0x19, 0x86, 0x6c, 0x44, 0x66, 0xc9, 0x80, 0x62,
	0xc4, 0x37, 0xa0, 0x5d, 0x98, 0x5d, 0x26, 0x16, 0x64, 0xd9, 0x0f, 0x62, 0xc9, 0xad, 0x23, 0x3f,
	0x57, 0x28, 0xdf, 0x1e, 0x39, 0x71, 0x2b, 0xfe, 0x6d, 0x7e, 0x77, 0xe7, 0xdf, 0x49, 0x07, 0xb7,
	0x9a, 0x44, 0x1f, 0x7c, 0xf4, 0x62, 0x8b, 0x44, 0xea, 0x03, 0x05, 0x75, 0xbe, 0xdf, 0x0c, 0x3b,
	0xdd, 0x8e, 0x32, 0x9b, 0x67, 0x9d, 0x2f, 0xe1, 0xfc, 0x60, 0xbd, 0x45, 0xbf, 0xe1, 0x9f, 0x3f,
	0xee, 0x35, 0xbb, 0x81, 0xd2, 0xd2, 0x93, 0x73, 0xd5, 0x49, 0x5d, 0x34, 0x0b, 0x39, 0x1d, 0xec,
	0x1e, 0x16, 0x7d, 0xf0, 0x66, 0xaf, 0x23, 0x55, 0x45, 0x3d, 0x6b, 0xce, 0x1e, 0xaf, 0xda, 0x4c,
	0x6c, 0xb3, 0x21, 0x8f, 0x09, 0x76, 0x07, 0x65, 0x62, 0x52, 0x35, 0x1b, 0xa3, 0x17, 0xc7, 0x68,
	0x52, 0xe5, 0xe4, 0xf1, 0x67, 0x38, 0x4d, 0x0f, 0xb6, 0x84, 0xc2, 0x9a, 0xb1, 0xad, 0x94, 0x85,
	0x35, 0xec, 0x01, 0xae, 0x95, 0x0b, 0xa8, 0xcc, 0x20, 0xf1, 0x3d, 0x20, 0x75, 0x2f, 0x76, 0x8b,
	0xa9, 0x35, 0x05, 0xfe, 0xb3, 0xf8, 0x1a, 0xe6, 0xb9, 0xfa, 0x0f, 0xac, 0x81, 0xcb, 0xfc, 0x63,
	0xb5, 0x1f, 0xbe, 0x83, 0x7e, 0xcb, 0x2b, 0xfe, 0x5a, 0xdb, 0x9d, 0x8d, 0x56, 0xb9, 0xd8, 0x05,
	0xc4, 0x69, 0x4d, 0xed, 0x9d, 0xd0, 0x74, 0xd8, 0xf4, 0x2b, 0x00, 0x00, 0xff, 0xff, 0x8b, 0x1d,
	0x11, 0x8c, 0x66, 0x01, 0x00, 0x00,
}