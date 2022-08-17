package smux

//go test -covermode=count -v -coverprofile=coverage.out -run=.
//go tool cover -html=coverage.out

import (
	"errors"
	"fmt"
	"github.com/sniperHW/flyfish/pkg/buffer"
	"github.com/stretchr/testify/assert"
	"initialthree/zaplogger"
	"net"
	"strings"
	"sync"
	"testing"
	"time"
)

func init() {
	logger := zaplogger.NewZapLogger("test", "log", "debug", 100, 14, true)
	zaplogger.InitLogger(logger)
}

type server struct {
	sync.Mutex
	sockets  map[*MuxSocket]bool
	ln       *net.TCPListener
	dontRecv bool
}

func enc(o interface{}, b *buffer.Buffer) error {
	b.AppendBytes(o.([]byte))
	return nil
}

func onDataServer(ms *MuxStream, data []byte) {
	fmt.Println("ondata", ms.ID(), len(data))
	msg := make([]byte, len(data))
	copy(msg, data)
	ms.Recv(onDataServer)
	ms.AsyncSend(msg, nil)
}

func newServer() *server {
	return &server{
		sockets: map[*MuxSocket]bool{},
	}
}

func (s *server) onNewStream(ms *MuxStream) {

	ms.SetCloseCallback(func(_ *MuxStream, err error) {
		fmt.Println("server muxstream close", err)
	})

	if !s.dontRecv {
		time.Sleep(time.Second)
		ms.SetRecvTimeout(time.Second * 10)
		fmt.Println("recv start", ms.ID())
		if err := ms.Recv(onDataServer); err != nil {
			fmt.Println(err)
		}
		ms.SetRecvTimeout(time.Second * 2)
	}
}

func (s *server) onSocketClose(ss *MuxSocket) {
	s.Lock()
	delete(s.sockets, ss)
	s.Unlock()
}

func (s *server) serve(service string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		panic(err.Error())
	}

	ln, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err.Error())
	}

	s.ln = ln

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
				ss := NewMuxSocketServer(conn, enc, s.onSocketClose)
				fmt.Println("on new server socket")
				s.Lock()
				s.sockets[ss] = true
				s.Unlock()
				ss.Listen(s.onNewStream)
			}
		}
	}()
}

