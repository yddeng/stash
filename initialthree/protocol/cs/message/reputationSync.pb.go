// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/reputationSync.proto

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

type ReputationSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReputationSyncToS) Reset()         { *m = ReputationSyncToS{} }
func (m *ReputationSyncToS) String() string { return proto.CompactTextString(m) }
func (*ReputationSyncToS) ProtoMessage()    {}
func (*ReputationSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_8aadeb9e09e87245, []int{0}
}

func (m *ReputationSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReputationSyncToS.Unmarshal(m, b)
}
func (m *ReputationSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReputationSyncToS.Marshal(b, m, deterministic)
}
func (m *ReputationSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReputationSyncToS.Merge(m, src)
}
func (m *ReputationSyncToS) XXX_Size() int {
	return xxx_messageInfo_ReputationSyncToS.Size(m)
}
func (m *ReputationSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_ReputationSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_ReputationSyncToS proto.InternalMessageInfo

type ReputationSyncToC struct {
	IsAll                *bool                `protobuf:"varint,1,req,name=isAll" json:"isAll,omitempty"`
	Reputations          []*Reputation        `protobuf:"bytes,2,rep,name=reputations" json:"reputations,omitempty"`
	ReputationItems      []*ReputationItem    `protobuf:"bytes,3,rep,name=reputationItems" json:"reputationItems,omitempty"`
	Refresh              []*ReputationRefresh `protobuf:"bytes,4,rep,name=refresh" json:"refresh,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ReputationSyncToC) Reset()         { *m = ReputationSyncToC{} }
func (m *ReputationSyncToC) String() string { return proto.CompactTextString(m) }
func (*ReputationSyncToC) ProtoMessage()    {}
func (*ReputationSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_8aadeb9e09e87245, []int{1}
}

func (m *ReputationSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReputationSyncToC.Unmarshal(m, b)
}
func (m *ReputationSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReputationSyncToC.Marshal(b, m, deterministic)
}
func (m *ReputationSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReputationSyncToC.Merge(m, src)
}
func (m *ReputationSyncToC) XXX_Size() int {
	return xxx_messageInfo_ReputationSyncToC.Size(m)
}
func (m *ReputationSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_ReputationSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_ReputationSyncToC proto.InternalMessageInfo

func (m *ReputationSyncToC) GetIsAll() bool {
	if m != nil && m.IsAll != nil {
		return *m.IsAll
	}
	return false
}

func (m *ReputationSyncToC) GetReputations() []*Reputation {
	if m != nil {
		return m.Reputations
	}
	return nil
}

func (m *ReputationSyncToC) GetReputationItems() []*ReputationItem {
	if m != nil {
		return m.ReputationItems
	}
	return nil
}

func (m *ReputationSyncToC) GetRefresh() []*ReputationRefresh {
	if m != nil {
		return m.Refresh
	}
	return nil
}

type ReputationItem struct {
	Id                   *int32   `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Count                *int32   `protobuf:"varint,2,req,name=count" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReputationItem) Reset()         { *m = ReputationItem{} }
func (m *ReputationItem) String() string { return proto.CompactTextString(m) }
func (*ReputationItem) ProtoMessage()    {}
func (*ReputationItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_8aadeb9e09e87245, []int{2}
}

func (m *ReputationItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReputationItem.Unmarshal(m, b)
}
func (m *ReputationItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReputationItem.Marshal(b, m, deterministic)
}
func (m *ReputationItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReputationItem.Merge(m, src)
}
func (m *ReputationItem) XXX_Size() int {
	return xxx_messageInfo_ReputationItem.Size(m)
}
func (m *ReputationItem) XXX_DiscardUnknown() {
	xxx_messageInfo_ReputationItem.DiscardUnknown(m)
}

var xxx_messageInfo_ReputationItem proto.InternalMessageInfo

func (m *ReputationItem) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *ReputationItem) GetCount() int32 {
	if m != nil && m.Count != nil {
		return *m.Count
	}
	return 0
}

type ReputationRefresh struct {
	Id                   *int32   `protobuf:"varint,1,req,name=id" json:"id,omitempty"`
	Times                *int32   `protobuf:"varint,2,req,name=times" json:"times,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ReputationRefresh) Reset()         { *m = ReputationRefresh{} }
func (m *ReputationRefresh) String() string { return proto.CompactTextString(m) }
func (*ReputationRefresh) ProtoMessage()    {}
func (*ReputationRefresh) Descriptor() ([]byte, []int) {
	return fileDescriptor_8aadeb9e09e87245, []int{3}
}

func (m *ReputationRefresh) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReputationRefresh.Unmarshal(m, b)
}
func (m *ReputationRefresh) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReputationRefresh.Marshal(b, m, deterministic)
}
func (m *ReputationRefresh) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReputationRefresh.Merge(m, src)
}
func (m *ReputationRefresh) XXX_Size() int {
	return xxx_messageInfo_ReputationRefresh.Size(m)
}
func (m *ReputationRefresh) XXX_DiscardUnknown() {
	xxx_messageInfo_ReputationRefresh.DiscardUnknown(m)
}

