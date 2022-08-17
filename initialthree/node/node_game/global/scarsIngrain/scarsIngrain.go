package scarsIngrain

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/node/common/db"
	"initialthree/node/common/serverType"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/node_game/global"
	constSI "initialthree/node/table/excel/ConstTable/ScarsIngrain"
	dataSIArea "initialthree/node/table/excel/DataTable/ScarsIngrainArea"
	"initialthree/node/table/excel/DataTable/ScarsIngrainBoss"
	"initialthree/pkg/json"
	"initialthree/pkg/timer"
	"initialthree/protocol/ss/rpc"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/rpc/rankCreate"
	"initialthree/zaplogger"
	"math/rand"
	"time"
)

/*
 战痕印刻
*/

const (
	status_none = 1
	status_busy = 2

	dbTable     = "global_data"
	dbKeyPrefix = "scars_ingrain"
	dbField     = "data"
)

var siMgr *ScarsIngrainMgr

type ScarsIngrainMgr struct {
	class map[int32]*ScarsIngrain
}

type ScarsIngrain struct {
	id             int32
	keyStr         string // id string
	status         int
	nextCreateTime int64 // 创建时间
	version        int64 // 数据库版本号
	pendFunc       []func()

	data *ScarsIngrainData
}

type ScarsIngrainData struct {
	RankLogic      string           `json:"rank_logic"`
	BeginTime      int64            `json:"begin_time"`
	EndTime        int64            `json:"end_time"`
	RankID         int32            `json:"rank_id"`
	RankIdx        int32            `json:"rank_idx"`
	BossChallenges []*BossChallenge `json:"boss_challenges"`
}

type BossChallenge struct {
	BossID      int32 `json:"boss_id"`
	ChallengeID int32 `json:"challenge_id"`
	Opened      bool  `json:"opened"`
}

func (this *ScarsIngrain) PendFunc(fn func()) {
	this.pendFunc = append(this.pendFunc, fn)
}

func (this *ScarsIngrain) doPendFunc() {
	if len(this.pendFunc) > 0 {
		fn := this.pendFunc[0]
		this.pendFunc = this.pendFunc[1:]
		fn()
	}
}

func (this *ScarsIngrain) setStatus(status int) {
	this.status = status
}

// 个位：类型。十位：等级区间。高位：自增ID
func (this *ScarsIngrain) genRankID() (lastID int32, rankID int32, newIdx int32) {
	newIdx = 1
	if this.data != nil {
		newIdx = this.data.RankIdx + 1
		lastID = this.data.RankID
	}
	rankID = newIdx*100 + this.id*10 + global.ScarsIngrain
	return
}

/*
 rank 上请求创建rank 实例。各个 game 竞争，由唯一ID判断是否已经创建。不为失败时，尝试保存数据库（已经存在、创建成功。理解为 rank创建成功由我完成的）
 数据保存也为竞争模式，由数据库 version 判断是否已经被更新。若果 当前版本号小于数据库版本号，则已经被更新，重新加载数据。
 出现情况： 如果rank上创建成功，但是数据落地失败。每一个game到达时间，都会重复以上过程，直到整个过程完成。
*/

func (this *ScarsIngrain) updateData() {

	lastID, rankID, newIdx := this.genRankID()
	beginTime := constSI.GetBeginTime().Unix()
	endTime := constSI.GetEndTime().Unix()
	// 第一次创建, 或者 中间几期没有开服，重启后正好在一期的中间位置
	if beginTime > endTime {
		beginTime = constSI.GetBeginTime().AddDate(0, 0, -7).Unix()
	}
	zaplogger.GetSugar().Infof("scarsIngrain updateData %d version %d beginTime %d endTime %d", this.id, this.version, beginTime, endTime)

	// 生成随机boss
	exits := map[int32]struct{}{}
	if this.data != nil {
		for _, v := range this.data.BossChallenges {
			exits[v.BossID] = struct{}{}
		}
	}
	cof := dataSIArea.GetID(this.id)
	bossIDs := cof.RandomBoss(exits, int(cof.BossCount))
	bossChallenges := make([]*BossChallenge, 0, len(bossIDs))
	for _, bossId := range bossIDs {
		bossDef := ScarsIngrainBoss.GetID(bossId)
		if bossDef == nil {
			panic(fmt.Sprintf("table ScarsIngrainBoss Id %d is not exist", bossId))
		}

		idx := rand.Int() % len(bossDef.ChallengeConfigArray)
		bossChallenges = append(bossChallenges, &BossChallenge{
			BossID:      bossId,
			ChallengeID: bossDef.ChallengeConfigArray[idx].ID,
		})
	}

	// rpc 创建rank
	this.createRank(rankID, lastID, beginTime, endTime, func(rankLogic string) {
		data := &ScarsIngrainData{
			RankLogic:      rankLogic,
			BeginTime:      beginTime,
			EndTime:        endTime,
			RankID:         rankID,
			RankIdx:        newIdx,
			BossChallenges: bossChallenges,
		}
		this.saveData(data)
	})
}

