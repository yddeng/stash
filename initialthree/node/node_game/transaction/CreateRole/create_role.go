package CreateRole

import (
	"fmt"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/db"
	"initialthree/node/common/idGenerator"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/base"
	"initialthree/zaplogger"
	"strconv"

	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionCreateRole struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	resp    *cs_message.CreateRoleToC
	msg     *cs_message.CreateRoleToS
	errcode cs_message.ErrCode
}

func (this *transactionCreateRole) GetModuleName() string {
	return "CreateRole"
}

func (this *transactionCreateRole) Begin() {
	this.resp = &cs_message.CreateRoleToC{}
	this.msg = this.req.GetData().(*cs_message.CreateRoleToS)
	zaplogger.GetSugar().Infof("%s CreateRoleToS %v", this.user.GetUserLogName(), this.msg)
	this.errcode = cs_message.ErrCode_OK

	if !(this.msg.GetSex() == 0 || this.msg.GetSex() == 1) {
		zaplogger.GetSugar().Infof("%s %s %d", this.user.GetUserID(), "CreateRole fail : sex err", this.msg.GetSex())
		this.errcode = cs_message.ErrCode_Request_Argument_Err
		this.EndTrans(this.resp, this.errcode)
		return
	}

	//验证名字合法性
	if !this.checkName(this.msg.GetName()) {
		zaplogger.GetSugar().Infof("%s %s %s", this.user.GetUserID(), "CreateRoleName name is fail :", this.msg.GetName())
		this.errcode = cs_message.ErrCode_Create_Role_Name_Error
		this.EndTrans(this.resp, this.errcode)
		return
	}

	if this.user.GetID() == 0 {
		this.initGameID()
	} else if this.user.GetName() == "" {
		this.setUserData()
	} else {
		// 玩家已经创建角色成功，直接返回
		zaplogger.GetSugar().Infof("user:%s already create role.", this.user.GetUserLogName())
		this.EndTrans(nil, this.errcode)
		return
	}
}

func (this *transactionCreateRole) initGameID() {
	this.AsynWrap(idGenerator.GetIDGen("game_id").GenID)(func(gameID int64, err error) {
		if err != nil {
			zaplogger.GetSugar().Infof("%s CreateRole error:%s.", this.user.GetUserLogName(), err)
			this.errcode = cs_message.ErrCode_ERROR
			this.EndTrans(this.resp, this.errcode)
			return
		}

		key := fmt.Sprintf("%s:%d", this.user.GetUserID(), this.user.ServerID)
		fields := map[string]interface{}{"id": gameID}
		set := db.GetFlyfishClient("game").Set("game_user", key, fields)
		set.AsyncExec(func(ret *flyfish.StatusResult) {
			if errcode.GetCode(ret.ErrCode) == errcode.Errcode_ok {
				this.user.SetID(uint64(gameID))
				this.setUserData()
			} else {
				zaplogger.GetSugar().Infof("%s CreateRole error:%s.", this.user.GetUserLogName(), errcode.GetErrorDesc(ret.ErrCode))
				this.errcode = cs_message.ErrCode_ERROR
				this.EndTrans(this.resp, this.errcode)
			}
		})
	})
}

func (this *transactionCreateRole) setUserData() {
	//验证名字唯一性
	fields := map[string]interface{}{}
	fields["owner"] = this.user.GetUserID() + ":" + strconv.FormatUint(this.user.GetID(), 10)

	set := db.GetFlyfishClient("game").SetNx("role_name", this.msg.GetName(), fields)
	set.AsyncExec(func(ret *flyfish.ValueResult) {
		if !this.IsTimeout() {
			if errcode.GetCode(ret.ErrCode) == errcode.Errcode_ok {
				baseModule := this.user.GetSubModule(module.Base).(*base.UserBase)
				baseModule.SetSex(this.msg.GetSex())
				baseModule.SetName(this.msg.GetName())
				this.EndTrans(this.resp, this.errcode)
				this.user.DoLoadPipeline()

			} else {
				zaplogger.GetSugar().Infof("%s name:%s repeated flyfish ret code:%s", this.user.GetUserID(), this.msg.GetName(), errcode.GetErrorDesc(ret.ErrCode))
				if errcode.GetCode(ret.ErrCode) == errcode.Errcode_record_exist {
					this.errcode = cs_message.ErrCode_Create_Role_Name_Repeat
				} else {
					this.errcode = cs_message.ErrCode_RETRY
				}
				this.EndTrans(this.resp, this.errcode)
			}
		} else {
			if errcode.GetCode(ret.ErrCode) == errcode.Errcode_ok {
				this.delName()
			}
		}
	})
}

func (this *transactionCreateRole) delName() {
	// 将保存的名字移除
	del := db.GetFlyfishClient("game").Del("role_name", this.msg.GetName())
	del.AsyncExec(func(ret *flyfish.StatusResult) {
		zaplogger.GetSugar().Infof("delete name:%s ret code:%s", this.msg.GetName(), errcode.GetErrorDesc(ret.ErrCode))
	})
}

func init() {
	//  CreateRole 特殊的
	user.RegisterTransStep(cmdEnum.CS_CreateRole, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionCreateRole{
			user: user,
			req:  msg,
		}
	}, user.StepIsMessageDisable, user.StepIsModuleDisable)
}
