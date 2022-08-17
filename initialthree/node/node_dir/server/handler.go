package server

import (
	"errors"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	codecs "initialthree/codec/cs"
	"initialthree/node/common/db"
	"initialthree/node/node_dir"
	"initialthree/protocol/cmdEnum"
	cs_msg "initialthree/protocol/cs/message"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/zaplogger"
	"reflect"
	"sync/atomic"
	"time"
)

func onReportStatus(_ addr.LogicAddr, msg proto.Message) {
	arg := msg.(*ssmessage.ReportStatus)
	for _, serverId := range arg.GetServerIds() {
		s := myServers.getServer(serverId)
		if s != nil {
			s.ServerStatus = 1
			s.ServerAddr = arg.GetServerAddr()
			s.timestamp = time.Now()
			s.playerNum = arg.GetPlayerNum()
		} else {
			zaplogger.GetSugar().Infof("serverMap losted serverId ==> %d", serverId)
		}
	}
}

func onUserLoginToDir(_ addr.LogicAddr, msg proto.Message) {
	arg := msg.(*ssmessage.UserLoginToDir)
	fields := map[string]interface{}{
		"lastlogin": arg.GetAreaID(),
	}
	set := db.GetFlyfishClient("dir").Set("user_dir_last_login", arg.GetUserID(), fields)
	set.AsyncExec(func(ret *flyfish.StatusResult) {
		if errcode.GetCode(ret.ErrCode) != errcode.Errcode_ok {
			zaplogger.GetSugar().Errorf("%s flyfish errcode:%s", arg.GetUserID(), errcode.GetErrorDesc(ret.ErrCode))
		}
	})
}

func onGetServerList(session *fnet.Socket, msg *codecs.Message) {

	seqNo := msg.GetSeriNo()
	arg := msg.GetData().(*cs_msg.ServerListToS)
	userID := arg.GetUserID()

	// 验证qq号
	//if !check.MatchQQ(userID){
	//  session.Send(codecs.ErrMessage(seqNo, msg.GetCmd(),uint16(cs_msg.ErrCode_ERROR)))
	//	session.Close(" Close", 1)
	//	return
	//}

	// 如果已经超过最大查询，拒绝
	if atomic.LoadInt32(&dbCurQueryCount) >= dbQueryMax {
		zaplogger.GetSugar().Debugf("dbCurQueryCount is max, %d ", dbQueryMax)
		session.Send(codecs.ErrMessage(seqNo, msg.GetCmd(), uint16(cs_msg.ErrCode_RETRY)))
		session.Close(errors.New(" Close"), 0)
		return
	}

	atomic.AddInt32(&dbCurQueryCount, 1)

	zaplogger.GetSugar().Infof("%s onGetServerList dbCurQueryCount:%d", userID, atomic.LoadInt32(&dbCurQueryCount))

	serverList := myServers.ServerList

	serverListToC := &cs_msg.ServerListToC{
		ServerList: make([]*cs_msg.Server, 0, len(serverList)),
	}
	for _, v := range serverList {
		server := &cs_msg.Server{}

		if v.ServerStatus == 1 {
			server.ServerType = cs_msg.ServerType_OPERATION.Enum()
		} else {
			server.ServerType = cs_msg.ServerType_SHUTOFF.Enum()
		}
		server.UserType = cs_msg.UserType_NONE.Enum()
		server.ServerId = proto.Int32(v.ServerId)
		server.ServerName = proto.String(v.ServerName)
		server.ServerAddr = proto.String(v.ServerAddr)
		server.PlayerNum = proto.Int32(v.playerNum)

		serverListToC.ServerList = append(serverListToC.ServerList, server)
	}

	getLastLogin(userID, func(areaID int32) {
		if areaID != 0 {
			for _, s := range serverListToC.ServerList {
				if s.GetServerId() == areaID {
					s.UserType = cs_msg.UserType_LAST_LOGIN.Enum()
				}
			}
		}
		atomic.AddInt32(&dbCurQueryCount, -1)
		zaplogger.GetSugar().Info(userID, serverListToC, areaID)
		session.Send(codecs.NewMessage(seqNo, serverListToC))
		session.Close(errors.New(" Close"), time.Second)
	})

}

