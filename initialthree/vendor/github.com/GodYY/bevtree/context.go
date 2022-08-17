package bevtree

import (
	"reflect"
	"strings"
	"time"

	"github.com/GodYY/gutils/assert"
	"github.com/pkg/errors"
)

// Context interface.
type Context interface {
	// Behavior tree
	Tree() Tree

	// Get user data.
	UserData() interface{}

	// Get the update serial number.
	UpdateSeri() uint32

	// Get data-set.
	DataSet() DataSet

	// Get Framework.
	framework() *Framework

	// Release the Context.
	release()

	// Reset the Context.
	reset()

	// Update.
	update()

	cloneWithTree(tree Tree, independentDataSet bool) Context

	internal
}

type context struct {
	_framework   *Framework
	tree         Tree
	userData     interface{}
	updateSeri   uint32
	dataSet      *dataSet
	dataSetOwner bool

	internalImpl
}

func newContext(framework *Framework, tree Tree, userData interface{}) *context {
	assert.Assert(framework != nil, "framework nil")
	assert.Assert(tree != nil, "tree nil")

	ctx := &context{
		_framework:   framework,
		tree:         tree,
		userData:     userData,
		dataSet:      newDataSet(),
		dataSetOwner: true,
	}

	return ctx
}

func (ctx *context) framework() *Framework { return ctx._framework }

func (ctx *context) Tree() Tree { return ctx.tree }

func (ctx *context) UserData() interface{} { return ctx.userData }

func (ctx *context) DataSet() DataSet { return ctx.dataSet }

func (ctx *context) UpdateSeri() uint32 { return ctx.updateSeri }

func (ctx *context) release() {
	if ctx.dataSetOwner {
		ctx.dataSet.Clear()
	}
	ctx.dataSet = nil
	ctx.userData = nil
	ctx.tree = nil
}

func (ctx *context) reset() {
	ctx.updateSeri = 0
	if ctx.dataSetOwner {
		ctx.dataSet.Clear()
	}
}

func (ctx *context) update() { ctx.updateSeri++ }

func (ctx *context) cloneWithTree(tree Tree, independentDataSet bool) Context {
	assert.Assert(tree != nil, "tree nil")

	cp := &context{
		_framework: ctx._framework,
		tree:       tree,
		userData:   ctx.userData,
		updateSeri: ctx.updateSeri,
	}

	if independentDataSet {
		cp.dataSet = newDataSet()
		cp.dataSetOwner = true
	} else {
		cp.dataSet = ctx.dataSet
		cp.dataSetOwner = false
	}

	return cp
}

// Return a error indicates that the value type of key is not
// wanted type.
func ErrGetValueType(key string, want reflect.Type, get interface{}) error {
	return errors.Errorf("Get%s(%s): %s", strings.Title(want.Name()), key, reflect.TypeOf(get).Name())
}

// Return a error indicate that the value of key is not exist with op.
func ErrValueNotExist(key, op string) error {
	return errors.Errorf("%s(%s): value not exist", op, key)
}

// DataSet interface.
type DataSet interface {
	Set(string, interface{})
	Get(string) interface{}
	Remove(string) interface{}
	Clear()

	SetInt8(string, int8)
	GetInt8(string) (int8, bool)
	AddInt8(string, int8) int8
	IncInt8(string) int8
	DecInt8(string) int8

	SetUint8(string, uint8)
	GetUint8(string) (uint8, bool)
	AddUint8(string, uint8) uint8
	SubUint8(string, uint8) uint8
	IncUint8(string) uint8
	DecUint8(string) uint8

	SetInt16(string, int16)
	GetInt16(string) (int16, bool)
	AddInt16(string, int16) int16
	IncInt16(string) int16
	DecInt16(string) int16

	SetUint16(string, uint16)
	GetUint16(string) (uint16, bool)
	AddUint16(string, uint16) uint16
	SubUint16(string, uint16) uint16
	IncUint16(string) uint16
	DecUint16(string) uint16

	SetInt32(string, int32)
	GetInt32(string) (int32, bool)
	AddInt32(string, int32) int32
	IncInt32(string) int32
	DecInt32(string) int32

	SetUint32(string, uint32)
	GetUint32(string) (uint32, bool)
	AddUint32(string, uint32) uint32
	SubUint32(string, uint32) uint32
	IncUint32(string) uint32
	DecUint32(string) uint32

	SetInt(string, int)
	GetInt(string) (int, bool)
	AddInt(string, int) int
	IncInt(string) int
	DecInt(string) int

	SetUint(string, uint)
	GetUint(string) (uint, bool)
	AddUint(string, uint) uint
	SubUint(string, uint) uint
	IncUint(string) uint
	DecUint(string) uint

	SetInt64(string, int64)
	GetInt64(string) (int64, bool)
	AddInt64(string, int64) int64
	IncInt64(string) int64
	DecInt64(string) int64

	SetUint64(string, uint64)
	GetUint64(string) (uint64, bool)
	AddUint64(string, uint64) uint64
	SubUint64(string, uint64) uint64
	IncUint64(string) uint64
	DecUint64(string) uint64

	SetFloat32(string, float32)
	GetFloat32(string) (float32, bool)
	AddFloat32(string, float32) float32

	SetFloat64(string, float64)
	GetFloat64(string) (float64, bool)
	AddFloat64(string, float64) float64

	SetDuration(string, time.Duration)
	GetDuration(string) (time.Duration, bool)

	SetTime(string, time.Time)
	GetTime(string) (time.Time, bool)
}

