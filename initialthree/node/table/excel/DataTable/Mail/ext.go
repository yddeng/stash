package Mail

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/protocol/cs/message"
	"strings"
	"time"
)

func BackpackFullMail(now time.Time, awards *message.Award, args ...string) *message.Mail {
	return GetID(1).makeMail(now, awards, args...)
}

func ScarsIngrainMail(now time.Time, awards *message.Award, args ...string) *message.Mail {
	return GetID(2).makeMail(now, awards, args...)
}

func DailyQuestMail(now time.Time, awards *message.Award, args ...string) *message.Mail {
	return GetID(3).makeMail(now, awards, args...)
}

func WeeklyQuestMail(now time.Time, awards *message.Award, args ...string) *message.Mail {
	return GetID(4).makeMail(now, awards, args...)
}

func NewbieGiftMail(now time.Time, awards *message.Award, args ...string) *message.Mail {
	return GetID(5).makeMail(now, awards, args...)
}

func (this *Mail) makeMail(now time.Time, awards *message.Award, args ...string) *message.Mail {
	if now.IsZero() {
		now = time.Now()
	}
	m := &message.Mail{
		Title:      proto.String(this.Title),
		Sender:     proto.String(this.Sender),
		CreateTime: proto.Int64(now.Unix()),
		ExpireTime: proto.Int64(0),
		Awards:     awards,
	}

	content := this.Content
	for i, arg := range args {
		content = strings.ReplaceAll(content, fmt.Sprintf("{%d}", i), arg)
	}
	m.Content = proto.String(content)

	if this.Expire != 0 {
		m.ExpireTime = proto.Int64(now.Add(time.Hour * time.Duration(int(this.Expire))).Unix())
	}
	return m
}
