package node_gate

//go test -covermode=count -v -coverprofile=coverage.out -run=.
//go tool cover -html=coverage.out

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/server/mock/kvnode"
	"github.com/sniperHW/kendynet"
	"github.com/stretchr/testify/assert"
	"initialthree/center"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	codecs "initialthree/codec/cs"
	"initialthree/codec/relaysc"
	"initialthree/cs"
	"initialthree/node/common/db"
	"initialthree/pkg/event"
	"initialthree/pkg/rpc"
	"initialthree/pkg/timer"
	"math/rand"
	//firewall2 "initialthree/ops/firewall"
	flyfishdb "github.com/sniperHW/flyfish/backend/db"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	ss_rpc "initialthree/protocol/ss/rpc"
	ss_msg "initialthree/protocol/ss/ssmessage"
	"initialthree/rpc/synctoken"
	"initialthree/zaplogger"
	"testing"
	"time"
)

var centerAddr string = "127.0.0.1:8000"
var center1 *center.Center
var encoder = codecs.NewEncoder("sc")
var receiver = codecs.NewReceiver("cs")
var kvserver = kvnode.New()

func init() {
	/*logger = util.NewLogger("./", "test_node_gate", 1024*1024*50)
	kendynet.InitLogger(logger)
	InitLogger(logger)
	cluster.InitLogger(logger)*/
	logger := zaplogger.NewZapLogger("test_node_gate.log", "log", "debug", 100, 14, true)
	zaplogger.InitLogger(logger)
	kendynet.InitLogger(zaplogger.GetSugar())
	cluster.InitLogger(zaplogger.GetSugar())

	rand.Seed(time.Now().Unix())
	node_gate_test = true
	center1 = center.New()
	go func() {
		center1.Start(centerAddr, zaplogger.GetSugar())
	}()

	def, _ := flyfishdb.CreateDbDefFromCsv([]string{"whitelist@ip:string:", "global_data@data:blob:"})

	kvserver.Start("localhost:12500", def)

	db.FlyfishInit([]string{"flyfish@login@localhost:12500"}, logger)

	//fmt.Println(firewall2.Init(db.GetFlyfishClient("login")))

	//fmt.Println(firewall2.SetStatus(firewall2.StatusClose))

	//启动gate
	gateaddr, _ := addr.MakeAddr("1.3.1", "127.0.0.1:8002")
	cluster.Start([]string{centerAddr}, gateaddr, uniLocker{})

	Init()

	InitGroups([]int32{1})

	Start("127.0.0.1:8003")

	ReportInit("127.0.0.1:8003")

}

type uniLocker struct {
}

func (this uniLocker) Lock(_ addr.Addr) bool {
	return true
}

func (this uniLocker) Unlock() {

}

type handler func()

type Dispatcher struct {
	evHandler *event.EventHandler
	queue     *event.EventQueue
}

func (this *Dispatcher) RegisterOnce(ev interface{}, callback interface{}) {
	switch ev.(type) {
	case string:
		this.evHandler.RegisterOnce(ev, callback)
		break
	case uint16:
		this.evHandler.RegisterOnce(ev, callback)
		break
	default:
		break
	}
}

func (this *Dispatcher) Register(ev interface{}, callback interface{}) {
	switch ev.(type) {
	case string:
		this.evHandler.Register(ev, callback)
		break
	case uint16:
		this.evHandler.Register(ev, callback)
		break
	default:
		break
	}
}

func (this *Dispatcher) Emit(e interface{}, args ...interface{}) {
	this.evHandler.EmitToEventQueue(event.EventQueueParam{
		Q: this.queue,
	}, e, args...)
}

func (this *Dispatcher) Dispatch(session kendynet.StreamSession, msg *codecs.Message) {
	if nil != msg {
		cmd := msg.GetCmd()
		if cs_msg.ErrCode(msg.GetErrCode()) != cs_msg.ErrCode_OK {
			fmt.Printf("(seqNo:%d cmd:%d) errCode %s", msg.GetSeriNo(), msg.GetCmd(), cs_msg.ErrCode(msg.GetErrCode()).String())
		}
		this.Emit(cmd, session, msg)
	}
}

