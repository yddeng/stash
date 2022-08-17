package behavior

import (
	"encoding/xml"
	codecs "initialthree/codec/cs"
	csmsg "initialthree/protocol/cs/message"
	"initialthree/robot/internal"
	"initialthree/robot/net"
	"initialthree/robot/robot/module"

	. "github.com/GodYY/bevtree"
	"github.com/golang/protobuf/proto"
)

type starAwardType int8

const (
	StarAward_Unknown = starAwardType(iota)
	StarAward_MainDungeon
)

var starAwardTypeNames = [...]string{
	StarAward_MainDungeon: "maindungeon",
}

func (t starAwardType) String() string { return starAwardTypeNames[t] }

var starAwardName2Types = map[string]starAwardType{}

func init() {
	for i, v := range starAwardTypeNames {
		starAwardName2Types[v] = starAwardType(i)
	}
}

func (t starAwardType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement(starAwardTypeNames[t], start)
}

func (t *starAwardType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var name string
	if err := d.DecodeElement(&name, &start); err != nil {
		return err
	}

	*t = starAwardName2Types[name]
	return nil
}

const claimStaraward = BevType("claimstaraward")

func init() {
	regBevType(claimStaraward, func() Bev { return new(BevClaimStaraward) })
}

type BevClaimStaraward struct {
	AwardType starAwardType `xml:"awardtype"`
}

func (BevClaimStaraward) BevType() BevType {
	return claimStaraward
}

func (b *BevClaimStaraward) CreateInstance() BevInstance {
	return &BevClaimStarAwardInst{
		BevClaimStaraward: b,
	}
}

func (b *BevClaimStaraward) DestroyInstance(bi BevInstance) {
	bi.(*BevClaimStarAwardInst).BevClaimStaraward = nil
}

type BevClaimStarAwardInst struct {
	bev
	*BevClaimStaraward
	msg proto.Message
}

// 行为类型
func (b *BevClaimStarAwardInst) BevType() BevType {
	return claimStaraward
}

func (b *BevClaimStarAwardInst) OnInit(ctx Context) bool {
	b.bev.OnInit(ctx)

	switch b.AwardType {
	case StarAward_MainDungeon:
		return b.performMainDungeon()

	default:
		b.player.Panicf("invalid staraward type: %s", b.AwardType)
		return false
	}
}

func (b *BevClaimStarAwardInst) OnTerminate(ctx Context) {
	b.msg = nil
	b.bev.OnTerminate(ctx)
}

func (b *BevClaimStarAwardInst) performMainDungeon() bool {
	moduleMainDungeon := b.player.GetModule(module.Module_MainDungeons).(*module.ModuleMainDungeons)
	chapterId, awardNo, ok := moduleMainDungeon.GetNextStarAwardNo()
	if chapterId == internal.InvalidID || awardNo == internal.InvalidID {
		b.player.Debugf("no next maindungeon staraward")
		return false
	}

	if !ok {
		b.player.Debugf("maindungeon staraward %d-%d could not be claimed yet", chapterId, awardNo)
		return false
	}

	b.msg = &csmsg.MainDungeonsGetChapterStarAwardToS{
		ChapterID: proto.Int32(chapterId),
		AwardNo:   proto.Int32(awardNo),
	}

	b.sendMessage(b.msg, b.onGetStarAward)

	b.player.Debugf("request to claim maindungeon staraward %d-%d", chapterId, awardNo)

	return true
}

func (b *BevClaimStarAwardInst) onGetStarAward(r player, msg *codecs.Message) bool {
	switch b.AwardType {
	case StarAward_MainDungeon:
		req := b.msg.(*csmsg.MainDungeonsGetChapterStarAwardToS)

		if !net.IsMessageOK(msg) {
			b.player.Errorf("claim maindungeon staraward %d-%d failed: %s", req.GetChapterID(), req.GetAwardNo(), net.GetErrCodeStr(msg.GetErrCode()))
			b.terminate(false)
			return false
		}

		b.player.Infof("claim maindungeon staraward %d-%d successfully", req.GetChapterID(), req.GetAwardNo())
		b.terminate(true)
		return false

	default:
		b.player.Panicf("invalid staraward type: %s", b.AwardType)
	}

	return false
}
