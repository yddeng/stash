package rank

import (
	"fmt"
	"initialthree/node/common/config"
	"initialthree/node/common/db"
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"strconv"
)

var dbClient *db.Client

func InitDBClient(dbConfig *config.Rank) error {
	c, err := db.NewClient(dbConfig.SqlType, dbConfig.DbHost, dbConfig.DbPort, dbConfig.DbDataBase, dbConfig.DbUser, dbConfig.DbPassword)
	if err != nil {
		zaplogger.GetSugar().Error(err)
		return err
	}

	c.SetMaxOpenConns(dbClientCount)
	dbClient = c
	return nil
}

func makeTableRole(id int32) string {
	return fmt.Sprintf("%s%d", tableRankRolePrefix, id)
}

/*
 根据排行榜起止时间， 创建、删除 玩家分数表
*/
var createSqlStr = `
CREATE TABLE "%s" (
	"__key__" varchar(255) NOT NULL,
    "__version__" int8 NOT NULL DEFAULT 0,
 	"%s" bytea NOT NULL DEFAULT '',
    "%s" int8 NOT NULL DEFAULT 0,
    "%s" int8 NOT NULL DEFAULT 0,
 	PRIMARY KEY ("__key__")
);
`

func createTableRankRole(id int32) error {
	_ = dropTableRankRole(id)
	tableName := makeTableRole(id)
	sql := fmt.Sprintf(createSqlStr, tableName, fieldRankRoleInfo, fieldRankRoleScore, fieldRankRoleIdx)
	return dbClient.Exec(sql)
}

func dropTableRankRole(id int32) error {
	tableName := makeTableRole(id)
	sql := fmt.Sprintf(`DROP TABLE IF EXISTS "%s";`, tableName)
	return dbClient.Exec(sql)
}

// 清理历史排行榜。排行榜创建成功后，异步调用
func clearHistoryRank(nowRankId int32) {
	idToLastId := map[int32]int32{}
	err := dbClient.GetAll(tableRankList, func(ret map[string]interface{}) error {
		id, err := strconv.ParseInt(ret[fieldKey].(string), 10, 32)
		if err != nil {
			return err
		}
		lastId := int32(ret[fieldRankListLastID].(int64))
		idToLastId[int32(id)] = lastId
		return nil
	})
	if err != nil {
		zaplogger.GetSugar().Error(err)
		return
	}

	// 移除掉阀值之前的所有历史版本
	links := makeLink(idToLastId)
	headNode, _ := findLink(nowRankId, links)
	if headNode != nil {
		clearIds := headNode.getNodes(keepHisListCount, 0)
		if len(clearIds) == 0 {
			return
		}
		n := len(clearIds)
		clearIds = clearIds[:n-1]
		zaplogger.GetSugar().Debug("clearHistoryRank", clearIds)
		// 移除数据库 历史版本
		for _, id := range clearIds {
			err := dropTableRankRole(id)
			if err != nil {
				zaplogger.GetSugar().Error(err)
			}

			err = deleteRankList(id)
			if err != nil {
				zaplogger.GetSugar().Error(err)
			}
		}

		// 移除实例
		remIds := headNode.getNodes(3, 0)
		for _, id := range remIds {
			Remove(id)
		}
	}

}

// 获取玩家详细信息
func getRoleInfo(rankID int32, roleId uint64) *message.RankRoleInfo {
	tableName := makeTableRole(rankID)
	key := fmt.Sprintf("%d", roleId)
	ret, err := dbClient.Get(tableName, key)
	if err != nil {
		return nil
	}
	info := ret[fieldRankRoleInfo]
	if info == nil || len(info.([]byte)) == 0 {
		return nil
	}

	var role *message.RankRoleInfo
	err = json.Unmarshal(info.([]byte), &role)
	if err != nil {
		return nil
	}

	return role
}

// 数据查询结算排名, rank,  exit
/*
 排行榜还在进行时，玩家只有分数信息没有排名信息。排名信息在排行榜到期时结算。
*/
type roleRankInfo struct {
	Rank  int32
	Score int32
}

func RankRoleRank(rankID int32, roleIDs []uint64) (ranks map[uint64]*roleRankInfo) {
	tableName := makeTableRole(rankID)

	keys := make([]string, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		keys = append(keys, fmt.Sprintf("%d", roleID))
	}

	ranks = make(map[uint64]*roleRankInfo, len(roleIDs))
	_ = dbClient.GetBatch(tableName, keys, func(row map[string]interface{}) error {
		if roleID, err := strconv.ParseUint(row[fieldKey].(string), 10, 64); err == nil {
			info := &roleRankInfo{
				Rank:  int32(row[fieldRankRoleIdx].(int64)),
				Score: int32(row[fieldRankRoleScore].(int64)),
			}
			ranks[roleID] = info
		}
		return nil
	}, fieldRankRoleScore, fieldRankRoleIdx)
	return
}

func RankRoleCount(rankID int32) int32 {
	tableName := makeTableRole(rankID)
	cnt, err := dbClient.Count(tableName)
	if err != nil {
		return 0
	}
	return int32(cnt)
}

// 查询排行榜的状态
func RankListStatus(rankId int32) (RankStatus, bool) {
	key := fmt.Sprintf("%d", rankId)

	ret, err := dbClient.Get(tableRankList, key)
	if err != nil || len(ret) == 0 {
		return StatusNoBegin, false
	}
	return RankStatus(ret[fieldRankListStatus].(int64)), true
}

// 查询排行榜创建的逻辑地址
func RankListLogicAddr(rankId int32) (string, bool) {
	key := fmt.Sprintf("%d", rankId)

	ret, err := dbClient.Get(tableRankList, key)
	if err != nil || len(ret) == 0 {
		return "", false
	}
	return ret[fieldRankListAddr].(string), true
}

// 获取上一期排行榜ID，地址
func RankListLastID(rankId int32) (int32, bool) {
	key := fmt.Sprintf("%d", rankId)

	ret, err := dbClient.Get(tableRankList, key)
	if err != nil || len(ret) == 0 {
		return 0, false
	}
	lastId := int32(ret[fieldRankListLastID].(int64))
	return lastId, lastId != 0
}
