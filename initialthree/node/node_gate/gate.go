package node_gate

import (
	"fmt"
	"github.com/sniperHW/flyfish/errcode"
	"github.com/sniperHW/flyfish/pkg/buffer"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/network/smux"
	"initialthree/node/common/db"
	"initialthree/node/common/serverType"
	"initialthree/node/common/taskPool"
	"initialthree/node/node_gate/monitor"
	"initialthree/zaplogger"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

func Start(externalAddr string) error {
	//初始化随机种子
	rand.Seed(time.Now().Unix())

	if err := initServerStatus(); nil != err {
		return err
	}

	t := strings.Split(externalAddr, ":")
	port, _ := strconv.Atoi(t[1])

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if nil != err {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		return err
	}

	go func() {
		zaplogger.GetSugar().Infof("node_gate start on %s", externalAddr)
		for {
			conn, e := listener.Accept()
			if e != nil {
				if ne, ok := e.(net.Error); ok && ne.Temporary() {
					zaplogger.GetSugar().Errorf("accept temp err: %v", ne)
					continue
				} else {
					return
				}
			} else {
				zaplogger.GetSugar().Info("new client")
				if err := loginPool.AddTask(func() {
					onNewClient(conn)
				}); err != nil {
					zaplogger.GetSugar().Info("gate: loginPool.AddTask error ", err)
					conn.Close()
				}
			}
		}
	}()

	return err
}

var (
	gameMap  = map[addr.LogicAddr]*gameSocket{}
	gameLock sync.Mutex

	loginPool = taskPool.NewTaskPool(runtime.NumCPU()*10, 2048)
)

type gameSocket struct {
	logic  addr.LogicAddr
	socket *smux.MuxSocket

	channels    map[int64]*channel
	channelLock sync.Mutex
}

func (g *gameSocket) addChannel(streamID int64, ch *channel) {
	g.channelLock.Lock()
	g.channels[streamID] = ch
	g.channelLock.Unlock()

	monitor.GameSocketStreamInc(g.logic.String())
	monitor.ChannelCountInc()

	ch.run(func(err error) {
		zaplogger.GetSugar().Infof("user:%s streamID:%d close reason :%s", ch.userID, streamID, err.Error())
		g.channelLock.Lock()
		delete(g.channels, streamID)
		g.channelLock.Unlock()

		monitor.GameSocketStreamDec(g.logic.String())
		monitor.ChannelCountDec()
	})

}

func modGame(userID string) (*gameSocket, *smux.MuxStream, error) {
	var err error
	var logic addr.LogicAddr

	// 用户登录过的game
	dbRet := db.GetFlyfishClient("game").Get("user_game_login", userID, "gameaddr").Exec()
	if errcode.GetCode(dbRet.ErrCode) == errcode.Errcode_ok {
		gameLogic := dbRet.Fields["gameaddr"].GetString()
		if gameLogic != "" {
			logic, _ = addr.MakeLogicAddr(gameLogic)
		}
	}

	if logic == addr.LogicAddr(0) {
		//if logic, err = cluster.Mod(serverType.Game, common.HashS(userID)); err != nil {
		//	return nil, nil, err
		//}
		if logic, err = cluster.LBMod(serverType.Game); err != nil {
			return nil, nil, err
		}
	}

	gameLock.Lock()

	g, ok := gameMap[logic]
	if !ok {
		g = &gameSocket{
			logic:    logic,
			channels: map[int64]*channel{},
		}
		gameMap[logic] = g
	}

	if g.socket == nil {
		socket, err := cluster.DialMuxSocket(logic, func(o interface{}, b *buffer.Buffer) error {
			b.AppendBytes(o.([]byte))
			return nil
		}, func(socket *smux.MuxSocket) {
			zaplogger.GetSugar().Infof("SmuxConn %s closed", logic.String())
			onSocketClose(logic)
			socket.Close()
		})

		if err != nil {
			gameLock.Unlock()
			return g, nil, err
		}
		g.socket = socket
	}
	gameLock.Unlock()

	ret := make(chan interface{}, 1)
	go func() {
		stream, err := g.socket.Dial(time.Second * 2)
		if err != nil {
			ret <- err
		} else {
			ret <- stream
		}
	}()

	i := <-ret
	switch i.(type) {
	case *smux.MuxStream:
		return g, i.(*smux.MuxStream), nil
	default:
		return g, nil, i.(error)
	}
}

func onSocketClose(logic addr.LogicAddr) {
	gameLock.Lock()
	if g, ok := gameMap[logic]; ok {
		g.channelLock.Lock()
		channel := g.channels
		g.channelLock.Unlock()
		delete(gameMap, logic)
		gameLock.Unlock()

		for _, c := range channel {
			c.close(fmt.Errorf("game socket %s closed. ", logic.String()))
		}
	} else {
		gameLock.Unlock()
	}
}

func broadcast(bytes []byte) {
	gameLock.Lock()
	defer gameLock.Unlock()
	for _, g := range gameMap {
		g.channelLock.Lock()
		for _, c := range g.channels {
			c.sendToTcp(bytes)
		}
		g.channelLock.Unlock()
	}
}

func init() {
	monitor.RegisterGameSocketCountFunc(func() float64 {
		gameLock.Lock()
		defer gameLock.Unlock()
		return float64(len(gameMap))
	})
}
