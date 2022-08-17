package network

import (
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"net"
)

func CreateSession(conn net.Conn) *fnet.Socket {
	return fnet.NewSocket(conn, fnet.OutputBufLimit{})
}
