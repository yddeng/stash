// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss/proto/rpc/teamPlayerGetFromGame.proto

package rpc

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	message "initialthree/protocol/cs/message"
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

type TeamPlayerGetFromGameReq struct {
	UID                  *string  `protobuf:"bytes,1,req,name=uID" json:"uID,omitempty"`
	TeamID               *uint32  `protobuf:"varint,2,req,name=teamID" json:"teamID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TeamPlayerGetFromGameReq) Reset()         { *m = TeamPlayerGetFromGameReq{} }
func (m *TeamPlayerGetFromGameReq) String() string { return proto.CompactTextString(m) }
func (*TeamPlayerGetFromGameReq) ProtoMessage()    {}
func (*TeamPlayerGetFromGameReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_e0edd0f0fdf05ab0, []int{0}
}

func (m *TeamPlayerGetFromGameReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamPlayerGetFromGameReq.Unmarshal(m, b)
}
func (m *TeamPlayerGetFromGameReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamPlayerGetFromGameReq.Marshal(b, m, deterministic)
}
func (m *TeamPlayerGetFromGameReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamPlayerGetFromGameReq.Merge(m, src)
}
func (m *TeamPlayerGetFromGameReq) XXX_Size() int {
	return xxx_messageInfo_TeamPlayerGetFromGameReq.Size(m)
}
func (m *TeamPlayerGetFromGameReq) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamPlayerGetFromGameReq.DiscardUnknown(m)
}

var xxx_messageInfo_TeamPlayerGetFromGameReq proto.InternalMessageInfo

func (m *TeamPlayerGetFromGameReq) GetUID() string {
	if m != nil && m.UID != nil {
		return *m.UID
	}
	return ""
}

func (m *TeamPlayerGetFromGameReq) GetTeamID() uint32 {
	if m != nil && m.TeamID != nil {
		return *m.TeamID
	}
	return 0
}

type TeamPlayerGetFromGameResp struct {
	ErrCode              *message.ErrCode    `protobuf:"varint,1,req,name=errCode,enum=message.ErrCode" json:"errCode,omitempty"`
	Player               *message.TeamPlayer `protobuf:"bytes,2,opt,name=player" json:"player,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TeamPlayerGetFromGameResp) Reset()         { *m = TeamPlayerGetFromGameResp{} }
func (m *TeamPlayerGetFromGameResp) String() string { return proto.CompactTextString(m) }
func (*TeamPlayerGetFromGameResp) ProtoMessage()    {}
func (*TeamPlayerGetFromGameResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_e0edd0f0fdf05ab0, []int{1}
}

func (m *TeamPlayerGetFromGameResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamPlayerGetFromGameResp.Unmarshal(m, b)
}
func (m *TeamPlayerGetFromGameResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamPlayerGetFromGameResp.Marshal(b, m, deterministic)
}
func (m *TeamPlayerGetFromGameResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamPlayerGetFromGameResp.Merge(m, src)
}
func (m *TeamPlayerGetFromGameResp) XXX_Size() int {
	return xxx_messageInfo_TeamPlayerGetFromGameResp.Size(m)
}
func (m *TeamPlayerGetFromGameResp) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamPlayerGetFromGameResp.DiscardUnknown(m)
}

var xxx_messageInfo_TeamPlayerGetFromGameResp proto.InternalMessageInfo

func (m *TeamPlayerGetFromGameResp) GetErrCode() message.ErrCode {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return message.ErrCode_OK
}

func (m *TeamPlayerGetFromGameResp) GetPlayer() *message.TeamPlayer {
	if m != nil {
		return m.Player
	}
	return nil
}

func init() {
	proto.RegisterType((*TeamPlayerGetFromGameReq)(nil), "rpc.teamPlayerGetFromGame_req")
	proto.RegisterType((*TeamPlayerGetFromGameResp)(nil), "rpc.teamPlayerGetFromGame_resp")
}

func init() {
	proto.RegisterFile("ss/proto/rpc/teamPlayerGetFromGame.proto", fileDescriptor_e0edd0f0fdf05ab0)
}

var fileDescriptor_e0edd0f0fdf05ab0 = []byte{
	// 224 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x90, 0x41, 0x4b, 0xc4, 0x30,
	0x10, 0x85, 0xd9, 0x2c, 0xac, 0x18, 0x51, 0x96, 0x08, 0x52, 0xab, 0x68, 0xd9, 0x53, 0x50, 0x68,
	0x60, 0x7f, 0x82, 0xb6, 0x96, 0xde, 0x24, 0x78, 0xf2, 0x22, 0x21, 0x0e, 0x5a, 0x68, 0x4c, 0x9c,
	0xa4, 0x07, 0xff, 0xbd, 0x34, 0x69, 0x7b, 0xd1, 0xbd, 0x25, 0xf3, 0xbe, 0x79, 0x6f, 0x66, 0x28,
	0xf7, 0x5e, 0x38, 0xb4, 0xc1, 0x0a, 0x74, 0x5a, 0x04, 0x50, 0xe6, 0xb9, 0x57, 0x3f, 0x80, 0x0d,
	0x84, 0x27, 0xb4, 0xa6, 0x51, 0x06, 0xca, 0xa8, 0xb3, 0x35, 0x3a, 0x9d, 0x5f, 0xe9, 0x19, 0x37,
	0xe0, 0xbd, 0xfa, 0x80, 0xd8, 0x92, 0x88, 0xfc, 0xf6, 0x8f, 0x08, 0x88, 0x80, 0x68, 0x31, 0x01,
	0xbb, 0x9a, 0x5e, 0xfe, 0x9b, 0xf0, 0x86, 0xf0, 0xcd, 0xb6, 0x74, 0x3d, 0xb4, 0x55, 0xb6, 0x2a,
	0x08, 0x3f, 0x96, 0xe3, 0x93, 0x5d, 0xd0, 0xcd, 0x88, 0xb7, 0x55, 0x46, 0x0a, 0xc2, 0x4f, 0xe5,
	0xf4, 0xdb, 0x0d, 0x34, 0x3f, 0x64, 0xe3, 0x1d, 0xbb, 0xa3, 0x47, 0x80, 0xf8, 0x68, 0xdf, 0x21,
	0x7a, 0x9d, 0xed, 0xb7, 0xe5, 0x34, 0x4e, 0x59, 0xa7, 0xba, 0x9c, 0x01, 0x76, 0x4f, 0x37, 0x2e,
	0xba, 0x64, 0xa4, 0x58, 0xf1, 0x93, 0xfd, 0xf9, 0x82, 0xbe, 0x2c, 0x01, 0x72, 0x42, 0x1e, 0x6e,
	0x5e, 0xaf, 0xbb, 0xaf, 0x2e, 0x74, 0xaa, 0x0f, 0x9f, 0x08, 0x90, 0x56, 0xd5, 0xb6, 0x17, 0xde,
	0x8f, 0xc7, 0xfb, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x9d, 0x36, 0x3e, 0x0e, 0x4b, 0x01, 0x00, 0x00,
}
