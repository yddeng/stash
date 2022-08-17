package serverstatus

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/node/common/db"
	"initialthree/node/common/omm/serverstatus"
	"initialthree/node/node_webservice/webservice"
	"initialthree/zaplogger"
)

type module struct {
	version int64
	status  serverstatus.Status
}

func (m *module) set(wait *webservice.WaitConn, req struct {
	Closed bool `json:"closed"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	status := serverstatus.StatusOpen
	if req.Closed {
		status = serverstatus.StatusMaintenance
	}

	if status == m.status {
		wait.Done()
		return
	}

	serverstatus.SetCmd(db.GetFlyfishClient("global"), status, m.version).AsyncExec(func(result *client.StatusResult) {
		if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
			_ = m.load()
		} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_version_mismatch {
			zaplogger.GetSugar().Infof("update server_status version %d miss match ,reload data", m.version)
			_ = m.load()
			wait.SetResult(errcode.GetErrorDesc(result.ErrCode), nil)
		} else {
			zaplogger.GetSugar().Errorf("update server_status failed: %s", errcode.GetErrorDesc(result.ErrCode))
			wait.SetResult(errcode.GetErrorDesc(result.ErrCode), nil)
		}
		wait.Done()
	})

}

func (m *module) load() error {
	result := serverstatus.LoadCmd(db.GetFlyfishClient("global"), m.version).Exec()
	if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
		status, err := serverstatus.Unmarshal(result.Fields)
		if err == nil {
			m.status = status
			if result.Version != nil {
				m.version = *result.Version
			}
			zaplogger.GetSugar().Infof("load server_status ok")
		} else {
			return err
		}
	} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist || errcode.GetCode(result.ErrCode) == errcode.Errcode_record_unchange {
		zaplogger.GetSugar().Errorf("load server_status %s", errcode.GetErrorDesc(result.ErrCode))
		m.status = serverstatus.DefaultStatue
	} else {
		zaplogger.GetSugar().Errorf("load server_status %s", errcode.GetErrorDesc(result.ErrCode))
		return errors.New(errcode.GetErrorDesc(result.ErrCode))
	}

	return nil
}

func (m *module) Init(app *gin.Engine) error {
	group := app.Group("/serverstatus")
	group.POST("/set", webservice.WarpHandle(m.set))
	return m.load()
}

func (m *module) Tick() {

}

var _module *module

func init() {
	_module = new(module)
	webservice.RegisterModule("serverstatus", _module)
}
