// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/levelFightEnd.proto

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

type LevelFightEndToS struct {
	FightID              *int64   `protobuf:"varint,1,req,name=fightID" json:"fightID,omitempty"`
	UseTime              *int32   `protobuf:"varint,2,opt,name=useTime" json:"useTime,omitempty"`
	BossCurHP            *float64 `protobuf:"fixed64,3,opt,name=bossCurHP" json:"bossCurHP,omitempty"`
	BossMaxHP            *float64 `protobuf:"fixed64,4,opt,name=bossMaxHP" json:"bossMaxHP,omitempty"`
	Pass                 *bool    `protobuf:"varint,5,opt,name=pass" json:"pass,omitempty"`
	Stars                []bool   `protobuf:"varint,6,rep,name=stars" json:"stars,omitempty"`
	BeHit                *int32   `protobuf:"varint,7,opt,name=beHit" json:"beHit,omitempty"`
	KillBossCount        *int32   `protobuf:"varint,8,opt,name=killBossCount" json:"killBossCount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LevelFightEndToS) Reset()         { *m = LevelFightEndToS{} }
func (m *LevelFightEndToS) String() string { return proto.CompactTextString(m) }
func (*LevelFightEndToS) ProtoMessage()    {}
func (*LevelFightEndToS) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ed5981b72009814, []int{0}
}

func (m *LevelFightEndToS) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LevelFightEndToS.Unmarshal(m, b)
}
func (m *LevelFightEndToS) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LevelFightEndToS.Marshal(b, m, deterministic)
}
func (m *LevelFightEndToS) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LevelFightEndToS.Merge(m, src)
}
func (m *LevelFightEndToS) XXX_Size() int {
	return xxx_messageInfo_LevelFightEndToS.Size(m)
}
func (m *LevelFightEndToS) XXX_DiscardUnknown() {
	xxx_messageInfo_LevelFightEndToS.DiscardUnknown(m)
}

var xxx_messageInfo_LevelFightEndToS proto.InternalMessageInfo

func (m *LevelFightEndToS) GetFightID() int64 {
	if m != nil && m.FightID != nil {
		return *m.FightID
	}
	return 0
}

func (m *LevelFightEndToS) GetUseTime() int32 {
	if m != nil && m.UseTime != nil {
		return *m.UseTime
	}
	return 0
}

func (m *LevelFightEndToS) GetBossCurHP() float64 {
	if m != nil && m.BossCurHP != nil {
		return *m.BossCurHP
	}
	return 0
}

func (m *LevelFightEndToS) GetBossMaxHP() float64 {
	if m != nil && m.BossMaxHP != nil {
		return *m.BossMaxHP
	}
	return 0
}

func (m *LevelFightEndToS) GetPass() bool {
	if m != nil && m.Pass != nil {
		return *m.Pass
	}
	return false
}

func (m *LevelFightEndToS) GetStars() []bool {
	if m != nil {
		return m.Stars
	}
	return nil
}

func (m *LevelFightEndToS) GetBeHit() int32 {
	if m != nil && m.BeHit != nil {
		return *m.BeHit
	}
	return 0
}

func (m *LevelFightEndToS) GetKillBossCount() int32 {
	if m != nil && m.KillBossCount != nil {
		return *m.KillBossCount
	}
	return 0
}

type LevelFightEndToC struct {
	DungeonID            *int32           `protobuf:"varint,1,req,name=dungeonID" json:"dungeonID,omitempty"`
	Pass                 *bool            `protobuf:"varint,2,opt,name=pass" json:"pass,omitempty"`
	Stars                []bool           `protobuf:"varint,3,rep,name=stars" json:"stars,omitempty"`
	AwardList            []*Award         `protobuf:"bytes,4,rep,name=awardList" json:"awardList,omitempty"`
	UseTime              *int32           `protobuf:"varint,5,opt,name=useTime" json:"useTime,omitempty"`
	ScarsIngrainEnd      *ScarsIngrainEnd `protobuf:"bytes,6,opt,name=scarsIngrainEnd" json:"scarsIngrainEnd,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *LevelFightEndToC) Reset()         { *m = LevelFightEndToC{} }
func (m *LevelFightEndToC) String() string { return proto.CompactTextString(m) }
func (*LevelFightEndToC) ProtoMessage()    {}
func (*LevelFightEndToC) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ed5981b72009814, []int{1}
}

func (m *LevelFightEndToC) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LevelFightEndToC.Unmarshal(m, b)
}
func (m *LevelFightEndToC) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LevelFightEndToC.Marshal(b, m, deterministic)
}
func (m *LevelFightEndToC) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LevelFightEndToC.Merge(m, src)
}
func (m *LevelFightEndToC) XXX_Size() int {
	return xxx_messageInfo_LevelFightEndToC.Size(m)
}
func (m *LevelFightEndToC) XXX_DiscardUnknown() {
	xxx_messageInfo_LevelFightEndToC.DiscardUnknown(m)
}

