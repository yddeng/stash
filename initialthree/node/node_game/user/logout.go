package user

import (
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/node/common/db"
	"initialthree/node/node_game/temporary"
	"initialthree/zaplogger"
)

func StopAndWaitAllUserLogout(saveTmp bool) {
	cluster.PostTask(func() {
		//stoped = true
		for _, v := range userMap {
			v.kick(saveTmp)
		}
	})

	cluster.WaitCondition(func() bool {
		zaplogger.GetSugar().Debugf("WaitCondition userMap length %d ", len(userMap))
		return len(userMap) == 0
	})
}

// 下线
func (this *User) logout(saveTmp bool) {
	zaplogger.GetSugar().Debugf("user %s logout %v", this.GetUserID(), saveTmp)
	saveTmp = true
	if !this.checkStatus(status_logout) {

		this.setStatus(status_logout)

		// 清理计时器
		this.removeLastTimer()

		//后面的逻辑必须所有transaction都执行完毕后才能执行
		this.transMgr.Close(func() {
			this.setStatus(status_wait_remove)

			final := func() {
				// 保存数据到数据库
				this.finalSave(func() {
					// 清理数据库登陆标记
					set := db.GetFlyfishClient("game").CompareAndSetNx("user_game_login", this.userID, "gameaddr", cluster.SelfAddr().Logic.String(), "")
					set.AsyncExec(func(ret *flyfish.ValueResult) {
						zaplogger.GetSugar().Infof("%s clear db logout code:%s ", this.userID, errcode.GetErrorDesc(ret.ErrCode))
						// 清理map
						deleteUser(this)
					})
				})
			}

			if saveTmp {
				temporary.Save(this, final)
			} else {
				// 执行玩家登出逻辑
				for _, v := range this.temporaryData {
					v.UserLogout()
				}
				final()
			}

		})
	}
}
