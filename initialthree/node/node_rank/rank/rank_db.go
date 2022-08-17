package rank

import (
	"fmt"
	"github.com/yddeng/sortedset"
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"strconv"
	"sync/atomic"
	"time"
)

/*
 save
*/

// 创建表的时候使用，唯一保证。存在则失败
func (r *Rank) insertRankList(logicAddr string, lastId int32) error {
	infoData, _ := json.Marshal(r.RankBase)
	return dbClient.Set(tableRankList, fmt.Sprintf("%d", r.ID), map[string]interface{}{
		fieldRankListInfo:   infoData,
		fieldRankListAddr:   logicAddr,
		fieldRankListLastID: lastId,
		fieldRankListStatus: r.Status,
	})
}

func deleteRankList(rankId int32) error {
	return dbClient.Delete(tableRankList, fmt.Sprintf("%d", rankId))
}

func (r *Rank) saveStatus() error {
	return dbClient.Upsert(tableRankList, fmt.Sprintf("%d", r.ID), map[string]interface{}{
		fieldRankListStatus: r.Status,
	})

}

func (r *Rank) saveSettledIdx() error {
	return dbClient.Upsert(tableRankList, fmt.Sprintf("%d", r.ID), map[string]interface{}{
		fieldRankListSettled: r.settledIdx,
	})
}

// 玩家分数
func (r *Rank) saveRoleScore(id int64, score int, roleInfo *message.RankRoleInfo) error {
	tableName := makeTableRole(r.ID)
	fields := map[string]interface{}{
		fieldRankRoleScore: score,
	}
	if roleInfo != nil {
		data, err := json.Marshal(roleInfo)
		if err != nil {
			return err
		}
		fields[fieldRankRoleInfo] = data
	}

	key := fmt.Sprintf("%d", id)
	return dbClient.Upsert(tableName, key, fields)
}

// 批量更新
func (r *Rank) saveRoleDirty() {
	if len(r.roleDirty) == 0 {
		return
	}

	//zaplogger.GetSugar().Infoln(r.ID, "saveRoleDirty", len(r.roleDirty))

	roles := r.roleDirty
	r.roleDirty = map[uint64]*message.RankRoleInfo{}

	// 异步保存
	go func() {
		atomic.AddInt32(&r.asynCount, 1)
		// 保存失败数据
		unSaved := map[uint64]*message.RankRoleInfo{}

		tableName := makeTableRole(r.ID)
		saving := make(map[uint64]*message.RankRoleInfo, dirtyBatchCount)
		keyFields := make(map[string]map[string]interface{}, dirtyBatchCount)
		for id, roleInfo := range roles {
			saving[id] = roleInfo
			score := roleInfo.GetScore()
			fields := map[string]interface{}{fieldRankRoleScore: score}
			if roleInfo != nil {
				data, err := json.Marshal(roleInfo)
				if err != nil {
					zaplogger.GetSugar().Error(err)
					unSaved[id] = roleInfo
					continue
				}
				fields[fieldRankRoleInfo] = data
			}
			keyFields[fmt.Sprintf("%d", id)] = fields

			if len(keyFields) >= dirtyBatchCount {
				err := dbClient.UpsertBatch(tableName, keyFields)
				if err != nil {
					zaplogger.GetSugar().Error(err)
					for k, v := range saving {
						unSaved[k] = v
					}
				}
				saving = make(map[uint64]*message.RankRoleInfo, dirtyBatchCount)
				keyFields = make(map[string]map[string]interface{}, dirtyBatchCount)
			}
		}
		if len(keyFields) > 0 {
			err := dbClient.UpsertBatch(tableName, keyFields)
			if err != nil {
				zaplogger.GetSugar().Error(err)
				for k, v := range saving {
					unSaved[k] = v
				}
			}
		}

		// 保存失败的数据 重设脏标记
		if len(unSaved) > 0 {
			r.PostBack(func() {
				for id, roleInfo := range unSaved {
					if _, ok := r.roleDirty[id]; !ok {
						r.roleDirty[id] = roleInfo
					}
				}
			})
		}

		atomic.AddInt32(&r.asynCount, -1)
	}()

}

func (r *Rank) dbSave(now time.Time) {
	if r.dirtyStatus {
		if err := r.saveStatus(); err == nil {
			r.dirtyStatus = false
		}
	}

	if len(r.roleDirty) >= saveUpper || now.After(r.nextSaveTime) {
		r.saveRoleDirty()
		r.nextSaveTime = now.Add(saveDuration)
	}

}

// 排行榜重构
func (r *Rank) load() error {
	if err := dbClient.GetAll(makeTableRole(r.ID), func(row map[string]interface{}) error {
		score := int32(row[fieldRankRoleScore].(int64))
		roleID := row[fieldKey].(string)
		key := sortedset.Key(roleID)
		r.zset.Set(key, Score(score))
		return nil
	}); err != nil {
		return err
	}

	r.zset.Range(1, topN, func(rank int, key sortedset.Key, value interface{}) bool {
		roleID, _ := strconv.Atoi(string(key))
		role := getRoleInfo(r.ID, uint64(roleID))
		r.topRoles[rank-1] = role
		r.topCount++
		return true
	})

	r.Total = int32(r.zset.Len())

	// 只生成展示，不需要实例
	//if r.Status == StatusDestroy {
	//	r.zset.Init()
	//}
	return nil
}
