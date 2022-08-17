package main

import (
	"fmt"
	"initialthree/protocol/ss/proto_def"
	"os"
	"strings"
)

var template string = `
package [s1]

import (
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/pkg/rpc"
	ss_rpc "initialthree/protocol/ss/rpc"
	"time"
)

type [s2] struct {
	replyer_ *rpc.RPCReplyer
}

func (this *[s2]) Reply(result *ss_rpc.[s3]) {
	this.replyer_.Reply(result,nil)
}

/*
func (this *[s2]) DropResponse() {
	this.replyer_.DropResponse()
}
*/

/*
func (this *[s2]) Error(err error) {
	this.replyer_.Reply(nil,err)
}
*/

func (this *[s2]) GetChannel() rpc.RPCChannel {
	return this.replyer_.GetChannel()
}


type [s5] interface {
	OnCall(*[s2],*ss_rpc.[s4])
}

func Register(methodObj [s5]) {
	f := func(r *rpc.RPCReplyer, arg interface{}) {
		replyer_ := &[s2]{replyer_:r}
		methodObj.OnCall(replyer_,arg.(*ss_rpc.[s4]))
	}

	cluster.RegisterMethod(&ss_rpc.[s4]{},f)
}

func AsynCall(peer addr.LogicAddr,arg *ss_rpc.[s4],timeout time.Duration,cb func(*ss_rpc.[s3],error)) {
	callback := func(r interface{},e error) {
		if nil != r {
			cb(r.(*ss_rpc.[s3]),e)
		} else {
			cb(nil,e)
		}
	}
	cluster.AsynCall(peer,arg,timeout,callback)
}

/*
func SyncCall(peer addr.LogicAddr,arg *ss_rpc.[s4],timeout time.Duration) (ret *ss_rpc.[s3], err error) {
	waitC := make(chan struct{})
	f := func(ret_ *ss_rpc.[s3], err_ error) {
		ret = ret_
		err = err_
		close(waitC)
	}
	AsynCall(peer,arg,timeout,f)
	<-waitC
	return
}
*/
`

func gen_rpc(array []proto_def.St) {

	for _, v := range array {

		path := v.Name
		filename := fmt.Sprintf("%s/%s.go", path, v.Name)
		os.MkdirAll(path, os.ModePerm)
		f, err := os.OpenFile(filename, os.O_RDWR, os.ModePerm)
		if err != nil {
			if os.IsNotExist(err) {
				f, err = os.Create(filename)
				if err != nil {
					fmt.Printf("------ error -------- create %s failed:%s", filename, err.Error())
					return
				}
			} else {
				fmt.Printf("------ error -------- open %s failed:%s", filename, err.Error())
				return
			}
		}
		err = os.Truncate(filename, 0)
		if err != nil {
			fmt.Printf("------ error -------- Truncate %s failed:%s", filename, err.Error())
			f.Close()
			return
		}

		content := template
		content = strings.Replace(content, "[s1]", v.Name, -1)
		content = strings.Replace(content, "[s2]", strings.Title(v.Name)+"Replyer", -1)
		content = strings.Replace(content, "[s3]", strings.Title(v.Name)+"Resp", -1)
		content = strings.Replace(content, "[s4]", strings.Title(v.Name)+"Req", -1)
		content = strings.Replace(content, "[s5]", strings.Title(v.Name), -1)
		_, err = f.WriteString(content)

		//fmt.Printf(content)

		if nil != err {
			fmt.Printf("------ error -------- %s Write error:%s\n", filename, err.Error())
			return
		} else {
			//fmt.Printf("%s Write ok\n", filename)
		}

		f.Close()
	}
}

func main() {
	fmt.Printf("gen_rpc\n")
	gen_rpc(proto_def.SS_rpc)
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("gen_rpc ok!\n")
}
