// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/enterSee.proto

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

type EnterSeeToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnterSeeToS) Reset()         { *m = EnterSeeToS{} }
func (m *EnterSeeToS) String() string { return proto.CompactTextString(m) }
func (*EnterSeeToS) ProtoMessage()    {}
func (*EnterSeeToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_9573704a7362d885, []int{0}
}

func (m *EnterSeeToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnterSeeToS.Unmarshal(m, b)
}
func (m *EnterSeeToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnterSeeToS.Marshal(b, m, deterministic)
}
func (m *EnterSeeToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnterSeeToS.Merge(m, src)
}
func (m *EnterSeeToS) XXX_Size() int {
	return xxx_messageInfo_EnterSeeToS.Size(m)
}
func (m *EnterSeeToS) XXX_DiscardUnknown() {
	xxx_messageInfo_EnterSeeToS.DiscardUnknown(m)
}

var xxx_messageInfo_EnterSeeToS proto.InternalMessageInfo

type EnterSeeToC struct {
	Objs                 []*ViewObj `protobuf:"bytes,1,rep,name=objs" json:"objs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *EnterSeeToC) Reset()         { *m = EnterSeeToC{} }
func (m *EnterSeeToC) String() string { return proto.CompactTextString(m) }
func (*EnterSeeToC) ProtoMessage()    {}
func (*EnterSeeToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_9573704a7362d885, []int{1}
}

func (m *EnterSeeToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnterSeeToC.Unmarshal(m, b)
}
func (m *EnterSeeToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnterSeeToC.Marshal(b, m, deterministic)
}
func (m *EnterSeeToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnterSeeToC.Merge(m, src)
}
func (m *EnterSeeToC) XXX_Size() int {
	return xxx_messageInfo_EnterSeeToC.Size(m)
}
func (m *EnterSeeToC) XXX_DiscardUnknown() {
	xxx_messageInfo_EnterSeeToC.DiscardUnknown(m)
}

var xxx_messageInfo_EnterSeeToC proto.InternalMessageInfo

func (m *EnterSeeToC) GetObjs() []*ViewObj {
	if m != nil {
		return m.Objs
	}
	return nil
}

func init() {
	proto.RegisterType((*EnterSeeToS)(nil), "message.enterSee_toS")
	proto.RegisterType((*EnterSeeToC)(nil), "message.enterSee_toC")
}

func init() { proto.RegisterFile("cs/proto/message/enterSee.proto", fileDescriptor_9573704a7362d885) }

var fileDescriptor_9573704a7362d885 = []byte{
	// 141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0xcd, 0x2b,
	0x49, 0x2d, 0x0a, 0x4e, 0x4d, 0xd5, 0x03, 0x0b, 0x0b, 0xb1, 0x43, 0xc5, 0xa5, 0x64, 0x31, 0x54,
	0x26, 0xe7, 0xe7, 0xe6, 0xe6, 0xe7, 0x41, 0xd4, 0x29, 0xf1, 0x71, 0xf1, 0xc0, 0x74, 0xc6, 0x97,
	0xe4, 0x07, 0x2b, 0x99, 0xa0, 0xf0, 0x9d, 0x85, 0x54, 0xb8, 0x58, 0xf2, 0x93, 0xb2, 0x8a, 0x25,
	0x18, 0x15, 0x98, 0x35, 0xb8, 0x8d, 0x04, 0xf4, 0xa0, 0x86, 0xe8, 0x95, 0x65, 0xa6, 0x96, 0xfb,
	0x27, 0x65, 0x05, 0x81, 0x65, 0x9d, 0x94, 0xa2, 0x14, 0x32, 0xf3, 0x32, 0x4b, 0x32, 0x13, 0x73,
	0x4a, 0x32, 0x8a, 0x52, 0x53, 0x21, 0x16, 0x26, 0xe7, 0xe7, 0xe8, 0x27, 0x17, 0xc3, 0xac, 0x05,
	0x04, 0x00, 0x00, 0xff, 0xff, 0xa2, 0x99, 0xcc, 0x7f, 0xb3, 0x00, 0x00, 0x00,
}
