package bigSecret

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/global/bigSecret"
	"initialthree/node/table/excel/ConstTable/BigSecret"
	"initialthree/zaplogger"
	"math/rand"

	"initialthree/node/node_game/module"
	"initialthree/protocol/cs/message"
	"time"
)

const (
	tableName     = "bigsecret"
	dataField     = "data"
	weaknessField = "weakness"
	passedFiled   = "pass"
	shopField     = "shop"
	timeField     = "timedata"
)

var tableFields = []string{dataField, weaknessField, passedFiled, shopField, timeField}

type Data struct {
	MaxLevel int32
	KeyCount int32

	WeaknessRefreshTimes int32

	BlessingLevel int32
	BlessingCount int32

	RankID int32 // 赛季ID
}

type BigSecretDungeon struct {
	*module.ModuleSaveBase
	userI module.UserI

	data         Data
	weakness     map[int32]*message.BigSecretWeakness // 每日刷新的弱点
	passed       map[int32]struct{}                   // 已通关的
	shopBuyTimes map[int32]int32                      // 商店道具购买次数
	timedata     map[string]int64

	dirtyData     bool
	dirtyPassed   map[int32]struct{}
	dirtyWeakness map[int32]*message.BigSecretWeakness
	dirtyShop     map[int32]int32

	competition *bigSecret.Data
	resetDate   int64 // 重置时间 到天，避免重复
}

func (this *BigSecretDungeon) GetKeyCount() int32 {
	return this.data.KeyCount
}

func (this *BigSecretDungeon) resetData() {
	def := BigSecret.GetID(1)
	if this.data.MaxLevel > def.ResetLv {
		this.data.MaxLevel = def.ResetLv
	} else {
		this.data.MaxLevel = 1
	}
	this.data.KeyCount = 1
	this.data.BlessingCount = 0
	this.data.BlessingLevel = 1

	this.dirtyData = true
	this.SetDirty(dataField)

	year, mon, day := time.Now().Date()
	this.resetDate = time.Date(year, mon, day, 0, 0, 0, 0, time.Local).Unix()
}

func (this *BigSecretDungeon) randomWeakness(level int32) *message.BigSecretWeakness {
	def := BigSecret.GetID(1)
	length := len(def.WeaknessArray)

	wk := &message.BigSecretWeakness{
		Level:    proto.Int32(level),
		Weakness: make([]int32, 0, 2),
	}
	if length > 2 {
		exist := map[int]bool{}
		for {
			if len(exist) >= 2 {
				break
			}
			idx := rand.Int() % length
			exist[idx] = true
		}

		for idx := range exist {
			wk.Weakness = append(wk.Weakness, def.WeaknessArray[idx].ID)
		}
	} else {
		for _, v := range def.WeaknessArray {
			wk.Weakness = append(wk.Weakness, v.ID)
		}
	}

	return wk
}

func (this *BigSecretDungeon) GetWeakness(level int32) *message.BigSecretWeakness {
	if wk, ok := this.weakness[level]; ok {
		return wk
	} else {
		wk = this.randomWeakness(level)
		this.weakness[level] = wk

		this.dirtyWeakness[level] = wk
		this.SetDirty(weaknessField)
		return wk
	}
}

func (this *BigSecretDungeon) GetWeaknessRefreshTimes() int32 {
	return this.data.WeaknessRefreshTimes
}

func (this *BigSecretDungeon) WeaknessRefresh(level int32) *message.BigSecretWeakness {
	wk := this.randomWeakness(level)
	this.weakness[level] = wk

	this.dirtyWeakness[level] = wk
	this.SetDirty(weaknessField)

	this.data.WeaknessRefreshTimes++
	this.dirtyData = true
	this.SetDirty(dataField)
	return wk
}

func (this *BigSecretDungeon) BestPassed() int32 {
	level := int32(0)
	for lv := range this.passed {
		if lv > level {
			level = lv
		}
	}
	return level
}

func (this *BigSecretDungeon) Passed(level int32) bool {
	if _, ok := this.passed[level]; ok {
		return true
	}
	return false
}

func (this *BigSecretDungeon) Unlocked(level int32) bool {
	if level <= this.data.MaxLevel {
		return true
	}
	return false
}

func (this *BigSecretDungeon) Pass(level int32, useTime int32) {
	def := BigSecret.GetID(1)
	unlock := int32(0)
	for _, v := range def.PassTimeUnlock {
		limit := v.PassTime * 60 // 分 -> 秒
		if useTime <= limit {
			unlock = v.Unlock
			break
		}
	}

	this.data.MaxLevel += unlock
	this.dirtyData = true
	this.SetDirty(dataField)

	this.passed[level] = struct{}{}
	this.dirtyPassed[level] = struct{}{}
	this.SetDirty(passedFiled)
}

func (this *BigSecretDungeon) GetBlessingCount() int32 {
	return this.data.BlessingCount
}
func (this *BigSecretDungeon) GetBlessingLevel() int32 {
	return this.data.BlessingLevel
}

func (this *BigSecretDungeon) BlessingLvUp(cost int32) {
	this.data.BlessingCount -= cost
	this.data.BlessingLevel += 1
	this.dirtyData = true
	this.SetDirty(dataField)
}

func (this *BigSecretDungeon) ModuleType() module.ModuleType {
	return module.BigSecret
}

