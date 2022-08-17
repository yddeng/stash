package scarsIngrain

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/global/scarsIngrain"
	"initialthree/pkg/json"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/table/excel/DataTable/ScarsIngrainArea"
	"initialthree/node/table/excel/DataTable/ScarsIngrainBoss"
	"initialthree/protocol/cs/message"
	"time"
)

const (
	dataField = "sidata"
	timeField = "timedata"
)

var tableFields = []string{dataField, timeField}

type SIData struct {
	SIID       int32                       `json:"siid"`    // 等级区间
	RankID     int32                       `json:"rank_id"` // 排行榜ID，用于判断是否过期
	TotalScore int32                       `json:"total_score"`
	BossGroup  []*message.ScarsIngrainBoss `json:"boss_group"`
	ScoreAward map[int32]struct{}          `json:"score_award"` // 记录已经领取的奖
}

type ScarsIngrain struct {
	userI    module.UserI
	data     SIData
	timedata map[string]int64 // 时间戳，定时器相关。 每日，每周

	// 零时数据
	scoreId   int32
	fightData map[int32]*FightData
	dirty     bool
	*module.ModuleSaveBase
}

func (this *ScarsIngrain) GetData() *SIData {
	return &this.data
}

func (this *ScarsIngrain) BossDifficultUnlock(bossID, difficult int32) bool {
	idx, _, score := this.BossDifficultScore(bossID, difficult)
	if idx == -1 || !this.data.BossGroup[idx].GetOpened() {
		return false
	}

	return score != -1
}

// 返回 当前boss所属的group，当前boss难度的分数
func (this *ScarsIngrain) BossDifficultScore(bossID, difficult int32) (int, *message.ScarsIngrainBoss, int32) {
	for idx, boss := range this.data.BossGroup {
		if boss.GetBossID() == bossID {
			for i, score := range boss.DifficultScore {
				if int32(i+1) == difficult {
					return idx, boss, score
				}
			}
			return idx, boss, -1
		}
	}
	return -1, nil, -1
}

// 返回除当前boss以外的驻守队伍
func (this *ScarsIngrain) GetDefendCharacter(bossID int32) (list map[int32]struct{}) {
	list = make(map[int32]struct{})
	for _, boss := range this.data.BossGroup {
		if boss.GetBossID() != bossID {
			for _, id := range boss.GetDefendCharacterList() {
				if id != 0 {
					list[id] = struct{}{}
				}
			}
		}
	}
	return
}

func (this *ScarsIngrain) SetScoreAwardGet(score int32) {
	this.data.ScoreAward[score] = struct{}{}
	this.SetDirty(dataField)
	this.dirty = true
}

func (this *ScarsIngrain) BossTeamPrefabSet(bossID int32, characterList []int32) bool {
	for _, boss := range this.data.BossGroup {
		if boss.GetBossID() == bossID {
			boss.PrefabCharacterList = characterList
			this.SetDirty(dataField)
			this.dirty = true
			return true
		}
	}
	return false
}

// 返回累积分数，本次结果是否需要同步
// 需要同步条件： 1. boss本难度本周第一次挑战，累积分数会加上本周最高分数 2.boss本难度刷新本周最高记录，累积分数重新计算
func (this *ScarsIngrain) FightEnd(bossId, difficult, score int32, bossDie bool, characterList []int32) (totalScore int32, newTotalScore bool) {

	_, bossGroup, curScore := this.BossDifficultScore(bossId, difficult)
	// 队伍预设
	bossGroup.PrefabCharacterList = characterList

	// 解锁下一等级
	if bossDie {
		bossDef := ScarsIngrainBoss.GetID(bossId)
		if difficult == int32(len(bossGroup.DifficultScore)) && difficult < int32(len(bossDef.BossDifficulty)) {
			// 为前最高难度且还有更高难度
			bossGroup.DifficultScore = append(bossGroup.DifficultScore, 0)
		}

		// 队伍驻守
		if difficult == int32(len(bossDef.BossDifficulty)) {
			bossGroup.DefendCharacterList = characterList
		}
	}

	// 该boss难度, 重新计算累积分数
	newTotalScore = score != curScore
	if newTotalScore {
		bossGroup.DifficultScore[difficult-1] = score
		// 累积分数=boss1难度1最高保存分数+boss1难度2最高保存分数+....  boss2难度5最高保存分数
		for _, boss := range this.data.BossGroup {
			for _, c := range boss.DifficultScore {
				totalScore += c
			}
		}
		// 累积分数
		this.data.TotalScore = totalScore
	}

	this.SetDirty(dataField)
	this.dirty = true
	return
}

func (this *ScarsIngrain) makeSiData(siID int32, rankID int32, bossChallenges []*scarsIngrain.BossChallenge) SIData {
	data := SIData{
		SIID:       siID,
		RankID:     rankID,
		TotalScore: 0,
		BossGroup:  make([]*message.ScarsIngrainBoss, 0, len(bossChallenges)),
		ScoreAward: map[int32]struct{}{},
	}

	for _, v := range bossChallenges {
		data.BossGroup = append(data.BossGroup, this.makeBossGroupIns(v.BossID, v.ChallengeID))
	}
	return data
}

