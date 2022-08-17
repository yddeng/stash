package types

import (
	"fmt"
	codecs "initialthree/codec/cs"
	"time"

	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/pkg/event"
)

type TimerID string

func (id TimerID) String() string {
	return string(id)
}

var timerIDMap = map[TimerID]struct{}{}

func NewTimerID(s string) TimerID {
	tid := TimerID(s)
	if _, ok := timerIDMap[tid]; ok {
		panic(fmt.Errorf("TimerID %s already exist", s))
	}

	timerIDMap[tid] = struct{}{}
	return tid
}

type Module interface {
	OnModuleSync(Robot, proto.Message)
}

type Robot interface {
	EventQueue() *event.EventQueue
	SetSession(*fnet.Socket)
	CloseSession()
	SendMessage(proto.Message, ...func(Robot, *codecs.Message) bool)
	UserID() string
	RoleName() string
	GetStatusValue(StatusID) StatusValue
	SetStatusValue(StatusID, StatusValue)
	IsStatus(StatusID) bool
	SetStatus(StatusID)
	UnsetStatus(StatusID)
	GetModule(int) Module
	AddTimer(TimerID, time.Duration, interface{}, func(Robot, interface{}))
	RemTimer(TimerID)
	PostNoWait(fn interface{}, arg ...interface{})
	Debugf(string, ...interface{})
	Debug(...interface{})
	Infof(string, ...interface{})
	Info(...interface{})
	Errorf(string, ...interface{})
	Error(...interface{})
	Fatalf(string, ...interface{})
	Fatal(...interface{})
	Panic(...interface{})
	Panicf(string, ...interface{})
}
