
package attr

type Attr struct {
	ID  int32
	Val int64
}

type AttrInfo struct {
	Min int64
	Max int64
}

var nameToIdx map[string]int32
var idxToName map[int32]string
var infoMap   map[int32]*AttrInfo

const(
	Level = 1 //等级
	CurrentExp = 2 //经验值
	CurrentFatigue = 3 //阻灵值
	CurrentTitle = 4 //当前称号ID
	Gold = 5 //法芙娜铸币
	Diamond = 6 //灵质黑盒
	DailyActiveness = 7 //日活跃度
	FatigueBuyCount = 8 //阻灵值购买次数
	IsInitAsset = 9 //是否初始化资源
	DailyOnLine = 10 //日在线时长
	EquipCap = 11 //装备背包容量
	WeaponCap = 12 //武器背包容量
	AccumulateLogin = 13 //累积登陆天数
	YuruCharacterID = 14 //看板娘DisplayID
	PassedPrologue = 15 //通关序章
	ScarsingrainSeal = 16 //战痕印章
	GoldBuyCount = 17 //法芙娜铸币购买次数
	FreeRenameTimes = 18 //免费改名次数
	SecretCoin = 19 //深渊奖章
	TrueDiamond = 20 //灵质立方
	NewbieGiftStartTime = 21 //七日任务开启时间
	NewbieGiftEndTime = 22 //七日任务结束时间
	WeeklyLogin = 23 //每周登陆次数
	WeeklyActiveness = 24 //每周活跃度
	AttrMax = 24
)

func init() {	
	nameToIdx = map[string]int32{}
	idxToName = map[int32]string{}
	infoMap   = map[int32]*AttrInfo{}

	nameToIdx["Level"] = 1
	idxToName[1] = "Level"
	nameToIdx["CurrentExp"] = 2
	idxToName[2] = "CurrentExp"
	nameToIdx["CurrentFatigue"] = 3
	idxToName[3] = "CurrentFatigue"
	nameToIdx["CurrentTitle"] = 4
	idxToName[4] = "CurrentTitle"
	nameToIdx["Gold"] = 5
	idxToName[5] = "Gold"
	nameToIdx["Diamond"] = 6
	idxToName[6] = "Diamond"
	nameToIdx["DailyActiveness"] = 7
	idxToName[7] = "DailyActiveness"
	nameToIdx["FatigueBuyCount"] = 8
	idxToName[8] = "FatigueBuyCount"
	nameToIdx["IsInitAsset"] = 9
	idxToName[9] = "IsInitAsset"
	nameToIdx["DailyOnLine"] = 10
	idxToName[10] = "DailyOnLine"
	nameToIdx["EquipCap"] = 11
	idxToName[11] = "EquipCap"
	nameToIdx["WeaponCap"] = 12
	idxToName[12] = "WeaponCap"
	nameToIdx["AccumulateLogin"] = 13
	idxToName[13] = "AccumulateLogin"
	nameToIdx["YuruCharacterID"] = 14
	idxToName[14] = "YuruCharacterID"
	nameToIdx["PassedPrologue"] = 15
	idxToName[15] = "PassedPrologue"
	nameToIdx["ScarsingrainSeal"] = 16
	idxToName[16] = "ScarsingrainSeal"
	nameToIdx["GoldBuyCount"] = 17
	idxToName[17] = "GoldBuyCount"
	nameToIdx["FreeRenameTimes"] = 18
	idxToName[18] = "FreeRenameTimes"
	nameToIdx["SecretCoin"] = 19
	idxToName[19] = "SecretCoin"
	nameToIdx["TrueDiamond"] = 20
	idxToName[20] = "TrueDiamond"
	nameToIdx["NewbieGiftStartTime"] = 21
	idxToName[21] = "NewbieGiftStartTime"
	nameToIdx["NewbieGiftEndTime"] = 22
	idxToName[22] = "NewbieGiftEndTime"
	nameToIdx["WeeklyLogin"] = 23
	idxToName[23] = "WeeklyLogin"
	nameToIdx["WeeklyActiveness"] = 24
	idxToName[24] = "WeeklyActiveness"

	infoMap[1] = &AttrInfo{
		Min: 0,
		Max: 120,
	}
	infoMap[2] = &AttrInfo{
		Min: 0,
		Max: 9999999999,
	}
	infoMap[3] = &AttrInfo{
		Min: 0,
		Max: 5000,
	}
	infoMap[4] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[5] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[6] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[7] = &AttrInfo{
		Min: 0,
		Max: 9999,
	}
	infoMap[8] = &AttrInfo{
		Min: 0,
		Max: 10,
	}
	infoMap[9] = &AttrInfo{
		Min: 0,
		Max: 1,
	}
	infoMap[10] = &AttrInfo{
		Min: 0,
		Max: 86400,
	}
	infoMap[11] = &AttrInfo{
		Min: 0,
		Max: 1000,
	}
	infoMap[12] = &AttrInfo{
		Min: 0,
		Max: 1000,
	}
	infoMap[13] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[14] = &AttrInfo{
		Min: 1,
		Max: 99999999999,
	}
	infoMap[15] = &AttrInfo{
		Min: 0,
		Max: 1,
	}
	infoMap[16] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[17] = &AttrInfo{
		Min: 0,
		Max: 10,
	}
	infoMap[18] = &AttrInfo{
		Min: 0,
		Max: 10,
	}
	infoMap[19] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[20] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[21] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[22] = &AttrInfo{
		Min: 0,
		Max: 99999999999,
	}
	infoMap[23] = &AttrInfo{
		Min: 0,
		Max: 7,
	}
	infoMap[24] = &AttrInfo{
		Min: 0,
		Max: 99999,
	}
}

func GetIdByName(name string) int32 {
	return nameToIdx[name]
}

func GetNameById(id int32) string {
	return idxToName[id]
}

func GetAttrInfo(idx int32) *AttrInfo {
	return infoMap[idx]
}
