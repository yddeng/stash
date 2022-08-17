package smux

import (
	"container/list"
	"encoding/binary"
	"errors"
	"github.com/sniperHW/flyfish/pkg/buffer"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var (
	Err_Close       = errors.New("stream closed")
	Err_SendTimeout = errors.New("send timeout")
	Err_DialTimeout = errors.New("dial timeout")
	Err_DialFail    = errors.New("dial fail")
	Err_RecvTimeout = errors.New("recv timeout")
	Err_NoFreeID    = errors.New("no free id")
	Err_Eof         = errors.New("EOF")
)

const (
	EOF  uint16 = 1
	SYN  uint16 = 2
	EST  uint16 = 3
	WANT uint16 = 4
	DATA uint16 = 5
)

const (
	MuxServer = 1
	MuxClient = 2
)

/*
 *    | type 2字节 |stream id 2字节| data len 4字节 |  data |
 */

type stSyn struct {
	ch chan interface{}
}

type stSend struct {
	s       *MuxStream
	o       interface{}
	timer   *time.Timer
	cb      func(error)
	listEle *list.Element
}

type stWant struct {
	onData   func(*MuxStream, []byte)
	deadline time.Time
}

type MuxSocket struct {
	mu            sync.Mutex
	streams       map[int64]*MuxStream
	pendingEst    map[int64]stSyn
	sendQueue     *BlockQueue
	conn          net.Conn
	enc           func(o interface{}, b *buffer.Buffer) error
	onSocketClose func(*MuxSocket)
	closed        int32
	recvTimeout   time.Duration
	mode          int
	onNewStream   func(*MuxStream)
	idCounter     int64
	dailOnce      sync.Once
	listenOnce    sync.Once
}

func NewMuxSocketServer(conn net.Conn, recvTimeout time.Duration, enc func(o interface{}, b *buffer.Buffer) error, onSocketClose func(*MuxSocket)) *MuxSocket {
	s := &MuxSocket{
		streams:       map[int64]*MuxStream{},
		conn:          conn,
		onSocketClose: onSocketClose,
		enc:           enc,
		sendQueue:     NewBlockQueue(),
		recvTimeout:   recvTimeout,
		mode:          MuxServer,
	}

	runtime.SetFinalizer(s, func(s *MuxSocket) {
		s.Close()
	})

	return s
}

func NewMuxSocketClient(conn net.Conn, recvTimeout time.Duration, enc func(o interface{}, b *buffer.Buffer) error, onSocketClose func(*MuxSocket)) *MuxSocket {
	s := &MuxSocket{
		streams:       map[int64]*MuxStream{},
		conn:          conn,
		onSocketClose: onSocketClose,
		enc:           enc,
		sendQueue:     NewBlockQueue(),
		recvTimeout:   recvTimeout,
		mode:          MuxClient,
		pendingEst:    map[int64]stSyn{},
	}

	runtime.SetFinalizer(s, func(s *MuxSocket) {
		s.Close()
	})

	return s
}

func (s *MuxSocket) Listen(onNewStream func(*MuxStream)) error {
	if s.mode != MuxServer {
		return errors.New("not smux server")
	}

	if nil == onNewStream {
		return errors.New("onNewStream = nil")
	}

	s.listenOnce.Do(func() {
		s.serve(onNewStream)
	})

	return nil
}

func (s *MuxSocket) Dial(timeout time.Duration) (*MuxStream, error) {
	if s.mode != MuxClient {
		return nil, errors.New("not smux client")
	}

	s.dailOnce.Do(func() {
		s.serve(nil)
	})

	s.mu.Lock()
	if len(s.streams)+len(s.pendingEst) >= 65535 {
		s.mu.Unlock()
		return nil, errors.New("busy")
	}

	ch := make(chan interface{}, 1)
	s.idCounter++
	id := s.idCounter
	s.pendingEst[id] = stSyn{
		ch: ch,
	}
	s.mu.Unlock()

	s.sendQueue.Add(s.makeSyn(id))

	ticker := time.NewTicker(timeout)

	select {
	case e := <-ch:
		ticker.Stop()
		switch e.(type) {
		case error:
			return nil, e.(error)
		default:
			return e.(*MuxStream), nil
		}
	case <-ticker.C:
		s.mu.Lock()
		delete(s.pendingEst, id)
		s.mu.Unlock()
		//通告对端关闭
		msg := make([]byte, 0, 10)
		msg = buffer.AppendUint16(msg, EOF)
		s.sendQueue.Add(buffer.AppendInt64(msg, id))
		return nil, Err_DialTimeout
	}
}

func (s *MuxSocket) serve(onNewStream func(*MuxStream)) {
	s.onNewStream = onNewStream
	go s.readloop()
	go s.sendloop()
}

func (s *MuxSocket) Close() {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		runtime.SetFinalizer(s, nil)
		s.mu.Lock()
		streams := s.streams
		pendingEst := s.pendingEst
		s.streams = nil
		s.pendingEst = nil
		s.mu.Unlock()

		s.sendQueue.Close()

		for _, v := range streams {
			v.onEof()
		}

		s.conn.Close()

		for _, v := range pendingEst {
			v.ch <- Err_DialFail
		}

		if nil != s.onSocketClose {
			s.onSocketClose(s)
		}
	}
}

