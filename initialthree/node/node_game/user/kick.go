package user

import (
	"errors"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/network/smux"
	"initialthree/node/common/db"
	"initialthree/pkg/timer"
	"initialthree/protocol/cs/message"
	ss "initialthree/protocol/ss/ssmessage"
	"initialthree/zaplogger"
	"time"
)

// 来源有 game , gm控制台
func onKickGameUser(from addr.LogicAddr, msg proto.Message) {
	kick := msg.(*ss.KickGameUser)
	userID := kick.GetUserID()

	u := userMap[userID]
	if nil != u {
		zaplogger.GetSugar().Debugf("onKick %v %s %d", userID, "user status ->", u.status)
		u.sendKickToC()
		u.kick(false)
	} else {

		/*
		 *  user_game_login上玩家所在服务器指向本服，实际上玩家对象在本服已经被清除。
		 *  需要将user_game_login上玩家所在服清除，玩家才能在其它服正常登录
		 */

		set := db.GetFlyfishClient("game").CompareAndSetNx("user_game_login", userID, "gameaddr", cluster.SelfAddr().Logic.String(), "")
		set.AsyncExec(func(ret *flyfish.ValueResult) {
			zaplogger.GetSugar().Infof("%s clear db logout code:%s ", userID, errcode.GetErrorDesc(ret.ErrCode))
		})

		zaplogger.GetSugar().Debugf("onKick %s %s", userID, "user is nil")
	}
}

func (this *User) kick(saveTmp bool) {
	zaplogger.GetSugar().Debugf("user kick %s %d", this.userID, this.status)
	if this.checkStatus(status_login) {
		//login为保护状态，只有在login操作执行完毕后才允许执行logout流程
		this.signalKick = true
	} else {
		this.logout(saveTmp)
	}
}

func (this *User) sendKickToC() {
	if this.stream != nil {
		this.Post(&message.KickToC{})
		cluster.RegisterTimerOnce(time.Second, func(timer *timer.Timer, i interface{}) {
			stream := i.(*smux.MuxStream)
			stream.Close(errors.New("kicked"))
		}, this.stream)
	}

}

func sendKickGameUser(gameAddr addr.LogicAddr, userID string) {
	//向目标服发出踢人请求，并通知玩家重试
	msg := &ss.KickGameUser{
		UserID: proto.String(userID),
	}
	cluster.PostMessage(gameAddr, msg)
}
