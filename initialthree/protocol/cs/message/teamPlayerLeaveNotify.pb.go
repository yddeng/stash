// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/teamPlayerLeaveNotify.proto

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

type TeamPlayerLeaveNotifyToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamPlayerLeaveNotifyToS) Reset()         { *m = TeamPlayerLeaveNotifyToS{} }
func (m *TeamPlayerLeaveNotifyToS) String() string { return proto.CompactTextString(m) }
func (*TeamPlayerLeaveNotifyToS) ProtoMessage()    {}
func (*TeamPlayerLeaveNotifyToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_30440b5737ae136f, []int{0}
}

func (m *TeamPlayerLeaveNotifyToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamPlayerLeaveNotifyToS.Unmarshal(m, b)
}
func (m *TeamPlayerLeaveNotifyToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamPlayerLeaveNotifyToS.Marshal(b, m, deterministic)
}
func (m *TeamPlayerLeaveNotifyToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamPlayerLeaveNotifyToS.Merge(m, src)
}
func (m *TeamPlayerLeaveNotifyToS) XXX_Size() int {
	return xxx_messageInfo_TeamPlayerLeaveNotifyToS.Size(m)
}
func (m *TeamPlayerLeaveNotifyToS) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamPlayerLeaveNotifyToS.DiscardUnknown(m)
}

var xxx_messageInfo_TeamPlayerLeaveNotifyToS proto.InternalMessageInfo

type TeamPlayerLeaveNotifyToC struct {
	Player               *TeamPlayer `protobuf:"bytes,1,req,name=player" json:"player,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *TeamPlayerLeaveNotifyToC) Reset()         { *m = TeamPlayerLeaveNotifyToC{} }
func (m *TeamPlayerLeaveNotifyToC) String() string { return proto.CompactTextString(m) }
func (*TeamPlayerLeaveNotifyToC) ProtoMessage()    {}
func (*TeamPlayerLeaveNotifyToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_30440b5737ae136f, []int{1}
}

func (m *TeamPlayerLeaveNotifyToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamPlayerLeaveNotifyToC.Unmarshal(m, b)
}
func (m *TeamPlayerLeaveNotifyToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamPlayerLeaveNotifyToC.Marshal(b, m, deterministic)
}
func (m *TeamPlayerLeaveNotifyToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamPlayerLeaveNotifyToC.Merge(m, src)
}
func (m *TeamPlayerLeaveNotifyToC) XXX_Size() int {
	return xxx_messageInfo_TeamPlayerLeaveNotifyToC.Size(m)
}
func (m *TeamPlayerLeaveNotifyToC) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamPlayerLeaveNotifyToC.DiscardUnknown(m)
}

var xxx_messageInfo_TeamPlayerLeaveNotifyToC proto.InternalMessageInfo

func (m *TeamPlayerLeaveNotifyToC) GetPlayer() *TeamPlayer {
	if m != nil {
		return m.Player
	}
	return nil
}

func init() {
	proto.RegisterType((*TeamPlayerLeaveNotifyToS)(nil), "message.teamPlayerLeaveNotify_toS")
	proto.RegisterType((*TeamPlayerLeaveNotifyToC)(nil), "message.teamPlayerLeaveNotify_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/teamPlayerLeaveNotify.proto", fileDescriptor_30440b5737ae136f)
}

var fileDescriptor_30440b5737ae136f = []byte{
	// 144 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x49, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x49, 0x4d,
	0xcc, 0x0d, 0xc8, 0x49, 0xac, 0x4c, 0x2d, 0xf2, 0x49, 0x4d, 0x2c, 0x4b, 0xf5, 0xcb, 0x2f, 0xc9,
	0x4c, 0xab, 0xd4, 0x03, 0xab, 0x11, 0x62, 0x87, 0x2a, 0x92, 0x92, 0xc6, 0xaa, 0x0d, 0xa2, 0x4a,
	0x49, 0x9a, 0x4b, 0x12, 0xab, 0x21, 0xf1, 0x25, 0xf9, 0xc1, 0x4a, 0x1e, 0xb8, 0x25, 0x9d, 0x85,
	0xb4, 0xb9, 0xd8, 0x0a, 0xc0, 0x12, 0x12, 0x8c, 0x0a, 0x4c, 0x1a, 0xdc, 0x46, 0xc2, 0x7a, 0x50,
	0xe3, 0xf5, 0x42, 0xe0, 0x7a, 0x82, 0xa0, 0x4a, 0x9c, 0x94, 0xa2, 0x14, 0x32, 0xf3, 0x32, 0x4b,
	0x32, 0x13, 0x73, 0x4a, 0x32, 0x8a, 0x52, 0x53, 0x21, 0xee, 0x49, 0xce, 0xcf, 0xd1, 0x4f, 0x2e,
	0x86, 0xb9, 0x0a, 0x10, 0x00, 0x00, 0xff, 0xff, 0x09, 0x34, 0x16, 0xc0, 0xdf, 0x00, 0x00, 0x00,
}
