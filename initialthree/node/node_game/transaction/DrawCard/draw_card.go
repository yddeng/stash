package DrawCard

import (
	"initialthree/node/common/inoutput"
	"initialthree/node/node_game/event"
	"initialthree/node/node_game/module/character"
	"initialthree/node/node_game/module/drawCard"
	"initialthree/node/node_game/module/weapon"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/DrawCardsPool"
	"initialthree/node/table/excel/DataTable/DropPool"
	"initialthree/zaplogger"
	"time"

	"github.com/gogo/protobuf/proto"

	codecs "initialthree/codec/cs"
	"initialthree/node/common/droppool"
	"initialthree/node/common/enumType"
	"initialthree/node/common/transaction"

	"initialthree/node/node_game/module"
	"initialthree/node/node_game/user"
	"initialthree/node/table/excel/DataTable/DrawCardsLib"
	"initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/protocol/cmdEnum"
	cs_message "initialthree/protocol/cs/message"
)

type transactionDrawCard struct {
	transaction.TransactionBase
	dc       *drawCard.DrawCard
	user     *user.User
	req      *codecs.Message
	resp     *cs_message.DrawCardDrawToC
	errCode  cs_message.ErrCode
	constDef *Global.Global
}

func (t *transactionDrawCard) GetModuleName() string {
	return "DrawCard"
}

func (t *transactionDrawCard) Begin() {
	defer func() { t.EndTrans(t.resp, t.errCode) }()
	msg := t.req.GetData().(*cs_message.DrawCardDrawToS)
	zaplogger.GetSugar().Infof("%s Call DrawCardDrawToS %v", t.user.GetUserID(), msg)
	t.resp = &cs_message.DrawCardDrawToC{LibID: proto.Int32(msg.GetLibID())}
	// 抽卡
	// 1. 卡池是否存在 & 卡池是否开放 & 消耗资源是否充足
	// 2. 增加保底值
	// 3. 保底库 ID 是否合规
	// 4. 保底值是否足够
	// 4.1 足够则抽取一次保底库
	// 4.2 不足则抽取一次常规库
	// 5. 如果是十连抽
	// 5.1 抽取 S 级卡牌超过 4 个，后续如果再次抽取 S 级卡牌，则丢弃重新抽取
	// 5.2 如果前 9 次抽取无 A 级及以上卡牌，则从 A 级保底库抽取一次
	// 6. 检查抽取的卡牌是否已存在，存在则转换为对应物品
	// 7. 添加抽卡结果到玩家&同步抽卡结果

	t.dc = t.user.GetSubModule(module.DrawCard).(*drawCard.DrawCard)
	poolIdx := t.dc.GetPoolIndex(msg.GetLibID())

	if t.constDef = Global.Get(); t.constDef == nil {
		t.errCode = cs_message.ErrCode_Config_NotExist
		return
	}

	if msg.GetDrawCount() != 1 && msg.GetDrawCount() != 10 {
		t.errCode = cs_message.ErrCode_DrawCard_Count // 请求抽卡次数错误
	} else if drawCardLib := DrawCardsLib.GetID(msg.GetLibID()); drawCardLib == nil {
		t.errCode = cs_message.ErrCode_Config_NotExist // 卡池不存在
	} else if !t.isDrawCardLibOpen(drawCardLib) {
		t.errCode = cs_message.ErrCode_DrawCard_NotOpen // 卡池未开放
	} else if !t.hasEnoughDrawCardConsumeToken(drawCardLib) {
		t.errCode = cs_message.ErrCode_DrawCard_NotEnoughCostAssets // 所需资源不足
	} else if drawCardTimes := t.dc.GetDailyTimes(drawCardLib.GuaranteeID); drawCardTimes+msg.GetDrawCount() > Global.Get().DrawCardDailyLimit {
		t.errCode = cs_message.ErrCode_DrawCard_TimesLimit // 每日抽卡次数上限
	} else if poolId, ok := t.getDrawCardsPoolID(drawCardLib, poolIdx); !ok { // 保底库 ID 错误，卡池中未包含此保底库
		t.errCode = cs_message.ErrCode_DrawCard_GuaranteeLibID
	} else {

		awardInfoList, guaranteeID, guaranteeCount := t.drawCard(drawCardLib, poolId, msg.GetDrawCount())

		useRes := []inoutput.ResDesc{{Type: enumType.IOType_Item, ID: drawCardLib.ConsumeTokenType, Count: drawCardLib.ConsumeTokenCount * msg.GetDrawCount()}}

		characterStatus, weaponStatus, outRes := t.drawCardBeforeStatus(awardInfoList)

		// 消耗对应的道具，道具模块会自动同步，不需要额外同步
		if t.errCode = t.user.DoInputOutput(useRes, outRes, true); t.errCode != cs_message.ErrCode_OK {
			zaplogger.GetSugar().Errorf("AddDrawCardAwardInfoList failed %s", t.errCode.String())
			return
		}

		t.packDrawCardAward(characterStatus, weaponStatus, awardInfoList)

		t.dc.AddHistory(msg.GetLibID(), t.resp.AwardList)
		t.dc.AddDailyTimes(drawCardLib.GuaranteeID, msg.GetDrawCount())

		t.user.EmitEvent(event.EventDrawCard, msg.GetLibID(), msg.GetDrawCount())

		t.resp.GuaranteeID, t.resp.GuaranteeCount = proto.Int32(guaranteeID), proto.Int32(guaranteeCount)
		t.dc.SetGuarantee(guaranteeID, guaranteeCount)
		zaplogger.GetSugar().Infof("%s %s", t.user.GetUserID(), "DrawCardDraw ok")
	}
}

