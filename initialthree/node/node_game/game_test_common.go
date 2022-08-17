//go:build game_test
// +build game_test

package node_game

import (
	"github.com/golang/protobuf/proto"
	"github.com/hqpko/hutils"
	"github.com/sniperHW/flyfish/server/mock/kvnode"
	"initialthree/center"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/db"
	"initialthree/node/common/functionSwitch"
	"initialthree/node/common/globalservice/logictimesys"
	"initialthree/node/common/timeDisposal"
	weathersys "initialthree/node/common/weathersystem/cluster"
	"initialthree/node/common/wordsFilter"
	"initialthree/node/node_game/global/scarsIngrain"
	"initialthree/pkg"
	"initialthree/pkg/event"
	"math/rand"
	"sync"
	"sync/atomic"

	_ "initialthree/node/node_game/module"
	_ "initialthree/node/node_game/transaction"
	"initialthree/node/table/excel"
	"initialthree/node/table/quest"
	"initialthree/protocol/cmdEnum"
	"initialthree/protocol/cs/message"
	_ "initialthree/protocol/cs/message"
	ss_rpc "initialthree/protocol/ss/rpc"
	ss_message "initialthree/protocol/ss/ssmessage"
	"initialthree/zaplogger"
	"time"
)

var centerAddr string = "127.0.0.1:8000"
var center1 *center.Center

var encoder = codecs.NewEncoder("cs")
var receiver = codecs.NewReceiver("sc")
var kvnode = mock_kvnode.New()

var seqno uint32

type testuniLocker struct {
}

func (this testuniLocker) Lock(_ addr.Addr) bool {
	return true
}

func (this testuniLocker) Unlock() {

}

type node_gate struct {
	c         *cluster.Cluster
	evHandler *event.EventHandler
	receiver  *codecs.Receiver
}

//处理发往客户端的消息
func (g *node_gate) Register(ev interface{}, callback interface{}) {
	g.evHandler.Register(ev, callback)
}

func (g *node_gate) RegisterOnce(ev interface{}, callback interface{}) {
	g.evHandler.RegisterOnce(ev, callback)
}

type eventWait struct {
	ch   chan interface{}
	drop int32
}

func (eW *eventWait) Pop() interface{} {
	return <-eW.ch
}

func (eW *eventWait) Drop() {
	atomic.AddInt32(&eW.drop, 1)
}

func (g *node_gate) RegisterWait(ev interface{}) *eventWait {
	eW := &eventWait{ch: make(chan interface{}, 2)}
	g.evHandler.Register(ev, func(ret interface{}) {
		if atomic.LoadInt32(&eW.drop) == 0 {
			eW.ch <- ret
		} else {
			atomic.AddInt32(&eW.drop, -1)
		}
	})
	return eW
}

//处理发往gate的消息
func (g *node_gate) RegisterMsgHandlerOnce(cmd uint16, callback interface{}) {
	g.evHandler.RegisterOnce(cmd, callback)
	g.c.Register(cmd, func(_ addr.LogicAddr, msg proto.Message) {
		g.evHandler.Emit(cmd, msg)
	})
}

func (g *node_gate) forwardClient(from addr.LogicAddr, msg proto.Message) {
	bytes := msg.(*ss_message.SsToGate).GetMessage()
	for _, v := range bytes {
		msg, err := receiver.DirectUnpack(v)
		if nil == err {
			g.evHandler.Emit(msg.(*codecs.Message).GetCmd(), msg.(*codecs.Message))
		}
	}
}

func (g *node_gate) PostMessageToGame(msg proto.Message) {
	g.c.PostMessage(cluster.SelfAddr().Logic, msg)
}

func (g *node_gate) Close() {
	g.c.Stop(nil, false)
}

func new_node_gate(address addr.Addr) *node_gate {
	gate := &node_gate{
		c:         cluster.NewCluster(),
		evHandler: event.NewEventHandler(),
		receiver:  codecs.NewReceiver("sc"),
	}

	if nil != gate.c.Start([]string{centerAddr}, address, testuniLocker{}) {
		return nil
	} else {
		return gate
	}
}

var gameInitOnce sync.Once

