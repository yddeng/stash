package node_gate

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/pkg/timer"
	_ "initialthree/protocol/ss" //触发pb注册
	ss_rpc "initialthree/protocol/ss/rpc"
	"initialthree/rpc/synctoken"
	"initialthree/zaplogger"
	"time"
)

type SyncToken struct {
	userTokenMap map[string]*userToken
}

type userToken struct {
	timer *timer.Timer
	token string
}

var syncToken *SyncToken = &SyncToken{
	userTokenMap: map[string]*userToken{},
}

func (this *SyncToken) OnCall(replyer *synctoken.SynctokenReplyer, arg *ss_rpc.SynctokenReq) {
	user := arg.GetUserid()
	token := arg.GetToken()

	token = this.addToken(user, token, time.Second*10)

	replyer.Reply(&ss_rpc.SynctokenResp{
		Token: proto.String(token),
	})

}

func (this *SyncToken) addToken(userID, token string, timeout time.Duration) string {
	t, ok := this.userTokenMap[userID]
	if !ok {
		t = &userToken{token: token}
		this.userTokenMap[userID] = t
	} else {
		token = t.token
	}

	if t.timer == nil || (t.timer != nil && !t.timer.ResetFireTime(timeout)) {
		t.timer = cluster.RegisterTimerOnce(timeout, func(tt *timer.Timer, _ interface{}) {
			t, ok := this.userTokenMap[userID]
			if ok && t.token == token && t.timer == tt {
				zaplogger.GetSugar().Debugf("removeToken u:%s,token:%s", userID, token)
				delete(this.userTokenMap, userID)

			}
		}, nil)
	}

	return token
}

func (this *SyncToken) checkToken(userID, token string) bool {
	zaplogger.GetSugar().Debugf("checkToken u:%s,token:%s", userID, token)
	t, ok := this.userTokenMap[userID]
	if !ok {
		return false
	} else {
		return t.token == token
	}
}

func CheckToken(userID, token string) bool {
	return syncToken.checkToken(userID, token)
}

func AddToken(userID, token string, timeout time.Duration) string {
	return syncToken.addToken(userID, token, timeout)
}

func init() {
	synctoken.Register(syncToken)
}
