package idGenerator

import (
	"fmt"
	"initialthree/cluster"
	"initialthree/node/common/db"
	"initialthree/zaplogger"
	"sync"
	"testing"
	"time"
)

func TestGen_Reserved(t *testing.T) {
	logger := zaplogger.NewZapLogger("idGen", ".", "debug", 100, 7, 5, true)
	if err := db.FlyfishInit("cluster", []string{"flyfish@global@10.128.2.123:10012"}, logger); err != nil {
		panic(err)
	}

	go cluster.GetEventQueue().Run()

	Register("test", db.GetFlyfishClient("global"), func(i int64) int64 {
		return i
	}, true)

	now := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		GetIDGen("test").GenID(func(i int64, e error) {
			wg.Done()
			fmt.Println(i, e)
		})
	}
	wg.Wait()
	fmt.Println(time.Now().Sub(now).String())

}