func (this *Dispatcher) OnClose(session kendynet.StreamSession, reason error) {
	fmt.Println("OnClose", reason)
	this.Emit("OnClose", session)
}

func (this *Dispatcher) OnEstablish(session kendynet.StreamSession) {
	fmt.Println("OnEstablish")
	this.Emit("Establish", session)
}

func (this *Dispatcher) OnConnectFailed(peerAddr string, err error) {
	fmt.Println("OnConnectFailed", err)
	this.Emit("ConnectFailed", peerAddr, err)
}

func NewDispatcher(processQueue *event.EventQueue) *Dispatcher {
	return &Dispatcher{
		evHandler: event.NewEventHandler(),
		queue:     processQueue,
	}
}

type client struct {
	uid      string
	token    string
	d        *Dispatcher
	queue    *event.EventQueue
	session  kendynet.StreamSession
	serverID int32
}

func new_client(uid string, token string) *client {
	queue := event.NewEventQueue()
	c := &client{
		uid:      uid,
		token:    token,
		queue:    queue,
		d:        NewDispatcher(queue),
		serverID: 1,
	}

	go c.queue.Run()

	return c
}

func (this *client) close() {
	this.queue.Close()
	if nil != this.session {
		this.session.Close(errors.New("none"), 0)
	}
}

func (this *client) send(seq uint32, msg proto.Message) error {
	return this.session.Send(codecs.NewMessage(seq, msg))
}

func (this *client) login(addr string) error {
	cc := make(chan interface{})
	timeout := make(chan interface{})

	this.d.RegisterOnce("Establish", func(session kendynet.StreamSession) {
		this.session = session
		cc <- struct{}{}
	})
	cs.DialTcp(addr, 0, this.d)

	<-cc

	//发送gameLogin

	gameLoginToS := &cs_msg.GameLoginToS{
		UserID:   proto.String(this.uid),
		Token:    proto.String(this.token),
		ServerID: proto.Int32(this.serverID),
	}

	this.d.RegisterOnce(cmdEnum.CS_GameLogin, func(session kendynet.StreamSession, msg *codecs.Message) {
		code := msg.GetErrCode()
		cc <- code
	})

	fmt.Println("send gamelogin")

	err := this.session.Send(codecs.NewMessage(uint32(0), gameLoginToS))
	fmt.Println("gameLoginToS", gameLoginToS, err)
	if nil != err {
		return err
	}

	timer.Once(time.Second*2, func(_ *timer.Timer, _ interface{}) {
		close(timeout)
	}, nil)

	select {
	case code := <-cc:
		if code.(uint16) != uint16(0) {
			return fmt.Errorf("login Error")
		} else {
			return nil
		}
	case <-timeout:
		return fmt.Errorf("login Error")
	}
}

type node_game struct {
	c *cluster.Cluster
}

func new_node_game(address addr.Addr) *node_game {

	game := &node_game{
		c: cluster.NewCluster(),
	}

	if nil != game.c.Start([]string{centerAddr}, address, uniLocker{}) {
		return nil
	} else {
		game.c.RegisterMethod(&ss_rpc.GateUserLoginReq{}, game.onLogin)
		game.c.RegisterMethod(&ss_rpc.ForwardUserMsgReq{}, game.onGateForwardMsg)

		return game
	}
}

func (this *node_game) onLogin(replyer *rpc.RPCReplyer, arg interface{}) {

	req := arg.(*ss_rpc.GateUserLoginReq)

	userID := req.GetUserID()

	fmt.Println("node_game.onLogin", userID)

	switch userID {
	case "testLoginTimeout":
		replyer.DropResponse()
	case "testLoginError":
		resp := &ss_rpc.GateUserLoginResp{
			Code:         ss_rpc.ErrCode_Error.Enum(),
			IsFirstLogin: proto.Bool(true),
		}
		replyer.Reply(resp, nil)
	case "delayLogin":
		time.Sleep(time.Second)
		resp := &ss_rpc.GateUserLoginResp{
			Code:         ss_rpc.ErrCode_OK.Enum(),
			IsFirstLogin: proto.Bool(true),
		}
		replyer.Reply(resp, nil)
	default:
		resp := &ss_rpc.GateUserLoginResp{
			Code:         ss_rpc.ErrCode_OK.Enum(),
			IsFirstLogin: proto.Bool(true),
		}
		replyer.Reply(resp, nil)
	}
}

