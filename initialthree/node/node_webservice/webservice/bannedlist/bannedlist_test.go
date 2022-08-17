package bannedlist

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	now := time.Now()
	e := time.Date(now.Year(), now.Month(), now.Day(), 17, 26, 0, 0, time.Local)

	version := map[string]interface{}{
		"userId":      "ydd015",
		"expiredTime": e.Unix(),
	}

	req, _ := dhttp.PostJson("http://10.128.2.123:41801/bannedlist/set", version)
	fmt.Println(req.ToString())
}
