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
				//cmd(lSession)
			} else {
				cmd(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("------ CreateRoleToC", data)
				cmd(lSession)
			})

			d.Register(cmdEnum.CS_DrawCardSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.DrawCardSyncToC)
				fmt.Println("------ DrawCardSyncToC", data)
			})

			d.Register(cmdEnum.CS_DrawCardDraw, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.DrawCardDrawToC)
				fmt.Println("------ DrawCardDrawToC", data)
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
		fmt.Println("fight input : id ,count ")
		var id, diff int32
		fmt.Scan(&id, &diff)
		req := &cs_msg.DrawCardDrawToS{
			LibID:     proto.Int32(id),
			DrawCount: proto.Int32(diff),
		}
		lSession.Send(req)

	}()
}
