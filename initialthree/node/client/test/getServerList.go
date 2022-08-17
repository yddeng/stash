package main

import (
	"fmt"
	codecs "initialthree/codec/cs"
	"initialthree/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/pkg/event"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"os"

	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
)

func main() {

	if len(os.Args) < 3 {
		fmt.Printf(" addr userID\n")
		return
	}

	processQueue := event.NewEventQueue()
	go func() {
		processQueue.Run()
		fmt.Println("queue break-------------")
	}()

	userID := os.Args[2]
	addr := os.Args[1]

	dis := dispatcher.New(processQueue)
	dis.RegisterOnce("Establish", func(session *fnet.Socket) {
		session.Send(codecs.NewMessage(uint32(0), &cs_msg.ServerListToS{
			UserID: proto.String(userID),
		}))
	})

	dis.Register(cmdEnum.CS_ServerList, func(session *fnet.Socket, msg *codecs.Message) {

		if msg.GetErrCode() == 0 {
			data := msg.GetData().(*cs_msg.ServerListToC)
			fmt.Println("ServerListToC", data)
		} else {
			fmt.Println("ServerListToC ErrCode", msg.GetErrCode())
		}
	})

	cs.DialTcp(addr, 0, dis)

	sigStop := make(chan bool)
	_, _ = <-sigStop
}
