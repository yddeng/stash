# hutils

> 不使用第三方库，不产生额外依赖


> 为了尽可能减少依赖，一些有用的代码贴在这里

```go
import "golang.org/x/crypto/pbkdf2"

// 密码加密
func EncodePassword(password, salt string) string {
	dk := pbkdf2.Key([]byte(password), []byte(salt), 4096, 32, sha256.New)
	return hex.EncodeToString(dk)
}

```