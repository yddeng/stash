package cs

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/sniperHW/flyfish/pkg/buffer"
	"initialthree/codec/pb"
	_ "initialthree/protocol/cs"
)

const (
	SizeLen   = 2
	SizeSeq   = 4
	SizeFlag  = 2
	SizeCmd   = 2
	SizeTeach = 2 // 教学引导ID
	SizeErr   = 2
	HeadSize  = SizeLen + SizeSeq + SizeFlag + SizeCmd + SizeTeach + SizeErr
)

/*
   length(2) | seqNo(4) | flag(2) | msgID(2) | teachID(2) | errcode(2) | content ...
*/

type Encoder struct {
	namespace string
	zipBuff   bytes.Buffer
	zipWriter *zlib.Writer
}

func NewEncoder(namespace string) *Encoder {
	return &Encoder{namespace: namespace}
}

func setCompressFlag(flag uint16) uint16 {
	return flag | 0x8000
}

func (this *Encoder) EnCode(o interface{}, buff *buffer.Buffer) error {
	switch o.(type) {
	case *Message:
		var pbbytes []byte
		var cmd uint32
		var err error

		msg := o.(*Message)
		data := msg.GetData()
		flag := uint16(0)
		seqNo := msg.GetSeriNo()

		if data == nil {
			cmd = uint32(msg.GetCmd())
			pbbytes = make([]byte, 0)
		} else {
			if pbbytes, cmd, err = pb.Marshal(this.namespace, msg.GetData()); err != nil {
				return err
			}
		}

		if msg.IsCompress() {
			flag = setCompressFlag(flag)
			//对pb数据执行压缩
			if nil == this.zipWriter {
				this.zipWriter = zlib.NewWriter(&this.zipBuff)
			} else {
				this.zipBuff.Reset()
				this.zipWriter.Reset(&this.zipBuff)
			}

			this.zipWriter.Write(pbbytes)
			this.zipWriter.Flush()
			pbbytes = this.zipBuff.Bytes()
		}

		totalLen := len(pbbytes) + HeadSize

		if totalLen > MaxPacketSize {
			return fmt.Errorf("packet too large totalLen:%d", totalLen)
		}

		//写payload大小
		buff.AppendUint16(uint16(totalLen - SizeLen))
		//写seq
		buff.AppendUint32(seqNo)
		//写flag
		buff.AppendUint16(flag)
		//写cmd
		buff.AppendUint16(uint16(cmd))
		// teachID
		buff.AppendUint16(msg.GetTeachID())
		//errCode
		buff.AppendUint16(msg.GetErrCode())
		//写数据
		buff.AppendBytes(pbbytes)
		return nil
	case []byte:
		//透传消息
		buff.AppendBytes(o.([]byte))
		return nil
	default:
		return fmt.Errorf("invaild object type")
	}
}

func (this *Encoder) ProtoMarshal(msg proto.Message) ([]byte, uint32, error) {
	return pb.Marshal(this.namespace, msg)
}

func (this *Encoder) ProtoUnmarshal(cmd uint16, data []byte) (proto.Message, error) {
	return pb.Unmarshal(this.namespace, uint32(cmd), data)
}
