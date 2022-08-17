package offlinemsg

import (
	"fmt"
	"github.com/sniperHW/flyfish/errcode"
	"github.com/stretchr/testify/assert"
	"initialthree/node/common/offlinemsg/testdb"
	"testing"
	"time"
)

func TestOfflineMsg(t *testing.T) {

	ItemCap = 10

	db := testdb.NewTestDB()
	{

		pushCh := make(chan struct{}, 1)
		cb := func(err errcode.Error, _ int64) {
			select {
			case pushCh <- struct{}{}:
			}
		}

		PushMsg(db, "sniperHW", "test", []byte("hello"), time.Second, cb)
		<-pushCh

		PushMsg(db, "sniperHW", "test", []byte("hello1"), time.Second, cb)
		<-pushCh

		ch := make(chan struct{}, 1)

		PullMsg(db, "sniperHW", "test", 0, time.Second, func(err errcode.Error, version int64, items []*Item) {
			assert.Equal(t, version, int64(2))
			assert.Equal(t, 2, len(items))
			select {
			case ch <- struct{}{}:
			}
		})

		<-ch

		PullMsg(db, "sniperHW", "test", 1, time.Second, func(err errcode.Error, version int64, items []*Item) {
			assert.Equal(t, version, int64(2))
			assert.Equal(t, 1, len(items))
			select {
			case ch <- struct{}{}:
			}
		})

		<-ch

		PullMsg(db, "sniperHW", "test", 2, time.Second, func(err errcode.Error, version int64, items []*Item) {
			assert.Equal(t, version, int64(2))
			assert.Equal(t, 0, len(items))
			select {
			case ch <- struct{}{}:
			}
		})

		<-ch

	}

	{
		pushCh := make(chan struct{}, 1)
		cb := func(err errcode.Error, _ int64) {
			select {
			case pushCh <- struct{}{}:
			}
		}

		for i := 0; i < 10; i++ {
			PushMsg(db, "sniperHW2", "test", []byte(fmt.Sprintf("hello_%d", i)), time.Second, cb)
			<-pushCh
		}

		ch := make(chan struct{}, 1)

		PullMsg(db, "sniperHW2", "test", 0, time.Second, func(err errcode.Error, version int64, items []*Item) {
			assert.Equal(t, version, int64(10))
			assert.Equal(t, 10, len(items))
			for i := 0; i < 10; i++ {
				assert.Equal(t, fmt.Sprintf("hello_%d", i), string(items[i].Content))
			}
			select {
			case ch <- struct{}{}:
			}
		})

		<-ch

		//将item_0覆盖
		PushMsg(db, "sniperHW2", "test", []byte(fmt.Sprintf("hello_%d", 10)), time.Second, cb)
		<-pushCh

		PullMsg(db, "sniperHW2", "test", 0, time.Second, func(err errcode.Error, version int64, items []*Item) {
			assert.Equal(t, version, int64(11))
			assert.Equal(t, 10, len(items))
			for i := 0; i < 10; i++ {
				assert.Equal(t, fmt.Sprintf("hello_%d", i+1), string(items[i].Content))
			}
			select {
			case ch <- struct{}{}:
			}
		})

		<-ch

		//将item_1覆盖
		PushMsg(db, "sniperHW2", "test", []byte(fmt.Sprintf("hello_%d", 11)), time.Second, cb)
		<-pushCh

		PullMsg(db, "sniperHW2", "test", 0, time.Second, func(err errcode.Error, version int64, items []*Item) {
			assert.Equal(t, version, int64(12))
			assert.Equal(t, 10, len(items))
			for i := 0; i < 10; i++ {
				assert.Equal(t, fmt.Sprintf("hello_%d", i+2), string(items[i].Content))
			}
			select {
			case ch <- struct{}{}:
			}
		})

		<-ch

		PullMsg(db, "sniperHW2", "test", 12, time.Second, func(err errcode.Error, version int64, items []*Item) {
			select {
			case ch <- struct{}{}:
			}
		})

		<-ch

		PullMsg(db, "sniperHW3", "test", 0, time.Second, func(err errcode.Error, version int64, items []*Item) {
			assert.Equal(t, version, int64(0))
			assert.Equal(t, 0, len(items))
			select {
			case ch <- struct{}{}:
			}
		})

		<-ch

	}

}
