package main

import (
	"fmt"
	"initialthree/protocol/cs/proto_def"
	"os"
	//"strings"
)

//产生协议注册文件
func gen_lua_id(out_path string) {

	f, err := os.OpenFile(out_path, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(out_path)
			if err != nil {
				fmt.Printf("------ error -------- create %s failed:%s", out_path, err.Error())
				return
			}
		} else {
			fmt.Printf("------ error -------- open %s failed:%s", out_path, err.Error())
			return
		}
	}

	err = os.Truncate(out_path, 0)

	if err != nil {
		fmt.Printf("------ error -------- Truncate %s failed:%s", out_path, err.Error())
		return
	}

	str := `
return {
`

	for _, v := range proto_def.CS_message {
		str = str + fmt.Sprintf(`	[%d] = {name='%s',desc='%s'},`, v.MessageID, v.Name, v.Desc) + "\n"
	}

	str = str + "}\n"

	_, err = f.WriteString(str)

	//fmt.Printf(str)

	if nil != err {
		fmt.Printf("------ error -------- %s Write error:%s\n", out_path, err.Error())
	} else {
		fmt.Printf("%s Write ok\n", out_path)
	}

	f.Close()

}

func main() {

	os.MkdirAll("../message", os.ModePerm)
	gen_lua_id("../messageID.lua")
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("cs gen_proto_lua ok!\n")
}
