package quest

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/node_game/event"
	attr2 "initialthree/node/node_game/module/attr"
	Quest2 "initialthree/node/table/excel/DataTable/Quest"
	"initialthree/pkg/json"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	cs_msg "initialthree/protocol/cs/message"
	"time"
)

const (
	timeField = "timedata"
)

var Int2String = map[int32]string{
	int32(enumType.QuestType_MainStory):   "slot0",
	int32(enumType.QuestType_Daily):       "slot1",
	int32(enumType.QuestType_Weekly):      "slot2",
	int32(enumType.QuestType_BigSecret):   "slot3",
	int32(enumType.QuestType_DailyReward): "slot5",
	int32(enumType.QuestType_NewbieGift):  "slot8",
}

type Condition struct {
	DoneTimes int32 `json:"dt"` // 已经完成多少次
	Complete  bool  `json:"c"`  // 完成

	condEvent CondEvent
}

type Quest struct {
	ID         int32             `json:"id"`
	State      cs_msg.QuestState `json:"st"`
	Conditions []*Condition      `json:"cond"`

	registered bool // 是否已经触发过 注册监听
	questTemp  *Quest2.Quest
}

func newQuest(temp *Quest2.Quest) *Quest {
	q := &Quest{
		ID:         temp.ID,
		State:      cs_msg.QuestState_Running,
		Conditions: make([]*Condition, 0, len(temp.Condition)),
		questTemp:  temp,
	}

	for i := 0; i < len(temp.Conditions()); i++ {
		q.Conditions = append(q.Conditions, &Condition{})
	}

	return q
}

func (this *Quest) GetConfig() *Quest2.Quest {
	return this.questTemp
}

func (this *Quest) IsFinished() bool {
	for _, v := range this.Conditions {
		if !v.Complete {
			return false
		}
	}
	return true
}

func (this *Quest) Pack() *cs_msg.Quest {
	q := &cs_msg.Quest{
		QuestID:         proto.Int32(this.ID),
		State:           this.State.Enum(),
		QuestConditions: make([]*cs_msg.QuestCondition, 0, len(this.Conditions)),
	}

	for _, v := range this.Conditions {
		q.QuestConditions = append(q.QuestConditions, &cs_msg.QuestCondition{
			DoneTimes: proto.Int32(v.DoneTimes),
			Complete:  proto.Bool(v.Complete),
		})
	}
	return q
}

type UserQuest struct {
	userI        module.UserI
	quests       map[int32]map[int32]*Quest // type -> (id -> quest)
	timedata     map[string]int64
	questDirty   map[int32]struct{}
	needRegister map[int32]*Quest // 需要注册的任务，由tick触发
	*module.ModuleSaveBase
}

func (this *UserQuest) GetQuest(id int32) *Quest {
	temp := Quest2.GetID(id)
	if temp == nil {
		return nil
	}
	return this.quests[temp.TypeEnum][id]
}

func (this *UserQuest) AddQuest(questID int32) {
	this.addQuest(questID)
}

func (this *UserQuest) addQuest(questID int32) {
	temp := Quest2.GetID(questID)
	if temp != nil {
		quests := this.quests[temp.TypeEnum]
		if _, ok := quests[questID]; !ok {
			q := newQuest(temp)
			quests[questID] = q
			this.SetDirty(temp.TypeEnum)
			this.questDirty[questID] = struct{}{}

			this.needRegister[questID] = q
		}
	}
}

func (this *UserQuest) tryAllCondComplete(q *Quest) {
	//log.GetLogger().Debugln("tryCondComplete", this.userI.GetID(), q.ID)
	if q.State == cs_msg.QuestState_Running && q.IsFinished() {
		zaplogger.GetSugar().Debugf("tryCondComplete %s finished quest %d", this.userI.GetUserID(), q.ID)
		q.State = cs_msg.QuestState_Finished

		// 添加后续任务
		if len(q.questTemp.UnlockQuestsArray) > 0 {
			for _, v := range q.questTemp.UnlockQuestsArray {
				this.addQuest(v.QuestID)
			}
		}

		// 清理注册事件
		for _, cond := range q.Conditions {
			if cond.condEvent != nil {
				cond.condEvent.UnTrigger()
			}
		}

		this.SetDirty(q.questTemp.TypeEnum)
		this.questDirty[q.ID] = struct{}{}
	}
}

func (this *UserQuest) End(q *Quest) {
	q.State = cs_msg.QuestState_End
	this.SetDirty(q.questTemp.TypeEnum)
	this.questDirty[q.ID] = struct{}{}
}