func (s *MuxSocket) onEof(id int64) {
	s.mu.Lock()
	ss := s.streams[id]
	if nil != ss {
		delete(s.streams, id)
		s.mu.Unlock()
		ss.onEof()
	} else {
		s.mu.Unlock()
	}
}

func (s *MuxSocket) createMuxStream(id int64) *MuxStream {
	ss := &MuxStream{
		id:        id,
		want:      list.New(),
		socket:    s,
		sendQueue: list.New(),
	}

	runtime.SetFinalizer(ss, func(ss *MuxStream) {
		ss.Close(Err_Close)
	})

	return ss
}

func (s *MuxSocket) onSyn(id int64) {
	if s.mode == MuxServer {
		s.mu.Lock()
		ss := s.streams[id]
		if nil == ss {
			ss = s.createMuxStream(id)
			s.streams[id] = ss
			s.mu.Unlock()
			s.sendQueue.Add(s.makeEst(id))
			s.onNewStream(ss)
		} else {
			s.mu.Unlock()
		}
	}
}

func (s *MuxSocket) onEst(id int64) {
	if s.mode == MuxClient {
		s.mu.Lock()
		defer s.mu.Unlock()
		st, ok := s.pendingEst[id]
		if ok {
			delete(s.pendingEst, id)
			ss := s.createMuxStream(id)
			s.streams[id] = ss
			st.ch <- ss
		}
	}
}

func (s *MuxSocket) onWantReq(id int64) {
	s.mu.Lock()
	ss := s.streams[id]
	s.mu.Unlock()
	if nil != ss {
		ss.onWantReq()
	}
}

func (s *MuxSocket) onData(id int64, b []byte) {
	s.mu.Lock()
	ss := s.streams[id]
	s.mu.Unlock()
	if nil != ss {
		ss.onData(b)
	}
}

func (s *MuxSocket) makeSyn(id int64) []byte {
	msg := make([]byte, 0, 10)
	msg = buffer.AppendUint16(msg, SYN)
	return buffer.AppendInt64(msg, id)
}

func (s *MuxSocket) makeEst(id int64) []byte {
	msg := make([]byte, 0, 10)
	msg = buffer.AppendUint16(msg, EST)
	return buffer.AppendInt64(msg, id)
}

func (s *MuxSocket) sendEOF(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.streams[id]; ok {
		delete(s.streams, id)
		msg := make([]byte, 0, 10)
		msg = buffer.AppendUint16(msg, EOF)
		s.sendQueue.Add(buffer.AppendInt64(msg, id))
	}
}

func (s *MuxSocket) sendWantReq(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.streams[id]; ok {
		msg := make([]byte, 0, 10)
		msg = buffer.AppendUint16(msg, WANT)
		s.sendQueue.Add(buffer.AppendInt64(msg, id))
		return nil
	} else {
		return Err_Close
	}
}

