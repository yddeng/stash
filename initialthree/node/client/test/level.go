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
			} else {
				cmd(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("------ CreateRoleToC", data)
				cmd(lSession)
			})

			d.Register(cmdEnum.CS_MaterialDungeonSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.MaterialDungeonSyncToC)
				fmt.Println("------ MaterialDungeonSyncToC", data)

			})
			d.Register(cmdEnum.CS_AttrSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.AttrSyncToC)
				fmt.Println("------ AttrSyncToC", data)

			})

			d.Register(cmdEnum.CS_LevelFight, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.LevelFightToC)
				fmt.Println("------  LevelFightToC", data)
				cmd(lSession)
			})
			d.Register(cmdEnum.CS_LevelFightEnd, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.LevelFightEndToC)
				fmt.Println(" ------ LevelFightEndToC", data)
				cmd(lSession)
			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func cmd(lSession *login.Session) {
	go func() {
		time.Sleep(time.Second)
		fmt.Println("1:fight 2:fightend")
		var key int
		fmt.Scan(&key)
		if key == 1 {
			fmt.Println("fight input : id ")
			var id int32
			fmt.Scan(&id)
			req := &cs_msg.LevelFightToS{
				DungeonID:     proto.Int32(id),
				CharacterTeam: &cs_msg.CharacterTeam{CharacterList: []int32{20, 0, 0}},
				Multiple:      proto.Int32(1),
			}
			lSession.Send(req)
		} else if key == 2 {
			fmt.Println("fightend input 1: FightID,pass ")
			var id, p int
			fmt.Scan(&id, &p)
			if p == 1 {
				fmt.Println("fightend input 2: UseTime, BossCurHp, BossMaxHp, behit")
				var t, cb, mb, bh int32
				fmt.Scan(&t, &cb, &mb, &bh)
				req := &cs_msg.LevelFightEndToS{
					FightID:       proto.Int64(int64(id)),
					UseTime:       proto.Int32(t),
					BossCurHP:     proto.Float64(float64(cb)),
					BossMaxHP:     proto.Float64(float64(mb)),
					BeHit:         proto.Int32(bh),
					Pass:          proto.Bool(true),
					KillBossCount: proto.Int32(5),
				}
				lSession.Send(req)
			} else {
				req := &cs_msg.LevelFightEndToS{
					FightID: proto.Int64(int64(id)),
					Pass:    proto.Bool(false),
				}
				lSession.Send(req)
			}

		}
	}()
}
