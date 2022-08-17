package behavior

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/enumType"
	"initialthree/node/table/excel/DataTable/Dungeon"
	"initialthree/node/table/excel/DataTable/MainChapter"
	"initialthree/node/table/excel/DataTable/MainDungeon"
	csmsg "initialthree/protocol/cs/message"
	"initialthree/robot/internal"
	"initialthree/robot/net"
	"initialthree/robot/robot/module"
	"initialthree/robot/types"
	"math/rand"
	"time"

	. "github.com/GodYY/bevtree"
	"github.com/golang/protobuf/proto"
)

const mainDungeon = BevType("maindungeon")

func init() {
	regBevType(mainDungeon, func() Bev { return new(BevMainDungeon) })
}

type BevMainDungeon struct {
	Random bool `xml:"random,attr"`
}

func (BevMainDungeon) BevType() BevType { return mainDungeon }

func (b *BevMainDungeon) CreateInstance() BevInstance {
	return &BevMainDungeonInst{BevMainDungeon: b}
}

func (b *BevMainDungeon) DestroyInstance(bi BevInstance) {
	bi.(*BevMainDungeonInst).BevMainDungeon = nil
}

type BevMainDungeonInst struct {
	bev
	*BevMainDungeon
	mainDungeonID  int32
	dungeonID      int32
	dungeonCfg     *Dungeon.Dungeon
	levelFightResp *csmsg.LevelFightToC
	useTime        int32
}

// 行为类型
func (b *BevMainDungeonInst) BevType() BevType {
	return mainDungeon
}

func (b *BevMainDungeonInst) OnInit(ctx Context) bool {
	b.bev.OnInit(ctx)

	var nextDungeon int32
	if !b.Random {
		nextDungeon = b.getNextDungeon()
		if nextDungeon == internal.InvalidID {
			return false
		}
	} else {
		nextDungeon = b.getRandomRemainCountDungeon()
		if nextDungeon == internal.InvalidID {
			return false
		}
	}

	msg := &csmsg.LevelFightToS{}

	mainDungeonCfg := MainDungeon.GetID(nextDungeon)
	msg.DungeonID = proto.Int32(mainDungeonCfg.DungeonID)
	dungeonCfg := Dungeon.GetID(mainDungeonCfg.DungeonID)

	msg.CharacterTeam = &csmsg.CharacterTeam{}
	moduleCharacter := b.player.GetModule(module.Module_Character).(*module.ModuleCharacter)
	switch dungeonCfg.TeamLimitTypeEnum {
	case enumType.DungeonTeamLimitType_CharacterCountLimit:
		msg.CharacterTeam.CharacterList = moduleCharacter.RandomCharacters(int(dungeonCfg.TeamLimitValue), nil)

	case enumType.DungeonTeamLimitType_AppointCharacter:
		characters := make([]int32, 3)
		characters[0] = dungeonCfg.TeamLimitValue
		extraCharacters := moduleCharacter.RandomCharacters(2, map[int32]struct{}{dungeonCfg.TeamLimitValue: struct{}{}})
		copy(characters[1:], extraCharacters)
		msg.CharacterTeam.CharacterList = characters
	}

	b.mainDungeonID = nextDungeon
	b.dungeonID = mainDungeonCfg.DungeonID
	b.dungeonCfg = dungeonCfg

	b.sendMessage(msg, b.onLevelFight)

	b.player.Debugf("request to maindungeon %d level fight", nextDungeon)

	return true
}

func (b *BevMainDungeonInst) OnTerminate(ctx Context) {
	b.mainDungeonID = internal.InvalidID
	b.dungeonID = internal.InvalidID
	b.dungeonCfg = nil
	b.levelFightResp = nil
	b.useTime = 0
	b.player.RemTimer(timerIDMainDungeonFight)
	b.bev.OnTerminate(ctx)
}

