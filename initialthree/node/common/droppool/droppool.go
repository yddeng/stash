package droppool

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/table/excel/DataTable/DropPool"
	"initialthree/pkg/util"
	"initialthree/protocol/cs/message"
)

func DropType2IOType(dropType int) int {
	var ioType int
	switch dropType {
	case enumType.DropType_Item:
		ioType = enumType.IOType_Item
	case enumType.DropType_Equip:
		ioType = enumType.IOType_Equip
	case enumType.DropType_Character:
		ioType = enumType.IOType_Character
	case enumType.DropType_UsualAttribute:
		ioType = enumType.IOType_UsualAttribute
	case enumType.DropType_Weapon:
		ioType = enumType.IOType_Weapon
	default:
		panic(fmt.Sprintf("droppool dropType trans to ioType %d not found", dropType))
	}
	return ioType
}

/*
 掉落池算法
*/

type AwardInfo struct {
	Type  int32
	ID    int32
	Count int32
}

func (ai *AwardInfo) ToMessage() *message.AwardInfo {
	return &message.AwardInfo{
		Type:  proto.Int32(ai.Type),
		ID:    proto.Int32(ai.ID),
		Count: proto.Int32(ai.Count),
	}
}

func (ai *AwardInfo) ToResDesc() inoutput.ResDesc {
	return inoutput.ResDesc{
		ID:    ai.ID,
		Count: ai.Count,
		Type:  DropType2IOType(int(ai.Type)),
	}
}

type Award struct {
	Infos map[string]*AwardInfo // key 为 Type:ID
}

func NewAward() *Award {
	return &Award{
		Infos: map[string]*AwardInfo{},
	}
}

func (this *Award) IsZero() bool {
	return this == nil || len(this.Infos) == 0
}

func (this *Award) ToMessageAward() *message.Award {
	msg := &message.Award{
		AwardInfos: make([]*message.AwardInfo, 0, len(this.Infos)),
	}
	for _, ai := range this.Infos {
		msg.AwardInfos = append(msg.AwardInfos, ai.ToMessage())
	}
	return msg
}

func (this *Award) ToResDesc() []inoutput.ResDesc {
	out := make([]inoutput.ResDesc, 0, len(this.Infos))
	for _, ai := range this.Infos {
		out = append(out, ai.ToResDesc())
	}
	return out
}

func (this *Award) AddInfo(dropType, dropID, count int32) {
	this.addInfo(dropType, dropID, count)
}

// 存在多种类型掉落物，防止不同类型物品的道具id相同
func makeKey(tt, id int32) string {
	return fmt.Sprintf("%d:%d", tt, id)
}

func parseKey(key string) (int32, int32) {
	var tt, id int32
	n, err := fmt.Sscanf(key, "%d:%d", &tt, &id)
	if n != 2 || err != nil {
		return 0, 0
	}
	return tt, id
}

func (this *Award) addInfo(dropType, dropID, count int32) {
	key := makeKey(dropType, dropID)
	if v, ok := this.Infos[key]; ok {
		v.Count += count
	} else {
		this.Infos[key] = &AwardInfo{
			Type:  dropType,
			ID:    dropID,
			Count: count,
		}
	}
}

func waveCount(dropCount, wave int32) int32 {
	if wave <= 0 {
		return dropCount
	}

	waveCount := util.Random(-wave, wave)
	count := dropCount + waveCount
	if count <= 0 {
		count = 1
	}

	return count
}

