package Sign

import (
	codecs "initialthree/codec/cs"
	"initialthree/node/common/droppool"
	"initialthree/node/common/enumType"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/module/sign"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/Sign"
	"initialthree/node/table/excel/DataTable/SignAward"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"time"
)

type transactionSignIn struct {
	transaction.TransactionBase
	user    *user.User
	req     *codecs.Message
	errcode cs_message.ErrCode
}

func (t *transactionSignIn) Begin() {
	defer func() { t.EndTrans(&cs_message.SignInToC{}, t.errcode) }()

	t.errcode = cs_message.ErrCode_OK
	reqMsg := t.req.GetData().(*cs_message.SignInToS)
	zaplogger.GetSugar().Infof("%s SignInToS %v", t.user.GetUserLogName(), reqMsg)

	def := Sign.GetID(reqMsg.GetId())
	if def == nil {
		t.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	signModule := t.user.GetSubModule(module.Sign).(*sign.UserSign)

	// 签到次数
	signTimes := signModule.SignTimes(reqMsg.GetId())
	if signTimes >= def.SignTimes {
		t.errcode = cs_message.ErrCode_Sign_TimesFinish
		return
	}

	// 签到时间合法性
	now := time.Now()
	switch def.TypeEnum {
	case enumType.SignType_Activity, enumType.SignType_Month:
		timeLimit := def.LimitDate()
		nowUnix := now.Unix()
		// zaplogger.GetSugar().Info(nowUnix, timeLimit, time.Unix(timeLimit.StartTime, 0).String(), time.Unix(timeLimit.EndTime, 0).String())
		if nowUnix < timeLimit.StartTime || nowUnix > timeLimit.EndTime {
			t.errcode = cs_message.ErrCode_Sign_TimeFailed
			return
		}

	case enumType.SignType_Beginner:
		// 新手7日登陆
		// 无持续时间，直到奖励领取完毕后关闭功能
		//baseModule := t.user.GetSubModule(module.Base).(*base.UserBase)
		//buildTime := time.Unix(baseModule.GetBuildTime(), 0)
		//startTime := time.Date(buildTime.Year(), buildTime.Month(), buildTime.Day(), 0, 0, 0, 0, time.Local)
		//endTime := time.Date(buildTime.Year(), buildTime.Month(), buildTime.Day()+7, 0, 0, 0, 0, time.Local)
		//if now.After(endTime) || now.Before(startTime) {
		//	t.errcode = cs_message.ErrCode_Sign_TimeFailed
		//	return
		//}
	}

	lastTime := signModule.LastTimeDate(reqMsg.GetId())
	if !lastTime.IsZero() {
		dailyTime := Global.Get().GetDailyRefreshTime()
		// 偏移时间
		lastTime = lastTime.Add(-(time.Duration(dailyTime.Hour)*time.Hour + time.Duration(dailyTime.Minute)*time.Minute))
		lastDay := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), 0, 0, 0, 0, time.Local)

		nowTime := now.Add(-(time.Duration(dailyTime.Hour)*time.Hour + time.Duration(dailyTime.Minute)*time.Minute))
		nowDay := time.Date(nowTime.Year(), nowTime.Month(), nowTime.Day(), 0, 0, 0, 0, time.Local)

		if lastDay.Unix() == nowDay.Unix() {
			t.errcode = cs_message.ErrCode_Sign_Already
			return
		}
	}

	// 奖励
	signAwardDef := SignAward.GetID(def.AwardRule*1000 + (signTimes + 1))
	if signAwardDef == nil {
		t.errcode = cs_message.ErrCode_Config_NotExist
		return
	}

	dropIds := []int32{signAwardDef.DropPoolID}
	_ = t.user.ApplyDropAward(droppool.DropAward(dropIds...))

	signModule.SignIn(reqMsg.GetId())
	zaplogger.GetSugar().Infof("%s SignInToS ok", t.user.GetUserLogName())
}

func (t *transactionSignIn) GetModuleName() string {
	return "Sign"
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_SignIn, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionSignIn{
			user: user,
			req:  msg,
		}
	})
}
