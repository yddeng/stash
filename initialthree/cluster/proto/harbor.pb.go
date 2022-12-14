// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/harbor.proto

package proto

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

type NotifyForginServicesH2SReq struct {
	Nodes                []uint32 `protobuf:"varint,1,rep,name=nodes" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyForginServicesH2SReq) Reset()         { *m = NotifyForginServicesH2SReq{} }
func (m *NotifyForginServicesH2SReq) String() string { return proto.CompactTextString(m) }
func (*NotifyForginServicesH2SReq) ProtoMessage()    {}
func (*NotifyForginServicesH2SReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_9678a441ab154bb9, []int{0}
}

func (m *NotifyForginServicesH2SReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyForginServicesH2SReq.Unmarshal(m, b)
}
func (m *NotifyForginServicesH2SReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyForginServicesH2SReq.Marshal(b, m, deterministic)
}
func (m *NotifyForginServicesH2SReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyForginServicesH2SReq.Merge(m, src)
}
func (m *NotifyForginServicesH2SReq) XXX_Size() int {
	return xxx_messageInfo_NotifyForginServicesH2SReq.Size(m)
}
func (m *NotifyForginServicesH2SReq) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyForginServicesH2SReq.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyForginServicesH2SReq proto.InternalMessageInfo

func (m *NotifyForginServicesH2SReq) GetNodes() []uint32 {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type NotifyForginServicesH2SResp struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyForginServicesH2SResp) Reset()         { *m = NotifyForginServicesH2SResp{} }
func (m *NotifyForginServicesH2SResp) String() string { return proto.CompactTextString(m) }
func (*NotifyForginServicesH2SResp) ProtoMessage()    {}
func (*NotifyForginServicesH2SResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_9678a441ab154bb9, []int{1}
}

func (m *NotifyForginServicesH2SResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyForginServicesH2SResp.Unmarshal(m, b)
}
func (m *NotifyForginServicesH2SResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyForginServicesH2SResp.Marshal(b, m, deterministic)
}
func (m *NotifyForginServicesH2SResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyForginServicesH2SResp.Merge(m, src)
}
func (m *NotifyForginServicesH2SResp) XXX_Size() int {
	return xxx_messageInfo_NotifyForginServicesH2SResp.Size(m)
}
func (m *NotifyForginServicesH2SResp) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyForginServicesH2SResp.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyForginServicesH2SResp proto.InternalMessageInfo

//harbor????????????
type NotifyForginServicesH2HReq struct {
	Nodes                []uint32 `protobuf:"varint,1,rep,name=nodes" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyForginServicesH2HReq) Reset()         { *m = NotifyForginServicesH2HReq{} }
func (m *NotifyForginServicesH2HReq) String() string { return proto.CompactTextString(m) }
func (*NotifyForginServicesH2HReq) ProtoMessage()    {}
func (*NotifyForginServicesH2HReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_9678a441ab154bb9, []int{2}
}

func (m *NotifyForginServicesH2HReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyForginServicesH2HReq.Unmarshal(m, b)
}
func (m *NotifyForginServicesH2HReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyForginServicesH2HReq.Marshal(b, m, deterministic)
}
func (m *NotifyForginServicesH2HReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyForginServicesH2HReq.Merge(m, src)
}
func (m *NotifyForginServicesH2HReq) XXX_Size() int {
	return xxx_messageInfo_NotifyForginServicesH2HReq.Size(m)
}
func (m *NotifyForginServicesH2HReq) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyForginServicesH2HReq.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyForginServicesH2HReq proto.InternalMessageInfo

func (m *NotifyForginServicesH2HReq) GetNodes() []uint32 {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type NotifyForginServicesH2HResp struct {
	Nodes                []uint32 `protobuf:"varint,1,rep,name=nodes" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NotifyForginServicesH2HResp) Reset()         { *m = NotifyForginServicesH2HResp{} }
func (m *NotifyForginServicesH2HResp) String() string { return proto.CompactTextString(m) }
func (*NotifyForginServicesH2HResp) ProtoMessage()    {}
func (*NotifyForginServicesH2HResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_9678a441ab154bb9, []int{3}
}

func (m *NotifyForginServicesH2HResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyForginServicesH2HResp.Unmarshal(m, b)
}
func (m *NotifyForginServicesH2HResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyForginServicesH2HResp.Marshal(b, m, deterministic)
}
func (m *NotifyForginServicesH2HResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyForginServicesH2HResp.Merge(m, src)
}
func (m *NotifyForginServicesH2HResp) XXX_Size() int {
	return xxx_messageInfo_NotifyForginServicesH2HResp.Size(m)
}
func (m *NotifyForginServicesH2HResp) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyForginServicesH2HResp.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyForginServicesH2HResp proto.InternalMessageInfo

func (m *NotifyForginServicesH2HResp) GetNodes() []uint32 {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func init() {
	proto.RegisterType((*NotifyForginServicesH2SReq)(nil), "proto.notifyForginServicesH2S_req")
	proto.RegisterType((*NotifyForginServicesH2SResp)(nil), "proto.notifyForginServicesH2S_resp")
	proto.RegisterType((*NotifyForginServicesH2HReq)(nil), "proto.notifyForginServicesH2H_req")
	proto.RegisterType((*NotifyForginServicesH2HResp)(nil), "proto.notifyForginServicesH2H_resp")
}

func init() { proto.RegisterFile("proto/harbor.proto", fileDescriptor_9678a441ab154bb9) }

var fileDescriptor_9678a441ab154bb9 = []byte{
	// 116 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x2a, 0x28, 0xca, 0x2f,
	0xc9, 0xd7, 0xcf, 0x48, 0x2c, 0x4a, 0xca, 0x2f, 0xd2, 0x03, 0x73, 0x84, 0x58, 0xc1, 0x94, 0x92,
	0x31, 0x97, 0x74, 0x5e, 0x7e, 0x49, 0x66, 0x5a, 0xa5, 0x5b, 0x7e, 0x51, 0x7a, 0x66, 0x5e, 0x70,
	0x6a, 0x51, 0x59, 0x66, 0x72, 0x6a, 0xb1, 0x87, 0x51, 0x70, 0x7c, 0x51, 0x6a, 0xa1, 0x90, 0x08,
	0x17, 0x6b, 0x5e, 0x7e, 0x4a, 0x6a, 0xb1, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0x6f, 0x10, 0x84, 0xa3,
	0x24, 0xc7, 0x25, 0x83, 0x5b, 0x53, 0x71, 0x01, 0x6e, 0x43, 0x3d, 0xf0, 0x18, 0x6a, 0x82, 0xcb,
	0x50, 0x0f, 0xb0, 0xa1, 0xd8, 0x75, 0x01, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4a, 0x5f, 0x6c, 0x03,
	0xdb, 0x00, 0x00, 0x00,
}
