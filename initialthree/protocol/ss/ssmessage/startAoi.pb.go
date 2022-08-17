// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss/proto/ssmessage/startAoi.proto

package ssmessage

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

type StartAoi struct {
	SceneIdx             *int32   `protobuf:"varint,1,opt,name=sceneIdx" json:"sceneIdx,omitempty"`
	UserID               *string  `protobuf:"bytes,2,opt,name=userID" json:"userID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartAoi) Reset()         { *m = StartAoi{} }
func (m *StartAoi) String() string { return proto.CompactTextString(m) }
func (*StartAoi) ProtoMessage()    {}
func (*StartAoi) Descriptor() ([]byte, []int) {
	return fileDescriptor_c117b19fdd48441c, []int{0}
}

func (m *StartAoi) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartAoi.Unmarshal(m, b)
}
func (m *StartAoi) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartAoi.Marshal(b, m, deterministic)
}
func (m *StartAoi) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartAoi.Merge(m, src)
}
func (m *StartAoi) XXX_Size() int {
	return xxx_messageInfo_StartAoi.Size(m)
}
func (m *StartAoi) XXX_DiscardUnknown() {
	xxx_messageInfo_StartAoi.DiscardUnknown(m)
}

var xxx_messageInfo_StartAoi proto.InternalMessageInfo

func (m *StartAoi) GetSceneIdx() int32 {
	if m != nil && m.SceneIdx != nil {
		return *m.SceneIdx
	}
	return 0
}

func (m *StartAoi) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func init() {
	proto.RegisterType((*StartAoi)(nil), "ssmessage.startAoi")
}

func init() { proto.RegisterFile("ss/proto/ssmessage/startAoi.proto", fileDescriptor_c117b19fdd48441c) }

var fileDescriptor_c117b19fdd48441c = []byte{
	// 131 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2c, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x2e, 0xce, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f,
	0x2e, 0x49, 0x2c, 0x2a, 0x71, 0xcc, 0xcf, 0xd4, 0x03, 0x4b, 0x08, 0x71, 0xc2, 0x65, 0x94, 0xec,
	0xb8, 0x38, 0x60, 0x92, 0x42, 0x52, 0x5c, 0x1c, 0xc5, 0xc9, 0xa9, 0x79, 0xa9, 0x9e, 0x29, 0x15,
	0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0xac, 0x41, 0x70, 0xbe, 0x90, 0x18, 0x17, 0x5b, 0x69, 0x71, 0x6a,
	0x91, 0xa7, 0x8b, 0x04, 0x93, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x94, 0xe7, 0xa4, 0x12, 0xa5, 0x94,
	0x99, 0x97, 0x59, 0x92, 0x99, 0x98, 0x53, 0x92, 0x51, 0x94, 0x9a, 0x0a, 0xb1, 0x39, 0x39, 0x3f,
	0x47, 0xbf, 0xb8, 0x18, 0x61, 0x3f, 0x20, 0x00, 0x00, 0xff, 0xff, 0xd1, 0x0c, 0x57, 0x25, 0x94,
	0x00, 0x00, 0x00,
}