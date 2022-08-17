package Base

import (
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	codecs "initialthree/codec/cs"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/db"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/transaction"
	"initialthree/node/common/wordsFilter"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/backpack"
	"initialthree/node/node_game/module/base"
	"initialthree/node/table/excel/ConstTable/Player"
	"initialthree/zaplogger"
	"strconv"
	"strings"

	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionBaseSetName struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.BaseSetNameToC
	errcode cs_message.ErrCode
}

func (this *transactionBaseSetName) GetModuleName() string {
	return "Base"
}

func (this *transactionBaseSetName) Begin() {
	this.errcode = cs_message.ErrCode_OK
	this.resp = &cs_message.BaseSetNameToC{}
	msg := this.req.GetData().(*cs_message.BaseSetNameToS)

	zaplogger.GetSugar().Infof("%s Call BaseSetNameToS %v", this.user.GetUserLogName(), msg)

	useRes := make([]inoutput.ResDesc, 0, 1)

	userBackpack := this.user.GetSubModule(module.Backpack).(*backpack.Backpack)
	if freeTimes := this.user.GetAttr(attr2.FreeRenameTimes); freeTimes > 0 { // 免费次数
		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_UsualAttribute, ID: int32(attr2.FreeRenameTimes), Count: 1})
	} else if itemCount := userBackpack.GetItemCountByTID(Player.Get().RenameCostItemID); itemCount > 0 { // 道具
		useRes = append(useRes, inoutput.ResDesc{Type: enumType.IOType_Item, ID: Player.Get().RenameCostItemID, Count: 1})
	} else {
		zaplogger.GetSugar().Infof("%s BaseSetNameToS name is failed : resource not enough", this.user.GetUserID())
		this.errcode = cs_message.ErrCode_Resource_NotEnough
		this.EndTrans(this.resp, this.errcode)
		return
	}

	//验证名字合法性
	if !this.checkName(msg.GetName()) {
		zaplogger.GetSugar().Infof("%s BaseSetNameToS name is failed : %s", this.user.GetUserID(), msg.GetName())
		this.errcode = cs_message.ErrCode_Create_Role_Name_Error
		this.EndTrans(this.resp, this.errcode)
		return
	}

	//验证名字唯一性
	fields := map[string]interface{}{}
	fields["owner"] = this.user.GetUserID() + ":" + strconv.FormatUint(this.user.GetID(), 10)

	set := db.GetFlyfishClient("game").SetNx("role_name", msg.GetName(), fields)
	set.AsyncExec(func(ret *flyfish.ValueResult) {
		if !this.IsTimeout() {
			if errcode.GetCode(ret.ErrCode) == errcode.Errcode_ok {
				// 资源消耗
				if this.errcode = this.user.DoInputOutput(useRes, nil, false); this.errcode != cs_message.ErrCode_OK {
					// 将保存的名字移除
					this.delName(msg.GetName())
				} else {
					baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
					oldName := baseModule.GetName()
					this.delName(oldName)
					baseModule.SetName(msg.GetName())

					zaplogger.GetSugar().Infof("%s Call BaseSetNameToS ok", this.user.GetUserLogName())
				}

			} else {
				zaplogger.GetSugar().Infof("%s name:%s flyfish ret code:%s", this.user.GetUserID(), msg.GetName(), errcode.GetErrorDesc(ret.ErrCode))
				if errcode.GetCode(ret.ErrCode) == errcode.Errcode_record_exist {
					this.errcode = cs_message.ErrCode_Create_Role_Name_Repeat
				} else {
					this.errcode = cs_message.ErrCode_RETRY
				}

			}
			this.EndTrans(this.resp, this.errcode)
		} else {
			if errcode.GetCode(ret.ErrCode) == errcode.Errcode_ok {
				// 将保存的名字移除
				this.delName(msg.GetName())
			}
		}

	})

}

func (this *transactionBaseSetName) delName(name string) {
	del := db.GetFlyfishClient("game").Del("role_name", name)
	del.AsyncExec(func(ret *flyfish.StatusResult) {
		zaplogger.GetSugar().Infof("BaseSetNameToS delete name:%s ret code:%s", name, errcode.GetErrorDesc(ret.ErrCode))
	})
}

func (this *transactionBaseSetName) checkName(name string) bool {
	str := []rune(name)
	if strings.TrimSpace(name) == "" || len(str) > 16 {
		return false
	}

	return !wordsFilter.Check(name)
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_BaseSetName, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionBaseSetName{
			user: user,
			req:  msg,
		}
	})
}
