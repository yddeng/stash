package announcement

import (
	"encoding/json"
	"fmt"
	"github.com/yddeng/dnet/dhttp"
	"testing"
)

func TestGet(t *testing.T) {
	version := map[string]interface{}{
		"version": 0,
	}

	req, _ := dhttp.PostJson("http://10.128.2.123:41801/announcement/get", version)
	fmt.Println(req.ToString())
}

func TestDel(t *testing.T) {

	req, _ := dhttp.PostJson("http://10.128.2.123:41801/announcement/delete", map[string]interface{}{
		"id": 1,
	})
	fmt.Println(req.ToString())
}

func TestAdd(t *testing.T) {

	content := `{
			"id": 1,
			"type": "AnnouncementType_System",
			"title": "第一行\n第二行",
			"smallTitle": "不知道叫什么",
			"startTime": 1649586412,
			"expireTime": 0,
			"remind": false,
			"content": [{
				"type": "0",
				"imageSkip": 1,
				"text": "",
				"image": "Announce_ad_1"
			}, {
				"type": "1",
				"imageSkip": 1,
				"text": "测试测试内容\n测试",
				"image": "Announce_banner_2"
			}]
		}`

	/*
		content := `{
			"id": 2,
			"type": "AnnouncementType_Activity",
			"title": "第二行\n第二行",
			"smallTitle": "不知道叫什么2",
			"startTime": 1649586412,
			"expireTime": 0,
			"remind": true,
			"content": [{
				"type": "0",
				"imageSkip": 0,
				"text": "",
				"image": "Announce_ad_1"
			}]
		}`

	*/

	var ann *Announcement

	err := json.Unmarshal([]byte(content), &ann)
	if err != nil {
		panic(err)
	}

	req, _ := dhttp.PostJson("http://10.128.2.123:41801/announcement/add", ann)
	fmt.Println(req.ToString())
}
