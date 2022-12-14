// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/baseSetSignature.proto

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

type BaseSetSignatureToS struct {
	Signature            *string  `protobuf:"bytes,1,opt,name=signature" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BaseSetSignatureToS) Reset()         { *m = BaseSetSignatureToS{} }
func (m *BaseSetSignatureToS) String() string { return proto.CompactTextString(m) }
func (*BaseSetSignatureToS) ProtoMessage()    {}
func (*BaseSetSignatureToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_6201ee52dc4418b6, []int{0}
}

func (m *BaseSetSignatureToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BaseSetSignatureToS.Unmarshal(m, b)
}
func (m *BaseSetSignatureToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BaseSetSignatureToS.Marshal(b, m, deterministic)
}
func (m *BaseSetSignatureToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BaseSetSignatureToS.Merge(m, src)
}
func (m *BaseSetSignatureToS) XXX_Size() int {
	return xxx_messageInfo_BaseSetSignatureToS.Size(m)
}
func (m *BaseSetSignatureToS) XXX_DiscardUnknown() {
	xxx_messageInfo_BaseSetSignatureToS.DiscardUnknown(m)
}

var xxx_messageInfo_BaseSetSignatureToS proto.InternalMessageInfo

func (m *BaseSetSignatureToS) GetSignature() string {
	if m != nil && m.Signature != nil {
		return *m.Signature
	}
	return ""
}

type BaseSetSignatureToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BaseSetSignatureToC) Reset()         { *m = BaseSetSignatureToC{} }
func (m *BaseSetSignatureToC) String() string { return proto.CompactTextString(m) }
func (*BaseSetSignatureToC) ProtoMessage()    {}
func (*BaseSetSignatureToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_6201ee52dc4418b6, []int{1}
}

func (m *BaseSetSignatureToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BaseSetSignatureToC.Unmarshal(m, b)
}
func (m *BaseSetSignatureToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BaseSetSignatureToC.Marshal(b, m, deterministic)
}
func (m *BaseSetSignatureToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BaseSetSignatureToC.Merge(m, src)
}
func (m *BaseSetSignatureToC) XXX_Size() int {
	return xxx_messageInfo_BaseSetSignatureToC.Size(m)
}
func (m *BaseSetSignatureToC) XXX_DiscardUnknown() {
	xxx_messageInfo_BaseSetSignatureToC.DiscardUnknown(m)
}

var xxx_messageInfo_BaseSetSignatureToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*BaseSetSignatureToS)(nil), "message.baseSetSignature_toS")
	proto.RegisterType((*BaseSetSignatureToC)(nil), "message.baseSetSignature_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/baseSetSignature.proto", fileDescriptor_6201ee52dc4418b6)
}

var fileDescriptor_6201ee52dc4418b6 = []byte{
	// 123 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4f, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x4a, 0x2c,
	0x4e, 0x0d, 0x4e, 0x2d, 0x09, 0xce, 0x4c, 0xcf, 0x4b, 0x2c, 0x29, 0x2d, 0x4a, 0xd5, 0x03, 0x4b,
	0x0b, 0xb1, 0x43, 0xe5, 0x95, 0x4c, 0xb8, 0x44, 0xd0, 0x95, 0xc4, 0x97, 0xe4, 0x07, 0x0b, 0xc9,
	0x70, 0x71, 0x16, 0xc3, 0x04, 0x24, 0x18, 0x15, 0x18, 0x35, 0x38, 0x83, 0x10, 0x02, 0x4a, 0x62,
	0x58, 0x75, 0x39, 0x3b, 0x29, 0x45, 0x29, 0x64, 0xe6, 0x65, 0x96, 0x64, 0x26, 0xe6, 0x94, 0x64,
	0x14, 0xa5, 0xa6, 0x42, 0xdc, 0x92, 0x9c, 0x9f, 0xa3, 0x9f, 0x5c, 0x0c, 0x73, 0x11, 0x20, 0x00,
	0x00, 0xff, 0xff, 0x86, 0x6c, 0xc9, 0xde, 0xa4, 0x00, 0x00, 0x00,
}
