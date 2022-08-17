package user

import (
	"github.com/golang/protobuf/proto"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/node_game/battleAtt"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/attr"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/module/equip"
	"initialthree/node/node_game/module/weapon"
	"initialthree/protocol/cs/message"
)

func (this *User) PackTeamPlayer() *message.TeamPlayer {
	attrModule := this.GetSubModule(module.Attr).(*attr.UserAttr)
	level, _ := attrModule.GetAttr(attr2.Level)

	online := false
	if this.checkStatus(status_playing) {
		online = true
	}

	return &message.TeamPlayer{
		UserID:      proto.String(this.GetUserID()),
		PlayerID:    proto.Uint64(this.GetID()),
		CharacterID: proto.Int32(100),
		PLevel:      proto.Int32(int32(level)),
		CombatPower: proto.Int32(1000),
		Name:        proto.String(this.GetName()),
		Portrait:    proto.Int32(0),
		OnLine:      proto.Bool(online),
	}
}

func (this *User) PackBattleCharacter(characterID int32) *battleAtt.BattleCharacter {
	if characterID == 0 {
		return nil
	}
	userCharacter := this.GetSubModule(module.Character).(*character.UserCharacter)
	chara := userCharacter.GetCharacter(characterID)
	if chara == nil { // 当前槽位角色被移除，清理
		return nil
	}

	c := &battleAtt.BattleCharacter{
		UserID:      this.GetUserID(),
		GameID:      this.GetID(),
		CharacterID: chara.CharacterID,
		Level:       chara.Level,
		GeneLevel:   chara.GeneLevel,
		BreakLevel:  chara.BreakLevel,
		Equips:      make([]*battleAtt.BattleEquip, 0, 6),
	}

	userWeapon := this.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	w := userWeapon.GetWeapon(chara.Weapon)
	if w != nil {
		c.Weapon = &battleAtt.BattleWeapon{
			ConfigID:   w.ConfigID,
			Level:      w.Level,
			Refine:     w.Refine,
			BreakTimes: w.BreakTimes,
		}
	}

	userEquip := this.GetSubModule(module.Equip).(*equip.UserEquip)
	for _, eId := range chara.GetEquipIDs() {
		if eId != 0 {
			equip_ := userEquip.GetEquip(eId)
			if equip_ != nil {
				e := &battleAtt.BattleEquip{
					ConfigID:     equip_.ConfigID,
					Level:        equip_.Level,
					RandomAttrID: equip_.RandomAttribId,
					Refine:       equip_.Refine,
				}
				c.Equips = append(c.Equips, e)
			}
		}
	}
	return c
}