func (s *MuxSocket) send(st *stSend) error {
	s.mu.Lock()
	if nil == s.streams[st.s.id] {
		s.mu.Unlock()
		return Err_Close
	} else {
		s.mu.Unlock()
	}
	return s.sendQueue.Add(st)
}

var RecvBuffSize int = 1024 * 4

const (
	m64K  = 64 * 1024
	m512K = 512 * 1024
)

func getPoolIndex(cap int) int {
	if cap <= m64K {
		return 0
	} else if cap <= m512K {
		return 1
	} else {
		return -1
	}
}

type recvBufferPool struct {
	fixedBufferPool []*sync.Pool
}

func (this *recvBufferPool) get(cap int) []byte {
	i := getPoolIndex(cap)
	if i < 0 {
		return make([]byte, cap)
	} else {
		return this.fixedBufferPool[i].Get().([]byte)
	}
}

func (this *recvBufferPool) put(b []byte) {
	if i := getPoolIndex(cap(b)); i >= 0 {
		this.fixedBufferPool[i].Put(b)
	}
}

var gRecvBufferPool recvBufferPool = recvBufferPool{
	fixedBufferPool: []*sync.Pool{
		&sync.Pool{
			New: func() interface{} {
				return make([]byte, m64K)
			},
		},
		&sync.Pool{
			New: func() interface{} {
				return make([]byte, m512K)
			},
		},
	},
}

func (s *MuxSocket) readloop() {
	recvbuff := make([]byte, RecvBuffSize)
	pRecvBuff := &recvbuff
	r := 0
	w := 0
	for {
		if s.recvTimeout > 0 {
			if err := s.conn.SetReadDeadline(time.Now().Add(s.recvTimeout)); err != nil {
				s.Close()
				return
			}
		}
		n, err := s.conn.Read((*pRecvBuff)[w:])
		if nil != err {
			s.Close()
			return
		} else {
			l := 0
			w += n
			for d := w - r; d >= 10; d = w - r {
				tt := binary.BigEndian.Uint16((*pRecvBuff)[r:])
				id := int64(binary.BigEndian.Uint64((*pRecvBuff)[r+2:]))
				if tt == DATA {
					if d >= 14 {
						totalSize := int(binary.BigEndian.Uint32((*pRecvBuff)[r+10:])) + 14
						if d >= totalSize {
							s.onData(id, (*pRecvBuff)[r+14:r+totalSize])
							r += totalSize
						} else {
							l = totalSize
							break
						}
					} else {
						break
					}
				} else {
					switch tt {
					case EOF:
						s.onEof(id)
					case SYN:
						s.onSyn(id)
					case EST:
						s.onEst(id)
					case WANT:
						s.onWantReq(id)
					default:
						s.Close()
						return
					}
					r += 10
				}
			}

			if w == r {
				if pRecvBuff != &recvbuff {
					gRecvBufferPool.put(*pRecvBuff)
				}
				pRecvBuff = &recvbuff
			} else {
				if l > cap(*pRecvBuff) {
					buffer := gRecvBufferPool.get(l)

					copy(buffer, (*pRecvBuff)[r:w])
					pRecvBuff = &buffer
				} else {
					copy(*pRecvBuff, (*pRecvBuff)[r:w])
				}
			}

			w = w - r
			r = 0
		}
	}
}

var MaxSendSize int = 1024 * 64

func (s *MuxSocket) sendloop() {

	localList := make([]interface{}, 0, 32)

	var closed bool

	var err error

	for {
		closed, localList = s.sendQueue.Get(localList)
		size := len(localList)
		if closed && size == 0 {
			break
		}

		b := buffer.Get()
		for i := 0; i < size; {
			for i < size {
				err = nil
				l := b.Len()

				switch localList[i].(type) {
				case []byte:
					b.AppendBytes(localList[i].([]byte))
				default:

					st := localList[i].(*stSend)
					stopTimerOK := false
					if nil == st.timer || st.timer.Stop() {
						stopTimerOK = true
					}

					if stopTimerOK {
						b.AppendUint16(DATA)
						b.AppendInt64(st.s.id)
						ll := b.Len()
						b.AppendUint32(0)
						if err = s.enc(st.o, b); nil == err {
							dataLen := b.Len() - ll - 4
							b.SetUint32(ll, uint32(dataLen))
						} else {
							//EnCode错误，这个包已经写入到b中的内容需要直接丢弃
							b.SetLen(l)
						}
						st.cb(err)
						if nil != err {
							st.s.onWantReq()
						}
					} else {
						st.s.onWantReq()
					}
				}

				localList[i] = nil
				i++

				if b.Len() >= MaxSendSize {
					break
				}
			}

			if b.Len() == 0 {
				b.Free()
				break
			}

			_, err = s.conn.Write(b.Bytes())

			if nil == err {
				b.Reset()
			} else {
				b.Free()
				s.Close()
				return
			}
		}
	}
}

