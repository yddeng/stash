// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/enterMap.proto

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

type EnterMapToS struct {
	MapID                *int32   `protobuf:"varint,1,opt,name=mapID" json:"mapID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnterMapToS) Reset()         { *m = EnterMapToS{} }
func (m *EnterMapToS) String() string { return proto.CompactTextString(m) }
func (*EnterMapToS) ProtoMessage()    {}
func (*EnterMapToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_63fd54538f46ced2, []int{0}
}

func (m *EnterMapToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnterMapToS.Unmarshal(m, b)
}
func (m *EnterMapToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnterMapToS.Marshal(b, m, deterministic)
}
func (m *EnterMapToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnterMapToS.Merge(m, src)
}
func (m *EnterMapToS) XXX_Size() int {
	return xxx_messageInfo_EnterMapToS.Size(m)
}
func (m *EnterMapToS) XXX_DiscardUnknown() {
	xxx_messageInfo_EnterMapToS.DiscardUnknown(m)
}

var xxx_messageInfo_EnterMapToS proto.InternalMessageInfo

func (m *EnterMapToS) GetMapID() int32 {
	if m != nil && m.MapID != nil {
		return *m.MapID
	}
	return 0
}

type EnterMapToC struct {
	MapID                *int32    `protobuf:"varint,1,req,name=mapID" json:"mapID,omitempty"`
	Pos                  *Position `protobuf:"bytes,2,req,name=pos" json:"pos,omitempty"`
	Angle                *int32    `protobuf:"varint,3,req,name=angle" json:"angle,omitempty"`
	Scene                *int32    `protobuf:"varint,4,req,name=scene" json:"scene,omitempty"`
	MapLogic             *uint32   `protobuf:"varint,5,req,name=mapLogic" json:"mapLogic,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *EnterMapToC) Reset()         { *m = EnterMapToC{} }
func (m *EnterMapToC) String() string { return proto.CompactTextString(m) }
func (*EnterMapToC) ProtoMessage()    {}
func (*EnterMapToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_63fd54538f46ced2, []int{1}
}

func (m *EnterMapToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnterMapToC.Unmarshal(m, b)
}
func (m *EnterMapToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnterMapToC.Marshal(b, m, deterministic)
}
func (m *EnterMapToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnterMapToC.Merge(m, src)
}
func (m *EnterMapToC) XXX_Size() int {
	return xxx_messageInfo_EnterMapToC.Size(m)
}
func (m *EnterMapToC) XXX_DiscardUnknown() {
	xxx_messageInfo_EnterMapToC.DiscardUnknown(m)
}

var xxx_messageInfo_EnterMapToC proto.InternalMessageInfo

func (m *EnterMapToC) GetMapID() int32 {
	if m != nil && m.MapID != nil {
		return *m.MapID
	}
	return 0
}

func (m *EnterMapToC) GetPos() *Position {
	if m != nil {
		return m.Pos
	}
	return nil
}

func (m *EnterMapToC) GetAngle() int32 {
	if m != nil && m.Angle != nil {
		return *m.Angle
	}
	return 0
}

func (m *EnterMapToC) GetScene() int32 {
	if m != nil && m.Scene != nil {
		return *m.Scene
	}
	return 0
}

func (m *EnterMapToC) GetMapLogic() uint32 {
	if m != nil && m.MapLogic != nil {
		return *m.MapLogic
	}
	return 0
}

func init() {
	proto.RegisterType((*EnterMapToS)(nil), "message.enterMap_toS")
	proto.RegisterType((*EnterMapToC)(nil), "message.enterMap_toC")
}

func init() { proto.RegisterFile("cs/proto/message/enterMap.proto", fileDescriptor_63fd54538f46ced2) }

var fileDescriptor_63fd54538f46ced2 = []byte{
	// 210 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x8f, 0x41, 0x4b, 0xc3, 0x30,
	0x18, 0x86, 0x69, 0x6a, 0x51, 0xa2, 0x1e, 0x0c, 0x1e, 0x42, 0x41, 0x0c, 0xd5, 0x43, 0x4f, 0x0d,
	0xf8, 0x13, 0xd4, 0x8b, 0xe0, 0x60, 0x74, 0xb7, 0x5d, 0x46, 0x08, 0x1f, 0x5d, 0xa0, 0xc9, 0x17,
	0x9a, 0xfc, 0x91, 0xfd, 0xe3, 0xd1, 0xa6, 0x2d, 0x1b, 0x3b, 0x3e, 0x4f, 0x9e, 0x97, 0xf0, 0xd1,
	0x77, 0x1d, 0xa4, 0x1f, 0x30, 0xa2, 0xb4, 0x10, 0x82, 0xea, 0x40, 0x82, 0x8b, 0x30, 0x6c, 0x94,
	0x6f, 0x26, 0xcd, 0xee, 0x67, 0x5f, 0xbe, 0xdd, 0x94, 0x1a, 0xad, 0x45, 0x97, 0xba, 0xea, 0x93,
	0x3e, 0x2d, 0xcb, 0x43, 0xc4, 0x1d, 0x7b, 0xa5, 0x85, 0x55, 0xfe, 0xef, 0x97, 0x67, 0x22, 0xab,
	0x8b, 0x36, 0x41, 0x75, 0xca, 0xae, 0xb2, 0x9f, 0xcb, 0x8c, 0xac, 0x19, 0xfb, 0xa0, 0xb9, 0xc7,
	0xc0, 0x89, 0x20, 0xf5, 0xe3, 0xd7, 0x4b, 0x33, 0x7f, 0xd8, 0x6c, 0x31, 0x98, 0x68, 0xd0, 0xb5,
	0xe3, 0xeb, 0x38, 0x55, 0xae, 0xeb, 0x81, 0xe7, 0x69, 0x3a, 0xc1, 0x68, 0x83, 0x06, 0x07, 0xfc,
	0x2e, 0xd9, 0x09, 0x58, 0x49, 0x1f, 0xac, 0xf2, 0xff, 0xd8, 0x19, 0xcd, 0x0b, 0x41, 0xea, 0xe7,
	0x76, 0xe5, 0xef, 0x6a, 0x2f, 0x8c, 0x33, 0xd1, 0xa8, 0x3e, 0x1e, 0x07, 0x80, 0x74, 0xa4, 0xc6,
	0x5e, 0xea, 0xb0, 0x9c, 0x7a, 0x0e, 0x00, 0x00, 0xff, 0xff, 0x98, 0x5e, 0xc2, 0x21, 0x27, 0x01,
	0x00, 0x00,
}
