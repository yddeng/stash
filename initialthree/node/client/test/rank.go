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
				getList(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
				getList(lSession)
			})

			d.Register(cmdEnum.CS_RankGetTopList, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.RankGetTopListToC)
				fmt.Println("RankGetTopListToC", data)
			})

			d.Register(cmdEnum.CS_RankGetRank, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.RankGetRankToC)
				fmt.Println("RankGetRankToC", data)
			})
		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func getList(lSession *login.Session) {
	fmt.Println("getList input rankId")
	var rankId uint32
	fmt.Scan(&rankId)
	lSession.Send(&cs_msg.RankGetTopListToS{
		Version: proto.Uint32(0),
		RankID:  proto.Uint32(rankId),
		GetLast: proto.Bool(false),
	})

	lSession.Send(&cs_msg.RankGetRankToS{
		RankID: proto.Uint32(rankId),
	})
}
