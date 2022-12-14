package main

import (
	"fmt"
	"initialthree/protocol/ss/proto_def"
	"os"
	"strings"
)

var message_template string = `syntax = "proto2";
package ssmessage;
option go_package = "initialthree/protocol/ss/ssmessage";

message %s {}`

var rpc_template string = `syntax = "proto2";
package rpc;
option go_package = "initialthree/protocol/ss/rpc";

message %s_req {}

message %s_resp {}`

var message = 1
var rpc = 2

func gen_proto(tt int, array []proto_def.St, out_path string) {

	if tt == message {
		fmt.Printf("gen_proto message ............\n")
	} else {
		fmt.Printf("gen_proto rpc ............\n")
	}

	for _, v := range array {
		filename := fmt.Sprintf("%s/%s.proto", out_path, v.Name)
		//检查文件是否存在，如果存在跳过不存在创建
		f, err := os.Open(filename)
		if nil != err && os.IsNotExist(err) {
			f, err = os.Create(filename)
			if nil == err {
				fmt.Printf("add %s.proto \n", v.Name)
				var content string
				if tt == message {
					content = fmt.Sprintf(message_template, v.Name)
				} else {
					content = fmt.Sprintf(rpc_template, v.Name, v.Name)
				}

				_, err = f.WriteString(content)

				if nil != err {
					fmt.Printf("------ error -------- %s Write error:%s\n", v.Name, err.Error())
				}

				f.Close()

			} else {
				fmt.Printf("------ error -------- %s Create error:%s\n", v.Name, err.Error())
			}
		} else if nil != f {
			//fmt.Printf("%s.proto exist skip\n", v.Name)
			f.Close()
		}
	}
}

var register_template string = `
package ss
import (
	"initialthree/codec/pb"
	"initialthree/protocol/ss/ssmessage"
	"initialthree/protocol/ss/rpc"
)

func init() {
	//普通消息
%s
	//rpc请求
%s
	//rpc响应
%s
}
`

//产生协议注册文件
func gen_register(out_path string) {

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
	id_start := 1000

	message_str := ""

	nameMap := map[string]bool{}
	idMap := map[int]bool{}

	for _, v := range proto_def.SS_message {

		if ok, _ := nameMap[v.Name]; ok {
			panic("duplicate ss message:" + v.Name)
		}

		if ok, _ := idMap[v.MessageID]; ok {
			panic(fmt.Sprintf("duplicate ss messageID: %d", v.MessageID))
		}

		nameMap[v.Name] = true
		idMap[v.MessageID] = true

		message_str = message_str + fmt.Sprintf(`	pb.Register("ss",&ssmessage.%s{},%d)`, strings.Title(v.Name), v.MessageID+id_start) + "\n"
	}

	rpc_req_str := ""
	rpc_resp_str := ""

	nameMap = map[string]bool{}
	idMap = map[int]bool{}

	for _, v := range proto_def.SS_rpc {

		if ok, _ := nameMap[v.Name]; ok {
			panic("duplicate rpc message:" + v.Name)
		}

		if ok, _ := idMap[v.MessageID]; ok {
			panic(fmt.Sprintf("duplicate rpc messageID: %d", v.MessageID))
		}

		nameMap[v.Name] = true
		idMap[v.MessageID] = true

		rpc_req_str = rpc_req_str + fmt.Sprintf(`	pb.Register("rpc_req",&rpc.%sReq{},%d)`, strings.Title(v.Name), v.MessageID+id_start) + "\n"
		rpc_resp_str = rpc_resp_str + fmt.Sprintf(`	pb.Register("rpc_resp",&rpc.%sResp{},%d)`, strings.Title(v.Name), v.MessageID+id_start) + "\n"
	}

	content := fmt.Sprintf(register_template, message_str, rpc_req_str, rpc_resp_str)

	_, err = f.WriteString(content)

	//fmt.Printf(content)

	if nil != err {
		fmt.Printf("------ error -------- %s Write error:%s\n", out_path, err.Error())
	} else {
		fmt.Printf("%s Write ok\n", out_path)
	}

	f.Close()

}

func main() {
	fmt.Printf("gen_ss\n")

	os.MkdirAll("../ssmessage", os.ModePerm)
	os.MkdirAll("../rpc", os.ModePerm)
	gen_proto(message, proto_def.SS_message, "../proto/ssmessage")
	gen_proto(rpc, proto_def.SS_rpc, "../proto/rpc")
	gen_register("../register.go")
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("ss gen_proto_go ok!\n")
}
