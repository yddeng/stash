package main

import (
	"bufio"
	"fmt"
	codecs "initialthree/codec/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/node/client/login"
	"initialthree/protocol/cmdEnum"
	"io"
	"os"
	"strings"

	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"

	cs_msg "initialthree/protocol/cs/message"
)

type backpackItem struct {
	id          uint32
	tid         int32
	count       int32
	acquireTime int64
	timeLimit   int64
}

var (
	stopCh        = make(chan struct{})
	nextCH        = make(chan struct{}, 1)
	backpackItems = make(map[uint32]*backpackItem)
	sess          *login.Session
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: addr userID\n")
		return
	}

	addr := os.Args[1]
	userID := os.Args[2]

	login.Login(userID, addr, func(lSession *login.Session, d *dispatcher.Dispatcher, msg *codecs.Message, err error) {
		if nil != err {
			fmt.Println(err)
			close(stopCh)
		} else {
			sess = lSession

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)

				fmt.Println("create role successfully:", data)

				// lSession.Send(&cs_msg.BackpackSyncToS{})

				nextCH <- struct{}{}
				go loop()

			})

			d.Register(cmdEnum.CS_GameMaster, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() != 0 {
					fmt.Println("add item failed:", cs_msg.ErrCode(msg.GetErrCode()).String())
				}

				nextCH <- struct{}{}
			})

			d.Register(cmdEnum.CS_BackpackSell, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() != 0 {
					fmt.Println("sell item failed:", cs_msg.ErrCode(msg.GetErrCode()).String())
				}

				nextCH <- struct{}{}
			})

			d.Register(cmdEnum.CS_BackpackUse, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() != 0 {
					fmt.Println("use item failed", cs_msg.ErrCode(msg.GetErrCode()).String())
				}

				nextCH <- struct{}{}
			})

			d.Register(cmdEnum.CS_BackpackRem, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() != 0 {
					fmt.Println("rem entity failed", cs_msg.ErrCode(msg.GetErrCode()).String())
				}

				nextCH <- struct{}{}
			})

			d.Register(cmdEnum.CS_AttrSync, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() == 0 {
					resp := msg.GetData().(*cs_msg.AttrSyncToC)
					fmt.Print("attr update:")
					for _, attr := range resp.Attrs {
						fmt.Printf(" {%d:%d}", attr.GetId(), attr.GetVal())
					}
					fmt.Println()
				}
			})

			d.Register(cmdEnum.CS_BackpackSync, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() == 0 {
					resp := msg.GetData().(*cs_msg.BackpackSyncToC)

					if resp.GetAll() {
						backpackItems = make(map[uint32]*backpackItem)

						for _, e := range resp.Entities {
							switch e.GetType() {
							case 1: // item
								backpackItems[e.GetId()] = &backpackItem{
									id:          e.GetId(),
									tid:         e.GetContent().GetItem().GetTid(),
									count:       e.GetCount(),
									acquireTime: e.GetAcquireTime(),
									timeLimit:   e.GetContent().GetItem().GetTimeLimit(),
								}
							}
						}

						fmt.Printf(
							`recv all backpack entities:
%s
`,
							formatBackpack(),
						)

						nextCH <- struct{}{}
						go loop()
					} else {
						fmt.Println("backpack update:")
						for _, e := range resp.Entities {
							switch e.GetType() {
							case 1:
								if e.GetCount() == 0 {
									delete(backpackItems, e.GetId())
									fmt.Printf("\titem %d removed\n", e.GetId())
								} else {
									item := backpackItems[e.GetId()]
									if item != nil {
										item.count = e.GetCount()
									} else {
										item = &backpackItem{
											id:          e.GetId(),
											tid:         e.GetContent().GetItem().GetTid(),
											count:       e.GetCount(),
											acquireTime: e.GetAcquireTime(),
											timeLimit:   e.GetContent().GetItem().GetTimeLimit(),
										}
										backpackItems[e.GetId()] = item
									}
									fmt.Printf("\titem(%d, %d, %d, %d) updated\n", item.id, item.tid, item.count, item.timeLimit)
								}
							}
						}
						fmt.Println()
					}
				}
			})

			if cs_msg.ErrCode(msg.GetErrCode()) == cs_msg.ErrCode_OK {
				data := msg.GetData().(*cs_msg.GameLoginToC)
				fmt.Println("game login:", data)

				if data.GetIsFirstLogin() {
					//fmt.Println("input role name ")
					req := &cs_msg.CreateRoleToS{
						Name: proto.String(userID),
					}
					lSession.Send(req)
				} else {
					nextCH <- struct{}{}
					go loop()
				}
			} else {
				fmt.Println("login failed", cs_msg.ErrCode(msg.GetErrCode()).String())
				close(stopCh)
			}
		}
	})

	_, _ = <-stopCh
}

