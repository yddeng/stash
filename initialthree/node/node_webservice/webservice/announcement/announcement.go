package announcement

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/node/common/db"
	"initialthree/node/common/idGenerator"
	"initialthree/node/node_webservice/webservice"
	"initialthree/pkg/json"
	"initialthree/zaplogger"
	"sort"
	"time"
)

const (
	tableName = "web_data"
	tableKey  = "announcement"
	slotCount = 100
)

func slotName(slotIdx int) string {
	return fmt.Sprintf("slot%d", slotIdx)
}

type Announcement struct {
	ID         int64                 `json:"id"`
	Type       string                `json:"type"`
	Title      string                `json:"title"`
	SmallTitle string                `json:"smallTitle"`
	Content    []AnnouncementContent `json:"content"`
	StartTime  int64                 `json:"startTime"`
	ExpireTime int64                 `json:"expireTime"` // 常驻显示， 值为0
	Remind     bool                  `json:"remind"`
}

type AnnouncementContent struct {
	Type      string `json:"type"`
	Image     string `json:"image"`
	ImageSkip int    `json:"imageSkip"`
	Text      string `json:"text"`
}

type module struct {
	dataVersion int64
	data        []*Announcement
	result      []*Announcement
}

func (m *module) findSlotIdx(id int64) int {
	for slotIdx, anno := range m.data {
		if anno != nil && anno.ID == id {
			return slotIdx
		}
	}
	return 0
}

func (m *module) load() error {
	m.data = make([]*Announcement, slotCount)
	result := db.GetFlyfishClient("global").GetAllWithVersion(tableName, tableKey, m.dataVersion).Exec()
	if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
		for i := 0; i < slotCount; i++ {
			fieldName := slotName(i)
			field, ok := result.Fields[fieldName]
			if ok && len(field.GetBlob()) > 0 {
				err := json.Unmarshal(field.GetBlob(), &m.data[i])
				if err != nil {
					zaplogger.GetSugar().Error(err.Error(), string(field.GetBlob()))
					return err
				}
			}
		}
		if result.Version != nil {
			m.dataVersion = *result.Version
		}
		m.updateResult()
		zaplogger.GetSugar().Infof("load announcement ok")
	} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist || errcode.GetCode(result.ErrCode) == errcode.Errcode_record_unchange {
		zaplogger.GetSugar().Errorf("load announcement %s", errcode.GetErrorDesc(result.ErrCode))
	} else {
		zaplogger.GetSugar().Errorf("load announcement %s", errcode.GetErrorDesc(result.ErrCode))
		return errors.New(errcode.GetErrorDesc(result.ErrCode))
	}
	return nil
}

func (m *module) update(anno *Announcement, count int, getSlotIdx func() int, cb func(err error)) {
	slotIdx := getSlotIdx()
	bytes := []byte{}
	if anno != nil {
		bytes, _ = json.Marshal(anno)
	}
	cmd := db.GetFlyfishClient("global").Set(tableName, tableKey, map[string]interface{}{
		slotName(slotIdx): bytes,
	}, m.dataVersion)
	cmd.AsyncExec(func(result *client.StatusResult) {
		if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
			//m.dataVersion = 1
			//m.data[slotIdx] = anno
			//m.updateResult()
			//cb(nil)
			err := m.load()
			cb(err)
		} else if errcode.GetCode(result.ErrCode) == errcode.Errcode_version_mismatch {
			zaplogger.GetSugar().Infof("update announcement version %d miss match ,reload data", m.dataVersion)
			if err := m.load(); err != nil {
				cb(err)
				return
			}
			if count <= 0 {
				cb(errors.New("retry"))
			} else {
				m.update(anno, count-1, getSlotIdx, cb)
			}
		} else {
			zaplogger.GetSugar().Errorf("update announcement %d failed: %s", slotIdx, errcode.GetErrorDesc(result.ErrCode))
			cb(errors.New(errcode.GetErrorDesc(result.ErrCode)))
		}
	})
}

