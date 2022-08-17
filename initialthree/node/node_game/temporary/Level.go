package temporary

import (
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"time"
)

var fightID int64

type LevelFightInfo struct {
	user           UserI
	Tos            *message.LevelFightToS
	Toc            *message.LevelFightToC //
	FightID        int64
	ResurrectCount int // 已复活次数
}

func NewLevelFightInfo(user UserI, tos *message.LevelFightToS, toc *message.LevelFightToC) *LevelFightInfo {
	fightID++
	return &LevelFightInfo{user: user, FightID: fightID, Tos: tos, Toc: toc}
}

func (m *LevelFightInfo) UserDisconnect() {
	//m.user.ClearTemporary(TempLevelFight)
}

func (m *LevelFightInfo) UserLogout() {
	m.user.ClearTemporary(TempLevelFight)
}

func (m *LevelFightInfo) Tick(now time.Time) {}

/*
 *
 */

type levelMarshaler struct{}

func (m *levelMarshaler) Marshal(temp TemporaryI) ([]byte, error) {
	info := temp.(*LevelFightInfo)
	return json.Marshal(info)
}

func (m *levelMarshaler) Unmarshal(user UserI, data []byte) (TemporaryI, error) {
	var info LevelFightInfo
	err := json.Unmarshal(data, &info)
	if err != nil {
		return nil, err
	}
	info.user = user
	return &info, nil
}

func init() {
	registerTempDataProcess(TempLevelFight, "level", &levelMarshaler{})
}
