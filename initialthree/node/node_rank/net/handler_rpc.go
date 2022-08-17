package net

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/node/node_rank/rank"
	"initialthree/pkg/timer"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/rpc"
	"initialthree/rpc/rankCreate"
	"initialthree/rpc/rankGetRank"
	"initialthree/rpc/rankGetTopList"
	"initialthree/rpc/rankSetScore"
	"initialthree/zaplogger"
	"time"
)

type rankCreate_ struct {
}

func (_ *rankCreate_) OnCall(replyer *rankCreate.RankCreateReplyer, arg *rpc.RankCreateReq) {
	// code = 1; // 0 失败，1 已经存在，2 成功, 3 服务上限

	// 其他状态，拒绝所有 请求
	if !rank.Mgr().Running() {
		replyer.Reply(&rpc.RankCreateResp{Code: proto.Int32(0)})
		return
	}

	zaplogger.GetSugar().Infof("rank create %v, gameLogic %s", arg, replyer.GetChannel().Name())

	logicAddr := cluster.SelfAddr().Logic.String()
	r := rank.Mgr().GetRank(arg.GetRankID())
	if r != nil {
		zaplogger.GetSugar().Errorf("rankCreate, rank %d is exit", arg.GetRankID())
		replyer.Reply(&rpc.RankCreateResp{Code: proto.Int32(1), LogicAddr: proto.String(logicAddr)})
		return
	}
	replyer.Reply(rank.Mgr().Create(arg, logicAddr))
}

type rankSetScore_ struct{}

func (_ *rankSetScore_) OnCall(replyer *rankSetScore.RankSetScoreReplyer, arg *rpc.RankSetScoreReq) {
	// 0 继续尝试，1 成功 2 失败

	// 其他状态，拒绝所有 请求
	if !rank.Mgr().Running() {
		replyer.Reply(&rpc.RankSetScoreResp{Code: proto.Int32(0)})
		return
	}

	//zaplogger.GetSugar().Infof("rank setRank %v ", arg)

	r := rank.Mgr().GetRank(arg.GetRankID())
	if r == nil {
		replyer.Reply(&rpc.RankSetScoreResp{Code: proto.Int32(2)})
		return
	}

	r.PostBack(func() {
		r.Set(arg.GetRoleID(), arg.GetScore(), arg.GetRoleInfo(), func(rank int32, total int32, code int32) {
			replyer.Reply(&rpc.RankSetScoreResp{
				Code:  proto.Int32(code),
				Rank:  proto.Int32(rank),
				Total: proto.Int32(total),
			})
		})
	})
}

type rankGetRank_ struct{}

func (_ *rankGetRank_) OnCall(replyer *rankGetRank.RankGetRankReplyer, arg *rpc.RankGetRankReq) {
	// 0 没有玩家数据， 1 排行榜正在进行， 2 排行榜完结，且有排名

	// 默认失败
	resp := &rpc.RankGetRankResp{
		Ok:      proto.Bool(false),
		RankID:  proto.Int32(arg.GetRankID()),
		Results: make([]*rpc.RankGetRankResult, 0, len(arg.GetRoleID())),
	}

	// 其他状态，拒绝所有 请求
	if !rank.Mgr().Running() {
		replyer.Reply(resp)
		return
	}

	//zaplogger.GetSugar().Infof("rank getRank %v ", arg)

	rankID := arg.GetRankID()
	roleIDs := arg.GetRoleID()

	r := rank.Mgr().GetRank(rankID)
	if r == nil {
		// 不在本服 从数据库加载

		// 设置超时时间，5s 没有返回就结束
		replied := false
		timeoutTimer := cluster.RegisterTimerOnce(time.Second*5, func(timer *timer.Timer, i interface{}) {
			if !replied {
				zaplogger.GetSugar().Infof("rank getRank %v timeout", arg)
				replyer.Reply(resp)
				replied = true
			}
		}, nil)

		selfLogicAddr := cluster.SelfAddr().Logic.String()

		go func() {
			// 查询服务地址
			rankAddr, ok := rank.RankListLogicAddr(rankID)
			if !ok {
				// 排行榜已移除
				cluster.PostTask(func() {
					if !replied {
						timeoutTimer.Cancel()
						replied = true
						resp.Ok = proto.Bool(true)
						for _, roleID := range roleIDs {
							resp.Results = append(resp.Results, &rpc.RankGetRankResult{
								Code:   proto.Int32(0),
								RoleID: proto.Uint64(roleID),
							})
						}
						replyer.Reply(resp)
					}
				})
			} else if selfLogicAddr == rankAddr {
				// 在本服服务
				roleRanks := rank.RankRoleRank(rankID, roleIDs)
				total := rank.RankRoleCount(rankID)
				cluster.PostTask(func() {
					if !replied {
						timeoutTimer.Cancel()
						replied = true
						resp.Ok = proto.Bool(true)
						for _, roleID := range roleIDs {
							result := &rpc.RankGetRankResult{
								Code:   proto.Int32(0),
								RoleID: proto.Uint64(roleID),
							}
							if info, exist := roleRanks[roleID]; exist {
								result.Score = proto.Int32(info.Score)
								if info.Rank == 0 {
									result.Code = proto.Int32(1)
								} else {
									result.Code = proto.Int32(2)
									result.Rank = proto.Int32(info.Rank)
									result.Total = proto.Int32(total)
								}
							}
							resp.Results = append(resp.Results, result)
						}
						replyer.Reply(resp)
					}
				})
			} else {
				// 不在本服服务，重定向地址
				cluster.PostTask(func() {
					if !replied {
						timeoutTimer.Cancel()
						replied = true
						resp.Ok = proto.Bool(true)
						resp.RedirectRankAddr = proto.String(rankAddr)
						replyer.Reply(resp)
					}
				})
			}
		}()
	} else {
		r.PostBack(func() {
			r.GetRank(roleIDs, func(results []*rpc.RankGetRankResult) {
				resp.Ok = proto.Bool(true)
				resp.Results = results
				replyer.Reply(resp)
			})
		})
	}

}