func (b *BevMainDungeonInst) getNextDungeon() int32 {
	moduleMainDungeon := b.player.GetModule(module.Module_MainDungeons).(*module.ModuleMainDungeons)
	nextDungeon := moduleMainDungeon.GetNextNormalDungeon()
	if nextDungeon == internal.InvalidID {
		b.player.Debug("no next maindungeon")
		return internal.InvalidID
	}

	moduleAttr := b.player.GetModule(module.Module_Attr).(*module.ModuleAttr)
	mainDungeonCfg := MainDungeon.GetID(nextDungeon)
	mainChapterCfg := MainChapter.GetID(mainDungeonCfg.ChapterID)
	if mainChapterCfg.PlayerLevelLimit > int32(moduleAttr.Level()) {
		b.player.Debug("level low")
		return internal.InvalidID
	}

	if moduleMainDungeon.DungeonCount(nextDungeon) <= 0 {
		b.player.Debug("no count")
		return internal.InvalidID
	}

	dungeonCfg := Dungeon.GetID(mainDungeonCfg.DungeonID)
	for _, v := range Dungeon.GetUnlock(mainDungeonCfg.DungeonID) {
		switch v.Type {
		case enumType.DungeonUnlockType_MainChapter:
			for _, vv := range v.Args {
				if !moduleMainDungeon.IsDungeonPass(vv) {
					b.player.Debug("unlock dungeon not pass")
					return internal.InvalidID
				}
			}

		case enumType.DungeonUnlockType_PlayerLevel:
			if moduleAttr.Level() < int64(v.Args[0]) {
				b.player.Debug("level low")
				return internal.InvalidID
			}
		}
	}

	if dungeonCfg.CostTypeEnum != 0 {
		switch dungeonCfg.CostTypeEnum {
		case enumType.DungeonCostType_FatigueCost:
			if moduleAttr.Fatigue() < int64(dungeonCfg.CostArgs) {
				b.player.Debug("fatigue not enough")
				return internal.InvalidID
			}
		}
	}

	moduleCharacter := b.player.GetModule(module.Module_Character).(*module.ModuleCharacter)
	switch dungeonCfg.TeamLimitTypeEnum {
	case enumType.DungeonTeamLimitType_CharacterCountLimit:
		if moduleCharacter.CharacterCount() < int(dungeonCfg.TeamLimitValue) {
			b.player.Debug("character count not enough")
			return internal.InvalidID
		}

	case enumType.DungeonTeamLimitType_AppointCharacter:
		if !moduleCharacter.HasCharacter(dungeonCfg.TeamLimitValue) {
			b.player.Debug("character not found")
			return internal.InvalidID
		}
	}

	return nextDungeon
}

func (b *BevMainDungeonInst) getRandomRemainCountDungeon() int32 {
	moduleMainDungeon := b.player.GetModule(module.Module_MainDungeons).(*module.ModuleMainDungeons)
	dungeonID := moduleMainDungeon.GetRandomDungeonCouldChallenged()
	if dungeonID == internal.InvalidID {
		b.player.Debugf("no remain count dungeon")
	}
	return dungeonID
}

var timerIDMainDungeonFight = types.NewTimerID("mainDungeonFight")

func (b *BevMainDungeonInst) onLevelFight(r player, msg *codecs.Message) bool {
	if !net.IsMessageOK(msg) {
		b.player.Errorf("maindungeon %d level fight failed: %s", b.mainDungeonID, net.GetErrCodeStr(msg.GetErrCode()))
		b.terminate(false)
		return false
	}

	b.player.Infof("maindungeon %d start fight", b.mainDungeonID)

	b.levelFightResp = msg.GetData().(*csmsg.LevelFightToC)

	useTime := int32(0)

	if b.dungeonCfg.ClearTimeLimit <= 0 {
		useTime = rand.Int31n(21) + 10
	} else {
		useTime = int32(float64(b.dungeonCfg.ClearTimeLimit) * float64(60+rand.Int31n(41)) / float64(100))
	}

	b.useTime = useTime
	b.player.AddTimer(timerIDMainDungeonFight, time.Duration(useTime)*time.Second, nil, b.onLevelFightTimer)

	return false
}

func (b *BevMainDungeonInst) onLevelFightTimer(r player, ctx interface{}) {
	b.player.Debugf("maindungeon %d fight end", b.mainDungeonID)
	b.levelFightEnd()
}

func (b *BevMainDungeonInst) levelFightEnd() {
	msg := &csmsg.LevelFightEndToS{
		FightID: proto.Int64(b.levelFightResp.GetFightID()),
		UseTime: proto.Int32(b.useTime),
		Pass:    proto.Bool(true),
		Stars:   []bool{true, true, true},
	}

	b.sendMessage(msg, b.onLevelFightEnd)
	b.player.Debugf("request to maindungeon %d level fight end", b.mainDungeonID)
}

func (b *BevMainDungeonInst) onLevelFightEnd(r player, msg *codecs.Message) bool {
	if !net.IsMessageOK(msg) {
		b.player.Errorf("maindungeon %d level fight end failed: %s", b.mainDungeonID, net.GetErrCodeStr(msg.GetErrCode()))
		b.terminate(false)
		return false
	}

	b.player.Infof("maindungeon %d fight end successfully", b.mainDungeonID)
	b.terminate(true)
	return false
}
