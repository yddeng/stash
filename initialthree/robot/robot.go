package robot

import (
	"fmt"
	"initialthree/node/table/excel"
	"initialthree/node/table/quest"
	"initialthree/pkg/event"
	"initialthree/pkg/timer"
	"initialthree/robot/behavior"

	"initialthree/robot/config"
	"initialthree/robot/eventqueue"
	"initialthree/robot/robot"
	"initialthree/robot/statistics"
	"initialthree/zaplogger"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
	"unsafe"

	"go.uber.org/zap"
)

func Start(cfgPath string) {
	cfg, err := config.Init(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config: %s", err)
		os.Exit(1)
	}

	zaplogger.InitLogger(zaplogger.NewZapLogger("robot.log", cfg.Log.Dir, cfg.Log.Level, cfg.Log.MaxSize, cfg.Log.MaxAge, int(cfg.Log.MaxBackups), cfg.Log.EnableStdOut))

	zaplogger.GetSugar().Info("config loaded.")

	if cfg.PProf.Enable {
		zaplogger.GetSugar().Infof("start pprof at *:%d", cfg.PProf.Port)
		startPProf(cfg.PProf.Port)
	}

	zaplogger.GetSugar().Info("load excel.")
	excel.Load(cfg.Resource.ExcelPath)
	zaplogger.GetLogger().Info("load quest.")
	quest.Load(cfg.Resource.QuestPath)

	zaplogger.GetSugar().Info("init behaviors.")
	if err := behavior.Init(cfg.Resource.BehaviorConfigPath); err != nil {
		zaplogger.GetSugar().Error("init behaviors:", err)
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())

	zaplogger.GetSugar().Info("initialize eventQueues.")
	eventqueue.Init()

	robot.Init(zaplogger.GetLogger())
	createRobots(cfg.Robot.Count, cfg.Robot.UserIDPrefix, cfg.Robot.InitialNo, cfg.Resource.BehaviorTree, time.Duration(cfg.TickCycle)*time.Millisecond)

	if cfg.DisconnectTest.Enable {
		startRandomDisconnect(time.Duration(cfg.DisconnectTest.Interval)*time.Second, int(cfg.DisconnectTest.Count))
	}

	zaplogger.GetSugar().Infof("expose statistics metrics at http://*:%d/metrics.", cfg.Statistics.MetricsPort)
	statistics.ExposeMetrics(cfg.Statistics.MetricsPort, zap.NewStdLog(zaplogger.GetLogger()))
	startOutputStatistics(cfg.GetStatisticsOutputInterval(), cfg.Statistics.OutputFile)

	zaplogger.GetSugar().Info("start eventQueues.")
	eventqueue.Start()

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	_ = <-c

	/*
		停机
	*/

	if cfg.DisconnectTest.Enable {
		stopRandomDisconnect()
	}

	destoryRobots()

	zaplogger.GetSugar().Info("stop eventQueues.")
	eventqueue.Stop()

	zaplogger.GetSugar().Info("conceal statistics metrics.")
	statistics.ConcealMetrics()
	stopOutputStatistics()

	zaplogger.GetSugar().Info("ROBOT stop.")
}

type robotSugar struct {
	*robot.Robot
}

func (r *robotSugar) stop(wait *sync.WaitGroup) {
	wait.Add(1)
	r.Robot.PostNoWait(func() {
		r.Robot.Stop()
		wait.Done()
	})
}

func (r *robotSugar) disconnect() {
	r.Robot.PostNoWait(r.Robot.ActiveDisconnect)
}

type robotGroup struct {
	state     int32
	eventQue  *event.EventQueue
	tickTimer *timer.Timer
	robots    []*robotSugar
}

func (rg *robotGroup) addRobot(robot *robotSugar) {
	rg.robots = append(rg.robots, robot)
}

func (rg *robotGroup) start(tickCycle time.Duration) {
	if atomic.CompareAndSwapInt32(&rg.state, 0, 1) {
		rg.tickTimer = timer.Repeat(tickCycle, func(t *timer.Timer, i interface{}) {
			if atomic.LoadInt32(&rg.state) == 1 {
				rg.eventQue.PostNoWait(0, rg.tick)
			}
		}, nil)
	}
}

func (rg *robotGroup) stop() {
	if atomic.CompareAndSwapInt32(&rg.state, 1, 2) {
		rg.tickTimer.Cancel()
		rg.tickTimer = nil
	}
}

func (rg *robotGroup) tick() {
	for _, robot := range rg.robots {
		robot.Tick()
	}
}

var mtxRobots sync.Mutex
var robotMap map[string]*robotSugar
var robotIDList []string
var robotGroups map[*event.EventQueue]*robotGroup

