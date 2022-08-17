package serverstatus

import (
	"errors"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/pkg/json"
	"initialthree/pkg/timer"
	"time"
)

type Status int8

const (
	StatusOpen = Status(iota)
	StatusMaintenance
)

var statusStrings = [...]string{
	StatusOpen:        "Open",
	StatusMaintenance: "Maintenance",
}

var DefaultStatue = StatusOpen

func (s Status) String() string {
	if s.Valid() {
		return statusStrings[s]
	} else {
		return "valid"
	}
}

func (s Status) Valid() bool {
	return s >= StatusOpen && s <= StatusMaintenance
}

type ServerStatus struct {
	Status `json:"status,omitempty"`
}

const (
	dbTable = "global_data"
	dbKey   = "serverstatus"
	dbField = "data"
)

func Unmarshal(fields map[string]*flyfish.Field) (Status, error) {
	var out ServerStatus
	if field, ok := fields[dbField]; ok && len(field.GetBlob()) > 0 {
		if err := json.Unmarshal(field.GetBlob(), &out); err == nil {
			return out.Status, nil
		} else {
			return DefaultStatue, err
		}
	}
	return DefaultStatue, errors.New("not exist")
}

func LoadCmd(fc *flyfish.Client, version ...int64) *flyfish.GetCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}

	var cmd *flyfish.GetCmd
	if len(version) > 0 {
		cmd = fc.GetWithVersion(dbTable, dbKey, version[0], dbField)
	} else {
		cmd = fc.Get(dbTable, dbKey, dbField)
	}
	return cmd
}

func SetCmd(fc *flyfish.Client, status Status, version int64) *flyfish.StatusCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}

	in := ServerStatus{Status: status}
	bytes, _ := json.Marshal(in)
	fields := map[string]interface{}{dbField: bytes}

	return fc.Set(dbTable, dbKey, fields, version)
}

const (
	reloadInterval = time.Second
)

type updateNotify struct {
	fc       *flyfish.Client
	callback func(status Status)

	version int64
	loading bool
}

func RegisterUpdateNotify(fc *flyfish.Client, callback func(status Status)) {
	l := &updateNotify{fc: fc, callback: callback}
	cluster.RegisterTimer(reloadInterval, func(timer *timer.Timer, i interface{}) {
		ll := i.(*updateNotify)
		reload(ll)
	}, l)
}

func reload(notify *updateNotify) {
	if notify.loading {
		return
	}

	notify.loading = true

	LoadCmd(notify.fc, notify.version).AsyncExec(func(result *flyfish.GetResult) {
		notify.loading = false
		if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
			var version int64
			if result.Version != nil {
				version = *result.Version
			}
			if status, err := Unmarshal(result.Fields); err == nil {
				notify.version = version
				notify.callback(status)
			}
		}
	})
}
