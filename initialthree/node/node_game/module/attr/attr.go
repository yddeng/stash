//角色属性相关
package attr

import (
	"fmt"
	comAttr "initialthree/node/common/attr"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module"
	"initialthree/pkg/json"
	cs_message "initialthree/protocol/cs/message"
	"time"

	flyfish "github.com/sniperHW/flyfish/client"

	"github.com/golang/protobuf/proto"
)

const (
	timeField  = "attr_time"
	dataField  = "attr"
	timePrefix = "time_"
)

type AttrData struct {
	Attrs []*comAttr.Attr `json:"As"`
}

type UserAttr struct {
	tt       module.ModuleType
	userI    module.UserI
	data     AttrData
	timedata map[string]int64
	dirty    map[int32]struct{}

	*module.ModuleSaveBase
}

func (this *UserAttr) pack(item *comAttr.Attr) *cs_message.Attr {
	ret := &cs_message.Attr{
		Id:  proto.Int32(item.ID),
		Val: proto.Int64(item.Val),
	}
	if nextTime, ok := this.timedata[timeId2Name(item.ID)]; ok {
		ret.NextRefreshTime = proto.Int64(nextTime)
	}
	return ret
}

func (this *UserAttr) ModuleType() module.ModuleType {
	return module.Attr
}

func (this *UserAttr) ReadOut() *module.ReadOutCommand {
	return &module.ReadOutCommand{
		Table:  "user_module_data",
		Key:    this.userI.GetIDStr(),
		Fields: []string{this.ModuleType().String(), timeField},
		Module: this,
	}
}

func (this *UserAttr) Tick(now time.Time) {
	this.clockTimer()
}

func (this *UserAttr) getLevel() int32 {
	level, _ := this.GetAttr(comAttr.Level)
	return int32(level)
}

func (this *UserAttr) GetAttr(idx int32) (int64, error) {
	if idx > 0 && idx < int32(len(this.data.Attrs)) {
		return this.data.Attrs[idx].Val, nil
	} else {
		return 0, fmt.Errorf("invaild attr")
	}
}

func (this *UserAttr) setAttr(idx int32, val int64, emitEvent bool) int64 {
	ele := this.data.Attrs[idx]
	oldVal := ele.Val
	ele.Val = val

	this.dirty[idx] = struct{}{}
	this.SetDirty(this.ModuleType().String())

	// 事件触发
	if oldVal != val && emitEvent {
		this.userI.EmitEvent(event.EventAttrChange, ele.ID, oldVal, val)
	}
	return oldVal
}

// set接口，选择是否属性变更事件. 返回旧值
func (this *UserAttr) SetAttr(idx int32, val int64, emitEvent bool) (int64, error) {
	if idx > 0 && idx < int32(len(this.data.Attrs)) {
		ele := this.data.Attrs[idx]
		if ele.Val == val {
			return val, nil
		}

		info := comAttr.GetAttrInfo(idx)
		if val < info.Min {
			val = info.Min
		}

		if val > info.Max {
			val = info.Max
		}

		return this.setAttr(idx, val, emitEvent), nil
	} else {
		return 0, fmt.Errorf("invaild attr")
	}
}

// 增量接口，触发属性变更事件
func (this *UserAttr) AddAttr(idx int32, val int64) (int64, error) {
	if idx > 0 && idx < int32(len(this.data.Attrs)) {
		ele := this.data.Attrs[idx]

		oldVal := ele.Val
		newVal := ele.Val + val

		info := comAttr.GetAttrInfo(idx)
		if newVal < info.Min {
			newVal = info.Min
		}

		if newVal > info.Max {
			newVal = info.Max
		}

		if oldVal == newVal {
			return newVal, nil
		}

		this.setAttr(idx, newVal, true)
		return newVal, nil
	} else {
		return 0, fmt.Errorf("invaild attr")
	}
}

