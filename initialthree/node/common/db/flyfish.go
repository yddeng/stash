package db

import (
	flyfish "github.com/sniperHW/flyfish/client"
	"go.uber.org/zap"
	"initialthree/cluster"
	"initialthree/cluster/priority"
	"strings"
)

var fclients map[string]*flyfish.Client = map[string]*flyfish.Client{}

func FlyfishInit(configs []string, l *zap.Logger) error {

	flyfish.InitLogger(l)

	//name@address
	for _, v := range configs {
		config := strings.Split(v, "@")
		if config[0] == "flyfish" {
			name := config[1]
			addr := config[2]

			clientCfg := flyfish.ClientConf{
				ClientType:     flyfish.FlyGate,
				NotifyPriority: priority.MID,
				NotifyQueue:    cluster.GetEventQueue(),
				PD:             strings.Split(addr, ";"),
			}

			c, err := flyfish.New(clientCfg)
			if err != nil {
				return err
			}
			fclients[name] = c
		}
	}
	return nil
}

func GetFlyfishClient(name string) *flyfish.Client {
	return fclients[name]
}
