package assert

import "log"

func Assert(b bool, msg interface{}) {
	if !b {
		log.Panic(msg)
	}
}

func AssertF(b bool, f string, args ...interface{}) {
	if !b {
		log.Panicf(f, args...)
	}
}

func Equal(a, b interface{}, msg interface{}) {
	Assert(a == b, msg)
}

func EqualF(a, b interface{}, f string, args ...interface{}) {
	AssertF(a == b, f, args...)
}

func NotEqual(a, b interface{}, msg interface{}) {
	Assert(a != b, msg)
}

func NotEqualF(a, b interface{}, f string, args ...interface{}) {
	AssertF(a != b, f, args...)
}