func (this *UserAttr) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  "user_module_data",
		Key:    this.userI.GetIDStr(),
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
		Module: this,
	}

	moduleName := this.ModuleType().String()
	for field := range fields {
		name := field.(string)
		switch name {
		case moduleName:
			data, err := json.Marshal(this.data)
			if nil != err {
				return nil
			}
			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  moduleName,
				Value: data,
			})
		case timeField:
			data, _ := json.Marshal(this.timedata)
			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  timeField,
				Value: data,
			})
		}
	}
	return cmd

}

func (this *UserAttr) FlushDirtyToClient() {
	if len(this.dirty) > 0 {
		msg := &cs_message.AttrSyncToC{
			IsAll: proto.Bool(false),
			Attrs: make([]*cs_message.Attr, 0, len(this.dirty)),
		}
		for idx := range this.dirty {
			msg.Attrs = append(msg.Attrs, this.pack(this.data.Attrs[idx]))
		}
		this.dirty = map[int32]struct{}{}
		this.userI.Post(msg)
	}
}

func (this *UserAttr) FlushAllToClient(seqNo ...uint32) {
	msg := &cs_message.AttrSyncToC{
		IsAll: proto.Bool(true),
		Attrs: make([]*cs_message.Attr, 0, len(this.dirty)),
	}
	for i := 1; i <= comAttr.AttrMax; i++ {
		msg.Attrs = append(msg.Attrs, this.pack(this.data.Attrs[i]))
	}
	this.dirty = map[int32]struct{}{}
	this.userI.Post(msg)
}

func (this *UserAttr) Init(fields map[string]*flyfish.Field) error {
	moduleName := this.ModuleType().String()
	// 属性数据
	field, ok := fields[moduleName]
	if ok && len(field.GetBlob()) != 0 {
		var data = AttrData{}
		err := json.Unmarshal(field.GetBlob(), &data)
		if err != nil {
			return fmt.Errorf("unmarshal: %s", err)
		} else {
			this.initOrRepair(&data)
		}
	} else {
		this.initOrRepair(&AttrData{})
	}

	// 时间戳数据
	field, ok = fields[timeField]
	if ok && len(field.GetBlob()) != 0 {
		err := json.Unmarshal(field.GetBlob(), &this.timedata)
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *UserAttr) AfterInitAll() error {
	this.userI.RegisterEvent(event.EventAttrChange, this.EventAttrChange)
	return nil
}

// 初始化或者修复（新增）
func (this *UserAttr) initOrRepair(dbData *AttrData) {
	//log.GetLogger().Infof("%s roleModule attr initOrRepair ", this.userI.GetUserID())

	attrs := make([]*comAttr.Attr, comAttr.AttrMax+1)
	copy(attrs, dbData.Attrs)

	//数据修复
	if len(dbData.Attrs) < comAttr.AttrMax+1 {
		for i := int32(len(dbData.Attrs)); i <= comAttr.AttrMax; i++ {
			attrs[i] = &comAttr.Attr{ID: i, Val: 0}
		}
		this.SetDirty(this.ModuleType().String())
	}
	this.data.Attrs = attrs
}

func (this *UserAttr) Query(arg *cs_message.QueryRoleInfoArg, ret *cs_message.QueryRoleInfoResult) error {
	ids := arg.GetAttrIDs()
	ret.Attrs = make([]*cs_message.Attr, 0, len(ids))
	for _, id := range ids {
		ret.Attrs = append(ret.Attrs, this.pack(this.data.Attrs[id]))
	}
	return nil
}

func init() {
	module.RegisterModule(module.Attr, func(userI module.UserI) module.ModuleI {
		roleAttr := &UserAttr{
			tt:       module.Attr,
			userI:    userI,
			dirty:    map[int32]struct{}{},
			timedata: map[string]int64{},
		}

		roleAttr.ModuleSaveBase = module.NewModuleSaveBase(roleAttr)

		return roleAttr
	})
}
