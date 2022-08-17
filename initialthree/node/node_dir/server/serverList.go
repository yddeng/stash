package server

import (
	"initialthree/cluster"
	"initialthree/node/common/config"
	"initialthree/pkg/timer"
	"time"
)

var dbQueryMax int32
var dbCurQueryCount int32 = 0

func InitDBQueryMax(num int32) {
	dbQueryMax = num
}

type Server struct {
	ServerId     int32  `toml:"serverID"`
	ServerName   string `toml:"serverName"`
	ServerAddr   string
	ServerStatus int //1运行中 0维护中
	timestamp    time.Time
	playerNum    int32
}

var (
	myServers *Servers
	timeout   = time.Second * 5
)

type Servers struct {
	ServerList []*Server `toml:"Servers"`
}

func LoadConfig(cfg *config.Dir) error {
	myServers = &Servers{ServerList: make([]*Server, len(cfg.Servers))}
	for i, v := range cfg.Servers {
		myServers.ServerList[i] = &Server{
			ServerId:   v.ServerId,
			ServerName: v.ServerName,
		}
	}

	cluster.RegisterTimer(time.Second, tick, nil)
	return nil
}

func (this *Servers) getServer(id int32) *Server {
	for _, v := range this.ServerList {
		if v.ServerId == id {
			return v
		}
	}
	return nil
}

func tick(t *timer.Timer, _ interface{}) {
	now := time.Now()
	for _, v := range myServers.ServerList {
		if !v.timestamp.IsZero() && v.timestamp.Add(timeout).Before(now) {
			v.ServerStatus = 0
		}
	}
}
