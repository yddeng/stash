package rankData

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common/droppool"
	"initialthree/node/node_game/global"
	"initialthree/node/table/excel/DataTable/ScarsIngrainArea"
	"initialthree/pkg/json"

	"initialthree/node/table/excel/DataTable/Mail"

	"initialthree/pkg/timer"
	"initialthree/rpc/rankGetRank"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/rankSetScore"
	"time"
)

type RankInfo struct {
	RankID    int32  `json:"Id"`
	RankLogic string `json:"rL"`
	GetAward  bool   `json:"ga"`
}

type RankData struct {
	userI module.UserI
	data  map[int32]*RankInfo
	*module.ModuleSaveBase
}

func (this *RankData) ModuleType() module.ModuleType {
	return module.RankData
}

func (this *RankData) GetRankInfo(rankID int32) *RankInfo {
	return this.data[rankID]
}

func (this *RankData) SetRank(RankID int32, Score int32, rankLogic string, RoleInfo *message.RankRoleInfo, callback func(rank, total int32, err error)) {
	rankAddr, err := addr.MakeLogicAddr(rankLogic)
	if err != nil {
		callback(0, 0, err)
		return
	}

	msg := &rpc.RankSetScoreReq{
		RankID:   proto.Int32(RankID),
		RoleID:   proto.Uint64(this.userI.GetID()),
		Score:    proto.Int32(Score),
		RoleInfo: RoleInfo,
	}

	this.updateScore(msg, rankAddr, time.Now().Add(time.Second*5), func(rank, total int32, err error) {
		if err != nil {
			callback(0, 0, err)
			return
		}

		info, ok := this.data[RankID]
		if !ok {
			info = &RankInfo{}
			this.data[RankID] = info
		}
		info.RankID = RankID
		info.RankLogic = rankLogic
		this.SetDirty(this.ModuleType().String())

		callback(rank, total, nil)
	})

}

func (this *RankData) updateScore(msg *rpc.RankSetScoreReq, rankAddr addr.LogicAddr, deadline time.Time, callback func(rank, total int32, err error)) {
	rankSetScore.AsynCall(rankAddr, msg, deadline.Sub(time.Now()), func(resp *rpc.RankSetScoreResp, e error) {
		// 0 继续尝试，1 成功 2 失败
		if e != nil || resp.GetCode() == 0 {
			cluster.RegisterTimerOnce(time.Second, func(timer *timer.Timer, i interface{}) {
				if time.Now().Before(deadline.Add(-time.Millisecond * 100)) {
					this.updateScore(msg, rankAddr, deadline, callback)
				} else {
					callback(0, 0, fmt.Errorf("rankData updateScore %s errror %v ", rankAddr.String(), e))
				}
			}, nil)

			return
		}

		switch resp.GetCode() {
		case 1:
			callback(resp.GetRank(), resp.GetTotal(), nil)
		case 2:
			callback(resp.GetRank(), resp.GetTotal(), fmt.Errorf("rankData updateScore %s randk %d not running", rankAddr.String(), msg.GetRankID()))
		default:
			callback(0, 0, fmt.Errorf("rankData updateScore %s code %d invaild", rankAddr.String(), resp.GetCode()))
		}
	})
}

func (this *RankData) DelRank(rankID int32) {
	delete(this.data, rankID)
	this.SetDirty(this.ModuleType().String())
}

func (this *RankData) SetApplyAward(rankID int32) {
	if info, ok := this.data[rankID]; ok {
		info.GetAward = true
		this.SetDirty(this.ModuleType().String())
	}
}

func (this *RankData) Init(fields map[string]*flyfish.Field) error {
	field, ok := fields[this.ModuleType().String()]

	if ok && len(field.GetBlob()) != 0 {
		var data map[int32]*RankInfo
		err := json.Unmarshal(field.GetBlob(), &data)
		if err != nil {
			return fmt.Errorf("unmarshal: %s", err)
		} else {
			this.data = data
		}
	}

	return nil
}

