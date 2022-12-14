// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss/proto/ssmessage/reportStatus.proto

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

type ReportStatus struct {
	ServerIds            []int32  `protobuf:"varint,1,rep,name=serverIds" json:"serverIds,omitempty"`
	ServerAddr           *string  `protobuf:"bytes,2,req,name=serverAddr" json:"serverAddr,omitempty"`
	PlayerNum            *int32   `protobuf:"varint,3,req,name=playerNum" json:"playerNum,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReportStatus) Reset()         { *m = ReportStatus{} }
func (m *ReportStatus) String() string { return proto.CompactTextString(m) }
func (*ReportStatus) ProtoMessage()    {}
func (*ReportStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_3238d799f3cca1b7, []int{0}
}

func (m *ReportStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReportStatus.Unmarshal(m, b)
}
func (m *ReportStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReportStatus.Marshal(b, m, deterministic)
}
func (m *ReportStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReportStatus.Merge(m, src)
}
func (m *ReportStatus) XXX_Size() int {
	return xxx_messageInfo_ReportStatus.Size(m)
}
func (m *ReportStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ReportStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ReportStatus proto.InternalMessageInfo

func (m *ReportStatus) GetServerIds() []int32 {
	if m != nil {
		return m.ServerIds
	}
	return nil
}

func (m *ReportStatus) GetServerAddr() string {
	if m != nil && m.ServerAddr != nil {
		return *m.ServerAddr
	}
	return ""
}

func (m *ReportStatus) GetPlayerNum() int32 {
	if m != nil && m.PlayerNum != nil {
		return *m.PlayerNum
	}
	return 0
}

func init() {
	proto.RegisterType((*ReportStatus)(nil), "ssmessage.reportStatus")
}

func init() {
	proto.RegisterFile("ss/proto/ssmessage/reportStatus.proto", fileDescriptor_3238d799f3cca1b7)
}

var fileDescriptor_3238d799f3cca1b7 = []byte{
	// 156 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2d, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0x2e, 0xce, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f,
	0x4a, 0x2d, 0xc8, 0x2f, 0x2a, 0x09, 0x2e, 0x49, 0x2c, 0x29, 0x2d, 0xd6, 0x03, 0x4b, 0x0a, 0x71,
	0xc2, 0x65, 0x95, 0xb2, 0xb8, 0x78, 0x90, 0x15, 0x08, 0xc9, 0x70, 0x71, 0x16, 0xa7, 0x16, 0x95,
	0xa5, 0x16, 0x79, 0xa6, 0x14, 0x4b, 0x30, 0x2a, 0x30, 0x6b, 0xb0, 0x06, 0x21, 0x04, 0x84, 0xe4,
	0xb8, 0xb8, 0x20, 0x1c, 0xc7, 0x94, 0x94, 0x22, 0x09, 0x26, 0x05, 0x26, 0x0d, 0xce, 0x20, 0x24,
	0x11, 0x90, 0xee, 0x82, 0x9c, 0xc4, 0xca, 0xd4, 0x22, 0xbf, 0xd2, 0x5c, 0x09, 0x66, 0x05, 0x26,
	0x90, 0x6e, 0xb8, 0x80, 0x93, 0x4a, 0x94, 0x52, 0x66, 0x5e, 0x66, 0x49, 0x66, 0x62, 0x4e, 0x49,
	0x46, 0x51, 0x6a, 0x2a, 0xc4, 0xa5, 0xc9, 0xf9, 0x39, 0xfa, 0xc5, 0xc5, 0x08, 0xf7, 0x02, 0x02,
	0x00, 0x00, 0xff, 0xff, 0x88, 0xde, 0x8a, 0xe5, 0xc4, 0x00, 0x00, 0x00,
}
