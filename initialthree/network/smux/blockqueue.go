package smux

import (
	"errors"
	"sync"
)

var (
	ErrQueueClosed = errors.New("queue closed")
)

const (
	initCap = 64
)

type BlockQueue struct {
	list        []interface{}
	listGuard   sync.Mutex
	emptyCond   *sync.Cond
	closed      bool
	emptyWaited int
}

func (self *BlockQueue) Add(item interface{}) error {
	self.listGuard.Lock()
	if self.closed {
		self.listGuard.Unlock()
		return ErrQueueClosed
	}

	self.list = append(self.list, item)

	needSignal := self.emptyWaited > 0
	self.listGuard.Unlock()
	if needSignal {
		self.emptyCond.Signal()
	}
	return nil
}

func (self *BlockQueue) Get(swaped []interface{}) (closed bool, datas []interface{}) {
	swaped = swaped[0:0]
	self.listGuard.Lock()
	for !self.closed && len(self.list) == 0 {
		self.emptyWaited++
		self.emptyCond.Wait()
		self.emptyWaited--
	}
	datas = self.list
	closed = self.closed
	self.list = swaped
	self.listGuard.Unlock()
	return
}

func (self *BlockQueue) Close() bool {
	self.listGuard.Lock()

	if self.closed {
		self.listGuard.Unlock()
		return false
	}

	self.closed = true
	self.listGuard.Unlock()
	self.emptyCond.Broadcast()

	return true
}

func NewBlockQueue() *BlockQueue {
	self := &BlockQueue{}
	self.closed = false
	self.emptyCond = sync.NewCond(&self.listGuard)
	self.list = make([]interface{}, 0, initCap)
	return self
}
