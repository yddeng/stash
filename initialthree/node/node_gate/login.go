package node_gate

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster"
	"initialthree/codec/cs"
	"initialthree/network/smux"
	"initialthree/node/common/db"
	"initialthree/node/common/omm/whitelist"
	"initialthree/node/common/serverType"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/zaplogger"
	"net"
	"reflect"
	"sync/atomic"
	"time"
)

var serverGroups []int32

func InitGroups(groups []int32) {
	serverGroups = groups
}

func checkGroup(serverId int32) bool {
	for _, id := range serverGroups {
		if id == serverId {
			return true
		}
	}
	return false
}

var (
	csReceiver = cs.NewReceiver("cs")
	scReceiver = cs.NewReceiver("sc")
	encoder    = cs.NewEncoder("sc")
)

type loginReply struct {
	conn   net.Conn
	fire   uint32
	cmd    uint16
	seriNo uint32
}

func (this *loginReply) reply(b []byte) error {
	if atomic.CompareAndSwapUint32(&this.fire, 0, 1) {
		if _, err := this.conn.Write(b); err != nil {
			return err
		}
	}
	return nil
}

func (this *loginReply) replyErr(code message.ErrCode) error {
	if atomic.CompareAndSwapUint32(&this.fire, 0, 1) {
		b, err := cs.Encode(cs.ErrMessage(this.seriNo, this.cmd, uint16(code)), encoder)
		if err != nil {
			return err
		}
		if _, err := this.conn.Write(b); err != nil {
			return err
		}
	}
	return nil
}

func onUserLogin(conn net.Conn) (*gameSocket, *smux.MuxStream, string, error) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 5))
	data, err := cs.ReadMessage(conn)
	if err != nil {
		return nil, nil, "", err
	}
	conn.SetReadDeadline(time.Time{})

	ret, err := csReceiver.DirectUnpack(data)
	if err != nil {
		return nil, nil, "", err
	}

	msg, ok := ret.(*cs.Message)
	if !ok {
		return nil, nil, "", fmt.Errorf("invaild message  type %s . ", reflect.TypeOf(ret).String())
	}

	req, ok := msg.GetData().(*message.GameLoginToS)
	if !ok {
		return nil, nil, "", fmt.Errorf("invaild type %s . ", reflect.TypeOf(msg.GetData()).String())
	}

	zaplogger.GetSugar().Infof("%s onLogin %v", req.GetUserID(), req)

	userID := req.GetUserID()
	replier := &loginReply{conn: conn, fire: 0, seriNo: msg.GetSeriNo(), cmd: msg.GetCmd()}

	//if !checkGroup(req.GetServerID()) {
	//	replier.replyErr(message.ErrCode_ERROR)
	//	return nil, nil, userID, fmt.Errorf("onLogin serverId:%d not service in this group", req.GetServerID())
	//}

	//首先验证令牌
	/*if !CheckToken(req.GetUserID(), req.GetToken()) {
		logger.Debugln("onLogin token failed", req.GetToken())
		reply.replyErr(cs_message.ErrCode_LOGIN_TOKEN_MISMATCH)
		return
	}*/

	// 服务器状态
	if !isServerOpen() {
		r := whitelist.GetCmd(db.GetFlyfishClient("global"), userID).Exec()
		if errcode.GetCode(r.ErrCode) != errcode.Errcode_ok {
			replier.replyErr(message.ErrCode_SERVER_MAINTAINED)
			return nil, nil, userID, fmt.Errorf("server maintained")
		}
	}

	g, stream, err := modGame(userID)
	if err != nil {
		replier.replyErr(message.ErrCode_RETRY)
		return g, stream, userID, err
	}

	// 直接将本次请求转发
	if err = stream.SyncSend(data, time.Second*2); err != nil {
		replier.replyErr(message.ErrCode_ERROR)
		return g, stream, userID, err
	}

	ch := make(chan interface{}, 1)
	stream.SetCloseCallback(func(stream *smux.MuxStream, err error) {
		select {
		case ch <- err:
		default:
		}
	})

	stream.SetRecvTimeout(time.Second * 5)
	stream.Recv(func(stream *smux.MuxStream, bytes []byte) {
		b := make([]byte, len(bytes))
		copy(b, bytes)
		select {
		case ch <- b:
		default:
		}
	})

	i := <-ch

	switch i.(type) {
	case error:
		replier.replyErr(message.ErrCode_ERROR)
		return g, stream, userID, i.(error)

	case []byte:
		stream.SetRecvTimeout(0)

		bytes := i.([]byte)
		ret, err := scReceiver.DirectUnpack(bytes)
		if err != nil {
			replier.replyErr(message.ErrCode_ERROR)
			return g, stream, userID, err
		}

		if err = replier.reply(bytes); err != nil {
			return g, stream, userID, err
		}

		m := ret.(*cs.Message)
		if m.GetErrCode() == uint16(message.ErrCode_OK) {
			return g, stream, userID, nil
		} else {
			return g, stream, userID, fmt.Errorf("gameLoginToC errCode %s. ", message.ErrCode(m.GetErrCode()).String())
		}
	default:
		return g, stream, userID, fmt.Errorf("invalid type %s. ", reflect.TypeOf(i).String())
	}
}

func userLoginToDir(userID string, serverID int32) {
	userLog := &ssmessage.UserLoginToDir{
		UserID: proto.String(userID),
		AreaID: proto.Int32(serverID),
	}
	cluster.Brocast(serverType.Dir, userLog)
}
