package otherStart

import (
	"strings"
)

const (
	PPROF      = "pprof"
	PROMETHEUS = "prometheus"
)

var (
	hasStarted = map[string]string{}
)

// 返回类型及地址
func DisposeArgs(args []string) map[string]string {
	hasStarted = map[string]string{}
	for _, arg := range args {
		k, v := dispose(arg)
		if k == "" {
			continue
		} else {
			hasStarted[k] = v
		}
	}
	return hasStarted
}

func dispose(arg string) (string, string) {
	s := strings.Split(arg, "@")
	if len(s) != 2 {
		return "", ""
	}

	switch s[0] {
	case PPROF, PROMETHEUS:
		return s[0], s[1]
	default:
		return "", ""
	}

}

func Has(name string) (string, bool) {
	v, ok := hasStarted[name]
	return v, ok
}
