// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/materialDungeonSync.proto

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

type MaterialDungeonSyncToS struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MaterialDungeonSyncToS) Reset()         { *m = MaterialDungeonSyncToS{} }
func (m *MaterialDungeonSyncToS) String() string { return proto.CompactTextString(m) }
func (*MaterialDungeonSyncToS) ProtoMessage()    {}
func (*MaterialDungeonSyncToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_266459aab6d32efa, []int{0}
}

func (m *MaterialDungeonSyncToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MaterialDungeonSyncToS.Unmarshal(m, b)
}
func (m *MaterialDungeonSyncToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MaterialDungeonSyncToS.Marshal(b, m, deterministic)
}
func (m *MaterialDungeonSyncToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MaterialDungeonSyncToS.Merge(m, src)
}
func (m *MaterialDungeonSyncToS) XXX_Size() int {
	return xxx_messageInfo_MaterialDungeonSyncToS.Size(m)
}
func (m *MaterialDungeonSyncToS) XXX_DiscardUnknown() {
	xxx_messageInfo_MaterialDungeonSyncToS.DiscardUnknown(m)
}

var xxx_messageInfo_MaterialDungeonSyncToS proto.InternalMessageInfo

type MaterialDungeonSyncToC struct {
	All                  *bool              `protobuf:"varint,1,opt,name=all" json:"all,omitempty"`
	MaterialDungeons     []*MaterialDungeon `protobuf:"bytes,2,rep,name=materialDungeons" json:"materialDungeons,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *MaterialDungeonSyncToC) Reset()         { *m = MaterialDungeonSyncToC{} }
func (m *MaterialDungeonSyncToC) String() string { return proto.CompactTextString(m) }
func (*MaterialDungeonSyncToC) ProtoMessage()    {}
func (*MaterialDungeonSyncToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_266459aab6d32efa, []int{1}
}

func (m *MaterialDungeonSyncToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MaterialDungeonSyncToC.Unmarshal(m, b)
}
func (m *MaterialDungeonSyncToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MaterialDungeonSyncToC.Marshal(b, m, deterministic)
}
func (m *MaterialDungeonSyncToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MaterialDungeonSyncToC.Merge(m, src)
}
func (m *MaterialDungeonSyncToC) XXX_Size() int {
	return xxx_messageInfo_MaterialDungeonSyncToC.Size(m)
}
func (m *MaterialDungeonSyncToC) XXX_DiscardUnknown() {
	xxx_messageInfo_MaterialDungeonSyncToC.DiscardUnknown(m)
}

var xxx_messageInfo_MaterialDungeonSyncToC proto.InternalMessageInfo

func (m *MaterialDungeonSyncToC) GetAll() bool {
	if m != nil && m.All != nil {
		return *m.All
	}
	return false
}

func (m *MaterialDungeonSyncToC) GetMaterialDungeons() []*MaterialDungeon {
	if m != nil {
		return m.MaterialDungeons
	}
	return nil
}

type MaterialDungeon struct {
	DungeonID            *int32   `protobuf:"varint,1,opt,name=dungeonID" json:"dungeonID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MaterialDungeon) Reset()         { *m = MaterialDungeon{} }
func (m *MaterialDungeon) String() string { return proto.CompactTextString(m) }
func (*MaterialDungeon) ProtoMessage()    {}
func (*MaterialDungeon) Descriptor() ([]byte, []int) {
	return fileDescriptor_266459aab6d32efa, []int{2}
}

func (m *MaterialDungeon) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MaterialDungeon.Unmarshal(m, b)
}
func (m *MaterialDungeon) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MaterialDungeon.Marshal(b, m, deterministic)
}
func (m *MaterialDungeon) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MaterialDungeon.Merge(m, src)
}
func (m *MaterialDungeon) XXX_Size() int {
	return xxx_messageInfo_MaterialDungeon.Size(m)
}
func (m *MaterialDungeon) XXX_DiscardUnknown() {
	xxx_messageInfo_MaterialDungeon.DiscardUnknown(m)
}

var xxx_messageInfo_MaterialDungeon proto.InternalMessageInfo

func (m *MaterialDungeon) GetDungeonID() int32 {
	if m != nil && m.DungeonID != nil {
		return *m.DungeonID
	}
	return 0
}

func init() {
	proto.RegisterType((*MaterialDungeonSyncToS)(nil), "message.materialDungeonSync_toS")
	proto.RegisterType((*MaterialDungeonSyncToC)(nil), "message.materialDungeonSync_toC")
	proto.RegisterType((*MaterialDungeon)(nil), "message.materialDungeon")
}

func init() {
	proto.RegisterFile("cs/proto/message/materialDungeonSync.proto", fileDescriptor_266459aab6d32efa)
}

var fileDescriptor_266459aab6d32efa = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x4a, 0x2e, 0xd6, 0x2f,
	0x28, 0xca, 0x2f, 0xc9, 0xd7, 0xcf, 0x4d, 0x2d, 0x2e, 0x4e, 0x4c, 0x4f, 0xd5, 0xcf, 0x4d, 0x2c,
	0x49, 0x2d, 0xca, 0x4c, 0xcc, 0x71, 0x29, 0xcd, 0x4b, 0x4f, 0xcd, 0xcf, 0x0b, 0xae, 0xcc, 0x4b,
	0xd6, 0x03, 0xab, 0x10, 0x62, 0x87, 0x2a, 0x51, 0x92, 0xe4, 0x12, 0xc7, 0xa2, 0x2a, 0xbe, 0x24,
	0x3f, 0x58, 0xa9, 0x10, 0x97, 0x94, 0xb3, 0x90, 0x00, 0x17, 0x73, 0x62, 0x4e, 0x8e, 0x04, 0xa3,
	0x02, 0xa3, 0x06, 0x47, 0x10, 0x88, 0x29, 0xe4, 0xc2, 0x25, 0x80, 0xa6, 0xb8, 0x58, 0x82, 0x49,
	0x81, 0x59, 0x83, 0xdb, 0x48, 0x42, 0x0f, 0x6a, 0x97, 0x1e, 0x9a, 0x82, 0x20, 0x0c, 0x1d, 0x4a,
	0xfa, 0x5c, 0xfc, 0x68, 0x62, 0x42, 0x32, 0x5c, 0x9c, 0x29, 0x10, 0xa6, 0xa7, 0x0b, 0xd8, 0x42,
	0xd6, 0x20, 0x84, 0x80, 0x93, 0x52, 0x94, 0x42, 0x66, 0x5e, 0x66, 0x49, 0x66, 0x62, 0x4e, 0x49,
	0x46, 0x51, 0x6a, 0x2a, 0xc4, 0xff, 0xc9, 0xf9, 0x39, 0xfa, 0xc9, 0xc5, 0xb0, 0x50, 0x00, 0x04,
	0x00, 0x00, 0xff, 0xff, 0x0d, 0xcb, 0x53, 0x00, 0x18, 0x01, 0x00, 0x00,
}
