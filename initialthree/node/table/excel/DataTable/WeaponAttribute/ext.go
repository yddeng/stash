package WeaponAttribute

import (
	"fmt"
	"initialthree/node/common/battleAttr"
)

func (this *WeaponAttribute) AttrID() int32 {
	id := battleAttr.GetIdByName(this.AttributeType)
	if id == 0 {
		panic(fmt.Errorf("id %s is not define", this.AttributeType))
	}
	return id
}

func (this *WeaponAttribute) GetBreakAttr(times int32) *BreakLevelAttr_ {
	idx := times - 1
	if idx < 0 || idx >= int32(len(this.BreakLevelAttr)) {
		return nil
	}
	return this.BreakLevelAttr[idx]
}
