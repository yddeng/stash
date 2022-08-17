package module

import (
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/table/excel/DataTable/Dungeon"
	"initialthree/node/table/excel/DataTable/MainChapter"
	"initialthree/node/table/excel/DataTable/MainDungeon"
	"initialthree/protocol/cmdEnum"
	csmsg "initialthree/protocol/cs/message"
	"initialthree/robot/internal"
	"initialthree/robot/types"
	"initialthree/zaplogger"
	"math/rand"

	"github.com/golang/protobuf/proto"
)

const (
	Module_Unknown = iota
	Module_Attr
	Module_Base
	Module_Character
	Module_Backpack
	Module_Quest
	Module_Equip
	Module_Team
	Module_MainDungeons
	Module_MaterialDungeons
	Module_ScarsIngrain
	Module_RewardQuest
	Module_Weapon
)

type Robot = types.Robot
type Module = types.Module

type ModuleDefine struct {
	Name          string
	AssociatedMsg []uint16
	Create        func() Module
}

var moduleDefines = [...]ModuleDefine{
	Module_Unknown:          {Name: "Unknown", AssociatedMsg: []uint16{}, Create: func() Module { return nil }},
	Module_Attr:             {Name: "Attr", AssociatedMsg: []uint16{cmdEnum.CS_AttrSync}, Create: func() Module { return new(ModuleAttr) }},
	Module_Base:             {Name: "Base", AssociatedMsg: []uint16{cmdEnum.CS_BaseSync}, Create: func() Module { return new(ModuleBase) }},
	Module_Character:        {Name: "Charater", AssociatedMsg: []uint16{cmdEnum.CS_CharacterSync, cmdEnum.CS_CharacterTeamPrefabSync}, Create: func() Module { return new(ModuleCharacter) }},
	Module_Backpack:         {Name: "Backpack", AssociatedMsg: []uint16{cmdEnum.CS_BackpackSync}, Create: func() Module { return new(ModuleBackpack) }},
	Module_Quest:            {Name: "Quest", AssociatedMsg: []uint16{cmdEnum.CS_QuestSync}, Create: func() Module { return new(ModuleQuest) }},
	Module_Equip:            {Name: "Equip", AssociatedMsg: []uint16{cmdEnum.CS_EquipSync}, Create: func() Module { return new(ModuleEquip) }},
	Module_Team:             {Name: "Team", AssociatedMsg: []uint16{cmdEnum.CS_TeamSync, cmdEnum.CS_TeamPosSync, cmdEnum.CS_TeamStatusSync}, Create: func() Module { return new(ModuleTeam) }},
	Module_MainDungeons:     {Name: "MainDungeons", AssociatedMsg: []uint16{cmdEnum.CS_MainDungeonsSync}, Create: func() Module { return new(ModuleMainDungeons) }},
	Module_MaterialDungeons: {Name: "MaterialDungeons", AssociatedMsg: []uint16{cmdEnum.CS_MaterialDungeonSync}, Create: func() Module { return new(ModuleMaterialDungeons) }},
	Module_ScarsIngrain:     {Name: "ScarsIngrain", AssociatedMsg: []uint16{cmdEnum.CS_ScarsIngrainSync}, Create: func() Module { return new(ModuleScarsIngrain) }},
	Module_RewardQuest:      {Name: "RewardQuest", AssociatedMsg: []uint16{cmdEnum.CS_RewardQuestSync}, Create: func() Module { return new(ModuleRewardQuest) }},
	Module_Weapon:           {Name: "Weapon", AssociatedMsg: []uint16{cmdEnum.CS_WeaponSync}, Create: func() Module { return new(ModuleWeapon) }},
}

func GetModuleName(m int) string {
	if m < 0 || m >= len(moduleDefines) {
		return ""
	}
	return moduleDefines[m].Name
}

func TraverseDefines(f func(int, *ModuleDefine) bool) {
	for id, v := range moduleDefines {
		if !f(id, &v) {
			break
		}
	}
}

func CreateModules() []Module {
	modules := make([]Module, len(moduleDefines))
	for i, v := range moduleDefines {
		modules[i] = v.Create()
	}
	return modules
}

type ModuleAttr struct {
	attrs []*csmsg.Attr
}

func (m *ModuleAttr) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.AttrSyncToC)
	if sync.GetIsAll() || m.attrs == nil {
		m.attrs = make([]*csmsg.Attr, attr.AttrMax)
	}

	for _, v := range sync.Attrs {
		idx := int(v.GetId() - 1)
		m.attrs[idx] = v
	}
}