func Test1(t *testing.T) {
	s := newServer()
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	var onData func(ms *MuxStream, data []byte)

	onData = func(ms *MuxStream, data []byte) {
		count := ms.GetUserData().(*int)
		ms.Recv(onData)
		(*count)++
		if *count < 10 {
			msg := make([]byte, len(data))
			copy(msg, data)
			ms.SyncSend(msg)
		}
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	c := make(chan struct{})

	ms, err := socket.Dial(time.Second * 5)
	if nil != err {
		fmt.Printf("dial error:%v", err)
	} else {

		ms.SetCloseCallback(func(_ *MuxStream, err error) {
			fmt.Println("client muxstream close", err)
			close(c)
		})

		count := 0
		ms.SetUserData(&count)
		ms.Recv(onData)
		ms.SyncSend([]byte("hello"))
	}

	<-c

	s.ln.Close()

}

func Test2(t *testing.T) {
	s := newServer()
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	c := make(chan struct{})

	ms, err := socket.Dial(time.Second * 5)
	if nil != err {
		fmt.Printf("dial error:%v", err)
	} else {

		ms.SetCloseCallback(func(_ *MuxStream, err error) {
			fmt.Println("client muxstream close", err)
			close(c)
		})
	}

	socket.Close()

	<-c

	s.ln.Close()

}

func Test3(t *testing.T) {
	s := newServer()
	s.dontRecv = true
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	//loop:
	ms, err := socket.Dial(time.Second * 5)
	if nil != err {
		fmt.Printf("dial error:%v", err)
	} else {
		ms.SetCloseCallback(func(_ *MuxStream, err error) {
			fmt.Println("client muxstream close", err)
		})

		err := ms.SyncSend([]byte("hello"), time.Second*1)
		assert.Equal(t, err, Err_SendTimeout)
		//goto loop
	}

	s.ln.Close()

}

func Test4(t *testing.T) {
	s := newServer()
	s.dontRecv = true
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	ms, err := socket.Dial(time.Second * 5)
	if nil != err {
		fmt.Printf("dial error:%v", err)
	} else {
		ms.SetCloseCallback(func(_ *MuxStream, err error) {
			fmt.Println("client muxstream close", err)
		})

		errCh := make(chan error)
		ms.AsyncSend([]byte("hello"), func(_ *MuxStream, e error) {
			errCh <- e
		}, time.Second*1)

		err := <-errCh
		assert.Equal(t, err, Err_SendTimeout)
		ms.Close(err)
	}

	s.ln.Close()

}

func Test5(t *testing.T) {
	s := newServer()
	s.dontRecv = true
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	//loop:
	ms, err := socket.Dial(time.Second * 5)
	if nil != err {
		fmt.Printf("dial error:%v", err)
	} else {
		ms.SetCloseCallback(func(_ *MuxStream, err error) {
			fmt.Println("client muxstream close", err)
		})

		go func() {
			time.Sleep(time.Second)
			ms.Close(errors.New("client active close"))
		}()

		err := ms.SyncSend([]byte("hello"), time.Second*2)
		fmt.Println(err)
	}

	s.ln.Close()

}

func Test6(t *testing.T) {
	s := newServer()
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	var onData func(ms *MuxStream, data []byte)

	onData = func(ms *MuxStream, data []byte) {
		count := ms.GetUserData().(*int)
		ms.Recv(onData)
		(*count)++
		if *count < 10 {
			msg := make([]byte, len(data))
			copy(msg, data)
			ms.SyncSend(msg)
		}
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	c := make(chan struct{})

	ms, err := socket.Dial(time.Second * 5)
	if nil != err {
		fmt.Printf("dial error:%v", err)
	} else {

		ms.SetCloseCallback(func(_ *MuxStream, err error) {
			fmt.Println("client muxstream close", err)
			close(c)
		})

		count := 0
		ms.SetUserData(&count)
		ms.Recv(onData)
		ms.SyncSend([]byte(strings.Repeat("s", RecvBuffSize)))
	}

	<-c

	s.ln.Close()

}

func Test7(t *testing.T) {
	s := newServer()
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	var onData func(ms *MuxStream, data []byte)

	onData = func(ms *MuxStream, data []byte) {
		count := ms.GetUserData().(*int)
		ms.Recv(onData)
		(*count)++
		if *count < 10 {
			msg := make([]byte, len(data))
			copy(msg, data)
			ms.SyncSend(msg)
		}
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	c := make(chan struct{})

	ms, err := socket.Dial(time.Second * 5)
	if nil != err {
		fmt.Printf("dial error:%v", err)
	} else {

		ms.SetCloseCallback(func(_ *MuxStream, err error) {
			fmt.Println("client muxstream close", err)
			close(c)
		})

		count := 0
		ms.SetUserData(&count)
		ms.Recv(onData)
		ms.SyncSend([]byte(strings.Repeat("s", m64K)))
	}

	<-c

	s.ln.Close()

}

func Test8(t *testing.T) {
	s := newServer()
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	var onData func(ms *MuxStream, data []byte)

	onData = func(ms *MuxStream, data []byte) {
		count := ms.GetUserData().(*int)
		ms.Recv(onData)
		(*count)++
		if *count < 10 {
			msg := make([]byte, len(data))
			copy(msg, data)
			ms.SyncSend(msg)
		}
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	c := make(chan struct{})

	ms, err := socket.Dial(time.Second * 5)
	if nil != err {
		fmt.Printf("dial error:%v", err)
	} else {

		ms.SetCloseCallback(func(_ *MuxStream, err error) {
			fmt.Println("client muxstream close", err)
			close(c)
		})

		count := 0
		ms.SetUserData(&count)
		ms.Recv(onData)
		ms.SyncSend([]byte(strings.Repeat("s", m512K)))
	}

	<-c

	s.ln.Close()

}

func Test9(t *testing.T) {
	s := newServer()
	s.dontRecv = true
	s.serve("localhost:8110")
	time.Sleep(time.Second)

	conn, err := net.Dial("tcp", "localhost:8110")
	if err != nil {
		t.Fatal(err)
	}

	socket := NewMuxSocketClient(conn, enc, nil)

	count := 100

	wait := sync.WaitGroup{}
	wait.Add(count)
	for j := 0; j < count; j++ {
		go func() {
			ms, err := socket.Dial(time.Second * 5)
			if nil != err {
				wait.Done()
			} else {
				ms.Close(nil)
				wait.Done()
			}
		}()
	}
	wait.Wait()
	//}

	s.ln.Close()
}
