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
				cmd(lSession)
			})

			d.Register(cmdEnum.CS_ChatMessageSend, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.ChatMessageSendToC)
				fmt.Println("ChatMessageSendToC", data)
			})

			d.Register(cmdEnum.CS_ChatMessageSync, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.ChatMessageSyncToC)
				fmt.Println("ChatMessageSyncToC", data)
				cmd(lSession)
			})

			data := msg.GetData().(*cs_msg.GameLoginToC)
			if data.GetIsFirstLogin() {
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				}

				lSession.Send(req)
			} else {
				cmd(lSession)
			}
		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func cmd(lSession *login.Session) {
	go func() {

		fmt.Printf("send==>")
		var tt string
		fmt.Scan(&tt)

		req := &cs_msg.ChatMessageSendToS{
			Message: proto.String(tt),
		}
		lSession.Send(req)

	}()

}
