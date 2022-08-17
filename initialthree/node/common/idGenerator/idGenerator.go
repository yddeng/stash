package idGenerator

import (
	"fmt"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
)

const (
	dbTable = "id_counter"
	dbField = "counter"
)

type IDGenerator struct {
	key     string
	flyCli  *client.Client
	genFunc func(int64) int64

	rsvModule *rsvModule
}

const (
	rsvStep        = 200 // 预留的宽度
	rsvIncrTrigger = 0.8 // 当使用比例达到该值，申请下一轮的预留值
)

type rsvModule struct {
	reserving   bool
	currentV    int64
	reservedV   int64
	incrTrigger int64
	callbacks   []func(int64, error)
}

func NewIDGenerator(key string, cli *client.Client, reserved bool, genFunc func(int64) int64) *IDGenerator {
	if cli == nil {
		panic("idGenerator:NewIDGenerator flyfish client is nil")
	}

	this := &IDGenerator{
		key:     key,
		flyCli:  cli,
		genFunc: genFunc,
	}

	if reserved {
		this.rsvModule = &rsvModule{
			callbacks: make([]func(int64, error), 0, rsvStep),
		}
		this.reserve()
	}

	return this
}

func (this IDGenerator) reserve() {
	this.rsvModule.reserving = true
	this.incrBy(rsvStep, func(i int64, e error) {
		this.rsvModule.reserving = false
		if e != nil {
			for _, callback := range this.rsvModule.callbacks {
				callback(0, e)
			}
			this.rsvModule.callbacks = this.rsvModule.callbacks[0:0]
		} else {
			this.rsvModule.reservedV = i
			this.rsvModule.currentV = this.rsvModule.reservedV - rsvStep
			this.rsvModule.incrTrigger = rsvStep*rsvIncrTrigger + this.rsvModule.currentV
			if this.rsvModule.incrTrigger > this.rsvModule.reservedV {
				this.rsvModule.incrTrigger = this.rsvModule.reservedV
			} else if this.rsvModule.incrTrigger < this.rsvModule.currentV {
				this.rsvModule.incrTrigger = this.rsvModule.currentV
			}
			//fmt.Println(this.rsvModule.currentV, this.rsvModule.incrTrigger, this.rsvModule.reservedV)

			callbacks := this.rsvModule.callbacks
			this.rsvModule.callbacks = make([]func(int64, error), 0, rsvStep)
			for _, callback := range callbacks {
				this.GenID(callback)
			}
		}
	})
}

func (this *IDGenerator) incrBy(value int64, callback func(int64, error)) {
	this.flyCli.IncrBy(dbTable, this.key, dbField, value).AsyncExec(func(result *client.ValueResult) {
		if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
			callback(result.Value.GetInt(), nil)
		} else {
			callback(0, fmt.Errorf("id generator gen failed %s", errcode.GetErrorDesc(result.ErrCode)))
		}
	})
}

func (this *IDGenerator) GenID(callback func(int64, error)) {
	if this.rsvModule != nil {
		if this.rsvModule.currentV < this.rsvModule.reservedV {
			this.rsvModule.currentV++
			if this.genFunc != nil {
				callback(this.genFunc(this.rsvModule.currentV), nil)
			} else {
				callback(this.rsvModule.currentV, nil)
			}

			if this.rsvModule.currentV >= this.rsvModule.incrTrigger && !this.rsvModule.reserving {
				this.reserve()
			}
		} else {
			this.rsvModule.callbacks = append(this.rsvModule.callbacks, callback)
		}
	} else {
		this.incrBy(1, func(i int64, e error) {
			if e != nil {
				callback(0, e)
			} else {
				if this.genFunc != nil {
					callback(this.genFunc(i), nil)
				} else {
					callback(i, nil)
				}
			}
		})
	}

}

var idGens = map[string]*IDGenerator{}

func Register(key string, cli *client.Client, genFunc func(int64) int64, reserved ...bool) {
	_reserved := false
	if len(reserved) != 0 && reserved[0] {
		_reserved = true
	}
	if _, ok := idGens[key]; !ok {
		idGen := NewIDGenerator(key, cli, _reserved, genFunc)
		idGens[key] = idGen
	}
}

func GetIDGen(key string) *IDGenerator {
	return idGens[key]
}
