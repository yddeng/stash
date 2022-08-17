package behavior

import (
	codecs "initialthree/codec/cs"
	csmsg "initialthree/protocol/cs/message"
	"initialthree/robot/config"
	"initialthree/robot/net"
	"initialthree/robot/statistics"
	"initialthree/robot/types"
	"time"

	. "github.com/GodYY/bevtree"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
)

const login = BevType("login")

func init() {
	regBevType(login, func() Bev { return new(BevLogin) })
}

type BevLogin struct {
	ReloginDelay uint `xml:"relogindelay"`
}

func (BevLogin) BevType() BevType { return login }

func (b *BevLogin) CreateInstance() BevInstance {
	return &BevLoginInstance{BevLogin: b}
}

func (b *BevLogin) DestroyInstance(bi BevInstance) {
	bi.(*BevLoginInstance).BevLogin = nil
}

type BevLoginInstance struct {
	bev
	*BevLogin
	startTime time.Time
	connector *net.Connector
	loginResp *csmsg.LoginToC
}

func (bd *BevLoginInstance) BevType() BevType { return login }

func (b *BevLoginInstance) OnInit(ctx Context) bool {
	b.bev.OnInit(ctx)

	if b.player.IsStatus(types.Status_IsLogin) {
		b.terminate(true)
		return true
	}

	b.startTime = time.Now()
	b.connectLogin()
	return true
}

func (b *BevLoginInstance) OnTerminate(ctx Context) {
	if b.connector != nil {
		b.connector.Stop()
		b.connector = nil
	}
	b.loginResp = nil
	b.player.RemTimer(timerIDRelogin)
	b.bev.OnTerminate(ctx)
}

func (b *BevLoginInstance) connectLogin() {
	b.player.Debugf("connect login")
	b.connector = net.ConnectService(config.GetConfig().Service, b.onConnectLogin, b.player.EventQueue())
}

func (b *BevLoginInstance) onConnectLogin(session *fnet.Socket, err error) {
	if err != nil {
		b.player.Errorf("connect login: %s", err)
		b.relogin()
		return
	}

	b.player.SetSession(session)

	loginMsg := &csmsg.LoginToS{
		UserID: proto.String(b.player.UserID()),
	}
	b.sendMessage(loginMsg, b.onLogin)

	b.player.Debugf("request to login")
}

func (b *BevLoginInstance) onLogin(r player, msg *codecs.Message) bool {
	if !net.IsMessageOK(msg) {
		b.player.Errorf("login failed: %s", net.GetErrCodeStr(msg.GetErrCode()))
		statistics.Login().Fail()
		b.relogin()
		return false
	}

	b.player.Infof("login successfully")
	b.player.CloseSession()
	b.loginResp = msg.GetData().(*csmsg.LoginToC)
	b.connectGate()

	return false
}

var timerIDRelogin = types.NewTimerID("relogin")

func (b *BevLoginInstance) reloginDelayDuration() time.Duration {
	return time.Duration(b.ReloginDelay) * time.Millisecond
}

func (b *BevLoginInstance) relogin() {
	b.player.CloseSession()

	if b.ReloginDelay > 0 {
		b.player.Infof("relogin %s later", b.reloginDelayDuration().String())
		b.player.AddTimer(timerIDRelogin, b.reloginDelayDuration(), nil, b.onRelogin)
	} else {
		b.player.Info("relogin")
		b.connectLogin()
	}
}

func (b *BevLoginInstance) onRelogin(r player, ctx interface{}) {
	b.connectLogin()
}

func (b *BevLoginInstance) connectGate() {
	b.player.Infof("connect gate:%s", b.loginResp.GetGame())
	b.connector = net.ConnectService(b.loginResp.GetGame(), b.onConnectGate, b.player.EventQueue())
}

func (b *BevLoginInstance) onConnectGate(session *fnet.Socket, err error) {
	if err != nil {
		b.player.Errorf("connect gate(%s): %s", b.loginResp.GetGame(), err)
		b.relogin()
		return
	}

	b.player.SetSession(session)

	b.gameLogin()
}

func (b *BevLoginInstance) onGameLogin(r player, msg *codecs.Message) bool {
	if !net.IsMessageOK(msg) {
		b.player.Errorf("gameLogin failed: %s", net.GetErrCodeStr(msg.GetErrCode()))

		if msg.GetErrCode() == uint16(csmsg.ErrCode_RETRY) {
			b.gameLogin()
		} else {
			statistics.Login().Fail()
			b.relogin()
		}

		return false
	}

	b.player.Infof("gameLogin successfully")

	statistics.Login().Success(time.Since(b.startTime))
	b.player.SetStatus(types.Status_IsLogin)
	resp := msg.GetData().(*csmsg.GameLoginToC)
	if resp.GetIsFirstLogin() {
		b.player.SetStatus(types.Status_IsFirstlogin)
	}

	// 登录成功，行为结束
	b.terminate(true)

	return false
}

func (b *BevLoginInstance) gameLogin() {
	gameLoginMsg := &csmsg.GameLoginToS{
		UserID:   proto.String(b.player.UserID()),
		Token:    proto.String(b.loginResp.GetToken()),
		ServerID: proto.Int32(int32(config.GetConfig().ServerID)),
		// SeqNo: ,
	}

	b.sendMessage(gameLoginMsg, b.onGameLogin)

	b.player.Debugf("request to gameLogin")
}
