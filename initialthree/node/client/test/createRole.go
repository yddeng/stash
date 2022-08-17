package main

import (
	"fmt"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/node/client/login"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"os"

	"github.com/golang/protobuf/proto"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("usage selectRole  addr userID\n")
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
			})

			data := msg.GetData().(*cs_msg.GameLoginToC)
			fmt.Println("GameLoginToc", data)

			if data.GetIsFirstLogin() {
				var n string
				fmt.Print("输入名字：")
				fmt.Scan(&n)
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(n),
				}
				lSession.Send(req)
			}

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}
