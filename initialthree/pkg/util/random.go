package util

import "math/rand"

func Random(min, max int32) int32 {
	return rand.Int31n(max-min+1) + min
}
