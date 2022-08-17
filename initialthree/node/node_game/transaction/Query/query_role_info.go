package Query

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/cluster/priority"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/db"
	"initialthree/node/common/transaction"
	"initialthree/node/node_game/module"
	"initialthree/node/node_game/user"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"strconv"
)

type transactionQueryRoleInfo struct {
	transaction.TransactionBase
	user      *user.User
	req       *codecs.Message
	reqMsg    *cs_message.QueryRoleInfoToS
	roleInfos map[uint64]*userInfo
	step      module.ModuleType
	errCode   cs_message.ErrCode
}

func (t *transactionQueryRoleInfo) Begin() {
	t.errCode = cs_message.ErrCode_OK
	t.reqMsg = t.req.GetData().(*cs_message.QueryRoleInfoToS)
	zaplogger.GetSugar().Debugf("%s Call QueryRoleInfoToS %v", t.user.GetUserLogName(), t.reqMsg)
	t.getRoleGameID()

	// 只有知道了userID才能查询是否在线
}

type Interface interface {
	Query(arg *cs_message.QueryRoleInfoArg, ret *cs_message.QueryRoleInfoResult) error
}

type userInfo struct {
	userID string
	gameID uint64
	idStr  string
	module.UserI
	*roleModuleData
	arg *cs_message.QueryRoleInfoArg
	ret *cs_message.QueryRoleInfoResult
}

func newUser(userID string, gameID uint64, arg *cs_message.QueryRoleInfoArg) *userInfo {
	r, ok := roleHotLru.Get(gameID)
	if !ok {
		r = newRoleModuleData()
		roleHotLru.Set(gameID, r)
	}
	return &userInfo{
		userID:         userID,
		gameID:         gameID,
		idStr:          fmt.Sprintf("%d", gameID),
		roleModuleData: r,
		arg:            arg,
		ret: &cs_message.QueryRoleInfoResult{
			UserID: proto.String(userID),
			GameID: proto.Uint64(gameID),
		},
	}
}

func (u *userInfo) GetID() uint64 {
	return u.gameID
}
func (u *userInfo) GetIDStr() string {
	return u.idStr
}
func (u *userInfo) GetUserID() string {
	return u.userID
}

func (t *transactionQueryRoleInfo) getRoleGameID() {
	t.roleInfos = make(map[uint64]*userInfo, len(t.reqMsg.GetQueryArgs()))

	gIDCmds := make([]*client.GetCmd, 0, len(t.reqMsg.GetQueryArgs()))
	for _, arg := range t.reqMsg.GetQueryArgs() {
		if arg.GetUserID() != "" && arg.GetGameID() == 0 {
			// 根据userID 查询 gameID
			gIDCmds = append(gIDCmds, db.GetFlyfishClient("game").
				Get("game_user", fmt.Sprintf("%s:%d", arg.GetUserID(), t.user.ServerID), "id"))
		}
	}

	if len(gIDCmds) > 0 {
		client.MGetWithEventQueue(priority.LOW, cluster.GetEventQueue(), gIDCmds...).AsyncExec(func(results []*client.GetResult) {
			ids := make(map[string]uint64, len(results))
			for _, ret := range results {
				field, ok := ret.Fields["id"]
				if ok {
					gameID := uint64(field.GetInt())
					ids[ret.Key] = gameID
				}
			}

			for _, arg := range t.reqMsg.GetQueryArgs() {
				gameID := uint64(0)
				if arg.GetUserID() == "" {
					if arg.GetGameID() != 0 {
						gameID = arg.GetGameID()
					}
				} else {
					if arg.GetGameID() != 0 {
						gameID = arg.GetGameID()
					} else {
						key := fmt.Sprintf("%s:%d", arg.GetUserID(), t.user.ServerID)
						gameID = ids[key]
					}
				}

				if gameID != 0 {
					t.roleInfos[gameID] = newUser(arg.GetUserID(), gameID, arg)
				}
			}
			t.getModuleData()
		})
	} else {
		for _, arg := range t.reqMsg.GetQueryArgs() {
			gameID := arg.GetGameID()
			if gameID != 0 {
				t.roleInfos[gameID] = newUser(arg.GetUserID(), gameID, arg)
			}
		}
		t.getModuleData()
	}

}

