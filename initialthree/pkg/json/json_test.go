package json

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"
)

type equip struct {
	InsID            uint32  `json:"ins_id"`
	ConfigID         int32   `json:"cfg_id"`
	Level            int32   `json:"level"`
	Exp              int32   `json:"exp"`              // 当前经验
	RandomAttribId   int32   `json:"random_attrib_id"` // 随机的技能ID，位置1
	Refine           []int32 `json:"refine"`           // 长度为2，第一个是配置表的固定技能ID，第二个是随机ID
	EquipCharacterId int32   `json:"e_c_id"`           //装备后的角色ID
	IsLock           bool    `json:"is_lock"`
	GetTime          int64   `json:"get_time"` // 获取时间
}

func BenchmarkStd(b *testing.B) {
	equips := []equip{}
	for i := 0; i < 1000; i++ {
		e := equip{
			InsID:            uint32(i + 1),
			ConfigID:         int32(rand.Int31()),
			Level:            int32(rand.Int31()),
			Exp:              int32(rand.Int31()),
			RandomAttribId:   int32(rand.Int31()),
			EquipCharacterId: int32(rand.Int31()),
			IsLock:           false,
			GetTime:          int64(time.Now().Unix()),
		}

		for j := 0; j < 10; j++ {
			e.Refine = append(e.Refine, int32(rand.Int31()))
		}

		equips = append(equips, e)
	}

	numLoops := b.N
	for i := 0; i < numLoops; i++ {
		b, _ := json.Marshal(equips)
		eequips := []equip{}
		json.Unmarshal(b, &eequips)
	}
}

func BenchmarkJsoniter(b *testing.B) {
	equips := []equip{}
	for i := 0; i < 1000; i++ {
		e := equip{
			InsID:            uint32(i + 1),
			ConfigID:         int32(rand.Int31()),
			Level:            int32(rand.Int31()),
			Exp:              int32(rand.Int31()),
			RandomAttribId:   int32(rand.Int31()),
			EquipCharacterId: int32(rand.Int31()),
			IsLock:           false,
			GetTime:          int64(time.Now().Unix()),
		}

		for j := 0; j < 10; j++ {
			e.Refine = append(e.Refine, int32(rand.Int31()))
		}

		equips = append(equips, e)
	}

	numLoops := b.N
	for i := 0; i < numLoops; i++ {
		b, _ := Marshal(equips)
		eequips := []equip{}
		Unmarshal(b, &eequips)
	}
}