func (this *UserQuest) AfterInitAll() error {
	// 七日任务清理
	attrModule := this.userI.GetSubModule(module.Attr).(*attr2.UserAttr)
	startTime, _ := attrModule.GetAttr(attr.NewbieGiftStartTime)
	if startTime != 0 {
		endTime, _ := attrModule.GetAttr(attr.NewbieGiftEndTime)
		if time.Now().Unix() > endTime {
			// 过期清零
			this.clear(enumType.QuestType_NewbieGift)
		}
	}
	return nil
}

func (this *UserQuest) tryRegister() {
	for _, q := range this.needRegister {
		this.registerEvent(q)
	}
}

func (this *UserQuest) Tick(now time.Time) {
	this.clockTimer()
	this.tryRegister()
}

func (this *UserQuest) FlushAllToClient(seqNo ...uint32) {
	msg := &cs_msg.QuestSyncToC{
		IsAll:  proto.Bool(true),
		Quests: make([]*cs_msg.Quest, 0, 16),
	}
	for _, hmap := range this.quests {
		for _, q := range hmap {
			msg.Quests = append(msg.Quests, q.Pack())
		}
	}

	this.questDirty = map[int32]struct{}{}
	this.userI.Post(msg)
}

func (this *UserQuest) FlushDirtyToClient() {
	if len(this.questDirty) > 0 {

		msg := &cs_msg.QuestSyncToC{
			IsAll:  proto.Bool(false),
			Quests: make([]*cs_msg.Quest, 0, len(this.questDirty)),
		}
		for id := range this.questDirty {
			q := this.GetQuest(id)
			msg.Quests = append(msg.Quests, q.Pack())

			// 当前任务已经完结，并且有后需，可删除。
			// 不删除的情况，没有后续任务，保留该任务， 策划后又添加新任务，可直接索引出
			//if q.State == cs_msg.QuestState_End /*&& len(tq.FollowedQuestArray) != 0 */ {
			//	zaplogger.GetSugar().Debugf("%v %d quest end delete", this.userI.GetID(), q.ID)
			//	delete(this.quests[q.questTemp.TypeEnum], q.ID)
			//	this.SetDirty(q.questTemp.TypeEnum)
			//}
		}

		this.questDirty = map[int32]struct{}{}
		this.userI.Post(msg)
	}
}

func (this *UserQuest) ModuleType() module.ModuleType {
	return module.Quest
}

func (this *UserQuest) Init(fields map[string]*flyfish.Field) error {
	for tt, s := range Int2String {
		field, ok := fields[s]

		if ok && len(field.GetBlob()) != 0 {
			var quests map[int32]*Quest
			err := json.Unmarshal(field.GetBlob(), &quests)
			if err != nil {
				return fmt.Errorf("unmarshal: %s", err)
			}

			// 绑定配置
			for _, q := range quests {
				temp := Quest2.GetID(q.ID)
				if temp == nil || tt != temp.TypeEnum ||
					(q.State == cs_msg.QuestState_Running && len(temp.Condition) != len(q.Conditions)) {
					// 任务正在执行中 需要绑定条件的配置
					//zaplogger.GetSugar().Errorf("init quest %d config is not equal", q.ID)
					// 移除任务
					delete(quests, q.ID)
					continue
				}
				q.questTemp = temp

				if q.State == cs_msg.QuestState_Running {
					this.needRegister[q.ID] = q
				}
			}

			this.quests[tt] = quests
		} else {
			this.quests[tt] = map[int32]*Quest{}
			this.SetDirty(s)
		}
	}

	field, ok := fields[timeField]
	if ok && field.GetBlob() != nil {
		err := json.Unmarshal(field.GetBlob(), &this.timedata)
		if err != nil {
			return fmt.Errorf("unmarshal: %s", err)
		}
	}

	return nil
}

func (this *UserQuest) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Fields: make([]string, 0, len(Int2String)+1),
		Module: this,
	}

	for _, field := range Int2String {
		cmd.Fields = append(cmd.Fields, field)
	}
	cmd.Fields = append(cmd.Fields, timeField)

	return cmd
}

func (this *UserQuest) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]*module.WriteBackFiled, 0, len(fields)+1),
	}

	for field := range fields {
		switch field.(type) {
		case int32:
			tt := field.(int32)
			data, _ := json.Marshal(this.quests[tt])
			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  Int2String[tt],
				Value: data,
			})
		case string:
			data, _ := json.Marshal(this.timedata)
			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  timeField,
				Value: data,
			})
		}
	}

	return cmd
}

func init() {
	module.RegisterModule(module.Quest, func(userI module.UserI) module.ModuleI {
		q := &UserQuest{
			userI:        userI,
			quests:       map[int32]map[int32]*Quest{},
			questDirty:   map[int32]struct{}{},
			timedata:     map[string]int64{},
			needRegister: map[int32]*Quest{},
		}
		q.ModuleSaveBase = module.NewModuleSaveBase(q)

		q.userI.RegisterEvent(event.EventBigSecret, q.EventResetBigSecret)

		return q
	})
}
