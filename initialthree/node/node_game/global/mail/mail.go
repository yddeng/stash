package mail

import (
	"errors"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/node/common/offlinemsg"
	"initialthree/zaplogger"
	"time"
)

const (
	userID = "global"
	topic  = "mail"
)

type GlobalMail struct {
	dbClient offlinemsg.DB
	version  int64
	items    []*offlinemsg.Item
	loading  bool
	loadFunc []func(items []*offlinemsg.Item, nowVersion int64)
}

func (this *GlobalMail) pull(callback func(err errcode.Error)) {
	this.loading = true
	offlinemsg.PullMsg(this.dbClient, userID, topic, 0, time.Second*6, func(err errcode.Error, version int64, items []*offlinemsg.Item) {
		if err == nil {
			this.version = version
			this.items = items
		} else {
			zaplogger.GetSugar().Errorf("load global mail %s", errcode.GetErrorDesc(err))
			if errcode.GetCode(err) == errcode.Errcode_record_notexist {
				err = nil
			}
		}
		this.loading = false

		if len(this.loadFunc) > 0 {
			callbacks := this.loadFunc
			this.loadFunc = this.loadFunc[0:0]
			for _, callback := range callbacks {
				callback(items, version)
			}
		}
		callback(err)
	})
}

func GetGlobalMail(version int64, callback func(items []*offlinemsg.Item, nowVersion int64)) {
	if globalMail.loading {
		globalMail.loadFunc = append(globalMail.loadFunc, callback)
	} else {
		items, version := globalMail.items, globalMail.version
		cluster.PostTask(func() {
			callback(items, version)
		})
	}
}

func OnUpdate() {
	globalMail.pull(func(err errcode.Error) {
		zaplogger.GetSugar().Errorf("load global mail %s", errcode.GetErrorDesc(err))
	})
}

var globalMail *GlobalMail

func InitGlobalMail(flyClient *client.Client) error {
	globalMail = &GlobalMail{
		dbClient: offlinemsg.NewFlyfishDB(flyClient),
	}

	errCh := make(chan error)
	globalMail.pull(func(err errcode.Error) {
		zaplogger.GetSugar().Errorf("load global mail %s", errcode.GetErrorDesc(err))
		if err == nil {
			errCh <- nil
		} else {
			errCh <- errors.New(errcode.GetErrorDesc(err))
		}
	})
	err := <-errCh
	return err
}
