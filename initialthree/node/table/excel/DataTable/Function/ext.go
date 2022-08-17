package Function

import (
	"errors"
	"fmt"
	"initialthree/node/common/enumType"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type Unlock struct {
	Type      int32
	Level     int
	DungeonID int
	Time      []int64
}

var unlock atomic.Value

func GetUnlock(id int32) []*Unlock {
	return unlock.Load().(map[int32][]*Unlock)[id]
}

func (this *Table) AfterLoad() {
	idMap := GetIDMap()
	values := make(map[int32][]*Unlock, len(idMap))
	var err error
	for id, v := range idMap {
		list := make([]*Unlock, 0, len(v.Unlock))
		for _, l := range v.Unlock {
			un := &Unlock{Type: l.Type}
			switch l.Type {
			case enumType.FunctionUnlockType_PlayerLevel:
				un.Level, err = strconv.Atoi(l.Arg)
			case enumType.FunctionUnlockType_DungeonPass:
				un.DungeonID, err = strconv.Atoi(l.Arg)
			case enumType.FunctionUnlockType_TimeRange:
				ss := strings.Split(l.Arg, ",")
				var year, month, day, hour int
				for _, s := range ss {
					_, e := fmt.Sscanf(s, "%d:%d:%d:%d", &year, &month, &day, &hour)
					if e != nil {
						err = e
						break
					}
					un.Time = append(un.Time, time.Date(year, time.Month(month), day, hour, 0, 0, 0, time.Local).Unix())
				}
			case enumType.FunctionUnlockType_Default:
			default:
				err = errors.New("function unlock type invalid")
			}

			if err != nil {
				//panic("DataTable/Function/ext : afterLoad" + err.Error())
			}
			list = append(list, un)
		}
		values[id] = list
	}

	unlock.Store(values)
}
