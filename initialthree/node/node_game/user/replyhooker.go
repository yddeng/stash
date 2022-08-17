package user

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/codec/cs"
	mail2 "initialthree/node/node_game/global/mail"
	"initialthree/node/node_game/global/scarsIngrain"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/mail"
	"initialthree/node/node_game/module/rankData"
	"initialthree/protocol/ss/rpc"
	"initialthree/zaplogger"
	"reflect"

	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/ss/ssmessage"
)

type Hooker func(addr.LogicAddr, *User, *cs.Message) bool

type relayMsgHooker struct {
	receiver  *cs.Receiver
	msgHooker map[string]Hooker
}

var relayHooker = relayMsgHooker{
	receiver:  cs.NewReceiver("sc"),
	msgHooker: map[string]Hooker{},
}

func hookMsg(msgType string, unpackMsg *cs.Message, u *User, logicAddr addr.LogicAddr) bool {
	h, ok := relayHooker.msgHooker[msgType]
	if ok {
		return h(logicAddr, u, unpackMsg)
	}
	return true
}

func RegisterHooker(msg proto.Message, h Hooker) {
	msgType := reflect.TypeOf(msg).String()
	_, ok := relayHooker.msgHooker[msgType]
	if !ok {
		relayHooker.msgHooker[msgType] = h
	}
}

// 服务器内部穿透消息
func onReply(from addr.LogicAddr, msg proto.Message) {
	req := msg.(*ssmessage.Relay)
	targets := req.GetTargets()

	var unpackMsg *cs.Message
	_, ok := relayHooker.msgHooker[req.GetMsgType()]
	if ok {
		p, err := relayHooker.receiver.DirectUnpack(req.GetMsg())
		if err != nil {
			zaplogger.GetSugar().Debugf("directUnpack err:%v", err)
			return
		}
		unpackMsg = p.(*cs.Message)
	}

	for _, v := range targets {
		u := userMap[v.GetUserID()]
		if u != nil && u.checkStatus(status_playing) {
			if hookMsg(req.GetMsgType(), unpackMsg, u, from) {
				u.Reply(unpackMsg.GetSeriNo(), unpackMsg.GetData())
			}
		}
	}
}

// 排行榜结束，通知给在线玩家，拉取奖励。
func onRankEnd(from addr.LogicAddr, msg proto.Message) {
	req := msg.(*ssmessage.RankEnd)
	rankId := req.GetRankID()

	zaplogger.GetSugar().Debugf("onRankEnd %d end", rankId)

	callGetRank := func(logicAddr addr.LogicAddr, rankID int32, roleIDs []uint64) {
		rankData.CallGetRank(logicAddr, rankId, roleIDs, func(results []*rpc.RankGetRankResult, err error) {
			if err != nil {
				zaplogger.GetSugar().Errorf("onRankEnd callGetRank err %s", err.Error())
				return
			}

			for _, result := range results {
				u := GetIDUser(result.GetRoleID())
				if u != nil {
					rankModule := u.GetSubModule(module.RankData).(*rankData.RankData)
					// 0 没有玩家数据， 1 排行榜正在进行， 2 排行榜完结，且有排名
					switch result.GetCode() {
					case 0:
						zaplogger.GetSugar().Debugf("%s rankData rank %d is not exit", u.GetUserID(), rankID)
						rankModule.DelRank(rankID)
					case 2:
						rank, total := result.GetRank(), result.GetTotal()
						if rank > 0 && total > 0 {
							rankModule.ApplyRankAward(rankID, result.GetScore(), rank, total)
						}
						rankModule.SetApplyAward(rankID)
					default:
					}
				}
			}
		})
	}

	roles := make([]uint64, 0, 100)
	for _, u := range userMap {
		rankModule := u.GetSubModule(module.RankData).(*rankData.RankData)
		if rankModule.GetRankInfo(rankId) != nil {
			roles = append(roles, u.GetID())
		}

		if len(roles) == 100 {
			callGetRank(from, rankId, roles)
			roles = roles[0:0]
		}
	}

	if len(roles) > 0 {
		callGetRank(from, rankId, roles)
	}
}

func onScarsIngrainUpdate(_ addr.LogicAddr, arg proto.Message) {
	msg := arg.(*ssmessage.ScarsIngrainUpdate)
	zaplogger.GetSugar().Debugf("ScarsIngrainUpdate %v", msg)

	// 更新数据
	scarsIngrain.OnUpdate(msg.GetId(), msg.GetVersion())
}

func onMailUpdate(from addr.LogicAddr, msg proto.Message) {
	req := msg.(*ssmessage.MailUpdate)
	zaplogger.GetSugar().Debugf("onMailUpdate %v", msg)

	for _, gameID := range req.GetGameID() {
		u := GetIDUser(gameID)
		if u != nil {
			mailModule := u.GetSubModule(module.Mail).(*mail.Mail)
			mailModule.PullOfflineMailData()
		}
	}

	if req.GetGlobal() {
		mail2.OnUpdate()
		for _, u := range userMap {
			mailModule := u.GetSubModule(module.Mail).(*mail.Mail)
			mailModule.PullGlobalMailData()
		}
	}
}

func init() {
	cluster.Register(cmdEnum.SS_Relay, onReply)
	cluster.Register(cmdEnum.SS_RankEnd, onRankEnd)
	cluster.Register(cmdEnum.SS_MailUpdate, onMailUpdate)
	cluster.Register(cmdEnum.SS_ScarsIngrainUpdate, onScarsIngrainUpdate)
}
