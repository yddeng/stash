// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/rewardQuestComplete.proto

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

type RewardQuestCompleteToS struct {
	QuestID              *int32   `protobuf:"varint,1,req,name=questID" json:"questID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RewardQuestCompleteToS) Reset()         { *m = RewardQuestCompleteToS{} }
func (m *RewardQuestCompleteToS) String() string { return proto.CompactTextString(m) }
func (*RewardQuestCompleteToS) ProtoMessage()    {}
func (*RewardQuestCompleteToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_44233fc308c58c45, []int{0}
}

func (m *RewardQuestCompleteToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RewardQuestCompleteToS.Unmarshal(m, b)
}
func (m *RewardQuestCompleteToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RewardQuestCompleteToS.Marshal(b, m, deterministic)
}
func (m *RewardQuestCompleteToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardQuestCompleteToS.Merge(m, src)
}
func (m *RewardQuestCompleteToS) XXX_Size() int {
	return xxx_messageInfo_RewardQuestCompleteToS.Size(m)
}
func (m *RewardQuestCompleteToS) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardQuestCompleteToS.DiscardUnknown(m)
}

var xxx_messageInfo_RewardQuestCompleteToS proto.InternalMessageInfo

func (m *RewardQuestCompleteToS) GetQuestID() int32 {
	if m != nil && m.QuestID != nil {
		return *m.QuestID
	}
	return 0
}

type RewardQuestCompleteToC struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RewardQuestCompleteToC) Reset()         { *m = RewardQuestCompleteToC{} }
func (m *RewardQuestCompleteToC) String() string { return proto.CompactTextString(m) }
func (*RewardQuestCompleteToC) ProtoMessage()    {}
func (*RewardQuestCompleteToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_44233fc308c58c45, []int{1}
}

func (m *RewardQuestCompleteToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RewardQuestCompleteToC.Unmarshal(m, b)
}
func (m *RewardQuestCompleteToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RewardQuestCompleteToC.Marshal(b, m, deterministic)
}
func (m *RewardQuestCompleteToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RewardQuestCompleteToC.Merge(m, src)
}
func (m *RewardQuestCompleteToC) XXX_Size() int {
	return xxx_messageInfo_RewardQuestCompleteToC.Size(m)
}
func (m *RewardQuestCompleteToC) XXX_DiscardUnknown() {
	xxx_messageInfo_RewardQuestCompleteToC.DiscardUnknown(m)
}

var xxx_messageInfo_RewardQuestCompleteToC proto.InternalMessageInfo

func init() {
	proto.RegisterType((*RewardQuestCompleteToS)(nil), "message.rewardQuestComplete_toS")
	proto.RegisterType((*RewardQuestCompleteToC)(nil), "message.rewardQuestComplete_toC")
}

func init() {
	proto.RegisterFile("cs/proto/message/rewardQuestComplete.proto", fileDescriptor_44233fc308c58c45)
}

var fileDescriptor_44233fc308c58c45 = []byte{
	// 129 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4a, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0x2f, 0x4a, 0x2d,
	0x4f, 0x2c, 0x4a, 0x09, 0x2c, 0x4d, 0x2d, 0x2e, 0x71, 0xce, 0xcf, 0x2d, 0xc8, 0x49, 0x2d, 0x49,
	0xd5, 0x03, 0xab, 0x10, 0x62, 0x87, 0x2a, 0x51, 0x32, 0xe6, 0x12, 0xc7, 0xa2, 0x2a, 0xbe, 0x24,
	0x3f, 0x58, 0x48, 0x82, 0x8b, 0xbd, 0x10, 0x24, 0xe8, 0xe9, 0x22, 0xc1, 0xa8, 0xc0, 0xa4, 0xc1,
	0x1a, 0x04, 0xe3, 0x2a, 0x49, 0xe2, 0xd2, 0xe4, 0xec, 0xa4, 0x14, 0xa5, 0x90, 0x99, 0x97, 0x59,
	0x92, 0x99, 0x98, 0x53, 0x92, 0x51, 0x94, 0x9a, 0x0a, 0x71, 0x50, 0x72, 0x7e, 0x8e, 0x7e, 0x72,
	0x31, 0xcc, 0x59, 0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0xcb, 0x6c, 0x76, 0x9b, 0xa9, 0x00, 0x00,
	0x00,
}
