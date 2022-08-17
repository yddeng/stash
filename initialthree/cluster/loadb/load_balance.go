package loadb

import (
	"github.com/cjbassi/gotop/src/utils"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// load balance
const (
	ModeAvg = iota
	ModeMax
	ModeMin
)

var (
	PriorityMax = 10 // 最大优先级
	PriorityMin = 1  // 最小优先级
)

// 收集信息
type LBCollector struct {
	maxToken int           // 最大承载量
	pCtor    *pidCollector // 运行状态收集

	priority        int
	tokenCount      int
	nextCollectTime time.Time
}

func New(maxToken int) *LBCollector {
	return &LBCollector{
		maxToken: maxToken,
		pCtor:    nil,
	}
}

func (this *LBCollector) cluPriority(count int) {
	if this.pCtor == nil {
		step := this.maxToken / PriorityMax
		this.priority = PriorityMax - count/step
	} else {
		used := this.pCtor.used
		// 每10%的资源为一个优先级
		this.priority = PriorityMax - int(used/10)
	}
	if this.priority < PriorityMin {
		this.priority = PriorityMin
	}
}

func (this *LBCollector) cluCount(count int) {
	if this.pCtor == nil {
		this.tokenCount = this.maxToken - count
	} else {
		used := this.pCtor.used * 100
		if count <= this.maxToken/10 || used == 0 {
			this.tokenCount = this.maxToken - count
		} else {
			this.tokenCount = int((10000 - used) / (used / float64(count)))
			if this.tokenCount > this.maxToken {
				this.tokenCount = this.maxToken - count
			}
		}
	}
	if this.tokenCount < 0 {
		this.tokenCount = 0
	}
}

// return priority, tokenCount
func (this *LBCollector) Get(count int) (int, int) {
	if this.nextCollectTime.IsZero() || time.Now().After(this.nextCollectTime) {
		if this.pCtor != nil {
			this.pCtor.collect()
		}
		this.cluPriority(count)
		this.cluCount(count)
		this.nextCollectTime = time.Now().Add(time.Second * 2)
	}

	return this.priority, this.tokenCount
}

func (this *LBCollector) SetPidCollector(pid int, mode int) {
	this.pCtor = &pidCollector{pid: pid, mode: mode}
}

type pidCollector struct {
	pid  int
	mode int
	cpu  float64
	mem  float64
	used float64
}

func (this *pidCollector) collect() {
	output, err := exec.Command("ps", "-p", strconv.Itoa(this.pid), "-o", "pcpu=12345,pmem=12345").Output()
	if err == nil {
		linesOfProcStrings := strings.Split(strings.TrimSpace(string(output)), "\n")
		if len(linesOfProcStrings) >= 2 {
			line := linesOfProcStrings[1]
			if cpu, err := strconv.ParseFloat(utils.ConvertLocalizedString(strings.TrimSpace(line[0:5])), 64); err == nil {
				this.cpu = cpu
			}
			if mem, err := strconv.ParseFloat(utils.ConvertLocalizedString(strings.TrimSpace(line[6:10])), 64); err == nil {
				this.mem = mem
			}

			this.value()
		}
	}
}

func (this *pidCollector) value() {
	switch this.mode {
	case ModeAvg:
		this.used = (this.cpu + this.mem) / 2
	case ModeMax:
		this.used = this.cpu
		if this.mem > this.used {
			this.used = this.mem
		}
	case ModeMin:
		this.used = this.cpu
		if this.mem < this.used {
			this.used = this.mem
		}
	}
}
