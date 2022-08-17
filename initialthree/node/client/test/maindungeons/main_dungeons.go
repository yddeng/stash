package main

import (
	"bufio"
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/node/client/dispatcher"
	"initialthree/node/client/login"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/config"
	"initialthree/node/table/excel"
	"initialthree/protocol/cmdEnum"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	codecs "initialthree/codec/cs"
	cs_message "initialthree/protocol/cs/message"

	ChapterTable "initialthree/node/table/excel/DataTable/MainChapter"
	DungeonTable "initialthree/node/table/excel/DataTable/MainDungeon"
)

type chapter struct {
	ID        int32
	AwardFlag []bool
}

func (c *chapter) String() string {
	starCount := 0

	chapterCfg := ChapterTable.GetID(c.ID)

	for _, v := range chapterCfg.DungeonsArray {
		dungeon := dungeons[v.ID]
		if dungeon != nil {
			starCount += dungeon.StarCount()
		}
	}

	return fmt.Sprintf("id:%d, awardFlag:%v starCount:%d", c.ID, c.AwardFlag, starCount)
}

type dungeon struct {
	ID    int32
	Stars []bool
	Count int32
}

func (d *dungeon) StarCount() int {
	n := 0
	for _, v := range d.Stars {
		if v {
			n++
		}
	}
	return n
}

func (d *dungeon) String() string {
	dungeonCfg := DungeonTable.GetID(d.ID)

	cid := dungeonCfg.ChapterID

	return fmt.Sprintf("id:%d, cid:%d, stars:%v, count:%d", d.ID, cid, d.Stars, d.Count)
}

type attr struct {
	ID  int32
	Val int64
}

func (a *attr) String() string {
	return fmt.Sprintf("%s(%d): %d", attr2.GetNameById(a.ID), a.ID, a.Val)
}

var (
	Session *login.Session

	mtx      sync.RWMutex
	chapters map[int32]*chapter = make(map[int32]*chapter)
	dungeons map[int32]*dungeon = make(map[int32]*dungeon)
	attrs    map[int32]*attr    = make(map[int32]*attr)

	sigStop = make(chan bool)
	nextCH  = make(chan struct{}, 1)
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: addr userID\n")
		return
	}

	cfg, err := config.LoadConfig("/Users/godyy/go/src/initialthree/node/config.toml")
	if err != nil {
		log.Fatalln(err)
	}

	excel.Load(cfg.Common.GetExcelPath())

	userID := os.Args[2]
	addr := os.Args[1]

	login.Login(userID, addr, func(lSession *login.Session, d *dispatcher.Dispatcher, msg *codecs.Message, err error) {
		Session = lSession

		if nil != err {
			fmt.Println(err)
		} else {
			data := msg.GetData().(*cs_message.GameLoginToC)
			if data.GetIsFirstLogin() {
				_ = lSession.Send(&cs_message.CreateRoleToS{
					Name: proto.String(userID),
				})
			} else {
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_message.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
			})

			d.Register(cmdEnum.CS_GameMaster, func(session *fnet.Socket, msg *codecs.Message) {

				fmt.Println("GameMaster", cs_message.ErrCode(msg.GetErrCode()))
				nextCH <- struct{}{}
			})

			d.Register(cmdEnum.CS_AttrSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_message.AttrSyncToC)

				mtx.Lock()

				if data.GetIsAll() {
					for _, v := range data.Attrs {
						attrs[v.GetId()] = &attr{
							ID:  v.GetId(),
							Val: v.GetVal(),
						}
					}
				} else {
					fmt.Println("attr update:")
					for _, v := range data.Attrs {
						attrs[v.GetId()] = &attr{
							ID:  v.GetId(),
							Val: v.GetVal(),
						}
						fmt.Println("\t", attrs[v.GetId()].String())
					}
				}

				mtx.Unlock()
			})

			d.Register(cmdEnum.CS_MainDungeonsSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_message.MainDungeonsSyncToC)

				mtx.Lock()

				if !data.GetAll() {
					fmt.Println("chapter update:")
				}
				for _, v := range data.Chapters {
					c := chapters[v.GetId()]
					if c == nil {
						c = &chapter{
							ID: v.GetId(),
						}
						chapters[v.GetId()] = c
					}

					c.AwardFlag = v.GetAwardFlag()

					if !data.GetAll() {
						fmt.Println("\t", c.String())
					}
				}

				if !data.GetAll() {
					fmt.Println("dungeon update:")
				}
				for _, v := range data.Dungeons {
					d := dungeons[v.GetId()]
					if d == nil {
						d = &dungeon{
							ID: v.GetId(),
						}
						dungeons[v.GetId()] = d
					}

					d.Stars = v.Stars
					d.Count = v.GetRemainCount()

					if !data.GetAll() {
						fmt.Println("\t", d.String())
					}
				}

				mtx.Unlock()

				if data.GetAll() {
					go loop()
					nextCH <- struct{}{}
				}
			})

			d.Register(cmdEnum.CS_MainDungeonsGetChapterStarAward, func(session *login.Session, msg *codecs.Message) {

				if cs_message.ErrCode(msg.GetErrCode()) {
					data := msg.GetData().(*cs_message.MainDungeonsGetChapterStarAwardToC)
					fmt.Println("chapter award:", *data)
				}
				nextCH <- struct{}{}
			})
		}
	})

	_, _ = <-sigStop
}

