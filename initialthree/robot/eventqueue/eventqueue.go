package eventqueue

import (
	"initialthree/pkg/event"
	"math/rand"
	"runtime"
	"sync"
)

var mtx sync.RWMutex
var status int32
var eventQueues []*event.EventQueue
var eventQueueCount uint

func Init() {
	mtx.Lock()
	defer mtx.Unlock()

	if status != 0 {
		return
	}

	eventQueueCount = uint(runtime.NumCPU())
	eventQueues = make([]*event.EventQueue, eventQueueCount)
	for i := range eventQueues {
		eventQueues[i] = event.NewEventQueue()
	}
	status = 1
}

func Start() {
	mtx.Lock()
	defer mtx.Unlock()

	if status != 1 {
		return
	}

	for i := range eventQueues {
		go eventQueues[i].Run()
	}
	status = 2
}

func Stop() {
	mtx.Lock()

	if status == 0 {
		mtx.Unlock()
		return
	}

	if status == 2 {
		status = 3
		mtx.Unlock()

		wait := &sync.WaitGroup{}
		for _, v := range eventQueues {
			wait.Add(1)
			v.PostNoWait(1, func() { wait.Done() })
			v.Close()
		}
		wait.Wait()

		mtx.Lock()
		defer mtx.Unlock()
		eventQueues = nil
		eventQueueCount = 0
		status = 0
	} else if status == 1 {
		eventQueues = nil
		eventQueueCount = 0
		status = 0
	}
}

func EventQueueCount() uint {
	mtx.RLock()
	defer mtx.RUnlock()

	return eventQueueCount
}

func ModEventQueue(n ...uint) *event.EventQueue {
	mtx.RLock()
	defer mtx.RUnlock()

	if status == 1 || status == 2 {
		var nn uint
		if len(n) > 0 {
			nn = n[0]
		} else {
			nn = uint(rand.Int())
		}

		return eventQueues[nn%eventQueueCount]
	}

	return nil
}

func ModPostNoWait(fn interface{}, args []interface{}, n ...uint) {
	if eventQueue := ModEventQueue(n...); eventQueue != nil {
		eventQueue.PostNoWait(1, fn, args...)
	}
}
