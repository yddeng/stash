package testdb

import (
	"fmt"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	flyproto "github.com/sniperHW/flyfish/proto"
	"sync"
)

//内存DB

type testRecord struct {
	version int64
	fields  map[string]interface{}
}

type testTable struct {
	sync.Mutex
	records map[string]*testRecord
}

func (this *testTable) get(owner string) *testRecord {
	this.Lock()
	defer this.Unlock()
	r, ok := this.records[owner]
	if ok {
		rr := &testRecord{
			version: r.version,
			fields:  map[string]interface{}{},
		}
		for k, v := range r.fields {
			rr.fields[k] = v
		}
		return rr
	} else {
		return nil
	}
}

func (this *testTable) set(owner string, version int64, fields map[string]interface{}) (errcode.Error, int64) {
	this.Lock()
	defer this.Unlock()
	r, ok := this.records[owner]
	if ok {
		if r.version != version {
			return errcode.New(errcode.Errcode_version_mismatch), 0
		} else {
			r.version++
			for k, v := range fields {
				r.fields[k] = v
			}
			return nil, r.version
		}
	} else {
		this.records[owner] = &testRecord{
			version: 1,
			fields:  fields,
		}
		return nil, 1
	}
}

type testDB struct {
	tables map[string]*testTable
	err    errcode.Error
}

func (this *testDB) SetErr(err errcode.Error) {
	this.err = err
}

func (this *testDB) get(table string, owner string) *testRecord {
	//fmt.Println("get", table, owner)
	t, ok := this.tables[table]
	if ok {
		return t.get(owner)
	} else {
		return nil
	}
}

func (this *testDB) set(table string, owner string, version int64, fields map[string]interface{}) (errcode.Error, int64) {
	//fmt.Println("set", table, owner)
	t, ok := this.tables[table]
	if ok {
		return t.set(owner, version, fields)
	} else {
		panic("invaild table")
		return nil, 0
	}
}

func (t *testDB) Load(table string, owner string, version int64, cb func(ret *flyfish.SliceResult)) {

	//fmt.Println("Load", table, owner)

	result := &flyfish.SliceResult{
		Table: table,
		Key:   owner,
	}

	if t.err != nil {
		result.ErrCode = t.err
		cb(result)
		return
	}

	r := t.get(table, owner)

	if nil == r {
		result.ErrCode = errcode.New(errcode.Errcode_record_notexist)
	} else if r.version == version {
		fmt.Println("Errcode_record_unchange", version, r.version)
		result.ErrCode = errcode.New(errcode.Errcode_record_unchange)
		result.Version = r.version
	} else {
		result.Version = r.version
		result.Fields = map[string]*flyfish.Field{}
		for k, v := range r.fields {
			result.Fields[k] = (*flyfish.Field)(flyproto.PackField(k, v))
		}
	}

	cb(result)

}

func (t *testDB) Update(table string, owner string, fields map[string]interface{}, version int64, cb func(ret *flyfish.StatusResult)) {

	if t.err != nil {
		cb(&flyfish.StatusResult{
			ErrCode: t.err,
			Table:   table,
			Key:     owner,
		})
		return
	}

	errCode, v := t.set(table, owner, version, fields)
	result := &flyfish.StatusResult{
		ErrCode: errCode,
		Table:   table,
		Key:     owner,
		Version: v,
	}

	if errCode != nil && errcode.GetCode(errCode) != errcode.Errcode_record_notexist {
		fmt.Println("db Update", errCode)
	}

	cb(result)
}

func NewTestDB() *testDB {
	tdb := &testDB{
		tables: map[string]*testTable{},
	}

	tdb.tables["offlinemsg"] = &testTable{
		records: map[string]*testRecord{},
	}

	return tdb
}
