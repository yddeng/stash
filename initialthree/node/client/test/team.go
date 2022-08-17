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
				fmt.Println("1 teamCreate, 2 teamJoinApply ")
				var b int
				fmt.Scan(&b)
				switch b {
				case 1:
					req := &cs_msg.TeamCreateToS{Target: &cs_msg.TeamTarget{LevelID: proto.Int32(1)}}
					_ = lSession.Send(req)
				case 2:
					fmt.Printf("input team id :")
					var id int
					fmt.Scan(&id)
					req := &cs_msg.TeamJoinApplyToS{
						TeamID: proto.Uint32(uint32(id)),
					}
					_ = lSession.Send(req)
				default:

				}
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)

			})

			d.Register(cmdEnum.CS_TeamSync, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamSyncToC)
				fmt.Println("TeamSyncToC", data)

				fmt.Println("1 踢人 ，2 解散，其他情况忽略")
				var b int
				fmt.Scan(&b)
				if b == 1 {
					fmt.Printf("input user id :")
					var id int
					fmt.Scan(&id)
					lSession.Send(&cs_msg.TeamKickPlayerToS{
						KickPlayerID: proto.Uint64(uint64(id)),
					})
				} else if b == 2 {
					lSession.Send(&cs_msg.TeamDismissToS{})
				}
			})
			d.Register(cmdEnum.CS_TeamDismiss, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamDismissToC)
				fmt.Println("TeamDismissToC", data)
			})
			d.Register(cmdEnum.CS_TeamDismissNotify, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamDismissNotifyToC)
				fmt.Println("TeamDismissNotifyToC", data)
			})
			d.Register(cmdEnum.CS_TeamPlayerJoinNotify, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamPlayerJoinNotifyToC)
				fmt.Println("TeamPlayerJoinNotifyToC", data)
			})
			d.Register(cmdEnum.CS_TeamPlayerLeaveNotify, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamPlayerLeaveNotifyToC)
				fmt.Println("TeamPlayerLeaveNotifyToC", data)
			})
			d.Register(cmdEnum.CS_TeamHeaderChangedNotify, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamHeaderChangedNotifyToC)
				fmt.Println("TeamHeaderChangedNotifyToC", data)
			})
			d.Register(cmdEnum.CS_TeamKickPlayer, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamKickPlayerToC)
				fmt.Println("TeamKickPlayerToC", data)
			})
			d.Register(cmdEnum.CS_TeamKickPlayerNotify, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamKickPlayerNotifyToC)
				fmt.Println("TeamKickPlayerNotifyToC", data)
			})

			d.Register(cmdEnum.CS_TeamCreate, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() == 0 {
					data := msg.GetData().(*cs_msg.TeamCreateToC)
					fmt.Println("TeamCreateToC", data)
				} else {
					fmt.Println("TeamCreateToC ErrCode", cs_msg.ErrCode(msg.GetErrCode()).String())
				}
			})

			d.Register(cmdEnum.CS_TeamJoinApply, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() == 0 {
					data := msg.GetData().(*cs_msg.TeamJoinApplyToC)
					fmt.Println("TeamJoinApplyToC", data)
				} else {
					fmt.Println("TeamJoinApplyToC ErrCode", cs_msg.ErrCode(msg.GetErrCode()).String())
				}
			})
			d.Register(cmdEnum.CS_TeamJoinApplyNotify, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamJoinApplyNotifyToC)
				fmt.Println("TeamJoinApplyNotifyToC", data)
				fmt.Println("1 同意，其他情况不同意")

				req := &cs_msg.TeamJoinReplyToS{
					Agree:   proto.Bool(false),
					AgreeID: proto.Uint64(data.GetPlayer().GetPlayerID()),
				}

				var b int
				fmt.Scan(&b)
				if b == 1 {
					req.Agree = proto.Bool(true)
				}

				_ = lSession.Send(req)

			})
			d.Register(cmdEnum.CS_TeamJoinReply, func(session *fnet.Socket, msg *codecs.Message) {

				if msg.GetErrCode() == 0 {
					data := msg.GetData().(*cs_msg.TeamJoinReplyToC)
					fmt.Println("TeamJoinReplyToC", data)
				} else {
					fmt.Println("TeamJoinReplyToC ErrCode", cs_msg.ErrCode(msg.GetErrCode()).String())
				}
			})
			d.Register(cmdEnum.CS_TeamJoinReplyNotify, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.TeamJoinReplyNotifyToC)
				fmt.Println("TeamJoinReplyNotifyToC", data)
			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}
