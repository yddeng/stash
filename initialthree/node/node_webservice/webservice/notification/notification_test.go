package notification

import (
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestGet(t *testing.T) {
	req, _ := dhttp.PostJson("http://10.128.2.123:41801/notification/get", struct {
	}{})
	fmt.Println(req.ToString())
}

func TestGetAll(t *testing.T) {
	req, _ := dhttp.PostJson("http://10.128.2.123:41801/notification/getAll", struct {
	}{})
	fmt.Println(req.ToString())
}

func TestSetClosed(t *testing.T) {
	req, _ := dhttp.PostJson("http://10.128.2.123:41801/notification/setClosed", struct {
		Closed bool `json:"closed"`
	}{
		Closed: false,
	})
	fmt.Println(req.ToString())
}

func TestUpdate(t *testing.T) {
	req, _ := dhttp.PostJson("http://10.128.2.123:41801/notification/update", struct {
		Type         string        `json:"type"`
		Notification *Notification `json:"notification"`
	}{
		Type: "online",
	})
	fmt.Println(req.ToString())
}
