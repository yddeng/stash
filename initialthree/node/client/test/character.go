package main

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
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
				Character(lSession)
			}

			d.Register(cmdEnum.CS_CreateRole, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CreateRoleToC)
				fmt.Println("CreateRoleToC", data)
				Character(lSession)
			})
			d.Register(cmdEnum.CS_CharacterLevelUp, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CharacterLevelUpToC)
				fmt.Println("CharacterLevelUpToC", data)
			})
			d.Register(cmdEnum.CS_CharacterEvolution, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CharacterEvolutionToC)
				fmt.Println("CharacterEvolutionToC", data)
			})

			d.Register(cmdEnum.CS_CharacterTeamPrefabSet, func(session *fnet.Socket, msg *codecs.Message) {

				data := msg.GetData().(*cs_msg.CharacterTeamPrefabSetToC)
				fmt.Println("CharacterTeamPrefabSetToC", data)
			})

		}
	})

	sigStop := make(chan bool)
	_, _ = <-sigStop
}

func Character(lSession *login.Session) {
	go func() {

		fmt.Println("1:levelUp 2:evolution 3:group 4:groupDefSet ")
		var key int
		fmt.Scan(&key)
		if key == 1 {
			fmt.Println("input : id itemId count")
			var id, itemId, count int
			fmt.Scan(&id, &itemId, &count)
			req := &cs_msg.CharacterLevelUpToS{
				CharacterID: proto.Int32(int32(id)),
				CostItems: []*cs_msg.CostItem{
					{ItemID: proto.Int32(int32(itemId)),
						Count: proto.Int32(int32(count))},
				},
			}
			lSession.Send(req)
		} else if key == 2 {
			fmt.Println("input : id ")
			var id int
			fmt.Scan(&id)
			req := &cs_msg.CharacterEvolutionToS{
				CharacterID: proto.Int32(int32(id)),
			}
			_ = lSession.Send(req)
		} else if key == 3 {
		} else if key == 4 {
			req := &cs_msg.CharacterTeam{}
			for i := 0; i < 3; i++ {
				fmt.Printf("input idx :%d character", i)
				var v int
				fmt.Scan(&v)
				req.CharacterList = append(req.CharacterList, int32(v))
			}
			_ = lSession.Send(&cs_msg.CharacterTeamPrefabSetToS{
				CharacterTeam: req,
				TeamPrefabIdx: proto.Int32(-1),
			})
		}
	}()
}
