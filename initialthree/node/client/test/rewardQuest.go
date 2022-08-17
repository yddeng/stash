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
				_ = lSession.Send(&cs_msg.CreateRoleToS{
					Name: proto.String(userID),
				})
			} else {
				cmd(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
				cmd(lSession)
			})

			d.Register(cmdEnum.CS_RewardQuestSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.RewardQuestSyncToC)
				fmt.Println("RewardQuestSyncToC", data)
			})
			d.Register(cmdEnum.CS_RewardQuestAccept, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.RewardQuestAcceptToC)
				fmt.Println("RewardQuestAcceptToC", data)
			})

			d.Register(cmdEnum.CS_RewardQuestComplete, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.RewardQuestCompleteToC)
				fmt.Println("RewardQuestCompleteToC", data)
			})

			d.Register(cmdEnum.CS_RewardQuestRefresh, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.RewardQuestRefreshToC)
				fmt.Println("RewardQuestRefreshToC", data)
			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func cmd(lSession *login.Session) {
	go func() {
		time.Sleep(time.Second)
		fmt.Println("1:accept 2:complete 3:refresh")
		var key int
		fmt.Scan(&key)
		if key == 1 {
			fmt.Println("input : id ")
			var id int
			fmt.Scan(&id)
			fmt.Println("input : roles ")
			var ids []int32
			fmt.Scan(&ids)
			req := &cs_msg.RewardQuestAcceptToS{
				QuestID:    proto.Int32(int32(id)),
				Characters: ids,
			}
			lSession.Send(req)
		} else if key == 2 {
			fmt.Println("input : id ")
			var id int
			fmt.Scan(&id)
			req := &cs_msg.RewardQuestCompleteToS{
				QuestID: proto.Int32(int32(id)),
			}
			lSession.Send(req)
		} else if key == 3 {
			fmt.Println("input : id ")
			var id int
			fmt.Scan(&id)
			req := &cs_msg.RewardQuestRefreshToS{
				QuestID: proto.Int32(int32(id)),
			}
			lSession.Send(req)
		}
	}()
}