func loop() {

	var (
		cmdDetail = `0: exit.
1: add item.
2: sell entity.
3: use entity.
4: rem entity.
5: show data.`

		cmdArgsDetail = map[int64]struct {
			detail string
			args   int
		}{
			0: {},
			1: {detail: "[item tid][count]", args: 2},
			2: {detail: "[entity id][count]", args: 2},
			3: {detail: "[entity id][count]", args: 2},
			4: {detail: "[entity id][count][type]", args: 3},
			5: {},
		}

		stop = false

		inputReader = bufio.NewReader(os.Stdin)

		waitNext = true

		cmd         int64
		inputString string
		args        []interface{}
		err         error
	)

loop:
	for !stop {

		if waitNext {
			_, _ = <-nextCH
			waitNext = false
		}

		fmt.Printf("input cmd (? for help): ")

		inputString, err = inputReader.ReadString('\n')
		if err != nil && err != io.EOF {
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

		argsDetail, ok := cmdArgsDetail[cmd]
		if !ok {
			fmt.Println("cmd error")
			continue
		}

		if argsDetail.args > 0 {
			fmt.Printf("input args %s: ", argsDetail.detail)

			inputString, err = inputReader.ReadString('\n')
			if err != nil && err != io.EOF {
				fmt.Println("input error:", err)
				continue
			}

			args = make([]interface{}, argsDetail.args)
			for i := 0; i < argsDetail.args; i++ {
				args[i] = new(int32)
			}
			if n, err := fmt.Sscan(inputString, args...); n != argsDetail.args || err != nil {
				fmt.Println("parse args error:", err)
				continue loop
			}
		}

		switch cmd {
		case 0:
			stop = true
			break

		case 1:
			tid := *args[0].(*int32)
			count := *args[1].(*int32)
			if count <= 0 {
				fmt.Println("count error")
				continue
			}

			msg := &cs_msg.GameMasterToS{
				Cmds: []*cs_msg.GmCmd{
					{
						Type:  proto.Int32(3),
						ID:    proto.Int32(tid),
						Count: proto.Int32(count),
					},
				},
			}
			sess.Send(msg)
			waitNext = true

		case 2:
			id := uint32(*args[0].(*int32))
			count := *args[1].(*int32)

			if count <= 0 {
				fmt.Println("count error")
				continue
			}

			msg := &cs_msg.BackpackSellToS{}
			msg.SellEntities = []*cs_msg.BackpackSellEntity{&cs_msg.BackpackSellEntity{
				Id:    proto.Uint32(id),
				Count: proto.Int32(count),
			}}
			sess.Send(msg)
			waitNext = true

		case 3:
			id := uint32(*args[0].(*int32))
			count := *args[1].(*int32)

			if count <= 0 {
				fmt.Println("count error")
				continue
			}

			msg := &cs_msg.BackpackUseToS{
				Id:    proto.Uint32(id),
				Count: proto.Int32(count),
			}
			sess.Send(msg)
			waitNext = true

		case 4:
			id := uint32(*args[0].(*int32))
			count := *args[1].(*int32)
			typ := *args[2].(*int32)

			if count <= 0 {
				fmt.Println("count error")
				continue
			}

			msg := &cs_msg.BackpackRemToS{
				RemEntities: []*cs_msg.BackpackRemEntity{
					&cs_msg.BackpackRemEntity{
						Type:  proto.Int32(typ),
						Id:    proto.Uint32(id),
						Count: proto.Int32(count),
					},
				},
			}
			sess.Send(msg)
			waitNext = true

		case 5:
			fmt.Printf("backpack:\n%s\n", formatBackpack())
		}
	}

	fmt.Println("exit")
	close(stopCh)
}

func formatBackpack() string {
	s := "items:\n"

	for _, it := range backpackItems {
		s += fmt.Sprintf("\t{id:%d, tid:%d, count:%d, tl:%d}\n", it.id, it.tid, it.count, it.timeLimit)
	}
	return s
}
