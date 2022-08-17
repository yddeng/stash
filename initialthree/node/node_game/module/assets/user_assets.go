package assets

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"initialthree/pkg/json"
	"initialthree/zaplogger"
	"time"

	flyfish "github.com/sniperHW/flyfish/client"

	"initialthree/node/node_game/module"
	"initialthree/protocol/cs/message"
)

const (
	userAssetsTableName = "user_assets"
	userAssetsField     = "assets"
)

type Asset struct {
	Kv    map[int32]int32    `json:"kv"` // id count
	dirty map[int32]struct{} `json:"_"`
}

type UserAssets struct {
	userI  module.UserI
	assets map[int32]*Asset   // AssetType -> Asset
	dirty  map[int32]struct{} // AssetType
	*module.ModuleSaveBase
}

func (ua *UserAssets) ModuleType() module.ModuleType {
	return module.Assets
}

func (ua *UserAssets) Init(fields map[string]*flyfish.Field) error {
	field, ok := fields[userAssetsField]
	if ok && len(field.GetBlob()) != 0 {
		if err := json.Unmarshal(field.GetBlob(), &ua.assets); err != nil {
			return fmt.Errorf("unmarshal: %s", err)
		}
	} else {
		ua.SetDirty(userAssetsField)
	}
	return nil
}

func (ua *UserAssets) ReadOut() *module.ReadOutCommand {
	return &module.ReadOutCommand{
		Table:  userAssetsTableName,
		Key:    ua.userI.GetIDStr(),
		Fields: []string{userAssetsField},
		Module: ua,
	}
}

func (ua *UserAssets) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	if data, err := json.Marshal(ua.assets); err != nil {
		zaplogger.GetSugar().Error(err.Error())
		return nil
	} else {
		return &module.WriteBackCommand{
			Table: userAssetsTableName,
			Key:   ua.userI.GetIDStr(),
			Fields: []*module.WriteBackFiled{{
				Name:  userAssetsField,
				Value: data,
			}},
			Module: ua,
		}
	}
}

func (ua *UserAssets) GetAssetCount(assetType int32, id int32) (count int32, ok bool) {
	if asset, exist := ua.assets[assetType]; exist {
		count, ok = asset.Kv[id]
		return
	}
	return 0, false
}

func (ua *UserAssets) GetAsset(assetType int32) *Asset {
	return ua.assets[assetType]
}

func (ua *UserAssets) AddAsset(assetType int32, id, dt int32) {
	count, _ := ua.GetAssetCount(assetType, id)
	ua.SetAsset(assetType, id, dt+count)
}

func (ua *UserAssets) SetAsset(assetType int32, id, count int32) {
	if asset, exist := ua.assets[assetType]; exist {
		asset.Kv[id] = count
		if asset.dirty == nil {
			asset.dirty = map[int32]struct{}{}
		}
		asset.dirty[id] = struct{}{}
	} else {
		ua.assets[assetType] = &Asset{
			Kv:    map[int32]int32{id: count},
			dirty: map[int32]struct{}{id: {}},
		}
	}
	ua.SetDirty(userAssetsField)

	ua.dirty[assetType] = struct{}{}

}

func (ua *UserAssets) FlushDirtyToClient() {
	if len(ua.dirty) > 0 {
		msg := &message.AssetSyncToC{
			IsAll:      proto.Bool(false),
			SyncAssets: make([]*message.Asset, 0, len(ua.dirty)),
		}
		for assetType := range ua.dirty {
			asset := ua.assets[assetType]
			if len(asset.dirty) > 0 {
				msgAsset := &message.Asset{
					Type:    proto.Int32(assetType),
					AssetKv: make([]*message.AssetValue, 0, len(asset.dirty)),
				}
				for id := range asset.dirty {
					count := asset.Kv[id]
					msgAsset.AssetKv = append(msgAsset.AssetKv, &message.AssetValue{
						ID:    proto.Int32(id),
						Count: proto.Int32(count),
					})
				}
				msg.SyncAssets = append(msg.SyncAssets, msgAsset)
			}
			asset.dirty = map[int32]struct{}{}
		}
		ua.dirty = map[int32]struct{}{}
		ua.userI.Post(msg)
	}
}
func (ua *UserAssets) FlushAllToClient(seqNo ...uint32) {
	msg := &message.AssetSyncToC{
		IsAll:      proto.Bool(true),
		SyncAssets: make([]*message.Asset, 0, len(ua.assets)),
	}
	for t, sub := range ua.assets {
		msgAsset := &message.Asset{
			Type:    proto.Int32(t),
			AssetKv: make([]*message.AssetValue, 0, len(sub.Kv)),
		}
		for id, count := range sub.Kv {
			msgAsset.AssetKv = append(msgAsset.AssetKv, &message.AssetValue{
				ID:    proto.Int32(id),
				Count: proto.Int32(count),
			})
		}
		msg.SyncAssets = append(msg.SyncAssets, msgAsset)
		sub.dirty = map[int32]struct{}{}
	}
	ua.dirty = map[int32]struct{}{}
	ua.userI.Post(msg)
}

func (ua *UserAssets) Tick(t time.Time) {}

func init() {
	module.RegisterModule(module.Assets, func(userI module.UserI) module.ModuleI {
		asset := &UserAssets{
			userI:  userI,
			assets: map[int32]*Asset{},
			dirty:  map[int32]struct{}{},
		}
		asset.ModuleSaveBase = module.NewModuleSaveBase(asset)
		return asset
	})
}
