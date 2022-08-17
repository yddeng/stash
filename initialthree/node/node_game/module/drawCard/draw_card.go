package drawCard

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/pkg/json"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/protocol/cs/message"
	"time"
)

const (
	tableName       = "drawcard"
	historyField    = "history"
	guaranteeField  = "guarantee"
	poolIndexField  = "poolindex"
	dailyTimesField = "dailytimes"
	timeField       = "timedata"
)

var tableFields = []string{poolIndexField, guaranteeField, historyField, dailyTimesField, timeField}

type DrawCard struct {
	userI      module.UserI
	guarantee  map[int32]int32                    // id -> value
	poolIndex  map[int32]int32                    // lidId -> value
	history    map[int32][]*message.DrawCardAward // lidId ->
	dailyTimes map[int32]int32                    // libId -> times 每日抽卡限制

	timedata map[string]int64

	guaranteeDirty  map[int32]struct{}
	poolIndexDirty  map[int32]struct{}
	dailyTimesDirty map[int32]struct{}
	*module.ModuleSaveBase
}

func (this *DrawCard) GetHistory(libID int32) []*message.DrawCardAward {
	return this.history[libID]
}

func (this *DrawCard) AddHistory(libID int32, awardList []*message.DrawCardAward) {
	history := this.history[libID]
	if history == nil {
		history = make([]*message.DrawCardAward, 0, 10)
	}
	history = append(awardList, history...)
	if len(history) > 10 { // 最多存放 10 条记录
		history = history[:10]
	}
	this.history[libID] = history
	this.SetDirty(historyField)
}

func (this *DrawCard) GetGuaranteeCount(guaranteeID int32) int32 {
	return this.guarantee[guaranteeID]
}

func (this *DrawCard) SetGuarantee(guaranteeID, guaranteeCount int32) {
	this.guarantee[guaranteeID] = guaranteeCount
	this.guaranteeDirty[guaranteeID] = struct{}{}
	this.SetDirty(guaranteeField)
}

func (this *DrawCard) GetDailyTimes(guaranteeID int32) int32 {
	return this.dailyTimes[guaranteeID]
}
func (this *DrawCard) AddDailyTimes(guaranteeID, times int32) {
	_times := this.dailyTimes[guaranteeID]
	this.dailyTimes[guaranteeID] = _times + times
	this.dailyTimesDirty[guaranteeID] = struct{}{}
	this.SetDirty(dailyTimesField)
}

func (this *DrawCard) GetPoolIndex(libID int32) int32 {
	return this.poolIndex[libID]
}

func (this *DrawCard) SetPoolIndex(libID, poolIndex int32) {
	this.poolIndex[libID] = poolIndex
	this.poolIndexDirty[libID] = struct{}{}
	this.SetDirty(poolIndexField)
}

func (this *DrawCard) ModuleType() module.ModuleType {
	return module.DrawCard
}

func (this *DrawCard) Init(fields map[string]*flyfish.Field) error {
	for _, name := range tableFields {
		field, ok := fields[name]
		if ok && len(field.GetBlob()) != 0 {
			var err error
			switch name {
			case historyField:
				err = json.Unmarshal(field.GetBlob(), &this.history)
			case guaranteeField:
				err = json.Unmarshal(field.GetBlob(), &this.guarantee)
			case poolIndexField:
				err = json.Unmarshal(field.GetBlob(), &this.poolIndex)
			case dailyTimesField:
				err = json.Unmarshal(field.GetBlob(), &this.dailyTimes)
			case timeField:
				err = json.Unmarshal(field.GetBlob(), &this.timedata)
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s initDrawCard name %s err %s ", this.userI.GetUserID(), name, err)
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}
	return nil
}

func (this *DrawCard) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Fields: tableFields,
		Module: this,
	}
	return cmd
}

