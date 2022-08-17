package util

import (
	"fmt"
	"strings"
)

func ParseAddress(address string) string {
	t := strings.Split(address, ":")
	return fmt.Sprintf("0.0.0.0:%s", t[1])
}
