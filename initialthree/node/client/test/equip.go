package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
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
			data := msg.GetData().(*cs_msg.GameLoginToC)
			if data.GetIsFirstLogin() {
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				}
				lSession.Send(req)
			} else {
				Equip(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
				Equip(lSession)
			})
			d.Register(cmdEnum.CS_EquipLock, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.EquipLockToC)
				fmt.Println("EquipLockToC", data)
			})
			d.Register(cmdEnum.CS_EquipSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.EquipSyncToC)
				fmt.Println("EquipSyncToC", data)
			})
			d.Register(cmdEnum.CS_EquipEquip, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.EquipEquipToC)
				fmt.Println("EquipEquipToC", data)
			})

			d.Register(cmdEnum.CS_EquipDemount, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.EquipDemountToC)
				fmt.Println("EquipDemountToC", data)
			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func Equip(lSession *login.Session) {
	go func() {

		fmt.Println("1:equip 2:demount")
		var key int
		fmt.Scan(&key)
		if key == 1 {
			fmt.Println("equip : characterId equipId ")
			var cid, eId int
			fmt.Scan(&cid, &eId)
			req := &cs_msg.EquipEquipToS{
				CharacterID: proto.Int32(int32(cid)),
				EquipID:     proto.Uint32(uint32(eId)),
			}
			lSession.Send(req)
		} else if key == 2 {
			fmt.Println("demount : equipId ")
			var eId int
			fmt.Scan(&eId)
			req := &cs_msg.EquipDemountToS{
				EquipID: []uint32{uint32(eId)},
			}
			_ = lSession.Send(req)
		}
	}()
}
