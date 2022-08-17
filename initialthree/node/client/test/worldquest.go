package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/node/client/login"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"os"
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

			d.Register(cmdEnum.CS_CreateRole, func(session *net.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
			})

			d.Register(cmdEnum.CS_AttrSync, func(session *net.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.AttrSyncToC)
				fmt.Println("AttrSyncToC", data)
			})

			d.Register(cmdEnum.CS_ReputationSync, func(session *net.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.ReputationSyncToC)
				fmt.Println("ReputationSyncToC", data)
			})
			d.Register(cmdEnum.CS_WorldQuestSync, func(session *net.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.WorldQuestSyncToC)
				fmt.Println("WorldQuestSyncToC", data)
			})

			d.Register(cmdEnum.CS_WorldQuestRefresh, func(session *net.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.WorldQuestRefreshToC)
				fmt.Println("WorldQuestRefreshToC", data)
			})

			d.Register(cmdEnum.CS_LevelFight, func(session *net.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.LevelFightToC)
				fmt.Println("------  LevelFightToC", data)

				req := &cs_msg.LevelFightEndToS{
					FightID:   proto.Int64(int64(data.GetFightID())),
					UseTime:   proto.Int32(10),
					BossCurHP: proto.Float64(float64(0)),
					BossMaxHP: proto.Float64(float64(100)),
					BeHit:     proto.Int32(11),
					Pass:      proto.Bool(true),
				}
				lSession.Send(req)
			})
			d.Register(cmdEnum.CS_LevelFightEnd, func(session *net.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.LevelFightEndToC)
				fmt.Println(" ------ LevelFightEndToC", data)

			})

			data := msg.GetData().(*cs_msg.GameLoginToC)
			if data.GetIsFirstLogin() {
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				}

				lSession.Send(req)
			} else {
				go world(lSession)
			}
		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func world(session *login.Session) {
	fmt.Println("refresh 1, fight 2")
	var k int
	fmt.Scan(&k)
	if k == 1 {
		req := &cs_msg.WorldQuestRefreshToS{}
		session.Send(req)
	} else {
		fmt.Println("fight id ")
		var id int32
		fmt.Scan(&id)
		req := &cs_msg.LevelFightToS{
			DungeonID:     proto.Int32(id),
			CharacterTeam: &cs_msg.CharacterTeam{CharacterList: []int32{2, 0, 0}},
		}
		session.Send(req)
	}
}
