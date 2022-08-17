package maindungeons

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/client"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module"
	MainChapterTable "initialthree/node/table/excel/DataTable/MainChapter"
	"initialthree/node/table/excel/DataTable/MainDungeon"
	"initialthree/pkg/json"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"time"
)

const (
	tableName          = "main_dungeons"
	chapterFieldPrefix = "chapter"
	chapterFieldOffset = 0
	chapterFieldCount  = 2
	dungeonFieldPrefix = "dungeon"
	dungeonFieldOffset = chapterFieldOffset + chapterFieldCount
	dungeonFieldCount  = 8
)

var (
	tableFields []string
)

func initFields() {
	tableFields = make([]string, chapterFieldCount+dungeonFieldCount)

	// init chapter field name.
	for i := 0; i < chapterFieldCount; i++ {
		tableFields[chapterFieldOffset+i] = fmt.Sprintf("%s%d", chapterFieldPrefix, i)
	}

	// init dungeon field name.
	for i := 0; i < dungeonFieldCount; i++ {
		tableFields[dungeonFieldOffset+i] = fmt.Sprintf("%s%d", dungeonFieldPrefix, i)
	}
}

// 关卡数据存在即完成
type Dungeon struct {
	ID       int32 `json:"id,omitempty"` // ID
	Finished bool  `json:"finished,omitempty"`
}

// 章节数据存在即完成
type Chapter struct {
	ID       int32 `json:"id,omitempty"`
	Finished bool  `json:"finished,omitempty"`
}

func (m *MainDungeons) ChapterRange(cb func(c *Chapter) bool) {
	for _, slot := range m.chapterSlot {
		for _, c := range slot {
			if !cb(c) {
				return
			}
		}
	}
}

type MainDungeons struct {
	user         module.UserI
	chapterSlot  []map[int32]*Chapter
	dirtyChapter map[int32]*Chapter
	dungeonSlot  []map[int32]*Dungeon
	dirtyDungeon map[int32]*Dungeon
	*module.ModuleSaveBase
}

func (m *MainDungeons) chapterSlotIdx(id int32) int {
	return int(id) % chapterFieldCount
}
func (m *MainDungeons) dungeonSlotIdx(id int32) int {
	return int(id) % dungeonFieldCount
}

func (m *MainDungeons) setChapterDirty(c *Chapter, sync bool) {
	if sync {
		m.dirtyChapter[c.ID] = c
	}
	m.SetDirty(m.chapterSlotIdx(c.ID) + chapterFieldOffset)
}

func (m *MainDungeons) setDungeonDirty(d *Dungeon, sync bool) {
	if sync {
		m.dirtyDungeon[d.ID] = d
	}
	m.SetDirty(m.dungeonSlotIdx(d.ID) + dungeonFieldOffset)
}

func (m *MainDungeons) IsDungeonPass(id int32) bool {
	slotIdx := m.dungeonSlotIdx(id)
	slot := m.dungeonSlot[slotIdx]
	return slot[id] != nil
}

func (m *MainDungeons) DungeonPass(id int32) bool {
	slotIdx := m.dungeonSlotIdx(id)
	slot := m.dungeonSlot[slotIdx]

	if _, ok := slot[id]; ok {
		return false
	}

	d := &Dungeon{
		ID:       id,
		Finished: true,
	}
	slot[d.ID] = d
	m.setDungeonDirty(d, true)
	m.chapterSet(id)
	return true
}

func (m *MainDungeons) chapterSet(dungeonID int32) {
	def := MainDungeon.GetID(dungeonID)
	chapterID := def.ChapterID
	if chapterID != 0 {
		slotIdx := m.chapterSlotIdx(chapterID)
		slot := m.chapterSlot[slotIdx]

		if _, ok := slot[chapterID]; ok {
			return
		}

		// 不存在，判断章节关卡是否都已经完成
		chapterDef := MainChapterTable.GetID(chapterID)
		for _, d := range chapterDef.DungeonsArray {
			if !m.IsDungeonPass(d.ID) {
				return
			}
		}

		c := &Chapter{
			ID:       chapterID,
			Finished: true,
		}
		slot[c.ID] = c
		m.setChapterDirty(c, true)

		m.user.EmitEvent(event.EventMainChapter, chapterID)
	}

}

func (m *MainDungeons) ModuleType() module.ModuleType {
	return module.MainDungeons
}

func (m *MainDungeons) Init(fields map[string]*client.Field) error {
	var err error
	for i, fieldName := range tableFields {
		field, ok := fields[fieldName]
		if i < dungeonFieldOffset {
			ci := i - chapterFieldOffset
			m.chapterSlot[ci] = map[int32]*Chapter{}
			if !ok || len(field.GetBlob()) == 0 {
				m.SetDirty(i)
			} else {
				if err = json.Unmarshal(field.GetBlob(), &m.chapterSlot[ci]); err != nil {
					return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
				}
			}
		} else {
			di := i - dungeonFieldOffset
			m.dungeonSlot[di] = map[int32]*Dungeon{}
			if !ok || len(field.GetBlob()) == 0 {
				m.SetDirty(i)
			} else {
				if err = json.Unmarshal(field.GetBlob(), &m.dungeonSlot[di]); err != nil {
					return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
				}
			}
		}
	}
	return nil
}