func (t *transactionDrawCard) drawCardBeforeStatus(awardInfoList []*droppool.AwardInfo) (map[int32]bool, map[int32]bool, []inoutput.ResDesc) {
	characterStatus := map[int32]bool{}
	weaponStatus := map[int32]bool{}
	outRes := make([]inoutput.ResDesc, 0, len(awardInfoList))

	userCharacter := t.user.GetSubModule(module.Character).(*character.UserCharacter)
	userWeapon := t.user.GetSubModule(module.Weapon).(*weapon.UserWeapon)
	for _, v := range awardInfoList {
		outRes = append(outRes, v.ToResDesc())
		switch v.Type {
		case enumType.DropType_Character:
			if c := userCharacter.GetCharacter(v.ID); c != nil {
				// 角色存在
				characterStatus[v.ID] = true
			} else {
				characterStatus[v.ID] = false
			}

		case enumType.DropType_Weapon:
			if userWeapon.ConfigIDIsExist(v.ID) {
				// 存在同名武器
				weaponStatus[v.ID] = true
			} else {
				weaponStatus[v.ID] = false
			}
		}
	}

	return characterStatus, weaponStatus, outRes
}

func (t *transactionDrawCard) isDrawCardLibOpen(drawCardLib *DrawCardsLib.DrawCardsLib) bool {
	if drawCardLib.LibOpenType == 1 { //  永久开放
		return true
	} else if drawCardLib.LibOpenType == 2 { //  指定日期开放
		if openTime, err := time.ParseInLocation("2006-01-02T15:04:05", drawCardLib.LibOpenTime, time.Local); err != nil {
			zaplogger.GetSugar().Errorf("%s DrawCardDraw %d ParseTime %s err %s", t.user.GetUserLogName(), drawCardLib.ID, drawCardLib.LibOpenTime, err)
			return false
		} else if endTime, err := time.ParseInLocation("2006-01-02T15:04:05", drawCardLib.LibEndTime, time.Local); err != nil {
			zaplogger.GetSugar().Errorf("%s DrawCardDraw %d ParseTime %s err %s", t.user.GetUserLogName(), drawCardLib.ID, drawCardLib.LibEndTime, err)
			return false
		} else {
			now := time.Now()
			return now.After(openTime) && now.Before(endTime)
		}
	}
	return false
}

func (t *transactionDrawCard) packDrawCardAward(characterStatus, weaponStatus map[int32]bool, awardInfoList []*droppool.AwardInfo) {
	t.resp.AwardList = make([]*cs_message.DrawCardAward, len(awardInfoList))

	timestamp := time.Now().Unix()
	userCharacter := t.user.GetSubModule(module.Character).(*character.UserCharacter)
	for i, v := range awardInfoList {
		award := &cs_message.DrawCardAward{
			AwardInfo: v.ToMessage(),
			Timestamp: proto.Int64(timestamp),
		}

		switch v.Type {
		case enumType.DropType_Character:
			exist := characterStatus[v.ID]
			if !exist {
				// 一次出现多次同ID角色，第一次将标记替换
				characterStatus[v.ID] = true
			} else {
				if c := userCharacter.GetCharacter(v.ID); c == nil {
					// 不存在
				} else if def := PlayerCharacter.GetID(c.CharacterID); c.HitTimes < def.DrawCardTimes {
					// 角色已存在，但不满命座要转成碎片
					award.State = proto.Int32(1)
				} else {
					// 角色满命座满
					award.State = proto.Int32(2)
				}
			}

		case enumType.DropType_Weapon:
			exist := weaponStatus[v.ID]
			if !exist {
				weaponStatus[v.ID] = true
			} else {
				award.State = proto.Int32(1)
			}
		}
		t.resp.AwardList[i] = award
	}
	return
}

