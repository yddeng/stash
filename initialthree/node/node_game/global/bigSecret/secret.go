package bigSecret

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/node/common/db"
	"initialthree/node/common/serverType"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/node_game/global"
	"initialthree/node/table/excel/DataTable/BigSecretCompetition"
	"initialthree/pkg/timer"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/rankCreate"
	"initialthree/zaplogger"
	"math"
	"time"
)

const (
	status_none = 1
	status_busy = 2

	dbTable = "global_data"
	dbKey   = "big_secret"
	dbField = "data"
)

type Data struct {
	RankLogic string `json:"rank_logic"`
	BeginTime int64  `json:"begin_time"`
	EndTime   int64  `json:"end_time"`
	RankID    int32  `json:"rank_id"`
	RankIdx   int32  `json:"rank_idx"`
}

type Secret struct {
	data           *Data
	status         int
	version        int64
	nextCreateTime int64
}

func (this *Secret) setStatus(status int) {
	this.status = status
}

// 个位：类型。高位：自增ID
func (this *Secret) genRankID() (lastID int32, rankID int32, newIdx int32) {
	newIdx = 1
	if this.data != nil {
		newIdx = this.data.RankIdx + 1
		lastID = this.data.RankID
	}
	rankID = newIdx*10 + global.BigSecret
	return
}

func (this *Secret) updateData() {
	lastID, rankID, newIdx := this.genRankID()
	// 读配置
	def := BigSecretCompetition.GetID(newIdx)
	if def == nil {
		zaplogger.GetSugar().Errorf("big secret competition %d is nil", newIdx)
		return
	}

	beginTime, err := time.ParseInLocation("2006-01-02 15:04:05", def.StartTime, time.Local)
	if err != nil {
		zaplogger.GetSugar().Errorf("big secret competition %d start time %s is failed, %s", newIdx, def.StartTime, err)
		return
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", def.EndTime, time.Local)
	if err != nil {
		zaplogger.GetSugar().Errorf("big secret competition %d end time %s is failed, %s", newIdx, def.EndTime, err)
		return
	}

	zaplogger.GetSugar().Infof("updateData big secret version %d beginTime %d endTime %d", this.version, beginTime, endTime)

	// rpc 创建rank
	this.createRank(rankID, lastID, beginTime.Unix(), endTime.Unix(), func(rankLogic string) {
		data := &Data{
			RankLogic: rankLogic,
			BeginTime: beginTime.Unix(),
			EndTime:   endTime.Unix(),
			RankID:    rankID,
			RankIdx:   newIdx,
		}
		this.saveData(data)
	})
}

func (this *Secret) createRank(rankID, lastID int32, beginT, endT int64, callback func(rankLogic string)) {
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

func (this *Secret) saveData(data *Data) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		zaplogger.GetSugar().Errorf(err.Error())
		return
	}

	this.setStatus(status_busy)
	set := db.GetFlyfishClient("game").Set(dbTable, dbKey, map[string]interface{}{
		dbField: dataBytes,
	}, this.version)
	set.AsyncExec(func(result *client.StatusResult) {
		this.setStatus(status_none)
		if result.ErrCode == nil {
			this.loadData()
		} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_version_mismatch {
			zaplogger.GetSugar().Infof("update big_secret failed code version %d miss match ,reload data", this.version)
			this.loadData()
		} else {
			zaplogger.GetSugar().Errorf("update big_secret failed: %s", errcode.GetErrorDesc(result.ErrCode))
		}
	})
}

func (this *Secret) loadData() {
	this.setStatus(status_busy)
	get := db.GetFlyfishClient("game").GetAll(dbTable, dbKey)
	get.AsyncExec(func(result *client.GetResult) {
		this.setStatus(status_none)
		if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
			var data *Data
			err := json.Unmarshal(result.Fields[dbField].GetBlob(), &data)
			if err != nil {
				zaplogger.GetSugar().Error(err.Error())
				return
			}
			this.version = *result.Version
			this.data = data
			this.nextCreateTime = 0
			zaplogger.GetSugar().Infof("load big_secret ok -> %v data %v", this, this.data)
		} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist {
			zaplogger.GetSugar().Errorf("load big_secret version %d ,code %s", this.version, errcode.GetErrorDesc(result.ErrCode))
		} else {
			zaplogger.GetSugar().Errorf("load big_secret version %d ,code %s", this.version, errcode.GetErrorDesc(result.ErrCode))
		}
	})
}

var secret *Secret

func (this *Secret) setNextCreateTime(now time.Time) {
	nextID := int32(1)
	if this.data != nil {
		nextID = this.data.RankIdx + 1
	}

	def := BigSecretCompetition.GetID(nextID)
	if def != nil {
		beginTime, err := time.ParseInLocation("2006-01-02 15:04:05", def.StartTime, time.Local)
		if err != nil {
			zaplogger.GetSugar().Errorf("big secret competition %d start time %s is failed, %s", nextID, def.StartTime, err)
			this.nextCreateTime = math.MaxInt64
			return
		}
		this.nextCreateTime = beginTime.Unix()
	} else {
		// 已确认没有最新配置，不用继续读配置了
		this.nextCreateTime = math.MaxInt64
	}
}

func (this *Secret) tick(now time.Time) {
	if this.nextCreateTime == 0 {
		this.setNextCreateTime(now)
	}

	if this.status == status_none && this.nextCreateTime != 0 && now.Unix() > this.nextCreateTime {
		this.updateData()
	}
}

func (this *Secret) Tick(t *timer.Timer, _ interface{}) {
	now := timeDisposal.Now()
	this.tick(now)
}

func GetData() *Data {
	return secret.data
}

func ParseRankID(rankID int32) (config, tt int32) {
	return rankID / 10, rankID % 10
}

func Launch() {
	secret = &Secret{data: nil}
	secret.loadData()

	cluster.RegisterTimer(time.Second, secret.Tick, nil)
}
