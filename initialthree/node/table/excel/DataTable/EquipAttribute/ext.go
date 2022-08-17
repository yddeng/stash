package EquipAttribute

import (
	"fmt"
	"initialthree/node/common/battleAttr"
)

func (this *EquipAttribute) AttrID() int32 {
	id := battleAttr.GetIdByName(this.AttributeType)
	if id == 0 {
		panic(fmt.Errorf("id %s is not define", this.AttributeType))
	}
	return id
}