func loop() {
	go func() {

		var (
			cmdDetail = `0: exit
1: add attr
2: dungeon pass
3: get chapter start award
4: show data`

			cmdArgsDetail = map[int64]struct {
				detail string
				args   int
			}{
				0: {},
				1: {detail: "[attr][count]", args: 2},
				2: {detail: "[dungeon id]", args: 1},
				3: {detail: "[chapter id][award No]", args: 2},
				4: {},
			}

			stop     = false
			waitNext = true

			inputReader = bufio.NewReader(os.Stdin)
			cmd         int64
			inputString string
			args        []interface{}
			err         error
		)

		for !stop {

			if waitNext {
				_, _ = <-nextCH
				waitNext = false
			}

			fmt.Printf("input cmd (? for help): ")

			if inputString, err = inputReader.ReadString('\n'); err != nil && err != io.EOF {
				fmt.Println("input error:", err)
				continue
			}

			if strings.HasPrefix(inputString, "?") {
				fmt.Println(cmdDetail)
				continue
			}

			if _, err = fmt.Sscanf(inputString, "%d\n", &cmd); err != nil {
				fmt.Println("parse cmd error:", err)
				continue
			}

			if d, ok := cmdArgsDetail[cmd]; !ok {
				fmt.Println("cmd error")
				continue
			} else if d.args > 0 {
				args = make([]interface{}, d.args)
				for i := 0; i < d.args; i++ {
					args[i] = new(int)
				}

				fmt.Printf("input args %s: ", d.detail)
				if inputString, err = inputReader.ReadString('\n'); err != nil && err != io.EOF {
					fmt.Println("input error:", err)
					continue
				}

				if n, err := fmt.Sscan(inputString, args...); n != d.args || err != nil {
					fmt.Println("parse args error:", err)
					continue
				}
			}

			switch cmd {
			case 0:
				stop = true
				continue

			case 1:
				attr := int32(*args[0].(*int))
				count := int32(*args[1].(*int))

				msg := &cs_message.GameMasterToS{
					Cmds: []*cs_message.GmCmd{
						{
							Type:  proto.Int32(2),
							ID:    proto.Int32(attr),
							Count: proto.Int32(count),
						},
					},
				}

				Session.Send(msg)
				waitNext = true

			case 2:
				dungeonID := int32(*args[0].(*int))

				msg := &cs_message.GameMasterToS{
					Cmds: []*cs_message.GmCmd{
						{
							Type: proto.Int32(4),
							ID:   proto.Int32(dungeonID),
						},
					},
				}

				Session.Send(msg)
				waitNext = true

			case 3:
				chapterID := int32(*args[0].(*int))
				awardNo := int32(*args[1].(*int))

				msg := &cs_message.MainDungeonsGetChapterStarAwardToS{
					ChapterID: proto.Int32(chapterID),
					AwardNo:   proto.Int32(awardNo),
				}

				Session.Send(msg)
				waitNext = true

			case 4:
				mtx.RLock()

				fmt.Println("attrs:")
				for _, v := range attrs {
					fmt.Println("\t", v.String())
				}

				fmt.Println("chapters:")
				for _, v := range chapters {
					fmt.Println("\t", v.String())
				}

				fmt.Println("dungeons:")
				for _, v := range dungeons {
					fmt.Println("\t", v.String())
				}

				mtx.RUnlock()
			}
		}

		fmt.Println("exit")
		close(sigStop)
	}()
}
