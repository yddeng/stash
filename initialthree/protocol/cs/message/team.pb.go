// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cs/proto/message/team.proto

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

type TeamStatus int32

const (
	TeamStatus_Unknown       TeamStatus = 0
	TeamStatus_Standby       TeamStatus = 1
	TeamStatus_Battle_Verify TeamStatus = 2
	TeamStatus_Battle_vote   TeamStatus = 3
	TeamStatus_Battle        TeamStatus = 4
)

var TeamStatus_name = map[int32]string{
	0: "Unknown",
	1: "Standby",
	2: "Battle_Verify",
	3: "Battle_vote",
	4: "Battle",
}

var TeamStatus_value = map[string]int32{
	"Unknown":       0,
	"Standby":       1,
	"Battle_Verify": 2,
	"Battle_vote":   3,
	"Battle":        4,
}

func (x TeamStatus) Enum() *TeamStatus {
	p := new(TeamStatus)
	*p = x
	return p
}

func (x TeamStatus) String() string {
	return proto.EnumName(TeamStatus_name, int32(x))
}

func (x *TeamStatus) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(TeamStatus_value, data, "TeamStatus")
	if err != nil {
		return err
	}
	*x = TeamStatus(value)
	return nil
}

func (TeamStatus) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_310eecdd299f0a1d, []int{0}
}

type Team struct {
	TeamID               *uint32       `protobuf:"varint,1,req,name=TeamID" json:"TeamID,omitempty"`
	TeamName             *string       `protobuf:"bytes,2,req,name=TeamName" json:"TeamName,omitempty"`
	Status               *TeamStatus   `protobuf:"varint,3,req,name=status,enum=message.TeamStatus" json:"status,omitempty"`
	Header               *uint64       `protobuf:"varint,4,req,name=Header" json:"Header,omitempty"`
	Players              []*TeamPlayer `protobuf:"bytes,5,rep,name=players" json:"players,omitempty"`
	Target               *TeamTarget   `protobuf:"bytes,6,opt,name=target" json:"target,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Team) Reset()         { *m = Team{} }
func (m *Team) String() string { return proto.CompactTextString(m) }
func (*Team) ProtoMessage()    {}
func (*Team) Descriptor() ([]byte, []int) {
	return fileDescriptor_310eecdd299f0a1d, []int{0}
}

func (m *Team) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Team.Unmarshal(m, b)
}
func (m *Team) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Team.Marshal(b, m, deterministic)
}
func (m *Team) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Team.Merge(m, src)
}
func (m *Team) XXX_Size() int {
	return xxx_messageInfo_Team.Size(m)
}
func (m *Team) XXX_DiscardUnknown() {
	xxx_messageInfo_Team.DiscardUnknown(m)
}

var xxx_messageInfo_Team proto.InternalMessageInfo

func (m *Team) GetTeamID() uint32 {
	if m != nil && m.TeamID != nil {
		return *m.TeamID
	}
	return 0
}

func (m *Team) GetTeamName() string {
	if m != nil && m.TeamName != nil {
		return *m.TeamName
	}
	return ""
}

func (m *Team) GetStatus() TeamStatus {
	if m != nil && m.Status != nil {
		return *m.Status
	}
	return TeamStatus_Unknown
}

func (m *Team) GetHeader() uint64 {
	if m != nil && m.Header != nil {
		return *m.Header
	}
	return 0
}

func (m *Team) GetPlayers() []*TeamPlayer {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *Team) GetTarget() *TeamTarget {
	if m != nil {
		return m.Target
	}
	return nil
}

type TeamPlayer struct {
	UserID               *string  `protobuf:"bytes,1,req,name=userID" json:"userID,omitempty"`
	PlayerID             *uint64  `protobuf:"varint,2,req,name=playerID" json:"playerID,omitempty"`
	CharacterID          *int32   `protobuf:"varint,3,req,name=characterID" json:"characterID,omitempty"`
	PLevel               *int32   `protobuf:"varint,4,req,name=pLevel" json:"pLevel,omitempty"`
	CombatPower          *int32   `protobuf:"varint,5,req,name=combatPower" json:"combatPower,omitempty"`
	Name                 *string  `protobuf:"bytes,6,req,name=name" json:"name,omitempty"`
	Portrait             *int32   `protobuf:"varint,7,req,name=portrait" json:"portrait,omitempty"`
	OnLine               *bool    `protobuf:"varint,8,req,name=onLine" json:"onLine,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamPlayer) Reset()         { *m = TeamPlayer{} }
func (m *TeamPlayer) String() string { return proto.CompactTextString(m) }
func (*TeamPlayer) ProtoMessage()    {}
func (*TeamPlayer) Descriptor() ([]byte, []int) {
	return fileDescriptor_310eecdd299f0a1d, []int{1}
}

func (m *TeamPlayer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamPlayer.Unmarshal(m, b)
}
func (m *TeamPlayer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamPlayer.Marshal(b, m, deterministic)
}
func (m *TeamPlayer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamPlayer.Merge(m, src)
}
func (m *TeamPlayer) XXX_Size() int {
	return xxx_messageInfo_TeamPlayer.Size(m)
}
func (m *TeamPlayer) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamPlayer.DiscardUnknown(m)
}

