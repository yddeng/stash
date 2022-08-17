package user

import (
	"errors"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/droppool"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/attr"
	"initialthree/node/node_game/module/backpack"
	"initialthree/node/node_game/module/base"
	"initialthree/node/node_game/module/mail"
	"initialthree/node/node_game/module/maindungeons"
	"initialthree/node/table/excel/DataTable/Function"
	"initialthree/node/table/excel/DataTable/Mail"
	"initialthree/protocol/cs/message"
	"time"
)

func (this *User) GenUID() uint32 {
	return this.GetSubModule(module.Base).(*base.UserBase).GenUID()
}

/* 道具 */
func (this *User) GetItemCountByTID(id int32) int32 {
	userBackpack := this.GetSubModule(module.Backpack).(*backpack.Backpack)
	return userBackpack.GetItemCountByTID(id)
}

/* 属性 */

func (this *User) GetAttr(attrID int32) int64 {
	userAttr := this.GetSubModule(module.Attr).(*attr.UserAttr)
	count, _ := userAttr.GetAttr(attrID)
	return count
}

func (this *User) GetLevel() int32 {
	userAttr := this.GetSubModule(module.Attr).(*attr.UserAttr)
	level, _ := userAttr.GetAttr(attr2.Level)
	return int32(level)
}

func (this *User) GetName() string {
	baseModule := this.GetSubModule(module.Base).(*base.UserBase)
	return baseModule.GetName()
}

/* 功能开启 */

func (this *User) FunctionUnlock(_func int32) bool {
	unlocks := Function.GetUnlock(_func)
	if unlocks == nil || len(unlocks) == 0 {
		return true
	}

	for _, v := range unlocks {
		switch v.Type {
		case enumType.FunctionUnlockType_PlayerLevel:
			if this.GetLevel() < int32(v.Level) {
				return false
			}
		case enumType.FunctionUnlockType_DungeonPass:
			userDungeon := this.GetSubModule(module.MainDungeons).(*maindungeons.MainDungeons)
			if !userDungeon.IsDungeonPass(int32(v.DungeonID)) {
				return false
			}
		case enumType.FunctionUnlockType_TimeRange:
			now := time.Now().Unix()
			length := len(v.Time)
			if length == 1 {
				if now < v.Time[0] {
					return false
				}
			} else if length == 2 {
				if !(now >= v.Time[0] && now <= v.Time[1]) {
					return false
				}
			}
		case enumType.FunctionUnlockType_Default:
			return false
		}
	}
	return true
}

func (this *User) SendMail(mails []*message.Mail) {
	mailModule := this.GetSubModule(module.Mail).(*mail.Mail)
	mailModule.AddMails(mails)
}

func (this *User) ApplyDropAward(award *droppool.Award) error {
	if award.IsZero() {
		return nil
	}

	if err := this.DoInputOutput(nil, award.ToResDesc(), true); err != message.ErrCode_OK {
		return errors.New(err.String())
	}

	return nil
}

// 返回cs的错误码
func (this *User) DoInputOutput(in, out []inoutput.ResDesc, fullSendMail bool) message.ErrCode {

	convert := func(err error) message.ErrCode {
		switch err {
		case inoutput.ErrCfgNotFound:
			return message.ErrCode_Config_NotExist
		case inoutput.ErrInputNotEnough:
			return message.ErrCode_Resource_NotEnough
		case inoutput.ErrSpaceNotEnough:
			return message.ErrCode_Backpack_SpaceNotEnough
		default: // case inoutput.ErrInvalidResCount,inoutput.ErrInvalidResType:
			return message.ErrCode_ERROR
		}
	}

	if err := inoutput.DoInputOutput(this, in, out); err != nil {
		if err == inoutput.ErrSpaceNotEnough && fullSendMail {
			if len(in) > 0 {
				if er := inoutput.DoInputOutput(this, in, nil); er != nil {
					return convert(er)
				}
			}

			// 本次全部奖励通过邮件发送
			fullAward := droppool.NewAward()
			for _, v := range out {
				fullAward.AddInfo(int32(inoutput.IOType2DropType(v.Type)), v.ID, v.Count)
			}
			if !fullAward.IsZero() {
				m := Mail.BackpackFullMail(time.Now(), fullAward.ToMessageAward())
				this.SendMail([]*message.Mail{m})
			}

		} else {
			return convert(err)
		}
	}
	return message.ErrCode_OK
}
