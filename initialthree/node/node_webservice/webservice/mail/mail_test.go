package mail

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet/dhttp"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/protocol/cs/message"
	"testing"
	"time"
)

type addArg struct {
	Type   string        `json:"type"`
	GameID []uint64      `json:"gameId"`
	Mail   *message.Mail `json:"mail"`
}

func TestAddMailGlobal(t *testing.T) {
	now := time.Now()
	arg := &addArg{
		Type:   "global",
		GameID: nil,
		Mail: &message.Mail{
			Title:      proto.String("全局测试邮件1"),
			Sender:     proto.String("global"),
			CreateTime: proto.Int64(now.Unix()),
			ExpireTime: proto.Int64(now.Add(time.Hour * 24 * 2).Unix()),
			Content:    proto.String("这是一封全局邮件,不带奖励，超时时间48小时"),
			Id:         proto.Uint32(0),
		},
	}

	//req, _ := dhttp.PostJson("http://212.129.131.27:40533/mail/add", arg)
	req, _ := dhttp.PostJson("http://10.128.2.123:41801/mail/add", arg)
	fmt.Println(req.ToString())

	arg = &addArg{
		Type:   "global",
		GameID: nil,
		Mail: &message.Mail{
			Title:      proto.String("全局测试邮件2"),
			Sender:     proto.String("global"),
			CreateTime: proto.Int64(now.Unix()),
			ExpireTime: proto.Int64(0),
			Content:    proto.String("这是一封全局邮件,无超时时间。带奖励,100金币"),
			Awards: &message.Award{
				AwardInfos: []*message.AwardInfo{
					{
						Type:  proto.Int32(int32(enumType.DropType_UsualAttribute)),
						ID:    proto.Int32(attr.Gold),
						Count: proto.Int32(100),
					},
				},
			},
		},
	}

	req, _ = dhttp.PostJson("http://10.128.2.123:41801/mail/add", arg)
	fmt.Println(req.ToString())
}

func TestAddMailUser(t *testing.T) {
	now := time.Now()
	arg := &addArg{
		Type:   "user",
		GameID: []uint64{100006},
		Mail: &message.Mail{
			Title:      proto.String("mail"),
			Sender:     proto.String("user"),
			CreateTime: proto.Int64(now.Unix()),
			ExpireTime: proto.Int64(now.Add(time.Hour * 24 * 2).Unix()),
			Content:    proto.String("content"),
			Id:         proto.Uint32(0),
		},
	}

	req, _ := dhttp.PostJson("http://10.128.2.123:41801/mail/add", arg)
	fmt.Println(req.ToString())
}