type rankGetTopList_ struct{}

func (_ *rankGetTopList_) OnCall(replyer *rankGetTopList.RankGetTopListReplyer, arg *rpc.RankGetTopListReq) {
	if !rank.Mgr().Running() {
		replyer.Reply(&rpc.RankGetTopListResp{Toc: &message.RankGetTopListToC{}})
		return
	}

	// 默认失败
	resp := &rpc.RankGetTopListResp{
		Toc: &message.RankGetTopListToC{
			Version: proto.Int32(0),
		},
	}

	getLast := arg.GetTos().GetGetLast()
	rankId := arg.GetTos().GetRankID()

	if getLast {
		if rankId2, ok := rank.RankListLastID(rankId); ok {
			rankId = rankId2
		} else {
			replyer.Reply(resp)
			return
		}
	}

	logic := cluster.SelfAddr().Logic.String()
	r := rank.Mgr().GetRank(rankId)
	if r == nil {
		if addr, ok := rank.RankListLogicAddr(rankId); ok {
			if logic != addr {
				// 重定向到服务节点
				resp.RedirectRankAddr = proto.String(addr)
			}
		}
		resp.CurRankID = proto.Int32(rankId)
		replyer.Reply(resp)
	} else {
		r.PostBack(func() {
			r.GetRankTop(arg.GetRoleID(), arg.GetTos().GetVersion(), func(top *message.RankGetTopListToC) {
				resp.Toc = top
				replyer.Reply(resp)
			})
		})
	}
}

/*
type rankDeleteScore_ struct{}

func (_ *rankDeleteScore_) OnCall(replyer *rankDeleteScore.RankDeleteScoreReplyer, arg *rpc.RankDeleteScoreReq) {
	resp := &rpc.RankDeleteScoreResp{Ok: proto.Bool(false)}

	// 其他状态，拒绝所有 请求
	if !rank2.Mgr().Running() {
		replyer.Reply(resp)
		zaplogger.GetSugar().Debug("rankDeleteScore_ rank stop")
		return
	}

	zaplogger.GetSugar().Infof("rank rankDeleteScore %v", arg)

	r := Mgr().GetRank(arg.GetRankID())
	if r != nil {
		zaplogger.GetSugar().Errorf("rankDeleteScore, rank %d is not exit", arg.GetRankID())
		replyer.Reply(resp)
		return
	}

	r.PostBack(func() {
		ok := r.Delete(arg.GetRoleID())
		resp.Ok = proto.Bool(ok)
		replyer.Reply(resp)
	})

}
*/

func init() {
	rankCreate.Register(&rankCreate_{})
	rankSetScore.Register(&rankSetScore_{})
	rankGetRank.Register(&rankGetRank_{})
	rankGetTopList.Register(&rankGetTopList_{})
	//rankDeleteScore.Register(new(rankDeleteScore_))
}
