package node_gate

import (
	"fmt"
	"initialthree/codec/cs"
	"initialthree/common"
	"initialthree/network/smux"
	"initialthree/node/node_gate/monitor"
	"initialthree/protocol/cmdEnum"
	"initialthree/zaplogger"
	"net"
	"sync/atomic"
	"time"
)

func onNewClient(conn net.Conn) {
	g, stream, userID, err := onUserLogin(conn)
	if err != nil {
		monitor.UserLoginFailedInc()
		zaplogger.GetSugar().Infof("onUserLogin user(%s) failed, %s", userID, err)
		conn.Close()
		if stream != nil {
			stream.Close(err)
		}
		return
	}

	monitor.UserLoginSuccessInc()

	zaplogger.GetSugar().Infof("new newConnection user(%s) login ok addBIO stream(%d)", userID, stream.ID())
	g.addChannel(stream.ID(), newChannel(stream, conn, userID))
}

type channel struct {
	userID  string
	stream  *smux.MuxStream
	tcpConn net.Conn
	tcpAddr *net.TCPAddr

	closed    int32
	chClose   chan struct{}
	closeFunc func(err error)
}

func newChannel(stream *smux.MuxStream, tcpConn net.Conn, userID string) *channel {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", tcpConn.RemoteAddr().String())
	return &channel{stream: stream, tcpConn: tcpConn, tcpAddr: tcpAddr, chClose: make(chan struct{}), userID: userID}
}

func (this *channel) close(err error) {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		close(this.chClose)
		this.stream.Close(err)
		this.tcpConn.Close()
		this.closeFunc(err)
	}
}

func (this *channel) run(closeFunc func(err error)) {
	this.closeFunc = closeFunc
	this.stream.SetCloseCallback(func(stream *smux.MuxStream, e error) {
		this.close(fmt.Errorf("stream %d closed %s. ", stream.ID(), e.Error()))
	})

	go this.stream2tcp()
	go this.tcp2stream()
}

func (this *channel) sendToTcp(bytes []byte) {
	_ = this.tcpConn.SetWriteDeadline(time.Now().Add(time.Second))
	_, _ = this.tcpConn.Write(bytes)
}

func (this *channel) stream2tcp() {
	ch := make(chan []byte, 10)
	for {
		this.stream.Recv(func(stream *smux.MuxStream, bytes []byte) {
			b := make([]byte, len(bytes))
			copy(b, bytes)

			select {
			case ch <- b:
			default:
			}
		})

		select {
		case <-this.chClose:
			return
		case bytes := <-ch:

			_, cmd, _ := cs.FetchSeqCmdCode(bytes)
			if cmd == cmdEnum.CS_ChatMessageSync {
				broadcast(bytes)
				break
			}

			_ = this.tcpConn.SetWriteDeadline(time.Now().Add(time.Second * time.Duration(5)))
			_, err := this.tcpConn.Write(bytes)
			if err != nil {
				this.close(err)
				return
			}
		}
	}
}

func (this *channel) tcp2stream() {
	for {
		_ = this.tcpConn.SetReadDeadline(time.Now().Add(common.HeartBeat_Timeout_Client))
		data, err := cs.ReadMessage(this.tcpConn)
		if err != nil {
			this.close(err)
			return
		}

		monitor.UserRequestInc()

		if err = this.stream.SyncSend(data); err != nil {
			this.close(err)
			return
		}

		select {
		case <-this.chClose:
			return
		default:
		}
	}
}
