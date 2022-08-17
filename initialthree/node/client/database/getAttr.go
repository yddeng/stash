package main

import (
	"fmt"
	"initialthree/node/client/database/base"
	"initialthree/node/common/attr"
	"initialthree/pkg/json"
	"os"
	"strconv"
)

func main() {
	userID := os.Args[1]

	id := base.GetID(userID)
	if id == 0 {
		fmt.Println("没有相关数据")
		return
	}

	cli := base.GetClient()
	ret, err := cli.Get("user_module_data", strconv.Itoa(int(id)))
	if err != nil {
		fmt.Println("查询 user_module_data 出错", err)
		return
	}

	data := ret["attr"].([]byte)
	if data == nil || len(data) == 0 {
		fmt.Println("没有相关数据")
		return
	}

	type AttrData struct {
		Attrs []*attr.Attr `json:"As"`
	}

	attrData := AttrData{}

	if err := json.Unmarshal(data, &attrData); err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range attrData.Attrs {
		fmt.Println(v.ID, v.Val)
	}
}
