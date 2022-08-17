package sign

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/module"
	Sign2 "initialthree/node/table/excel/DataTable/Sign"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"time"
)

const (
	tableName = "sign"
	dataField = "data"
	timeField = "timedata"
)

var tableFields = []string{timeField, dataField}

const (
	slotCount = 10
)

func slotName(slotIdx int) string {
	return fmt.Sprintf("slot%d", slotIdx)
}

func calcSlotIdx(id int) int {
	return id % slotCount
}

type Sign struct {
	ID           int32 `json:"id"`
	SignTimes    int32 `json:"signTimes"`
	LastSignTime int64 `json:"lastSignTime"`
}

type UserSign struct {
	userI     module.UserI
	data      map[int32]*Sign
	timedata  map[string]int64
	dataDirty map[int32]*Sign
	*module.ModuleSaveBase
}

func (this *UserSign) LastTimeDate(id int32) time.Time {
	s := this.data[id]
	if s == nil {
		return time.Time{}
	}
	return time.Unix(s.LastSignTime, 0)
}

func (this *UserSign) SignTimes(id int32) int32 {
	s := this.data[id]
	if s == nil {
		return 0
	}
	return s.SignTimes
}

func (this *UserSign) Clear(id int32) {
	delete(this.data, id)

	this.SetDirty(dataField)
	this.dataDirty[id] = nil
}

func (this *UserSign) Reset(id int32) {
	s := this.data[id]
	if s == nil {
		s = &Sign{ID: id}
		this.data[id] = s
	}
	s.SignTimes = 0

	this.SetDirty(dataField)
	this.dataDirty[id] = nil
}

func (this *UserSign) SignIn(id int32) {
	s := this.data[id]
	if s == nil {
		s = &Sign{ID: id}
		this.data[id] = s
	}

	s.SignTimes++
	s.LastSignTime = time.Now().Unix()

	this.SetDirty(dataField)
	this.dataDirty[id] = s
}

func (this *UserSign) ModuleType() module.ModuleType {
	return module.Sign
}

func (this *UserSign) Init(fields map[string]*flyfish.Field) error {
	for _, name := range tableFields {
		field, ok := fields[name]
		if ok && len(field.GetBlob()) != 0 {
			var err error
			switch name {
			case timeField:
				err = json.Unmarshal(field.GetBlob(), &this.timedata)
			case dataField:
				if err = json.Unmarshal(field.GetBlob(), &this.data); err == nil {
					// 清理已经过期的签到信息
					for _, v := range this.data {
						def := Sign2.GetID(v.ID)
						timeLimit := def.LimitDate()
						nowUnix := time.Now().Unix()
						if nowUnix < timeLimit.StartTime || nowUnix > timeLimit.EndTime {
							delete(this.data, v.ID)
						}
					}
				}
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s init sign name %s err %s ", this.userI.GetUserID(), name, err)
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}

	return nil
}

func (this *UserSign) ReadOut() *module.ReadOutCommand {
	out := &module.ReadOutCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: tableFields,
	}
	return out
}

func (this *UserSign) Tick(now time.Time) {
	// this.clockTimer()
}

func (this *UserSign) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
	}

	for field := range fields {
		name := field.(string)
		var data []byte
		switch name {
		case timeField:
			data, _ = json.Marshal(this.timedata)
		case dataField:
			data, _ = json.Marshal(this.data)
		default:
			continue
		}
		cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
			Name:  name,
			Value: data,
		})
	}

	return cmd
}

func (this *UserSign) FlushAllToClient(seqNo ...uint32) {
	mg := &message.SignSyncToC{
		IsAll:    proto.Bool(true),
		SignList: make([]*message.Sign, 0, len(this.data)),
	}

	for _, s := range this.data {
		mg.SignList = append(mg.SignList, &message.Sign{
			Id:           proto.Int32(s.ID),
			SignTimes:    proto.Int32(s.SignTimes),
			LastSignTime: proto.Int64(s.LastSignTime),
		})
	}

	this.userI.Post(mg)
	this.dataDirty = map[int32]*Sign{}
}

func (this *UserSign) FlushDirtyToClient() {
	if len(this.dataDirty) > 0 {
		mg := &message.SignSyncToC{
			IsAll:    proto.Bool(false),
			SignList: make([]*message.Sign, 0, len(this.dataDirty)),
		}

		for id, s := range this.dataDirty {
			if s != nil {
				mg.SignList = append(mg.SignList, &message.Sign{
					Id:           proto.Int32(s.ID),
					SignTimes:    proto.Int32(s.SignTimes),
					LastSignTime: proto.Int64(s.LastSignTime),
				})
			} else {
				mg.SignList = append(mg.SignList, &message.Sign{
					Id:           proto.Int32(id),
					SignTimes:    proto.Int32(0),
					LastSignTime: proto.Int64(0),
				})
			}
		}

		this.userI.Post(mg)
		this.dataDirty = map[int32]*Sign{}
	}

}

func init() {
	module.RegisterModule(module.Sign, func(userI module.UserI) module.ModuleI {
		m := &UserSign{
			userI:     userI,
			data:      map[int32]*Sign{},
			timedata:  map[string]int64{},
			dataDirty: map[int32]*Sign{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
