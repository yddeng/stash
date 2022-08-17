package timer

import (
	"fmt"
	"testing"
	"time"
)

func tick() {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case now := <-ticker.C:
			Loop(now.Unix())
		}
	}
}

func TestOnce(t *testing.T) {
	go tick()
	now := time.Now().Unix()
	Once(now+1, func() {
		fmt.Println(1, now, time.Now().Unix())
	})
	s1 := Once(now+3, func() {
		fmt.Println(3, now, time.Now().Unix())
	})
	s1.Reset(now + 5)

	s := Once(now+3, func() {
		fmt.Println(4, now, time.Now().Unix())
	})
	go func() {
		time.Sleep(time.Second * 4)
		s.Reset(now + 1)
		s.Reset(now + 3)
	}()

	select {}
}