func (this *DrawCard) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
		Module: this,
	}

	for field := range fields {
		name := field.(string)
		var data []byte
		switch name {
		case poolIndexField:
			data, _ = json.Marshal(this.poolIndex)
		case guaranteeField:
			data, _ = json.Marshal(this.guarantee)
		case historyField:
			data, _ = json.Marshal(this.history)
		case dailyTimesField:
			data, _ = json.Marshal(this.dailyTimes)
		case timeField:
			data, _ = json.Marshal(this.timedata)
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

func (this *DrawCard) Tick(now time.Time) {
	this.clockTimer()
}

func (this *DrawCard) FlushDirtyToClient() {
	if len(this.poolIndexDirty) > 0 || len(this.guaranteeDirty) > 0 || len(this.dailyTimesDirty) > 0 {
		msg := &message.DrawCardSyncToC{
			IsAll:      proto.Bool(false),
			PoolIndex:  make([]*message.DrawCardPool, 0, len(this.poolIndexDirty)),
			Guarantee:  make([]*message.DrawCardGuarantee, 0, len(this.guaranteeDirty)),
			DailyTimes: make([]*message.DrawCardDailyTimes, 0, len(this.dailyTimesDirty)),
		}
		for id := range this.poolIndexDirty {
			msg.PoolIndex = append(msg.PoolIndex, &message.DrawCardPool{
				LibID:     proto.Int32(id),
				PoolIndex: proto.Int32(this.GetPoolIndex(id)),
			})
		}
		for id := range this.dailyTimesDirty {
			msg.DailyTimes = append(msg.DailyTimes, &message.DrawCardDailyTimes{
				LibID: proto.Int32(id),
				Times: proto.Int32(this.GetDailyTimes(id)),
			})
		}
		for id := range this.guaranteeDirty {
			msg.Guarantee = append(msg.Guarantee, &message.DrawCardGuarantee{
				GuaranteeID:    proto.Int32(id),
				GuaranteeCount: proto.Int32(this.GetGuaranteeCount(id)),
			})
		}

		this.userI.Post(msg)
		this.poolIndexDirty = map[int32]struct{}{}
		this.guaranteeDirty = map[int32]struct{}{}
		this.dailyTimesDirty = map[int32]struct{}{}
	}
}

func (this *DrawCard) FlushAllToClient(seqNo ...uint32) {
	msg := &message.DrawCardSyncToC{
		IsAll:      proto.Bool(true),
		PoolIndex:  make([]*message.DrawCardPool, 0, len(this.poolIndex)),
		Guarantee:  make([]*message.DrawCardGuarantee, 0, len(this.guarantee)),
		DailyTimes: make([]*message.DrawCardDailyTimes, 0, len(this.dailyTimes)),
	}
	for id := range this.poolIndex {
		msg.PoolIndex = append(msg.PoolIndex, &message.DrawCardPool{
			LibID:     proto.Int32(id),
			PoolIndex: proto.Int32(this.GetPoolIndex(id)),
		})
	}
	for id := range this.dailyTimes {
		msg.DailyTimes = append(msg.DailyTimes, &message.DrawCardDailyTimes{
			LibID: proto.Int32(id),
			Times: proto.Int32(this.GetDailyTimes(id)),
		})
	}
	for id := range this.guarantee {
		msg.Guarantee = append(msg.Guarantee, &message.DrawCardGuarantee{
			GuaranteeID:    proto.Int32(id),
			GuaranteeCount: proto.Int32(this.GetGuaranteeCount(id)),
		})
	}
	this.userI.Post(msg)
	this.poolIndexDirty = map[int32]struct{}{}
	this.guaranteeDirty = map[int32]struct{}{}
	this.dailyTimesDirty = map[int32]struct{}{}
}

func init() {
	module.RegisterModule(module.DrawCard, func(userI module.UserI) module.ModuleI {
		m := &DrawCard{
			userI:           userI,
			history:         map[int32][]*message.DrawCardAward{},
			guarantee:       map[int32]int32{},
			poolIndex:       map[int32]int32{},
			guaranteeDirty:  map[int32]struct{}{},
			poolIndexDirty:  map[int32]struct{}{},
			dailyTimes:      map[int32]int32{},
			dailyTimesDirty: map[int32]struct{}{},
			timedata:        map[string]int64{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
