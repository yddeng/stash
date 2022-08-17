package behavior

import (
	"encoding/xml"
	"fmt"
	codecs "initialthree/codec/cs"
	csmsg "initialthree/protocol/cs/message"
	"initialthree/robot/net"

	. "github.com/GodYY/bevtree"
	"github.com/golang/protobuf/proto"
)

type GMCmdType int8

const (
	GMCmd_Unknown = GMCmdType(iota)
	GMCmd_AddAttr
	GMCmd_AddCharacter
	GMCmd_AddItem
	GMCmd_PassMainDungeon
	GMCmd_AddEquip
	GMCmd_AddWeapon
)

var gmCmdNames = [...]string{
	GMCmd_Unknown:         "unknown",
	GMCmd_AddAttr:         "addattr",
	GMCmd_AddCharacter:    "addcharacter",
	GMCmd_AddItem:         "additem",
	GMCmd_PassMainDungeon: "passmaindungeon",
	GMCmd_AddEquip:        "addequip",
	GMCmd_AddWeapon:       "addweapon",
}

var gmCmdName2Types = map[string]GMCmdType{}

func init() {
	for i, v := range gmCmdNames {
		gmCmdName2Types[v] = GMCmdType(i)
	}
}

func (t GMCmdType) String() string {
	return gmCmdNames[t]
}

func (t GMCmdType) MarshalText() ([]byte, error) {
	return ([]byte)(t.String()), nil
}

func (t *GMCmdType) UnmarshalText(bytes []byte) error {
	name := string(bytes)
	if tt, ok := gmCmdName2Types[name]; !ok {
		return xml.UnmarshalError(fmt.Sprintf("invalid GMCmdType name: %s", name))
	} else {
		*t = tt
		return nil
	}
}

type GMCmd struct {
	Type  GMCmdType `xml:"type,attr"`
	ID    int32     `xml:"id,attr,omitempty"`
	Count int32     `xml:"count,attr,omitempty"`

	data *csmsg.GmCmd
}

func (c *GMCmd) pack() *csmsg.GmCmd {
	if c.data == nil {
		c.data = &csmsg.GmCmd{
			Type:  proto.Int32(int32(c.Type)),
			ID:    proto.Int32(c.ID),
			Count: proto.Int32(c.Count),
		}
	}
	return c.data
}

const gameMaster = BevType("gamemaster")

func init() {
	regBevType(gameMaster, func() Bev { return new(BevGameMaster) })
}

type BevGameMaster struct {
	Desc string  `xml:"desc"`
	Cmd  []GMCmd `xml:"cmd>one"`
}

func (BevGameMaster) BevType() BevType { return gameMaster }

func (b *BevGameMaster) CreateInstance() BevInstance {
	return &BevGameMasterInst{
		BevGameMaster: b,
	}
}

func (b *BevGameMaster) DestroyInstance(bi BevInstance) {
	bi.(*BevGameMasterInst).BevGameMaster = nil
}

type BevGameMasterInst struct {
	bev
	*BevGameMaster
}

// 行为类型
func (b *BevGameMasterInst) BevType() BevType {
	return gameMaster
}

func (b *BevGameMasterInst) OnInit(ctx Context) bool {
	b.bev.OnInit(ctx)

	msg := &csmsg.GameMasterToS{
		Cmds: make([]*csmsg.GmCmd, len(b.Cmd)),
	}
	for i, v := range b.Cmd {
		msg.Cmds[i] = v.pack()
	}

	b.sendMessage(msg, b.onGameMaster)

	b.player.Infof("request to gamemaster \"%s\"", b.Desc)

	return true
}

func (b *BevGameMasterInst) OnTerminate(ctx Context) {
	b.bev.OnTerminate(ctx)
}

func (b *BevGameMasterInst) onGameMaster(r player, msg *codecs.Message) bool {
	if !net.IsMessageOK(msg) {
		b.player.Errorf("gamemaster \"%s\" failed: %s", b.Desc, net.GetMsgErrcodeStr(msg))
		b.terminate(false)
		return false
	}

	b.player.Infof("gamemaster \"%s\" successfully", b.Desc)
	b.terminate(true)

	return false
}
