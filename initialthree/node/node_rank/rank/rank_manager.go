package rank

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/pkg/json"
	"initialthree/protocol/ss/rpc"
	"initialthree/zaplogger"
)

var mgr *Manager

func Mgr() *Manager {
	return mgr
}

type Manager struct {
	ranks   map[int32]*Rank
	running bool
	stopCh  chan bool
}

func (this *Manager) GetRank(rankID int32) *Rank {
	return this.ranks[rankID]
}

func (this *Manager) Running() bool {
	return this.running
}

func (this *Manager) Create(arg *rpc.RankCreateReq, logicAddr string) *rpc.RankCreateResp {
	base := &RankBase{
		ID:        arg.GetRankID(),
		BeginTime: arg.GetBeginTime(),
		EndTime:   arg.GetEndTime(),
	}

	r := newRank(base)
	r.Status = StatusNoBegin
	r.LastID = arg.GetLastRankID()
	if err := r.insertRankList(logicAddr, r.LastID); err != nil {
		// 已经存在的主键，获取创建的地址
		ret, err := dbClient.Get(tableRankList, fmt.Sprintf("%d", base.ID))
		if err != nil || len(ret) == 0 {
			zaplogger.GetSugar().Infof("rank create %d is exist,not in three", base.ID)
			return &rpc.RankCreateResp{Code: proto.Int32(0)}
		} else {
			logicAddr = ret[fieldRankListAddr].(string)
			zaplogger.GetSugar().Infof("rank create %d is exist, in three", base.ID)
			return &rpc.RankCreateResp{Code: proto.Int32(1), LogicAddr: proto.String(logicAddr)}
		}
	}

	if err := createTableRankRole(base.ID); err != nil {
		_ = deleteRankList(base.ID)
		zaplogger.GetSugar().Infof("rank create %d failed, %s", base.ID, err)
		return &rpc.RankCreateResp{Code: proto.Int32(0)}
	}

	r.start()
	this.ranks[base.ID] = r

	go clearHistoryRank(base.ID)

	return &rpc.RankCreateResp{Code: proto.Int32(2), LogicAddr: proto.String(logicAddr)}
}

func Remove(id int32) {
	cluster.PostTask(func() {
		delete(Mgr().ranks, id)
	})
}

// 仅初始化，当前正在运行的、上一期的排行榜。更早的历史记录不初始话
func (this *Manager) loadRankList() error {
	logic := cluster.SelfAddr().Logic.String()

	ranks := map[int32]*Rank{}
	idToLastId := map[int32]int32{}
	idToLogic := map[int32]string{}

	if err := dbClient.GetAll(tableRankList, func(row map[string]interface{}) error {
		infodata := row[fieldRankListInfo].([]byte)
		settledIdx := row[fieldRankListSettled].(int64)
		logicAddr := row[fieldRankListAddr].(string)
		lastID := int32(row[fieldRankListLastID].(int64))
		status := int32(row[fieldRankListStatus].(int64))

		if len(infodata) == 0 {
			return nil
		}
		var base *RankBase
		if err := json.Unmarshal(infodata, &base); err != nil {
			return err
		}

		r := newRank(base)
		r.Status = RankStatus(status)
		r.LastID = lastID
		r.settledIdx = int(settledIdx)

		idToLastId[base.ID] = lastID
		idToLogic[base.ID] = logicAddr
		ranks[base.ID] = r

		zaplogger.GetSugar().Info("getAll rank ", base, settledIdx, logicAddr, lastID)
		return nil
	}); err != nil {
		return err
	}

	// 每个link只初始化前两个，即 正在运作，上一期
	links := makeLink(idToLastId)
	for _, l := range links {
		//l.show()
		ids := l.getNodes(0, 2)
		for _, id := range ids {
			if id == 0 {
				continue
			}
			if logic_ := idToLogic[id]; logic != logic_ {
				zaplogger.GetSugar().Debugf("rank %d logicAddr %s is not self %s", id, logic_, logic)
				continue
			}

			r := ranks[id]
			if err := r.load(); err != nil {
				zaplogger.GetSugar().Errorf(err.Error())
			}
			this.ranks[r.ID] = r
			r.start()
		}
	}

	return nil
}

func Shutdown() {
	cluster.PostTask(func() {
		mgr.running = false
		mgr.stopCh = make(chan bool, len(mgr.ranks))
		for _, r := range mgr.ranks {
			insRank := r
			insRank.PostBack(func() { insRank.stop(mgr.stopCh) })
		}
	})

	cluster.WaitCondition(func() bool {
		return len(Mgr().stopCh) == len(Mgr().ranks)
	})

}

func InitManager() error {
	mgr = &Manager{
		ranks: map[int32]*Rank{},
	}

	// 启动时，重建排行榜
	if err := mgr.loadRankList(); err != nil {
		return err
	}

	zaplogger.GetSugar().Debug("rank serveCount", len(mgr.ranks))
	mgr.running = true

	return nil
}
