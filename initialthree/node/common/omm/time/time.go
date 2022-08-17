package time

import (
	"encoding/json"
	"errors"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/pkg/timer"
	"time"
)

const (
	dbTable = "global_data"
	dbKey   = "time"
	dbField = "data"
)

type Data struct {
	Offset time.Duration `json:"offset,omitempty"`
}

func Unmarshal(fields map[string]*flyfish.Field) (time.Duration, error) {
	var out Data
	if field, ok := fields[dbField]; ok && len(field.GetBlob()) > 0 {
		if err := json.Unmarshal(field.GetBlob(), &out); err == nil {
			return out.Offset, nil
		} else {
			return 0, err
		}
	}
	return 0, errors.New("not exist")
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

func SetCmd(fc *flyfish.Client, offset time.Duration, version int64) *flyfish.StatusCmd {
	if fc == nil {
		panic("fc(*flyfish.Client) nil")
	}

	in := Data{Offset: offset}
	bytes, _ := json.Marshal(in)
	fields := map[string]interface{}{dbField: bytes}

	return fc.Set(dbTable, dbKey, fields, version)
}

const (
	reloadInterval = time.Second
)

type updateNotify struct {
	fc       *flyfish.Client
	callback func(offset time.Duration)

	version int64
	loading bool
}

func RegisterUpdateNotify(fc *flyfish.Client, callback func(offset time.Duration)) {
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
			if offset, err := Unmarshal(result.Fields); err == nil {
				notify.version = version
				notify.callback(offset)
			}
		}
	})
}
