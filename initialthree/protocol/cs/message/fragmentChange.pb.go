// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/fragmentChange.proto

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

type FragmentChangeToS struct {
	Id                   *int32   `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Count                *int32   `protobuf:"varint,2,opt,name=count" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FragmentChangeToS) Reset()         { *m = FragmentChangeToS{} }
func (m *FragmentChangeToS) String() string { return proto.CompactTextString(m) }
func (*FragmentChangeToS) ProtoMessage()    {}
func (*FragmentChangeToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0f9cd01dc8c83f, []int{0}
}

func (m *FragmentChangeToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FragmentChangeToS.Unmarshal(m, b)
}
func (m *FragmentChangeToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FragmentChangeToS.Marshal(b, m, deterministic)
}
func (m *FragmentChangeToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FragmentChangeToS.Merge(m, src)
}
func (m *FragmentChangeToS) XXX_Size() int {
	return xxx_messageInfo_FragmentChangeToS.Size(m)
}
func (m *FragmentChangeToS) XXX_DiscardUnknown() {
	xxx_messageInfo_FragmentChangeToS.DiscardUnknown(m)
}

var xxx_messageInfo_FragmentChangeToS proto.InternalMessageInfo

func (m *FragmentChangeToS) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *FragmentChangeToS) GetCount() int32 {
	if m != nil && m.Count != nil {
		return *m.Count
	}
	return 0
}

type FragmentChangeToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FragmentChangeToC) Reset()         { *m = FragmentChangeToC{} }
func (m *FragmentChangeToC) String() string { return proto.CompactTextString(m) }
func (*FragmentChangeToC) ProtoMessage()    {}
func (*FragmentChangeToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_3f0f9cd01dc8c83f, []int{1}
}

func (m *FragmentChangeToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FragmentChangeToC.Unmarshal(m, b)
}
func (m *FragmentChangeToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FragmentChangeToC.Marshal(b, m, deterministic)
}
func (m *FragmentChangeToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FragmentChangeToC.Merge(m, src)
}
func (m *FragmentChangeToC) XXX_Size() int {
	return xxx_messageInfo_FragmentChangeToC.Size(m)
}
func (m *FragmentChangeToC) XXX_DiscardUnknown() {
	xxx_messageInfo_FragmentChangeToC.DiscardUnknown(m)
}

var xxx_messageInfo_FragmentChangeToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*FragmentChangeToS)(nil), "message.fragmentChange_toS")
	proto.RegisterType((*FragmentChangeToC)(nil), "message.fragmentChange_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/fragmentChange.proto", fileDescriptor_3f0f9cd01dc8c83f)
}

var fileDescriptor_3f0f9cd01dc8c83f = []byte{
	// 136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4d, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x2b, 0x4a,
	0x4c, 0xcf, 0x4d, 0xcd, 0x2b, 0x71, 0xce, 0x48, 0xcc, 0x4b, 0x4f, 0xd5, 0x03, 0x4b, 0x0a, 0xb1,
	0x43, 0x65, 0x95, 0xac, 0xb8, 0x84, 0x50, 0x15, 0xc4, 0x97, 0xe4, 0x07, 0x0b, 0xf1, 0x71, 0x31,
	0x65, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0xb0, 0x06, 0x31, 0x65, 0xa6, 0x08, 0x89, 0x70, 0xb1,
	0x26, 0xe7, 0x97, 0xe6, 0x95, 0x48, 0x30, 0x81, 0x85, 0x20, 0x1c, 0x25, 0x11, 0x2c, 0x7a, 0x9d,
	0x9d, 0x94, 0xa2, 0x14, 0x32, 0xf3, 0x32, 0x4b, 0x32, 0x13, 0x73, 0x4a, 0x32, 0x8a, 0x52, 0x53,
	0x21, 0xae, 0x49, 0xce, 0xcf, 0xd1, 0x4f, 0x2e, 0x86, 0xb9, 0x09, 0x10, 0x00, 0x00, 0xff, 0xff,
	0xba, 0x09, 0x4e, 0x41, 0xa6, 0x00, 0x00, 0x00,
}
