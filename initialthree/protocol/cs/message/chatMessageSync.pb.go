// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/chatMessageSync.proto

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

type ChatMessageSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChatMessageSyncToS) Reset()         { *m = ChatMessageSyncToS{} }
func (m *ChatMessageSyncToS) String() string { return proto.CompactTextString(m) }
func (*ChatMessageSyncToS) ProtoMessage()    {}
func (*ChatMessageSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_89d196ddeda7679b, []int{0}
}

func (m *ChatMessageSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChatMessageSyncToS.Unmarshal(m, b)
}
func (m *ChatMessageSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChatMessageSyncToS.Marshal(b, m, deterministic)
}
func (m *ChatMessageSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChatMessageSyncToS.Merge(m, src)
}
func (m *ChatMessageSyncToS) XXX_Size() int {
	return xxx_messageInfo_ChatMessageSyncToS.Size(m)
}
func (m *ChatMessageSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_ChatMessageSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_ChatMessageSyncToS proto.InternalMessageInfo

type ChatMessageSyncToC struct {
	UserID               *string  `protobuf:"bytes,1,opt,name=userID" json:"userID,omitempty"`
	GameID               *uint64  `protobuf:"varint,2,opt,name=gameID" json:"gameID,omitempty"`
	Name                 *string  `protobuf:"bytes,3,opt,name=name" json:"name,omitempty"`
	Portrait             *int32   `protobuf:"varint,4,opt,name=Portrait" json:"Portrait,omitempty"`
	PortraitFrame        *int32   `protobuf:"varint,5,opt,name=PortraitFrame" json:"PortraitFrame,omitempty"`
	Message              *string  `protobuf:"bytes,6,opt,name=message" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChatMessageSyncToC) Reset()         { *m = ChatMessageSyncToC{} }
func (m *ChatMessageSyncToC) String() string { return proto.CompactTextString(m) }
func (*ChatMessageSyncToC) ProtoMessage()    {}
func (*ChatMessageSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_89d196ddeda7679b, []int{1}
}

func (m *ChatMessageSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChatMessageSyncToC.Unmarshal(m, b)
}
func (m *ChatMessageSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChatMessageSyncToC.Marshal(b, m, deterministic)
}
func (m *ChatMessageSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChatMessageSyncToC.Merge(m, src)
}
func (m *ChatMessageSyncToC) XXX_Size() int {
	return xxx_messageInfo_ChatMessageSyncToC.Size(m)
}
func (m *ChatMessageSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_ChatMessageSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_ChatMessageSyncToC proto.InternalMessageInfo

func (m *ChatMessageSyncToC) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *ChatMessageSyncToC) GetGameID() uint64 {
	if m != nil && m.GameID != nil {
		return *m.GameID
	}
	return 0
}

func (m *ChatMessageSyncToC) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *ChatMessageSyncToC) GetPortrait() int32 {
	if m != nil && m.Portrait != nil {
		return *m.Portrait
	}
	return 0
}

func (m *ChatMessageSyncToC) GetPortraitFrame() int32 {
	if m != nil && m.PortraitFrame != nil {
		return *m.PortraitFrame
	}
	return 0
}

func (m *ChatMessageSyncToC) GetMessage() string {
	if m != nil && m.Message != nil {
		return *m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*ChatMessageSyncToS)(nil), "message.chatMessageSync_toS")
	proto.RegisterType((*ChatMessageSyncToC)(nil), "message.chatMessageSync_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/chatMessageSync.proto", fileDescriptor_89d196ddeda7679b)
}

var fileDescriptor_89d196ddeda7679b = []byte{
	// 199 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4b, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0xce, 0x48,
	0x2c, 0xf1, 0x85, 0xb0, 0x83, 0x2b, 0xf3, 0x92, 0xf5, 0xc0, 0xb2, 0x42, 0xec, 0x50, 0x69, 0x25,
	0x51, 0x2e, 0x61, 0x34, 0x15, 0xf1, 0x25, 0xf9, 0xc1, 0x4a, 0x5b, 0x19, 0xb1, 0x89, 0x3b, 0x0b,
	0x89, 0x71, 0xb1, 0x95, 0x16, 0xa7, 0x16, 0x79, 0xba, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06,
	0x41, 0x79, 0x20, 0xf1, 0xf4, 0xc4, 0xdc, 0x54, 0x4f, 0x17, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x96,
	0x20, 0x28, 0x4f, 0x48, 0x88, 0x8b, 0x25, 0x2f, 0x31, 0x37, 0x55, 0x82, 0x19, 0xac, 0x1a, 0xcc,
	0x16, 0x92, 0xe2, 0xe2, 0x08, 0xc8, 0x2f, 0x2a, 0x29, 0x4a, 0xcc, 0x2c, 0x91, 0x60, 0x51, 0x60,
	0xd4, 0x60, 0x0d, 0x82, 0xf3, 0x85, 0x54, 0xb8, 0x78, 0x61, 0x6c, 0xb7, 0x22, 0x90, 0x46, 0x56,
	0xb0, 0x02, 0x54, 0x41, 0x21, 0x09, 0x2e, 0x98, 0xfb, 0x25, 0xd8, 0xc0, 0x06, 0xc3, 0xb8, 0x4e,
	0x4a, 0x51, 0x0a, 0x99, 0x79, 0x99, 0x25, 0x99, 0x89, 0x39, 0x25, 0x19, 0x45, 0xa9, 0xa9, 0x90,
	0xb0, 0x48, 0xce, 0xcf, 0xd1, 0x4f, 0x2e, 0x86, 0x85, 0x08, 0x20, 0x00, 0x00, 0xff, 0xff, 0xfb,
	0x67, 0xc5, 0x7d, 0x24, 0x01, 0x00, 0x00,
}
