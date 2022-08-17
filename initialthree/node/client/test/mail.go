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

			d.Register(cmdEnum.CS_MailSync, func(session *fnet.Socket, msg *codecs.Message) {
				data := msg.GetData().(*cs_msg.MailSyncToC)
				fmt.Println("MailSyncToC", data, len(data.GetMails()))

				bytes, _ := proto.Marshal(data)
				fmt.Println("MailSyncToC length", len(bytes))

				for _, m := range data.GetMails() {
					fmt.Println("------------------------")
					fmt.Println(m.GetTitle())
					fmt.Println(m.GetSender())
					if m.GetExpireTime() != 0 {
						createTime := time.Unix(m.GetCreateTime(), 0)
						expireTime := time.Unix(m.GetExpireTime(), 0)
						fmt.Printf("%s %s\n", createTime.String(), expireTime.String())
					} else {
						createTime := time.Unix(m.GetCreateTime(), 0)
						fmt.Printf("%s %s\n", createTime.String(), "-")
					}
					fmt.Println(m.GetContent())
					//fmt.Println(m.GetAwards())
					fmt.Println("------------------------")
				}

			})
			d.Register(cmdEnum.CS_MailRead, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.MailReadToC)
				fmt.Println("MailReadToC", data)

				cmd(lSession)
			})
			d.Register(cmdEnum.CS_MailDelete, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.MailDeleteToC)
				fmt.Println("MailDeleteToC", data)
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
		fmt.Println("1:read 2:delete")
		var key int
		fmt.Scan(&key)
		if key == 1 {
			fmt.Println("input : id ")
			var id uint32
			fmt.Scan(&id)
			req := &cs_msg.MailReadToS{
				MailIDs: []uint32{id},
			}
			lSession.Send(req)
		} else if key == 2 {
			fmt.Println("input : id ")
			var id uint32
			fmt.Scan(&id)
			req := &cs_msg.MailDeleteToS{
				MailIDs: []uint32{id},
			}
			lSession.Send(req)
		}
	}()
}
