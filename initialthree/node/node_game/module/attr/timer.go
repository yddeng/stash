package attr

import (
	"fmt"
	"initialthree/node/common/attr"
	"initialthree/node/common/timeDisposal"
	"initialthree/node/node_game/module"
	"initialthree/node/table/excel/ConstTable/Player"
	"initialthree/node/table/excel/DataTable/PlayerLevel"
	"time"
)

func (this *UserAttr) TimeEqual(name string, unix int64) bool {
	if timestamp, ok := this.timedata[name]; !ok {
		return false
	} else {
		return timestamp == unix
	}
}

func (this *UserAttr) tryClock(name string, nowUnix int64, fn func(nowUnix int64)) {
	if timestamp, ok := this.timedata[name]; ok {
		if nowUnix >= timestamp {
			fn(nowUnix)
		}
	} else {
		fn(nowUnix)
	}
}

func (this *UserAttr) clockTimer() {
	now := time.Now().Unix()

	this.tryClock(module.WeeklyTimeName, now, this.weeklyClock)
	this.tryClock(module.DailyTimeName, now, this.dailyClock)

	// 在线时长
	dailyOnlineName := timeId2Name(attr.DailyOnLine)
	if timestamp, ok := this.timedata[dailyOnlineName]; ok {
		if now >= timestamp {
			this.dailyOnlineClock(now)
		}
	} else {
		this.timedata[dailyOnlineName] = now + 60
		this.SetDirty(timeField)
	}

	fatigueName := timeId2Name(attr.CurrentFatigue)
	if timestamp, ok := this.timedata[fatigueName]; ok {
		if now >= timestamp {
			this.fatigueClock(now)
		}
	} else {
		this.fatigueCheck()
	}

}

func timeId2Name(id int32) string {
	return fmt.Sprintf("%s%d", timePrefix, id)
}

// 日更新
func (this *UserAttr) dailyClock(now int64) {
	// 这样设置，是防止属性变化触发事件
	this.data.Attrs[attr.FatigueBuyCount].Val = 0
	this.data.Attrs[attr.DailyActiveness].Val = 0
	this.data.Attrs[attr.DailyOnLine].Val = 0
	this.data.Attrs[attr.GoldBuyCount].Val = 0

	this.dirty[attr.FatigueBuyCount] = struct{}{}
	this.dirty[attr.DailyActiveness] = struct{}{}
	this.dirty[attr.DailyOnLine] = struct{}{}
	this.dirty[attr.GoldBuyCount] = struct{}{}

	this.SetDirty(this.ModuleType().String())

	// 累积登陆的事件触发
	this.AddAttr(attr.AccumulateLogin, 1)

	this.timedata[module.DailyTimeName] = module.CalDailyTime().Unix()
	this.SetDirty(timeField)
}

// 周更新
func (this *UserAttr) weeklyClock(now int64) {
	this.data.Attrs[attr.WeeklyActiveness].Val = 0
	this.data.Attrs[attr.WeeklyLogin].Val = 0

	this.dirty[attr.WeeklyActiveness] = struct{}{}
	this.dirty[attr.WeeklyLogin] = struct{}{}

	this.SetDirty(this.ModuleType().String())

	this.timedata[module.WeeklyTimeName] = module.CalWeeklyTime().Unix()
	this.SetDirty(timeField)
}

func (this *UserAttr) fatigueCheck() {
	name := timeId2Name(attr.CurrentFatigue)
	if _, ok := this.timedata[name]; ok {
		return
	}

	// 体力未满， 需要添加定时器
	levelMax := int64(PlayerLevel.GetID(this.getLevel()).MaxFatigueValue)
	val, _ := this.GetAttr(attr.CurrentFatigue)
	if val < levelMax {
		interval := int64(time.Duration(Player.Get().FatigueValueRecoverTimeUnit) * 60)
		timestamp := timeDisposal.NowUnix() + interval
		this.timedata[name] = timestamp
		this.dirty[attr.CurrentFatigue] = struct{}{}
		this.SetDirty(timeField)
	}
}

// 体力更新
func (this *UserAttr) fatigueClock(now int64) {
	//log.GetLogger().Debugln("clockFatigue")
	item := this.data.Attrs[attr.CurrentFatigue]
	def := Player.Get()
	levelMax := int64(PlayerLevel.GetID(this.getLevel()).MaxFatigueValue)
	interval := int64(time.Duration(Player.Get().FatigueValueRecoverTimeUnit) * 60)

	name := timeId2Name(attr.CurrentFatigue)

	nextTime := this.timedata[name]
	if item.Val >= levelMax {
		nextTime = 0
	} else {
		add := int64(0)
		for now >= nextTime {
			add += int64(def.FatigueValueRecoverValueEveryTime)
			nextTime += interval
			if item.Val+add >= levelMax {
				nextTime = 0
				break
			}
		}
		if add != 0 {
			this.AddAttr(attr.CurrentFatigue, add)
		}

	}

	if nextTime == 0 {
		delete(this.timedata, name)
	} else {
		this.dirty[attr.CurrentFatigue] = struct{}{}
		this.timedata[name] = nextTime
	}
	this.SetDirty(timeField)
}

// 日在线时长更新,60 秒一次
func (this *UserAttr) dailyOnlineClock(now int64) {
	name := timeId2Name(attr.DailyOnLine)
	// 角色状态正确，在线
	if this.userI.StatusOk() {
		_, _ = this.AddAttr(attr.DailyOnLine, 60)
	}

	timestamp := now + 60
	this.timedata[name] = timestamp
	this.SetDirty(timeField)
}
