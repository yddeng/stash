package whitelist

import (
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
)

const (
	dbTable            = "whitelist"
	dbFieldKey         = "__key__"
	dbFieldAccessTimes = "access_times"
)

func SetCmd(fc *flyfish.Client, userID string) *flyfish.StatusCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}
	if userID == "" {
		panic("userID is nil")
	}

	fields := map[string]interface{}{dbFieldAccessTimes: 0}
	return fc.Set(dbTable, userID, fields)
}

func DelCmd(fc *flyfish.Client, userID string) *flyfish.StatusCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}
	if userID == "" {
		panic("userID is nil")
	}

	return fc.Del(dbTable, userID)
}

func GetCmd(fc *flyfish.Client, userID string) *flyfish.GetCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}
	if userID == "" {
		panic("userID is nil")
	}

	return fc.Get(dbTable, userID, dbFieldAccessTimes)
}

type UserIDInfo struct {
	UserID      string `json:"user_id"`
	AccessTimes int64  `json:"access_times"`
}

func Unmarshal(userID string, fields map[string]*flyfish.Field) *UserIDInfo {
	var out UserIDInfo
	out.UserID = userID
	out.AccessTimes = fields[dbFieldAccessTimes].GetInt()
	return &out
}

func AuthUserID(fc *flyfish.Client, userID string, cb func(exist bool)) {
	if fc == nil {
		panic("fc(*flyfish.client) nil")
	}
	if userID == "" {
		panic("userID is nil")
	}

	GetCmd(fc, userID).AsyncExec(func(r *flyfish.GetResult) {
		if r.ErrCode == nil {
			fc.IncrBy(dbTable, userID, dbFieldAccessTimes, 1).AsyncExec(func(result *flyfish.ValueResult) {})
			cb(true)
		} else if r.ErrCode.Code == errcode.Errcode_record_notexist {
			cb(false)
		} else {
			cb(false)
		}
	})

}
