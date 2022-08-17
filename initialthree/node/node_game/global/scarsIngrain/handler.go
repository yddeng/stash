package scarsIngrain

import (
	"initialthree/zaplogger"
)

func OnUpdate(ID int32, version int64) {
	class, ok := siMgr.class[ID]
	if !ok {
		zaplogger.GetSugar().Errorf("ScarsIngrain %d is not exit", ID)
		return
	}
	class.PendFunc(func() {
		class.loadData()
	})

}
