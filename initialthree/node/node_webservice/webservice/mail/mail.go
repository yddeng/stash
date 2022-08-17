package mail

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/node/common/db"
	"initialthree/node/common/offlinemsg"
	"initialthree/node/common/serverType"
	"initialthree/node/node_webservice/webservice"
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/zaplogger"
	"time"
)

type module struct {
	dbClient offlinemsg.DB
}

func (m *module) addMail(wait *webservice.WaitConn, req struct {
	Type   string        `json:"type"`
	GameID []uint64      `json:"gameId"`
	Mail   *message.Mail `json:"mail"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	if req.Mail == nil {
		wait.SetResult("user add mail, but not set mail", nil)
		wait.Done()
		return
	}

	global := true
	if req.Type == "user" {
		if len(req.GameID) == 0 {
			wait.SetResult("user add mail, but not set gameID", nil)
			wait.Done()
			return
		}
		global = false
	}

	req.Mail.CreateTime = proto.Int64(time.Now().Unix())

	bytes, _ := json.Marshal(req.Mail)

	cb := func(err errcode.Error, global bool, gameID []uint64) {
		if err == nil {
			msg := &ssmessage.MailUpdate{
				Global: proto.Bool(global),
				GameID: gameID,
			}
			cluster.Brocast(serverType.Game, msg)
		} else {
			wait.SetResult(errcode.GetErrorDesc(err), nil)
		}
		wait.Done()
	}

	if global {
		offlinemsg.PushMsg(m.dbClient, "global", "mail", bytes, time.Second*6, func(i errcode.Error, i2 int64) {
			cb(i, true, []uint64{})
		})
	} else {
		gameIDs := make([]uint64, 0, len(req.GameID))
		doneTimes := len(req.GameID)
		for _, gid := range req.GameID {
			gameID := gid
			name := fmt.Sprintf("%d", gameID)
			offlinemsg.PushMsg(m.dbClient, name, "mail", bytes, time.Second*6, func(i errcode.Error, i2 int64) {
				if i == nil {
					gameIDs = append(gameIDs, gameID)
				}
				doneTimes--
				if doneTimes == 0 {
					cb(nil, false, gameIDs)
				}
			})
		}
	}

}

func (m *module) Init(app *gin.Engine) error {
	group := app.Group("/mail")
	group.POST("/add", webservice.WarpHandle(_module.addMail))

	m.dbClient = offlinemsg.NewFlyfishDB(db.GetFlyfishClient("global"))

	// test

	return nil
}

func (m *module) Tick() {
}

var _module *module

func init() {
	_module = new(module)
	webservice.RegisterModule("mail", _module)
}