func (this *RankData) ReadOut() *module.ReadOutCommand {
	return &module.ReadOutCommand{
		Table:  "user_module_data",
		Key:    this.userI.GetIDStr(),
		Fields: []string{this.ModuleType().String()},
		Module: this,
	}
}

func (this *RankData) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
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

// 初始化后, 批量验证排行榜状态。
func (this *RankData) checkRank() {
	for _, v := range this.data {
		logicAddr, err := addr.MakeLogicAddr(v.RankLogic)
		if err != nil {
			zaplogger.GetSugar().Errorf("rankData check logic %s make err %s", v.RankLogic, err.Error())
			this.DelRank(v.RankID)
			continue
		}

		rankID := v.RankID

		CallGetRank(logicAddr, rankID, []uint64{this.userI.GetID()}, func(results []*rpc.RankGetRankResult, err error) {
			if err != nil {
				zaplogger.GetSugar().Errorf("rankData callGetRank err %s", err.Error())
				return
			}

			result := results[0]
			// 0 没有玩家数据， 1 排行榜正在进行， 2 排行榜完结，且有排名
			switch result.GetCode() {
			case 0:
				zaplogger.GetSugar().Debugf("%s rankData rank %d is not exit", this.userI.GetUserID(), rankID)
				this.DelRank(rankID)
			case 2:
				rank, total := result.GetRank(), result.GetTotal()
				if rank > 0 && total > 0 {
					this.ApplyRankAward(rankID, result.GetScore(), rank, total)
				}
				this.SetApplyAward(rankID)
			default:
			}

		})
	}
}

func callGetRank(rankAddr addr.LogicAddr, msg *rpc.RankGetRankReq, callback func(result []*rpc.RankGetRankResult, err error)) {
	rankGetRank.AsynCall(rankAddr, msg, time.Second*6, func(resp *rpc.RankGetRankResp, e error) {
		if e != nil || !resp.GetOk() {
			callback(nil, fmt.Errorf("call get rank err %v ", e))
			return
		}

		if resp.GetRedirectRankAddr() != "" {
			logicAddr, err := addr.MakeLogicAddr(resp.GetRedirectRankAddr())
			if err != nil {
				callback(nil, fmt.Errorf("call get rank make logicAddr %s err %v ", resp.GetRedirectRankAddr(), e))
				return
			}
			callGetRank(logicAddr, msg, callback)
		} else {
			callback(resp.GetResults(), nil)
		}

	})
}

func CallGetRank(rankAddr addr.LogicAddr, rankID int32, roleIDs []uint64, callback func(result []*rpc.RankGetRankResult, err error)) {
	req := &rpc.RankGetRankReq{
		RoleID: roleIDs,
		RankID: proto.Int32(rankID),
	}

	callGetRank(rankAddr, req, callback)
}

func (this *RankData) ApplyRankAward(rankID, score, rank, total int32) {
	rankInfo, ok := this.data[rankID]
	if ok && !rankInfo.GetAward {
		awardPoolID := global.RankAwardPoolID(rankID, rank, total)
		if awardPoolID != 0 {
			award := droppool.DropAward(awardPoolID)
			if !award.IsZero() {
				tt := global.ParseRankID(rankID)
				switch tt {
				case global.ScarsIngrain:
					area := (rankID / 10) % 10
					sia := ScarsIngrainArea.GetID(area)
					m := Mail.ScarsIngrainMail(time.Now(), award.ToMessageAward(),
						sia.AreaName,
						fmt.Sprintf("%d", score),
						fmt.Sprintf("%d%%", int32(global.RankPercent(rank, total)*100)),
					)
					this.userI.SendMail([]*message.Mail{m})
				case global.BigSecret:
				default:

				}
			}
		}
	}
}

func (this *RankData) AfterInitAll() error {
	this.checkRank()
	return nil
}

func (this *RankData) Tick(now time.Time)               {}
func (this *RankData) FlushDirtyToClient()              {}
func (this *RankData) FlushAllToClient(seqNo ...uint32) {}

func init() {
	module.RegisterModule(module.RankData, func(userI module.UserI) module.ModuleI {
		m := &RankData{
			userI: userI,
			data:  map[int32]*RankInfo{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
