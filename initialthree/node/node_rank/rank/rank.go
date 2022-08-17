package rank

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/sortedset"
	"initialthree/cluster"
	"initialthree/cluster/priority"
	"initialthree/node/common/serverType"
	"initialthree/pkg/event"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/zaplogger"
	"strconv"
	"sync/atomic"
	"time"
)

type RankBase struct {
	ID        int32
	BeginTime int64
	EndTime   int64
}

type Rank struct {
	RankBase
	Status     RankStatus
	LastID     int32 // 上期ID
	Total      int32 // 总人数
	zset       *sortedset.SortedSet
	eventQueue *event.EventQueue
	// top 榜玩家详细数据
	topCount   int
	topRoles   []*message.RankRoleInfo
	topVersion int32
	// 结算
	settledFired bool // 结算函数防止重复调用
	settledIdx   int  // 结算索引

	// 脏数据
	dirtyStatus  bool
	roleDirty    map[uint64]*message.RankRoleInfo
	nextSaveTime time.Time

	stopSign  chan bool //停止信号
	asynCount int32     // 正在异步调用的数量。停止时 需判断为 0, atomic
}

func newRank(base *RankBase) *Rank {
	return &Rank{
		RankBase:     *base,
		zset:         sortedset.New(),
		eventQueue:   event.NewEventQueue(),
		topRoles:     make([]*message.RankRoleInfo, topN),
		topVersion:   1,
		roleDirty:    map[uint64]*message.RankRoleInfo{},
		nextSaveTime: Now().Add(saveDuration),
		stopSign:     make(chan bool),
	}
}

func (r *Rank) Reset(startT, endT int64) *Rank {
	r.BeginTime = startT
	r.EndTime = endT
	r.Status = StatusNoBegin
	r.zset.Init()
	return r
}

func (r *Rank) start() {
	go r.eventQueue.Run()
	r.PostBack(r.tick)
}

func (r *Rank) stop(stopCh chan<- bool) {
	zaplogger.GetSugar().Infof("rank %d do stop\n", r.ID)
	close(r.stopSign)

	// 立即触发保存
	r.saveRoleDirty()
	_ = r.saveStatus()

	fn := func() bool {
		if len(r.roleDirty) == 0 && atomic.LoadInt32(&r.asynCount) == 0 {
			return true
		}
		zaplogger.GetSugar().Debugf("rank %d wait stop dirtyLen(%d) asynCount(%d) \n", r.ID, len(r.roleDirty), atomic.LoadInt32(&r.asynCount))
		return false
	}

	go func() {
		stoped := int32(0)
		for atomic.LoadInt32(&stoped) == 0 {
			time.Sleep(time.Millisecond * 100)
			r.PostBack(func() {
				if fn() {
					atomic.StoreInt32(&stoped, 1)
				}
			})
		}

		zaplogger.GetSugar().Infof("rank %d stop ok\n", r.ID)
		stopCh <- true
	}()
}

func (r *Rank) PostBack(fn interface{}, args ...interface{}) {
	r.eventQueue.PostNoWait(priority.MID, fn, args...)
}

func (r *Rank) changeStatus(now int64) {
	if r.Status == StatusNoBegin && now >= r.BeginTime {
		r.Status = StatusBegin
		r.dirtyStatus = true
	}
	if r.Status == StatusBegin && now >= r.EndTime {
		r.Status = StatusEnd
		r.dirtyStatus = true
		// 排名期结束，保存所有玩家数据
		r.saveRoleDirty()
	}
	if r.Status == StatusEnd {
		r.Status = StatusSettling
		r.dirtyStatus = true
	}
}

func (r *Rank) tick() {
	now := Now()
	r.changeStatus(now.Unix())

	// 结算期， 奖励分发
	if r.Status == StatusSettling && !r.settledFired {
		r.settle()
	}

	r.dbSave(now)

	time.AfterFunc(tickDuration, func() {
		r.PostBack(r.tick)
	})
}

// 通知 game 在线玩家拉取奖励
func (r *Rank) broadcastRankEnd() {
	cluster.Brocast(serverType.Game, &ssmessage.RankEnd{
		RankID: proto.Int32(r.ID),
	})
}