var xxx_messageInfo_LevelFightEndToC proto.InternalMessageInfo

func (m *LevelFightEndToC) GetDungeonID() int32 {
	if m != nil && m.DungeonID != nil {
		return *m.DungeonID
	}
	return 0
}

func (m *LevelFightEndToC) GetPass() bool {
	if m != nil && m.Pass != nil {
		return *m.Pass
	}
	return false
}

func (m *LevelFightEndToC) GetStars() []bool {
	if m != nil {
		return m.Stars
	}
	return nil
}

func (m *LevelFightEndToC) GetAwardList() []*Award {
	if m != nil {
		return m.AwardList
	}
	return nil
}

func (m *LevelFightEndToC) GetUseTime() int32 {
	if m != nil && m.UseTime != nil {
		return *m.UseTime
	}
	return 0
}

func (m *LevelFightEndToC) GetScarsIngrainEnd() *ScarsIngrainEnd {
	if m != nil {
		return m.ScarsIngrainEnd
	}
	return nil
}

type ScarsIngrainEnd struct {
	HistoryScore         *int32   `protobuf:"varint,1,opt,name=historyScore" json:"historyScore,omitempty"`
	ScoreID              *int32   `protobuf:"varint,2,opt,name=scoreID" json:"scoreID,omitempty"`
	TotalScore           *int32   `protobuf:"varint,3,opt,name=totalScore" json:"totalScore,omitempty"`
	TimeScore            *int32   `protobuf:"varint,4,opt,name=timeScore" json:"timeScore,omitempty"`
	DamageScore          *int32   `protobuf:"varint,5,opt,name=damageScore" json:"damageScore,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ScarsIngrainEnd) Reset()         { *m = ScarsIngrainEnd{} }
func (m *ScarsIngrainEnd) String() string { return proto.CompactTextString(m) }
func (*ScarsIngrainEnd) ProtoMessage()    {}
func (*ScarsIngrainEnd) Descriptor() ([]byte, []int) {
	return fileDescriptor_8ed5981b72009814, []int{2}
}

func (m *ScarsIngrainEnd) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ScarsIngrainEnd.Unmarshal(m, b)
}
func (m *ScarsIngrainEnd) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ScarsIngrainEnd.Marshal(b, m, deterministic)
}
func (m *ScarsIngrainEnd) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ScarsIngrainEnd.Merge(m, src)
}
func (m *ScarsIngrainEnd) XXX_Size() int {
	return xxx_messageInfo_ScarsIngrainEnd.Size(m)
}
func (m *ScarsIngrainEnd) XXX_DiscardUnknown() {
	xxx_messageInfo_ScarsIngrainEnd.DiscardUnknown(m)
}

var xxx_messageInfo_ScarsIngrainEnd proto.InternalMessageInfo

func (m *ScarsIngrainEnd) GetHistoryScore() int32 {
	if m != nil && m.HistoryScore != nil {
		return *m.HistoryScore
	}
	return 0
}

func (m *ScarsIngrainEnd) GetScoreID() int32 {
	if m != nil && m.ScoreID != nil {
		return *m.ScoreID
	}
	return 0
}

func (m *ScarsIngrainEnd) GetTotalScore() int32 {
	if m != nil && m.TotalScore != nil {
		return *m.TotalScore
	}
	return 0
}

func (m *ScarsIngrainEnd) GetTimeScore() int32 {
	if m != nil && m.TimeScore != nil {
		return *m.TimeScore
	}
	return 0
}

func (m *ScarsIngrainEnd) GetDamageScore() int32 {
	if m != nil && m.DamageScore != nil {
		return *m.DamageScore
	}
	return 0
}

func init() {
	proto.RegisterType((*LevelFightEndToS)(nil), "message.levelFightEnd_toS")
	proto.RegisterType((*LevelFightEndToC)(nil), "message.levelFightEnd_toC")
	proto.RegisterType((*ScarsIngrainEnd)(nil), "message.scarsIngrainEnd")
}

func init() {
	proto.RegisterFile("cs/proto/message/levelFightEnd.proto", fileDescriptor_8ed5981b72009814)
}

var fileDescriptor_8ed5981b72009814 = []byte{
	// 403 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcb, 0x6e, 0xd4, 0x30,
	0x14, 0x86, 0xe5, 0x49, 0xd2, 0x99, 0x7a, 0xb8, 0x08, 0x8b, 0x85, 0x85, 0x00, 0x59, 0x51, 0x17,
	0x59, 0xa0, 0x89, 0xd4, 0x37, 0x60, 0xda, 0xa2, 0x19, 0x09, 0xa4, 0xca, 0x65, 0xc5, 0x06, 0xb9,
	0x89, 0xc9, 0x58, 0x24, 0x76, 0xe5, 0xe3, 0x70, 0x79, 0x2a, 0x1e, 0x0b, 0xf1, 0x16, 0xc8, 0x4e,
	0xd2, 0x24, 0xed, 0xec, 0xfc, 0xff, 0xdf, 0x39, 0x3e, 0x17, 0x1b, 0x9f, 0x15, 0x90, 0xdf, 0x59,
	0xe3, 0x4c, 0xde, 0x48, 0x00, 0x51, 0xc9, 0xbc, 0x96, 0x3f, 0x64, 0xfd, 0x41, 0x55, 0x07, 0x77,
	0xa5, 0xcb, 0x4d, 0x60, 0x64, 0xd9, 0xc3, 0x57, 0x6f, 0x1e, 0x85, 0x17, 0xa6, 0x69, 0x8c, 0xee,
	0xe2, 0xd2, 0x7f, 0x08, 0xbf, 0x98, 0xe5, 0x7f, 0x75, 0xe6, 0x86, 0x50, 0xbc, 0xfc, 0xe6, 0xf5,
	0xfe, 0x92, 0x22, 0xb6, 0xc8, 0x22, 0x3e, 0x48, 0x4f, 0x5a, 0x90, 0x9f, 0x55, 0x23, 0xe9, 0x82,
	0xa1, 0x2c, 0xe1, 0x83, 0x24, 0xaf, 0xf1, 0xe9, 0xad, 0x01, 0xb8, 0x68, 0xed, 0xee, 0x9a, 0x46,
	0x0c, 0x65, 0x88, 0x8f, 0xc6, 0x40, 0x3f, 0x89, 0x5f, 0xbb, 0x6b, 0x1a, 0x8f, 0x34, 0x18, 0x84,
	0xe0, 0xf8, 0x4e, 0x00, 0xd0, 0x84, 0xa1, 0x6c, 0xc5, 0xc3, 0x99, 0xbc, 0xc4, 0x09, 0x38, 0x61,
	0x81, 0x9e, 0xb0, 0x28, 0x5b, 0xf1, 0x4e, 0x78, 0xf7, 0x56, 0xee, 0x94, 0xa3, 0xcb, 0x50, 0xbd,
	0x13, 0xe4, 0x0c, 0x3f, 0xfd, 0xae, 0xea, 0x7a, 0xeb, 0xcb, 0x99, 0x56, 0x3b, 0xba, 0x0a, 0x74,
	0x6e, 0xa6, 0x7f, 0x8f, 0xcc, 0x7a, 0xe1, 0x3b, 0x2b, 0x5b, 0x5d, 0x49, 0xa3, 0xfb, 0x69, 0x13,
	0x3e, 0x1a, 0xf7, 0x9d, 0x2d, 0x8e, 0x75, 0x16, 0x4d, 0x3b, 0x7b, 0x87, 0x4f, 0xc5, 0x4f, 0x61,
	0xcb, 0x8f, 0x0a, 0x1c, 0x8d, 0x59, 0x94, 0xad, 0xcf, 0x9f, 0x6d, 0xfa, 0x9d, 0x6f, 0xde, 0x7b,
	0xc2, 0xc7, 0x80, 0xe9, 0x1e, 0x93, 0xf9, 0x1e, 0xb7, 0xf8, 0x39, 0x14, 0xc2, 0xc2, 0x5e, 0x57,
	0x56, 0x28, 0x7d, 0xa5, 0x4b, 0x7a, 0xc2, 0x50, 0xb6, 0x3e, 0xa7, 0xf7, 0xb7, 0x3d, 0xe0, 0xfc,
	0x61, 0x42, 0xfa, 0x07, 0x3d, 0xba, 0x84, 0xa4, 0xf8, 0xc9, 0x41, 0x81, 0x33, 0xf6, 0xf7, 0x4d,
	0x61, 0xac, 0xa4, 0x28, 0x94, 0x9d, 0x79, 0xbe, 0x2b, 0xf0, 0x87, 0xfd, 0xe5, 0xf0, 0xba, 0xbd,
	0x24, 0x6f, 0x31, 0x76, 0xc6, 0x89, 0xba, 0xcb, 0x8d, 0x02, 0x9c, 0x38, 0x7e, 0x8b, 0x4e, 0x35,
	0xb2, 0xc3, 0x71, 0xc0, 0xa3, 0x41, 0x18, 0x5e, 0x97, 0xa2, 0x11, 0x55, 0xcf, 0xbb, 0x89, 0xa7,
	0xd6, 0x36, 0xfd, 0xc2, 0x94, 0x56, 0x4e, 0x89, 0xda, 0x1d, 0xac, 0x94, 0xdd, 0x97, 0x2d, 0x4c,
	0x9d, 0x17, 0x30, 0x7c, 0xdc, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x17, 0x0f, 0x96, 0xfa,
	0x02, 0x00, 0x00,
}
