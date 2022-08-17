// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/rankGetRank.proto

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

type RankGetRankToS struct {
	RankID               *int32   `protobuf:"varint,1,req,name=rankID" json:"rankID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RankGetRankToS) Reset()         { *m = RankGetRankToS{} }
func (m *RankGetRankToS) String() string { return proto.CompactTextString(m) }
func (*RankGetRankToS) ProtoMessage()    {}
func (*RankGetRankToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0b020fca72a1c61, []int{0}
}

func (m *RankGetRankToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RankGetRankToS.Unmarshal(m, b)
}
func (m *RankGetRankToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RankGetRankToS.Marshal(b, m, deterministic)
}
func (m *RankGetRankToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RankGetRankToS.Merge(m, src)
}
func (m *RankGetRankToS) XXX_Size() int {
	return xxx_messageInfo_RankGetRankToS.Size(m)
}
func (m *RankGetRankToS) XXX_DiscardUnknown() {
	xxx_messageInfo_RankGetRankToS.DiscardUnknown(m)
}

var xxx_messageInfo_RankGetRankToS proto.InternalMessageInfo

func (m *RankGetRankToS) GetRankID() int32 {
	if m != nil && m.RankID != nil {
		return *m.RankID
	}
	return 0
}

type RankGetRankToC struct {
	Rank                 *int32   `protobuf:"varint,1,opt,name=rank" json:"rank,omitempty"`
	Total                *int32   `protobuf:"varint,2,opt,name=total" json:"total,omitempty"`
	RankID               *int32   `protobuf:"varint,3,opt,name=rankID" json:"rankID,omitempty"`
	Percent              *int32   `protobuf:"varint,4,opt,name=percent" json:"percent,omitempty"`
	Score                *int32   `protobuf:"varint,5,opt,name=score" json:"score,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RankGetRankToC) Reset()         { *m = RankGetRankToC{} }
func (m *RankGetRankToC) String() string { return proto.CompactTextString(m) }
func (*RankGetRankToC) ProtoMessage()    {}
func (*RankGetRankToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0b020fca72a1c61, []int{1}
}

func (m *RankGetRankToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RankGetRankToC.Unmarshal(m, b)
}
func (m *RankGetRankToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RankGetRankToC.Marshal(b, m, deterministic)
}
func (m *RankGetRankToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RankGetRankToC.Merge(m, src)
}
func (m *RankGetRankToC) XXX_Size() int {
	return xxx_messageInfo_RankGetRankToC.Size(m)
}
func (m *RankGetRankToC) XXX_DiscardUnknown() {
	xxx_messageInfo_RankGetRankToC.DiscardUnknown(m)
}

var xxx_messageInfo_RankGetRankToC proto.InternalMessageInfo

func (m *RankGetRankToC) GetRank() int32 {
	if m != nil && m.Rank != nil {
		return *m.Rank
	}
	return 0
}

func (m *RankGetRankToC) GetTotal() int32 {
	if m != nil && m.Total != nil {
		return *m.Total
	}
	return 0
}

func (m *RankGetRankToC) GetRankID() int32 {
	if m != nil && m.RankID != nil {
		return *m.RankID
	}
	return 0
}

func (m *RankGetRankToC) GetPercent() int32 {
	if m != nil && m.Percent != nil {
		return *m.Percent
	}
	return 0
}

func (m *RankGetRankToC) GetScore() int32 {
	if m != nil && m.Score != nil {
		return *m.Score
	}
	return 0
}

func init() {
	proto.RegisterType((*RankGetRankToS)(nil), "message.rankGetRank_toS")
	proto.RegisterType((*RankGetRankToC)(nil), "message.rankGetRank_toC")
}

func init() { proto.RegisterFile("cs/proto/message/rankGetRank.proto", fileDescriptor_f0b020fca72a1c61) }

var fileDescriptor_f0b020fca72a1c61 = []byte{
	// 183 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x8f, 0xb1, 0xca, 0xc2, 0x30,
	0x14, 0x46, 0x49, 0xff, 0xe6, 0x2f, 0x64, 0x11, 0x82, 0x48, 0xc6, 0x92, 0xa9, 0x2e, 0xe6, 0x1d,
	0x54, 0x10, 0xd7, 0xba, 0xb9, 0x48, 0x08, 0x17, 0x2d, 0xad, 0x49, 0x49, 0xee, 0x1b, 0xf8, 0xe2,
	0xd2, 0xdb, 0x14, 0x8a, 0x5b, 0xce, 0x39, 0xe1, 0x83, 0x2b, 0xb4, 0x4b, 0x66, 0x8c, 0x01, 0x83,
	0x79, 0x43, 0x4a, 0xf6, 0x09, 0x26, 0x5a, 0xdf, 0x5f, 0x00, 0x5b, 0xeb, 0xfb, 0x03, 0x15, 0x59,
	0xe5, 0xa4, 0xf7, 0x62, 0xb3, 0xaa, 0x0f, 0x0c, 0x37, 0xb9, 0x13, 0xff, 0x93, 0xba, 0x9e, 0x15,
	0xab, 0x8b, 0x86, 0xb7, 0x99, 0xf4, 0x87, 0xfd, 0xfe, 0x3d, 0x49, 0x29, 0xca, 0x49, 0x29, 0x56,
	0xb3, 0x86, 0xb7, 0xf4, 0x96, 0x5b, 0xc1, 0x31, 0xa0, 0x1d, 0x54, 0x41, 0x72, 0x86, 0xd5, 0xea,
	0x1f, 0xe9, 0x4c, 0x52, 0x89, 0x6a, 0x84, 0xe8, 0xc0, 0xa3, 0x2a, 0x29, 0x2c, 0x38, 0xed, 0x24,
	0x17, 0x22, 0x28, 0x3e, 0xef, 0x10, 0x1c, 0xf5, 0xbd, 0xee, 0x7c, 0x87, 0x9d, 0x1d, 0xf0, 0x15,
	0x01, 0xe6, 0x4b, 0x5d, 0x18, 0x8c, 0x4b, 0xcb, 0xbd, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc2,
	0x7f, 0x83, 0x79, 0x02, 0x01, 0x00, 0x00,
}
