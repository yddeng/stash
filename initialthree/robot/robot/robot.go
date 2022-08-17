package robot

import (
	"container/list"
	"initialthree/pkg/event"
	"initialthree/robot/behavior"
	"initialthree/robot/robot/module"
	"initialthree/robot/types"
	"time"

	"github.com/GodYY/bevtree"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"go.uber.org/zap"
)

type Module = types.Module
type RobotI = types.Robot

type Robot struct {
	userID   string
	roleName string

	eventQueue *event.EventQueue

	status types.Status

	bevEntity bevtree.Entity

	logger *zap.SugaredLogger

	// 网络会话
	session *fnet.Socket
	// 消息序号
	msgSeriNo uint32
	// 上一次发送消息的时间
	lastSendMsgTime time.Time
	// 回调列表，用以实现类似RPC的功能
	// 正常情况下不会存在太多RPC等待，并且消息是按顺序回调的，
	// 固使用链表来维护消息回调
	msgCBList *list.List

	// 数据模块
	modules []Module

	timers map[TimerID]*timer
}

func CreateRobot(userID string, tree string, eventQueue *event.EventQueue) (*Robot, error) {
	robot := &Robot{
		userID:     userID,
		roleName:   userID,
		eventQueue: eventQueue,
		logger:     sugaredLogger.With("robot", userID),
		timers:     make(map[TimerID]*timer),
		msgCBList:  list.New(),
		modules:    module.CreateModules(),
	}

	bevEntity, err := behavior.CreateEntity(tree, robot)
	if err != nil {
		return nil, err
	}

	robot.bevEntity = bevEntity

	return robot, nil
}

func (r *Robot) SetBevtree(tree string) error {
	if bevEntity, err := behavior.CreateEntity(tree, r); err != nil {
		return err
	} else {
		r.bevEntity = bevEntity
		return nil
	}
}

func (r *Robot) Tick() {
	if r.bevEntity != nil {
		r.bevEntity.Update()
	}
}

func (r *Robot) Stop() {
	r.stopSession()
	r.status.Reset()
	r.clearTimer()
	r.bevEntity.Stop()
	r.Debugf("stop")
}

func (r *Robot) UserID() string {
	return r.userID
}

func (r *Robot) RoleName() string {
	return r.roleName
}

func (r *Robot) GetStatusValue(s types.StatusID) types.StatusValue {
	return r.status.GetStatusValue(s)
}

func (r *Robot) SetStatusValue(s types.StatusID, val types.StatusValue) {
	r.status.SetStatusValue(s, val)
}

func (r *Robot) IsStatus(s types.StatusID) bool {
	return r.status.Is(s)
}

func (r *Robot) SetStatus(s types.StatusID) {
	r.status.Set(s)
}

func (r *Robot) UnsetStatus(s types.StatusID) {
	r.status.Unset(s)
}

func (r *Robot) IsLogin() bool {
	return r.session != nil
}

func (r *Robot) GetModule(moduleID int) Module {
	if moduleID < 0 || moduleID >= len(r.modules) {
		return nil
	}
	return r.modules[moduleID]
}

func (r *Robot) EventQueue() *event.EventQueue {
	return r.eventQueue
}

func (r *Robot) PostNoWait(fn interface{}, args ...interface{}) {
	r.eventQueue.PostNoWait(1, fn, args...)
}

func (r *Robot) Debug(args ...interface{}) {
	r.logger.Debug(args...)
}

func (r *Robot) Debugf(f string, args ...interface{}) {
	r.logger.Debugf(f, args...)
}

func (r *Robot) Info(args ...interface{}) {
	r.logger.Debug(args...)
}

func (r *Robot) Infof(f string, args ...interface{}) {
	r.logger.Infof(f, args...)
}

func (r *Robot) Error(args ...interface{}) {
	r.logger.Error(args...)
}

func (r *Robot) Errorf(f string, args ...interface{}) {
	r.logger.Errorf(f, args...)
}

func (r *Robot) Fatal(args ...interface{}) {
	r.logger.Fatal(args...)
}

func (r *Robot) Fatalf(f string, args ...interface{}) {
	r.logger.Fatalf(f, args...)
}

func (r *Robot) Panic(args ...interface{}) {
	r.logger.Panic(args...)
}

func (r *Robot) Panicf(f string, args ...interface{}) {
	r.logger.Panicf(f, args...)
}

func Init(logger *zap.Logger) {
	initLogger(logger)
}
