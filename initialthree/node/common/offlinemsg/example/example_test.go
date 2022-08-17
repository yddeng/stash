package example

import (
	"fmt"
	"github.com/sniperHW/flyfish/errcode"
	"github.com/stretchr/testify/assert"
	"initialthree/node/common/offlinemsg"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	offlinemsg.ItemCap = 10

	pushCh := make(chan struct{}, 1)

	offlinemsg.PushMsg(db, "sniperHW", "test", []byte("hello"), time.Second, func(err errcode.Error, _ int64) {
		select {
		case pushCh <- struct{}{}:
		}
	})
	<-pushCh

	q := NewMsgqueue("sniperHW", "test")

	msg := q.pop()

	assert.Equal(t, "hello", string(msg))

	for i := 0; i < 10; i++ {
		q.push([]byte(fmt.Sprintf("haha:%d", i)))
	}

	for i := 0; i < 10; i++ {
		msg := q.pop()
		assert.Equal(t, fmt.Sprintf("haha:%d", i), string(msg))
	}
}