func (m *ModuleAttr) GetAttr(id int32) *csmsg.Attr {
	if id <= 0 || int(id) > len(m.attrs) {
		return nil
	}

	return m.attrs[id-1]
}

func (m *ModuleAttr) Level() int64 {
	return m.attrs[attr.Level-1].GetVal()
}

func (m *ModuleAttr) Fatigue() int64 {
	return m.attrs[attr.CurrentFatigue-1].GetVal()
}

type ModuleBase struct {
	data *csmsg.BaseSyncToC
}

func (m *ModuleBase) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.BaseSyncToC)
	m.data = sync
}

type ModuleCharacter struct {
	charaters   map[int32]*csmsg.Character
	teamPrefabs []*csmsg.CharacterTeamPrefab
}

func (m *ModuleCharacter) OnModuleSync(r Robot, msg proto.Message) {
	switch sync := msg.(type) {
	case *csmsg.CharacterSyncToC:
		if sync.GetIsAll() || m.charaters == nil {
			m.charaters = map[int32]*csmsg.Character{}
		}

		for _, v := range sync.GetCharacters() {
			m.charaters[v.GetCharacterID()] = v
		}

	case *csmsg.CharacterTeamPrefabSyncToC:
		m.teamPrefabs = sync.CharacterTeamPrefabs
	}
}

func (m *ModuleCharacter) CharacterCount() int {
	return len(m.charaters)
}

func (m *ModuleCharacter) HasCharacter(cid int32) bool {
	return m.charaters[cid] != nil
}

func (m *ModuleCharacter) RandomCharacters(n int, except map[int32]struct{}) []int32 {
	N := len(m.charaters) - len(except)
	if n > N {
		return nil
	}

	characters := make([]int32, n)
	i := 0
	for k, _ := range m.charaters {
		if _, ok := except[k]; ok {
			continue
		}

		if rand.Intn(N-i) > n {
			characters[i] = k
			i++
			n--
			if n <= 0 {
				break
			}
		}
	}

	return characters
}

type ModuleBackpack struct {
	items map[uint32]*csmsg.BackpackEntity
}

func (m *ModuleBackpack) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.BackpackSyncToC)

	if sync.GetAll() || m.items == nil {
		m.items = map[uint32]*csmsg.BackpackEntity{}
	}

	for _, v := range sync.GetEntities() {
		if v.GetCount() == 0 {
			delete(m.items, v.GetId())
		} else {
			m.items[v.GetId()] = v
		}
	}
}

type ModuleQuest struct {
	quests map[int32]*csmsg.Quest
}

func (m *ModuleQuest) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.QuestSyncToC)
	if sync.GetIsAll() || m.quests == nil {
		m.quests = map[int32]*csmsg.Quest{}
	}

	for _, v := range sync.Quests {
		m.quests[v.GetQuestID()] = v
	}
}

type ModuleEquip struct {
	equips map[uint32]*csmsg.Equip
}

func (m *ModuleEquip) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.EquipSyncToC)
	if sync.GetIsAll() || m.equips == nil {
		m.equips = map[uint32]*csmsg.Equip{}
	}

	for _, v := range sync.GetEquips() {
		m.equips[v.GetID()] = v
	}
}

type ModuleTeam struct {
	team *csmsg.Team
}

func (m *ModuleTeam) OnModuleSync(r Robot, msg proto.Message) {
	switch sync := msg.(type) {
	case *csmsg.TeamSyncToC:
		m.team = sync.GetUpdateTeam()

	case *csmsg.TeamPosSyncToC:

	case *csmsg.TeamStatusSyncToC:
	}
}

type mainChapter struct {
	*csmsg.MainChapter
	stars int32
}

type mainDungeon struct {
	*csmsg.MainDungeon
	stars int32
}

type ModuleMainDungeons struct {
	chapters                    map[int32]*mainChapter
	dungeons                    map[int32]*mainDungeon
	normalProgressDungeon       int32
	dungeonsCouldChallenged     map[int32]int
	dungeonCouldChallengedArray []int32

	nextStarAward struct {
		chapterID int32
		awardNo   int32
	}
}

