package kick

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common/db"
	"initialthree/node/common/serverType"
	"initialthree/node/node_webservice/webservice"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/zaplogger"
)

type module struct {
}

func (m *module) kickUser(wait *webservice.WaitConn, req struct {
	UserID []string `json:"userID"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	if len(req.UserID) == 0 {
		wait.SetResult("userID is nil", nil)
		wait.Done()
		return
	}

	dbCli := db.GetFlyfishClient("game")
	cmdCount := len(req.UserID)

	for _, value := range req.UserID {
		userID := value
		dbCli.Get("user_game_login", userID, "gameaddr").AsyncExec(func(result *client.GetResult) {
			if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
				gameaddr := result.Fields["gameaddr"].GetString()
				if gameaddr == "" {
					zaplogger.GetSugar().Infof("kick: user %s not online.", userID)
				} else {
					if gameLogicAddr, err := addr.MakeLogicAddr(gameaddr); err != nil {
						zaplogger.GetSugar().Errorf("kick user %s: get user gamelogicaddr %s.", userID, err)
					} else {
						zaplogger.GetSugar().Infof("kick: notify game:%s to kick user %s.", gameaddr, userID)
						cluster.PostMessage(gameLogicAddr, &ssmessage.KickGameUser{UserID: proto.String(userID)})
					}
				}
			} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist {
				zaplogger.GetSugar().Infof("kick: user %s not online.", userID)
			} else {
				zaplogger.GetSugar().Infof("load user's gameaddr error: %s.", errcode.GetErrorDesc(result.ErrCode))
			}

			cmdCount--
			if cmdCount <= 0 {
				wait.Done()
			}
		})
	}
}

func (m *module) kickIP(wait *webservice.WaitConn, req struct {
	IP []string `json:"ip"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	if len(req.IP) == 0 {
		wait.SetResult("ip is nil", nil)
	} else {
		kickIPMsg := &ssmessage.KickIPGate{
			RegexpIPs: req.IP,
		}
		cluster.Brocast(serverType.Gate, kickIPMsg)
	}

	wait.Done()
}

func (m *module) Init(app *gin.Engine) error {
	group := app.Group("/kick")
	group.POST("/user", webservice.WarpHandle(m.kickUser))
	group.POST("/ip", webservice.WarpHandle(m.kickIP))
	return nil
}

func (m *module) Tick() {

}

var _module *module

func init() {
	_module = new(module)
	webservice.RegisterModule("kick", _module)
}