type MuxStream struct {
	muRecv      sync.Mutex
	id          int64
	closed      int32
	want        *list.List
	recvtimeout time.Duration
	timer       *time.Timer
	//
	muSend    sync.Mutex
	wantReq   int
	sendQueue *list.List
	//
	onClose func(*MuxStream, error)

	muUd   sync.Mutex
	ud     interface{}
	socket *MuxSocket
}

func (s *MuxStream) SetRecvTimeout(timeout time.Duration) {
	s.muRecv.Lock()
	defer s.muRecv.Unlock()

	if nil != s.timer {
		s.timer.Stop()
	}

	s.recvtimeout = timeout

	now := time.Now()

	for e := s.want.Front(); nil != e; e = e.Next() {
		if timeout != 0 {
			e.Value.(*stWant).deadline = now.Add(timeout)
		} else {
			e.Value.(*stWant).deadline = time.Time{}
		}
	}

	if timeout != 0 && s.want.Len() != 0 {
		s.timer = time.AfterFunc(timeout, s.onTimeout)
	}
}

func (s *MuxStream) onTimeout() {
	if atomic.LoadInt32(&s.closed) == 1 {
		return
	}

	now := time.Now()
	s.muRecv.Lock()
	for e := s.want.Front(); nil != e; {
		if now.After(e.Value.(*stWant).deadline) {
			s.muRecv.Unlock()
			s.Close(Err_RecvTimeout)
			return
		} else {
			break
		}
	}

	e := s.want.Front()

	if nil != e {
		s.timer = time.AfterFunc(e.Value.(*stWant).deadline.Sub(now), s.onTimeout)
	} else {
		s.timer = nil
	}

	s.muRecv.Unlock()
}

func (s *MuxStream) onEof() {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		runtime.SetFinalizer(s, nil)

		s.muSend.Lock()

		if nil != s.timer {
			s.timer.Stop()
		}

		for n := s.sendQueue.Front(); nil != n; n = s.sendQueue.Front() {
			v := s.sendQueue.Remove(n).(*stSend)
			if nil != v.timer {
				v.timer.Stop()
			}
			v.cb(Err_Eof)
		}

		s.muSend.Unlock()

		s.muRecv.Lock()
		for f := s.want.Front(); nil != f; f = s.want.Front() {
			s.want.Remove(f)
		}
		s.muRecv.Unlock()

		if nil != s.onClose {
			s.onClose(s, Err_Eof)
		}
	}
}

func (s *MuxStream) Recv(onData func(*MuxStream, []byte)) error {
	if atomic.LoadInt32(&s.closed) == 1 {
		return Err_Close
	}

	s.muRecv.Lock()
	defer s.muRecv.Unlock()

	st := &stWant{
		onData: onData,
	}

	if s.recvtimeout != 0 {
		st.deadline = time.Now().Add(s.recvtimeout)
		if nil == s.timer {
			s.timer = time.AfterFunc(s.recvtimeout, s.onTimeout)
		}
	}

	s.want.PushBack(st)

	return s.socket.sendWantReq(s.id)
}

func (s *MuxStream) onData(data []byte) {
	if atomic.LoadInt32(&s.closed) == 1 {
		return
	}

	s.muRecv.Lock()

	if s.want.Len() == 0 {
		s.muRecv.Unlock()
		return
	}

	f := s.want.Front()
	s.want.Remove(f)
	s.muRecv.Unlock()
	f.Value.(*stWant).onData(s, data)
}

