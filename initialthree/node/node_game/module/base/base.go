package base

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/module"
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"time"
)

type userBaseData struct {
	UserID        string  `json:"user_id"`
	Name          string  `json:"name"` // 角色名
	Sex           int32   `json:"sex"`
	BuildTime     int64   `json:"build_time"` // 创建账号的时间
	Birthday      string  `json:"birthday"`
	Signature     string  `json:"signature"`
	CharacterList []int32 `json:"character_list"`
	Portrait      int32   `json:"portrait"`
	PortraitFrame int32   `json:"portrait_frame"`
	Card          int32   `json:"card"`
	UIDCounter    uint32  `json:"uc"`
}

type UserBase struct {
	userI module.UserI
	data  userBaseData
	dirty bool
	*module.ModuleSaveBase
}

func (this *UserBase) setDirty() {
	this.SetDirty(this.ModuleType().String())
	this.dirty = true
}

func (this *UserBase) GetUserID() string {
	return this.data.UserID
}

func (this *UserBase) GenUID() uint32 {
	this.data.UIDCounter++
	this.SetDirty(this.ModuleType().String())
	return this.data.UIDCounter
}

func (this *UserBase) GetCard() int32 {
	return this.data.Card
}

func (this *UserBase) GetBuildTime() int64 {
	return this.data.BuildTime
}

func (this *UserBase) SetCard(v int32) {
	this.data.Card = v
	this.setDirty()
}

func (this *UserBase) SetName(name string) {
	this.data.Name = name
	this.setDirty()
}

func (this *UserBase) GetBirthday() string {
	return this.data.Birthday
}

func (this *UserBase) SetBirthday(v string) {
	this.data.Birthday = v
	this.setDirty()
}

func (this *UserBase) GetCharacterList() []int32 {
	return this.data.CharacterList
}

func (this *UserBase) SetPortraitFrame(v int32) {
	this.data.PortraitFrame = v
	this.setDirty()
}

func (this *UserBase) GetSignature() string {
	return this.data.Signature
}

func (this *UserBase) SetSignature(v string) {
	this.data.Signature = v
	this.setDirty()
}

func (this *UserBase) GetPortraitFrame() int32 {
	return this.data.PortraitFrame
}

func (this *UserBase) SetCharacterList(v []int32) {
	this.data.CharacterList = v
	this.setDirty()
}

func (this *UserBase) GetName() string {
	return this.data.Name
}

func (this *UserBase) GetSex() int32 {
	return this.data.Sex
}

func (this *UserBase) SetSex(v int32) {
	this.data.Sex = v
	this.setDirty()
}

func (this *UserBase) GetPortrait() int32 {
	return this.data.Portrait
}

func (this *UserBase) SetPortrait(v int32) {
	this.data.Portrait = v
	this.setDirty()
}

func (this *UserBase) ModuleType() module.ModuleType {
	return module.Base
}

func (this *UserBase) Init(fields map[string]*client.Field) error {
	field, ok := fields[this.ModuleType().String()]
	if ok {
		if err := json.Unmarshal(field.GetBlob(), &this.data); err != nil {
			return fmt.Errorf("unmarshal: %s", err)
		}

	} else {
		this.SetDirty(this.ModuleType().String())
	}

	return nil
}

func (this *UserBase) ReadOut() *module.ReadOutCommand {
	return &module.ReadOutCommand{
		Table:  "user_module_data",
		Key:    this.userI.GetIDStr(),
		Fields: []string{this.ModuleType().String()},
		Module: this,
	}
}

func (this *UserBase) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	data, err := json.Marshal(this.data)
	if nil != err {
		return nil
	}

	cmd := &module.WriteBackCommand{
		Table: "user_module_data",
		Key:   this.userI.GetIDStr(),
		Fields: []*module.WriteBackFiled{{
			Name:  this.ModuleType().String(),
			Value: data,
		}},
		Module: this,
	}

	return cmd
}

func (this *UserBase) pack() *message.BaseSyncToC {
	return &message.BaseSyncToC{
		UserID:        proto.String(this.GetUserID()),
		GameID:        proto.Uint64(this.userI.GetID()),
		Name:          proto.String(this.data.Name),
		Sex:           proto.Int32(this.data.Sex),
		Birthday:      proto.String(this.data.Birthday),
		Signature:     proto.String(this.data.Signature),
		CharacterList: this.data.CharacterList,
		Portrait:      proto.Int32(this.data.Portrait),
		PortraitFrame: proto.Int32(this.data.PortraitFrame),
		Card:          proto.Int32(this.data.Card),
	}
}

func (this *UserBase) FlushDirtyToClient() {
	if this.dirty {
		msg := this.pack()
		this.userI.Post(msg)
		this.dirty = false
	}
}

func (this *UserBase) FlushAllToClient(seqNo ...uint32) {
	msg := this.pack()
	this.dirty = false
	this.userI.Post(msg)

}

func (this *UserBase) Tick(time.Time) {
}

func (this *UserBase) Query(arg *message.QueryRoleInfoArg, ret *message.QueryRoleInfoResult) error {
	ret.UserID = proto.String(this.GetUserID())
	ret.BaseResp = this.pack()
	return nil
}

func init() {
	module.RegisterModule(module.Base, func(u module.UserI) module.ModuleI {
		m := &UserBase{
			userI: u,
			data: userBaseData{
				UserID:        u.GetUserID(),
				BuildTime:     time.Now().Unix(),
				Portrait:      1,
				PortraitFrame: 1,
				Card:          1,
			},
		}
		m.ModuleSaveBase = module.NewModuleSaveBase(m)

		return m
	})
}
