package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/node/common/db"
	"initialthree/pkg/json"
	"initialthree/zaplogger"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
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

	k := os.Args[1]

	rand.Seed(time.Now().UnixNano())

	c, err := db.NewClient("pgsql", "10.50.31.13", 5432, "initialthree2", "sgzr2", "sniperHWfeiyu2019")
	if err != nil {
		panic(err)
	}

	logger := zaplogger.NewZapLogger("weapon.log", "log", "debug", 100, 14, 10, true)
	zaplogger.InitLogger(logger)
	db.FlyfishInit([]string{"flyfish@dir@localhost:10050"}, logger)
	f := db.GetFlyfishClient("dir")

	w := make([]map[uint32]*Weapon, 4)

	num := 0

	if k == "1" {
		c.GetAll("game_user", func(i map[string]interface{}) error {
			key := i["__key__"].(string)
			id := i["id"].(int64)

			ret, err := c.Get("weapon", fmt.Sprintf("%d", id))
			if err != nil {
				zaplogger.GetSugar().Infof("%s %d weapon get failed %s", key, id, err)
				return nil
			}

			count := 0
			num++
			for i := 0; i < 4; i++ {
				dbname := fmt.Sprintf("%s%d", "slice", i)
				field, ok := ret[dbname]
				if !ok {
					zaplogger.GetSugar().Infof("%s %s not found \n", id, dbname)
					continue
				}
				err := json.Unmarshal(field.([]byte), &w[i])
				if err != nil {
					zaplogger.GetSugar().Infof("unmarshal:%s %s", string(field.([]byte)), err)
					continue
				}
				count += len(w[i])
			}

			//if count == 0 {
			zaplogger.GetSugar().Infof("%s %d weapon count %d\n", key, id, count)
			//}
			if count == 0 {
				zaplogger.GetSugar().Infof("%s %d weapon count 0 \n", key, id)
			}

			result := f.GetAll("weapon", fmt.Sprintf("%d", id)).Exec()
			if errcode.GetCode(result.ErrCode) != errcode.Errcode_ok {
				zaplogger.GetSugar().Infof("%s %d  errcode %s", key, id, errcode.GetErrorDesc(result.ErrCode))
			}

			return nil
		})
	} else if k == "2" {
		c.GetAll("weapon", func(ret map[string]interface{}) error {
			id := ret["__key__"].(string)

			count := 0
			num++
			for i := 0; i < 4; i++ {
				dbname := fmt.Sprintf("%s%d", "slice", i)
				field, ok := ret[dbname]
				if !ok {
					zaplogger.GetSugar().Infof("%s %s not found \n", id, dbname)
					continue
				}
				err := json.Unmarshal(field.([]byte), &w[i])
				if err != nil {
					zaplogger.GetSugar().Infof("unmarshal:%s %s", string(field.([]byte)), err)
					continue
				}
				count += len(w[i])
			}

			if count == 0 {
				zaplogger.GetSugar().Infof("%s weapon count 0 \n", id)
			}

			//result := f.GetAll("weapon", id).Exec()
			//if errcode.GetCode(result.ErrCode) != errcode.Errcode_ok {
			//	zaplogger.GetSugar().Infof("%s  errcode %s", id, errcode.GetErrorDesc(result.ErrCode))
			//}

			return nil
		})
	} else if k == "3" {
		wg := sync.WaitGroup{}
		for i := 0; i < 5; i++ {
			cdb, err := sqlOpen("pgsql", "10.50.31.13", 5432, "initialthree2", "sgzr2", "sniperHWfeiyu2019")
			if err != nil {
				panic(err)
			}
			start := i*2000 + 1
			wg.Add(1)
			go func(s int) {
				for i := s; i < s+2000; {
					num := rand.Int()%5 + 5
					keys := make([]string, 0, num)
					for j := 0; j < num && i+j < 10000; j++ {
						keys = append(keys, fmt.Sprintf("robot1_%d:1", i+j))
					}
					i += num

					checkKey(cdb, "game_user", keys)
				}
				wg.Done()
			}(start)
		}
		wg.Wait()
	}

	zaplogger.GetSugar().Infof("end %d", num)

}

func checkKey(cdb *sqlx.DB, table string, keys []string) {

	sqlStr := `
SELECT __key__ FROM "%s" 
WHERE `

	s := fmt.Sprintf(sqlStr, table)
	ss := make([]string, 0, len(keys))
	for _, k := range keys {
		ss = append(ss, fmt.Sprintf("__key__ = '%s'", k))
	}

	s += strings.Join(ss, " or ") + ";"

	//fmt.Println(s)

	rows, err := cdb.Query(s)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	findk := make([]string, 0, len(keys))
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			panic(err)
		}

		findk = append(findk, key)
		//fmt.Println(findk, key)

		find := false
		for _, k := range keys {
			if k == key {
				find = true
				break
			}
		}

		if !find {
			fmt.Printf(" %v 多读到key %s \n", keys, key)
		}
	}

	for _, k := range keys {
		find := false
		for _, kk := range findk {
			if kk == k {
				find = true
				break
			}
		}
		if !find {
			fmt.Printf(" %v 少读到key %s \n", keys, k)
		}
	}

}

func pgsqlOpen(host string, port int, dbname string, user string, password string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", host, port, dbname, user, password)
	return sqlx.Open("postgres", connStr)
}

func mysqlOpen(host string, port int, dbname string, user string, password string) (*sqlx.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbname)
	return sqlx.Open("mysql", connStr)
}

func sqlOpen(sqlType string, host string, port int, dbname string, user string, password string) (*sqlx.DB, error) {
	if sqlType == "mysql" {
		return mysqlOpen(host, port, dbname, user, password)
	} else {
		return pgsqlOpen(host, port, dbname, user, password)
	}
}
