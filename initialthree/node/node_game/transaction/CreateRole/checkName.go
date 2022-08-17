package CreateRole

import (
	"initialthree/node/common/wordsFilter"
	"strings"
)

func (this *transactionCreateRole) checkName(name string) bool {
	str := []rune(name)
	if strings.TrimSpace(name) == "" || len(str) > 16 {
		return false
	}

	return !wordsFilter.Check(name)
}
