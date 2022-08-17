package ScarsIngrain

import (
	"github.com/golang/protobuf/proto"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/global"
	scarsIngrain2 "initialthree/node/node_game/global/scarsIngrain"
	"initialthree/node/node_game/module/base"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/rankData"
	"initialthree/node/node_game/module/scarsIngrain"
	"initialthree/node/node_game/user"
	DungeonTable "initialthree/node/table/excel/DataTable/Dungeon"
	"initialthree/node/table/excel/DataTable/ScarsIngrainBossInstance"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionScarsIngrainSaveScore struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
	resp    *cs_message.ScarsIngrainSaveScoreToC
}

func (this *transactionScarsIngrainSaveScore) GetModuleName() string {
	return "ScarsIngrain"
}

func (this *transactionScarsIngrainSaveScore) Begin() {
	this.resp = &cs_message.ScarsIngrainSaveScoreToC{}
	msg := this.req.GetData().(*cs_message.ScarsIngrainSaveScoreToS)
	zaplogger.GetSugar().Infof("%s ScarsIngrainSaveScoreToS %v ", this.user.GetUserLogName(), msg)

	siModule := this.user.GetSubModule(module.ScarsIngrain).(*scarsIngrain.ScarsIngrain)
	siData := siModule.GetData()

	data, ok := siModule.GetFightData(msg.GetScoreID())
	if !ok && msg.GetIsSave() {
		zaplogger.GetSugar().Infof("%s ScarsIngrainSaveScoreToS failed, scoreID %d is nil", this.user.GetUserLogName(), msg.GetScoreID())
		this.errcode = cs_message.ErrCode_ScarsIngrain_SaveScoreIDErr
		this.EndTrans(this.resp, this.errcode)
		return
	}

	if !msg.GetIsSave() {
		siModule.ClearFightData()
		this.errcode = cs_message.ErrCode_OK
		this.EndTrans(this.resp, this.errcode)
		return
	}

	// 结算
	dungeonCfg := DungeonTable.GetID(data.Tos.GetDungeonID())
	bIns := ScarsIngrainBossInstance.GetID(dungeonCfg.SystemConfig)
	team := data.Tos.GetCharacterTeam()
	bossId, difficult := bIns.BossDifficult()
	totalScore, newTotalScore := siModule.FightEnd(bossId, difficult, data.Score, data.BossDie, team.GetCharacterList())

	// 移除零时数据
	siModule.ClearFightData()

	// 事件触发
	this.user.EmitEvent(event.EventInstanceSucceed, data.Tos.GetDungeonID(), int32(0), map[int32]int32{})

	if newTotalScore {
		// 该分数重新计算后需要同步到排行榜
		baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
		gSiData := scarsIngrain2.GetIDData(siData.SIID)
		rankRoleInfo := &cs_message.RankRoleInfo{
			ID:            proto.Uint64(this.user.GetID()),
			Name:          proto.String(baseModule.GetName()),
			Level:         proto.Int32(this.user.GetLevel()),
			Score:         proto.Int32(totalScore),
			CharacterList: team.GetCharacterList(),
		}

		rankModule := this.user.GetSubModule(module.RankData).(*rankData.RankData)
		rankModule.SetRank(gSiData.RankID, totalScore, gSiData.RankLogic, rankRoleInfo, func(rank, total int32, err error) {
			if err != nil {
				zaplogger.GetSugar().Debugf("%s ScarsIngrainSaveScoreToS failed, %v ", this.user.GetUserLogName(), err)
				this.errcode = cs_message.ErrCode_ScarsIngrain_RankIsEnd
				this.EndTrans(this.resp, this.errcode)
				return
			}
			zaplogger.GetSugar().Debugf("%s ScarsIngrainSaveScoreToS, rank (idx %d, total %d)", this.user.GetUserLogName(), rank, total)
			this.resp.Rank = proto.Int32(rank)
			this.resp.Total = proto.Int32(total)
			if rank > 0 && total > 0 {
				this.resp.Percent = proto.Int32(int32(global.RankPercent(rank, total) * 100))
			}
			this.errcode = cs_message.ErrCode_OK
			this.EndTrans(this.resp, this.errcode)
		})
	} else {
		this.errcode = cs_message.ErrCode_OK
		this.EndTrans(this.resp, this.errcode)
	}
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_ScarsIngrainSaveScore, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionScarsIngrainSaveScore{
			user: user,
			req:  msg,
		}
	})
}
