package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/node/client/login"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("usage addr userID\n")
		return
	}

	addr := os.Args[1]
	userID := os.Args[2]

	login.Login(userID, addr, func(lSession *login.Session, d *dispatcher.Dispatcher, msg *codecs.Message, err error) {
		if nil != err {
			fmt.Println(err)
		} else {

			data := msg.GetData().(*cs_msg.GameLoginToC)
			if data.GetIsFirstLogin() {
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				}
				lSession.Send(req)
				//cmd(lSession)
			} else {
				//cmd(lSession)
				timesQuery(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("------ CreateRoleToC", data)
				cmd(lSession)
			})
			d.Register(cmdEnum.CS_BaseSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.BaseSyncToC)
				fmt.Println("------ BaseSyncToC", data)
			})

			d.Register(cmdEnum.CS_QueryRoleInfo, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.QueryRoleInfoToC)
				fmt.Println("------ QueryRoleInfoToC", data)
				respTimes++
				checkResp()
				//cmd(lSession)
			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

var times int
var respTimes int
var startTime time.Time

func checkResp() {
	if respTimes == times {
		fmt.Println(time.Now().Sub(startTime))
		// 1000 7.597226436s
		// 1000 2.58s  优化后
	}
}

func timesQuery(lSession *login.Session) {
	fmt.Println(" input : times ")
	fmt.Scan(&times)
	respTimes = 0
	startTime = time.Now()

	for i := 0; i < times; i++ {
		req := &cs_msg.QueryRoleInfoToS{
			QueryArgs: []*cs_msg.QueryRoleInfoArg{
				{
					GameID:       proto.Uint64(100018),
					Base:         proto.Bool(true),
					AttrIDs:      []int32{1, 2, 3, 4},
					CharacterIDs: []int32{20},
					WeaponIDs:    []uint32{13},
				},
				{
					//GameID:       proto.Uint64(uint64(id)),
					UserID:       proto.String("ydd0026"),
					Base:         proto.Bool(true),
					AttrIDs:      []int32{1, 2, 3, 4},
					CharacterIDs: []int32{20},
					WeaponIDs:    []uint32{13},
				},
			},
		}
		lSession.Send(req)

	}
}

func cmd(lSession *login.Session) {
	go func() {
		time.Sleep(time.Second)
		fmt.Println("fight input : userID ")
		var id string
		fmt.Scan(&id)
		fmt.Println("fight input : gameID ")
		var gid uint64
		fmt.Scan(&gid)
		req := &cs_msg.QueryRoleInfoToS{
			QueryArgs: []*cs_msg.QueryRoleInfoArg{
				{
					GameID:       proto.Uint64(gid),
					Base:         proto.Bool(true),
					AttrIDs:      []int32{1, 2, 3, 4},
					CharacterIDs: []int32{20},
					WeaponIDs:    []uint32{13},
				},
				{
					//GameID:       proto.Uint64(uint64(id)),
					UserID:       proto.String(id),
					Base:         proto.Bool(true),
					AttrIDs:      []int32{1, 2, 3, 4},
					CharacterIDs: []int32{20},
					WeaponIDs:    []uint32{13},
				},
			},
		}
		lSession.Send(req)

	}()
}
