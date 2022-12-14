// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/signSync.proto

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

type SignSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignSyncToS) Reset()         { *m = SignSyncToS{} }
func (m *SignSyncToS) String() string { return proto.CompactTextString(m) }
func (*SignSyncToS) ProtoMessage()    {}
func (*SignSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba9baef5f1c16be1, []int{0}
}

func (m *SignSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignSyncToS.Unmarshal(m, b)
}
func (m *SignSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignSyncToS.Marshal(b, m, deterministic)
}
func (m *SignSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignSyncToS.Merge(m, src)
}
func (m *SignSyncToS) XXX_Size() int {
	return xxx_messageInfo_SignSyncToS.Size(m)
}
func (m *SignSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_SignSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_SignSyncToS proto.InternalMessageInfo

type SignSyncToC struct {
	IsAll                *bool    `protobuf:"varint,1,req,name=isAll" json:"isAll,omitempty"`
	SignList             []*Sign  `protobuf:"bytes,2,rep,name=signList" json:"signList,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignSyncToC) Reset()         { *m = SignSyncToC{} }
func (m *SignSyncToC) String() string { return proto.CompactTextString(m) }
func (*SignSyncToC) ProtoMessage()    {}
func (*SignSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba9baef5f1c16be1, []int{1}
}

func (m *SignSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignSyncToC.Unmarshal(m, b)
}
func (m *SignSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignSyncToC.Marshal(b, m, deterministic)
}
func (m *SignSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignSyncToC.Merge(m, src)
}
func (m *SignSyncToC) XXX_Size() int {
	return xxx_messageInfo_SignSyncToC.Size(m)
}
func (m *SignSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_SignSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_SignSyncToC proto.InternalMessageInfo

func (m *SignSyncToC) GetIsAll() bool {
	if m != nil && m.IsAll != nil {
		return *m.IsAll
	}
	return false
}

func (m *SignSyncToC) GetSignList() []*Sign {
	if m != nil {
		return m.SignList
	}
	return nil
}

type Sign struct {
	Id                   *int32   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	SignTimes            *int32   `protobuf:"varint,2,opt,name=signTimes" json:"signTimes,omitempty"`
	LastSignTime         *int64   `protobuf:"varint,3,opt,name=lastSignTime" json:"lastSignTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Sign) Reset()         { *m = Sign{} }
func (m *Sign) String() string { return proto.CompactTextString(m) }
func (*Sign) ProtoMessage()    {}
func (*Sign) Descriptor() ([]byte, []int) {
	return fileDescriptor_ba9baef5f1c16be1, []int{2}
}

func (m *Sign) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Sign.Unmarshal(m, b)
}
func (m *Sign) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Sign.Marshal(b, m, deterministic)
}
func (m *Sign) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Sign.Merge(m, src)
}
func (m *Sign) XXX_Size() int {
	return xxx_messageInfo_Sign.Size(m)
}
func (m *Sign) XXX_DiscardUnknown() {
	xxx_messageInfo_Sign.DiscardUnknown(m)
}

var xxx_messageInfo_Sign proto.InternalMessageInfo

func (m *Sign) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *Sign) GetSignTimes() int32 {
	if m != nil && m.SignTimes != nil {
		return *m.SignTimes
	}
	return 0
}

func (m *Sign) GetLastSignTime() int64 {
	if m != nil && m.LastSignTime != nil {
		return *m.LastSignTime
	}
	return 0
}

func init() {
	proto.RegisterType((*SignSyncToS)(nil), "message.signSync_toS")
	proto.RegisterType((*SignSyncToC)(nil), "message.signSync_toC")
	proto.RegisterType((*Sign)(nil), "message.sign")
}

func init() { proto.RegisterFile("cs/proto/message/signSync.proto", fileDescriptor_ba9baef5f1c16be1) }

var fileDescriptor_ba9baef5f1c16be1 = []byte{
	// 202 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0xcf, 0x3f, 0x4b, 0xc5, 0x30,
	0x14, 0x05, 0x70, 0x92, 0xfa, 0xf0, 0x79, 0x7d, 0xbe, 0x21, 0x38, 0x64, 0x10, 0x0c, 0x99, 0xe2,
	0xd2, 0x80, 0xdf, 0x40, 0x5d, 0x05, 0x21, 0x75, 0x10, 0x17, 0x29, 0x79, 0xa1, 0x5e, 0x48, 0x1b,
	0xe9, 0xcd, 0xe2, 0xb7, 0x97, 0xf4, 0x8f, 0x8f, 0x6e, 0x39, 0xbf, 0x13, 0x2e, 0x1c, 0xb8, 0xf7,
	0x64, 0x7f, 0xc6, 0x94, 0x93, 0xed, 0x03, 0x51, 0xdb, 0x05, 0x4b, 0xd8, 0x0d, 0xcd, 0xef, 0xe0,
	0xeb, 0x89, 0xc5, 0xe5, 0xe2, 0xfa, 0x08, 0x87, 0xb5, 0xfa, 0xca, 0xa9, 0xd1, 0x6f, 0x9b, 0xfc,
	0x22, 0x6e, 0x61, 0x87, 0xf4, 0x14, 0xa3, 0x64, 0x8a, 0x9b, 0xbd, 0x9b, 0x83, 0x78, 0x80, 0x7d,
	0xf9, 0xf5, 0x8a, 0x94, 0x25, 0x57, 0x95, 0xb9, 0x7e, 0xbc, 0xa9, 0x97, 0x8b, 0x75, 0x29, 0xdc,
	0x7f, 0xad, 0x3f, 0xe0, 0xa2, 0xbc, 0xc5, 0x11, 0x38, 0x9e, 0x24, 0x53, 0xcc, 0xec, 0x1c, 0xc7,
	0x93, 0xb8, 0x83, 0xab, 0xe2, 0xef, 0xd8, 0x07, 0x92, 0x7c, 0xe2, 0x33, 0x08, 0x0d, 0x87, 0xd8,
	0x52, 0x6e, 0x16, 0x90, 0x95, 0x62, 0xa6, 0x72, 0x1b, 0x7b, 0xd6, 0x9f, 0x0a, 0x07, 0xcc, 0xd8,
	0xc6, 0xfc, 0x3d, 0x86, 0x30, 0x0f, 0xf6, 0x29, 0x5a, 0x4f, 0xeb, 0xec, 0xbf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0xcd, 0x17, 0x7f, 0x77, 0x09, 0x01, 0x00, 0x00,
}