// dataSet is used to store key-values, like blackboard.
type dataSet struct {
	keyValues map[string]interface{}
}

func newDataSet() *dataSet {
	return &dataSet{
		keyValues: make(map[string]interface{}),
	}
}

func (dc *dataSet) Set(key string, val interface{}) {
	dc.keyValues[key] = val
}

func (dc *dataSet) Get(key string) interface{} {
	return dc.keyValues[key]
}

func (dc *dataSet) Remove(key string) interface{} {
	if val := dc.keyValues[key]; val == nil {
		return nil
	} else {
		delete(dc.keyValues, key)
		return val
	}
}

func (dc *dataSet) Clear() {
	dc.keyValues = map[string]interface{}{}
}

func (dc *dataSet) SetInt8(key string, val int8) { dc.Set(key, val) }

func (dc *dataSet) GetInt8(key string) (int8, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(int8); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(int8(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddInt8(key string, d int8) int8 {
	if v, ok := dc.GetInt8(key); ok {
		v += d
		dc.SetInt8(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddInt8"))
	}
}

func (dc *dataSet) IncInt8(key string) int8 {
	return dc.AddInt8(key, 1)
}

func (dc *dataSet) DecInt8(key string) int8 {
	return dc.AddInt8(key, -1)
}

func (dc *dataSet) SetUint8(key string, val uint8) { dc.Set(key, val) }

func (dc *dataSet) GetUint8(key string) (uint8, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(uint8); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(uint8(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddUint8(key string, d uint8) uint8 {
	if v, ok := dc.GetUint8(key); ok {
		v += d
		dc.SetUint8(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddUint8"))
	}
}

func (dc *dataSet) SubUint8(key string, d uint8) uint8 {
	if v, ok := dc.GetUint8(key); ok {
		v -= d
		dc.SetUint8(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "SubUint8"))
	}
}

func (dc *dataSet) IncUint8(key string) uint8 {
	return dc.AddUint8(key, 1)
}

func (dc *dataSet) DecUint8(key string) uint8 {
	return dc.SubUint8(key, 1)
}

func (dc *dataSet) SetInt16(key string, val int16) { dc.Set(key, val) }

func (dc *dataSet) GetInt16(key string) (int16, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(int16); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(int16(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddInt16(key string, d int16) int16 {
	if v, ok := dc.GetInt16(key); ok {
		v += d
		dc.SetInt16(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddInt16"))
	}
}

func (dc *dataSet) IncInt16(key string) int16 {
	return dc.AddInt16(key, 1)
}

func (dc *dataSet) DecInt16(key string) int16 {
	return dc.AddInt16(key, -1)
}

func (dc *dataSet) SetUint16(key string, val uint16) { dc.Set(key, val) }

func (dc *dataSet) GetUint16(key string) (uint16, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(uint16); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(uint16(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddUint16(key string, d uint16) uint16 {
	if v, ok := dc.GetUint16(key); ok {
		v += d
		dc.SetUint16(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddUint16"))
	}
}

func (dc *dataSet) SubUint16(key string, d uint16) uint16 {
	if v, ok := dc.GetUint16(key); ok {
		v -= d
		dc.SetUint16(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "SubUint16"))
	}
}

func (dc *dataSet) IncUint16(key string) uint16 {
	return dc.AddUint16(key, 1)
}

func (dc *dataSet) DecUint16(key string) uint16 {
	return dc.SubUint16(key, 1)
}

func (dc *dataSet) SetInt32(key string, val int32) { dc.Set(key, val) }

func (dc *dataSet) GetInt32(key string) (int32, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(int32); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(int32(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddInt32(key string, d int32) int32 {
	if v, ok := dc.GetInt32(key); ok {
		v += d
		dc.SetInt32(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddInt32"))
	}
}

func (dc *dataSet) IncInt32(key string) int32 {
	return dc.AddInt32(key, 1)
}

func (dc *dataSet) DecInt32(key string) int32 {
	return dc.AddInt32(key, -1)
}

func (dc *dataSet) SetInt(key string, val int) { dc.Set(key, val) }

func (dc *dataSet) GetInt(key string) (int, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(int); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(int(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddInt(key string, d int) int {
	if v, ok := dc.GetInt(key); ok {
		v += d
		dc.SetInt(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddInt"))
	}
}

func (dc *dataSet) IncInt(key string) int {
	return dc.AddInt(key, 1)
}

func (dc *dataSet) DecInt(key string) int {
	return dc.AddInt(key, -1)
}

func (dc *dataSet) SetUint(key string, val uint) { dc.Set(key, val) }

func (dc *dataSet) GetUint(key string) (uint, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(uint); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(uint(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddUint(key string, d uint) uint {
	if v, ok := dc.GetUint(key); ok {
		v += d
		dc.SetUint(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddUint"))
	}
}

func (dc *dataSet) SubUint(key string, d uint) uint {
	if v, ok := dc.GetUint(key); ok {
		v -= d
		dc.SetUint(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "SubUint"))
	}
}

func (dc *dataSet) IncUint(key string) uint {
	return dc.AddUint(key, 1)
}

func (dc *dataSet) DecUint(key string) uint {
	return dc.SubUint(key, 1)
}

func (dc *dataSet) SetUint32(key string, val uint32) { dc.Set(key, val) }

func (dc *dataSet) GetUint32(key string) (uint32, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(uint32); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(uint32(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddUint32(key string, d uint32) uint32 {
	if v, ok := dc.GetUint32(key); ok {
		v += d
		dc.SetUint32(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddUint32"))
	}
}

func (dc *dataSet) SubUint32(key string, d uint32) uint32 {
	if v, ok := dc.GetUint32(key); ok {
		v -= d
		dc.SetUint32(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "SubUint32"))
	}
}

func (dc *dataSet) IncUint32(key string) uint32 {
	return dc.AddUint32(key, 1)
}

func (dc *dataSet) DecUint32(key string) uint32 {
	return dc.SubUint32(key, 1)
}

func (dc *dataSet) SetInt64(key string, val int64) { dc.Set(key, val) }

func (dc *dataSet) GetInt64(key string) (int64, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(int64); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(int64(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddInt64(key string, d int64) int64 {
	if v, ok := dc.GetInt64(key); ok {
		v += d
		dc.SetInt64(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddInt64"))
	}
}

func (dc *dataSet) IncInt64(key string) int64 {
	return dc.AddInt64(key, 1)
}

func (dc *dataSet) DecInt64(key string) int64 {
	return dc.AddInt64(key, -1)
}

func (dc *dataSet) SetUint64(key string, val uint64) { dc.Set(key, val) }

func (dc *dataSet) GetUint64(key string) (uint64, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(uint64); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(uint64(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddUint64(key string, d uint64) uint64 {
	if v, ok := dc.GetUint64(key); ok {
		v += d
		dc.SetUint64(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddUint64"))
	}
}

func (dc *dataSet) SubUint64(key string, d uint64) uint64 {
	if v, ok := dc.GetUint64(key); ok {
		v -= d
		dc.SetUint64(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "SubUint64"))
	}
}

func (dc *dataSet) IncUint64(key string) uint64 {
	return dc.AddUint64(key, 1)
}

func (dc *dataSet) DecUint64(key string) uint64 {
	return dc.SubUint64(key, 1)
}

func (dc *dataSet) SetFloat32(key string, val float32) { dc.Set(key, val) }

func (dc *dataSet) GetFloat32(key string) (float32, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(float32); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(float32(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddFloat32(key string, d float32) float32 {
	if v, ok := dc.GetFloat32(key); ok {
		v += d
		dc.SetFloat32(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddFloat32"))
	}
}

func (dc *dataSet) SetFloat64(key string, val float64) { dc.Set(key, val) }

func (dc *dataSet) GetFloat64(key string) (float64, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(float64); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(float64(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) AddFloat64(key string, d float64) float64 {
	if v, ok := dc.GetFloat64(key); ok {
		v += d
		dc.SetFloat64(key, v)
		return v
	} else {
		panic(ErrValueNotExist(key, "AddFloat64"))
	}
}

func (dc *dataSet) SetDuration(key string, val time.Duration) { dc.Set(key, val) }

func (dc *dataSet) GetDuration(key string) (time.Duration, bool) {
	if val := dc.Get(key); val == nil {
		return 0, false
	} else if v, ok := val.(time.Duration); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(time.Duration(0)), reflect.TypeOf(val).Name()))
	}
}

func (dc *dataSet) SetTime(key string, val time.Time) { dc.Set(key, val) }

func (dc *dataSet) GetTime(key string) (time.Time, bool) {
	if val := dc.Get(key); val == nil {
		return time.Time{}, false
	} else if v, ok := val.(time.Time); ok {
		return v, true
	} else {
		panic(ErrGetValueType(key, reflect.TypeOf(time.Time{}), reflect.TypeOf(val).Name()))
	}
}
