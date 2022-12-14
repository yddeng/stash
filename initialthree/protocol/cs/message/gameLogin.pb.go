// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/gameLogin.proto

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

type GameLoginToS struct {
	UserID               *string  `protobuf:"bytes,1,req,name=userID" json:"userID,omitempty"`
	Token                *string  `protobuf:"bytes,2,req,name=token" json:"token,omitempty"`
	ServerID             *int32   `protobuf:"varint,3,req,name=serverID" json:"serverID,omitempty"`
	SeqNo                *int32   `protobuf:"varint,4,opt,name=seqNo" json:"seqNo,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GameLoginToS) Reset()         { *m = GameLoginToS{} }
func (m *GameLoginToS) String() string { return proto.CompactTextString(m) }
func (*GameLoginToS) ProtoMessage()    {}
func (*GameLoginToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_58444d112079827c, []int{0}
}

func (m *GameLoginToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameLoginToS.Unmarshal(m, b)
}
func (m *GameLoginToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameLoginToS.Marshal(b, m, deterministic)
}
func (m *GameLoginToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameLoginToS.Merge(m, src)
}
func (m *GameLoginToS) XXX_Size() int {
	return xxx_messageInfo_GameLoginToS.Size(m)
}
func (m *GameLoginToS) XXX_DiscardUnknown() {
	xxx_messageInfo_GameLoginToS.DiscardUnknown(m)
}

var xxx_messageInfo_GameLoginToS proto.InternalMessageInfo

func (m *GameLoginToS) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *GameLoginToS) GetToken() string {
	if m != nil && m.Token != nil {
		return *m.Token
	}
	return ""
}

func (m *GameLoginToS) GetServerID() int32 {
	if m != nil && m.ServerID != nil {
		return *m.ServerID
	}
	return 0
}

func (m *GameLoginToS) GetSeqNo() int32 {
	if m != nil && m.SeqNo != nil {
		return *m.SeqNo
	}
	return 0
}

type GameLoginToC struct {
	IsFirstLogin         *bool    `protobuf:"varint,1,opt,name=isFirstLogin" json:"isFirstLogin,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GameLoginToC) Reset()         { *m = GameLoginToC{} }
func (m *GameLoginToC) String() string { return proto.CompactTextString(m) }
func (*GameLoginToC) ProtoMessage()    {}
func (*GameLoginToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_58444d112079827c, []int{1}
}

func (m *GameLoginToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GameLoginToC.Unmarshal(m, b)
}
func (m *GameLoginToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GameLoginToC.Marshal(b, m, deterministic)
}
func (m *GameLoginToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GameLoginToC.Merge(m, src)
}
func (m *GameLoginToC) XXX_Size() int {
	return xxx_messageInfo_GameLoginToC.Size(m)
}
func (m *GameLoginToC) XXX_DiscardUnknown() {
	xxx_messageInfo_GameLoginToC.DiscardUnknown(m)
}

var xxx_messageInfo_GameLoginToC proto.InternalMessageInfo

func (m *GameLoginToC) GetIsFirstLogin() bool {
	if m != nil && m.IsFirstLogin != nil {
		return *m.IsFirstLogin
	}
	return false
}

func init() {
	proto.RegisterType((*GameLoginToS)(nil), "message.gameLogin_toS")
	proto.RegisterType((*GameLoginToC)(nil), "message.gameLogin_toC")
}

func init() { proto.RegisterFile("cs/proto/message/gameLogin.proto", fileDescriptor_58444d112079827c) }

var fileDescriptor_58444d112079827c = []byte{
	// 189 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x48, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x4f, 0xcc,
	0x4d, 0xf5, 0xc9, 0x4f, 0xcf, 0xcc, 0xd3, 0x03, 0x8b, 0x0b, 0xb1, 0x43, 0x25, 0x94, 0xf2, 0xb9,
	0x78, 0xe1, 0x72, 0xf1, 0x25, 0xf9, 0xc1, 0x42, 0x62, 0x5c, 0x6c, 0xa5, 0xc5, 0xa9, 0x45, 0x9e,
	0x2e, 0x12, 0x8c, 0x0a, 0x4c, 0x1a, 0x9c, 0x41, 0x50, 0x9e, 0x90, 0x08, 0x17, 0x6b, 0x49, 0x7e,
	0x76, 0x6a, 0x9e, 0x04, 0x13, 0x58, 0x18, 0xc2, 0x11, 0x92, 0xe2, 0xe2, 0x28, 0x4e, 0x2d, 0x2a,
	0x03, 0xab, 0x67, 0x56, 0x60, 0xd2, 0x60, 0x0d, 0x82, 0xf3, 0x41, 0x3a, 0x8a, 0x53, 0x0b, 0xfd,
	0xf2, 0x25, 0x58, 0x14, 0x18, 0x35, 0x58, 0x83, 0x20, 0x1c, 0x25, 0x63, 0x54, 0x0b, 0x9d, 0x85,
	0x94, 0xb8, 0x78, 0x32, 0x8b, 0xdd, 0x32, 0x8b, 0x8a, 0x4b, 0xc0, 0x62, 0x12, 0x8c, 0x0a, 0x8c,
	0x1a, 0x1c, 0x41, 0x28, 0x62, 0x4e, 0x4a, 0x51, 0x0a, 0x99, 0x79, 0x99, 0x25, 0x99, 0x89, 0x39,
	0x25, 0x19, 0x45, 0xa9, 0xa9, 0x10, 0xcf, 0x25, 0xe7, 0xe7, 0xe8, 0x27, 0x17, 0xc3, 0xbc, 0x08,
	0x08, 0x00, 0x00, 0xff, 0xff, 0x95, 0xf8, 0xdb, 0x3c, 0xf5, 0x00, 0x00, 0x00,
}