func (this *BigSecretDungeon) Init(fields map[string]*flyfish.Field) error {
	for _, name := range tableFields {
		field, ok := fields[name]
		if ok && len(field.GetBlob()) != 0 {
			var err error
			switch name {
			case dataField:
				err = json.Unmarshal(field.GetBlob(), &this.data)
			case weaknessField:
				err = json.Unmarshal(field.GetBlob(), &this.weakness)
			case passedFiled:
				err = json.Unmarshal(field.GetBlob(), &this.passed)
			case timeField:
				err = json.Unmarshal(field.GetBlob(), &this.timedata)
			case shopField:
				err = json.Unmarshal(field.GetBlob(), &this.shopBuyTimes)
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s initBigSecretDungeon name %s err %s ", this.userI.GetUserID(), name, err)
				return fmt.Errorf("unmarshal: %s", err)
			}
		} else {
			switch name {
			case dataField:
				this.resetData()
			}
		}
	}
	return nil
}

func (this *BigSecretDungeon) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Fields: tableFields,
		Module: this,
	}
	return cmd
}

func (this *BigSecretDungeon) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
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
		case dataField:
			data, _ = json.Marshal(this.data)
		case timeField:
			data, _ = json.Marshal(this.timedata)
		case weaknessField:
			data, _ = json.Marshal(this.weakness)
		case passedFiled:
			data, _ = json.Marshal(this.passed)
		case shopField:
			data, _ = json.Marshal(this.shopBuyTimes)
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

func (this *BigSecretDungeon) Tick(now time.Time) {
	this.checkData()
	this.clockTimer()
}

func (this *BigSecretDungeon) checkData() {
	data := bigSecret.GetData()
	if data != nil && data != this.competition {
		if this.data.RankID != data.RankID {
			// 重置赛季
			this.data.RankID = data.RankID
			this.resetData()

			cfg, _ := bigSecret.ParseRankID(this.data.RankID)
			this.userI.EmitEvent(event.EventBigSecret, cfg)
		}

		this.competition = data
		this.flushCompetition()
	}
}

func (this *BigSecretDungeon) flushCompetition() {
	msg := &message.BigSecretCompetitionSyncToC{
		BeginTime: proto.Int64(this.competition.BeginTime),
		EndTime:   proto.Int64(this.competition.EndTime),
		RankID:    proto.Int32(this.competition.RankID),
		Config:    proto.Int32(this.competition.RankIdx),
	}
	this.userI.Post(msg)
}

func (this *BigSecretDungeon) FlushDirtyToClient() {
	if this.dirtyData || len(this.dirtyWeakness) > 0 || len(this.dirtyPassed) > 0 {
		msg := &message.BigSecretSyncToC{
			IsAll:    proto.Bool(false),
			Passed:   make([]int32, 0, len(this.dirtyPassed)),
			Weakness: make([]*message.BigSecretWeakness, 0, len(this.dirtyWeakness)),
		}

		for _, wk := range this.dirtyWeakness {
			msg.Weakness = append(msg.Weakness, wk)
		}
		for id := range this.dirtyPassed {
			msg.Passed = append(msg.Passed, id)
		}
		if this.dirtyData {
			msg.Data = &message.BigSecretData{
				MaxLevel:             proto.Int32(this.data.MaxLevel),
				KeyCount:             proto.Int32(this.data.KeyCount),
				WeaknessRefreshTimes: proto.Int32(this.data.WeaknessRefreshTimes),
				BlessingLevel:        proto.Int32(this.data.BlessingLevel),
				BlessingCount:        proto.Int32(this.data.BlessingCount),
			}
		}

		this.userI.Post(msg)
		this.dirtyData = false
		this.dirtyWeakness = map[int32]*message.BigSecretWeakness{}
		this.dirtyPassed = map[int32]struct{}{}
	}

}
func (this *BigSecretDungeon) FlushAllToClient(seqNo ...uint32) {
	msg := &message.BigSecretSyncToC{
		IsAll: proto.Bool(true),
		Data: &message.BigSecretData{
			MaxLevel:             proto.Int32(this.data.MaxLevel),
			KeyCount:             proto.Int32(this.data.KeyCount),
			WeaknessRefreshTimes: proto.Int32(this.data.WeaknessRefreshTimes),
			BlessingLevel:        proto.Int32(this.data.BlessingLevel),
			BlessingCount:        proto.Int32(this.data.BlessingCount),
		},
		Passed:   make([]int32, 0, len(this.passed)),
		Weakness: make([]*message.BigSecretWeakness, 0, len(this.weakness)),
	}

	for _, wk := range this.weakness {
		msg.Weakness = append(msg.Weakness, wk)
	}

	for id := range this.passed {
		msg.Passed = append(msg.Passed, id)
	}

	this.userI.Post(msg)
	this.dirtyData = false
	this.dirtyWeakness = map[int32]*message.BigSecretWeakness{}
	this.dirtyPassed = map[int32]struct{}{}

	if this.competition != nil {
		this.flushCompetition()
	}
}

func init() {
	module.RegisterModule(module.BigSecret, func(userI module.UserI) module.ModuleI {
		m := &BigSecretDungeon{
			userI:         userI,
			passed:        map[int32]struct{}{},
			weakness:      map[int32]*message.BigSecretWeakness{},
			shopBuyTimes:  map[int32]int32{},
			timedata:      map[string]int64{},
			dirtyWeakness: map[int32]*message.BigSecretWeakness{},
			dirtyPassed:   map[int32]struct{}{},
			dirtyShop:     map[int32]int32{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