//异步发送
func (s *MuxStream) AsyncSend(o interface{}, cb func(*MuxStream, error), timeout ...time.Duration) error {

	if atomic.LoadInt32(&s.closed) == 1 {
		return Err_Close
	}

	s.muSend.Lock()
	defer s.muSend.Unlock()

	st := &stSend{
		s: s,
		o: o,
	}

	st.cb = func(err error) {
		if atomic.LoadInt32(&s.closed) == 0 {
			if err == Err_SendTimeout {
				s.muSend.Lock()
				if nil != st.listEle {
					s.sendQueue.Remove(st.listEle)
				}
				s.muSend.Unlock()
			}
			if nil != cb {
				cb(s, err)
			}
		}
	}

	if len(timeout) > 0 && timeout[0] > 0 {
		st.timer = time.AfterFunc(timeout[0], func() {
			st.cb(Err_SendTimeout)
		})
	}

	if s.wantReq > 0 {
		s.wantReq--
		return s.socket.send(st)
	} else {
		st.listEle = s.sendQueue.PushBack(st)
		return nil
	}
}

//同步发送
func (s *MuxStream) SyncSend(o interface{}, timeout ...time.Duration) error {
	if atomic.LoadInt32(&s.closed) == 1 {
		return Err_Close
	}

	s.muSend.Lock()

	var err error

	ch := make(chan struct{})

	st := &stSend{
		s: s,
		o: o,
	}

	st.cb = func(e error) {
		err = e
		close(ch)
		if atomic.LoadInt32(&s.closed) == 0 {
			if err == Err_SendTimeout {
				s.muSend.Lock()
				if nil != st.listEle {
					s.sendQueue.Remove(st.listEle)
				}
				s.muSend.Unlock()
			}
		}
	}

	if len(timeout) > 0 && timeout[0] > 0 {
		st.timer = time.AfterFunc(timeout[0], func() {
			st.cb(Err_SendTimeout)
		})
	}

	if s.wantReq > 0 {
		s.wantReq--
		s.muSend.Unlock()
		return s.socket.send(st)
	} else {
		st.listEle = s.sendQueue.PushBack(st)
		s.muSend.Unlock()
	}

	<-ch

	return err
}

func (s *MuxStream) onWantReq() {
	if atomic.LoadInt32(&s.closed) == 0 {
		s.muSend.Lock()
		defer s.muSend.Unlock()
		s.wantReq++
		for s.wantReq > 0 {
			f := s.sendQueue.Front()
			if nil != f {
				v := s.sendQueue.Remove(f).(*stSend)
				v.listEle = nil
				if err := s.socket.send(v); nil != err {
					return
				}
				s.wantReq--
			} else {
				return
			}
		}
	}
}

func (s *MuxStream) Close(err error) {
	if atomic.CompareAndSwapInt32(&s.closed, 0, 1) {
		s.socket.sendEOF(s.id)

		runtime.SetFinalizer(s, nil)

		s.muSend.Lock()

		if nil != s.timer {
			s.timer.Stop()
		}

		for n := s.sendQueue.Front(); nil != n; n = s.sendQueue.Front() {
			v := s.sendQueue.Remove(n).(*stSend)
			if nil != v.timer {
				v.timer.Stop()
			}
			v.cb(err)
		}

		s.muSend.Unlock()

		s.muRecv.Lock()
		for f := s.want.Front(); nil != f; f = s.want.Front() {
			s.want.Remove(f)
		}
		s.muRecv.Unlock()

		if nil != s.onClose {
			s.onClose(s, err)
		}
	}
}

func (s *MuxStream) SetCloseCallback(cb func(*MuxStream, error)) {
	s.onClose = cb
}

func (s *MuxStream) SetUserData(ud interface{}) {
	s.muUd.Lock()
	defer s.muUd.Unlock()
	s.ud = ud
}

func (s *MuxStream) GetUserData() interface{} {
	s.muUd.Lock()
	defer s.muUd.Unlock()
	return s.ud
}

func (s *MuxStream) ID() int64 {
	return s.id
}
