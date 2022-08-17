package behavior

import (
	codecs "initialthree/codec/cs"
	"initialthree/robot/types"

	. "github.com/GodYY/bevtree"

	"github.com/golang/protobuf/proto"
)

type player = types.Robot
type timerID = types.TimerID

type bev struct {
	player player
	result Result
}

func (b *bev) OnInit(ctx Context) bool {
	b.player = ctx.UserData().(player)
	b.result = Running
	return true
}

func (b *bev) OnUpdate(ctx Context) Result {
	return b.result
}

func (b *bev) OnTerminate(ctx Context) {
	b.player = nil
}

func (b *bev) terminate(result bool) {
	if result {
		b.result = Success
	} else {
		b.result = Failure
	}
}

type msgCallback = func(player, *codecs.Message) bool

func (b *bev) sendMessageVar(msg proto.Message, cb ...msgCallback) {
	n := len(cb)
	copy(cb[1:], cb[0:n])
	cb[0] = b.msgMiddleWare
	b.player.SendMessage(msg, cb...)
}

func (b *bev) sendMessage(msg proto.Message, cb msgCallback) {
	b.player.SendMessage(msg, b.msgMiddleWare, cb)
}

func (b *bev) msgMiddleWare(p player, msg *codecs.Message) bool {
	if p != b.player || b.result != Running {
		return false
	}

	return true
}

var framework = NewFramework()

func regNodeType(nodeType NodeType, nodeCreator func() Node, taskCreator func() Task) {
	framework.RegisterNodeType(nodeType, nodeCreator, taskCreator)
}

func regBevType(bevType BevType, creator func() Bev) {
	framework.RegisterBevType(bevType, creator)
}

func CreateEntity(tree string, player player) (Entity, error) {
	return framework.CreateEntity(tree, player)
}

func CreateExporter() *Exporter {
	return NewExporter(framework)
}

func Init(configPath string) error {
	return framework.Init(configPath)
}