//随机掉落 种类在最大与最小值之间
func (this *Award) dropPoolRand(dropDef *DropPool.DropPool, up ...*DropPool.UpWeight) {
	//掉落数量限制
	var dropCount int32
	var repeated = dropDef.Repeatable

	if dropDef.MinCount <= 0 {
		dropCount = 1
	} else {
		if dropDef.MaxCount <= 0 || dropDef.MaxCount == dropDef.MinCount {
			dropCount = dropDef.MinCount
		} else {
			dropCount = util.Random(dropDef.MinCount, dropDef.MaxCount)
		}
	}
	// 不可重复,掉落数大于配置的物品数量
	if !repeated && dropCount > int32(len(dropDef.DropList)) {
		dropCount = int32(len(dropDef.DropList))
	}

	if dropCount == 0 {
		return
	}

	var weight, totalWeight int32 // 权重
	dropList := make([]*DropPool.DropList_, len(dropDef.DropList))
	for i, v := range dropDef.DropList {
		if rate, ok := DropPool.FindUpRate(v.Type, v.ID, up); ok {
			weight = v.Weight * rate
		} else {
			weight = v.Weight
		}
		dropList[i] = &DropPool.DropList_{
			Type:   v.Type,
			ID:     v.ID,
			Count:  v.Count,
			Wave:   v.Wave,
			Weight: weight,
		}
		totalWeight += weight
	}

	if totalWeight <= 0 {
		return
	}

	var curCount int32 = 0 //掉落数计数
	for curCount < dropCount {
		if totalWeight <= 0 {
			break
		}
		weight = util.Random(1, totalWeight)
		for idx, v := range dropList {
			weight -= v.Weight
			if weight <= 0 {
				if repeated {
					if v.Type == enumType.DropType_Pool {
						this.dropPoolPool(v.ID, waveCount(v.Count, v.Wave))
					} else {
						this.addInfo(v.Type, v.ID, waveCount(v.Count, v.Wave))
					}
					curCount += 1
				} else {
					if v.Type == enumType.DropType_Pool {
						this.dropPoolPool(v.ID, waveCount(v.Count, v.Wave))
					} else {
						this.addInfo(v.Type, v.ID, waveCount(v.Count, v.Wave))
					}
					curCount += 1
					// 移除已经掉落的
					dropList = append(dropList[:idx], dropList[idx+1:]...)
					totalWeight -= v.Weight
				}
				break
			}
		}
	}
}

//固定掉落
func (this *Award) dropPoolFixed(dropDef *DropPool.DropPool) {
	for _, v := range dropDef.DropList {
		if v.Type == enumType.DropType_Pool {
			this.dropPoolPool(v.ID, waveCount(v.Count, v.Wave))
		} else {
			this.addInfo(v.Type, v.ID, waveCount(v.Count, v.Wave))
		}
	}
}

//池中池的处理
func (this *Award) dropPoolPool(id, count int32) {
	for i := int32(0); i < count; i++ {
		poolAward := DropAward(id)
		for id, info := range poolAward.Infos {
			info_ := this.Infos[id]
			if info_ != nil {
				info_.Count += info.Count
			} else {
				this.Infos[id] = info
			}
		}
	}
}

func DropWithPool(pool *DropPool.DropPool, up ...*DropPool.UpWeight) (award *Award) {
	award = NewAward()
	if pool == nil {
		return
	}

	//判读掉落池类型
	switch pool.TypeEnum {
	case enumType.DropPoolType_Rand:
		award.dropPoolRand(pool, up...)

	case enumType.DropPoolType_Fixed: //固定掉落
		award.dropPoolFixed(pool)

	default:
		panic("not implemented")
	}

	return award
}

func DropWithID(poolID int32, up ...*DropPool.UpWeight) (award *Award) {
	dropDef := DropPool.GetID(poolID)
	return DropWithPool(dropDef, up...)
}

//奖励计算
func DropAward(poolIDs ...int32) *Award {
	awards := make([]*Award, 0, len(poolIDs))
	for _, id := range poolIDs {
		award := DropWithID(id)
		awards = append(awards, award)
	}

	return AllAwardsToOne(awards)
}

//多个奖励合并为一个奖励
func AllAwardsToOne(awards []*Award) *Award {
	award := NewAward()
	for _, a := range awards {
		if !a.IsZero() {
			for id, info := range a.Infos {
				info_ := award.Infos[id]
				if info_ != nil {
					info_.Count += info.Count
				} else {
					award.Infos[id] = info
				}
			}
		}
	}
	return award
}