func (this *node_game) sendToClient(guid uint64, seqNo uint32, msg proto.Message) {
	toCliMsg := codecs.NewMessage(seqNo, msg)

	gateaddr, _ := this.c.Random(3)

	cluster.Post(gateaddr, relaysc.NewMessage(toCliMsg, encoder, false, guid))

	/*buffer, err := encoder.EnCode(toCliMsg)
	if nil == err {
		ssMsg := &ss_msg.SsToGate{
			GateUsers: []uint64{guid},
			Message:   [][]byte{buffer.Bytes()},
		}
		this.c.PostMessage(gateaddr, ssMsg)
	}*/
}

func (this *node_game) broadcastToClient(msg proto.Message) {
	toCliMsg := codecs.NewMessage(0, msg)

	gateaddr, _ := this.c.Random(3)

	cluster.Post(gateaddr, relaysc.NewMessage(toCliMsg, encoder, true))
}

func (this *node_game) processMessage(replyer *rpc.RPCReplyer, uid string, guid uint64, seqNo uint32, message string) {

	fmt.Println("processMessage", message)

	switch message {
	case "echo":
		this.sendToClient(guid, seqNo, &cs_msg.EchoToC{Msg: proto.String("echo")})
	case "broadcast":
		this.broadcastToClient(&cs_msg.EchoToC{Msg: proto.String("broadcast")})
	case "forword error1":
		replyer.Reply(nil, fmt.Errorf("error"))
		return
	case "forword error2":
		replyer.Reply(&ss_rpc.ForwardUserMsgResp{
			Code: ss_rpc.ErrCode_Error.Enum(),
		}, nil)
		return
	case "forword timeout":
		replyer.DropResponse()
		return
	case "delay forword":
		time.Sleep(time.Millisecond * 200)

	case "kick gate user1":

		replyer.Reply(&ss_rpc.ForwardUserMsgResp{
			Code: ss_rpc.ErrCode_OK.Enum(),
		}, nil)

		msg := &ss_msg.KickGateUser{
			UserID: proto.String(uid),
		}

		gateaddr, _ := this.c.Random(3)

		this.c.PostMessage(gateaddr, msg)

		return

	case "kick gate user2":

		msg := &ss_msg.KickGateUser{
			UserID: proto.String(uid),
		}

		gateaddr, _ := this.c.Random(3)

		this.c.PostMessage(gateaddr, msg)

	case "synctoken":
		gateaddr, _ := this.c.Random(3)
		arg := &ss_rpc.SynctokenReq{
			Userid: proto.String("test"),
			Token:  proto.String("test"),
		}
		synctoken.AsynCall(gateaddr, arg, time.Second*10, func(result *ss_rpc.SynctokenResp, err error) {
		})
	default:
	}

	replyer.Reply(&ss_rpc.ForwardUserMsgResp{
		Code: ss_rpc.ErrCode_OK.Enum(),
	}, nil)

}

func (this *node_game) onGateForwardMsg(replyer *rpc.RPCReplyer, arg interface{}) {

	fmt.Println("onGateForwardMsg")

	msg, err := receiver.DirectUnpack(arg.(*ss_rpc.ForwardUserMsgReq).GetMessages())
	if nil == err {
		seqNo := msg.(*codecs.Message).GetSeriNo()
		guid := arg.(*ss_rpc.ForwardUserMsgReq).GetGateUserID()
		message := msg.(*codecs.Message).GetData().(*cs_msg.EchoToS).GetMsg()
		this.processMessage(replyer, arg.(*ss_rpc.ForwardUserMsgReq).GetUserID(), guid, seqNo, message)
	}
}

