package bannedlist

import (
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"time"
)

const (
	dbTableName        = "bannedlist"
	dbFieldKey         = "__key__"
	dbFieldExpiredTime = "expired_time"
)

func SetCmd(fc *flyfish.Client, userID string, expiredTime int64) *flyfish.StatusCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}
	if userID == "" {
		panic("userID is nil")
	}

	fields := map[string]interface{}{dbFieldExpiredTime: expiredTime}
	return fc.Set(dbTableName, userID, fields)
}

func DelCmd(fc *flyfish.Client, userID string) *flyfish.StatusCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}
	if userID == "" {
		panic("userID is nil")
	}

	return fc.Del(dbTableName, userID)
}

func GetCmd(fc *flyfish.Client, userID string) *flyfish.GetCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}
	if userID == "" {
		panic("userID is nil")
	}

	return fc.Get(dbTableName, userID, dbFieldExpiredTime)
}

type UserIDInfo struct {
	UserID      string
	ExpiredTime int64
}

func Unmarshal(userID string, fields map[string]*flyfish.Field) *UserIDInfo {
	var out UserIDInfo
	out.UserID = userID
	out.ExpiredTime = fields[dbFieldExpiredTime].GetInt()
	return &out
}

func AuthUserID(fc *flyfish.Client, userID string, cb func(banned bool)) {
	if fc == nil {
		panic("fc(*flyfish.client) nil")
	}
	if userID == "" {
		panic("userID is nil")
	}

	GetCmd(fc, userID).AsyncExec(func(r *flyfish.GetResult) {
		if r.ErrCode == nil {
			userInfo := Unmarshal(userID, r.Fields)

			if userInfo.ExpiredTime == 0 {
				// 永久封禁
				cb(true)
			} else {
				nowUnix := time.Now().Unix()
				if nowUnix <= userInfo.ExpiredTime {
					cb(true)
				} else {
					cb(false)
				}
			}

		} else if r.ErrCode.Code == errcode.Errcode_record_notexist {
			cb(false)
		} else {
			cb(false)
		}
	})

}
