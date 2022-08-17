package hutils

import (
	"math/rand"
	"time"
)

const strBase = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var strMaxLen = len(strBase)

func InitRandSeed() {
	rand.Seed(time.Now().UnixNano())
}

func RandString(n int) string {
	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		buf[i] = strBase[rand.Intn(strMaxLen)]
	}
	return string(buf)
}
