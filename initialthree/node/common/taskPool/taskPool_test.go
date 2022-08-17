package taskPool

import (
	"testing"
	"time"
)

func TestNewTaskPool(t *testing.T) {
	p := NewTaskPool(0, 100)
	p.AddTask(func() { t.Log("f1") })
	p.AddTask(func() { t.Log("f2") })
	p.AddTask(func() { t.Log("f3") })
	p.AddTask(func() { t.Log("f4") })
	p.AddTask(func() { t.Log("f5") })

	t.Log(p.Running())
	time.Sleep(time.Millisecond)
	t.Log(p.Running())

	p.AddTask(func() { t.Log("f6") })
	p.AddTask(func() { t.Log("f7") })
	p.AddTask(func() { t.Log("f8") })

	t.Log(p.Running())
	time.Sleep(time.Millisecond)
	t.Log(p.Running())

	p.Stop()
	t.Log(p.AddTask(func() { t.Log("f10") }))
}

func TestNewTaskPool2(t *testing.T) {
	p := NewTaskPool(2, 100)
	p.AddTask(func() { t.Log("f1") })
	p.AddTask(func() { t.Log("f2") })
	p.AddTask(func() { t.Log("f3") })
	p.AddTask(func() { t.Log("f4") })
	p.AddTask(func() { t.Log("f5") })
	t.Log(p.Running())

	time.Sleep(time.Millisecond)
	t.Log(p.Running())

	p.Stop()
	t.Log(p.AddTask(func() { t.Log("f10") }))
}