func (this *ScarsIngrain) createRank(rankID, lastID int32, beginT, endT int64, callback func(rankLogic string)) {
	this.setStatus(status_busy)
	req := &rpc.RankCreateReq{
		RankID:     proto.Int32(rankID),
		BeginTime:  proto.Int64(beginT),
		EndTime:    proto.Int64(endT),
		LastRankID: proto.Int32(lastID),
	}

	logicAddr, err := cluster.Random(serverType.Rank)
	if err != nil {
		zaplogger.GetSugar().Errorf("%s service not found,%s", serverType.Type2Name(serverType.Rank), err.Error())
		this.setStatus(status_none)
	} else {
		rankCreate.AsynCall(logicAddr, req, time.Second*8, func(resp *rpc.RankCreateResp, e error) {
			this.setStatus(status_none)
			if e != nil {
				zaplogger.GetSugar().Error(e.Error())
				return
			}

			// code = 1; // 0 失败，1 已经存在，2 成功, 3 服务上限
			code := resp.GetCode()
			if code == 1 {
				this.loadData()
			} else if code == 2 {
				callback(resp.GetLogicAddr())
			}
		})
	}
}

func (this *ScarsIngrain) saveData(data *ScarsIngrainData) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		zaplogger.GetSugar().Errorf(err.Error())
		return
	}

	this.setStatus(status_busy)
	set := db.GetFlyfishClient("game").Set(dbTable, this.keyStr, map[string]interface{}{
		dbField: dataBytes,
	}, this.version)
	set.AsyncExec(func(result *client.StatusResult) {
		this.setStatus(status_none)
		if result.ErrCode == nil {
			//this.data = data
			//this.version = result.Version
			//this.nextCreateTime = time.Unix(data.BeginTime, 0).AddDate(0, 0, 7).Unix() - 10 // 下一次开始时间
			this.loadData()
			this.notifyUpdate()
		} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_version_mismatch {
			zaplogger.GetSugar().Infof("update prison class %d failed code version %d miss match ,reload data", this.id, this.version)
			this.loadData()
		} else {
			zaplogger.GetSugar().Errorf("update prison class %d failed: %s", this.id, errcode.GetErrorDesc(result.ErrCode))
		}
	})
}

// 通知game(不包括自己)， 版本已更新，需要重新拉取
func (this *ScarsIngrain) notifyUpdate() {
	msg := &ssmessage.ScarsIngrainUpdate{
		Version: proto.Int64(this.version),
		Id:      proto.Int32(this.id),
	}
	cluster.Brocast(serverType.Game, msg, true)
}

func (this *ScarsIngrain) loadData() {
	this.setStatus(status_busy)
	get := db.GetFlyfishClient("game").GetAll(dbTable, this.keyStr)
	get.AsyncExec(func(result *client.GetResult) {
		this.setStatus(status_none)
		if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
			var data *ScarsIngrainData
			err := json.Unmarshal(result.Fields[dbField].GetBlob(), &data)
			if err != nil {
				zaplogger.GetSugar().Error(err.Error())
				return
			}
			this.version = *result.Version
			this.data = data
			this.nextCreateTime = time.Unix(data.BeginTime, 0).AddDate(0, 0, 7).Unix() - 10 // 下一次开始时间
			zaplogger.GetSugar().Infof("load ScarsIngrain %d ok -> %v data %v", this.id, this, this.data)
		} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist {
			zaplogger.GetSugar().Errorf("load ScarsIngrain %d version %d ,code %s", this.id, this.version, errcode.GetErrorDesc(result.ErrCode))
		} else {
			zaplogger.GetSugar().Errorf("load ScarsIngrain %d version %d ,code %s", this.id, this.version, errcode.GetErrorDesc(result.ErrCode))
		}

	})
}

func (this *ScarsIngrain) updateOpened() {
	if this.data == nil {
		return
	}
	for idx, v := range this.data.BossChallenges {
		if !v.Opened {
			v.Opened = constSI.BossOpen(idx)
		}
	}
}

func (this *ScarsIngrain) tick(now time.Time) {
	if this.status == status_none {
		this.doPendFunc()
	}

	if this.status == status_none && now.Unix() > this.nextCreateTime {
		this.updateData()
	}

	this.updateOpened()
}

func (this *ScarsIngrainMgr) Tick(t *timer.Timer, _ interface{}) {
	now := timeDisposal.Now()
	for _, c := range this.class {
		c.tick(now)
	}
}

func Launch() {
	siMgr = &ScarsIngrainMgr{
		class: map[int32]*ScarsIngrain{},
	}

	siMap := dataSIArea.GetIDMap()
	for id := range siMap {
		class := &ScarsIngrain{
			id:     id,
			keyStr: fmt.Sprintf("%s_%d", dbKeyPrefix, id),
		}

		class.loadData()

		siMgr.class[id] = class
	}

	cluster.RegisterTimer(time.Second, siMgr.Tick, nil)
}

// 返回当前等级区间的 数据
func GetIDData(id int32) *ScarsIngrainData {
	v, ok := siMgr.class[id]
	if ok && v.data != nil {
		return v.data
	}
	return nil
}

func GetIDClass(id int32) *ScarsIngrain {
	v, ok := siMgr.class[id]
	if ok {
		return v
	}
	return nil
}

// 是否正在进行. 可能在竞争期间
// export 一个副本的时间， 例： 一场战斗进行 120s，但是本期只有 60s 结束。则判断超期，不允许进入
func (this *ScarsIngrain) IsRunning(export int64) bool {
	if this.version == 0 {
		return false
	}
	now := timeDisposal.Now().Unix()
	if now >= this.data.BeginTime && now+export < this.data.EndTime {
		return true
	}
	return false
}

/*
 如果rank先时间到期，这里还没有到期。玩家看到可以上分，但实际已经不计入排名
*/