func createRobots(count uint, idPrefix string, initailNo int, bevtree string, tickCycle time.Duration) {
	mtxRobots.Lock()
	defer mtxRobots.Unlock()

	if robotMap != nil {
		return
	}

	zaplogger.GetLogger().Info("create robots.")

	robotMap = make(map[string]*robotSugar, count)
	robotIDList = make([]string, count)

	eventQueueCount := eventqueue.EventQueueCount()
	robotGroups = make(map[*event.EventQueue]*robotGroup, eventQueueCount)
	for i := uint(0); i < eventQueueCount; i++ {
		eventQueue := eventqueue.ModEventQueue(i)
		robotGroup := &robotGroup{
			eventQue: eventQueue,
		}
		robotGroup.start(tickCycle)
		robotGroups[eventQueue] = robotGroup
	}

	for i := uint(0); i < count; i++ {
		rid := fmt.Sprintf("%s%d", idPrefix, initailNo+int(i))
		robot, err := robot.CreateRobot(rid, bevtree, eventqueue.ModEventQueue(i))
		if err != nil {
			zaplogger.GetSugar().Error("create No.%d robot \"%s\":", err)
			os.Exit(1)
		}

		robotSugar := &robotSugar{Robot: robot}
		robotMap[rid] = robotSugar
		robotIDList[i] = rid
		robotGroups[robot.EventQueue()].addRobot(robotSugar)
	}
}

func destoryRobots() {
	mtxRobots.Lock()
	defer mtxRobots.Unlock()

	if robotMap == nil {
		return
	}

	zaplogger.GetSugar().Info("destory robots.")

	for _, robotGroup := range robotGroups {
		robotGroup.stop()
	}

	var wait sync.WaitGroup
	for _, robot := range robotMap {
		robot.stop(&wait)
	}
	wait.Wait()

	robotMap = nil
	robotIDList = nil
	robotGroups = nil
}

var timerRandomDisconnect *timer.Timer

func startRandomDisconnect(interval time.Duration, disCount int) {
	if interval <= 0 || disCount <= 0 {
		return
	}

	zaplogger.GetSugar().Info("start random disconnect.")

	timerRandomDisconnect = timer.Repeat(interval, func(t *timer.Timer, i interface{}) {
		timer := (*timer.Timer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&timerRandomDisconnect))))
		if t != timer {
			return
		}

		zaplogger.GetSugar().Info("timed random disconnection")

		mtxRobots.Lock()
		defer mtxRobots.Unlock()

		robotCount := len(robotIDList)
		if robotCount == 0 {
			return
		}

		if disCount >= robotCount {
			for _, v := range robotMap {
				v.disconnect()
			}
		} else {
			robotIDs := make([]string, robotCount)
			copy(robotIDs, robotIDList)
			n := robotCount
			for m := 0; m < disCount; m++ {
				k := rand.Intn(n)
				robot := robotMap[robotIDs[k]]
				robot.disconnect()
				robotIDs[k] = robotIDs[n-1]
				n--
			}
		}

	}, nil)
}

func stopRandomDisconnect() {
	timer := (*timer.Timer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&timerRandomDisconnect))))
	if timer != nil {
		zaplogger.GetSugar().Info("stop random disconnect.")
		timer.Cancel()
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&timerRandomDisconnect)), nil)
	}
}

var timerStatistics *timer.Timer
var fileStatistics *os.File

func startOutputStatistics(interval time.Duration, file string) {
	var err error
	fileStatistics, err = os.Create(file)
	if err != nil {
		zaplogger.GetSugar().Fatal("create statistics file:", err)
		os.Exit(1)
	}

	zaplogger.GetSugar().Info("start output statistics.")

	timerStatistics = timer.Repeat(interval, func(t *timer.Timer, i interface{}) {
		timer := (*timer.Timer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&timerStatistics))))

		if t != timer {
			return
		}

		file := (*os.File)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&fileStatistics))))
		if file == nil {
			return
		}

		zaplogger.GetSugar().Info("output statistics.")

		file.Seek(0, 0)
		file.Truncate(0)

		if err := statistics.Statistics().Output(file); err != nil {
			zaplogger.GetSugar().Error("output statistics:", err)
		} else {
			file.Sync()
		}
	}, nil)
}

func stopOutputStatistics() {
	timer := (*timer.Timer)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&timerStatistics))))
	if timer != nil {
		zaplogger.GetSugar().Info("stop output statistics.")
		timer.Cancel()
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&timerStatistics)), nil)
	}

	file := (*os.File)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&fileStatistics))))
	if file != nil {
		zaplogger.GetSugar().Info("output final statistics.")
		file.Seek(0, 0)
		file.Truncate(0)
		if err := statistics.Statistics().Output(file); err != nil {
			zaplogger.GetSugar().Error("output final statistics:", err)
		} else {
			file.Sync()
		}

		file.Close()
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&fileStatistics)), nil)
	}
}

func startPProf(port int) {
	go func() {
		if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
			zaplogger.GetSugar().Fatal("start pprof:", err)
		}
	}()
}
