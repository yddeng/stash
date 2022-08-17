package transaction

//go test -covermode=count -v -coverprofile=coverage.out -run=.
//go tool cover -html=coverage.out

import (
	"fmt"
	"initialthree/pkg/timer"
	"testing"
	"time"
)

func AsynFunc1(args string, callback func(string)) {
	timer.Once(time.Second*time.Duration(3), func(*timer.Timer, interface{}) {
		callback("AsynFunc1 " + args)
	}, nil)
}

type transactionTest1 struct {
	TransactionBase
}

func (this *transactionTest1) GetModuleName() string {
	return "test1"
}

func (this *transactionTest1) Begin() {
	this.AsynWrap(AsynFunc1)("test1", func(msg string) {
		fmt.Println("AsynFunc1 callback", msg)
		this.EndTrans()
	})
}

func (this *transactionTest1) End() {
}

func (this *transactionTest1) Timeout() {
	fmt.Println("transactionTest1 timeout")
}

type st struct {
	data string
}

func (this *st) AsynFunc2(args string, callback func(string)) {
	timer.Once(time.Second*time.Duration(1), func(*timer.Timer, interface{}) {
		callback("AsynFunc2 " + args + " " + this.data)
	}, nil)
}

type transactionTest2 struct {
	TransactionBase
	mgr *TransactionMgr
}

func (this *transactionTest2) GetModuleName() string {
	return "test2"
}

func (this *transactionTest2) Begin() {

	st_ := &st{
		data: "i'm st",
	}

	this.AsynWrap(st_.AsynFunc2)("test2", func(msg string) {
		fmt.Println("AsynFunc2 callback", msg)
		this.EndTrans()
	})
}

func (this *transactionTest2) End() {
}

func (this *transactionTest2) Timeout() {
	fmt.Println("transactionTest2 timeout")
}

type transactionTest3 struct {
	TransactionBase
}

func (this *transactionTest3) GetModuleName() string {
	return "test3"
}

func (this *transactionTest3) Begin() {
	this.AsynWrap(AsynFunc1)("test1", func(msg int, v int) {
		fmt.Println("AsynFunc1 callback", msg)
		this.EndTrans()
	})
}

func (this *transactionTest3) End() {
}

func (this *transactionTest3) Timeout() {
	fmt.Println("transactionTest3 timeout")
}

func TestTrans(t *testing.T) {
	mgr := New(nil)

	trans1 := &transactionTest1{}
	trans1.setTransBase(mgr, trans1)

	mgr.PushTrans(trans1, time.Second*time.Duration(2))

	for nil != mgr.current {
		mgr.Tick(time.Now())
		time.Sleep(time.Second)
	}

	trans2 := &transactionTest2{}
	trans2.setTransBase(mgr, trans2)

	mgr.PushTrans(trans2, time.Second*time.Duration(2))
	for nil != mgr.current {
		mgr.Tick(time.Now())
		time.Sleep(time.Second)
	}

	//这里会触发panic

	trans3 := &transactionTest3{}
	trans3.setTransBase(mgr, trans3)

	mgr.PushTrans(trans3, time.Second*time.Duration(4))

	for nil != mgr.current {
		mgr.Tick(time.Now())
		time.Sleep(time.Second)
	}

}
