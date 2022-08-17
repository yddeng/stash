package cs

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"github.com/golang/protobuf/proto"
	kbuffer "github.com/sniperHW/flyfish/pkg/buffer"
	"initialthree/codec/pb"
	"io"
	"io/ioutil"
)

const (
	MaxPacketSize int = 65535
)

type Receiver struct {
	buffer    []byte
	w         int
	r         int
	namespace string
	zipBuff   bytes.Buffer
	unpackMsg map[uint16]bool
}

func isCompress(flag uint16) bool {
	return flag&0x8000 > 0
}

func (this *Receiver) DirectUnpack(buff []byte) (interface{}, error) {
	if len(buff) > MaxPacketSize {
		return nil, fmt.Errorf("packet too large totalLen:%d", len(buff))
	}
	msg, _, err := this.unpack(buff, 0, len(buff))
	return msg, err
}

func (this *Receiver) ProtoMarshal(msg proto.Message) ([]byte, uint32, error) {
	return pb.Marshal(this.namespace, msg)
}

func (this *Receiver) ProtoUnmarshal(cmd uint16, data []byte) (proto.Message, error) {
	return pb.Unmarshal(this.namespace, uint32(cmd), data)
}

func (this *Receiver) unpack(buffer []byte, r int, w int) (ret interface{}, packetSize int, err error) {
	unpackSize := w - r
	if unpackSize >= HeadSize {

		var payload uint16
		var seqNo uint32
		var flag uint16
		var cmd uint16
		var teachID uint16
		var errCode uint16
		var buff []byte
		var msg proto.Message

		reader := kbuffer.NewReader(buffer[r : r+unpackSize])

		if payload, err = reader.CheckGetUint16(); err != nil {
			return
		}

		fullSize := int(payload) + SizeLen

		if fullSize >= MaxPacketSize {
			err = fmt.Errorf("packet too large %d", fullSize)
			return
		}

		if payload == 0 {
			err = fmt.Errorf("zero packet")
			return
		}

		packetSize = fullSize

		if fullSize <= unpackSize {
			if seqNo, err = reader.CheckGetUint32(); err != nil {
				return
			}

			if flag, err = reader.CheckGetUint16(); err != nil {
				return
			}

			if cmd, err = reader.CheckGetUint16(); err != nil {
				return
			}

			if teachID, err = reader.CheckGetUint16(); err != nil {
				return
			}

			if errCode, err = reader.CheckGetUint16(); err != nil {
				return
			}

			size := payload - (HeadSize - SizeLen)
			if buff, err = reader.CheckGetBytes(int(size)); err != nil {
				return
			}

			if this.unpackMsg != nil {
				if _, ok := this.unpackMsg[cmd]; !ok {
					//透传消息
					//message := NewBytesMassage(this.buffer[r : r+fullSize])
					message := make([]byte, fullSize)
					copy(message, this.buffer[r:r+fullSize])
					ret = message
					return
				}
			}

			if isCompress(flag) {
				this.zipBuff.Reset()
				this.zipBuff.Write(buff)
				var r io.ReadCloser
				r, err = zlib.NewReader(&this.zipBuff)
				if err != nil {
					return
				}

				buff, err = ioutil.ReadAll(r)
				r.Close()
				if err != nil {
					if err != io.ErrUnexpectedEOF && err != io.EOF {
						return
					}
				}
			}

			if errCode == 0 {
				if msg, err = pb.Unmarshal(this.namespace, uint32(cmd), buff); err != nil {
					return
				}
			}

			ret = &Message{
				seriNO:  seqNo,
				data:    msg,
				cmd:     cmd,
				teachID: teachID,
				errCode: errCode,
			}
			return
		} else {
			return
		}
	}
	return
}

/*
 * 不设置 unpackMag 时，全部拆包
 * 设置时，遇到 unpackMag 注册的消息，拆包。其他消息拆分为字节消息
 */
func NewReceiver(namespace string, unpackMsg ...map[uint16]bool) *Receiver {
	receiver := &Receiver{}
	receiver.namespace = namespace
	receiver.buffer = make([]byte, MaxPacketSize*2)
	if len(unpackMsg) > 0 {
		receiver.unpackMsg = unpackMsg[0]
	}
	return receiver
}

func (this *Receiver) GetRecvBuff() []byte {
	return this.buffer[this.w:]
}

func (this *Receiver) OnData(data []byte) {
	this.w += len(data)
}

func (this *Receiver) Unpack() (interface{}, error) {
	if this.r == this.w {
		return nil, nil
	} else {
		msg, packetSize, err := this.unpack(this.buffer, this.r, this.w)
		if nil != msg {
			this.r += packetSize
			if this.r == this.w {
				this.r = 0
				this.w = 0
			}
		} else if nil == err {
			if packetSize > cap(this.buffer) {
				buffer := make([]byte, packetSize)
				copy(buffer, this.buffer[this.r:this.w])
				this.buffer = buffer
			} else {
				//空间足够容纳下一个包，
				copy(this.buffer, this.buffer[this.r:this.w])
			}
			this.w = this.w - this.r
			this.r = 0
		}
		return msg, err
	}
}

//SizeLen  = 2
//SizeSeqNo  = 4
//SizeFlag = 2
//SizeCmd  = 2
//SizeTeach = 2
//SizeErr  = 2
func FetchSeqCmdCode(buff []byte) (uint32, uint16, uint16) {
	seqno := binary.BigEndian.Uint32(buff[2 : 2+4])
	cmd := binary.BigEndian.Uint16(buff[8 : 8+2])
	code := binary.BigEndian.Uint16(buff[12 : 12+2])
	return seqno, cmd, code
}
