package server

import (
	"fmt"
	"github.com/sniperHW/kendynet/buffer"
	"initialthree/network/smux"
	"net"
	"sync"
	"time"
)

type server struct {
	sync.Mutex
	sockets map[*smux.MuxSocket]bool
	Ln      *net.TCPListener
}

func Enc(o interface{}, b *buffer.Buffer) error {
	b.AppendBytes(o.([]byte))
	return nil
}

func onDataServer(ms *smux.MuxStream, data []byte) {
	msg := make([]byte, len(data))
	copy(msg, data)
	ms.Recv(onDataServer)
	ms.AsyncSend(msg, nil)
}

func NewServer() *server {
	return &server{
		sockets: map[*smux.MuxSocket]bool{},
	}
}

func (s *server) onNewStream(ms *smux.MuxStream) {
	ms.SetRecvTimeout(time.Second * 5)
	ms.SetCloseCallback(func(_ *smux.MuxStream, err error) {
		fmt.Println("server muxstream close", err)
	})
	ms.Recv(onDataServer)

}

func (s *server) onSocketClose(ss *smux.MuxSocket) {
	s.Lock()
	delete(s.sockets, ss)
	s.Unlock()
}

func (s *server) Serve() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err.Error())
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err.Error())
	}

	s.Ln = ln

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					fmt.Printf("accept temp err: %v\n", ne)
					continue
				} else {
					return
				}
			} else {
				ss := smux.NewMuxSocketServer(conn, Enc, s.onSocketClose)
				fmt.Println("on new server socket")
				s.Lock()
				s.sockets[ss] = true
				s.Unlock()
				ss.Listen(s.onNewStream)
			}
		}
	}()
}
