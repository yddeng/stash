package offlinemsg

import (
	"fmt"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"github.com/sniperHW/flyfish/pkg/buffer"
	"sort"
	"time"
)

var ItemCap = 100

const Table = "offlinemsg"

type Item struct {
	Idx     int
	Version int64
	Content []byte
}

type DB interface {
	Load(table string, unikey string, version int64, cb func(ret *flyfish.GetResult))
	Update(table string, unikey string, fields map[string]interface{}, version int64, cb func(ret *flyfish.StatusResult))
}

type flyfishDB struct {
	cli *flyfish.Client
}

func NewFlyfishDB(cli *flyfish.Client) *flyfishDB {
	return &flyfishDB{
		cli: cli,
	}
}

func (db *flyfishDB) Load(table string, unikey string, version int64, cb func(ret *flyfish.GetResult)) {
	db.cli.GetAllWithVersion(table, unikey, version).AsyncExec(cb)
}

func (db *flyfishDB) Update(table string, unikey string, fields map[string]interface{}, version int64, cb func(ret *flyfish.StatusResult)) {
	//zaplogger.GetSugar().Debug(table, unikey, version)
	if version == 0 {
		db.cli.SetNx(table, unikey, fields).AsyncExec(func(result *flyfish.ValueResult) {
			//zaplogger.GetSugar().Debug(result)
			cb(&flyfish.StatusResult{
				ErrCode: result.ErrCode,
				Table:   result.Table,
				Key:     result.Key,
			})
		})
	} else {
		db.cli.Set(table, unikey, fields, version).AsyncExec(cb)
	}
}

func (it *Item) pack() []byte {
	bu := make([]byte, 0, 10+len(it.Content))
	bu = buffer.AppendUint16(bu, uint16(it.Idx))
	bu = buffer.AppendInt64(bu, it.Version)
	bu = buffer.AppendBytes(bu, it.Content)
	return bu
}

func (it *Item) unpack(bu []byte) *Item {
	r := buffer.NewReader(bu)
	it.Idx = int(r.GetUint16())
	it.Version = r.GetInt64()
	it.Content = bu[10:]
	return it
}

func load(db DB, unikey string, version int64, cb func(errcode.Error, *int64, []*Item)) {
	db.Load(Table, unikey, version, func(ret *flyfish.GetResult) {
		var items []*Item
		if ret.ErrCode == nil {
			items = make([]*Item, len(ret.Fields))
			c := 0
			for _, v := range ret.Fields {
				data := v.GetBlob()
				if len(data) > 0 {
					it := &Item{}
					it.unpack(data)
					items[it.Idx] = it
					c++
				}
			}
			items = items[:c]
		}
		cb(ret.ErrCode, ret.Version, items)
	})
}

type push struct {
	db          DB
	unikey      string
	content     []byte
	cb          func(errcode.Error, int64)
	deadline    time.Time
	itemVersion int64
}

func (p *push) updateCb(err errcode.Error, version int64, fields map[string]interface{}) {
	if nil == err {
		p.cb(nil, p.itemVersion)
	} else {
		switch errcode.GetCode(err) {
		case errcode.Errcode_version_mismatch, errcode.Errcode_record_exist:
			if errcode.GetCode(err) == errcode.Errcode_version_mismatch {
				if time.Now().Before(p.deadline) {
					time.AfterFunc(time.Millisecond*10, func() {
						load(p.db, p.unikey, 0, p.loadCb)
					})
				} else {
					p.cb(errcode.New(errcode.Errcode_timeout), 0)
				}
			} else {
				load(p.db, p.unikey, 0, p.loadCb)
			}
		case errcode.Errcode_retry:
			if time.Now().Before(p.deadline) {
				time.AfterFunc(time.Millisecond*50, func() {
					p.db.Update(Table, p.unikey, fields, version, func(ret *flyfish.StatusResult) {
						p.updateCb(ret.ErrCode, version, fields)
					})
				})
			} else {
				p.cb(errcode.New(errcode.Errcode_timeout), 0)
			}
		default:
			p.cb(err, 0)
		}
	}
}

func (p *push) loadCb(err errcode.Error, pversion *int64, items []*Item) {

	var version int64

	if nil != pversion {
		version = *pversion
	}

	if nil == err || errcode.GetCode(err) == errcode.Errcode_record_notexist {
		last := -1
		for i, v := range items {
			if v.Version == version {
				last = i
				break
			}
		}
		tail := (last + 1) % ItemCap

		var it *Item
		if tail >= len(items) {
			it = &Item{
				Idx: tail,
			}
		} else {
			it = items[tail]
		}

		it.Version = version + 1
		it.Content = p.content
		fields := map[string]interface{}{}
		fields[fmt.Sprintf("item_%d", it.Idx)] = it.pack()
		p.itemVersion = it.Version
		p.db.Update(Table, p.unikey, fields, version, func(ret *flyfish.StatusResult) {
			p.updateCb(ret.ErrCode, version, fields)
		})

	} else if errcode.GetCode(err) == errcode.Errcode_retry && time.Now().Before(p.deadline) {
		time.AfterFunc(time.Millisecond*50, func() {
			load(p.db, p.unikey, 0, p.loadCb)
		})
	}
}

type pull struct {
	db       DB
	unikey   string
	deadline time.Time
	version  int64
	cb       func(err errcode.Error, version int64, items []*Item)
}

func (p *pull) loadCb(err errcode.Error, pversion *int64, items []*Item) {
	var version int64

	if nil != pversion {
		version = *pversion
	}
	if nil == err || errcode.GetCode(err) == errcode.Errcode_record_unchange || errcode.GetCode(err) == errcode.Errcode_record_notexist {
		outItems := []*Item{}
		for _, v := range items {
			if v.Version > p.version {
				outItems = append(outItems, v)
			}
		}

		if len(outItems) > 1 {
			sort.Slice(outItems, func(i, j int) bool {
				return outItems[i].Version < outItems[j].Version
			})
		}

		p.cb(nil, version, outItems)
	} else if errcode.GetCode(err) == errcode.Errcode_retry {
		if time.Now().Before(p.deadline) {
			time.AfterFunc(time.Millisecond*50, func() {
				load(p.db, p.unikey, p.version, p.loadCb)
			})
		} else {
			p.cb(errcode.New(errcode.Errcode_timeout), 0, nil)
		}
	} else {
		p.cb(err, 0, nil)
	}
}

func PushMsg(db DB, userID string, topic string, content []byte, timeout time.Duration, cb func(errcode.Error, int64)) {
	p := &push{
		db:       db,
		unikey:   fmt.Sprintf("%s:%s", userID, topic),
		content:  content,
		cb:       cb,
		deadline: time.Now().Add(timeout),
	}

	load(p.db, p.unikey, 0, p.loadCb)
}

func PullMsg(db DB, userID string, topic string, version int64, timeout time.Duration, cb func(err errcode.Error, version int64, items []*Item)) {
	p := &pull{
		db:       db,
		unikey:   fmt.Sprintf("%s:%s", userID, topic),
		cb:       cb,
		deadline: time.Now().Add(timeout),
		version:  version,
	}

	load(p.db, p.unikey, version, p.loadCb)
}