func (t *transactionQueryRoleInfo) isGet(arg *cs_message.QueryRoleInfoArg) bool {
	switch t.step {
	case module.Base:
		return arg.GetQueryBase()
	case module.Attr:
		return len(arg.GetAttrIDs()) > 0
	case module.Character:
		return len(arg.GetCharacterIDs()) > 0
	case module.Weapon:
		return len(arg.GetWeaponIDs()) > 0
	case module.Equip:
		return len(arg.GetEquipIDs()) > 0
	default:
		return false
	}
}

func (t *transactionQueryRoleInfo) next() module.ModuleType {
	switch t.step {
	case module.Base:
		t.step = module.Attr
	case module.Attr:
		t.step = module.Character
	case module.Character:
		t.step = module.Weapon
	case module.Weapon:
		t.step = module.Equip
	case module.Equip:
		t.step = module.End
	default:
		t.step = module.Base
	}
	return t.step
}

func (t *transactionQueryRoleInfo) readout(u *userInfo, moduleType module.ModuleType) *client.GetCmd {
	creator := module.Modules[moduleType]
	m := creator(u)
	u.roleModuleData.modules[moduleType] = m
	readOut := m.ReadOut()
	return db.GetFlyfishClient("game").Get(readOut.Table, readOut.Key, readOut.Fields...)
}

func (t *transactionQueryRoleInfo) initModule(u *userInfo, moduleType module.ModuleType, fields map[string]*client.Field) error {
	m, ok := u.modules[moduleType]
	if !ok {
		return fmt.Errorf("query role info doQuery moduleType %s not exist", moduleType.String())
	}
	return m.Init(fields)
}

func (t *transactionQueryRoleInfo) doQuery(u *userInfo, moduleType module.ModuleType) error {
	m, ok := u.modules[moduleType]
	if !ok {
		return fmt.Errorf("query role info doQuery moduleType %s not exist", moduleType.String())
	}

	mQuery, ok := m.(Interface)
	if !ok {
		return fmt.Errorf("query role info moduleType %s canot to Interface", moduleType.String())
	}

	return mQuery.Query(u.arg, u.ret)
}

func (t *transactionQueryRoleInfo) getModuleData() {
	moduleType := t.next()
	if moduleType == module.End {
		t.end()
		return
	}

	var ok bool
	cmds := make([]*client.GetCmd, 0, len(t.roleInfos))
	for _, u := range t.roleInfos {
		if ok = t.isGet(u.arg); ok {
			if _, ok = u.modules[moduleType]; !ok {
				cmds = append(cmds, t.readout(u, moduleType))
			}
		}
	}

	if len(cmds) > 0 {
		client.MGetWithEventQueue(priority.LOW, cluster.GetEventQueue(), cmds...).AsyncExec(func(results []*client.GetResult) {
			for _, ret := range results {
				if errcode.GetCode(ret.ErrCode) == errcode.Errcode_ok {
					gameID, err := strconv.ParseUint(ret.Key, 10, 64)
					if err != nil {
						zaplogger.GetSugar().Errorf("query role info parseUint %s failed:%s", ret.Key, err.Error())
						continue
					}

					u := t.roleInfos[gameID]
					if err = t.initModule(u, moduleType, ret.Fields); err != nil {
						zaplogger.GetSugar().Error(err.Error())
					}
				}
			}

			for _, u := range t.roleInfos {
				if err := t.doQuery(u, moduleType); err != nil {
					zaplogger.GetSugar().Error(err.Error())
				}
			}
			t.getModuleData()
		})
	} else {
		for _, u := range t.roleInfos {
			if err := t.doQuery(u, moduleType); err != nil {
				zaplogger.GetSugar().Error(err.Error())
			}
		}
		t.getModuleData()
	}
}

func (t *transactionQueryRoleInfo) end() {
	resp := &cs_message.QueryRoleInfoToC{
		QueryResults: make([]*cs_message.QueryRoleInfoResult, len(t.roleInfos)),
	}

	for _, u := range t.roleInfos {
		resp.QueryResults = append(resp.QueryResults, u.ret)
	}

	t.EndTrans(resp, t.errCode)
}

func (t *transactionQueryRoleInfo) GetModuleName() string {
	return "Query"
}

func init() {
	//  ServerTime 特殊的 trans，没有前置检查
	user.RegisterTransStep(cmdEnum.CS_QueryRoleInfo, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionQueryRoleInfo{
			user: user,
			req:  msg,
		}
	})
}