func (r *Rank) settle() {
	if r.settledFired {
		return
	}

	zaplogger.GetSugar().Debugf("rank %d len(%d) settle start idx %d ", r.ID, r.zset.Len(), r.settledIdx)
	r.settledFired = true

	go func() {
		atomic.AddInt32(&r.asynCount, 1)

		errTimes := 0
		tableName := makeTableRole(r.ID)
		length := r.zset.Len()

		keyFields := make(map[string]map[string]interface{}, batchCount)
		var start, end int = 1, length
		var count int
		for r.settledIdx < length {
			select {
			case <-r.stopSign:
				goto endFor
			default:

			}

			if r.settledIdx != 0 {
				start = r.settledIdx + 1
			}
			end = start + batchCount - 1
			count = 0
			r.zset.Range(start, end, func(idx int, key sortedset.Key, value interface{}) bool {
				keyFields[string(key)] = map[string]interface{}{
					fieldRankRoleIdx: idx,
				}
				count++
				return true
			})

			if err := dbClient.UpsertBatch(tableName, keyFields); err == nil {
				r.settledIdx += count
				_ = r.saveSettledIdx()
			} else {
				zaplogger.GetSugar().Error(err)
				errTimes++
				// 失败尝试3次
				if errTimes >= 3 {
					r.settledIdx += count
					_ = r.saveSettledIdx()
					errTimes = 0
				}
			}
		}

	endFor:

		if r.settledIdx == length {
			r.PostBack(func() {
				r.Status = StatusDestroy
				_ = r.saveStatus()
				r.broadcastRankEnd()
				//r.zset.Init()
				zaplogger.GetSugar().Debugf("rank %d settle end size(%d)", r.ID, length)
			})
		}
		atomic.AddInt32(&r.asynCount, -1)
	}()

}

// set
func (r *Rank) Set(roleID uint64, score int32, roleInfo *message.RankRoleInfo, callback func(rank int32, total int32, code int32)) {
	if r.Status != StatusBegin {
		callback(0, 0, 2)
		return
	}

	r.roleDirty[roleID] = roleInfo

	key := sortedset.Key(strconv.FormatUint(roleID, 10))
	oldIdx := r.zset.GetRank(key) - 1
	newIdx := r.zset.Set(key, Score(score)) - 1

	if oldIdx >= 0 {
		r.Total--
	}
	r.Total++

	if oldIdx != newIdx {
		topUpdate := false
		if oldIdx >= 0 && oldIdx < topN {
			topUpdate = true
			r.topCount--
			r.topRoles = append(r.topRoles[:oldIdx], r.topRoles[oldIdx+1:]...)
		}

		if newIdx >= 0 && newIdx < topN {
			topUpdate = true
			r.topCount++
			list := r.topRoles
			r.topRoles = make([]*message.RankRoleInfo, 0, topN)
			r.topRoles = append(r.topRoles, list[:newIdx]...)
			r.topRoles = append(r.topRoles, roleInfo)
			r.topRoles = append(r.topRoles, list[newIdx:]...)
		}

		if r.topCount > topN {
			r.topCount = topN
			r.topRoles = r.topRoles[:topN-1]
		}
		if topUpdate {
			r.topVersion++
		}
	} else {
		r.topRoles[oldIdx] = roleInfo
		r.topVersion++
	}

	zaplogger.GetSugar().Debug("setScore ", r.ID, score, oldIdx, newIdx, r.topCount, r.topVersion)

	callback(int32(newIdx)+1, r.Total, 1)
}

func (r *Rank) GetRank(roleIDs []uint64, callback func(results []*rpc.RankGetRankResult)) {
	// 0 没有玩家数据， 1 排行榜正在进行， 2 排行榜完结，且有排名
	results := make([]*rpc.RankGetRankResult, 0, len(roleIDs))

	for _, roleID := range roleIDs {
		result := &rpc.RankGetRankResult{
			Code:   proto.Int32(0),
			RoleID: proto.Uint64(roleID),
			Rank:   proto.Int32(0),
			Total:  proto.Int32(0),
		}
		if r.Status == StatusNoBegin {
			result.Code = proto.Int32(0)
		} else {
			key := sortedset.Key(strconv.FormatUint(roleID, 10))
			rank := r.zset.GetRank(key)
			if rank == 0 {
				// 无数据
				result.Code = proto.Int32(0)
			} else {
				result.Rank = proto.Int32(int32(rank))
				result.Total = proto.Int32(r.Total)
				result.Code = proto.Int32(1)
				score, _ := r.zset.GetValue(key)
				result.Score = proto.Int32(int32(score.(Score)))
				if r.Status == StatusDestroy {
					result.Code = proto.Int32(2)
				}
			}
		}
		results = append(results, result)
	}

	callback(results)
}

func (r *Rank) GetRankTop(roleID uint64, version int32, callback func(top *message.RankGetTopListToC)) {
	// 默认失败
	toc := &message.RankGetTopListToC{Version: proto.Int32(0)}

	if r.Status != StatusNoBegin {
		key := sortedset.Key(strconv.FormatUint(roleID, 10))
		rank := r.zset.GetRank(key)
		toc.Total = proto.Int32(r.Total)
		toc.Rank = proto.Int32(int32(rank))
	}

	if r.topVersion != version {
		toc.Version = proto.Int32(r.topVersion)
		toc.Roles = r.topRoles[:r.topCount]
		toc.RankInfo = &message.RankInfo{
			RankID:    proto.Int32(r.ID),
			BeginTime: proto.Int64(r.BeginTime),
			EndTime:   proto.Int64(r.EndTime),
		}
	}
	callback(toc)
}