var xxx_messageInfo_TeamPlayer proto.InternalMessageInfo

func (m *TeamPlayer) GetUserID() string {
	if m != nil && m.UserID != nil {
		return *m.UserID
	}
	return ""
}

func (m *TeamPlayer) GetPlayerID() uint64 {
	if m != nil && m.PlayerID != nil {
		return *m.PlayerID
	}
	return 0
}

func (m *TeamPlayer) GetCharacterID() int32 {
	if m != nil && m.CharacterID != nil {
		return *m.CharacterID
	}
	return 0
}

func (m *TeamPlayer) GetPLevel() int32 {
	if m != nil && m.PLevel != nil {
		return *m.PLevel
	}
	return 0
}

func (m *TeamPlayer) GetCombatPower() int32 {
	if m != nil && m.CombatPower != nil {
		return *m.CombatPower
	}
	return 0
}

func (m *TeamPlayer) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *TeamPlayer) GetPortrait() int32 {
	if m != nil && m.Portrait != nil {
		return *m.Portrait
	}
	return 0
}

func (m *TeamPlayer) GetOnLine() bool {
	if m != nil && m.OnLine != nil {
		return *m.OnLine
	}
	return false
}

type TeamTarget struct {
	LevelID              *int32   `protobuf:"varint,1,req,name=levelID" json:"levelID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamTarget) Reset()         { *m = TeamTarget{} }
func (m *TeamTarget) String() string { return proto.CompactTextString(m) }
func (*TeamTarget) ProtoMessage()    {}
func (*TeamTarget) Descriptor() ([]byte, []int) {
	return fileDescriptor_310eecdd299f0a1d, []int{2}
}

func (m *TeamTarget) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamTarget.Unmarshal(m, b)
}
func (m *TeamTarget) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamTarget.Marshal(b, m, deterministic)
}
func (m *TeamTarget) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamTarget.Merge(m, src)
}
func (m *TeamTarget) XXX_Size() int {
	return xxx_messageInfo_TeamTarget.Size(m)
}
func (m *TeamTarget) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamTarget.DiscardUnknown(m)
}

var xxx_messageInfo_TeamTarget proto.InternalMessageInfo

func (m *TeamTarget) GetLevelID() int32 {
	if m != nil && m.LevelID != nil {
		return *m.LevelID
	}
	return 0
}

func init() {
	proto.RegisterEnum("message.TeamStatus", TeamStatus_name, TeamStatus_value)
	proto.RegisterType((*Team)(nil), "message.Team")
	proto.RegisterType((*TeamPlayer)(nil), "message.TeamPlayer")
	proto.RegisterType((*TeamTarget)(nil), "message.TeamTarget")
}

func init() { proto.RegisterFile("cs/proto/message/team.proto", fileDescriptor_310eecdd299f0a1d) }

var fileDescriptor_310eecdd299f0a1d = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x52, 0x4d, 0x6f, 0xd4, 0x30,
	0x10, 0x25, 0xd9, 0x7c, 0x6c, 0x27, 0x2a, 0x04, 0x23, 0x21, 0x0b, 0x2e, 0x56, 0x0e, 0x28, 0x02,
	0xb1, 0x91, 0xf8, 0x09, 0x55, 0x0f, 0x54, 0xaa, 0x50, 0xe5, 0x96, 0x1e, 0xb8, 0x20, 0x37, 0x1d,
	0xda, 0x88, 0xc4, 0x5e, 0xd9, 0xd3, 0x5d, 0xed, 0x6f, 0xe5, 0xca, 0x0f, 0x41, 0x76, 0xbc, 0xbb,
	0x08, 0x4e, 0xf1, 0x7b, 0xf3, 0xe6, 0xcd, 0xbc, 0xd8, 0xf0, 0xb6, 0x77, 0xdd, 0xda, 0x1a, 0x32,
	0xdd, 0x84, 0xce, 0xa9, 0x07, 0xec, 0x08, 0xd5, 0xb4, 0x0a, 0x14, 0x2b, 0x23, 0xd7, 0xfc, 0x4a,
	0x20, 0xbb, 0x41, 0x35, 0xb1, 0xd7, 0x50, 0xf8, 0xef, 0xc5, 0x39, 0x4f, 0x44, 0xda, 0x9e, 0xca,
	0x88, 0xd8, 0x1b, 0x58, 0xfa, 0xd3, 0x17, 0x35, 0x21, 0x4f, 0x45, 0xda, 0x9e, 0xc8, 0x03, 0x66,
	0x1f, 0xa0, 0x70, 0xa4, 0xe8, 0xc9, 0xf1, 0x85, 0x48, 0xdb, 0xe7, 0x9f, 0x5e, 0xad, 0xa2, 0xed,
	0xca, 0x4b, 0xae, 0x43, 0x49, 0x46, 0x89, 0x1f, 0xf0, 0x19, 0xd5, 0x3d, 0x5a, 0x9e, 0x89, 0xb4,
	0xcd, 0x64, 0x44, 0xec, 0x23, 0x94, 0xeb, 0x51, 0xed, 0xd0, 0x3a, 0x9e, 0x8b, 0x45, 0x5b, 0xfd,
	0xe3, 0x72, 0x15, 0x6a, 0x72, 0xaf, 0xf1, 0x33, 0x49, 0xd9, 0x07, 0x24, 0x5e, 0x88, 0xe4, 0x3f,
	0xf5, 0x4d, 0x28, 0xc9, 0x28, 0x69, 0x7e, 0x27, 0x00, 0x47, 0x13, 0xbf, 0xc2, 0x93, 0x43, 0x1b,
	0x33, 0x9e, 0xc8, 0x88, 0x7c, 0xc6, 0xd9, 0xfe, 0xe2, 0x3c, 0x64, 0xcc, 0xe4, 0x01, 0x33, 0x01,
	0x55, 0xff, 0xa8, 0xac, 0xea, 0x29, 0x94, 0x7d, 0xd0, 0x5c, 0xfe, 0x4d, 0x79, 0xd7, 0xf5, 0x25,
	0x6e, 0x70, 0x0c, 0xc1, 0x72, 0x19, 0x51, 0xe8, 0x34, 0xd3, 0x9d, 0xa2, 0x2b, 0xb3, 0x45, 0xcb,
	0xf3, 0xd8, 0x79, 0xa4, 0x18, 0x83, 0x4c, 0xfb, 0xff, 0x5a, 0x84, 0x6d, 0xc2, 0x39, 0xec, 0x62,
	0x2c, 0x59, 0x35, 0x10, 0x2f, 0x43, 0xcb, 0x01, 0xfb, 0x49, 0x46, 0x5f, 0x0e, 0x1a, 0xf9, 0x52,
	0xa4, 0xed, 0x52, 0x46, 0xd4, 0xbc, 0x9b, 0x53, 0xce, 0xe1, 0x19, 0x87, 0x72, 0xf4, 0x0b, 0xc4,
	0x98, 0xb9, 0xdc, 0xc3, 0xf7, 0xb7, 0xb3, 0x6e, 0xbe, 0x18, 0x56, 0x41, 0xf9, 0x55, 0xff, 0xd4,
	0x66, 0xab, 0xeb, 0x67, 0x1e, 0x5c, 0x93, 0xd2, 0xf7, 0x77, 0xbb, 0x3a, 0x61, 0x2f, 0xe1, 0xf4,
	0x4c, 0x11, 0x8d, 0xf8, 0xfd, 0x16, 0xed, 0xf0, 0x63, 0x57, 0xa7, 0xec, 0x05, 0x54, 0x91, 0xda,
	0x18, 0xc2, 0x7a, 0xc1, 0x00, 0x8a, 0x99, 0xa8, 0xb3, 0xb3, 0xe6, 0x9b, 0x18, 0xf4, 0x40, 0x83,
	0x1a, 0xe9, 0xd1, 0x22, 0xce, 0xcf, 0xae, 0x37, 0x63, 0xd7, 0xbb, 0xfd, 0xe3, 0xfb, 0x13, 0x00,
	0x00, 0xff, 0xff, 0x16, 0x07, 0x98, 0x70, 0x8f, 0x02, 0x00, 0x00,
}