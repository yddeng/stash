package cluster

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/sniperHW/flyfish/pkg/crypto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/cluster/addr"
	"io"
	"net"
	"time"
)

var key []byte = []byte("feiyu_tech_2022")

type loginReq struct {
	Mux       bool   `json:"Mux,omitempty"`
	LogicAddr uint32 `json:"LogicAddr,omitempty"`
	NetAddr   string `json:"NetAddr,omitempty"`
}

func (this *Cluster) login(end *endPoint, conn net.Conn, ismux bool) error {

	j, err := json.Marshal(&loginReq{
		Mux:       ismux,
		LogicAddr: uint32(this.serverState.selfAddr.Logic),
		NetAddr:   conn.LocalAddr().String(),
	})

	if nil != err {
		return err
	}

	if j, err = crypto.AESCBCEncrypt(key, j); nil != err {
		return err
	}

	b := make([]byte, 4+len(j))
	binary.BigEndian.PutUint32(b, uint32(len(j)))
	copy(b[4:], j)

	conn.SetWriteDeadline(time.Now().Add(time.Second))
	_, err = conn.Write(b)
	conn.SetWriteDeadline(time.Time{})

	if nil != err {
		return err
	} else {
		buffer := make([]byte, 4)
		conn.SetReadDeadline(time.Now().Add(time.Second))
		_, err = io.ReadFull(conn, buffer)
		conn.SetReadDeadline(time.Time{})
		if nil != err {
			return err
		}

		if binary.BigEndian.Uint32(buffer) != 0 {
			return ERR_AUTH
		}
	}
	return nil
}

func (this *Cluster) auth(conn net.Conn) error {

	bLen := make([]byte, 4)
	conn.SetReadDeadline(time.Now().Add(time.Second))
	defer conn.SetReadDeadline(time.Time{})

	_, err := io.ReadFull(conn, bLen)
	if nil != err {
		return err
	}

	datasize := int(binary.BigEndian.Uint32(bLen))

	if datasize > 128 {
		return errors.New("packet too large")
	}

	b := make([]byte, datasize)

	_, err = io.ReadFull(conn, b)
	if nil != err {
		return err
	}

	if b, err = crypto.AESCBCDecrypter(key, b); nil != err {
		return err
	}

	var req loginReq

	if err = json.Unmarshal(b, &req); nil != err {
		return err
	}

	//lIdx1 := strings.LastIndex(req.NetAddr, ":")
	//fromAddr := conn.RemoteAddr().String()
	//lIdx2 := strings.LastIndex(fromAddr, ":")

	//if req.NetAddr[lIdx1+1:] != fromAddr[lIdx2+1:] {
	//	return errors.New("illegal connection: " + req.NetAddr + "<=>" + conn.RemoteAddr().String())
	//}

	//if req.NetAddr != conn.RemoteAddr().String() {
	//	return errors.New("illegal connection: " + req.NetAddr + "<=>" + conn.RemoteAddr().String())
	//}

	ismux := req.Mux
	logicAddr := req.LogicAddr

	end := this.serviceMgr.getEndPoint(addr.LogicAddr(logicAddr))
	if nil == end {
		return ERR_INVAILD_ENDPOINT
	}

	if !ismux {
		end.Lock()
		if end.dialing {
			//当前节点同时正在向对端dialing,逻辑地址小的一方放弃接受连接
			if this.serverState.selfAddr.Logic < end.addr.Logic {
				end.Unlock()
				logger.Sugar().Errorf("(self:%v) (other:%v) both side connectting", this.serverState.selfAddr.Logic, end.addr.Logic)
				return errors.New("both side connectting")
			}
		}

		if nil != end.session {
			logger.Sugar().Infof("(self:%v) auth duplicate %v\n", this.serverState.selfAddr.Logic, end.addr.Logic)
			err = ERR_DUP_CONN
		} else {
			this.onEstablishServer(end, fnet.NewSocket(conn, fnet.OutputBufLimit{
				OutPutLimitSoft:        512 * 1024,
				OutPutLimitSoftSeconds: 10,
				OutPutLimitHard:        8 * 1024 * 1024,
			}))
		}
		end.Unlock()
	}

	resp := []byte{0, 0, 0, 0}
	binary.BigEndian.PutUint32(resp, 0)

	conn.SetWriteDeadline(time.Now().Add(time.Second))
	defer conn.SetWriteDeadline(time.Time{})

	_, sendErr := conn.Write(resp)

	if nil == sendErr && ismux {
		if this.onNewMuxConn == nil {
			logger.Sugar().Infof("auth muxConn %s , but onNewMuxConn is nil", this.serverState.selfAddr.Logic.String())
			return errors.New("onNewMuxConn is nil")
		} else {
			this.onNewMuxConn(addr.LogicAddr(logicAddr), conn)
		}
	}
	return sendErr
}
