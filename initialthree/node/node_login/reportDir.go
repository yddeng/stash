package node_login

import (
	"initialthree/cluster"
	"initialthree/node/common/serverType"
	ss_meg "initialthree/protocol/ss/ssmessage"
	"time"

	"github.com/golang/protobuf/proto"
)

func reportStatus(serverGroups []int32) {
	//开启一个go程 向dir服上报状态
	go func() {
		for {
			time.Sleep(time.Second * 1)
			cluster.PostTask(func() {
				report := &ss_meg.ReportStatus{
					ServerAddr: proto.String(loginAddr),
					ServerIds:  serverGroups,
					PlayerNum:  proto.Int32(0),
				}

				//Infoln("ReportStatus", report)
				cluster.Brocast(serverType.Dir, report)
			})
		}
	}()
}
