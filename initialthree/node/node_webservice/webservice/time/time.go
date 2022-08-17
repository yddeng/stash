package time

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/node/common/db"
	time2 "initialthree/node/common/omm/time"
	"initialthree/node/node_webservice/webservice"
	"initialthree/zaplogger"
	"time"
)

type module struct {
	version int64
	offsetV time.Duration
}

func (m *module) offset(wait *webservice.WaitConn, req struct {
	TargetTime string `json:"targetTime"`
	Reset      bool   `json:"reset"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	offsetV := time.Duration(0)
	if !req.Reset {
		var year, month, day, hour, min, sec int
		if n, err := fmt.Sscanf(req.TargetTime, "%d-%d-%d %d:%d:%d", &year, &month, &day, &hour, &min, &sec); n != 6 || err != nil {
			zaplogger.GetSugar().Errorf("[time/offset]: the targetTime:\"%s\" is invalid.", req.TargetTime)
			wait.SetResult("targetTime is invalid", nil)
			wait.Done()
			return
		}

		offsetV = time.Duration(time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local).Unix()-time.Now().Unix()) * time.Second
	}

	time2.SetCmd(db.GetFlyfishClient("global"), offsetV, m.version).AsyncExec(func(result *client.StatusResult) {
		if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
			_ = m.load()
		} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_version_mismatch {
			zaplogger.GetSugar().Infof("[time/offset]: update version %d miss match ,reload data", m.version)
			_ = m.load()
			wait.SetResult(errcode.GetErrorDesc(result.ErrCode), nil)
		} else {
			zaplogger.GetSugar().Errorf("[time/offset]: update failed: %s", errcode.GetErrorDesc(result.ErrCode))
			wait.SetResult(errcode.GetErrorDesc(result.ErrCode), nil)
		}
		wait.Done()
	})

}

func (m *module) load() error {
	result := time2.LoadCmd(db.GetFlyfishClient("global"), m.version).Exec()
	if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
		offsetV, err := time2.Unmarshal(result.Fields)
		if err == nil {
			m.offsetV = offsetV
			if result.Version != nil {
				m.version = *result.Version
			}
			zaplogger.GetSugar().Infof("load time offset ok")
		} else {
			return err
		}
	} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist || errcode.GetCode(result.ErrCode) == errcode.Errcode_record_unchange {
		zaplogger.GetSugar().Errorf("load time offset %s", errcode.GetErrorDesc(result.ErrCode))
	} else {
		zaplogger.GetSugar().Errorf("load time offset %s", errcode.GetErrorDesc(result.ErrCode))
		return errors.New(errcode.GetErrorDesc(result.ErrCode))
	}

	return nil
}

func (m *module) Init(app *gin.Engine) error {
	group := app.Group("/time")
	group.POST("/offset", webservice.WarpHandle(m.offset))
	return m.load()
}

func (m *module) Tick() {}

var _module *module

func init() {
	_module = new(module)
	webservice.RegisterModule("time", _module)
}
