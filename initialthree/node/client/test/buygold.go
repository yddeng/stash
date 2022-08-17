package main

import (
	"fmt"
	codecs "initialthree/codec/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/node/client/login"
	"initialthree/node/common/attr"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"os"

	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("usage: addr userID\n")
		return
	}

	addr := os.Args[1]
	userID := os.Args[2]

	stopCh := make(chan struct{})

	login.Login(userID, addr, func(lSession *login.Session, d *dispatcher.Dispatcher, msg *codecs.Message, err error) {
		if err != nil {
			fmt.Println("login", err)
		} else {
			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {
				if msg.GetErrCode() == 0 {
					data := msg.GetData().(*cs_msg.CreateRoleToC)

					fmt.Println("create role successfully:", data)

					lSession.Send(&cs_msg.BuyGoldToS{})
				} else {
					fmt.Println("create role failed", cs_msg.ErrCode(msg.GetErrCode()).String())
					close(stopCh)
				}

			})

			d.Register(cmdEnum.CS_AttrSync, func(session *fnet.Socket, msg *codecs.Message) {
				if msg.GetErrCode() == 0 {
					resp := msg.GetData().(*cs_msg.AttrSyncToC)
					fmt.Println("attr:")
					for _, a := range resp.Attrs {
						switch a.GetId() {
						case attr.Gold, attr.GoldBuyCount, attr.Diamond:
							fmt.Printf("\t%s: %d\n", attr.GetNameById(a.GetId()), a.GetVal())
						}
					}

				}
			})

			d.Register(cmdEnum.CS_BuyGold, func(session *fnet.Socket, msg *codecs.Message) {
				if msg.GetErrCode() == 0 {
					fmt.Println("buy gold successfully")
				} else {
					fmt.Println("buy gold failed", cs_msg.ErrCode(msg.GetErrCode()).String())
					close(stopCh)
				}
			})

			if cs_msg.ErrCode(msg.GetErrCode()) == cs_msg.ErrCode_OK {
				data := msg.GetData().(*cs_msg.GameLoginToC)
				fmt.Println("game login:", data)

				if data.GetIsFirstLogin() {
					//fmt.Println("input role name ")
					req := &cs_msg.CreateRoleToS{
						Name: proto.String(userID),
					}
					lSession.Send(req)
				} else {
					lSession.Send(&cs_msg.BuyGoldToS{})
				}
			} else {
				fmt.Println("login failed", cs_msg.ErrCode(msg.GetErrCode()).String())
				close(stopCh)
			}
		}
	})

	_ = <-stopCh
}
