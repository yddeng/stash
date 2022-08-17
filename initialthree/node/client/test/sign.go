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
				_ = lSession.Send(&cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				})
			} else {
				cmd(lSession)
			}
			d.Register(cmdEnum.CS_BaseSync, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.BaseSyncToC)
				fmt.Println("BaseSyncToC", data)
			})
			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
				cmd(lSession)
			})
			d.Register(cmdEnum.CS_AttrSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.AttrSyncToC)
				fmt.Println("AttrSyncToC", data)
			})

			d.Register(cmdEnum.CS_SignSync, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.SignSyncToC)
				fmt.Println("SignSyncToC", data)
			})
			d.Register(cmdEnum.CS_SignIn, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.SignInToC)
				fmt.Println("SignInToC", data)

				cmd(lSession)
			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func cmd(lSession *login.Session) {
	go func() {
		fmt.Println("input : id ")
		var id int32
		fmt.Scan(&id)
		req := &cs_msg.SignInToS{
			Id: proto.Int32(id),
		}
		lSession.Send(req)

	}()
}
