package user

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"initialthree/cluster"
	"initialthree/node/common/attr"
	"initialthree/node/common/enumType"
	"initialthree/node/common/inoutput"
	"initialthree/node/common/pipeline"
	"initialthree/node/node_game/module/character"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/NewbieGift"
	"initialthree/node/table/excel/DataTable/Quest"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"

	"initialthree/node/node_game/module"
	attr2 "initialthree/node/node_game/module/attr"
	"initialthree/node/node_game/module/quest"
	TableQuest "initialthree/node/table/excel/ConstTable/Quest"
	"initialthree/node/table/excel/DataTable/AccountInitializeAssets"
	"initialthree/node/table/excel/DataTable/MainDungeon"
	"time"
)

var initStop = errors.New("pipeline step end")

const (
	step_waitCreateRole = 0 // 等待创建账户, 只允许CreateRole协议
	step_initAsset      = 1 // 初始化资源, 出错踢人
	step_startUserLogic = 2 // 启动用户逻辑，出错踢人
	step_ok             = 3 // 所有步骤完成。启动玩家定时器
)

type userLoadPipeline struct {
	user *User
	step int
}

func (this *User) loadStep() int {
	return this.loadPipeline.step
}

func (this *User) DoLoadPipeline() {
	if this.loadPipeline == nil {
		this.loadPipeline = &userLoadPipeline{
			user: this,
			step: step_waitCreateRole,
		}
	}

	if this.loadPipeline.step != step_ok {
		out, err := loadPipeline.RunStartStep(this.loadPipeline, this.loadPipeline.step)
		if err != nil && err != initStop {
			zaplogger.GetSugar().Errorf("%s doLoadPipeline %s", this.GetUserLogName(), err.Error())
			this.kick(false)
			return
		}
		this.loadPipeline = out.(*userLoadPipeline)
	}

	if this.loadPipeline.step == step_ok {
		this.FlushAllToClient()
	}
}

func (this *User) initAsset() error {

	// 任务
	questIds := make([]int32, 0, 32)
	// enumType.QuestType_MainStory:
	questIds = append(questIds, TableQuest.GetMainQuest()...)
	// enumType.QuestType_Daily:
	questIds = append(questIds, TableQuest.GetDailyQuest(int(time.Now().Weekday()))...)
	// enumType.QuestType_Weekly:
	questIds = append(questIds, TableQuest.GetWeekdayQuest()...)
	// enumType.QuestType_Event:

	// enumType.QuestType_DailyReward:
	questIds = append(questIds, TableQuest.GetDailyRewardQuest()...)
	// enumType.QuestType_InstanceReward:
	questIds = append(questIds, MainDungeon.GetInstanceQuest()...)
	// enumType.QuestType_NewbieGift:
	for day := int32(1); day <= 7; day++ {
		def := NewbieGift.GetID(day)
		if def != nil {
			for _, v := range def.QuestIDListArray {
				if v.QuestID != 0 {
					questIds = append(questIds, v.QuestID)
				}
			}
			questIds = append(questIds, def.GroupRewardQuest)
		}
	}

	for _, questId := range questIds {
		if Quest.GetID(questId) == nil {
			return fmt.Errorf("initAsset: InitAsset quest %d is nil", questId)
		}
	}

	// 资源
	assets := AccountInitializeAssets.GetIDMap()
	out := make([]inoutput.ResDesc, 0, len(assets)+1)
	rollAttrFunc := []func(){}

	attrModule := this.GetSubModule(module.Attr).(*attr2.UserAttr)
	for _, v := range assets {
		// 属性相关直接设置值，用ioput会触发相关事件
		if v.AssetTypeEnum == enumType.IOType_UsualAttribute {
			oldVal, err := attrModule.SetAttr(v.AssetID, int64(v.Count), false)
			if err != nil {
				id := v.AssetID
				rollAttrFunc = append(rollAttrFunc, func() {
					attrModule.SetAttr(id, oldVal, false)
				})
			}
		}
	}
	for _, v := range assets {
		if v.AssetTypeEnum != enumType.IOType_UsualAttribute {
			out = append(out, inoutput.ResDesc{Type: int(v.AssetTypeEnum), ID: v.AssetID, Count: v.Count})
		}
	}

	if err := inoutput.DoInputOutput(this, nil, out); err != nil {
		zaplogger.GetSugar().Errorf("%s %s %s", this.GetUserLogName(), "initAsset inoutput", err.Error())
		for _, fn := range rollAttrFunc {
			fn()
		}
		return err
	}

	// 资源标记
	attrModule.SetAttr(attr.IsInitAsset, int64(1), false)

	// 初始任务
	questModule := this.GetSubModule(module.Quest).(*quest.UserQuest)
	for _, questId := range questIds {
		questModule.AddQuest(questId)
	}

	// 默认编队
	defaultTeam := Global.Get().DefaultCharacterTeamArray
	if len(defaultTeam) == 3 {
		charaModule := this.GetSubModule(module.Character).(*character.UserCharacter)
		prefab := &message.CharacterTeamPrefab{
			Name:          proto.String("默认编队"),
			CharacterList: make([]int32, 0, 3),
		}
		for _, v := range defaultTeam {
			if v.ID != 0 {
				if charaModule.GetCharacter(v.ID) == nil {
					return fmt.Errorf("default character team %d is not find", v.ID)
				}
			}
			prefab.CharacterList = append(prefab.CharacterList, v.ID)
		}
		charaModule.GroupDefSet(2, prefab)
	}

	return nil
}

