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

			data := msg.GetData().(*cs_msg.GameLoginToC)
			if data.GetIsFirstLogin() {
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				}
				lSession.Send(req)
			} else {
				assets(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
				assets(lSession)
			})
			d.Register(cmdEnum.CS_Assets, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.AssetsToC)
				fmt.Println("AssetsToC", data)
			})

			d.Register(cmdEnum.CS_GameMaster, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() != 0 {
					fmt.Println("CS_GameMaster ErrCode", cs_msg.ErrCode(msg.GetErrCode()).String())
				}
			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func assets(lSession *login.Session) {
	fmt.Printf("----------------------type -->  assetsSync:1 assetsSet:2 \n")
	go func() {
		var t int
		fmt.Print("type=> ")
		fmt.Scan(&t)
		switch t {
		case 1:
			var id int32
			fmt.Print("DrawCardToken  = 1; ClientDatabase = 2; =>")
			fmt.Scan(&id)
			lSession.Send(&cs_msg.AssetsToS{
				AssetType: proto.Int32(id),
			})

		}
	}()
}
