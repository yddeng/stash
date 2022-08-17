package regexp

import (
	"regexp"
)

// 验证qq号
// 首位不能是0
// 5-11 位的数字
var regQQ = regexp.MustCompile(`^[1-9]\d{4,10}$`)

func MatchQQ(qq string) bool {
	return regQQ.Match([]byte(qq))
}

// 微信号
// 以字母开头（不区分大小写）
// 6—20个字母、数字、下划线和减号
var regWX = regexp.MustCompile(`^[a-zA-Z][-_a-zA-Z0-9]{5,19}$`)

func MatchWX(wx string) bool {
	return regWX.Match([]byte(wx))
}

var regEmail = regexp.MustCompile(`^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`)

func MatchEmail(email string) bool {
	return regEmail.Match([]byte(email))
}
