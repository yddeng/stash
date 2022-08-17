package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/node/client/login"
	table "initialthree/node/table/excel"
	"initialthree/node/table/excel/DataTable/ScarsIngrainBossInstance"
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
	table.Load("/Users/yidongdeng/svn/TheInitial3_Dev/TheInitial3_Dev_Config/Excel")

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
				cmd(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("------ CreateRoleToC", data)
				cmd(lSession)
			})
			d.Register(cmdEnum.CS_RankGetTopList, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.RankGetTopListToC)
				fmt.Println("------ RankGetTopListToC", data)
			})
			d.Register(cmdEnum.CS_ScarsIngrainSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.ScarsIngrainSyncToC)
				fmt.Println("------ ScarsIngrainSyncToC", data)
				if data.GetCurrentRankId() != 0 {
					_ = lSession.Send(&cs_msg.RankGetTopListToS{
						Version: proto.Int32(0),
						RankID:  proto.Int32(data.GetCurrentRankId()),
						GetLast: proto.Bool(false),
					})
					_ = lSession.Send(&cs_msg.RankGetRankToS{
						RankID: proto.Int32(data.GetCurrentRankId()),
					})
				}
			})
			d.Register(cmdEnum.CS_CharacterSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CharacterSyncToC)
				fmt.Println(" ------ CharacterSyncToC", data)
			})
			d.Register(cmdEnum.CS_RankGetRank, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.RankGetRankToC)
				fmt.Println(" ------ RankGetRankToC", data)
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
			d.Register(cmdEnum.CS_ScarsIngrainSaveScore, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.ScarsIngrainSaveScoreToC)
				fmt.Println("------ ScarsIngrainSaveScoreToC", data)
				cmd(lSession)
			})

			d.Register(cmdEnum.CS_ScarsIngrainGetScoreAward, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.ScarsIngrainGetScoreAwardToC)
				fmt.Println("------ ScarsIngrainGetScoreAwardToC", data)
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
		fmt.Println("1:fight 2:fightend 3:saveScore 4:getAward ")
		var key int
		fmt.Scan(&key)
		if key == 1 {
			fmt.Println("fight input : bossId ,difficult ")
			var id, diff int32
			fmt.Scan(&id, &diff)
			b := ScarsIngrainBossInstance.GetID(id*10 + diff)
			fmt.Println(b)
			fmt.Println("fight team -> (1,0,0) ", "DungeonID ID", b.DungeonID)
			req := &cs_msg.LevelFightToS{
				DungeonID:     proto.Int32(b.DungeonID),
				CharacterTeam: &cs_msg.CharacterTeam{CharacterList: []int32{10, 0, 0}},
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
					FightID:   proto.Int64(int64(id)),
					UseTime:   proto.Int32(t),
					BossCurHP: proto.Float64(float64(cb)),
					BossMaxHP: proto.Float64(float64(mb)),
					BeHit:     proto.Int32(bh),
					Pass:      proto.Bool(true),
				}
				lSession.Send(req)
			} else {
				req := &cs_msg.LevelFightEndToS{
					FightID: proto.Int64(int64(id)),
					Pass:    proto.Bool(false),
				}
				lSession.Send(req)
			}

		} else if key == 3 {
			fmt.Println("saveScore input : ScoreID,Save ")
			var id, score int
			fmt.Scan(&id, &score)
			req := &cs_msg.ScarsIngrainSaveScoreToS{
				ScoreID: proto.Int32(int32(id)),
				IsSave:  proto.Bool(score == 1),
			}
			lSession.Send(req)
		} else if key == 4 {
			fmt.Println("getAward input : Score ")
			var id int32
			fmt.Scan(&id)
			req := &cs_msg.ScarsIngrainGetScoreAwardToS{
				Score: proto.Int32(id),
			}
			lSession.Send(req)
		}
	}()
}
