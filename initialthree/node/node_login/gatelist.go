package node_login

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common"
	"initialthree/protocol/cmdEnum"
	ss_rpc "initialthree/protocol/ss/rpc"
	ss_msg "initialthree/protocol/ss/ssmessage"
	"initialthree/rpc/synctoken"
	"initialthree/zaplogger"
	"math/rand"
	"sort"
	"time"
)

var gateMap map[addr.LogicAddr]*Gate
var gateList []*Gate //当前可用的网关列表

type Gate struct {
	clusterAddr addr.LogicAddr
	addr        string //对外服务地址
	playerNum   int32
	deadline    time.Time
}

//产生令牌
func genToken(user string) string {
	data := []byte(fmt.Sprintf("%s:%d", user, rand.Int()))
	sum := md5.Sum(data)
	return base64.StdEncoding.EncodeToString(sum[:])
}

func sortGateList() {
	sort.Slice(gateList, func(i, j int) bool {
		return uint32(gateList[i].clusterAddr) < uint32(gateList[j].clusterAddr)
	})
}

func GetGate(user string, callback func(error, string, string)) {
	hashCode := common.HashS(user)
	size := len(gateList)
	if size > 0 {
		i := hashCode % size
		gate := gateList[i]
		arg := &ss_rpc.SynctokenReq{
			Userid: proto.String(user),
			Token:  proto.String(genToken(user)),
		}

		fmt.Println("synctocken to ", gate.clusterAddr)

		//将令牌同步给选中的gate
		synctoken.AsynCall(gate.clusterAddr, arg, time.Second*10, func(result *ss_rpc.SynctokenResp, err error) {
			if nil != err {
				callback(err, "", "")
			} else {
				callback(nil, gate.addr, result.GetToken())
			}
		})
	} else {
		//没用可供使用的gate
		callback(fmt.Errorf("no avalilable gate"), "", "")
	}
}

func onReportGate(_ addr.LogicAddr, msg proto.Message) {
	report := msg.(*ss_msg.ReportGate)

	//kendynet.Infoln("onReportGate:", report.GetPeerID(), report.GetExternalAddr(), report.GetPlayerNum())

	peerAddr, err := addr.MakeLogicAddr(report.GetPeerID())

	if nil != err {
		zaplogger.GetSugar().Info("invaild:", report.GetPeerID())
		return
	}

	g, ok := gateMap[peerAddr]

	if ok {
		g.playerNum = report.GetPlayerNum()
		g.deadline = time.Now().Add(time.Second * 20)
	} else {
		g = &Gate{}
		g.clusterAddr = peerAddr
		g.addr = report.GetExternalAddr()
		g.playerNum = report.GetPlayerNum()
		g.deadline = time.Now().Add(time.Second * 20)
		gateMap[peerAddr] = g
		gateList = append(gateList, g)
		sortGateList()
	}
}

func init() {

	gateMap = map[addr.LogicAddr]*Gate{}
	gateList = []*Gate{}

	cluster.Register(cmdEnum.SS_ReportGate, onReportGate)
	go func() {
		for {
			//交给cluster的任务队列单线程执行
			cluster.PostTask(func() {
				gateList = gateList[0:0]
				now := time.Now()
				for k, v := range gateMap {
					if v.deadline.After(now) {
						gateList = append(gateList, v)
					} else {
						delete(gateMap, k)
					}
				}
				sortGateList()
			})
			time.Sleep(time.Second * 5)
		}
	}()
}
