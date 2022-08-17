package rank

import (
	"initialthree/node/common/timeDisposal"
	"time"
)

type RankOrder int

const (
	OrderDescending RankOrder = 0 // 降序
	OrderAscending  RankOrder = 1 // 升序
)

type RankUpdate int

const (
	UpdateLast   RankUpdate = 0 //
	UpdateBetter RankUpdate = 1
	UpdateSum    RankUpdate = 2
)

type RankStatus int

const (
	StatusNoBegin  = RankStatus(0) // 未开始
	StatusBegin    = RankStatus(1) // 排行期开始
	StatusEnd      = RankStatus(2) // 排行期结束
	StatusSettling = RankStatus(3) // 正在结算
	StatusDestroy  = RankStatus(4) // 销毁排行榜
)

const tickDuration = time.Millisecond * 500

const (
	topN             = 100              // 排行榜top榜人数
	saveUpper        = 10               // 脏标记阀值
	saveDuration     = time.Second * 30 // 脏标记定时储存
	batchCount       = 500              // db批量更新数量
	dirtyBatchCount  = 200              // db批量插入更新数量
	dbClientCount    = 15               // db连接池 连接上限
	keepHisListCount = 10               // 历史版本保留数量, >= 1
)

// 数据库表
const (
	tableRankList        = "rank_list"
	fieldRankListInfo    = "rank_info"
	fieldRankListStatus  = "rank_status"
	fieldRankListSettled = "settled_idx"
	fieldRankListAddr    = "logic_addr"
	fieldRankListLastID  = "last_id"

	tableRankRolePrefix = "rank_role_"
	fieldRankRoleScore  = "role_score"
	fieldRankRoleInfo   = "role_info"
	fieldRankRoleIdx    = "role_idx" // 排名，排行榜结束时设置

	fieldKey = "__key__"
)

func Now() time.Time {
	return timeDisposal.Now()
}