func (this *ScarsIngrain) makeBossGroupIns(bossID, challenge int32) *message.ScarsIngrainBoss {
	bossDef := ScarsIngrainBoss.GetID(bossID)
	bossIns := &message.ScarsIngrainBoss{
		BossID:         proto.Int32(bossID),
		ChallengeID:    proto.Int32(challenge),
		DifficultScore: make([]int32, 0, len(bossDef.BossDifficulty)),
	}
	bossIns.DifficultScore = append(bossIns.DifficultScore, 0)
	return bossIns
}

func (this *ScarsIngrain) check() {
	id := this.data.SIID
	if id == 0 {
		area := ScarsIngrainArea.GetPrisonByLevel(this.userI.GetLevel())
		// 获取本期信息。不在进行时间内，可能正在创建、竞争。
		ok := scarsIngrain.GetIDClass(area.ID).IsRunning(0)
		if ok {
			siData := scarsIngrain.GetIDData(area.ID)
			this.data = this.makeSiData(area.ID, siData.RankID, siData.BossChallenges)
			this.SetDirty(dataField)
			this.dirty = true
		}
	} else {
		siData := scarsIngrain.GetIDData(id)
		if siData.RankID == this.data.RankID {
			// boss开启
			for idx, v := range siData.BossChallenges {
				if this.data.BossGroup[idx].GetOpened() != v.Opened {
					this.data.BossGroup[idx].Opened = proto.Bool(v.Opened)
					this.dirty = true
				}
			}

		} else {
			// 副本已经刷新，重新计算分区
			area := ScarsIngrainArea.GetPrisonByLevel(this.userI.GetLevel())
			// 获取本期信息。不在进行时间内，可能正在创建、竞争。
			ok := scarsIngrain.GetIDClass(area.ID).IsRunning(0)
			if ok {
				siData = scarsIngrain.GetIDData(area.ID)
				this.data = this.makeSiData(area.ID, siData.RankID, siData.BossChallenges)
				this.SetDirty(dataField)
				this.dirty = true
			}
		}
	}
	//util.Logger().Debugf("%s check %v data %v %v %v %v", this.userI.GetUserLogName(), this.checkData, this.data, this.dailyTimes, this.roleTimes, this.timedata)
}

func (this *ScarsIngrain) Tick(now time.Time) {
	this.check()
	this.clockTimer()
}

func (this *ScarsIngrain) ModuleType() module.ModuleType {
	return module.ScarsIngrain
}

func (this *ScarsIngrain) Init(fields map[string]*flyfish.Field) error {

	for _, name := range tableFields {
		field, ok := fields[name]
		if ok && len(field.GetBlob()) != 0 {
			var err error
			switch name {
			case dataField:
				err = json.Unmarshal(field.GetBlob(), &this.data)
			case timeField:
				err = json.Unmarshal(field.GetBlob(), &this.timedata)
			}
			if err != nil {
				zaplogger.GetSugar().Errorf("%s initScarsIngrain name %s data %s ", this.userI.GetUserID(), name, string(field.GetBlob()))
				return fmt.Errorf("unmarshal: %s", err)
			}
		}
	}

	return nil
}

func (this *ScarsIngrain) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Fields: tableFields,
		Module: this,
	}

	return cmd
}

func (this *ScarsIngrain) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  this.ModuleType().String(),
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

func (this *ScarsIngrain) pack() *message.ScarsIngrainSyncToC {
	if this.data.SIID == 0 {
		return &message.ScarsIngrainSyncToC{
			CurrentID: proto.Int32(this.data.SIID),
		}
	}
	siMap := ScarsIngrainArea.GetIDMap()
	resp := &message.ScarsIngrainSyncToC{
		ScarsIngrainGroup: make([]*message.ScarsIngrain, 0, len(siMap)),
		CurrentID:         proto.Int32(this.data.SIID),
		TotalScore:        proto.Int32(this.data.TotalScore),
		CurrentRankId:     proto.Int32(this.data.RankID),
		BossGroup:         this.data.BossGroup,
		ScoreRewardGet:    make([]int32, 0, len(this.data.ScoreAward)),
	}

	// 分数领奖情况
	for score := range this.data.ScoreAward {
		resp.ScoreRewardGet = append(resp.ScoreRewardGet, score)
	}

	// 全区段基础信息
	for _, v := range siMap {
		gSiData := scarsIngrain.GetIDData(v.ID)
		bossIds := make([]int32, 0, len(gSiData.BossChallenges))
		for _, v := range gSiData.BossChallenges {
			bossIds = append(bossIds, v.BossID)
		}
		si := &message.ScarsIngrain{
			ScarsIngrainID: proto.Int32(v.ID),
			RankID:         proto.Int32(gSiData.RankID),
			BossIDs:        bossIds,
		}
		resp.ScarsIngrainGroup = append(resp.ScarsIngrainGroup, si)

	}

	return resp
}

func (this *ScarsIngrain) FlushDirtyToClient() {
	if this.dirty {
		this.userI.Post(this.pack())
		this.dirty = false
	}
}

func (this *ScarsIngrain) FlushAllToClient(seqNo ...uint32) {
	this.userI.Post(this.pack())
}

func init() {
	module.RegisterModule(module.ScarsIngrain, func(userI module.UserI) module.ModuleI {
		m := &ScarsIngrain{
			userI:     userI,
			timedata:  map[string]int64{},
			fightData: map[int32]*FightData{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
