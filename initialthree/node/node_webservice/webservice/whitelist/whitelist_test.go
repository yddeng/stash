package whitelist

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestSet(t *testing.T) {
	version := map[string]interface{}{
		"userId": "ydd015",
	}

	req, _ := dhttp.PostJson("http://10.128.2.123:41801/whitelist/set", version)
	fmt.Println(req.ToString())
}