func (m *ModuleMainDungeons) OnModuleSync(r Robot, msg proto.Message) {
	startChapterOfRefreshStartAward := int32(internal.InvalidID)

	sync := msg.(*csmsg.MainDungeonsSyncToC)
	if sync.GetAll() {
		m.chapters = map[int32]*mainChapter{}
		m.dungeons = map[int32]*mainDungeon{}
		m.normalProgressDungeon = internal.InvalidID
		m.dungeonsCouldChallenged = map[int32]int{}
		m.dungeonCouldChallengedArray = nil
		m.nextStarAward.chapterID = internal.InvalidID
		m.nextStarAward.awardNo = internal.InvalidID
		startChapterOfRefreshStartAward = 1
	} else {
		if m.chapters == nil {
			m.chapters = map[int32]*mainChapter{}
		}

		if m.dungeons == nil {
			m.dungeons = map[int32]*mainDungeon{}
		}

		if m.dungeonsCouldChallenged == nil {
			m.dungeonsCouldChallenged = map[int32]int{}
		}
	}

	for _, v := range sync.Chapters {
		m.updateChapter(v)

		chapterID := v.GetId()
		if chapterID == m.nextStarAward.chapterID && v.GetAwardFlag()[m.nextStarAward.awardNo] {
			// 更新下一个星级奖励章节
			if m.nextStarAward.awardNo == int32(len(v.GetAwardFlag())-1) {
				startChapterOfRefreshStartAward = chapterID + 1
			} else {
				m.nextStarAward.awardNo++
			}
		}
	}

	// 刷新下一个领奖章节
	if startChapterOfRefreshStartAward != internal.InvalidID {
		m.refreshNextStarAwardNo(startChapterOfRefreshStartAward)
	}

	for _, v := range sync.Dungeons {
		m.updateDungeon(v)
	}
}

func (m *ModuleMainDungeons) updateChapter(data *csmsg.MainChapter) {
	chapterID := data.GetId()
	chapter := m.chapters[chapterID]
	if chapter == nil {
		chapter = new(mainChapter)
		m.chapters[chapterID] = chapter
	}
	chapter.MainChapter = data
}

func (m *ModuleMainDungeons) addChapterStars(chapterId, stars int32) {
	chapter := m.chapters[chapterId]
	if chapter == nil {
		chapter = new(mainChapter)
		m.chapters[chapterId] = chapter
	}
	chapter.stars += stars
}

func (m *ModuleMainDungeons) updateDungeon(data *csmsg.MainDungeon) {
	dungeonID := data.GetId()
	dungeon := m.dungeons[dungeonID]
	if dungeon == nil {
		dungeon = new(mainDungeon)
		m.dungeons[dungeonID] = dungeon
	}
	dungeon.MainDungeon = data
	oldStars := dungeon.stars
	dungeon.stars = getMainDungeonStars(data)

	if dungeon.GetRemainCount() > 0 {
		if _, ok := m.dungeonsCouldChallenged[dungeonID]; !ok {
			idx := len(m.dungeonCouldChallengedArray)
			m.dungeonCouldChallengedArray = append(m.dungeonCouldChallengedArray, dungeonID)
			m.dungeonsCouldChallenged[dungeonID] = idx
		}
	} else {
		if idx, ok := m.dungeonsCouldChallenged[dungeonID]; ok {
			n := len(m.dungeonCouldChallengedArray) - 1
			m.dungeonCouldChallengedArray[idx] = m.dungeonCouldChallengedArray[n]
			m.dungeonCouldChallengedArray = m.dungeonCouldChallengedArray[0:n]
			delete(m.dungeonsCouldChallenged, dungeonID)
		}
	}

	dungeonCfg := MainDungeon.GetID(dungeonID)
	chapterCfg := MainChapter.GetID(dungeonCfg.ChapterID)
	switch chapterCfg.GetChapterType() {
	case enumType.MainChapterType_Normal:
		if m.normalProgressDungeon < dungeonID {
			m.normalProgressDungeon = dungeonID
		}

		m.addChapterStars(dungeonCfg.ChapterID, dungeon.stars-oldStars)
	}
}

func (m *ModuleMainDungeons) GetNextNormalDungeon() int32 {
	if m.normalProgressDungeon == 0 {
		return internal.InvalidID
	}

	if m.normalProgressDungeon == internal.InvalidID {
		chapterCfg := MainChapter.GetID(1)
		return chapterCfg.DungeonsArray[0].ID
	}

	dungeonCfg := MainDungeon.GetID(m.normalProgressDungeon)
	if dungeonCfg.NextMainDungeonID != internal.InvalidID {
		return dungeonCfg.NextMainDungeonID
	}

	chapterCfg := MainChapter.GetID(dungeonCfg.ChapterID + 1)
	if chapterCfg != nil {
		return chapterCfg.DungeonsArray[0].ID
	}

	return internal.InvalidID
}

