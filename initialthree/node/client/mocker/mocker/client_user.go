package mocker

import (
	"fmt"
	"log"
	"sync"

	codecs "initialthree/codec/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/node/client/login"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"

	"github.com/gogo/protobuf/proto"
)

type ClientUser struct {
	w          sync.WaitGroup
	userID     string
	addr       string
	session    *login.Session
	dispatcher *dispatcher.Dispatcher

	callMap sync.Map
}

func NewClientUser(userID, addr string) *ClientUser {
	cu := &ClientUser{userID: userID, addr: addr}
	cu.login()
	return cu
}

func (cu *ClientUser) registerCommands(d *dispatcher.Dispatcher, session *login.Session) {
	d.Register(cmdEnum.CS_CreateRole, cu.handlerResp)
	d.Register(cmdEnum.CS_DrawCardStoreIn, cu.handlerResp)
}

func (cu *ClientUser) Run() {
	cu.actionDrawCardStoreIn()
}

func (cu *ClientUser) handlerResp(args ...interface{}) {
	msg := args[1].(*codecs.Message)
	mustBeZero(msg.GetErrCode())
	if v, ok := cu.callMap.Load(msg.GetSeriNo()); ok {
		c := v.(*call)
		c.doneWithReply(msg)
	}
}

func (cu *ClientUser) call(req proto.Message) *codecs.Message {
	seq := newSeq()
	call := newCall()
	cu.callMap.Store(seq, call)
	if err := cu.session.Send(req); err != nil {
		cu.callMap.Delete(seq)
		call.doneWithErr(err)
	}
	return must(call.Done()).(*codecs.Message)
}

func (cu *ClientUser) login() {
	cu.w.Add(1)
	login.Login(cu.userID, cu.addr, func(session *login.Session, d *dispatcher.Dispatcher, message *codecs.Message, err error) {
		log.Printf("client user login response ...")
		cu.dispatcher = d
		fmt.Println(d)
		cu.session = session
		cu.registerCommands(d, session)
		must(nil, err)
		mustBeZero(message.GetErrCode())
		if message.GetData().(*cs_msg.GameLoginToC).GetIsFirstLogin() {
			log.Printf("client user create role ...")
			resp := cu.call(&cs_msg.CreateRoleToS{Name: proto.String(cu.userID)})
			data := resp.GetData().(*cs_msg.CreateRoleToC)
			log.Printf("client user create role resp:%v", data)
			mustBeZero(resp.GetErrCode())
			must(nil, session.Send(newMessage(&cs_msg.CreateRoleToS{Name: proto.String(cu.userID)}).GetData()))
		}
		cu.w.Done()
	})
	cu.w.Wait()
}

func (cu *ClientUser) do(f func(w *sync.WaitGroup)) {
	cu.w.Add(1)
	f(&cu.w)
	cu.w.Wait()
}