func getLastLogin(userID string, cb func(areaID int32)) {
	get := db.GetFlyfishClient("dir").Get("user_dir_last_login", userID, "lastlogin")
	get.AsyncExec(func(ret *flyfish.GetResult) {
		if errcode.GetCode(ret.ErrCode) == errcode.Errcode_ok {
			id := ret.Fields["lastlogin"].GetInt()
			cb(int32(id))
		} else {
			zaplogger.GetSugar().Errorf("%s flyfish errcode:%s", userID, errcode.GetErrorDesc(ret.ErrCode))
			cb(0)
		}
	})
}

func onGetAnnouncement(session *fnet.Socket, msg *codecs.Message) {
	seqNo := msg.GetSeriNo()
	//arg := msg.GetData().(*cs_msg.AnnouncementToS)
	session.Send(codecs.NewMessage(seqNo, &cs_msg.AnnouncementToC{
		Announcement: proto.String("fsd"),
	}))
	/*
		oldVersion := arg.GetVersion()

		zaplogger.GetSugar().Infof("get announcement version:%d ", oldVersion)

		const (
			tableName = "web_data"
			tableKey  = "announcement"
			slotCount = 100
		)
		db.GetFlyfishClient("global").GetAllWithVersion(tableName, tableKey, announceVersion).AsyncExec(func(result *flyfish.SliceResult) {
			if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok ||
				errcode.GetCode(result.ErrCode) == errcode.Errcode_record_notexist ||
				errcode.GetCode(result.ErrCode) == errcode.Errcode_record_unchange {

				if errcode.GetCode(result.ErrCode) == errcode.Errcode_ok {
					nowUnix := time.Now().Unix()
					announcementResult = make([]*cs_msg.Announcement, 0, slotCount)
					for i := 0; i < slotCount; i++ {
						fieldName := fmt.Sprintf("slot%d", i)
						field, ok := result.Fields[fieldName]
						var ann cs_msg.Announcement
						if ok && len(field.GetBlob()) > 0 {
							if err := proto.Unmarshal(field.GetBlob(), &ann); err == nil {
								if ann.GetExpireTime() == 0 || nowUnix < ann.GetExpireTime() {
									announcementResult = append(announcementResult, &ann)
								}
							}

						}
					}

					announceVersion = result.Version
				}

				resp := &cs_msg.AnnouncementToC{
					Version: proto.Int64(announceVersion),
					Announcements: []*cs_msg.Announcement{
						{
							Id:         proto.Int32(1),
							Type:       proto.String("AnnouncementType_System"),
							Title:      proto.String("fsfs"),
							SmallTitle: proto.String("small"),
							StartTime:  proto.Int64(9),
							ExpireTime: proto.Int64(0),
							Remind:     proto.Bool(false),
							Content: []*cs_msg.AnnouncementContent{
								&cs_msg.AnnouncementContent{
									Type:      proto.String("1"),
									ImageSkip: proto.Int32(0),
									Text:      proto.String("test"),
									Image:     proto.String("Announce_banner_1"),
								},
							},
						},
					},
				}
				//if oldVersion != announceVersion {
				//	resp.Announcements = announcementResult
				//}

				session.Send(codecs.NewMessage(seqNo, resp))
				zaplogger.GetSugar().Infof("load announcement %d %v ok", announceVersion, resp)

			} else {
				zaplogger.GetSugar().Errorf("get announcement failed, %d", result.ErrCode.Code)
				session.Send(codecs.ErrMessage(seqNo, cmdEnum.CS_Announcement, uint16(cs_msg.ErrCode_RETRY)))
			}
		})

	*/

}

var (
	announceVersion int64
	//announcementResult []*cs_msg.Announcement
)

func init() {
	cluster.Register(cmdEnum.SS_ReportStatus, onReportStatus)
	cluster.Register(cmdEnum.SS_UserLoginToDir, onUserLoginToDir)
	node_dir.RegisterHandler(cmdEnum.CS_ServerList, onGetServerList)
	node_dir.RegisterHandler(cmdEnum.CS_Announcement, onGetAnnouncement)
	node_dir.RegisterHandler(cmdEnum.CS_Heartbeat, func(session *fnet.Socket, message *codecs.Message) {
		zaplogger.GetSugar().Debug("heartbeat", reflect.ValueOf(message.GetData()).String())
	})
}
