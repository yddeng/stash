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
		fmt.Printf("usage selectRole addr userID \n")
		return
	}

	userID := os.Args[2]
	addr := os.Args[1]

	login.Login(userID, addr, func(lSession *login.Session, d *dispatcher.Dispatcher, msg *codecs.Message, err error) {
		if nil != err {
			fmt.Println(err)
		} else {

			data := msg.GetData().(*cs_msg.GameLoginToC)
			fmt.Println("GameLoginToc", data)

			if data.GetIsFirstLogin() {
				fmt.Println("input role name ")
			} else {
				lSession.Send(&cs_msg.SelectRoleToS{})
			}

			d.Register(cmdEnum.CS_EnterMap, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.EnterMapToC)
				fmt.Println("EnterMapToC", data)

				starAoi(lSession)

			})

			d.Register(cmdEnum.CS_StartAoi, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.StartAoiToC)
				fmt.Println("StartAoiToC", data)

				go move(lSession)

			})

			d.Register(cmdEnum.CS_UpdatePos, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.UpdatePosToC)
				fmt.Println("UpdatePosToC", data)

			})

			d.Register(cmdEnum.CS_EnterSee, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.EnterSeeToC)
				fmt.Println("EnterSeeToC", data)

			})

			d.Register(cmdEnum.CS_LeaveSee, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.LeaveSeeToC)
				fmt.Println("LeaveSeeToC", data)

			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func enterMap(lSession *login.Session) {
	fmt.Println("Send EnterMapToS")
	lSession.Send(&cs_msg.EnterMapToS{
		MapID: proto.Int32(int32(1)),
	})
}

func starAoi(lSession *login.Session) {

	fmt.Println("Send StartAoiToS")
	lSession.Send(&cs_msg.StartAoiToS{})
}

func move(lSession *login.Session) {
	x := 10
	_var := 100
	for {
		time.Sleep(500 * time.Millisecond)
		x += _var
		if x > 2000 || x < 10 {
			_var = -_var
		}

		lSession.Send(&cs_msg.MoveToS{
			Pos: &cs_msg.Position{
				X: proto.Int32(int32(x)),
				Y: proto.Int32(100),
				Z: proto.Int32(10),
			},
		})
	}

}
