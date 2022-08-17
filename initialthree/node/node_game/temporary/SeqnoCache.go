package temporary

import (
	"fmt"
	codecs "initialthree/codec/cs"
	"initialthree/pkg/json"
	"initialthree/zaplogger"

	"time"
)

const reqCacheCount = 10

type SeqCache struct {
	userI     UserI
	respMsg   []*codecs.Message // 客户端请求的缓存
	nextSeqNo uint32
	headSeqNo uint32
}

func (this *SeqCache) AddReqCache(seqNo uint32, msg *codecs.Message) {
	for _, v := range this.respMsg {
		if v.GetSeriNo() == seqNo {
			return
		}
	}

	if len(this.respMsg) >= reqCacheCount {
		this.respMsg = append(this.respMsg[:0], this.respMsg[1:]...)
	}
	this.respMsg = append(this.respMsg, msg)
}

func (this *SeqCache) GetReqCache(seqNo uint32) (*codecs.Message, error) {

	seqLen := len(this.respMsg)
	if seqLen == 0 {
		this.nextSeqNo = seqNo + 1
		this.headSeqNo = seqNo
		zaplogger.GetSugar().Debugf("%s seqLen=%d curSeqNo=%d nextSeqNo=%d first", this.userI.GetUserID(), seqLen, seqNo, this.nextSeqNo)
		return nil, nil
	}

	zaplogger.GetSugar().Debugf("%s seqLen=%d curSeqNo=%d nextSeqNo=%d", this.userI.GetUserID(), seqLen, seqNo, this.nextSeqNo)

	if seqNo == this.nextSeqNo {
		if this.nextSeqNo-this.headSeqNo >= reqCacheCount {
			this.headSeqNo++
		}
		this.nextSeqNo++
		return nil, nil
	}

	// 超出缓存范围
	if seqNo < this.headSeqNo || seqNo > this.nextSeqNo {
		return nil, fmt.Errorf("seqNo:%d error,out range [%d %d]", seqNo, this.headSeqNo, this.nextSeqNo)
	}

	for _, v := range this.respMsg {
		if v.GetSeriNo() == seqNo {
			return v, nil
		}
	}

	return nil, nil
}

func NewSeqCache(userI UserI) *SeqCache {
	return &SeqCache{
		userI:   userI,
		respMsg: make([]*codecs.Message, 0, reqCacheCount),
	}
}

func (this *SeqCache) UserDisconnect() {
	this.respMsg = this.respMsg[:0]
}

func (this *SeqCache) UserLogout() {
	this.respMsg = this.respMsg[:0]
}

func (this *SeqCache) Tick(now time.Time) {}

/*
 *
 */

type seqnoMarshaler struct{}

func (m *seqnoMarshaler) Marshal(temp TemporaryI) ([]byte, error) {
	info := temp.(*SeqCache)
	cache := seqnoCache{
		NextSeqNo: info.nextSeqNo,
		HeadSeqNo: info.headSeqNo,
		Msgs:      make([]seqnoMsg, 0, len(info.respMsg)),
	}
	for _, msg := range info.respMsg {
		sMsg := seqnoMsg{
			SeqNo: msg.GetSeriNo(),
			Cmd:   msg.GetCmd(),
			Err:   msg.GetErrCode(),
		}
		if sMsg.Err == 0 {
			data, _, err := encoder.ProtoMarshal(msg.GetData())
			if err != nil {
				return nil, err
			}
			sMsg.Data = data
		}
		cache.Msgs = append(cache.Msgs, sMsg)

	}
	return json.Marshal(cache)
}

func (m *seqnoMarshaler) Unmarshal(user UserI, data []byte) (TemporaryI, error) {
	var cache seqnoCache
	err := json.Unmarshal(data, &cache)
	if err != nil {
		return nil, err
	}

	info := &SeqCache{
		userI:     user,
		respMsg:   make([]*codecs.Message, 0, len(cache.Msgs)),
		nextSeqNo: cache.NextSeqNo,
		headSeqNo: cache.HeadSeqNo,
	}
	for _, v := range cache.Msgs {
		if v.Err == 0 {
			if pMsg, err := encoder.ProtoUnmarshal(v.Cmd, v.Data); err != nil {
				return nil, fmt.Errorf("%d %d %s", v.SeqNo, v.Cmd, err.Error())
			} else {
				info.respMsg = append(info.respMsg, codecs.NewMessage(v.SeqNo, pMsg))
			}
		} else {
			info.respMsg = append(info.respMsg, codecs.ErrMessage(v.SeqNo, v.Cmd, v.Err))
		}
	}
	return info, nil
}

type seqnoCache struct {
	NextSeqNo uint32     `json:"next_seq_no"`
	HeadSeqNo uint32     `json:"head_seq_no"`
	Msgs      []seqnoMsg `json:"msgs"`
}

type seqnoMsg struct {
	SeqNo uint32 `json:"seq_no"`
	Cmd   uint16 `json:"cmd"`
	Err   uint16 `json:"err,omitempty"`
	Data  []byte `json:"data,omitempty"`
}

func init() {
	registerTempDataProcess(TempSeqCache, "seqno", &seqnoMarshaler{})
}
