package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gim/im/protocol"
	"initialthree/node/client/database/base"
	attr2 "initialthree/node/common/attr"
	"initialthree/node/common/db"
	"initialthree/node/node_game/module/attr"
	"initialthree/node/node_game/module/character"
	"strings"
	"time"
)

func dispose(data map[string]interface{}) map[string]interface{} {
	ret := make(map[string]interface{}, len(data)-2)
	for k, v := range data {
		if !(k == "__key__" || k == "__version__") {
			ret[k] = v
		}
	}
	//fmt.Println(data, ret)
	return ret
}

var flyClient *db.Client

func copyTableRow(tabName, oldID, newID string, skipField ...string) {
	ret, _ := flyClient.Get(tabName, oldID)
	fields := dispose(ret)
	if len(fields) != 0 {
		if len(skipField) != 0 {
			// 跳过
			if _, ok := ret[skipField[0]]; ok {
				ret2, _ := flyClient.Get(tabName, newID)
				fields[skipField[0]] = ret2[skipField[0]]
			}
		}

		if err := flyClient.Upsert(tabName, newID, fields); err != nil {
			panic(tabName + err.Error())
		}
		//fmt.Printf("%s copy ok \n", tabName)
	}
}

func getID(u string) (string, bool) {
	ret, err := flyClient.Get("game_user_0", u)
	if err != nil || len(ret) == 0 {
		return "", false
	}

	return fmt.Sprintf("%d", ret["id_1"].(int64)), true
}

func copyUser(o, n string) {
	area := 1

	oldUser := fmt.Sprintf("%s:%d", o, area)
	newUser := fmt.Sprintf("%s:%d", n, area)

	oldID, oldExist := getID(oldUser)
	if !oldExist {
		fmt.Printf("==============> old user %s not exist\n", n)
		return
	}
	newID, newExist := getID(newUser)
	if !newExist {
		fmt.Printf("==============> new user %s not exist\n", n)
		return
	}

	copyTableRow("user_module_data_0", oldID, newID, "base_1")
	copyTableRow("character_0", oldID, newID)
	copyTableRow("quest_0", oldID, newID)
	copyTableRow("user_assets_0", oldID, newID)
	copyTableRow("backpack_0", oldID, newID)
	copyTableRow("equip_0", oldID, newID)
	copyTableRow("weapon_0", oldID, newID)
	copyTableRow("drawcard_0", oldID, newID)
	copyTableRow("scarsingrain_0", oldID, newID)
	copyTableRow("main_dungeons_0", oldID, newID)
	copyTableRow("rewardquest_0", oldID, newID)
	copyTableRow("materialdungeon_0", oldID, newID)
	copyTableRow("temporary_0", oldID, newID)
	copyTableRow("shop_0", oldID, newID)
	copyTableRow("worldquest_0", oldID, newID)
	copyTableRow("sign_0", oldID, newID)

	fmt.Printf("copy (%s,%s) to (%s,%s) ok\n", oldUser, oldID, newUser, newID)
}

func insertImUser(uID string) {
	attrRet, _ := flyClient.Get("user_module_data_0", uID)
	var userAttr attr.AttrData
	json.Unmarshal(attrRet["attr_1"].([]byte), &userAttr)

	type userBaseData struct {
		UserID        string  `json:"user_id"`
		Name          string  `json:"name"` // 角色名
		Sex           int32   `json:"sex"`
		BuildTime     int64   `json:"build_time"` // 创建账号的时间
		Birthday      string  `json:"birthday"`
		Signature     string  `json:"signature"`
		CharacterList []int32 `json:"character_list"`
		Portrait      int32   `json:"portrait"`
		PortraitFrame int32   `json:"portrait_frame"`
		Card          int32   `json:"card"`
		UIDCounter    uint32  `json:"uc"`
	}
	var userBase userBaseData
	json.Unmarshal(attrRet["base_1"].([]byte), &userBase)

	charaRet, _ := flyClient.Get("character_0", uID)
	var userChara = make([]map[int32]*character.Character, 10)
	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("slot%d_1", i)
		json.Unmarshal(charaRet[name].([]byte), &userChara[i])
	}

	find := func(id int32) int32 {
		for _, v := range userChara {
			for k, c := range v {
				if k == id {
					return c.Level
				}
			}
		}
		return 1
	}

	charas := []string{}
	for _, v := range userBase.CharacterList {
		charas = append(charas, fmt.Sprintf("%d#%d", v, find(v)))
	}

	extra := []*protocol.Extra{
		{Key: proto.String("UserID"), Value: proto.String(uID)},
		{Key: proto.String("名字"), Value: proto.String(userBase.Name)},
		{Key: proto.String("头像"), Value: proto.String(fmt.Sprintf("%d", userBase.Portrait))},
		{Key: proto.String("生日"), Value: proto.String(fmt.Sprintf("%s", userBase.Birthday))},
		{Key: proto.String("签名"), Value: proto.String(fmt.Sprintf("%s", userBase.Signature))},
		{Key: proto.String("等级"), Value: proto.String(fmt.Sprintf("%d", userAttr.Attrs[attr2.Level].Val))},
		{Key: proto.String("展示角色"), Value: proto.String(strings.Join(charas, ","))},
	}

	fmt.Println(uID, extra)

	sqlStatement := `
	  INSERT INTO "users" (id,create_at,update_at,extra)
	  VALUES($1, $2, $3, $4)
	  ON conflict(id) DO
	  UPDATE SET create_at = $2, update_at = $3, extra = $4;`

	create_at := time.Now().Unix()
	data, _ := json.Marshal(extra)
	if err := flyClient.Exec(sqlStatement, uID, create_at, create_at, data); err != nil {
		panic(err.Error())
	}

}

