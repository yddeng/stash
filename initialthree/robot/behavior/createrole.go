package behavior

import (
	codecs "initialthree/codec/cs"
	csmsg "initialthree/protocol/cs/message"
	"initialthree/robot/net"
	"initialthree/robot/types"

	. "github.com/GodYY/bevtree"
	"github.com/golang/protobuf/proto"
)

const createRole = BevType("createrole")

func init() {
	regBevType(createRole, func() Bev { return new(BevCreateRole) })
}

type BevCreateRole struct{}

func (BevCreateRole) BevType() BevType { return createRole }

func (b *BevCreateRole) CreateInstance() BevInstance {
	return &BevCreateRoleInst{}
}

func (b *BevCreateRole) DestroyInstance(bi BevInstance) {
}

type BevCreateRoleInst struct {
	bev
}

func (b *BevCreateRoleInst) BevType() BevType { return createRole }

func (b *BevCreateRoleInst) OnInit(ctx Context) bool {
	b.bev.OnInit(ctx)

	if !b.player.IsStatus(types.Status_IsFirstlogin) {
		b.player.Debugf("not first login")
		b.terminate(true)
		return true
	}

	creareRoleMsg := &csmsg.CreateRoleToS{
		Name: proto.String(b.player.RoleName()),
	}
	b.sendMessage(creareRoleMsg, b.onCreateRole)

	b.player.Debugf("request to create role")

	return true
}

func (b *BevCreateRoleInst) OnTerminate(ctx Context) {
	b.bev.OnTerminate(ctx)
}

func (b *BevCreateRoleInst) onCreateRole(r player, msg *codecs.Message) bool {
	if !net.IsMessageOK(msg) {
		// 创角失败
		r.Errorf("create role failed: %s", net.GetErrCodeStr(msg.GetErrCode()))
		b.terminate(false)
		return false
	}

	// 创角成功
	r.Infof("create role successfully")
	b.player.UnsetStatus(types.Status_IsFirstlogin)
	b.terminate(true)
	return false
}
