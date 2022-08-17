package example

import (
	"container/list"
	"fmt"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/node/common/offlinemsg"
	"initialthree/node/common/offlinemsg/testdb"
	"sync"
	"time"
)

var db offlinemsg.DB = testdb.NewTestDB()

type msgqueue struct {
	mu       sync.Mutex
	cond     *sync.Cond
	version  int64
	msgs     *list.List
	userID   string
	topic    string
	notifyCh chan int64
}

func NewMsgqueue(userID, topic string) *msgqueue {
	q := &msgqueue{
		userID:   userID,
		topic:    topic,
		msgs:     list.New(),
		notifyCh: make(chan int64, 1),
	}
	q.cond = sync.NewCond(&q.mu)
	//先同步一次
	q.pull(true, 0)

	go func() {
		//监听远程队列，如果远程队列添加了消息，同步到本地
		for v := range q.notifyCh {
			fmt.Println("notify", v)
			q.pull(false, v)
		}
	}()

	return q
}

func (q *msgqueue) pushCb(err errcode.Error, itemVersion int64) {
	fmt.Println("push", err, itemVersion)
	if nil == err {
		select {
		case q.notifyCh <- itemVersion:
		}
	}
}

//跟远程队列同步消息
func (q *msgqueue) pull(force bool, itemVersion int64) {
	/*
	 *  对于连续发布的item 1,2。如果接收者先收到item 1的通告，接收者同步时将把item 1，2都同步到本地。
	 *  因此接到item 2的通告时，就不需要再次同步了。
	 */
	if force || q.version < itemVersion {
		offlinemsg.PullMsg(db, q.userID, q.topic, q.version, time.Second, func(err errcode.Error, version int64, items []*offlinemsg.Item) {
			if nil == err {
				q.version = version
				if len(items) > 0 {
					q.mu.Lock()
					fmt.Printf("got %d items\n", len(items))
					for _, v := range items {
						q.msgs.PushBack(v)
					}
					q.mu.Unlock()
					q.cond.Broadcast()
				}
			}
		})
	}
}

//从本地队列弹出一条消息
func (q *msgqueue) pop() []byte {
	q.mu.Lock()
	for q.msgs.Len() == 0 {
		q.cond.Wait()
	}

	m := q.msgs.Remove(q.msgs.Front()).(*offlinemsg.Item).Content
	q.mu.Unlock()
	return m
}

//向远程队列投递一条消息
func (q *msgqueue) push(m []byte) {
	offlinemsg.PushMsg(db, q.userID, q.topic, m, time.Second, q.pushCb)
}
