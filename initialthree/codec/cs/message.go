package cs

import (
	"github.com/golang/protobuf/proto"
	"initialthree/codec/pb"
	"reflect"
)

type Message struct {
	seriNO   uint32
	data     proto.Message
	cmd      uint16
	compress bool   //是否压缩
	errCode  uint16 // 错误码
	teachID  uint16
}

func NewMessage(seriNO uint32, data proto.Message) *Message {
	return &Message{seriNO: seriNO, data: data}
}

func ErrMessage(seriNO uint32, cmd uint16, errCode uint16) *Message {
	return &Message{seriNO: seriNO, cmd: cmd, errCode: errCode}
}

func (this *Message) SetCompress() *Message {
	this.compress = true
	return this
}

func (this *Message) IsCompress() bool {
	return this.compress
}

func (this *Message) SetTeachID(teachID uint16) *Message {
	this.teachID = teachID
	return this
}

func (this *Message) GetTeachID() uint16 {
	return this.teachID
}

func (this *Message) GetData() proto.Message {
	return this.data
}

func (this *Message) GetSeriNo() uint32 {
	return this.seriNO
}

func (this *Message) GetErrCode() uint16 {
	return this.errCode
}

func (this *Message) GetCmd() uint16 {
	if this.data != nil {
		name := reflect.TypeOf(this.data).String()
		return uint16(pb.GetCmdByName(name))
	}
	return this.cmd
}

/*
// 字节消息
type BytesMassage struct {
	bytes []byte
}

func NewBytesMassage(bytes []byte) *BytesMassage {
	msg := &BytesMassage{
		bytes: make([]byte, len(bytes)),
	}
	copy(msg.bytes, bytes)
	return msg
}

func (this *BytesMassage) GetData() []byte {
	return this.bytes
}
*/
