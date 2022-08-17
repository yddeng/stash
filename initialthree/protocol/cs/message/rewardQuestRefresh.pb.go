// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/rewardQuestRefresh.proto

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

type RewardQuestRefreshToS struct {
	QuestID              *int32   `protobuf:"varint,1,req,name=questID" json:"questID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RewardQuestRefreshToS) Reset()         { *m = RewardQuestRefreshToS{} }
func (m *RewardQuestRefreshToS) String() string { return proto.CompactTextString(m) }
func (*RewardQuestRefreshToS) ProtoMessage()    {}
func (*RewardQuestRefreshToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5f5efd47b4f54f3, []int{0}
}

func (m *RewardQuestRefreshToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RewardQuestRefreshToS.Unmarshal(m, b)
}
func (m *RewardQuestRefreshToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RewardQuestRefreshToS.Marshal(b, m, deterministic)
}
func (m *RewardQuestRefreshToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardQuestRefreshToS.Merge(m, src)
}
func (m *RewardQuestRefreshToS) XXX_Size() int {
	return xxx_messageInfo_RewardQuestRefreshToS.Size(m)
}
func (m *RewardQuestRefreshToS) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardQuestRefreshToS.DiscardUnknown(m)
}

var xxx_messageInfo_RewardQuestRefreshToS proto.InternalMessageInfo

func (m *RewardQuestRefreshToS) GetQuestID() int32 {
	if m != nil && m.QuestID != nil {
		return *m.QuestID
	}
	return 0
}

type RewardQuestRefreshToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RewardQuestRefreshToC) Reset()         { *m = RewardQuestRefreshToC{} }
func (m *RewardQuestRefreshToC) String() string { return proto.CompactTextString(m) }
func (*RewardQuestRefreshToC) ProtoMessage()    {}
func (*RewardQuestRefreshToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_c5f5efd47b4f54f3, []int{1}
}

func (m *RewardQuestRefreshToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RewardQuestRefreshToC.Unmarshal(m, b)
}
func (m *RewardQuestRefreshToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RewardQuestRefreshToC.Marshal(b, m, deterministic)
}
func (m *RewardQuestRefreshToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardQuestRefreshToC.Merge(m, src)
}
func (m *RewardQuestRefreshToC) XXX_Size() int {
	return xxx_messageInfo_RewardQuestRefreshToC.Size(m)
}
func (m *RewardQuestRefreshToC) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardQuestRefreshToC.DiscardUnknown(m)
}

var xxx_messageInfo_RewardQuestRefreshToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RewardQuestRefreshToS)(nil), "message.rewardQuestRefresh_toS")
	proto.RegisterType((*RewardQuestRefreshToC)(nil), "message.rewardQuestRefresh_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/rewardQuestRefresh.proto", fileDescriptor_c5f5efd47b4f54f3)
}

var fileDescriptor_c5f5efd47b4f54f3 = []byte{
	// 128 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4c, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x4a, 0x2d,
	0x4f, 0x2c, 0x4a, 0x09, 0x2c, 0x4d, 0x2d, 0x2e, 0x09, 0x4a, 0x4d, 0x2b, 0x4a, 0x2d, 0xce, 0xd0,
	0x03, 0x2b, 0x10, 0x62, 0x87, 0xaa, 0x50, 0x32, 0xe2, 0x12, 0xc3, 0x54, 0x14, 0x5f, 0x92, 0x1f,
	0x2c, 0x24, 0xc1, 0xc5, 0x5e, 0x08, 0x12, 0xf3, 0x74, 0x91, 0x60, 0x54, 0x60, 0xd2, 0x60, 0x0d,
	0x82, 0x71, 0x95, 0x24, 0x70, 0xe8, 0x71, 0x76, 0x52, 0x8a, 0x52, 0xc8, 0xcc, 0xcb, 0x2c, 0xc9,
	0x4c, 0xcc, 0x29, 0xc9, 0x28, 0x4a, 0x4d, 0x85, 0xb8, 0x26, 0x39, 0x3f, 0x47, 0x3f, 0xb9, 0x18,
	0xe6, 0x26, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xfd, 0x6c, 0xe3, 0x8a, 0xa6, 0x00, 0x00, 0x00,
}