func (this *User) startUserLogic() error {

	// 模块相关功能启动
	for _, m := range this.modules {
		allInit, ok := m.(module.AfterModuleInitAll)
		if ok {
			if err := allInit.AfterInitAll(); err != nil {
				zaplogger.GetSugar().Errorf("user %s module %s AfterInitAll %s", this.GetUserID(), m.ModuleType().String(), err)
				return err
			}
		}
	}

	cluster.RegisterTimerOnce(time.Second, this.Tick, nil)
	return nil
}

var loadPipeline *pipeline.Pipeline

func waitCreateRole(in interface{}) (interface{}, error) {
	out := in.(*userLoadPipeline)
	if out.user.GetID() == 0 || out.user.GetName() == "" {
		return out, initStop
	}
	zaplogger.GetSugar().Infof("%s %s", out.user.userID, "loadPipeline step waitCreateRole ok")
	out.step = step_initAsset
	return out, nil
}

func initAsset(in interface{}) (interface{}, error) {
	out := in.(*userLoadPipeline)
	if out.user.GetAttr(attr.IsInitAsset) == 1 {
		out.step = step_startUserLogic
		return out, nil
	}

	// 初始化资源
	if err := out.user.initAsset(); err != nil {
		zaplogger.GetSugar().Errorf("%v %s %s", out.user.GetUserID(), "initAsset failed", err.Error())
		return out, err
	}
	zaplogger.GetSugar().Infof("%v %s", out.user.userID, "loadPipeline step initAsset ok")
	out.step = step_startUserLogic
	return out, nil
}

func startUserLogic(in interface{}) (interface{}, error) {
	out := in.(*userLoadPipeline)
	if err := out.user.startUserLogic(); err != nil {
		zaplogger.GetSugar().Errorf("%s %s %s", out.user.GetUserID(), "startUserLogic failed", err.Error())
		return out, err
	}
	zaplogger.GetSugar().Infof("%s %s", out.user.userID, "loadPipeline step startUserLogic ok -> step_ok")
	out.step = step_ok
	return out, nil
}

func init() {
	loadPipeline = pipeline.NewPipeline()
	loadPipeline.AddStep(waitCreateRole, initAsset, startUserLogic)
}