func TestID(t *testing.T) {
	ids := []uint64{}
	nums := make([]uint16, cap(gidGenerator.pool))

	for i := 0; i < cap(gidGenerator.pool); i++ {
		id, err := gidGenerator.Get()
		assert.Nil(t, err)
		n := id >> 48
		assert.Equal(t, true, n <= uint64(cap(gidGenerator.pool)))
		ids = append(ids, id)
		nums[n-1] = uint16(n)
	}

	id, err := gidGenerator.Get()
	assert.Equal(t, id, uint64(0))
	assert.Equal(t, err, noFreeIDError)

	for i, v := range nums {
		assert.Equal(t, uint16(i+1), v)
	}

	for _, v := range ids {
		gidGenerator.Release(v)
	}
}

func TestNodeGate(t *testing.T) {

	cluster.PostTask(func() {
		AddToken("sniperHW", "sniperHW", time.Second*60)
		AddToken("sniperHW1", "sniperHW1", time.Second*60)
		AddToken("testLoginTimeout", "testLoginTimeout", time.Second*60)
		AddToken("testLoginError", "testLoginError", time.Second*60)
		AddToken("oneSecond", "oneSecond", time.Second)
		AddToken("delayLogin", "delayLogin", time.Second*60)
	})

	{
		c := new_client("sniperHW", "sniperHW")

		err := c.login("127.0.0.1:8003")

		assert.NotNil(t, err)

	}

	userCount()

	assert.Equal(t, false, CheckToken("a", "b"))

	//启动game
	gameaddr, _ := addr.MakeAddr("1.4.1", "127.0.0.1:8001")
	game := new_node_game(gameaddr)

	time.Sleep(time.Second * 2)

	fmt.Println("test1")

	{
		/*c := new_client("sniperHW", "bad")

		err := c.login("127.0.0.1:8003")

		fmt.Println(err)*/

		c := new_client("sniperHW", "bad")
		c.serverID = 2

		err := c.login("127.0.0.1:8003")

		fmt.Println(err)
	}

	fmt.Println("test2")

	{
		c := new_client("sniperHW", "sniperHW")

		err := c.login("127.0.0.1:8003")

		assert.Nil(t, err)

		cc := make(chan interface{})

		c.d.Register(cmdEnum.CS_Echo, func(session kendynet.StreamSession, msg *codecs.Message) {
			data := msg.GetData().(*cs_msg.EchoToC)
			fmt.Println(data)
			cc <- struct{}{}
		})

		c.send(1, &cs_msg.EchoToS{Msg: proto.String("echo")})

		c.send(1, &cs_msg.EchoToS{Msg: proto.String("broadcast")})

		<-cc

		<-cc

		c.session.Close(errors.New("none"), 0)

		time.Sleep(time.Second * 1)
	}
	fmt.Println("test3")
	{
		c := new_client("testLoginTimeout", "testLoginTimeout")

		err := c.login("127.0.0.1:8003")

		assert.NotNil(t, err)

		c.session.Close(errors.New("none"), 0)

		time.Sleep(time.Second * 1)
	}
	fmt.Println("test4")
	{
		c := new_client("testLoginError", "testLoginError")

		err := c.login("127.0.0.1:8003")

		assert.NotNil(t, err)

		c.session.Close(errors.New("none"), 0)

		time.Sleep(time.Second * 1)
	}
	fmt.Println("test5")
	//测试id用完
	{
		ids := []uint64{}

		for i := 0; i < cap(gidGenerator.pool); i++ {
			id, _ := gidGenerator.Get()
			ids = append(ids, id)
		}

		c := new_client("sniperHW", "sniperHW")

		err := c.login("127.0.0.1:8003")

		assert.NotNil(t, err)

		for _, v := range ids {
			gidGenerator.Release(v)
		}
	}
	fmt.Println("test6")
	{
		c := new_client("sniperHW", "sniperHW")

		err := c.login("127.0.0.1:8003")

		assert.Nil(t, err)

		c.send(1, &cs_msg.EchoToS{Msg: proto.String("forword error1")})

		time.Sleep(time.Second * 1)

		c.session.Close(errors.New("none"), 0)

		time.Sleep(time.Second * 1)
	}
	fmt.Println("test7")
	{
		c := new_client("sniperHW", "sniperHW")

		err := c.login("127.0.0.1:8003")

		assert.Nil(t, err)

		c.send(1, &cs_msg.EchoToS{Msg: proto.String("forword error2")})

		time.Sleep(time.Second * 1)

		c.session.Close(errors.New("none"), 0)

		time.Sleep(time.Second * 1)
	}
	fmt.Println("test8")
	{
		c := new_client("sniperHW", "sniperHW")

		err := c.login("127.0.0.1:8003")

		assert.Nil(t, err)

		cc := make(chan interface{})

		c.d.Register("OnClose", func(session kendynet.StreamSession) {
			cc <- struct{}{}
		})

		c.send(1, &cs_msg.EchoToS{Msg: proto.String("forword timeout")})

		<-cc

	}

	/*fmt.Println("test9")
	{

		back_maxSendQueueBytes := maxSendQueueBytes

		maxSendQueueBytes = 10

		c := new_client("sniperHW", "sniperHW")

		err := c.login("127.0.0.1:8003")

		assert.Nil(t, err)

		cc := make(chan interface{})

		c.d.Register(cmdEnum.CS_Echo, func(session kendynet.StreamSession, msg *codecs.Message) {
			assert.Equal(t, cs_msg.ErrCode_RETRY, cs_msg.ErrCode(msg.GetErrCode()))
			cc <- struct{}{}
		})

		//发送心跳
		c.send(0, &cs_msg.HeartbeatToS{})

		c.send(3, &cs_msg.EchoToS{Msg: proto.String("fasdfasdfasfasfsfasfasfasfasfdasdfasfasfasfasfasdffasdfasdf")})

		<-cc

		maxSendQueueBytes = back_maxSendQueueBytes

		c.session.Close(errors.New("none"), 0)

	}*/
	fmt.Println("test10")
	{

		c1 := new_client("sniperHW", "sniperHW")

		err := c1.login("127.0.0.1:8003")

		assert.Nil(t, err)

		c2 := new_client("sniperHW", "sniperHW")

		err = c2.login("127.0.0.1:8003")

		assert.NotNil(t, err)

		c1.session.Close(errors.New("none"), 0)
		c2.session.Close(errors.New("none"), 0)

	}
	fmt.Println("test11")
	{

		cc := make(chan interface{})

		c1 := new_client("delayLogin", "delayLogin")

		c2 := new_client("delayLogin", "delayLogin")

		go func() {

			err := c1.login("127.0.0.1:8003")

			assert.NotNil(t, err)

			cc <- struct{}{}
		}()

		go func() {

			err := c2.login("127.0.0.1:8003")

			assert.NotNil(t, err)
			cc <- struct{}{}
		}()

		<-cc
		<-cc

		c1.session.Close(errors.New("none"), 0)
		c2.session.Close(errors.New("none"), 0)

	}
	fmt.Println("test12")
	{
		cc := make(chan interface{})

		c1 := new_client("delayLogin", "delayLogin")

		c1.d.Register("OnClose", func(session kendynet.StreamSession) {
			cc <- struct{}{}
		})

		go func() {

			c1.login("127.0.0.1:8003")

		}()

		time.Sleep(10 * time.Millisecond)
		c1.session.Close(errors.New("none"), 0)

		<-cc

	}
	fmt.Println("test13")
	{
		cc := make(chan interface{})

		c1 := new_client("sniperHW", "sniperHW")

		c1.d.Register("OnClose", func(session kendynet.StreamSession) {
			cc <- struct{}{}
		})

		c1.login("127.0.0.1:8003")

		c1.send(3, &cs_msg.EchoToS{Msg: proto.String("delay forword")})

		time.Sleep(10 * time.Millisecond)
		c1.session.Close(errors.New("none"), 0)

		<-cc

		time.Sleep(time.Second)

	}
	fmt.Println("test14")
	{
		cc := make(chan interface{})

		c1 := new_client("sniperHW", "sniperHW")

		c1.d.Register("OnClose", func(session kendynet.StreamSession) {
			cc <- struct{}{}
		})

		c1.login("127.0.0.1:8003")

		c1.send(3, &cs_msg.EchoToS{Msg: proto.String("synctoken")})

		time.Sleep(10 * time.Millisecond)
		c1.session.Close(errors.New("none"), 0)

		<-cc

		time.Sleep(time.Second)

	}

	fmt.Println("test15")
	{
		cc := make(chan interface{})

		c1 := new_client("sniperHW", "sniperHW")

		c1.d.Register("OnClose", func(session kendynet.StreamSession) {
			cc <- struct{}{}
		})

		assert.Nil(t, c1.login("127.0.0.1:8003"))

		c1.send(3, &cs_msg.EchoToS{Msg: proto.String("kick gate user1")})

		<-cc

		time.Sleep(time.Second)
	}
	fmt.Println("test16")
	{
		cc := make(chan interface{})

		c1 := new_client("sniperHW", "sniperHW")

		c1.d.Register("OnClose", func(session kendynet.StreamSession) {
			cc <- struct{}{}
		})

		assert.Nil(t, c1.login("127.0.0.1:8003"))

		c1.send(3, &cs_msg.EchoToS{Msg: proto.String("kick gate user2")})

		<-cc

		time.Sleep(time.Second)
	}
	/*
		fmt.Println("test17")
		{
			cc := make(chan interface{})

			c1 := new_client("sniperHW", "sniperHW")

			c1.d.Register("OnClose", func(session kendynet.StreamSession) {
				cc <- struct{}{}
			})

			assert.Nil(t, c1.login("127.0.0.1:8003"))

			fmt.Println(firewall2.SetStatus(firewall2.StatusAll))

			<-cc

		}

			fmt.Println("test18")

			{
				c1 := new_client("sniperHW", "sniperHW")
				assert.NotNil(t, c1.login("127.0.0.1:8003"))

				fmt.Println(firewall2.SetStatus(firewall2.StatusClose))
				for {

					if firewallStatus == firewall2.StatusClose {
						break
					} else {
						time.Sleep(time.Second)
					}

				}

			}

			fmt.Println("test19")

			{

				cc := make(chan interface{})

				c1 := new_client("sniperHW", "sniperHW")

				c1.d.Register("OnClose", func(session kendynet.StreamSession) {
					cc <- struct{}{}
				})

				assert.Nil(t, c1.login("127.0.0.1:8003"))

				firewall2.UpdateBlacklist([]string{"127.0.0.1"}, nil)

				fmt.Println(firewall2.SetStatus(firewall2.StatusBlacklist))

				<-cc

				fmt.Println(firewall2.SetStatus(firewall2.StatusClose))
				for {

					if firewallStatus == firewall2.StatusClose {
						break
					} else {
						time.Sleep(time.Second)
					}

				}

			}
	*/
	/*fmt.Println("test20")

	{
		cc := make(chan interface{})

		c1 := new_client("sniperHW", "sniperHW")
		c2 := new_client("sniperHW1", "sniperHW1")

		c1.d.Register("OnClose", func(session kendynet.StreamSession) {
			cc <- struct{}{}
		})

		c2.d.Register("OnClose", func(session kendynet.StreamSession) {
			cc <- struct{}{}
		})

		assert.Nil(t, c1.login("127.0.0.1:8003"))
		assert.Nil(t, c2.login("127.0.0.1:8003"))

		c1.send(3, &cs_msg.EchoToS{Msg: proto.String("forword timeout")})

		time.Sleep(10 * time.Millisecond)

		game.c.Stop(nil, true)

		//c2.send(3, &cs_msg.EchoToS{Msg: proto.String("kick gate user2")})

		<-cc
		<-cc

		time.Sleep(time.Second)
	}*/

	game.c.Stop(nil, true)

}