func GameInit() {

	gameInitOnce.Do(func() {

		//logger = util.NewLogger("./", "test_user", 1024*1024*50)
		logger := zaplogger.NewZapLogger("game_test.log", "log", "debug", 100, 14, true)
		zaplogger.InitLogger(logger)

		kendynet.InitLogger(zaplogger.GetSugar())
		cluster.InitLogger(zaplogger.GetSugar())

		rand.Seed(time.Now().Unix())

		center1 = center.New()
		go func() {
			center1.Start(centerAddr, zaplogger.GetSugar())
		}()

		kvnode.Start("localhost:12500", []string{
			"game_user@userdata:blob:,id:int:0",
			"user_game_login@gameaddr:string:",
			"role_name@owner:string:",
			"user_dir_last_login@lastlogin:int:0",
			"user_module_data@attr:blob:,attr_time:blob:,mapdata:blob:,base:blob:,rankdata:blob:",
			"character@slice0:blob:,slice1:blob:,slice2:blob:,slice3:blob:,slice4:blob:,slice5:blob:,slice6:blob:,slice7:blob:,slice8:blob:,slice9:blob:,group_ids:blob:,group_default:blob:",
			"function_switch@modules:string:,messages:string:",
			"global_data@data:blob:",
			"whitelist@ip:string:",
			"quest@timedata:blob:,main_story:blob:,daily:blob:,weekly:blob:,event:blob:,instance_rd:blob:,daily_rd:blob:,character_g:blob:,",
			"user_assets@assets:blob:",
			"backpack@capacity:blob:,slot0:blob:,slot1:blob:,slot2:blob:,slot3:blob:,slot4:blob:,slot5:blob:,slot6:blob:,slot7:blob:,slot8:blob:,slot9:blob:",
			"equip@slice0:blob:,slice1:blob:,slice2:blob:,slice3:blob:,groups:blob:",
			"msgbox@boxdata:blob:,slot1:blob:,slot2:blob:,slot3:blob:,slot4:blob:,slot5:blob:,slot6:blob:,slot7:blob:,slot8:blob:,slot9:blob:,slot10:blob:,slot11:blob:,slot12:blob:,slot13:blob:,slot14:blob:,slot15:blob:,slot16:blob:,slot17:blob:,slot18:blob:,slot19:blob:,slot20:blob:",
			"mailbox@boxdata:blob:,slot1:blob:,slot2:blob:,slot3:blob:,slot4:blob:,slot5:blob:,slot6:blob:,slot7:blob:,slot8:blob:,slot9:blob:,slot10:blob:,slot11:blob:,slot12:blob:,slot13:blob:,slot14:blob:,slot15:blob:,slot16:blob:,slot17:blob:,slot18:blob:,slot19:blob:,slot20:blob:",
			"main_dungeons@base:blob:,chapter0:blob:,chapter1:blob:,dungeon0:blob:,dungeon1:blob:,dungeon2:blob:,dungeon3:blob:,dungeon4:blob:,dungeon5:blob:,dungeon6:blob:,dungeon7:blob:",
			"draw_card_history@history:string:{}",
			"scarsingrain@sidata:blob:,daily_times:blob:,role_times:blob:,timedata:blob:,high_score:blob:",
			"materialdungeon@dungeon_data:blob:",
			"rewardquest@base:blob:,data:blob:,timedata:blob:",
			"temporary@level:blob:,seqno:blob:,",
			"nodelock@phyaddr:string:",
		})

		db.FlyfishInit([]string{
			"flyfish@login@localhost:12500",
			"flyfish@dir@localhost:12500",
			"flyfish@game@localhost:12500",
			"flyfish@global@localhost:12500",
			"flyfish@gm@localhost:12500",
			"flyfish@nodelock@localhost:12500",
		}, zaplogger.GetSugar())

		excel.Load("../../configs/Excel")

		gameaddr, _ := addr.MakeAddr("1.4.1", "127.0.0.1:8002")
		hutils.Must(nil, cluster.Start([]string{centerAddr}, gameaddr, testuniLocker{}))

		hutils.Must(nil, timeDisposal.Init(db.GetFlyfishClient("game"), zaplogger.GetSugar()))

		hutils.Must(nil, logictimesys.Start(true, "../../configs/TimeConfig.asset", db.GetFlyfishClient("game"), zaplogger.GetSugar()))

		hutils.Must(nil, weathersys.Launch("../../configs/WeatherConfig.asset", logictimesys.TimeMgr(), db.GetFlyfishClient("game"), zaplogger.GetSugar()))

		hutils.Must(nil, functionSwitch.Init(false, db.GetFlyfishClient("game"), zaplogger.GetSugar()))

		hutils.Must(nil, quest.Load("../../configs/quest"))

		wordsFilter.Init("../../configs/wordsFilter/wordsFilter.txt")

		scarsIngrain.Launch()
	})
}

func NewGate() *node_gate {
	gateaddr, _ := addr.MakeAddr("1.3.1", "127.0.0.1:8003")
	gate := new_node_gate(gateaddr)
	gate.c.Register(cmdEnum.SS_SsToGate, gate.forwardClient)
	return gate
}

func LoginAndCreateRole(gate *node_gate, userID string) (interface{}, error) {
	r, err := LoginUser(gate, userID)
	if nil != err {
		return nil, err
	}

	ret := r.(*ss_rpc.GateUserLoginResp)

	if ret.GetIsFirstLogin() {
		req := &message.CreateRoleToS{
			Name: proto.String(userID),
		}
		return ForwardMsg(gate, userID, req)
	} else {
		return nil, nil
	}
}

type rpcResult struct {
	r   interface{}
	err error
}

func LoginUser(gate *node_gate, userID string) (interface{}, error) {
	time.Sleep(time.Second)
	respCh := make(chan rpcResult)
	gate.c.AsynCall(cluster.SelfAddr().Logic, &ss_rpc.GateUserLoginReq{
		UserID:     proto.String(userID),
		GateUserID: proto.Uint64(1),
		ServerID:   proto.Int32(1),
		IsUserReq:  proto.Bool(true),
	}, time.Second*7, func(r interface{}, e error) {
		respCh <- rpcResult{
			r:   r,
			err: e,
		}
	})
	r := <-respCh
	return r.r, r.err
}

func ForwardMsg(gate *node_gate, userID string, msg proto.Message) (interface{}, error) {
	seqno := atomic.AddUint32(&seqno, 1)
	b, err := encoder.EnCode(codecs.NewMessage(seqno, msg))
	if err != nil {
		return nil, err
	}

	respCh := make(chan rpcResult)
	gate.c.AsynCall(cluster.SelfAddr().Logic, &ss_rpc.ForwardUserMsgReq{
		UserID:     proto.String(userID),
		GateUserID: proto.Uint64(1),
		Messages:   b.Bytes(),
		SeqNo:      proto.Int64(1),
	}, time.Second*5, func(r interface{}, e error) {
		respCh <- rpcResult{
			r:   r,
			err: e,
		}
	})
	r := <-respCh
	return r.r, r.err
}