func (t *transactionDrawCard) hasEnoughDrawCardConsumeToken(drawCardLib *DrawCardsLib.DrawCardsLib) bool {
	return t.user.GetItemCountByTID(drawCardLib.ConsumeTokenType) >= drawCardLib.ConsumeTokenCount
}

func (t *transactionDrawCard) getDrawCardsPoolID(drawCardLib *DrawCardsLib.DrawCardsLib, poolIdx int32) (int32, bool) {
	if int32(len(drawCardLib.DrawCardsPoolArray)) <= poolIdx {
		return 0, false
	}
	v := drawCardLib.DrawCardsPoolArray[poolIdx]
	if v.PoolID == 0 {
		return 0, false
	}
	return v.PoolID, true
}

// 默认配置文件正确，确保抽取的卡牌配置有效 & 十连抽抽取 S 级卡牌的几率不会太高导致过量 for 循环
func (t *transactionDrawCard) drawCard(drawCardLib *DrawCardsLib.DrawCardsLib, poolID, drawCount int32) ([]*droppool.AwardInfo, int32, int32) {
	guaranteeCount, sCardCount, aCardCount := t.dc.GetGuaranteeCount(drawCardLib.GuaranteeID), 0, 0
	guaranteeCount += drawCardLib.GuaranteeIncrease * drawCount

	tenGuarantee := drawCount == 10
	awardInfoList := make([]*droppool.AwardInfo, 0)
	if guaranteeCount >= drawCardLib.GuaranteeThreshold {
		// 抽取一次保底库
		awardInfoList, sCardCount, aCardCount = t.drawCardOnce(DrawCardsPool.GetGuaranteePool(poolID), awardInfoList, sCardCount, aCardCount)
		guaranteeCount -= drawCardLib.GuaranteeThreshold
		drawCount--
	}

	for drawCount > 0 {
		oldSCount := sCardCount
		if tenGuarantee && drawCount == 1 && sCardCount+aCardCount == 0 {
			// 10连抽 最后一次抽取时，没有 A 级以上卡牌，使用 A 级保底卡库
			awardInfoList, sCardCount, aCardCount = t.drawCardOnce(DrawCardsPool.GetTenGuaranteePool(poolID), awardInfoList, sCardCount, aCardCount)
		} else {
			awardInfoList, sCardCount, aCardCount = t.drawCardOnce(DrawCardsPool.GetDrawCardsPool(poolID), awardInfoList, sCardCount, aCardCount)
		}
		// 抽到了S卡，保底值清0
		drawCount--
		if sCardCount > oldSCount {
			guaranteeCount = drawCount
		}
	}
	return awardInfoList, drawCardLib.GuaranteeID, guaranteeCount
}

func (t *transactionDrawCard) drawCardOnce(pool *DropPool.DropPool, awardInfoList []*droppool.AwardInfo, sCardCount, aCardCount int) ([]*droppool.AwardInfo, int, int) {
	award := droppool.DropWithPool(pool)
	if len(award.Infos) != 1 {
		zaplogger.GetSugar().Errorf("%s DrawCardDraw award info failed, length %d", t.user.GetUserLogName(), len(award.Infos))
		return awardInfoList, sCardCount, aCardCount
	}
	sCount, aCount := 0, 0
	for _, awardInfo := range award.Infos { // 确保一次抽取只会抽取一张卡牌，否则可能出错
		if awardInfo.Type == enumType.DropType_Character {
			// 仅相关角色
			if def := PlayerCharacter.GetID(awardInfo.ID); def != nil {
				if def.RarityEnum >= t.constDef.DrawCardGuaranteeRarityEnum {
					sCount++
				}
				if def.RarityEnum >= t.constDef.DrawCardTenthGuaranteeRarityEnum {
					aCount++
				}
			}
		}
	}

	if sCardCount+sCount > 4 {
		// 如果已抽取 4 张及以上的 S 级卡牌，再次抽取 S 级无效，重新抽取
		return t.drawCardOnce(DrawCardsPool.GetDrawCardsPool(pool.ID), awardInfoList, sCardCount, aCardCount)
	}

	sCardCount += sCount
	aCardCount += aCount
	for _, awardInfo := range award.Infos {
		awardInfoList = append(awardInfoList, awardInfo)
	}

	return awardInfoList, sCardCount, aCardCount
}

func init() {
	user.RegisterTransFunc(cmdEnum.CS_DrawCardDraw, func(user *user.User, msg *codecs.Message) transaction.Transaction {
		return &transactionDrawCard{
			user: user,
			req:  msg,
		}
	})
}