func (m *module) updateResult() {
	m.result = make([]*Announcement, 0, slotCount)
	nowUnix := time.Now().Unix()
	for _, v := range m.data {
		if v != nil && (v.ExpireTime == 0 || nowUnix < v.ExpireTime) {
			m.result = append(m.result, v)
		}
	}

	// 创建时间升序、已到期、空
	sort.Slice(m.result, func(i, j int) bool {
		annoi, annoj := m.result[i], m.result[j]
		if annoi != nil && annoj == nil {
			return true
		} else if annoi != nil && annoj != nil {
			if (annoi.ExpireTime == 0 || annoi.ExpireTime > nowUnix) && (annoj.ExpireTime != 0 && annoj.ExpireTime <= nowUnix) {
				// i 到期，j 未到期
				return true
			} else if (annoi.ExpireTime == 0 || annoi.ExpireTime > nowUnix) && (annoj.ExpireTime == 0 && annoj.ExpireTime > nowUnix) {
				// i,j 都未到期
				if annoi.StartTime < annoj.StartTime {
					// 升序
					return true
				}
			}
		}
		return false
	})
}

func (m *module) addAnnouncement(wait *webservice.WaitConn, req *Announcement) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	if req.Type == "" {
		wait.SetResult("缺少必要参数", nil)
		wait.Done()
		return
	}

	if req.StartTime == 0 {
		req.StartTime = time.Now().Unix()
	}

	idGenerator.GetIDGen(tableKey).GenID(func(i int64, e error) {
		if e != nil {
			wait.SetResult(e.Error(), nil)
			wait.Done()
		} else {
			req.ID = i

			m.update(req, 3, func() int {
				nowUnix := time.Now().Unix()
				slotIdx := 0
				for i, v := range m.data {
					if v == nil {
						// 当前位置为空
						return i
					} else {
						if v.ExpireTime != 0 && nowUnix > v.ExpireTime {
							// 已经过期
							return i
						}
						if v.StartTime < m.data[slotIdx].StartTime {
							// 记录最早创建的ID
							slotIdx = i
						}
					}
				}
				return slotIdx
			}, func(err error) {
				if err != nil {
					wait.SetResult(err.Error(), nil)
				}
				wait.Done()
			})
		}
	})
}

func (m *module) updateAnnouncement(wait *webservice.WaitConn, req *Announcement) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	if req.Type == "" {
		wait.SetResult("缺少必要参数", nil)
		wait.Done()
		return
	}

	if req.StartTime == 0 {
		req.StartTime = time.Now().Unix()
	}

	idGenerator.GetIDGen(tableKey).GenID(func(i int64, e error) {
		if e != nil {
			wait.SetResult(e.Error(), nil)
			wait.Done()
		} else {
			slotIdx := m.findSlotIdx(req.ID)
			req.ID = i
			m.update(req, 3, func() int {
				return slotIdx
			}, func(err error) {
				if err != nil {
					wait.SetResult(err.Error(), nil)
				}
				wait.Done()
			})
		}
	})

}

func (m *module) delAnnouncement(wait *webservice.WaitConn, req struct {
	ID int64 `json:"id"`
}) {
	zaplogger.GetSugar().Infof("%s %v", wait.GetRoute(), req)

	m.update(nil, 2, func() int {
		return m.findSlotIdx(req.ID)
	}, func(err error) {
		if err != nil {
			wait.SetResult(err.Error(), nil)
		}
		wait.Done()
	})

}

type getResult struct {
	Version      int64           `json:"version"`
	Announcement []*Announcement `json:"announcement"`
}

func (m *module) getAnnouncement(wait *webservice.WaitConn, req struct {
	Version int64 `json:"version"`
}) {

	ret := getResult{
		Version: m.dataVersion,
	}
	if req.Version != m.dataVersion {
		ret.Announcement = m.result
	}
	wait.SetResult("", ret)
	wait.Done()
}

func (m *module) Init(app *gin.Engine) error {
	group := app.Group("/announcement")
	group.POST("/add", webservice.WarpHandle(_module.addAnnouncement))
	group.POST("/update", webservice.WarpHandle(_module.updateAnnouncement))
	group.POST("/delete", webservice.WarpHandle(_module.delAnnouncement))
	group.POST("/get", webservice.WarpHandle(_module.getAnnouncement))

	idGenerator.Register(tableKey, db.GetFlyfishClient("global"), func(i int64) int64 {
		return i
	})

	return m.load()
}

func (m *module) Tick() {
	nowUnix := time.Now().Unix()
	for i, v := range m.result {
		if v.ExpireTime != 0 && nowUnix > v.ExpireTime {
			m.result = append(m.result[:i], m.result[i+1:]...)
			break
		}
	}
}

var _module *module

func init() {
	_module = new(module)
	webservice.RegisterModule("announcement", _module)
}
