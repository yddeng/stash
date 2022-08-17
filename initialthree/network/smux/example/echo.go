package main

import (
	"fmt"
	"initialthree/network/smux"
	"initialthree/network/smux/example/server"
	"net"
	"time"
)

func main() {
	s := server.NewServer()
	s.Serve()

	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", s.Ln.Addr().String())
	if err != nil {
		panic(err)
	}

	socket := smux.NewMuxSocketClient(conn, server.Enc, nil)

	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:8110")
	if err != nil {
		panic(err.Error())
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("gate start on", ln.Addr().String())

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
			fmt.Println("on new client")
			go func() {
				ms, err := socket.Dial(time.Second * 5)
				if nil != err {
					fmt.Printf("dial error:%v", err)
					conn.Close()
				} else {
					ms.SetCloseCallback(func(_ *smux.MuxStream, err error) {
						fmt.Println("client muxstream close", err)
						conn.Close()
					})

					readbuff := make([]byte, 4096)

					for {
						if n, err := conn.Read(readbuff); nil == err {
							//注意:recv callback在socket的接收线程调用，如果conn.Write阻塞，将导致所有MuxStream都消息都被阻塞
							ms.Recv(func(_ *smux.MuxStream, data []byte) {
								conn.Write(data)
							})
							if nil != ms.SyncSend(readbuff[:n]) {
								return
							}
						} else {
							ms.Close(err)
							break
						}
					}
				}
			}()
		}
	}
}
