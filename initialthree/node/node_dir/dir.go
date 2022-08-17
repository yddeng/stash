package node_dir

import (
	"fmt"
	"initialthree/cs"
	"initialthree/zaplogger"
	"strconv"
	"strings"
)

func Start(externalAddr string) error {

	//if err := firewall.Init(db.GetFlyfishClient("dir"), new(firewallListener), zaplogger.GetSugar()); err != nil {
	//	return fmt.Errorf("init firewall: %s", err)
	//}

	t := strings.Split(externalAddr, ":")
	port, _ := strconv.Atoi(t[1])
	return cs.StartTcpServer("tcp", fmt.Sprintf("0.0.0.0:%d", port), &gDispatcher)
}

type firewallListener struct{}

//func (f firewallListener) OnUpdated(firewall.UpdatedArgs) {
//	fwSt := firewall.GetStatus()
//	zaplogger.GetSugar().Infof("firewall updated: %s", fwSt.String())
//}

func (f firewallListener) OnReloadError(err error) {
	zaplogger.GetSugar().Errorf("firewall reload: %s", err)
}
