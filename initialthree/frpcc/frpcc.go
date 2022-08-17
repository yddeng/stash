package main

import (
	"flag"
	"fmt"
	"github.com/fatedier/frp/server/verify"
	frpIo "github.com/fatedier/golib/io"
	"net"
)

type proxy struct {
	name    string
	svrAddr string
	svrPort int
}

func (p *proxy) handleUserConn(usrConn net.Conn) {

	fmt.Printf("handleUserConn %s remote:%v", p.name, usrConn.RemoteAddr())

	dialer := &net.Dialer{}
	serConn, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", p.svrAddr, p.svrPort))
	if nil != err {
		usrConn.Close()
		fmt.Println(err)
		return
	}

	if err = verify.Auth(serConn); nil != err {
		fmt.Println(err)
		usrConn.Close()
		serConn.Close()
		return
	}

	defer usrConn.Close()
	defer serConn.Close()

	frpIo.Join(usrConn, serConn)

}

func (p *proxy) run() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", p.svrPort))
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		panic(err)
	}

	fmt.Printf("%s start listen at:%d\n", p.name, p.svrPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		} else {
			go p.handleUserConn(conn)
		}
	}
}

func main() {

	config := flag.String("config", "frpcc.toml", "config")

	conf, err := LoadConfig(*config)
	if nil != err {
		panic(err)
	}

	for _, v := range conf.Proxy {
		p := &proxy{
			name:    v.Name,
			svrAddr: v.SvrAddr,
			svrPort: v.SvrPort,
		}
		go p.run()
	}

	c := make(chan struct{})

	<-c
}
