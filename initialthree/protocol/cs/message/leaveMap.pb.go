// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/leaveMap.proto

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

type LeaveMapToS struct {
	ReqMapID             *int32    `protobuf:"varint,1,opt,name=ReqMapID" json:"ReqMapID,omitempty"`
	ReqPosition          *Position `protobuf:"bytes,2,opt,name=ReqPosition" json:"ReqPosition,omitempty"`
	ReqAngle             *int32    `protobuf:"varint,3,opt,name=ReqAngle" json:"ReqAngle,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *LeaveMapToS) Reset()         { *m = LeaveMapToS{} }
func (m *LeaveMapToS) String() string { return proto.CompactTextString(m) }
func (*LeaveMapToS) ProtoMessage()    {}
func (*LeaveMapToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_12462c8dbde496ff, []int{0}
}

func (m *LeaveMapToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaveMapToS.Unmarshal(m, b)
}
func (m *LeaveMapToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaveMapToS.Marshal(b, m, deterministic)
}
func (m *LeaveMapToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaveMapToS.Merge(m, src)
}
func (m *LeaveMapToS) XXX_Size() int {
	return xxx_messageInfo_LeaveMapToS.Size(m)
}
func (m *LeaveMapToS) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaveMapToS.DiscardUnknown(m)
}

var xxx_messageInfo_LeaveMapToS proto.InternalMessageInfo

func (m *LeaveMapToS) GetReqMapID() int32 {
	if m != nil && m.ReqMapID != nil {
		return *m.ReqMapID
	}
	return 0
}

func (m *LeaveMapToS) GetReqPosition() *Position {
	if m != nil {
		return m.ReqPosition
	}
	return nil
}

func (m *LeaveMapToS) GetReqAngle() int32 {
	if m != nil && m.ReqAngle != nil {
		return *m.ReqAngle
	}
	return 0
}

type LeaveMapToC struct {
	MapID                *int32    `protobuf:"varint,1,opt,name=mapID" json:"mapID,omitempty"`
	Position             *Position `protobuf:"bytes,2,opt,name=position" json:"position,omitempty"`
	Angle                *int32    `protobuf:"varint,3,opt,name=angle" json:"angle,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *LeaveMapToC) Reset()         { *m = LeaveMapToC{} }
func (m *LeaveMapToC) String() string { return proto.CompactTextString(m) }
func (*LeaveMapToC) ProtoMessage()    {}
func (*LeaveMapToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_12462c8dbde496ff, []int{1}
}

func (m *LeaveMapToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaveMapToC.Unmarshal(m, b)
}
func (m *LeaveMapToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaveMapToC.Marshal(b, m, deterministic)
}
func (m *LeaveMapToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaveMapToC.Merge(m, src)
}
func (m *LeaveMapToC) XXX_Size() int {
	return xxx_messageInfo_LeaveMapToC.Size(m)
}
func (m *LeaveMapToC) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaveMapToC.DiscardUnknown(m)
}

var xxx_messageInfo_LeaveMapToC proto.InternalMessageInfo

func (m *LeaveMapToC) GetMapID() int32 {
	if m != nil && m.MapID != nil {
		return *m.MapID
	}
	return 0
}

func (m *LeaveMapToC) GetPosition() *Position {
	if m != nil {
		return m.Position
	}
	return nil
}

func (m *LeaveMapToC) GetAngle() int32 {
	if m != nil && m.Angle != nil {
		return *m.Angle
	}
	return 0
}

func init() {
	proto.RegisterType((*LeaveMapToS)(nil), "message.leaveMap_toS")
	proto.RegisterType((*LeaveMapToC)(nil), "message.leaveMap_toC")
}

func init() { proto.RegisterFile("cs/proto/message/leaveMap.proto", fileDescriptor_12462c8dbde496ff) }

var fileDescriptor_12462c8dbde496ff = []byte{
	// 203 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4f, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0xcf, 0x49, 0x4d,
	0x2c, 0x4b, 0xf5, 0x4d, 0x2c, 0xd0, 0x03, 0x0b, 0x0b, 0xb1, 0x43, 0xc5, 0xa5, 0x64, 0x31, 0x54,
	0x26, 0xe7, 0xe7, 0xe6, 0xe6, 0xe7, 0x41, 0xd4, 0x29, 0x55, 0x73, 0xf1, 0xc0, 0x74, 0xc6, 0x97,
	0xe4, 0x07, 0x0b, 0x49, 0x71, 0x71, 0x04, 0xa5, 0x16, 0xfa, 0x26, 0x16, 0x78, 0xba, 0x48, 0x30,
	0x2a, 0x30, 0x6a, 0xb0, 0x06, 0xc1, 0xf9, 0x42, 0xc6, 0x5c, 0xdc, 0x41, 0xa9, 0x85, 0x01, 0xf9,
	0xc5, 0x99, 0x25, 0x99, 0xf9, 0x79, 0x12, 0x4c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0x82, 0x7a, 0x50,
	0x73, 0xf5, 0x60, 0x12, 0x41, 0xc8, 0xaa, 0xa0, 0x06, 0x3a, 0xe6, 0xa5, 0xe7, 0xa4, 0x4a, 0x30,
	0xc3, 0x0d, 0x04, 0xf3, 0x95, 0x32, 0x51, 0x2c, 0x77, 0x16, 0x12, 0xe1, 0x62, 0xcd, 0x45, 0xb2,
	0x19, 0xc2, 0x11, 0xd2, 0xe5, 0xe2, 0x28, 0x20, 0x68, 0x27, 0x5c, 0x09, 0xc8, 0x90, 0x44, 0x24,
	0xdb, 0x20, 0x1c, 0x27, 0xa5, 0x28, 0x85, 0xcc, 0xbc, 0xcc, 0x92, 0xcc, 0xc4, 0x9c, 0x92, 0x8c,
	0xa2, 0xd4, 0x54, 0x48, 0x90, 0x24, 0xe7, 0xe7, 0xe8, 0x27, 0x17, 0xc3, 0x02, 0x06, 0x10, 0x00,
	0x00, 0xff, 0xff, 0x4c, 0x72, 0x47, 0xf1, 0x55, 0x01, 0x00, 0x00,
}
