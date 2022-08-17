package mapData

import (
	"fmt"
	flyfish "github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/module"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/Map"
	"initialthree/pkg/json"
	"time"
)

type mapdata struct {
	//所处地图及坐标
	MapID int32 `json:"id"`
	X     int32 `json:"x"`
	Y     int32 `json:"y"`
	Z     int32 `json:"z"`
	Angle int32 `json:"ag"`
}

type UserMapData struct {
	userI module.UserI
	data  mapdata
	*module.ModuleSaveBase
}

func (this *UserMapData) SetMapID(mapID int32) {
	if this.data.MapID != mapID {
		this.data.MapID = mapID
		this.SetDirty(this.ModuleType().String())
	}
}

func (this *UserMapData) GetMapID() int32 {
	return this.data.MapID
}

func (this *UserMapData) SetPos(x, y, z, angle int32) {
	this.data.X = x
	this.data.Y = y
	this.data.Z = z
	this.data.Angle = angle
	this.SetDirty(this.ModuleType().String())

}

func (this *UserMapData) GetPos() (int32, int32, int32, int32) {
	return this.data.X, this.data.Y, this.data.Z, this.data.Angle
}

func (this *UserMapData) FlushAllToClient(seqNo ...uint32) {}

func (this *UserMapData) ModuleType() module.ModuleType {
	return module.MapData
}

func (this *UserMapData) Init(fields map[string]*flyfish.Field) error {
	field, ok := fields[this.ModuleType().String()]

	if ok && field.GetBlob() != nil {
		var data mapdata
		err := json.Unmarshal(field.GetBlob(), &data)
		if err != nil {
			return fmt.Errorf("unmarshal: %s", err)
		} else {
			this.data = data
		}
	} else {
		//设置主城ID
		this.data.MapID = Global.GetID(1).MapID
		mapDef := Map.GetID(this.data.MapID)
		this.data.X = mapDef.DefaultPositionStruct.X
		this.data.Y = mapDef.DefaultPositionStruct.Y
		this.data.Z = mapDef.DefaultPositionStruct.Z
		this.data.Angle = mapDef.DefaultRotation
		this.SetDirty(this.ModuleType().String())
	}

	return nil
}

func (this *UserMapData) ReadOut() *module.ReadOutCommand {
	return &module.ReadOutCommand{
		Table:  "user_module_data",
		Key:    this.userI.GetIDStr(),
		Fields: []string{this.ModuleType().String()},
		Module: this,
	}
}

func (this *UserMapData) Tick(now time.Time) {}

//打包道具格子
func (this *UserMapData) PackAll() interface{} {
	return nil
}

func (this *UserMapData) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	data, err := json.Marshal(this.data)
	if nil != err {
		return nil
	}

	cmd := &module.WriteBackCommand{
		Table: "user_module_data",
		Key:   this.userI.GetIDStr(),
		Fields: []*module.WriteBackFiled{{
			Name:  this.ModuleType().String(),
			Value: data,
		}},
		Module: this,
	}
	return cmd
}

func (this *UserMapData) FlushDirtyToClient() {}

func init() {
	module.RegisterModule(module.MapData, func(userI module.UserI) module.ModuleI {
		m := &UserMapData{
			userI: userI,
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
