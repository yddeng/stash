// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/rankGetTopList.proto

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

type RankGetTopListToS struct {
	Version              *int32   `protobuf:"varint,1,req,name=version" json:"version,omitempty"`
	RankID               *int32   `protobuf:"varint,2,req,name=rankID" json:"rankID,omitempty"`
	GetLast              *bool    `protobuf:"varint,3,opt,name=getLast" json:"getLast,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RankGetTopListToS) Reset()         { *m = RankGetTopListToS{} }
func (m *RankGetTopListToS) String() string { return proto.CompactTextString(m) }
func (*RankGetTopListToS) ProtoMessage()    {}
func (*RankGetTopListToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0109f0479f32891, []int{0}
}

func (m *RankGetTopListToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RankGetTopListToS.Unmarshal(m, b)
}
func (m *RankGetTopListToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RankGetTopListToS.Marshal(b, m, deterministic)
}
func (m *RankGetTopListToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RankGetTopListToS.Merge(m, src)
}
func (m *RankGetTopListToS) XXX_Size() int {
	return xxx_messageInfo_RankGetTopListToS.Size(m)
}
func (m *RankGetTopListToS) XXX_DiscardUnknown() {
	xxx_messageInfo_RankGetTopListToS.DiscardUnknown(m)
}

var xxx_messageInfo_RankGetTopListToS proto.InternalMessageInfo

func (m *RankGetTopListToS) GetVersion() int32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *RankGetTopListToS) GetRankID() int32 {
	if m != nil && m.RankID != nil {
		return *m.RankID
	}
	return 0
}

func (m *RankGetTopListToS) GetGetLast() bool {
	if m != nil && m.GetLast != nil {
		return *m.GetLast
	}
	return false
}

type RankGetTopListToC struct {
	Version              *int32          `protobuf:"varint,1,req,name=version" json:"version,omitempty"`
	RankInfo             *RankInfo       `protobuf:"bytes,2,opt,name=rankInfo" json:"rankInfo,omitempty"`
	Roles                []*RankRoleInfo `protobuf:"bytes,3,rep,name=roles" json:"roles,omitempty"`
	Rank                 *int32          `protobuf:"varint,4,opt,name=rank" json:"rank,omitempty"`
	Total                *int32          `protobuf:"varint,5,opt,name=total" json:"total,omitempty"`
	Percent              *int32          `protobuf:"varint,6,opt,name=percent" json:"percent,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *RankGetTopListToC) Reset()         { *m = RankGetTopListToC{} }
func (m *RankGetTopListToC) String() string { return proto.CompactTextString(m) }
func (*RankGetTopListToC) ProtoMessage()    {}
func (*RankGetTopListToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0109f0479f32891, []int{1}
}

func (m *RankGetTopListToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RankGetTopListToC.Unmarshal(m, b)
}
func (m *RankGetTopListToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RankGetTopListToC.Marshal(b, m, deterministic)
}
func (m *RankGetTopListToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RankGetTopListToC.Merge(m, src)
}
func (m *RankGetTopListToC) XXX_Size() int {
	return xxx_messageInfo_RankGetTopListToC.Size(m)
}
func (m *RankGetTopListToC) XXX_DiscardUnknown() {
	xxx_messageInfo_RankGetTopListToC.DiscardUnknown(m)
}

var xxx_messageInfo_RankGetTopListToC proto.InternalMessageInfo

func (m *RankGetTopListToC) GetVersion() int32 {
	if m != nil && m.Version != nil {
		return *m.Version
	}
	return 0
}

func (m *RankGetTopListToC) GetRankInfo() *RankInfo {
	if m != nil {
		return m.RankInfo
	}
	return nil
}

func (m *RankGetTopListToC) GetRoles() []*RankRoleInfo {
	if m != nil {
		return m.Roles
	}
	return nil
}

func (m *RankGetTopListToC) GetRank() int32 {
	if m != nil && m.Rank != nil {
		return *m.Rank
	}
	return 0
}

func (m *RankGetTopListToC) GetTotal() int32 {
	if m != nil && m.Total != nil {
		return *m.Total
	}
	return 0
}

func (m *RankGetTopListToC) GetPercent() int32 {
	if m != nil && m.Percent != nil {
		return *m.Percent
	}
	return 0
}

type RankInfo struct {
	RankID               *int32   `protobuf:"varint,1,opt,name=RankID" json:"RankID,omitempty"`
	BeginTime            *int64   `protobuf:"varint,2,opt,name=BeginTime" json:"BeginTime,omitempty"`
	EndTime              *int64   `protobuf:"varint,3,opt,name=EndTime" json:"EndTime,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RankInfo) Reset()         { *m = RankInfo{} }
func (m *RankInfo) String() string { return proto.CompactTextString(m) }
func (*RankInfo) ProtoMessage()    {}
func (*RankInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0109f0479f32891, []int{2}
}

func (m *RankInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RankInfo.Unmarshal(m, b)
}
func (m *RankInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RankInfo.Marshal(b, m, deterministic)
}
func (m *RankInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RankInfo.Merge(m, src)
}
func (m *RankInfo) XXX_Size() int {
	return xxx_messageInfo_RankInfo.Size(m)
}
func (m *RankInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RankInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RankInfo proto.InternalMessageInfo

func (m *RankInfo) GetRankID() int32 {
	if m != nil && m.RankID != nil {
		return *m.RankID
	}
	return 0
}

func (m *RankInfo) GetBeginTime() int64 {
	if m != nil && m.BeginTime != nil {
		return *m.BeginTime
	}
	return 0
}

func (m *RankInfo) GetEndTime() int64 {
	if m != nil && m.EndTime != nil {
		return *m.EndTime
	}
	return 0
}

type RankRoleInfo struct {
	ID                   *uint64  `protobuf:"varint,1,opt,name=ID" json:"ID,omitempty"`
	Name                 *string  `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Level                *int32   `protobuf:"varint,3,opt,name=level" json:"level,omitempty"`
	Score                *int32   `protobuf:"varint,4,opt,name=score" json:"score,omitempty"`
	CharacterList        []int32  `protobuf:"varint,5,rep,name=characterList" json:"characterList,omitempty"`
	Avatar               *int32   `protobuf:"varint,6,opt,name=avatar" json:"avatar,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RankRoleInfo) Reset()         { *m = RankRoleInfo{} }
func (m *RankRoleInfo) String() string { return proto.CompactTextString(m) }
func (*RankRoleInfo) ProtoMessage()    {}
func (*RankRoleInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_d0109f0479f32891, []int{3}
}

func (m *RankRoleInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RankRoleInfo.Unmarshal(m, b)
}
func (m *RankRoleInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RankRoleInfo.Marshal(b, m, deterministic)
}
func (m *RankRoleInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RankRoleInfo.Merge(m, src)
}
func (m *RankRoleInfo) XXX_Size() int {
	return xxx_messageInfo_RankRoleInfo.Size(m)
}
func (m *RankRoleInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_RankRoleInfo.DiscardUnknown(m)
}

var xxx_messageInfo_RankRoleInfo proto.InternalMessageInfo

func (m *RankRoleInfo) GetID() uint64 {
	if m != nil && m.ID != nil {
		return *m.ID
	}
	return 0
}

func (m *RankRoleInfo) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *RankRoleInfo) GetLevel() int32 {
	if m != nil && m.Level != nil {
		return *m.Level
	}
	return 0
}

func (m *RankRoleInfo) GetScore() int32 {
	if m != nil && m.Score != nil {
		return *m.Score
	}
	return 0
}

func (m *RankRoleInfo) GetCharacterList() []int32 {
	if m != nil {
		return m.CharacterList
	}
	return nil
}

func (m *RankRoleInfo) GetAvatar() int32 {
	if m != nil && m.Avatar != nil {
		return *m.Avatar
	}
	return 0
}

func init() {
	proto.RegisterType((*RankGetTopListToS)(nil), "message.rankGetTopList_toS")
	proto.RegisterType((*RankGetTopListToC)(nil), "message.rankGetTopList_toC")
	proto.RegisterType((*RankInfo)(nil), "message.RankInfo")
	proto.RegisterType((*RankRoleInfo)(nil), "message.RankRoleInfo")
}

func init() {
	proto.RegisterFile("cs/proto/message/rankGetTopList.proto", fileDescriptor_d0109f0479f32891)
}

var fileDescriptor_d0109f0479f32891 = []byte{
	// 360 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcf, 0x4a, 0x33, 0x31,
	0x14, 0xc5, 0x99, 0x99, 0xa6, 0x7f, 0xd2, 0xef, 0x13, 0x0c, 0x2a, 0x59, 0xb8, 0x18, 0x06, 0x85,
	0x01, 0xb1, 0x03, 0x7d, 0x84, 0x5a, 0x91, 0x42, 0x57, 0xb1, 0xab, 0x6e, 0x34, 0xc4, 0x6b, 0x1b,
	0x4c, 0x93, 0x92, 0x84, 0x3e, 0x8d, 0xcf, 0xe4, 0x33, 0x49, 0x32, 0x19, 0xb5, 0x20, 0xee, 0xf2,
	0xbb, 0xe7, 0xe4, 0xe6, 0x9e, 0x24, 0xf8, 0x5a, 0xb8, 0x66, 0x6f, 0x8d, 0x37, 0xcd, 0x0e, 0x9c,
	0xe3, 0x1b, 0x68, 0x2c, 0xd7, 0x6f, 0x0f, 0xe0, 0x57, 0x66, 0xbf, 0x94, 0xce, 0x4f, 0xa2, 0x48,
	0x06, 0x49, 0xad, 0x9e, 0x31, 0x39, 0x36, 0x3c, 0x79, 0xf3, 0x48, 0x28, 0x1e, 0x1c, 0xc0, 0x3a,
	0x69, 0x34, 0xcd, 0xca, 0xbc, 0x46, 0xac, 0x43, 0x72, 0x81, 0xfb, 0xc1, 0xbf, 0x98, 0xd3, 0x3c,
	0x0a, 0x89, 0xc2, 0x8e, 0x0d, 0xf8, 0x25, 0x77, 0x9e, 0x16, 0x65, 0x56, 0x0f, 0x59, 0x87, 0xd5,
	0x47, 0xf6, 0xcb, 0x11, 0x77, 0x7f, 0x1c, 0x71, 0x8b, 0x87, 0xb1, 0xa9, 0x7e, 0x35, 0x34, 0x2f,
	0xb3, 0x7a, 0x3c, 0x3d, 0x9d, 0xa4, 0x71, 0x27, 0x2c, 0x09, 0xec, 0xcb, 0x42, 0x6e, 0x30, 0xb2,
	0x46, 0x81, 0xa3, 0x45, 0x59, 0xd4, 0xe3, 0xe9, 0xf9, 0x91, 0x97, 0x19, 0x05, 0xd1, 0xdf, 0x7a,
	0x08, 0xc1, 0xbd, 0xb0, 0x91, 0xf6, 0xca, 0xac, 0x46, 0x2c, 0xae, 0xc9, 0x19, 0x46, 0xde, 0x78,
	0xae, 0x28, 0x8a, 0xc5, 0x16, 0xc2, 0x7c, 0x7b, 0xb0, 0x02, 0xb4, 0xa7, 0xfd, 0x58, 0xef, 0xb0,
	0x5a, 0xe3, 0x61, 0x37, 0x46, 0xb8, 0x0e, 0xd6, 0x5e, 0x47, 0x16, 0x4d, 0x89, 0xc8, 0x25, 0x1e,
	0xcd, 0x60, 0x23, 0xf5, 0x4a, 0xee, 0x20, 0x86, 0x28, 0xd8, 0x77, 0x21, 0xf4, 0xbe, 0xd7, 0x2f,
	0x51, 0x2b, 0xa2, 0xd6, 0x61, 0xf5, 0x9e, 0xe1, 0x7f, 0x3f, 0xe7, 0x26, 0x27, 0x38, 0x4f, 0xcd,
	0x7b, 0x2c, 0x5f, 0xcc, 0x43, 0x00, 0xcd, 0x53, 0xcf, 0x11, 0x8b, 0xeb, 0x10, 0x40, 0xc1, 0x01,
	0x54, 0x6c, 0x86, 0x58, 0x0b, 0xa1, 0xea, 0x84, 0xb1, 0x90, 0xb2, 0xb6, 0x40, 0xae, 0xf0, 0x7f,
	0xb1, 0xe5, 0x96, 0x0b, 0x0f, 0x36, 0xbc, 0x05, 0x45, 0x65, 0x51, 0x23, 0x76, 0x5c, 0x0c, 0xb1,
	0xf8, 0x81, 0x7b, 0x6e, 0x53, 0xf6, 0x44, 0xb3, 0x6a, 0x5d, 0x4a, 0x2d, 0xbd, 0xe4, 0xca, 0x6f,
	0x2d, 0x40, 0xfb, 0xd3, 0x84, 0x51, 0x8d, 0x70, 0xdd, 0x7f, 0xfb, 0x0c, 0x00, 0x00, 0xff, 0xff,
	0xab, 0x80, 0xac, 0xe8, 0x82, 0x02, 0x00, 0x00,
}