// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/levelFight.proto

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

type LevelFightToS struct {
	DungeonID            *int32         `protobuf:"varint,1,req,name=dungeonID" json:"dungeonID,omitempty"`
	CharacterTeam        *CharacterTeam `protobuf:"bytes,2,req,name=characterTeam" json:"characterTeam,omitempty"`
	Multiple             *int32         `protobuf:"varint,3,opt,name=multiple" json:"multiple,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *LevelFightToS) Reset()         { *m = LevelFightToS{} }
func (m *LevelFightToS) String() string { return proto.CompactTextString(m) }
func (*LevelFightToS) ProtoMessage()    {}
func (*LevelFightToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_2725419b2b1a22e3, []int{0}
}

func (m *LevelFightToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LevelFightToS.Unmarshal(m, b)
}
func (m *LevelFightToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LevelFightToS.Marshal(b, m, deterministic)
}
func (m *LevelFightToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LevelFightToS.Merge(m, src)
}
func (m *LevelFightToS) XXX_Size() int {
	return xxx_messageInfo_LevelFightToS.Size(m)
}
func (m *LevelFightToS) XXX_DiscardUnknown() {
	xxx_messageInfo_LevelFightToS.DiscardUnknown(m)
}

var xxx_messageInfo_LevelFightToS proto.InternalMessageInfo

func (m *LevelFightToS) GetDungeonID() int32 {
	if m != nil && m.DungeonID != nil {
		return *m.DungeonID
	}
	return 0
}

func (m *LevelFightToS) GetCharacterTeam() *CharacterTeam {
	if m != nil {
		return m.CharacterTeam
	}
	return nil
}

func (m *LevelFightToS) GetMultiple() int32 {
	if m != nil && m.Multiple != nil {
		return *m.Multiple
	}
	return 0
}

type LevelFightToC struct {
	FightID              *int64             `protobuf:"varint,1,req,name=fightID" json:"fightID,omitempty"`
	CharacterTeam        *CharacterTeam     `protobuf:"bytes,2,req,name=characterTeam" json:"characterTeam,omitempty"`
	BattleAttrSlice      []*BattleAttrSlice `protobuf:"bytes,3,rep,name=battleAttrSlice" json:"battleAttrSlice,omitempty"`
	StartTime            *int64             `protobuf:"varint,4,req,name=startTime" json:"startTime,omitempty"`
	ComboSkill           *int32             `protobuf:"varint,5,opt,name=comboSkill" json:"comboSkill,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *LevelFightToC) Reset()         { *m = LevelFightToC{} }
func (m *LevelFightToC) String() string { return proto.CompactTextString(m) }
func (*LevelFightToC) ProtoMessage()    {}
func (*LevelFightToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_2725419b2b1a22e3, []int{1}
}

func (m *LevelFightToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LevelFightToC.Unmarshal(m, b)
}
func (m *LevelFightToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LevelFightToC.Marshal(b, m, deterministic)
}
func (m *LevelFightToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LevelFightToC.Merge(m, src)
}
func (m *LevelFightToC) XXX_Size() int {
	return xxx_messageInfo_LevelFightToC.Size(m)
}
func (m *LevelFightToC) XXX_DiscardUnknown() {
	xxx_messageInfo_LevelFightToC.DiscardUnknown(m)
}

var xxx_messageInfo_LevelFightToC proto.InternalMessageInfo

func (m *LevelFightToC) GetFightID() int64 {
	if m != nil && m.FightID != nil {
		return *m.FightID
	}
	return 0
}

func (m *LevelFightToC) GetCharacterTeam() *CharacterTeam {
	if m != nil {
		return m.CharacterTeam
	}
	return nil
}

func (m *LevelFightToC) GetBattleAttrSlice() []*BattleAttrSlice {
	if m != nil {
		return m.BattleAttrSlice
	}
	return nil
}

func (m *LevelFightToC) GetStartTime() int64 {
	if m != nil && m.StartTime != nil {
		return *m.StartTime
	}
	return 0
}

func (m *LevelFightToC) GetComboSkill() int32 {
	if m != nil && m.ComboSkill != nil {
		return *m.ComboSkill
	}
	return 0
}

func init() {
	proto.RegisterType((*LevelFightToS)(nil), "message.levelFight_toS")
	proto.RegisterType((*LevelFightToC)(nil), "message.levelFight_toC")
}

func init() { proto.RegisterFile("cs/proto/message/levelFight.proto", fileDescriptor_2725419b2b1a22e3) }

var fileDescriptor_2725419b2b1a22e3 = []byte{
	// 277 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x51, 0xcd, 0x4a, 0xc3, 0x40,
	0x10, 0x26, 0x89, 0xa1, 0xba, 0x45, 0x85, 0x3d, 0xc8, 0x12, 0x44, 0x62, 0x50, 0xc8, 0x29, 0x81,
	0x9e, 0xbd, 0x18, 0x45, 0xf0, 0x9a, 0xf4, 0xe4, 0x45, 0xb6, 0xeb, 0x98, 0x2c, 0x6e, 0xb2, 0x25,
	0x3b, 0x15, 0x7c, 0x03, 0x9f, 0xd5, 0xa7, 0x90, 0xa4, 0x69, 0x7e, 0xda, 0x63, 0x8f, 0xdf, 0xcf,
	0xf0, 0x7d, 0x33, 0x43, 0x6e, 0x85, 0x89, 0xd7, 0xb5, 0x46, 0x1d, 0x97, 0x60, 0x0c, 0xcf, 0x21,
	0x56, 0xf0, 0x0d, 0xea, 0x45, 0xe6, 0x05, 0x46, 0xad, 0x40, 0x67, 0x9d, 0xe2, 0xdd, 0x1d, 0x78,
	0x45, 0xc1, 0x6b, 0x2e, 0x10, 0xea, 0xec, 0xa7, 0x12, 0x5b, 0xbb, 0x77, 0x7f, 0xe0, 0x5a, 0x71,
	0x44, 0x05, 0x8f, 0x88, 0x23, 0x5b, 0xf0, 0x6b, 0x91, 0x8b, 0x21, 0xea, 0x1d, 0x75, 0x46, 0xaf,
	0xc9, 0xd9, 0xc7, 0xa6, 0xca, 0x41, 0x57, 0xaf, 0xcf, 0xcc, 0xf2, 0xed, 0xd0, 0x4d, 0x07, 0x82,
	0x3e, 0x90, 0xf3, 0x3e, 0x6e, 0x09, 0xbc, 0x64, 0xb6, 0x6f, 0x87, 0xf3, 0xc5, 0x55, 0xd4, 0xc5,
	0x44, 0x13, 0x35, 0x9d, 0x9a, 0xa9, 0x47, 0x4e, 0xcb, 0x8d, 0x42, 0xb9, 0x56, 0xc0, 0x1c, 0xdf,
	0x0a, 0xdd, 0xb4, 0xc7, 0xc1, 0xdf, 0x7e, 0x95, 0x27, 0xca, 0xc8, 0xec, 0xb3, 0x01, 0x5d, 0x11,
	0x27, 0xdd, 0xc1, 0x23, 0x6b, 0x24, 0xe4, 0x72, 0x74, 0x0d, 0x25, 0x45, 0xd3, 0xc6, 0x09, 0xe7,
	0x0b, 0xd6, 0xcf, 0x27, 0x53, 0x3d, 0xdd, 0x1f, 0x68, 0xce, 0x64, 0x90, 0xd7, 0xb8, 0x94, 0x25,
	0xb0, 0x93, 0xb6, 0xdd, 0x40, 0xd0, 0x1b, 0x42, 0x84, 0x2e, 0x57, 0x3a, 0xfb, 0x92, 0x4a, 0x31,
	0xb7, 0x5d, 0x75, 0xc4, 0x24, 0xc1, 0x9b, 0x2f, 0x2b, 0x89, 0x92, 0x2b, 0x2c, 0x6a, 0x80, 0xed,
	0xab, 0x84, 0x56, 0xb1, 0x30, 0xbb, 0x87, 0xfd, 0x07, 0x00, 0x00, 0xff, 0xff, 0xa7, 0xbe, 0x46,
	0xe5, 0x15, 0x02, 0x00, 0x00,
}