var xxx_messageInfo_ReputationRefresh proto.InternalMessageInfo

func (m *ReputationRefresh) GetId() int32 {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return 0
}

func (m *ReputationRefresh) GetTimes() int32 {
	if m != nil && m.Times != nil {
		return *m.Times
	}
	return 0
}

type Reputation struct {
	CampType             *int32   `protobuf:"varint,1,opt,name=campType" json:"campType,omitempty"`
	ReputationLevel      *int32   `protobuf:"varint,2,opt,name=reputationLevel" json:"reputationLevel,omitempty"`
	CurrentReputation    *int32   `protobuf:"varint,3,opt,name=currentReputation" json:"currentReputation,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Reputation) Reset()         { *m = Reputation{} }
func (m *Reputation) String() string { return proto.CompactTextString(m) }
func (*Reputation) ProtoMessage()    {}
func (*Reputation) Descriptor() ([]byte, []int) {
	return fileDescriptor_8aadeb9e09e87245, []int{4}
}

func (m *Reputation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Reputation.Unmarshal(m, b)
}
func (m *Reputation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Reputation.Marshal(b, m, deterministic)
}
func (m *Reputation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Reputation.Merge(m, src)
}
func (m *Reputation) XXX_Size() int {
	return xxx_messageInfo_Reputation.Size(m)
}
func (m *Reputation) XXX_DiscardUnknown() {
	xxx_messageInfo_Reputation.DiscardUnknown(m)
}

var xxx_messageInfo_Reputation proto.InternalMessageInfo

func (m *Reputation) GetCampType() int32 {
	if m != nil && m.CampType != nil {
		return *m.CampType
	}
	return 0
}

func (m *Reputation) GetReputationLevel() int32 {
	if m != nil && m.ReputationLevel != nil {
		return *m.ReputationLevel
	}
	return 0
}

func (m *Reputation) GetCurrentReputation() int32 {
	if m != nil && m.CurrentReputation != nil {
		return *m.CurrentReputation
	}
	return 0
}

func init() {
	proto.RegisterType((*ReputationSyncToS)(nil), "message.reputationSync_toS")
	proto.RegisterType((*ReputationSyncToC)(nil), "message.reputationSync_toC")
	proto.RegisterType((*ReputationItem)(nil), "message.reputationItem")
	proto.RegisterType((*ReputationRefresh)(nil), "message.reputationRefresh")
	proto.RegisterType((*Reputation)(nil), "message.reputation")
}

func init() {
	proto.RegisterFile("cs/proto/message/reputationSync.proto", fileDescriptor_8aadeb9e09e87245)
}

var fileDescriptor_8aadeb9e09e87245 = []byte{
	// 295 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0x4f, 0x4b, 0xc3, 0x30,
	0x18, 0xc6, 0x69, 0x66, 0xd9, 0x78, 0x07, 0x93, 0xc5, 0x81, 0x61, 0xa7, 0x12, 0x10, 0x7a, 0x90,
	0x16, 0x44, 0x05, 0x8f, 0xd3, 0x93, 0xe0, 0x29, 0xf3, 0xe4, 0x45, 0x4a, 0xf6, 0xea, 0x02, 0x6d,
	0x53, 0x92, 0x54, 0xd8, 0xcd, 0xcf, 0xe9, 0xa7, 0x91, 0xa5, 0xeb, 0xba, 0x3f, 0x3d, 0x3e, 0xcf,
	0xf3, 0x7b, 0xf2, 0xf2, 0x86, 0x17, 0x6e, 0xa4, 0x4d, 0x2b, 0xa3, 0x9d, 0x4e, 0x0b, 0xb4, 0x36,
	0xfb, 0xc6, 0xd4, 0x60, 0x55, 0xbb, 0xcc, 0x29, 0x5d, 0x2e, 0x37, 0xa5, 0x4c, 0x7c, 0x48, 0x87,
	0xbb, 0x94, 0xcf, 0x80, 0x1e, 0x03, 0x9f, 0x4e, 0x2f, 0xf9, 0x5f, 0xd0, 0x63, 0xbf, 0xd0, 0x19,
	0x84, 0xca, 0x2e, 0xf2, 0x9c, 0x05, 0x11, 0x89, 0x47, 0xa2, 0x11, 0xf4, 0x01, 0xc6, 0x1d, 0x6b,
	0x19, 0x89, 0x06, 0xf1, 0xf8, 0xee, 0x2a, 0xd9, 0x4d, 0x48, 0xba, 0x4c, 0x1c, 0x72, 0x74, 0x01,
	0x97, 0x9d, 0x7c, 0x75, 0x58, 0x58, 0x36, 0xf0, 0xd5, 0xeb, 0x9e, 0xea, 0x36, 0x17, 0xa7, 0x3c,
	0xbd, 0x87, 0xa1, 0xc1, 0x2f, 0x83, 0x76, 0xcd, 0x2e, 0x7c, 0x75, 0xde, 0x37, 0xb5, 0x21, 0x44,
	0x8b, 0xf2, 0x47, 0x98, 0x1c, 0x3f, 0x44, 0x27, 0x40, 0xd4, 0xca, 0x2f, 0x15, 0x0a, 0xa2, 0x56,
	0xdb, 0x3d, 0xa5, 0xae, 0x4b, 0xc7, 0x88, 0xb7, 0x1a, 0xc1, 0x9f, 0x60, 0x7a, 0xf6, 0x6a, 0x5f,
	0xd5, 0xa9, 0x02, 0x6d, 0x5b, 0xf5, 0x82, 0xff, 0x06, 0x00, 0x5d, 0x97, 0xce, 0x61, 0x24, 0xb3,
	0xa2, 0x7a, 0xdf, 0x54, 0xc8, 0x82, 0x28, 0x88, 0x43, 0xb1, 0xd7, 0x34, 0x3e, 0xfc, 0x96, 0x37,
	0xfc, 0xc1, 0x9c, 0x11, 0x8f, 0x9c, 0xda, 0xf4, 0x16, 0xa6, 0xb2, 0x36, 0x06, 0x4b, 0x27, 0xf6,
	0x09, 0x1b, 0x78, 0xf6, 0x3c, 0x78, 0xe6, 0x1f, 0x91, 0x2a, 0x95, 0x53, 0x59, 0xee, 0xd6, 0x06,
	0xb1, 0x39, 0x12, 0xa9, 0xf3, 0x54, 0xda, 0xf6, 0x54, 0xfe, 0x03, 0x00, 0x00, 0xff, 0xff, 0x77,
	0x20, 0x8e, 0x6f, 0x3d, 0x02, 0x00, 0x00,
}