func (m *MainDungeons) ReadOut() *module.ReadOutCommand {
	cmd := &module.ReadOutCommand{
		Table:  tableName,
		Key:    m.user.GetIDStr(),
		Fields: tableFields,
		Module: m,
	}
	return cmd
}

func (m *MainDungeons) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  tableName,
		Key:    m.user.GetIDStr(),
		Module: m,
	}

	wbFields := make([]*module.WriteBackFiled, 0, len(fields))

	for f, _ := range fields {
		switch ff := f.(type) {
		case int:
			var fieldName string
			var bytes []byte

			if ff >= chapterFieldOffset && ff < dungeonFieldOffset {
				fieldName = tableFields[ff]
				ci := ff - chapterFieldOffset
				bytes, _ = json.Marshal(&m.chapterSlot[ci])

			} else if ff >= dungeonFieldOffset && ff < dungeonFieldOffset+dungeonFieldCount {
				fieldName = tableFields[ff]
				di := ff - dungeonFieldOffset
				bytes, _ = json.Marshal(&m.dungeonSlot[di])

			} else {
				zaplogger.GetSugar().Errorf("MainDungeons: user(%s, %d) do not implemented field_index %d", m.user.GetUserID(), m.user.GetID(), ff)
				continue
			}

			wbFields = append(wbFields, &module.WriteBackFiled{
				Name:  fieldName,
				Value: bytes,
			})

		default:
			panic("not implemented")
		}
	}

	cmd.Fields = wbFields

	return cmd
}

func (m *MainDungeons) FlushDirtyToClient() {
	if len(m.dirtyChapter) == 0 && len(m.dirtyDungeon) == 0 {
		return
	}

	msg := &cs_message.MainDungeonsSyncToC{
		All: proto.Bool(false),
	}

	chapters := make([]*cs_message.MainChapter, 0, len(m.dirtyChapter))
	for _, c := range m.dirtyChapter {
		chapters = append(chapters, &cs_message.MainChapter{
			Id:       proto.Int32(c.ID),
			Finished: proto.Bool(c.Finished),
		})
	}
	msg.Chapters = chapters
	m.dirtyChapter = make(map[int32]*Chapter)

	dungeons := make([]*cs_message.MainDungeon, 0, len(m.dirtyDungeon))
	for _, d := range m.dirtyDungeon {
		dungeons = append(dungeons, &cs_message.MainDungeon{
			Id:       proto.Int32(d.ID),
			Finished: proto.Bool(d.Finished),
		})
	}
	msg.Dungeons = dungeons
	m.dirtyDungeon = make(map[int32]*Dungeon)

	m.user.Post(msg)
}

func (m *MainDungeons) FlushAllToClient(seqNo ...uint32) {
	msg := &cs_message.MainDungeonsSyncToC{
		All: proto.Bool(true),
	}

	chapters := make([]*cs_message.MainChapter, 0, len(m.chapterSlot)*2)
	for _, slot := range m.chapterSlot {
		for _, c := range slot {
			chapters = append(chapters, &cs_message.MainChapter{
				Id:       proto.Int32(c.ID),
				Finished: proto.Bool(c.Finished),
			})
		}
	}
	msg.Chapters = chapters

	dungeons := make([]*cs_message.MainDungeon, 0, len(m.dungeonSlot)*2)
	for _, slot := range m.dungeonSlot {
		for _, d := range slot {
			dungeons = append(dungeons, &cs_message.MainDungeon{
				Id:       proto.Int32(d.ID),
				Finished: proto.Bool(d.Finished),
			})
		}
	}
	msg.Dungeons = dungeons

	seq := uint32(0)
	if len(seqNo) > 0 {
		seq = seqNo[0]
	}
	m.user.Reply(seq, msg)

	m.dirtyChapter = make(map[int32]*Chapter)
	m.dirtyDungeon = make(map[int32]*Dungeon)
}

func (m *MainDungeons) Tick(now time.Time) {}

func init() {
	initFields()

	module.RegisterModule(module.MainDungeons, func(user module.UserI) module.ModuleI {
		m := &MainDungeons{
			user:         user,
			chapterSlot:  make([]map[int32]*Chapter, chapterFieldCount),
			dirtyChapter: map[int32]*Chapter{},
			dungeonSlot:  make([]map[int32]*Dungeon, dungeonFieldCount),
			dirtyDungeon: map[int32]*Dungeon{},
		}

		m.ModuleSaveBase = module.NewModuleSaveBase(m)

		return m
	})
}
