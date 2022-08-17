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

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
				go GM(lSession)
			})
			d.Register(cmdEnum.CS_ServerTime, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.ServerTimeToC)
				fmt.Println("ServerTimeToC", data)
			})
			d.Register(cmdEnum.CS_AttrSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.AttrSyncToC)
				fmt.Println("AttrSyncToC", data)
			})
			d.Register(cmdEnum.CS_DrawCardSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.DrawCardSyncToC)
				fmt.Println("DrawCardSyncToC", data)
			})
			d.Register(cmdEnum.CS_BattleAttrSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.BattleAttrSyncToC)
				fmt.Println("BattleAttrSyncToC", data)
			})
			d.Register(cmdEnum.CS_BackpackSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.BackpackSyncToC)
				fmt.Println("BackpackSyncToC", data)
			})
			d.Register(cmdEnum.CS_EquipSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.EquipSyncToC)
				fmt.Println("EquipSyncToC", data)
			})
			d.Register(cmdEnum.CS_WeaponSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.WeaponSyncToC)
				fmt.Println("WeaponSyncToC", data)
			})
			d.Register(cmdEnum.CS_QuestSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.QuestSyncToC)
				fmt.Println("QuestSyncToC", data)
			})
			d.Register(cmdEnum.CS_LevelFight, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.LevelFightToC)
				fmt.Println("LevelFightToC", data)
			})
			d.Register(cmdEnum.CS_GameMaster, func(session *fnet.Socket, msg *codecs.Message) {

				fmt.Println("CS_GameMaster ErrCode", cs_msg.ErrCode(msg.GetErrCode()).String())
			})
			d.Register(cmdEnum.CS_BaseSync, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.BaseSyncToC)
				fmt.Println("BaseSyncToC", data)
			})
			d.Register(cmdEnum.CS_AssetSync, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.AssetSyncToC)
				fmt.Println("AssetSyncToC", data)
			})

			d.Register(cmdEnum.CS_MainDungeonsSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.MainDungeonsSyncToC)
				fmt.Println("MainDungeonsSyncToC", data)
			})

			data := msg.GetData().(*cs_msg.GameLoginToC)
			fmt.Println("GameLoginToC", data)
			if data.GetIsFirstLogin() {
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				}

				lSession.Send(req)
			} else {
				go GM(lSession)
			}
		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func GM(lSession *login.Session) {
	for {
		fmt.Printf("----------------------type -->  attr:1 character:2 item:3 主线关卡:4 装备:5\n")

		fmt.Printf("添加type:")
		var tt int
		fmt.Scan(&tt)
		fmt.Printf("添加ID:")
		var cid int
		fmt.Scan(&cid)
		fmt.Printf("添加数量:")
		var count int
		fmt.Scan(&count)

		req := &cs_msg.GameMasterToS{
			Cmds: []*cs_msg.GmCmd{{
				Type:  proto.Int32(int32(tt)),
				ID:    proto.Int32(int32(cid)),
				Count: proto.Int32(int32(count)),
			}},
		}
		lSession.Send(req)
	}
}
