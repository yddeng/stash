package main

import (
	"fmt"
	"initialthree/node/common/db"
	"initialthree/pkg/json"
)

type Weapon struct {
	InsID            uint32 `json:"ins_id"`
	ConfigID         int32  `json:"cfg_id"`
	Level            int32  `json:"level"`
	Exp              int32  `json:"exp"` // 当前经验
	Refine           int32  `json:"refine"`
	BreakTimes       int32  `json:"break_times"` // 突破次数
	EquipCharacterID int32  `json:"e_c_id"`      //装备后的角色ID
	IsLock           bool   `json:"is_lock"`
	GetTime          int64  `json:"get_time"` // 获取时间
}

func main() {
	c, err := db.NewClient("pgsql", "10.50.31.13", 5432, "initialthree2", "sgzr2", "sniperHWfeiyu2019")
	if err != nil {
		panic(err)
	}

	w := []map[uint32]*Weapon{}

	num := 0
	err = c.GetAll("game_user", func(ret map[string]interface{}) error {
		key := ret["__key__"].(string)
		id := ret["__key__"].(int64)
		count := 0
		num++
		for i := 0; i < 4; i++ {
			dbname := fmt.Sprintf("%s%d", "slice", i)
			field, ok := ret[dbname]
			if !ok {
				fmt.Printf("%s %s not found \n", id, dbname)
			}
			err := json.Unmarshal(field.([]byte), &w[i])
			if err != nil {
				fmt.Printf("unmarshal:%s %s", string(field.([]byte)), err)
			}
			count += len(w[i])
		}

		if count == 0 {
			fmt.Printf("%s weapon count is 0", id)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("end", num)
}
