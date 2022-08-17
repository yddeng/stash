// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/drawCardHistory.proto

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

type DrawCardHistoryToS struct {
	LibID                *int32   `protobuf:"varint,1,opt,name=libID" json:"libID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DrawCardHistoryToS) Reset()         { *m = DrawCardHistoryToS{} }
func (m *DrawCardHistoryToS) String() string { return proto.CompactTextString(m) }
func (*DrawCardHistoryToS) ProtoMessage()    {}
func (*DrawCardHistoryToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a2e14e723d35d2a, []int{0}
}

func (m *DrawCardHistoryToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DrawCardHistoryToS.Unmarshal(m, b)
}
func (m *DrawCardHistoryToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DrawCardHistoryToS.Marshal(b, m, deterministic)
}
func (m *DrawCardHistoryToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DrawCardHistoryToS.Merge(m, src)
}
func (m *DrawCardHistoryToS) XXX_Size() int {
	return xxx_messageInfo_DrawCardHistoryToS.Size(m)
}
func (m *DrawCardHistoryToS) XXX_DiscardUnknown() {
	xxx_messageInfo_DrawCardHistoryToS.DiscardUnknown(m)
}

var xxx_messageInfo_DrawCardHistoryToS proto.InternalMessageInfo

func (m *DrawCardHistoryToS) GetLibID() int32 {
	if m != nil && m.LibID != nil {
		return *m.LibID
	}
	return 0
}

type DrawCardHistoryToC struct {
	History              []*DrawCardAward `protobuf:"bytes,1,rep,name=history" json:"history,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *DrawCardHistoryToC) Reset()         { *m = DrawCardHistoryToC{} }
func (m *DrawCardHistoryToC) String() string { return proto.CompactTextString(m) }
func (*DrawCardHistoryToC) ProtoMessage()    {}
func (*DrawCardHistoryToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_8a2e14e723d35d2a, []int{1}
}

func (m *DrawCardHistoryToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DrawCardHistoryToC.Unmarshal(m, b)
}
func (m *DrawCardHistoryToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DrawCardHistoryToC.Marshal(b, m, deterministic)
}
func (m *DrawCardHistoryToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DrawCardHistoryToC.Merge(m, src)
}
func (m *DrawCardHistoryToC) XXX_Size() int {
	return xxx_messageInfo_DrawCardHistoryToC.Size(m)
}
func (m *DrawCardHistoryToC) XXX_DiscardUnknown() {
	xxx_messageInfo_DrawCardHistoryToC.DiscardUnknown(m)
}

var xxx_messageInfo_DrawCardHistoryToC proto.InternalMessageInfo

func (m *DrawCardHistoryToC) GetHistory() []*DrawCardAward {
	if m != nil {
		return m.History
	}
	return nil
}

func init() {
	proto.RegisterType((*DrawCardHistoryToS)(nil), "message.drawCardHistory_toS")
	proto.RegisterType((*DrawCardHistoryToC)(nil), "message.drawCardHistory_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/drawCardHistory.proto", fileDescriptor_8a2e14e723d35d2a)
}

var fileDescriptor_8a2e14e723d35d2a = []byte{
	// 164 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x4b, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x4f, 0x29, 0x4a,
	0x2c, 0x77, 0x4e, 0x2c, 0x4a, 0xf1, 0xc8, 0x2c, 0x2e, 0xc9, 0x2f, 0xaa, 0xd4, 0x03, 0xcb, 0x0a,
	0xb1, 0x43, 0xa5, 0xa5, 0x94, 0x71, 0x6a, 0x70, 0x29, 0x4a, 0x2c, 0x87, 0xa8, 0x56, 0xd2, 0xe6,
	0x12, 0x46, 0x33, 0x26, 0xbe, 0x24, 0x3f, 0x58, 0x48, 0x84, 0x8b, 0x35, 0x27, 0x33, 0xc9, 0xd3,
	0x45, 0x82, 0x51, 0x81, 0x51, 0x83, 0x35, 0x08, 0xc2, 0x51, 0x72, 0xc7, 0xa6, 0xd8, 0x59, 0xc8,
	0x80, 0x8b, 0x3d, 0x03, 0xc2, 0x95, 0x60, 0x54, 0x60, 0xd6, 0xe0, 0x36, 0x12, 0xd3, 0x83, 0xda,
	0xa8, 0x07, 0x53, 0xee, 0x58, 0x9e, 0x58, 0x94, 0x12, 0x04, 0x53, 0xe6, 0xa4, 0x14, 0xa5, 0x90,
	0x99, 0x97, 0x59, 0x92, 0x99, 0x98, 0x53, 0x92, 0x51, 0x94, 0x9a, 0x0a, 0x71, 0x66, 0x72, 0x7e,
	0x8e, 0x7e, 0x72, 0x31, 0xcc, 0xb1, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xf4, 0x6d, 0xb7, 0xdc,
	0xf0, 0x00, 0x00, 0x00,
}
