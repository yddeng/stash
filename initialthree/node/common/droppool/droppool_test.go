package droppool

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	table "initialthree/node/table/excel"
)

func load() {
	table.Load("../../configs/GameDesigner/Excel")
	rand.Seed(time.Now().Unix())
}

func TestDropAward(t *testing.T) {
	load()

	award1 := DropAward(1)
	fmt.Println("award1")
	for _, v := range award1.Infos {
		fmt.Println(v)
	}

	award2 := DropAward(2)
	fmt.Println("award2")
	for _, v := range award2.Infos {
		fmt.Println(v)
	}

	award := AllAwardsToOne([]*Award{award1, award2})
	fmt.Println("award")
	for _, v := range award.Infos {
		fmt.Println(v)
	}
}

func TestDropCard(t *testing.T) {
	load()

	award1 := DropCard(10001)
	fmt.Println("award1")
	for _, v := range award1.Infos {
		fmt.Println(v)
	}
}
