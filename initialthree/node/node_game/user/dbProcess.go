package user

import (
	"fmt"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/cluster/priority"
	"initialthree/node/common/db"
	"initialthree/node/node_game/module"
	"initialthree/pkg/timer"
	"initialthree/zaplogger"
	"time"
)

type readOutClusterEle struct {
	uniKey  string
	table   string
	key     string
	fields  []string
	modules []module.ModuleI
	result  *flyfish.GetResult
}

type readOutCluster struct {
	readOuts map[string]*readOutClusterEle
}

func (roc *readOutCluster) addCommand(cmd *module.ReadOutCommand) {
	if len(cmd.Fields) == 0 {
		zaplogger.GetSugar().Errorf("readout %s in table %s without fields", cmd.Key, cmd.Table)
		return
	}

	uniKey := cmd.Table + ":" + cmd.Key
	readout := roc.readOuts[uniKey]
	if readout == nil {
		readout = &readOutClusterEle{
			uniKey:  uniKey,
			table:   cmd.Table,
			key:     cmd.Key,
			fields:  cmd.Fields,
			modules: []module.ModuleI{cmd.Module},
		}
		roc.readOuts[uniKey] = readout
	} else {
		readout.fields = append(readout.fields, cmd.Fields...)
		readout.modules = append(readout.modules, cmd.Module)
	}
}

func (roc *readOutCluster) do(user *User, done func()) {
	if len(roc.readOuts) == 0 {
		done()
		return
	}

	gets := make([]*flyfish.GetCmd, 0, len(roc.readOuts))
	for _, v := range roc.readOuts {
		var cmdGet *flyfish.GetCmd
		if v.fields[0] == "all" {
			cmdGet = db.GetFlyfishClient("game").GetAll(v.table, v.key)
		} else {
			cmdGet = db.GetFlyfishClient("game").Get(v.table, v.key, v.fields...)
		}

		gets = append(gets, cmdGet)
	}

	flyfish.MGetWithEventQueue(priority.LOW, cluster.GetEventQueue(), gets...).AsyncExec(func(results []*flyfish.GetResult) {
		for _, result := range results {
			uniKey := result.Table + ":" + result.Key

			readOut := roc.readOuts[uniKey]
			if readOut == nil {
				panic(fmt.Errorf("readout %s not found", uniKey))
			}

			readOut.result = result
		}

		done()
	})
}

type writeBackClusterEle struct {
	table   string
	key     string
	unikey  string
	fields  []*module.WriteBackFiled
	modules map[module.ModuleI]struct{}
	result  *flyfish.StatusResult
}

type writeBackCluster struct {
	writebacks map[string]*writeBackClusterEle
}

func (this *writeBackCluster) addCommand(c *module.WriteBackCommand) {
	unikey := c.Table + ":" + c.Key
	ele, ok := this.writebacks[unikey]
	if !ok {
		ele = &writeBackClusterEle{
			table:   c.Table,
			key:     c.Key,
			unikey:  unikey,
			modules: map[module.ModuleI]struct{}{c.Module: {}},
			fields:  []*module.WriteBackFiled{},
		}
		this.writebacks[unikey] = ele
	}

	if _, ok := ele.modules[c.Module]; !ok {
		ele.modules[c.Module] = struct{}{}
	}
	ele.fields = append(ele.fields, c.Fields...)
}

func (this *writeBackCluster) do(user *User, callback func(bool)) {
	fn := func() {
		allOK := true
		for _, v := range this.writebacks {
			if v.result.ErrCode != nil {
				for vv := range v.modules {
					vv.WriteBackRet(false)
				}
				allOK = false
				zaplogger.GetSugar().Errorf("dbSave %s Error:%s %s", user.GetUserLogName(), errcode.GetErrorDesc(v.result.ErrCode), v.result.Table)
			} else {
				for vv := range v.modules {
					vv.WriteBackRet(true)
				}
			}
		}

		callback(allOK)
	}

	counter := 0
	for _, v := range this.writebacks {
		ele := v
		fields := map[string]interface{}{}
		for _, vv := range ele.fields {
			fields[vv.Name] = vv.Value
		}
		cmd := db.GetFlyfishClient("game").Set(ele.table, ele.key, fields)
		cmd.AsyncExec(func(ret *flyfish.StatusResult) {
			ele.result = ret
			counter++
			if counter == len(this.writebacks) {
				fn()
			}
		})
	}
}

type saveProcessor struct {
	u             *User
	lock          bool
	pending       []bool
	finalCallback func()
}

func (this *saveProcessor) doSave(finalSave bool, callback func()) {
	if this.lock {
		this.pending = append(this.pending, finalSave)
		if finalSave {
			zaplogger.GetSugar().Debugf("%s finalSave ,but lock ", this.u.GetUserLogName())
			this.finalCallback = callback
		}
		return
	}

	if len(this.pending) > 0 {
		for _, v := range this.pending {
			if v {
				finalSave = true
				callback = this.finalCallback
				break
			}
		}
		this.pending = this.pending[:0]
	}

	this.lock = true
	this.do(finalSave, callback, 1)
}

func (this *saveProcessor) do(finalSave bool, callback func(), saveCount int) {
	wbCluster := &writeBackCluster{
		writebacks: map[string]*writeBackClusterEle{},
	}

	for _, v := range this.u.modules {
		if v.IsDirty() {
			writeBack := v.WriteBack()
			if writeBack != nil {
				wbCluster.addCommand(writeBack)
			}
		}
	}

	if len(wbCluster.writebacks) == 0 {
		this.lock = false
		if finalSave {
			callback()
		}
		return
	}

	//zaplogger.GetSugar().Debugf("%s saveProcessor do,len %d ", this.u.GetUserLogName(), len(wbCluster.writebacks))
	wbCluster.do(this.u, func(ok bool) {
		if finalSave {
			if ok {
				callback()
			} else {
				saveCount++
				if saveCount >= 4 {
					callback()
				} else {
					//try again
					cluster.RegisterTimerOnce(time.Second*(time.Duration(saveCount-1)), func(_ *timer.Timer, _ interface{}) {
						zaplogger.GetSugar().Errorf("%s %s %d", this.u.GetUserLogName(), "finalSave try again", saveCount)
						this.do(finalSave, callback, saveCount)
					}, nil)
				}
			}
		} else {
			this.lock = false
			if len(this.pending) > 0 {
				this.doSave(false, nil)
			}
		}
	})
}

// 定时存储
func (this *User) durationSave() {
	this.saveP.doSave(false, nil)
}

// 下线存储
func (this *User) finalSave(callback func()) {
	this.saveP.doSave(true, callback)
}
