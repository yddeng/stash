package whitelist

import (
	"github.com/gin-gonic/gin"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/node/common/db"
	"initialthree/node/common/omm/whitelist"
	"initialthree/node/node_webservice/webservice"
	"initialthree/zaplogger"
)

type module struct{}

func (m *module) set(wait *webservice.WaitConn, req struct {
	UserID      string `json:"userId"`
	ExpiredTime int64  `json:"expiredTime"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	if req.UserID == "" {
		wait.SetResult("userID is nil", nil)
		wait.Done()
		return
	}

	whitelist.SetCmd(db.GetFlyfishClient("global"), req.UserID).AsyncExec(func(result *client.StatusResult) {
		if errcode.GetCode(result.ErrCode) != errcode.Errcode_ok {
			zaplogger.GetSugar().Errorf("update white list failed: %s", errcode.GetErrorDesc(result.ErrCode))
			wait.SetResult(errcode.GetErrorDesc(result.ErrCode), nil)
		}
		wait.Done()
	})

}

func (m *module) del(wait *webservice.WaitConn, req struct {
	UserID string `json:"userId"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	if req.UserID == "" {
		wait.SetResult("userID is nil", nil)
		wait.Done()
		return
	}

	whitelist.DelCmd(db.GetFlyfishClient("global"), req.UserID).AsyncExec(func(result *client.StatusResult) {
		if !(errcode.GetCode(result.ErrCode) == errcode.Errcode_ok || errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist) {
			zaplogger.GetSugar().Errorf("delete white list failed: %s", errcode.GetErrorDesc(result.ErrCode))
			wait.SetResult(errcode.GetErrorDesc(result.ErrCode), nil)
		}
		wait.Done()
	})

}

func (m *module) Init(app *gin.Engine) error {
	group := app.Group("/whitelist")
	group.POST("/set", webservice.WarpHandle(m.set))
	group.POST("/del", webservice.WarpHandle(m.set))

	return nil
}

func (m *module) Tick() {}

var _module *module

func init() {
	_module = new(module)
	webservice.RegisterModule("whitelist", _module)
}
