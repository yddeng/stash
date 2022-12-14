// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/mailSync.proto

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

type MailSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MailSyncToS) Reset()         { *m = MailSyncToS{} }
func (m *MailSyncToS) String() string { return proto.CompactTextString(m) }
func (*MailSyncToS) ProtoMessage()    {}
func (*MailSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_985396339cd00e82, []int{0}
}

func (m *MailSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MailSyncToS.Unmarshal(m, b)
}
func (m *MailSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MailSyncToS.Marshal(b, m, deterministic)
}
func (m *MailSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MailSyncToS.Merge(m, src)
}
func (m *MailSyncToS) XXX_Size() int {
	return xxx_messageInfo_MailSyncToS.Size(m)
}
func (m *MailSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_MailSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_MailSyncToS proto.InternalMessageInfo

type MailSyncToC struct {
	IsAll                *bool    `protobuf:"varint,1,opt,name=isAll" json:"isAll,omitempty"`
	Mails                []*Mail  `protobuf:"bytes,2,rep,name=mails" json:"mails,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MailSyncToC) Reset()         { *m = MailSyncToC{} }
func (m *MailSyncToC) String() string { return proto.CompactTextString(m) }
func (*MailSyncToC) ProtoMessage()    {}
func (*MailSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_985396339cd00e82, []int{1}
}

func (m *MailSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MailSyncToC.Unmarshal(m, b)
}
func (m *MailSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MailSyncToC.Marshal(b, m, deterministic)
}
func (m *MailSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MailSyncToC.Merge(m, src)
}
func (m *MailSyncToC) XXX_Size() int {
	return xxx_messageInfo_MailSyncToC.Size(m)
}
func (m *MailSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_MailSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_MailSyncToC proto.InternalMessageInfo

func (m *MailSyncToC) GetIsAll() bool {
	if m != nil && m.IsAll != nil {
		return *m.IsAll
	}
	return false
}

func (m *MailSyncToC) GetMails() []*Mail {
	if m != nil {
		return m.Mails
	}
	return nil
}

func init() {
	proto.RegisterType((*MailSyncToS)(nil), "message.mailSync_toS")
	proto.RegisterType((*MailSyncToC)(nil), "message.mailSync_toC")
}

func init() { proto.RegisterFile("cs/proto/message/mailSync.proto", fileDescriptor_985396339cd00e82) }

var fileDescriptor_985396339cd00e82 = []byte{
	// 155 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0xcf, 0x4d, 0xcc,
	0xcc, 0x09, 0xae, 0xcc, 0x4b, 0xd6, 0x03, 0x0b, 0x0b, 0xb1, 0x43, 0xc5, 0xa5, 0x64, 0x31, 0x54,
	0x26, 0xe7, 0xe7, 0xe6, 0xe6, 0xe7, 0x41, 0xd4, 0x29, 0xf1, 0x71, 0xf1, 0xc0, 0x74, 0xc6, 0x97,
	0xe4, 0x07, 0x2b, 0x79, 0xa2, 0xf0, 0x9d, 0x85, 0x44, 0xb8, 0x58, 0x33, 0x8b, 0x1d, 0x73, 0x72,
	0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x82, 0x20, 0x1c, 0x21, 0x65, 0x2e, 0x56, 0x90, 0xaa, 0x62,
	0x09, 0x26, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x5e, 0x3d, 0xa8, 0xd9, 0x7a, 0xbe, 0x89, 0x99, 0x39,
	0x41, 0x10, 0x39, 0x27, 0xa5, 0x28, 0x85, 0xcc, 0xbc, 0xcc, 0x92, 0xcc, 0xc4, 0x9c, 0x92, 0x8c,
	0xa2, 0xd4, 0x54, 0x88, 0x2b, 0x92, 0xf3, 0x73, 0xf4, 0x93, 0x8b, 0x61, 0x6e, 0x01, 0x04, 0x00,
	0x00, 0xff, 0xff, 0x4e, 0x19, 0x32, 0x0d, 0xc8, 0x00, 0x00, 0x00,
}
