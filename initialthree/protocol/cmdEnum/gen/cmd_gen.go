package main

import (
	"fmt"
	csproto_def "initialthree/protocol/cs/proto_def"
	ssproto_def "initialthree/protocol/ss/proto_def"
	"os"
	"strings"
)

var cmd_tmp string = `
package cmdEnum

const (
	//cs消息
%s
	//rpc消息
%s
	//ss消息
%s
)
`

//产生协议注册文件
func gen_cmd(out_path string) {

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
	//用户定义ID开始区段
	ss_start := 1000

	cs_str, rpc_str, ss_str := "", "", ""
	for _, v := range csproto_def.CS_message {
		cs_str += fmt.Sprintf("    CS_%s uint16 = %d\n", strings.Title(v.Name), v.MessageID)
	}
	for _, v := range ssproto_def.SS_rpc {
		rpc_str += fmt.Sprintf("    RPC_%s uint16 = %d\n", strings.Title(v.Name), v.MessageID)
	}
	for _, v := range ssproto_def.SS_message {
		ss_str += fmt.Sprintf("    SS_%s uint16 = %d\n", strings.Title(v.Name), v.MessageID+ss_start)
	}

	content := fmt.Sprintf(cmd_tmp, cs_str, rpc_str, ss_str)

	_, err = f.WriteString(content)

	if nil != err {
		fmt.Printf("------ error -------- %s Write error:%s\n", out_path, err.Error())
	} else {
		fmt.Printf("%s Write ok\n", out_path)
	}

	f.Close()

}

func main() {
	fmt.Printf("start gen_cmd ...\n")

	gen_cmd("../cmdEnum.go")
	fmt.Printf("gen_cmd ok  \n")
}
