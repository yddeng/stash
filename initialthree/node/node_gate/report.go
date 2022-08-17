package node_gate

import (
	"initialthree/cluster"
	"initialthree/node/common/serverType"
	ss_msg "initialthree/protocol/ss/ssmessage"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

var report_once sync.Once

func ReportInit(externalAddr string) {
	report_once.Do(func() {
		//启动一个go程，每10秒上报自己的服务状态
		go func() {
			report := &ss_msg.ReportGate{
				PeerID:       proto.String(cluster.SelfAddr().Logic.String()),
				ExternalAddr: proto.String(externalAddr),
			}
			for {
				report.PlayerNum = proto.Int32(0)
				//Infoln("report", report)
				cluster.Brocast(serverType.Login, report)
				time.Sleep(time.Second * 10)
			}
		}()
	})
}
