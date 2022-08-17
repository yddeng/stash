package verify

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"
)

const fixedStr string = "flyfish"

func Auth(svrConn net.Conn) error {
	randNum := rand.Uint64()
	token := make([]byte, 8)
	binary.BigEndian.PutUint64(token, randNum)
	token = append(token, fixedStr...)

	sum := md5.Sum(token)
	buff := make([]byte, 24)
	binary.BigEndian.PutUint64(buff, randNum)
	for i := 0; i < len(sum); i++ {
		buff[i+8] = sum[i]
	}

	svrConn.SetWriteDeadline(time.Now().Add(time.Second * 5))
	_, err := svrConn.Write(buff)
	svrConn.SetWriteDeadline(time.Time{})

	return err
}

func Verify(userConn net.Conn) error {
	buf := make([]byte, 24)
	userConn.SetReadDeadline(time.Now().Add(time.Second * 5))
	_, err := io.ReadFull(userConn, buf)
	if nil != err {
		return fmt.Errorf("the user conn [%s] was rejected, err:%v", userConn.RemoteAddr().String(), err)
	}
	userConn.SetReadDeadline(time.Time{})
	randNum := binary.BigEndian.Uint64(buf)
	token := make([]byte, 8)
	binary.BigEndian.PutUint64(token, randNum)
	token = append(token, fixedStr...)
	sum := md5.Sum(token)
	for i := 0; i < len(sum); i++ {
		if sum[i] != buf[8+i] {
			return fmt.Errorf("the user conn [%s] was rejected, Token missmatch", userConn.RemoteAddr().String())
		}
	}
	return nil
}
