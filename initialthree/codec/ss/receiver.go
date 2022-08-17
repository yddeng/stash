package ss

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	kbuffer "github.com/sniperHW/flyfish/pkg/buffer"
	"initialthree/cluster/addr"
	"initialthree/cluster/rpcerr"
	"initialthree/codec/pb"
	"initialthree/pkg/rpc"
)

const (
	maxPacketSize int = 65535 * 4
	minSize       int = sizeLen + sizeFlag
)

type Receiver struct {
	buffer      []byte
	w           int
	r           int
	ns_msg      string
	ns_req      string
	ns_resp     string
	lastRecved  int
	recvCount   int
	unpackCount int
	selfAddr    addr.LogicAddr
}

func (this *Receiver) isTarget(to addr.LogicAddr) bool {
	return this.selfAddr == to
}

func (this *Receiver) unpack(buffer []byte, r int, w int) (ret interface{}, packetSize int, err error) {
	unpackSize := w - r
	if unpackSize >= minSize {
		var totalSize int
		var payload int
		var flag byte
		var cmd uint16
		var buff []byte
		var errStr string
		var msg proto.Message
		var to addr.LogicAddr
		var from addr.LogicAddr
		var t uint32
		var f uint32
		var addrSize int
		reader := kbuffer.NewReader(buffer[r : r+unpackSize])
		if payload, err = reader.CheckGetInt(); err != nil {
			return
		}

		if payload == 0 {
			err = fmt.Errorf("zero payload")
			return
		}

		totalSize = int(payload + sizeLen)

		if totalSize > maxPacketSize {
			err = fmt.Errorf("large packet %d", totalSize)
			return
		}

		packetSize = totalSize

		if totalSize <= unpackSize {

			if flag, err = reader.CheckGetByte(); err != nil {
				return
			}

			if isRelay(flag) {

				if t, err = reader.CheckGetUint32(); err != nil {
					return
				}
				if f, err = reader.CheckGetUint32(); err != nil {
					return
				}
				to = addr.LogicAddr(t)
				from = addr.LogicAddr(f)
				addrSize = sizeTo + sizeFrom
			}

			if (isRelay(flag) && this.isTarget(to)) || !isRelay(flag) {

				if cmd, err = reader.CheckGetUint16(); err != nil {
					return
				}

				tt := getMsgType(flag)
				if tt == MESSAGE {
					//普通消息
					size := payload - (sizeCmd + sizeFlag + addrSize)
					if buff, err = reader.CheckGetBytes(size); err != nil {
						return
					}
					if msg, err = pb.Unmarshal(this.ns_msg, uint32(cmd), buff); err != nil {
						return
					}
					ret = NewMessage(msg, to, from)
				} else {
					var seqNO uint64
					if seqNO, err = reader.CheckGetUint64(); err != nil {
						return
					}

					size := payload - (sizeCmd + sizeFlag + sizeRPCSeqNo + addrSize)
					if tt == RPCERR {
						//RPC响应错误信息
						if errStr, err = reader.CheckGetString(size); err != nil {
							return
						}

						ret = NewMessage(&rpc.RPCResponse{Seq: seqNO, Err: rpcerr.GetErrorByShortStr(errStr)}, to, from)
					} else if tt == RPCRESP {
						//RPC响应
						if buff, err = reader.CheckGetBytes(size); err != nil {
							return
						}
						if msg, err = pb.Unmarshal(this.ns_resp, uint32(cmd), buff); err != nil {
							return
						}
						ret = NewMessage(&rpc.RPCResponse{Seq: seqNO, Ret: msg}, to, from)
					} else if tt == RPCREQ {
						//RPC请求
						if buff, err = reader.CheckGetBytes(size); err != nil {
							return
						}
						if msg, err = pb.Unmarshal(this.ns_req, uint32(cmd), buff); err != nil {
							return
						}
						ret = NewMessage(&rpc.RPCRequest{
							Seq:      seqNO,
							Method:   pb.GetNameByID(this.ns_req, uint32(cmd)),
							NeedResp: getNeedRPCResp(flag),
							Arg:      msg,
						}, to, from)
					} else {
						err = fmt.Errorf("invaild message type")
					}
				}
			} else {
				ret = NewRelayMessage(to, from, buffer[r:r+totalSize])
			}
		}
	}
	return
}

func NewReceiver(ns_msg, ns_req, ns_resp string, selfAddr ...addr.LogicAddr) *Receiver {
	receiver := &Receiver{}
	receiver.ns_msg = ns_msg
	receiver.ns_req = ns_req
	receiver.ns_resp = ns_resp
	receiver.buffer = make([]byte, maxPacketSize*2)

	if len(selfAddr) > 0 {
		receiver.selfAddr = selfAddr[0]
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
