package uniLocker

import (
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/cluster/addr"
	"time"
)

const (
	table = "nodelock"
	field = "phyaddr"
)

type UniLokcer struct {
	flyfishClient *flyfish.Client
	a             addr.Addr
}

func New(flyfishClient *flyfish.Client) *UniLokcer {
	return &UniLokcer{
		flyfishClient: flyfishClient,
	}
}

func (this *UniLokcer) Lock(a addr.Addr) bool {
	this.a = a
	set := this.flyfishClient.CompareAndSetNx(table, a.Logic.String(), field, "", a.Net.String())

	for i := 0; i < 3; i++ {
		ret := set.Exec()
		switch errcode.GetCode(ret.ErrCode) {
		case 0:
			return true
		case errcode.Errcode_cas_not_equal:
			if a.Net.String() == ret.Value.GetString() {
				return true
			} else {
				return false
			}
		default:
			time.Sleep(time.Second)
		}
	}

	return false
}

func (this *UniLokcer) Unlock() {
	set := this.flyfishClient.CompareAndSetNx(table, this.a.Logic.String(), field, this.a.Net.String(), "")
	for i := 0; i < 3; i++ {
		ret := set.Exec()
		switch errcode.GetCode(ret.ErrCode) {
		case 0:
			return
		case errcode.Errcode_cas_not_equal:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
