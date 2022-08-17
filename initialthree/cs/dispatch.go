package cs

import (
	fnet "github.com/sniperHW/flyfish/pkg/net"
	codecs "initialthree/codec/cs"
)

type ServerDispatcher interface {
	Dispatch(*fnet.Socket, *codecs.Message)
	OnClose(*fnet.Socket, error)
	OnNewClient(*fnet.Socket)
	OnAuthenticate(*fnet.Socket) bool
}

type ClientDispatcher interface {
	Dispatch(*fnet.Socket, *codecs.Message)
	OnClose(*fnet.Socket, error)
	OnEstablish(*fnet.Socket)
	OnConnectFailed(peerAddr string, err error)
}
