package hutils

import (
	"sync"
	"time"
)

type Watcher struct {
	lock  sync.Mutex
	waits []*WaitTimeout
}

func NewWatcher() *Watcher {
	return &Watcher{waits: []*WaitTimeout{}}
}

func (w *Watcher) Watch(timeout time.Duration) chan bool {
	waiter := NewWaitTimeout().Add(1)
	w.addWaiter(waiter)
	return waiter.Wait(timeout)
}

func (w *Watcher) addWaiter(waiter *WaitTimeout) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.waits = append(w.waits, waiter)
}

func (w *Watcher) Notify() {
	w.lock.Lock()
	defer w.lock.Unlock()
	for _, waiter := range w.waits {
		waiter.Done()
	}
	w.waits = []*WaitTimeout{}
}