func (m *ModuleMainDungeons) DungeonCount(dungeonID int32) int32 {
	dungeon := m.dungeons[dungeonID]
	if dungeon == nil {

		mainDungeonCfg := MainDungeon.GetID(dungeonID)
		if mainDungeonCfg == nil {
			return 0
		}

		dungeonCfg := Dungeon.GetID(mainDungeonCfg.DungeonID)
		if dungeonCfg == nil {
			return 0
		}

		return dungeonCfg.TimesLimit
	}

	return dungeon.GetRemainCount()
}

func (m *ModuleMainDungeons) GetRandomDungeonCouldChallenged() int32 {
	n := len(m.dungeonCouldChallengedArray)
	if n == 0 {
		return internal.InvalidID
	}

	return m.dungeonCouldChallengedArray[rand.Intn(n)]
}

func (m *ModuleMainDungeons) refreshNextStarAwardNo(startChapter int32) {
	m.nextStarAward.chapterID = internal.InvalidID
	m.nextStarAward.awardNo = internal.InvalidID

	for chapterID, chapterCfg := startChapter, MainChapter.GetID(startChapter); chapterCfg != nil; chapterID, chapterCfg = chapterID+1, MainChapter.GetID(chapterID+1) {
		chapter := m.chapters[chapterID]
		if chapter == nil {
			m.nextStarAward.chapterID = chapterID
			m.nextStarAward.awardNo = 0
			return
		} else {
			for i, v := range chapter.GetAwardFlag() {
				if !v {
					m.nextStarAward.chapterID = chapterID
					m.nextStarAward.awardNo = int32(i)
					return
				}
			}
		}
	}
}

func getMainDungeonStars(dungeon *csmsg.MainDungeon) int32 {
	if dungeon == nil {
		return 0
	}

	stars := int32(0)
	for _, v := range dungeon.GetStars() {
		if v {
			stars++
		}
	}

	return stars
}

func (m *ModuleMainDungeons) GetNextStarAwardNo() (chapterID, awardNo int32, couldBeClaimed bool) {
	if m.nextStarAward.chapterID == internal.InvalidID {
		return internal.InvalidID, internal.InvalidID, false
	}

	chapterCfg := MainChapter.GetID(m.nextStarAward.chapterID)
	if chapterCfg == nil {
		zaplogger.GetSugar().Panicf("MainChapter %d config not found", m.nextStarAward.chapterID)
		return m.nextStarAward.chapterID, m.nextStarAward.awardNo, false
	}

	award := chapterCfg.GetStarAward(m.nextStarAward.awardNo)
	if award == nil {
		zaplogger.GetSugar().Panicf("MainChapter %d do not have No.%d star award", m.nextStarAward.chapterID, m.nextStarAward.awardNo+1)
		return m.nextStarAward.chapterID, m.nextStarAward.awardNo, false
	}

	chapterStars := int32(0)
	chapter := m.chapters[m.nextStarAward.chapterID]
	if chapter != nil {
		chapterStars = chapter.stars
	}

	return m.nextStarAward.chapterID, m.nextStarAward.awardNo, chapterStars >= award.Stars
}

func (m *ModuleMainDungeons) IsDungeonPass(dungeonID int32) bool {
	return m.dungeons[dungeonID] != nil
}

type ModuleMaterialDungeons struct {
	dungeons []*csmsg.MaterialDungeon
}

func (m *ModuleMaterialDungeons) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.MaterialDungeonSyncToC)
	if sync.GetAll() {
		m.dungeons = sync.MaterialDungeons
	}

	m.dungeons = append(m.dungeons, sync.MaterialDungeons...)
}

type ModuleScarsIngrain struct {
	data *csmsg.ScarsIngrainSyncToC
}

func (m *ModuleScarsIngrain) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.ScarsIngrainSyncToC)
	m.data = sync
}

type ModuleRewardQuest struct {
	quests map[int32]*csmsg.RewardQuest
}

func (m *ModuleRewardQuest) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.RewardQuestSyncToC)
	if sync.GetIsAll() || m.quests == nil {
		m.quests = map[int32]*csmsg.RewardQuest{}
	}

	for _, v := range sync.Quests {
		m.quests[v.GetQuestID()] = v
	}
}

type ModuleWeapon struct {
	weapons map[uint32]*csmsg.Weapon
}

func (m *ModuleWeapon) OnModuleSync(r Robot, msg proto.Message) {
	sync := msg.(*csmsg.WeaponSyncToC)
	if sync.GetIsAll() || m.weapons == nil {
		m.weapons = map[uint32]*csmsg.Weapon{}
	}

	for _, v := range sync.Weapons {
		m.weapons[v.GetID()] = v
	}
}
