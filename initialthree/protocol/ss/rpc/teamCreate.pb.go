// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ss/proto/rpc/teamCreate.proto

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

type TeamCreateReq struct {
	Player               *message.TeamPlayer `protobuf:"bytes,1,req,name=player" json:"player,omitempty"`
	Target               *message.TeamTarget `protobuf:"bytes,2,req,name=target" json:"target,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *TeamCreateReq) Reset()         { *m = TeamCreateReq{} }
func (m *TeamCreateReq) String() string { return proto.CompactTextString(m) }
func (*TeamCreateReq) ProtoMessage()    {}
func (*TeamCreateReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8f687580c083b4, []int{0}
}

func (m *TeamCreateReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamCreateReq.Unmarshal(m, b)
}
func (m *TeamCreateReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamCreateReq.Marshal(b, m, deterministic)
}
func (m *TeamCreateReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamCreateReq.Merge(m, src)
}
func (m *TeamCreateReq) XXX_Size() int {
	return xxx_messageInfo_TeamCreateReq.Size(m)
}
func (m *TeamCreateReq) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamCreateReq.DiscardUnknown(m)
}

var xxx_messageInfo_TeamCreateReq proto.InternalMessageInfo

func (m *TeamCreateReq) GetPlayer() *message.TeamPlayer {
	if m != nil {
		return m.Player
	}
	return nil
}

func (m *TeamCreateReq) GetTarget() *message.TeamTarget {
	if m != nil {
		return m.Target
	}
	return nil
}

type TeamCreateResp struct {
	ErrCode              *message.ErrCode `protobuf:"varint,1,req,name=errCode,enum=message.ErrCode" json:"errCode,omitempty"`
	TeamID               *uint32          `protobuf:"varint,2,req,name=teamID" json:"teamID,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *TeamCreateResp) Reset()         { *m = TeamCreateResp{} }
func (m *TeamCreateResp) String() string { return proto.CompactTextString(m) }
func (*TeamCreateResp) ProtoMessage()    {}
func (*TeamCreateResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_ec8f687580c083b4, []int{1}
}

func (m *TeamCreateResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TeamCreateResp.Unmarshal(m, b)
}
func (m *TeamCreateResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TeamCreateResp.Marshal(b, m, deterministic)
}
func (m *TeamCreateResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TeamCreateResp.Merge(m, src)
}
func (m *TeamCreateResp) XXX_Size() int {
	return xxx_messageInfo_TeamCreateResp.Size(m)
}
func (m *TeamCreateResp) XXX_DiscardUnknown() {
	xxx_messageInfo_TeamCreateResp.DiscardUnknown(m)
}

var xxx_messageInfo_TeamCreateResp proto.InternalMessageInfo

func (m *TeamCreateResp) GetErrCode() message.ErrCode {
	if m != nil && m.ErrCode != nil {
		return *m.ErrCode
	}
	return message.ErrCode_OK
}

func (m *TeamCreateResp) GetTeamID() uint32 {
	if m != nil && m.TeamID != nil {
		return *m.TeamID
	}
	return 0
}

func init() {
	proto.RegisterType((*TeamCreateReq)(nil), "rpc.teamCreate_req")
	proto.RegisterType((*TeamCreateResp)(nil), "rpc.teamCreate_resp")
}

func init() { proto.RegisterFile("ss/proto/rpc/teamCreate.proto", fileDescriptor_ec8f687580c083b4) }

var fileDescriptor_ec8f687580c083b4 = []byte{
	// 219 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0xd0, 0xcf, 0x6b, 0xc2, 0x30,
	0x14, 0x07, 0x70, 0xd6, 0x41, 0x07, 0x19, 0xeb, 0x46, 0x06, 0xa3, 0x74, 0x3f, 0xe9, 0x69, 0x6c,
	0xd0, 0x40, 0xff, 0x84, 0x75, 0x1e, 0xbc, 0x49, 0xa8, 0x17, 0x2f, 0x12, 0xe2, 0xa3, 0x56, 0x5a,
	0x13, 0x5f, 0x72, 0xf1, 0xbf, 0x97, 0x26, 0xa9, 0x22, 0x7a, 0xec, 0xfb, 0x7e, 0xde, 0xfb, 0xd2,
	0x90, 0x77, 0x63, 0x98, 0x46, 0x65, 0x15, 0x43, 0x2d, 0x99, 0x05, 0xd1, 0x57, 0x08, 0xc2, 0x42,
	0xe1, 0x86, 0xf4, 0x16, 0xb5, 0xcc, 0x5e, 0xe5, 0x68, 0x7a, 0x30, 0x46, 0x34, 0xe0, 0x9c, 0x17,
	0xd9, 0xe7, 0x45, 0x08, 0x88, 0x80, 0xa8, 0xd0, 0x83, 0x7c, 0x43, 0x92, 0xd3, 0xd9, 0x25, 0xc2,
	0x8e, 0xfe, 0x92, 0x58, 0x77, 0x62, 0x0f, 0x98, 0xde, 0x7c, 0x45, 0xdf, 0xf7, 0xe5, 0x73, 0x11,
	0x56, 0x8b, 0x1a, 0x44, 0x3f, 0x73, 0x11, 0x0f, 0x64, 0xc0, 0x56, 0x60, 0x03, 0x36, 0x8d, 0xae,
	0xe0, 0xda, 0x45, 0x3c, 0x90, 0x7c, 0x4e, 0x1e, 0xcf, 0xba, 0x8c, 0xa6, 0x3f, 0xe4, 0x0e, 0x10,
	0x2b, 0xb5, 0x02, 0xd7, 0x96, 0x94, 0x4f, 0xc7, 0x03, 0x13, 0x3f, 0xe7, 0x23, 0xa0, 0x2f, 0x24,
	0x1e, 0xd6, 0xa7, 0xff, 0xae, 0xeb, 0x81, 0x87, 0xaf, 0xbf, 0x8f, 0xc5, 0x5b, 0xbb, 0x6d, 0x6d,
	0x2b, 0x3a, 0xbb, 0x46, 0x00, 0xff, 0xbf, 0x52, 0x75, 0xcc, 0x98, 0xe1, 0xd9, 0x0e, 0x01, 0x00,
	0x00, 0xff, 0xff, 0x2d, 0x30, 0x48, 0x8c, 0x45, 0x01, 0x00, 0x00,
}
