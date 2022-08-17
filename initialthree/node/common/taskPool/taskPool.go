package taskPool

import (
	"errors"
	"initialthree/pkg/pcall"
	"runtime"
	"sync"
)

const (
	defaultTaskSize = 1024
)

type TaskPool struct {
	workers    int
	workerSize int
	workerLock sync.Mutex

	taskChan chan *funcTask

	die     chan struct{}
	dieOnce sync.Once
}

func (this *TaskPool) Running() int {
	this.workerLock.Lock()
	defer this.workerLock.Unlock()
	return this.workers
}

func (this *TaskPool) AddTask(fn interface{}, args ...interface{}) error {
	select {
	case <-this.die:
		return errors.New("taskPool:AddTask pool is stopped")
	default:
	}

	task := &funcTask{fn: fn, args: args}

	var taskChan chan *funcTask
	if this.workerSize == 0 {
		taskChan = make(chan *funcTask, 1)
	} else {
		taskChan = this.taskChan
	}

	select {
	case taskChan <- task:
	default:
		return errors.New("taskPool:AddTask task channel is full")
	}

	this.workerLock.Lock()
	defer this.workerLock.Unlock()

	if this.workerSize == 0 || this.workers < this.workerSize {
		this.workers++
		this.goWorker(taskChan)
	}
	return nil
}

func (this *TaskPool) Stop() {
	this.dieOnce.Do(func() {
		close(this.die)
	})
}

// NewTaskPool
// workerSize > 0 , 限制goroutine的数量; workerSize = 0 , 不限制
func NewTaskPool(workerSize, taskSize int) *TaskPool {
	if taskSize < defaultTaskSize {
		taskSize = defaultTaskSize
	}
	if workerSize < 0 {
		workerSize = 0
	}

	pool := new(TaskPool)
	pool.die = make(chan struct{})
	pool.workerSize = workerSize
	pool.taskChan = make(chan *funcTask, taskSize)

	return pool
}

type funcTask struct {
	fn   interface{}
	args []interface{}
}

func (this *TaskPool) goWorker(taskC chan *funcTask) {
	go func() {
		defer func() {
			this.workerLock.Lock()
			this.workers--
			this.workerLock.Unlock()
		}()

		for {
			select {
			case task := <-taskC:
				_, _ = pcall.Call(task.fn, task.args...)
			default:
				return
			}
		}
	}()
}

var (
	defaultTaskPool *TaskPool
	createOnce      sync.Once
)

func AddTask(fn interface{}, args ...interface{}) error {
	createOnce.Do(func() {
		defaultTaskPool = NewTaskPool(runtime.NumCPU()*2, defaultTaskSize)
	})
	return defaultTaskPool.AddTask(fn, args...)
}
