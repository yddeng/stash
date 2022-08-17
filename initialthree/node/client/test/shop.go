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
				getShopSync(lSession)
			})

			d.Register(cmdEnum.CS_AttrSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.AttrSyncToC)
				fmt.Println("AttrSyncToC", data)
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
			d.Register(cmdEnum.CS_ShopBuy, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.ShopBuyToC)
				fmt.Println("ShopBuy", data)
				shop(lSession)
			})
			d.Register(cmdEnum.CS_ShopRefresh, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.ShopRefreshToC)
				fmt.Println("ShopRefreshToC", data)
				shop(lSession)
			})
			d.Register(cmdEnum.CS_ShopSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.ShopSyncToC)
				fmt.Println("ShopSyncToC", data)

				if data.GetIsAll() {
					shop(lSession)
				}
			})
			d.Register(cmdEnum.CS_AssetSync, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.AssetSyncToC)
				fmt.Println("AssetSyncToC", data)
			})
			d.Register(cmdEnum.CS_ShopPay, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.ShopPayToC)
				fmt.Println("ShopPayToC", data)
				shop(lSession)
			})

			data := msg.GetData().(*cs_msg.GameLoginToC)
			if data.GetIsFirstLogin() {
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				}

				lSession.Send(req)
			} else {
				getShopSync(lSession)
			}
		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func getShopSync(lSession *login.Session) {
	lSession.Send(&cs_msg.ShopSyncToS{})
}

func shop(lSession *login.Session) {
	fmt.Printf("----------------------type -->  buy:1 refresh:2 pay:3 \n")
	go func() {
		var key int
		fmt.Scan(&key)
		if key == 1 {
			fmt.Println("input : pid ,count ")
			var id, diff int32
			fmt.Scan(&id, &diff)

			msg := &cs_msg.ShopBuyToS{
				Id:    proto.Int32(id),
				Count: proto.Int32(diff),
			}
			lSession.Send(msg)
		} else if key == 2 {
			fmt.Println("input : shopID ")
			var id int32
			fmt.Scan(&id)

			msg := &cs_msg.ShopRefreshToS{
				ShopID: proto.Int32(id),
			}
			lSession.Send(msg)
		} else if key == 3 {
			fmt.Println("input : payID ")
			var id int32
			fmt.Scan(&id)

			msg := &cs_msg.ShopPayToS{
				PayID: proto.Int32(id),
			}
			lSession.Send(msg)
		}

	}()

}
