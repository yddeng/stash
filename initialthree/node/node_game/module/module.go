package module

import (
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/pkg/event"
	"initialthree/protocol/cs/message"
	"time"
)

type ModuleType int32

const (
	Invaild         = ModuleType(0)
	Base            = ModuleType(1)
	Attr            = ModuleType(2)
	Weapon          = ModuleType(3)
	Character       = ModuleType(4)
	Backpack        = ModuleType(5)
	Assets          = ModuleType(6)
	MapData         = ModuleType(7)
	Equip           = ModuleType(8)
	Quest           = ModuleType(9)
	Shop            = ModuleType(10)
	MainDungeons    = ModuleType(11)
	RankData        = ModuleType(12)
	MaterialDungeon = ModuleType(13)
	DrawCard        = ModuleType(14)
	ScarsIngrain    = ModuleType(15)
	RewardQuest     = ModuleType(16)
	WorldQuest      = ModuleType(17)
	TrialDungeon    = ModuleType(18)
	Mail            = ModuleType(19)
	Sign            = ModuleType(20)
	BigSecret       = ModuleType(21)
	Talent          = ModuleType(22)
	End             = ModuleType(23)
)

var typeToStr = map[ModuleType]string{
	Base:            "base",
	Attr:            "attr",
	Weapon:          "weapon",
	Character:       "character",
	Backpack:        "backpack",
	Assets:          "assets",
	MapData:         "mapdata",
	Equip:           "equip",
	Quest:           "quest",
	Shop:            "shop",
	MainDungeons:    "maindungeons",
	RankData:        "rankdata",
	MaterialDungeon: "materialdungeon",
	DrawCard:        "drawcard",
	ScarsIngrain:    "scarsingrain",
	RewardQuest:     "rewardquest",
	WorldQuest:      "worldquest",
	TrialDungeon:    "trialdungeon",
	Mail:            "mail",
	Sign:            "sign",
	BigSecret:       "bigsecret",
	Talent:          "talent",
}

func (t ModuleType) String() string {
	return typeToStr[t]
}

func (t ModuleType) Int32() int32 {
	return int32(t)
}

type ReadOutCommand struct {
	Table  string
	Key    string
	Fields []string
	Module ModuleI
}

type WriteBackFiled struct {
	Name  string
	Value interface{}
}

type WriteBackCommand struct {
	Table  string
	Key    string
	Fields []*WriteBackFiled
	Module ModuleI
}

type UserI interface {
	Post(proto.Message)
	Reply(seqNo uint32, msg proto.Message)
	GetSubModule(moduleType ModuleType) ModuleI
	FlushAllDirtyToClient()
	FlushAllToClient()
	GetID() uint64
	GetIDStr() string
	GetUserID() string
	GetLevel() int32
	StatusOk() bool

	EventI

	SendMail(mails []*message.Mail)
}

type EventI interface {
	RegisterEvent(event interface{}, fn interface{}) event.Handle
	RegisterEventOnce(event interface{}, fn interface{}) event.Handle
	UnRegisterEvent(h event.Handle)
	ClearEvent(event interface{})
	EmitEvent(event interface{}, args ...interface{})
}

type ModuleI interface {
	ModuleType() ModuleType
	Init(map[string]*flyfish.Field) error

	//存储相关
	ReadOut() *ReadOutCommand
	WriteCommand(fields map[interface{}]struct{}) *WriteBackCommand
	ModuleDBSaveI

	//同步客户端相关
	FlushDirtyToClient()
	FlushAllToClient(seqNo ...uint32)
	Tick(time.Time)
}

// 所有模块初始化成功后执行的行为
type AfterModuleInitAll interface {
	AfterInitAll() error
}

type Creator func(UserI) ModuleI

var Modules = map[ModuleType]Creator{}

//func GetCreatorByName(name string) func(UserI) ModuleI {
//	return Modules[strToType[name]]
//}
//
//func GetModuleTypeByName(name string) ModuleType {
//	return strToType[name]
//}

func RegisterModule(tt ModuleType, creator func(UserI) ModuleI) {
	if tt <= Invaild || tt >= End {
		panic("invaild SubmoduleType")
	}

	if tt.String() == "" {
		panic("invaild submodule name")
	}

	if _, ok := Modules[tt]; ok {
		panic("duplicate submodule")
	}

	Modules[tt] = creator
}

type ModuleDBSaveI interface {
	SetDirty(fields ...interface{})
	IsDirty() bool
	WriteBackRet(ok bool) // 调用存储的结果
	WriteBack() *WriteBackCommand
}

type ModuleSaveBase struct {
	moduleI      ModuleI
	savingFields map[interface{}]struct{}
	dirtyFields  map[interface{}]struct{}
}

func NewModuleSaveBase(moduleI ModuleI) *ModuleSaveBase {
	return &ModuleSaveBase{
		moduleI:      moduleI,
		savingFields: map[interface{}]struct{}{},
		dirtyFields:  map[interface{}]struct{}{},
	}
}

func (this *ModuleSaveBase) SetDirty(fields ...interface{}) {
	for _, k := range fields {
		this.dirtyFields[k] = struct{}{}
	}
}

func (this *ModuleSaveBase) IsDirty() bool {
	return len(this.dirtyFields) != 0
}

func (this *ModuleSaveBase) WriteBackRet(ok bool) {
	if !ok {
		for k := range this.savingFields {
			this.dirtyFields[k] = struct{}{}
		}
	}
	this.savingFields = map[interface{}]struct{}{}
}

func (this *ModuleSaveBase) WriteBack() *WriteBackCommand {
	cmd := this.moduleI.WriteCommand(this.dirtyFields)
	if cmd == nil || len(cmd.Fields) == 0 {
		return nil
	}

	for k := range this.dirtyFields {
		this.savingFields[k] = struct{}{}
	}
	this.dirtyFields = map[interface{}]struct{}{}

	return cmd
}

const (
	DailyTimeName   = "dailyTime"
	WeeklyTimeName  = "weeklyTime"
	MonthlyTimeName = "monthlyTime"
)

func CalDailyTime() time.Time {
	rt := Global.Get().GetDailyRefreshTime()
	dailyTime := timeDisposal.CalcLatestTimeAfter(rt.Hour, rt.Minute, 0)
	return dailyTime
}

func CalWeeklyTime() time.Time {
	rt := Global.Get().GetWeeklyRefreshTime()
	weeklyTime := timeDisposal.CalcLatestWeekTimeAfter(rt.Weekday, rt.Hour, rt.Minute, 0)
	return weeklyTime
}

func CalMonthlyTime() time.Time {
	rt := Global.Get().GetMonthlyRefreshTime()
	monthlyTime := timeDisposal.CalcLatestMonthlyTimeAfter(int(rt.Day), rt.Hour, rt.Minute, 0)
	return monthlyTime
}