func addFriend(u1, u2 string) {
	if u1 > u2 {
		u1, u2 = u2, u1
	}

	sqlStatement := `
INSERT INTO "friend" (id,user1_id,user2_id,create_at,status)
VALUES($1, $2, $3, $4,$5)
ON conflict(id) DO
UPDATE SET  create_at = $4, status = $5;`
	key := fmt.Sprintf("%s_%s", u1, u2)
	create_at := time.Now().Unix()
	statue := int(4)

	if err := flyClient.Exec(sqlStatement, key, u1, u2, create_at, statue); err != nil {
		panic(err.Error())
	}
}

func main() {
	flyClient = base.InitClient("pgsql", "10.50.31.13", 5432, "teacher2", "sgzr2", "sniperHWfeiyu2019")
	//flyClient = base.InitClient("pgsql", "127.0.0.1", 5432, "yidongdeng", "dbuser", "123456")

	tmp1 := "test001"
	c1 := []string{
		"dbc5807cc6a94b8c12bda585bb3759de",
		"583890fa0acfe5eb2418a16c446e21c7",
		"40a3c3bb062d5486f5ef9c7bc46555ee",
		"427adb4b0c780dd996216cf8531bfe9e",
		"10fdb8c8c00c9da6c9aed6f3123eca61",
		"df6d7408413a31553babd275cd2bebec",
		"7359860e7626a71d6a08bca408e70a86",
		"6862c87b174b55bfe0b259862afc4cc6",
		"640c2b38cc71cd31d6c3ef1105ff24ea",
	}
	for _, c := range c1 {
		copyUser(tmp1, c)
	}

	tmp2 := "test002"
	c2 := []string{
		"177c9b3310e44219f63240c1ad6caec1",
		"73eb0fa8122435c80cc8f7905e1fad84",
		"7711d3a3be329e901daebb8f7724f53b",
		"60c086e03fbd9a289db3753bea7f5d4a",
		"3d229c46272f97b053eb7d3f8033ac56",
		"eae018483dc84138b917f0a9288414c9",
		"56f87e0c26bdfe830f38e9bbf453b77e",
		"47d5d575b8bc1e9af5f3850630c072fd",
		"fe7de7381dcb1cb6895fa5f911f072f4",
	}
	for _, c := range c2 {
		copyUser(tmp2, c)
	}

	tmp3 := "test003"
	c3 := []string{
		"aa06234bf4dc1c927b1063aa23eb7bbf",
		"3cb919e33e322a93bf1b93b8bf0bec27",
		"fb4cd05bd164aebd42e4e9254d9a96c0",
		"11eef5e0cefd763fe78261d1a0405502",
		"729380782dfcbd7f62ead55d3df969c3",
		"a58d5ccb79baf20e721c7e8dd56b7838",
		"7ebd980da028170286764a2a329437fe",
		"a39d2173f1d435e0c2bc63edd478a417",
		"8707f502d4adafd3fd7a8bfd1cb11045",
	}
	for _, c := range c3 {
		copyUser(tmp3, c)
	}

	// addFriend

	frinds := []string{
		"dbc5807cc6a94b8c12bda585bb3759de",
		"583890fa0acfe5eb2418a16c446e21c7",
		"40a3c3bb062d5486f5ef9c7bc46555ee",
		"427adb4b0c780dd996216cf8531bfe9e",
		"10fdb8c8c00c9da6c9aed6f3123eca61",
		"df6d7408413a31553babd275cd2bebec",
		"7359860e7626a71d6a08bca408e70a86",
		"6862c87b174b55bfe0b259862afc4cc6",
		"640c2b38cc71cd31d6c3ef1105ff24ea",
	}

	ids := make([]string, 0, len(frinds))
	for _, id := range frinds {
		uid, ok := getID(fmt.Sprintf("%s:%d", id, 1))
		if ok {
			ids = append(ids, uid)
			insertImUser(uid)
		}
	}

	for i := 0; i < len(ids); i++ {
		for j := 0; j < len(ids); j++ {
			if i != j {
				addFriend(ids[i], ids[j])
			}
		}
	}

}
