package attr

import (
	"initialthree/node/common/attr"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/base"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/ConstTable/Quest"
	"initialthree/node/table/excel/DataTable/PlayerLevel"
	"time"
)

func (this *UserAttr) EventAttrChange(id int32, oldVal, newVal int64) {
	//log.GetLogger().Debugln("EventAttrChange", id, oldVal, newVal)
	switch id {
	case attr.CurrentExp:
		//经验改变触发升级
		curLevel := this.getLevel()
		oldLevel := curLevel

		levelDef := PlayerLevel.GetID(curLevel)
		if levelDef == nil {
			return
		}

		curExp := newVal
		needExp := int64(levelDef.MaxExpValue)
		nextLevelDef := PlayerLevel.GetID(curLevel + 1)
		for nextLevelDef != nil && needExp != 0 && curExp >= needExp {
			curExp -= needExp
			curLevel += 1

			needExp = int64(PlayerLevel.GetID(curLevel).MaxExpValue)
			nextLevelDef = PlayerLevel.GetID(curLevel + 1)
		}

		// 到达最大等级
		if nextLevelDef == nil && curExp >= needExp && needExp != 0 {
			curExp = needExp - 1
		}

		if newVal != curExp {
			this.data.Attrs[attr.CurrentExp].Val = curExp
			this.dirty[attr.CurrentExp] = struct{}{}
			this.SetDirty(this.ModuleType().String())
		}
		if oldLevel != curLevel {
			this.SetAttr(attr.Level, int64(curLevel), true)
		}
	case attr.Level:
		addFatigue := int64(0)
		for lev := oldVal; lev < newVal; lev++ {
			def := PlayerLevel.GetID(int32(lev))
			if def == nil {
				break
			}
			addFatigue += int64(def.GiveFatigue)
		}
		if addFatigue > 0 {
			this.AddAttr(attr.CurrentFatigue, addFatigue)
		}
		// 初始化7日任务时间
		this.initNewbieTime(newVal)

	case attr.DailyActiveness:
		if newVal > oldVal {
			this.AddAttr(attr.WeeklyActiveness, newVal-oldVal)
		}
	case attr.AccumulateLogin:
		this.AddAttr(attr.WeeklyLogin, newVal-oldVal)
	}
}

func (this *UserAttr) initNewbieTime(level int64) {
	startTime, _ := this.GetAttr(attr.NewbieGiftStartTime)
	if startTime == 0 {
		def := Quest.GetID(1)
		if level >= int64(def.NewbieGiftUnlockLevel) {
			baseModule := this.userI.GetSubModule(module.Base).(*base.UserBase)
			buildTime := time.Unix(baseModule.GetBuildTime(), 0)
			now := time.Now()
			endTime := buildTime.Add(time.Duration(def.NewbieGiftUnlockLessTime) * time.Hour * 24)
			daily := Global.Get().GetDailyRefreshTime()
			endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), int(daily.Hour), int(daily.Minute), 0, 0, time.Local)
			this.SetAttr(attr.NewbieGiftStartTime, now.Unix(), false)
			this.SetAttr(attr.NewbieGiftEndTime, endTime.Unix(), false)
		}
	}
}
