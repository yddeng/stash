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
				baseSet(userID, lSession)
			})

			d.Register(cmdEnum.CS_BaseSync, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.BaseSyncToC)
				fmt.Println("BaseSyncToC", data)
			})

			d.Register(cmdEnum.CS_BaseSetSignature, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.BaseSetSignatureToC)
				fmt.Println("BaseSetSignatureToC", data)
			})

			data := msg.GetData().(*cs_msg.GameLoginToC)
			if data.GetIsFirstLogin() {
				req := &cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				}

				lSession.Send(req)
			} else {
				baseSet(userID, lSession)
			}
		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func baseSet(userID string, lSession *login.Session) {
	req := &cs_msg.BaseSetSignatureToS{
		Signature: proto.String(userID),
	}
	lSession.Send(req)
}
