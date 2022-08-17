package temporary

import (
	"fmt"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/db"
	"initialthree/zaplogger"

	"time"
)

type TemporaryI interface {
	UserDisconnect()
	UserLogout()
	Tick(now time.Time)
}

type UserI interface {
	GetTemporary(tt TemporaryType) TemporaryI
	SetTemporary(tt TemporaryType, cache TemporaryI)
	ClearTemporary(tt TemporaryType)
	GetID() uint64
	GetUserID() string
}

type TemporaryType int

const (
	Invaild        = TemporaryType(0)
	TempMapInfo    = TemporaryType(1)
	TempSeqCache   = TemporaryType(2)
	TempTeamInfo   = TemporaryType(3)
	TempLevelFight = TemporaryType(4)
)

/*
 * game 停机更新，重启后需要还原的临时数据
 */
type TempProcess interface {
	Marshal(temp TemporaryI) ([]byte, error)
	Unmarshal(user UserI, data []byte) (TemporaryI, error)
}

var (
	//receiver = codecs.NewReceiver("cs")
	encoder = codecs.NewEncoder("sc")

	tempData    = map[TemporaryType]*tempProcess{}
	tableName   = "temporary"
	tableFields []string
)

type tempProcess struct {
	tableField string
	process    TempProcess
}

func registerTempDataProcess(tt TemporaryType, field string, tp TempProcess) {
	if _, ok := tempData[tt]; ok {
		panic(fmt.Sprintf("temp data process %d is already register", tt))
	} else {
		tempData[tt] = &tempProcess{
			tableField: field,
			process:    tp,
		}
		tableFields = append(tableFields, field)
	}
}

//
func Load(user UserI) {
	key := fmt.Sprintf("%d", user.GetID())
	zaplogger.GetSugar().Infof("user %s %d load temporary %v", user.GetUserID(), user.GetID(), tableFields)
	ret := db.GetFlyfishClient("game").Get(tableName, key, tableFields...).Exec()

	if nil == ret.ErrCode {
		for tt, tp := range tempData {
			field, ok := ret.Fields[tp.tableField]
			if ok && len(field.GetBlob()) != 0 {
				if temp, err := tp.process.Unmarshal(user, field.GetBlob()); err == nil {
					zaplogger.GetSugar().Debugf("%s load temporary %s ", user.GetUserID(), tp.tableField)
					user.SetTemporary(tt, temp)
				}
			}
		}
	}

}

func Save(user UserI, fn func()) {
	fields := map[string]interface{}{}
	for tt, tp := range tempData {
		temp := user.GetTemporary(tt)
		if temp != nil {
			if data, err := tp.process.Marshal(temp); err == nil {
				//log.GetLogger().Debugln("get", tt, temp)
				fields[tp.tableField] = data
			} else {
				zaplogger.GetSugar().Errorf("save temporary %d err %s", tt, err.Error())
			}
		}
	}

	if len(fields) > 0 {
		key := fmt.Sprintf("%d", user.GetID())
		set := db.GetFlyfishClient("game").Set(tableName, key, fields)
		set.AsyncExec(func(ret *flyfish.StatusResult) {
			zaplogger.GetSugar().Debugf("user %s %d save temporary %s", user.GetUserID(), user.GetID(), errcode.GetErrorDesc(ret.ErrCode))
			fn()
		})
	} else {
		fn()
	}
}
