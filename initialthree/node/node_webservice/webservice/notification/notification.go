package notification

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/node/common/db"
	"initialthree/node/common/omm/serverstatus"
	"initialthree/node/node_webservice/webservice"
	"initialthree/zaplogger"
)

// 登陆界面的公告
// 常规公告、停服公告。二选一下发

// 判断是否停服，检测集群中 login 服是否存在

const (
	tableName      = "web_data"
	tableKey       = "notification"
	onlineSlotIdx  = 0
	offlineSlotIdx = 1
)

var slots = []int{onlineSlotIdx, offlineSlotIdx}

func slotName(slotIdx int) string {
	return fmt.Sprintf("slot%d", slotIdx)
}

type Notification struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type module struct {
	version       int64
	Notifications []*Notification `json:"notifications"` // 长度为2 第一位是正常公告，第二位是停服公告
	IsClosed      bool            `json:"isClosed"`
	serverStatus  serverstatus.Status
}

func (m *module) updateClosed(status serverstatus.Status) {
	m.serverStatus = status
	m.IsClosed = true
	if m.serverStatus == serverstatus.StatusOpen {
		m.IsClosed = false
	}
}

func (m *module) update(wait *webservice.WaitConn, req struct {
	Type         string        `json:"type"`
	Notification *Notification `json:"notification"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	bytes, _ := json.Marshal(req.Notification)
	fieldName := slotName(onlineSlotIdx)
	slotIdx := onlineSlotIdx
	if req.Type == "offline" {
		slotIdx = offlineSlotIdx
		fieldName = slotName(offlineSlotIdx)
	}

	db.GetFlyfishClient("global").Set(tableName, tableKey, map[string]interface{}{
		fieldName: bytes,
	}).AsyncExec(func(result *client.StatusResult) {
		if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
			m.Notifications[slotIdx] = req.Notification
			wait.Done()
		} else {
			wait.SetResult(errcode.GetErrorDesc(result.ErrCode), nil)
			wait.Done()
		}
	})
}

func (m *module) getAll(wait *webservice.WaitConn) {
	//zaplogger.GetSugar().Infof("%s", wait.GetRoute())

	wait.SetResult("", m)
	wait.Done()
}

type getResult struct {
	IsClosedGameServer bool   `json:"isClosedGameServer"`
	Title              string `json:"title"`
	Content            string `json:"content"`
}

func (m *module) get(wait *webservice.WaitConn) {
	//zaplogger.GetSugar().Infof("%s", wait.GetRoute())

	get := &getResult{
		IsClosedGameServer: m.IsClosed,
	}
	if m.IsClosed {
		get.Title = m.Notifications[offlineSlotIdx].Title
		get.Content = m.Notifications[offlineSlotIdx].Content
	} else {
		get.Title = m.Notifications[onlineSlotIdx].Title
		get.Content = m.Notifications[onlineSlotIdx].Content
	}

	wait.SetResult("", get)
	wait.Done()
}

func (m *module) load() error {
	result := db.GetFlyfishClient("global").GetAllWithVersion(tableName, tableKey, m.version).Exec()
	if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok || errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist {
		for _, i := range slots {
			fieldName := slotName(i)
			field, ok := result.Fields[fieldName]
			if ok && len(field.GetBlob()) > 0 {
				err := json.Unmarshal(field.GetBlob(), &m.Notifications[i])
				if err != nil {
					zaplogger.GetSugar().Error(err.Error(), string(field.GetBlob()))
					return err
				}
			} else {
				m.Notifications[i] = &Notification{
					Title:   "无公告",
					Content: "",
				}
			}
		}
		if result.Version != nil {
			m.version = *result.Version
		}
		zaplogger.GetSugar().Infof("load notifications ok")
	} else {
		zaplogger.GetSugar().Errorf("load notifications %s", errcode.GetErrorDesc(result.ErrCode))
		return errors.New(errcode.GetErrorDesc(result.ErrCode))
	}
	return nil
}

func (m *module) Init(app *gin.Engine) error {
	group := app.Group("/notification")
	group.POST("/update", webservice.WarpHandle(m.update))
	group.POST("/getAll", webservice.WarpHandle(m.getAll))
	group.POST("/get", webservice.WarpHandle(m.get))

	m.Notifications = make([]*Notification, 2)
	m.updateClosed(serverstatus.DefaultStatue)
	serverstatus.RegisterUpdateNotify(db.GetFlyfishClient("global"), func(status serverstatus.Status) {
		m.updateClosed(status)
	})

	return m.load()
}

func (m *module) Tick() {

}

var _module *module

func init() {
	_module = new(module)
	webservice.RegisterModule("notification", _module)
}
