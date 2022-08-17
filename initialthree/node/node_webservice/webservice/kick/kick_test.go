package kick

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestKickUser(t *testing.T) {
	req, _ := dhttp.PostJson("http://10.128.2.123:41801/kick/user", struct {
		UserID []string `json:"userID"`
	}{UserID: []string{"ydd001"}})

	fmt.Println(req.ToString())
}

func TestKickIP(t *testing.T) {
	req, _ := dhttp.PostJson("http://10.128.2.123:41801/kick/ip", struct {
		IP []string `json:"ip"`
	}{IP: []string{"10.128.2.123"}})

	fmt.Println(req.ToString())
}
