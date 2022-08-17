package node_login

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/cs"
	"initialthree/node/common/config"
	"initialthree/node/common/db"
	"initialthree/node/common/omm/bannedlist"
	"initialthree/node/common/omm/whitelist"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"strconv"
	"strings"
	"time"

	"initialthree/zaplogger"
)

var loginAddr string

func Start(externalAddr string, comm *config.Common) error {
	loginAddr = externalAddr

	//if err := firewall.Init(db.GetFlyfishClient("login"), new(firewallListener), zaplogger.GetSugar()); err != nil {
	//	return fmt.Errorf("init firewall: %s", err)
	//}

	if err := initServerStatus(); err != nil {
		return err
	}

	reportStatus(comm.ServerGroups)

	t := strings.Split(externalAddr, ":")
	port, _ := strconv.Atoi(t[1])
	return cs.StartTcpServer("tcp", fmt.Sprintf("0.0.0.0:%d", port), &gDispatcher)
}

func init() {
	gDispatcher.handlers = map[uint16]handler{}
	gDispatcher.Register(cmdEnum.CS_Login, func(session *fnet.Socket, msg *codecs.Message) {
		login := msg.GetData().(*cs_message.LoginToS)
		zaplogger.GetSugar().Info(login.GetUserID(), " OnLogin")
		seq := msg.GetSeriNo()

		f := func() {
			bannedlist.AuthUserID(db.GetFlyfishClient("global"), login.GetUserID(), func(banned bool) {
				if banned {
					zaplogger.GetSugar().Infof("user %s OnLogin, was banned.", login.GetUserID())
					_ = session.Send(codecs.ErrMessage(seq, msg.GetCmd(), uint16(cs_message.ErrCode_Banned)))
					session.Close(errors.New("user-id banned"), time.Second)
					return
				}

				GetGate(login.GetUserID(), func(err error, gate string, token string) {
					loginResp := &cs_message.LoginToC{}
					if nil == err {
						zaplogger.GetSugar().Infof("user %s login ok", login.GetUserID())
						loginResp.Game = proto.String(gate)
						loginResp.Token = proto.String(token)
						loginResp.UserID = proto.String(login.GetUserID())
						_ = session.Send(codecs.NewMessage(seq, loginResp))
					} else {
						zaplogger.GetSugar().Infof(err.Error())
						errCode := cs_message.ErrCode_ERROR
						_ = session.Send(codecs.ErrMessage(seq, msg.GetCmd(), uint16(errCode)))
					}

					session.Close(nil, time.Second)
				})
			})
		}

		if !isServerOpen() {
			whitelist.AuthUserID(db.GetFlyfishClient("global"), login.GetUserID(), func(exist bool) {
				if !exist {
					zaplogger.GetSugar().Infof("user %s OnLogin, server not open and whitelist auth not pass", login.GetUserID())
					_ = session.Send(codecs.ErrMessage(seq, msg.GetCmd(), uint16(cs_message.ErrCode_SERVER_MAINTAINED)))
					session.Close(errors.New("server not open"), time.Millisecond*100)
					return
				}
				f()
			})
		} else {
			f()
		}

	})

}

type firewallListener struct{}

//func (f firewallListener) OnUpdated(firewall.UpdatedArgs) {
//	fwSt := firewall.GetStatus()
//	zaplogger.GetSugar().Infof("firewall updated: %s", fwSt.String())
//}

func (f firewallListener) OnReloadError(err error) {
	zaplogger.GetSugar().Errorf("firewall reload: %s", err)
}
